--- <string.h> helpers

local function read_byte(ptr, offset)
    return buffer.readu8(memory, ptr + offset)
end

local function write_byte(ptr, offset, value)
    buffer.writeu8(memory, ptr + offset, bit32.band(value, 0xFF))
end

local function c_strlen(ptr)
    if ptr == 0 then
        return 0
    end
    local len = 0
    while read_byte(ptr, len) ~= 0 do
        len += 1
    end
    return len
end

local function copy_string(dest, src, max_length, pad)
    pad = pad or false
    local i = 0
    while i < max_length do
        local byte = read_byte(src, i)
        write_byte(dest, i, byte)
        i += 1
        if byte == 0 then
            break
        end
    end
    if pad and i < max_length then
        while i < max_length do
            write_byte(dest, i, 0)
            i += 1
        end
    end
    return dest
end

local function is_delimiter(byte, delim)
    if delim == "" then
        return false
    end
    local char = string.char(byte)
    return delim:find(char, 1, true) ~= nil
end

local strtok_state = {
    next = 0,
    delim = "",
}

functions["strlen"] = function()
    local ptr = get_args()
    push_args(c_strlen(ptr))
end

functions["strcpy"] = function()
    local dest, src = get_args()
    copy_string(dest, src, c_strlen(src) + 1)
    push_args(dest)
end

functions["strncpy"] = function()
    local dest, src, count = get_args()
    copy_string(dest, src, count, true)
    push_args(dest)
end

functions["strcat"] = function()
    local dest, src = get_args()
    local dest_len = c_strlen(dest)
    copy_string(dest + dest_len, src, c_strlen(src) + 1)
    push_args(dest)
end

functions["strncat"] = function()
    local dest, src, count = get_args()
    local dest_len = c_strlen(dest)
    local i = 0
    while i < count do
        local byte = read_byte(src, i)
        write_byte(dest, dest_len + i, byte)
        if byte == 0 then
            break
        end
        i += 1
    end
    if i == count then
        write_byte(dest, dest_len + i, 0)
    end
    push_args(dest)
end

local function compare_bytes(a, b, limit, case_insensitive)
    local i = 0
    while i < limit do
        local byteA = read_byte(a, i)
        local byteB = read_byte(b, i)
        if case_insensitive then
            byteA = string.byte(string.lower(string.char(byteA)))
            byteB = string.byte(string.lower(string.char(byteB)))
        end
        if byteA ~= byteB then
            return byteA - byteB
        end
        if byteA == 0 then
            return 0
        end
        i += 1
    end
    return 0
end

functions["strcmp"] = function()
    local a, b = get_args()
    push_args(compare_bytes(a, b, math.max(c_strlen(a), c_strlen(b)) + 1, false))
end

functions["strncmp"] = function()
    local a, b, count = get_args()
    push_args(compare_bytes(a, b, count, false))
end

functions["strcasecmp"] = function()
    local a, b = get_args()
    local limit = math.max(c_strlen(a), c_strlen(b)) + 1
    push_args(compare_bytes(a, b, limit, true))
end

functions["strchr"] = function()
    local ptr, value = get_args()
    value = bit32.band(value, 0xFF)
    local offset = 0
    while true do
        local byte = read_byte(ptr, offset)
        if byte == value then
            push_args(ptr + offset)
            return
        end
        if byte == 0 then
            push_args(0)
            return
        end
        offset += 1
    end
end

functions["strrchr"] = function()
    local ptr, value = get_args()
    value = bit32.band(value, 0xFF)
    local last = 0
    local offset = 0
    while true do
        local byte = read_byte(ptr, offset)
        if byte == value then
            last = ptr + offset
        end
        if byte == 0 then
            break
        end
        offset += 1
    end
    push_args(last)
end

functions["strstr"] = function()
    local haystack, needle = get_args()
    if needle == 0 or read_byte(needle, 0) == 0 then
        push_args(haystack)
        return
    end
    local needle_len = c_strlen(needle)
    local offset = 0
    while read_byte(haystack, offset) ~= 0 do
        if compare_bytes(haystack + offset, needle, needle_len, false) == 0 then
            push_args(haystack + offset)
            return
        end
        offset += 1
    end
    push_args(0)
end

functions["strtok"] = function()
    local str, delim_ptr = get_args()
    if str ~= 0 then
        strtok_state.next = str
        if delim_ptr ~= 0 then
            strtok_state.delim = read_string(delim_ptr)
        end
    end

    local current = strtok_state.next
    if current == 0 then
        push_args(0)
        return
    end

    while is_delimiter(read_byte(current, 0), strtok_state.delim) and read_byte(current, 0) ~= 0 do
        current += 1
    end

    if read_byte(current, 0) == 0 then
        strtok_state.next = 0
        push_args(0)
        return
    end

    local start = current
    while read_byte(current, 0) ~= 0 and not is_delimiter(read_byte(current, 0), strtok_state.delim) do
        current += 1
    end

    if read_byte(current, 0) == 0 then
        strtok_state.next = 0
    else
        write_byte(current, 0, 0)
        strtok_state.next = current + 1
    end

    push_args(start)
end
