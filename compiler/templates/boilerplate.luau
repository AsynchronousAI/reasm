--!strict
--!native
--!optimize 2

-- Compiled from RISC-V assembly.
-- API
local mem: number = --{memory here} -- bytes of RAM
local memory: buffer = buffer.create(mem) -- our memory!
local temp: buffer = buffer.create(128) -- just used for reinterpreting
local r1: number, r2: number, r3: number, r4: number, r5: number, r6: number, r7: number, r8: number, r9: number, r10: number, r11: number, r12: number, r13: number, r14: number, r15: number, r16: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
local r17: number, r18: number, r19: number, r20: number, r21: number, r22: number, r23: number, r24: number, r25: number, r26: number, r27: number, r28: number, r29: number, r30: number, r31: number, r32: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
local r33: number, r34: number, r35: number, r36: number, r37: number, r38: number, r39: number, r40: number, r41: number, r42: number, r43: number, r44: number, r45: number, r46: number, r47: number, r48: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
local r49: number, r50: number, r51: number, r52: number, r53: number, r54: number, r55: number, r56: number, r57: number, r58: number, r59: number, r60: number, r61: number, r62: number, r63: number, r64: number = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0

-- Variables
local PC: number = 1 -- current position
local mallocDepth: number = 0
local stdout_cache: string = ""
local errno: number = 0
local fflags: number = 0 -- floating-point exception flags (NV,DZ,OF,UF,NX)
local files: {[string]: buffer} = {}
local open_files: {[number]: {buffer: buffer, pos: number, name: string}} = {}
local next_fd: number = 1

-- Utility
--- Base (can be found in generated code)
local function two_words_to_double(highWord: number, lowWord: number): number
    buffer.writei32(temp, 0, ((math.floor(lowWord % 0x100000000) + 0x100000000) % 0x100000000))
    buffer.writei32(temp, 4, ((math.floor(highWord % 0x100000000) + 0x100000000) % 0x100000000))
    return buffer.readf64(temp, 0)
end
function reset_registers(): ()
    r1,r2,r3,r4,r5,r6,r7,r8,r9,r10,r11,r12,r13,r14,r15,r16 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r17,r18,r19,r20,r21,r22,r23,r24,r25,r26,r27,r28,r29,r30,r31,r32 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r33,r34,r35,r36,r37,r38,r39,r40,r41,r42,r43,r44,r45,r46,r47,r48 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r49,r50,r51,r52,r53,r54,r55,r56,r57,r58,r59,r60,r61,r62,r63,r64 = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
    r3 = mem - 1024 -- x2/sp starts at top of stack minus padding
    fflags = 0
end

--- Strings
local function read_string(ptr: number): string
    local startPointer = ((math.floor(ptr % 0x100000000) + 0x100000000) % 0x100000000)
    local endPointer = startPointer
    while endPointer < mem do
        if buffer.readu8(memory, endPointer) == 0 then
            return buffer.readstring(memory, startPointer, endPointer - startPointer)
        end
        endPointer += 1
    end
    error("Exceeded buffer size when reading string.")
end
function format_string(fmt: number, args: {number}): string
    local fmt_ptr = ((math.floor(fmt % 0x100000000) + 0x100000000) % 0x100000000)
    local fmtString: string = read_string(fmt_ptr)
    local arg_index: number = 1

    return fmtString:gsub("%%([%d%.]*[lL]?[dfseixXpus])", function(spec)
        local fmtSpec = spec:gsub("[lL]", "")
        local type_char = spec:sub(-1)
        if type_char == "d" or type_char == "i" then
            local val = args[arg_index]
            arg_index += 1

           	local n = --[[u32]] ((math.floor(val % 0x100000000) + 0x100000000) % 0x100000000);
            return string.format("%"..fmtSpec, --[[i32]] (n >= 0x80000000) and (n - 0x100000000) or n)
        elseif type_char == "x" or type_char == "X" or type_char == "p" or type_char == "u" then
            local val = args[arg_index]
            arg_index += 1
            if type_char == "p" then 
                fmtSpec = fmtSpec:gsub("p", "x")
            end
            return string.format("%"..fmtSpec, --[[u32]]((math.floor(val % 0x100000000) + 0x100000000) % 0x100000000))
        elseif type_char == "f" then
            if arg_index % 2 == 1 then
                arg_index += 1 -- 64-bit varargs are aligned to even a-registers (skip a1/a3/...)
            end
            local low = args[arg_index]; arg_index += 1
            local high = args[arg_index]; arg_index += 1
            local float = two_words_to_double(high, low)
            return string.format("%"..fmtSpec, float)
        elseif type_char == "s" then
            local val = args[arg_index]
            arg_index += 1
            return string.format("%"..fmtSpec, read_string(val))
        elseif type_char == "e" then
            if arg_index % 2 == 1 then
                arg_index += 1 -- 64-bit varargs are aligned to even a-registers (skip a1/a3/...)
            end
            local low = args[arg_index]; arg_index += 1
            local high = args[arg_index]; arg_index += 1
            local float = two_words_to_double(high, low)
            return string.format("%"..fmtSpec, float)
        else
            return spec
        end
    end)
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
local function malloc(size: number): number
    mallocDepth += size
    return buffer.len(memory) - mallocDepth
end


-- Functions
local functions = {}

local syscalls: {[number]: () -> ()} = {}

-- Extensions
--{extentions here}

-- Localized Functions
local FUNCS: {[number]: () -> boolean} = {}
---- Auto generated code starts here
--{code here}
