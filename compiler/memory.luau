--- <memory.h> helpers

functions["memcpy"] = function()
    local dest, src, count = get_args()
    buffer.copy(memory, dest, memory, src, count)
end

functions["memset"] = function()
    local dest, value, count = get_args()
    buffer.fill(memory, dest, bit32.band(value, 0xFF), count)
end
