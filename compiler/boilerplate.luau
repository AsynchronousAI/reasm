-- Compiled from RISC-V assembly.
-- API
local mem = 2048 -- 2KB of RAM
local memory = buffer.create(mem) -- our memory!
local registers = table.create(64,0)

-- Variables
local PC = 1 -- current position
local mallocDepth = 0

-- Utility
--- Base (can be found in generated code)
local function idiv_trunc(a, b)
    if b == 0 then error("division by zero") end
    if a >= 0 then
        return (a - (a % b)) // b
    else
        return -((-a) - ((-a) % b)) // b
    end
end
local function float_to_int(f)
    local i = string.unpack("i", string.pack("f", f))
    return i
end
local function int_to_float(i)
    local packed = string.pack("I4", i)
    local f = string.unpack("f", packed)
    return f
end
local function float_to_double(f)
    -- Pack float into 4 bytes
    local packed_f = string.pack("f", f)

    -- Pad with 4 zero bytes to make 8 bytes
    local padded = packed_f .. ("\0\0\0\0")

    -- Unpack as double
    local d = string.unpack("d", padded)
    return d
end
local function two_words_to_double(lowWord, highWord)
    local packed = string.pack("I4I4", lowWord, highWord)
    local d = string.unpack("d", packed)
    return d
end
local function hi(addr)
    return bit32.lshift(bit32.rshift(addr, 12), 12)
end
local function lo(addr)
    return bit32.band(addr, 0xFFF)
end
function fclass(x)
    local result = 0

    if x ~= x then
        -- NaN
        if x < 0 then
            result = bit32.bor(result, bit32.lshift(1, 0)) -- -NaN
        else
            result = bit32.bor(result, bit32.lshift(1, 9)) -- +NaN
        end
    elseif x == math.huge then
        result = bit32.bor(result, bit32.lshift(1, 8)) -- +Inf
    elseif x == -math.huge then
        result = bit32.bor(result, bit32.lshift(1, 1)) -- -Inf
    elseif x == 0 then
        if 1/x == math.huge then
            result = bit32.bor(result, bit32.lshift(1, 5)) -- +Zero
        else
            result = bit32.bor(result, bit32.lshift(1, 4)) -- -Zero
        end
    else
        local absx = math.abs(x)
        local min_normal = 2.2250738585072014e-308 -- 2^-1022
        if absx < min_normal then
            if x > 0 then
                result = bit32.bor(result, bit32.lshift(1, 6)) -- +Subnormal
            else
                result = bit32.bor(result, bit32.lshift(1, 3)) -- -Subnormal
            end
        else
            if x > 0 then
                result = bit32.bor(result, bit32.lshift(1, 7)) -- +Normal
            else
                result = bit32.bor(result, bit32.lshift(1, 2)) -- -Normal
            end
        end
    end

    return result
end
function reset_registers()
    registers = table.create(64,0)
end

--- Strings
local function read_string(startPointer)
    -- read null terminated strings from memory
    local pointer = startPointer
    local str = ""
    local byte
    repeat
        byte = buffer.readbits(memory, pointer * 8, 8)
        if byte == 0 then break end
        str = str .. string.char(byte)
        pointer = pointer + 1
        if pointer >= mem then error("Exceeded buffer size when reading string.") end
    until false
    return str
end
local function format_string(fmt, args)
    -- custom format string because strings are stored as pointers, and floats as ints.
    local arg_index = 1

    local result = fmt:gsub("%%([%d%.]*[dfseX])", function(spec)
        local val = args[arg_index]
        arg_index += 1

        if spec:sub(-1) == "d" then
            return string.format("%d", val)
        elseif spec:sub(-1) == "X" then
            return string.format("%X", val)
        elseif spec:sub(-1) == "f" then
            local low = args[arg_index]; arg_index+=1;
            local high = args[arg_index]; arg_index+=1;
            local float = two_words_to_double(high, math.abs(low))

            return string.format("%"..spec, float)
        elseif spec:sub(-1) == "s" then
            return read_string(val)
        elseif spec:sub(-1) == "e" then
            arg_index+=1 -- skip for alignment
            local low = args[arg_index]; arg_index+=1;
            local high = args[arg_index]; arg_index+=1;
            local float = two_words_to_double(high, low)

            return string.format("%"..spec, float)
        else
            return spec
        end
    end)

    return result
end

--- Memory
local function malloc(size)
    mallocDepth+=size
    return buffer.len(memory)-mallocDepth
end

--- Args
local function get_args()
   return
       registers[11],
       registers[12],
       registers[13],
       registers[14],
       registers[15],
       registers[16],
       registers[17],
       registers[18]
end
local function push_args(...)
    local vals = {...}
    for i = 1, #vals do
        registers[i + 10] = vals[i]
    end
end

local function get_f_args()
   return
       registers[43],
       registers[44],
       registers[45],
       registers[46],
       registers[47],
       registers[48],
       registers[49],
       registers[50]
end
local function push_f_args(...)
    local vals = {...}
    for i = 1, #vals do
        registers[i + 42] = vals[i]
    end
end


-- Functions
local functions = {
    ["printf"] = function()
        local args = {get_args()}
        local str = args[1]
        table.remove(args, 1)
        print(format_string(read_string(str), args))
    end,
    ["puts"] = function()
        print(read_string(get_args()))
    end,
    ["memcpy"] = function()
        local dest,src,count = get_args()

        buffer.copy(memory, dest, memory, src, count)
    end,
    ["malloc"] = function()
        local size = get_args()
        local dest = malloc(size)

        push_args(dest)
    end
}

-- Extensions
--{extentions here}

-- Localized Functions
--- bit32
local band = bit32.band
local bor = bit32.bor
local lshift = bit32.lshift
local rshift = bit32.rshift
local arshift = bit32.arshift
local bxor = bit32.bxor

--- buffer
local writestring = buffer.writestring
local writei8 = buffer.writei8
local writei16 = buffer.writei16
local writei32 = buffer.writei32
local readi8 = buffer.readi8
local readi16 = buffer.readi16
local readi32 = buffer.readi32
local writeu8 = buffer.writeu8
local writeu16 = buffer.writeu16
local writeu32 = buffer.writeu32
local readu8 = buffer.readu8
local readu16 = buffer.readu16
local readu32 = buffer.readu32
local writef32 = buffer.writef32
local writef64 = buffer.writef64
local readf32 = buffer.readf32
local readf64 = buffer.readf64
local fill = buffer.fill

--- math
local sqrt = math.sqrt
local abs = math.abs
local sgn = math.sign
local min = math.min
local max = math.max
local floor = math.floor

local FUNCS = {}
---- Auto generated code starts here
--{code here}
