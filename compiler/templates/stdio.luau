--- <stdio.h> helpers

functions["putchar"] = function()
    local c = get_args()
    local char = string.char(c)
    if char == "\n" then
        flush_stdout()
    else
        stdout_cache = stdout_cache .. char
    end
end

functions["puts"] = function()
    local fmt = get_args()
    local str = read_string(fmt)
    stdout_cache = stdout_cache .. str
    flush_stdout()
end

functions["printf"] = function()
    local args = {get_args()}
    local fmt_ptr = args[1]
    table.remove(args, 1)

    local formatted = format_string(fmt_ptr, args)
    for i = 1, #formatted do
        local ch = formatted:sub(i, i)
        if ch == "\n" then
            flush_stdout()
        else
            stdout_cache = stdout_cache .. ch
        end
    end
end

local strerror_map = {
    [0] = "No error",
    [1] = "Operation not permitted",
    [2] = "No such file or directory",
    [6] = "No such device or address",
    [12] = "Out of memory",
    [17] = "File exists",
    [22] = "Invalid argument",
    [34] = "Result too large",
}

local function strerror_text(code)
    return strerror_map[code] or ("Unknown error " .. tostring(code))
end

local function write_string_to_memory(source)
    local ptr = malloc(#source + 1)
    for i = 1, #source do
        buffer.writeu8(memory, ptr + i - 1, string.byte(source, i))
    end
    buffer.writeu8(memory, ptr + #source, 0)
    return ptr
end

functions["strerror"] = function()
    local errnum = get_args()
    if errnum == 0 then
        errnum = errno
    end
    push_args(write_string_to_memory(strerror_text(errnum)))
end

functions["perror"] = function()
    local message_ptr = get_args()
    local prefix = ""
    if message_ptr ~= 0 then
        prefix = read_string(message_ptr) .. ": "
    end
    local formatted = prefix .. strerror_text(errno)
    stdout_cache = stdout_cache .. formatted .. "\n"
    flush_stdout()
end
