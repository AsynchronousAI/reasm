--- <memory.h> helpers

functions["memcpy"] = function()
    local dest, src, count = get_args()
    buffer.copy(memory, dest, memory, src, count)
end

functions["memset"] = function()
    local dest, value, count = get_args()
    buffer.fill(memory, dest, bit32.band(value, 0xFF), count)
end

functions["memmove"] = function()
    local dest, src, count = get_args()
    if count <= 0 or dest == src then
        push_args(dest)
        return
    end

    if src < dest and src + count > dest then
        local i = count - 1
        while i >= 0 do
            buffer.writeu8(memory, dest + i, buffer.readu8(memory, src + i))
            i -= 1
        end
    else
        local i = 0
        while i < count do
            buffer.writeu8(memory, dest + i, buffer.readu8(memory, src + i))
            i += 1
        end
    end

    push_args(dest)
end

functions["memcmp"] = function()
    local a, b, count = get_args()
    local i = 0
    while i < count do
        local byteA = buffer.readu8(memory, a + i)
        local byteB = buffer.readu8(memory, b + i)
        if byteA ~= byteB then
            push_args(byteA - byteB)
            return
        end
        i += 1
    end
    push_args(0)
end
