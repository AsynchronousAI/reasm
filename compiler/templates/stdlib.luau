--- <stdlib.h> helpers

local allocations: {[number]: number} = {}
local rand_state: number = 1
local env: {[string]: string} = {}
local atexit_callbacks: {number} = {}

local function record_allocation(ptr: number, size: number)
    allocations[ptr] = size
    return ptr
end

local function allocate(size: number): number
    if size <= 0 then
        size = 1
    end
    mallocDepth += size
    local ptr = buffer.len(memory) - mallocDepth
    return record_allocation(ptr, size)
end

local function write_pointer(destination: number, value: number)
    buffer.writeu32(memory, destination, bit32.band(value, 0xFFFFFFFF))
end

local function read_pointer(source: number)
    return buffer.readu32(memory, source)
end

local function copy_c_string(dest: number, value: string)
    for i = 1, #value do
        buffer.writeu8(memory, dest + i - 1, string.byte(value, i))
    end
    buffer.writeu8(memory, dest + #value, 0)
end

local function align_up(value: number, alignment: number)
    local mask = alignment - 1
    return bit32.band(value + mask, bit32.bnot(mask))
end

local function native_malloc(size: number)
    return allocate(size)
end

malloc = native_malloc

functions["malloc"] = function()
    local size = get_args()
    push_args(native_malloc(size))
end

functions["calloc"] = function()
    local nmemb, size = get_args()
    local total = nmemb * size
    local ptr = native_malloc(total)
    buffer.fill(memory, ptr, 0, total)
    push_args(ptr)
end

functions["realloc"] = function()
    local ptr, size = get_args()
    if ptr == 0 then
        push_args(native_malloc(size))
        return
    end
    local old_size = allocations[ptr] or 0
    local new_ptr = native_malloc(size)
    local copy_size = math.min(old_size, size)
    buffer.copy(memory, new_ptr, memory, ptr, copy_size)
    allocations[ptr] = nil
    push_args(new_ptr)
end

functions["free"] = function()
    local ptr = get_args()
    allocations[ptr] = nil
end

functions["aligned_alloc"] = function()
    local alignment, size = get_args()
    if alignment == 0 or bit32.band(alignment, alignment - 1) ~= 0 then
        errno = 22
        push_args(0)
        return
    end
    local raw = native_malloc(size + alignment)
    allocations[raw] = nil
    local aligned = align_up(raw, alignment)
    record_allocation(aligned, size)
    push_args(aligned)
end

functions["posix_memalign"] = function()
    local memptr, alignment, size = get_args()
    if bit32.band(alignment, alignment - 1) ~= 0 then
        errno = 22
        push_args(22)
        return
    end
    local raw = native_malloc(size + alignment)
    allocations[raw] = nil
    local aligned = align_up(raw, alignment)
    record_allocation(aligned, size)
    write_pointer(memptr, aligned)
    errno = 0
    push_args(0)
end

functions["rand"] = function()
    rand_state = bit32.band(rand_state * 1103515245 + 12345, 0x7fffffff)
    push_args(rand_state)
end

functions["srand"] = function()
    local seed = get_args()
    rand_state = seed
end

functions["abs"] = function()
    local value = get_args()
    push_args(math.abs(value))
end

functions["labs"] = function()
    local value = get_args()
    push_args(math.abs(value))
end

local function record_division_result(a, b)
    if b == 0 then
        errno = 22
        push_args(0, 0)
        return true
    end
    local quot = math.floor(a / b)
    local rem = a - quot * b
    push_args(quot, rem)
    return false
end

functions["div"] = function()
    local a, b = get_args()
    if not record_division_result(a, b) then
        errno = 0
    end
end

functions["ldiv"] = function()
    local a, b = get_args()
    if not record_division_result(a, b) then
        errno = 0
    end
end

functions["atexit"] = function()
    local func = get_args()
    table.insert(atexit_callbacks, func)
end

functions["abort"] = function()
    errno = 6
    error("abort called")
end

local function exit_with_code(return_code: number)
    local code = return_code or 0
    flush_stdout()
    error("exit(" .. tostring(code) .. ")")
end

functions["exit"] = function()
    local code = get_args()
    exit_with_code(code)
end

functions["_Exit"] = functions["exit"]

functions["getenv"] = function()
    local name = read_string(get_args())
    local value = env[name]
    if not value then
        push_args(0)
        return
    end
    local ptr = native_malloc(#value + 1)
    copy_c_string(ptr, value)
    push_args(ptr)
end

functions["setenv"] = function()
    local name, value, overwrite = get_args()
    local key = read_string(name)
    local val = read_string(value)
    if env[key] and overwrite == 0 then
        errno = 17
        push_args(-1)
        return
    end
    env[key] = val
    errno = 0
    push_args(0)
end

functions["putenv"] = function()
    local pair = read_string(get_args())
    local separator = pair:find("=")
    if not separator then
        errno = 22
        push_args(-1)
        return
    end
    local key = pair:sub(1, separator - 1)
    local value = pair:sub(separator + 1)
    env[key] = value
    errno = 0
    push_args(0)
end
