--!strict
--!native
--!optimize 2

-- Compiled from RISC-V assembly.
-- API
local mem: number = 2048 -- 2KB of RAM
local memory: buffer = buffer.create(mem) -- our memory!
local r1: number, r2: number, r3: number, r4: number, r5: number, r6: number, r7: number, r8: number, r9: number, r10: number, r11: number, r12: number, r13: number, r14: number, r15: number, r16: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
local r17: number, r18: number, r19: number, r20: number, r21: number, r22: number, r23: number, r24: number, r25: number, r26: number, r27: number, r28: number, r29: number, r30: number, r31: number, r32: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
local r33: number, r34: number, r35: number, r36: number, r37: number, r38: number, r39: number, r40: number, r41: number, r42: number, r43: number, r44: number, r45: number, r46: number, r47: number, r48: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
local r49: number, r50: number, r51: number, r52: number, r53: number, r54: number, r55: number, r56: number, r57: number, r58: number, r59: number, r60: number, r61: number, r62: number, r63: number, r64: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0

-- Variables
local PC: number = 1 -- current position
local mallocDepth: number = 0
local stdout_cache: string = ""

-- Utility
--- 32-bit helpers
local function u32(v: number): number
    local n = bit32.band(v, 0xFFFFFFFF)
    if n < 0 then
        return n + 0x100000000
    end
    return n
end
local function i32(v)
    local n = u32(v)
    if n >= 0x80000000 then
        return n - 0x100000000
    end
    return n
end

--- Base (can be found in generated code)
local function idiv_trunc(a: number, b: number): number
    if b == 0 then error("division by zero") end
    if a >= 0 then
        return (a - (a % b)) // b
    else
        return -((-a) - ((-a) % b)) // b
    end
end
local function float_to_int(f: number): number
    return string.unpack("i", string.pack("f", f))
end
local function int_to_float(i: number): number
    local packed = string.pack("I4", i)
    return string.unpack("f", packed)
end
local function float_to_double(f: number): number
    local packed_f = string.pack("f", f)
    local padded = packed_f .. ("\0\0\0\0")
    return string.unpack("d", padded)
end
local function two_words_to_double(highWord: number, lowWord: number): number
    local packed = string.pack("<I4I4", u32(lowWord), u32(highWord))
    return string.unpack("d", packed)
end
local function hi(addr: number): number
    return bit32.lshift(bit32.rshift(addr, 12), 12)
end
local function lo(addr: number): number
    return bit32.band(addr, 0xFFF)
end
function fclass(x: number): number
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
function reset_registers(): ()
    r1,r2,r3,r4,r5,r6,r7,r8,r9,r10,r11,r12,r13,r14,r15,r16 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r17,r18,r19,r20,r21,r22,r23,r24,r25,r26,r27,r28,r29,r30,r31,r32 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r33,r34,r35,r36,r37,r38,r39,r40,r41,r42,r43,r44,r45,r46,r47,r48 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r49,r50,r51,r52,r53,r54,r55,r56,r57,r58,r59,r60,r61,r62,r63,r64 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
end

--- Strings
local function read_string(startPointer: number): string
    -- read null terminated strings from memory
    local pointer: number = startPointer
    local str: string = ""
    local byte: number = 0
    repeat
        byte = buffer.readbits(memory, pointer * 8, 8)
        if byte == 0 then break end
        str = str .. string.char(byte)
        pointer = pointer + 1
        if pointer >= mem then error("Exceeded buffer size when reading string.") end
    until false
    return str
end
function format_string(fmt: number, args: {number}): string
    local fmtString: string = read_string(fmt)
    local arg_index: number = 1

    return fmtString:gsub("%%([%d%.]*[dfseX])", function(spec)
        local val = args[arg_index]
        arg_index += 1

        if spec:sub(-1) == "d" then
            return tostring(i32(val))
        elseif spec:sub(-1) == "X" then
            return string.format("%X", u32(val))
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
end


--- Memory
local function malloc(size: number): number
    mallocDepth+=size
    return buffer.len(memory)-mallocDepth
end

--- Args
local function get_args(): (number, number, number, number, number, number, number, number)
    return r11, r12, r13, r14, r15, r16, r17, r18
end
local function push_args(a1: number?, a2: number?, a3: number?, a4: number?, a5: number?, a6: number?, a7: number?, a8: number?)
    r11 = a1 or 0
    r12 = a2 or 0
    r13 = a3 or 0
    r14 = a4 or 0
    r15 = a5 or 0
    r16 = a6 or 0
    r17 = a7 or 0
    r18 = a8 or 0
end
local function get_f_args(): (number, number, number, number, number, number, number, number)
    return r43, r44, r45, r46, r47, r48, r49, r50
end
local function push_f_args(a1: number?, a2: number?, a3: number?, a4: number?, a5: number?, a6: number?, a7: number?, a8: number?)
    r43 = a1 or 0
    r44 = a2 or 0
    r45 = a3 or 0
    r46 = a4 or 0
    r47 = a5 or 0
    r48 = a6 or 0
    r49 = a7 or 0
    r50 = a8 or 0
end

--- IO
local function flush_stdout()
    if #stdout_cache > 0 then
        print(stdout_cache)
        stdout_cache = ""
    end
end


-- Functions
local functions = {
    ["memcpy"] = function()
        local dest,src,count = get_args()

        buffer.copy(memory, dest, memory, src, count)
    end,
    ["memset"] = function()
        local dest, value, count = get_args()
        buffer.fill(memory, dest, bit32.band(value, 0xFF), count)
    end,
    ["malloc"] = function()
        local size = get_args()
        local dest = malloc(size)

        push_args(dest)
    end,
    ["putchar"] = function()
        local c = get_args()
        local char = string.char(c)
        if char == "\n" then
            flush_stdout()
        else
            stdout_cache = stdout_cache .. char
        end
    end,
    ["puts"] = function()
        local fmt = get_args()
        local str = read_string(fmt)
        stdout_cache = stdout_cache .. str
        flush_stdout()
    end,
    ["printf"] = function()
        local args = {get_args()}
        local fmt_ptr = args[1]
        table.remove(args, 1)

        local formatted = format_string(fmt_ptr, args)
        for i = 1, #formatted do
            local ch = formatted:sub(i,i)
            if ch == "\n" then
                flush_stdout()
            else
                stdout_cache = stdout_cache .. ch
            end
        end
    end,
}

-- Extensions
--{extentions here}

-- Localized Functions
--- buffer
local writei8, writei16, writei32 = buffer.writei8, buffer.writei16, buffer.writei32
local readi8, readi16, readi32 = buffer.readi8, buffer.readi16, buffer.readi32
local writeu8, writeu16, writeu32, writestring, fill = buffer.writeu8, buffer.writeu16, buffer.writeu32, buffer.writestring, buffer.fill
local readu8, readu16, readu32, readf32, readf64, writef32, writef64 = buffer.readu8, buffer.readu16, buffer.readu32, buffer.readf32, buffer.readf64, buffer.writef32, buffer.writef64

local FUNCS: {[number]: () -> boolean} = {}
---- Auto generated code starts here
--{code here}
