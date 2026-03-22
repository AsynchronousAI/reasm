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
    stdout_cache = stdout_cache .. str .. "\n"
    flush_stdout()
end

functions["fputs"] = function()
    local str_ptr, stream = get_args()
    local str = read_string(str_ptr)
    stdout_cache = stdout_cache .. str
    flush_stdout()
    push_args(0)
end

functions["fflush"] = function()
    flush_stdout()
    push_args(0)
end

functions["fopen"] = function()
    local path_ptr, mode_ptr = get_args()
    -- local path = read_string(path_ptr)
    -- local mode = read_string(mode_ptr)
    -- Just return 0 (NULL) for all file open requests for now
    push_args(0)
end

functions["fclose"] = function()
    push_args(0)
end

functions["fread"] = function()
    push_args(0)
end

functions["fseek"] = function()
    push_args(-1)
end

functions["ftell"] = function()
    push_args(-1)
end

functions["feof"] = function()
    push_args(1)
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

functions["fprintf"] = function()
    local args = {get_args()}
    -- a0 is FILE*, a1 is fmt, a2... are args
    local fmt_ptr = args[2]
    table.remove(args, 1) -- remove stream
    table.remove(args, 1) -- remove fmt

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

functions["vfprintf"] = function()
    local args = {get_args()}
    local fmt_ptr = args[2]
    local va_list = args[3]
    
    -- va_list is a pointer to memory where arguments are stored
    local format_args = {}
    local current_va = va_list
    -- Read up to 10 arguments as a heuristic, assuming 32-bit words
    for i = 1, 10 do
        table.insert(format_args, buffer.readi32(memory, current_va))
        current_va += 4
    end

    local formatted = format_string(fmt_ptr, format_args)
    for i = 1, #formatted do
        local ch = formatted:sub(i, i)
        if ch == "\n" then
            flush_stdout()
        else
            stdout_cache = stdout_cache .. ch
        end
    end
end

functions["vsnprintf"] = function()
    local args = {get_args()}
    local str = args[1]
    local size = args[2]
    local fmt_ptr = args[3]
    local va_list = args[4]
    
    local format_args = {}
    local current_va = va_list
    for i = 1, 10 do
        table.insert(format_args, buffer.readi32(memory, current_va))
        current_va += 4
    end

    if size > 0 then
        local formatted = format_string(fmt_ptr, format_args)
        local write_len = math.min(#formatted, size - 1)
        for i = 1, write_len do
            buffer.writeu8(memory, str + i - 1, string.byte(formatted, i))
        end
        buffer.writeu8(memory, str + write_len, 0)
        push_args(write_len)
    else
        push_args(0)
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
