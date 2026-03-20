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
