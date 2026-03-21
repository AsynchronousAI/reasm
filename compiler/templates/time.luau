--- <time.h> helpers

local tm_size: number = 9 * 4
local tm_cache: number = 0

local function ensure_tm_cache()
    if tm_cache == 0 then
        tm_cache = malloc(tm_size)
    end
    return tm_cache
end

local function read_tm(ptr: number)
    if ptr == 0 then
        return nil
    end
    return {
        sec = buffer.readi32(memory, ptr),
        min = buffer.readi32(memory, ptr + 4),
        hour = buffer.readi32(memory, ptr + 8),
        mday = buffer.readi32(memory, ptr + 12),
        mon = buffer.readi32(memory, ptr + 16),
        year = buffer.readi32(memory, ptr + 20),
        wday = buffer.readi32(memory, ptr + 24),
        yday = buffer.readi32(memory, ptr + 28),
        isdst = buffer.readi32(memory, ptr + 32),
    }
end

local function write_tm(ptr: number, tm)
    if ptr == 0 or not tm then
        return
    end
    buffer.writei32(memory, ptr, tm.sec or 0)
    buffer.writei32(memory, ptr + 4, tm.min or 0)
    buffer.writei32(memory, ptr + 8, tm.hour or 0)
    buffer.writei32(memory, ptr + 12, tm.mday or 0)
    buffer.writei32(memory, ptr + 16, (tm.mon or 1) - 1)
    buffer.writei32(memory, ptr + 20, (tm.year or 0) - 1900)
    buffer.writei32(memory, ptr + 24, ((tm.wday or 1) - 1))
    buffer.writei32(memory, ptr + 28, ((tm.yday or 1) - 1))
    buffer.writei32(memory, ptr + 32, tm.isdst and 1 or 0)
end

local function table_from_tm(ptr: number)
    local tm = read_tm(ptr)
    if not tm then
        return nil
    end
    return {
        sec = tm.sec,
        min = tm.min,
        hour = tm.hour,
        day = tm.mday,
        month = tm.mon + 1,
        year = tm.year + 1900,
        yday = tm.yday + 1,
        wday = tm.wday,
        isdst = tm.isdst,
    }
end

functions["time"] = function()
    local target = get_args()
    local now = os.time()
    if target ~= 0 then
        buffer.writei32(memory, target, now)
    end
    push_args(now)
end

functions["clock"] = function()
    push_args(os.clock())
end

functions["difftime"] = function()
    local a, b = get_args()
    push_args(a - b)
end

functions["mktime"] = function()
    local tm_ptr = get_args()
    local tm_value = read_tm(tm_ptr)
    if not tm_value then
        push_args(-1)
        return
    end
    local timestamp = os.time({
        year = (tm_value.year or 0) + 1900,
        month = (tm_value.mon or 0) + 1,
        day = tm_value.mday or 1,
        hour = tm_value.hour or 0,
        min = tm_value.min or 0,
        sec = tm_value.sec or 0,
    })
    local normalized = os.date("*t", timestamp)
    write_tm(tm_ptr, normalized)
    push_args(timestamp)
end

local function fill_tm_from_timestamp(ptr: number, timestamp: number, utc: boolean)
    local format = utc and "!*t" or "*t"
    local tm_table = os.date(format, timestamp)
    write_tm(ptr, tm_table)
end

functions["gmtime"] = function()
    local timer = get_args()
    local timestamp = 0
    if timer ~= 0 then
        timestamp = buffer.readi32(memory, timer)
    end
    if timestamp == 0 then
        timestamp = os.time()
    end
    local result_ptr = ensure_tm_cache()
    fill_tm_from_timestamp(result_ptr, timestamp, true)
    push_args(result_ptr)
end

functions["localtime"] = function()
    local timer = get_args()
    local timestamp = 0
    if timer ~= 0 then
        timestamp = buffer.readi32(memory, timer)
    end
    if timestamp == 0 then
        timestamp = os.time()
    end
    local result_ptr = ensure_tm_cache()
    fill_tm_from_timestamp(result_ptr, timestamp, false)
    push_args(result_ptr)
end

functions["strftime"] = function()
    local buf, maxsize, fmt_ptr, tm_ptr = get_args()
    local format = read_string(fmt_ptr)
    local tm_table = table_from_tm(tm_ptr)
    if not tm_table then
        push_args(0)
        return
    end
    local timestamp = os.time({
        year = tm_table.year,
        month = tm_table.month,
        day = tm_table.day,
        hour = tm_table.hour,
        min = tm_table.min,
        sec = tm_table.sec,
    })
    local formatted = os.date(format, timestamp)
    local capacity = math.max(0, maxsize - 1)
    local length = math.min(#formatted, capacity)
    for i = 1, length do
        buffer.writeu8(memory, buf + i - 1, string.byte(formatted, i))
    end
    buffer.writeu8(memory, buf + length, 0)
    push_args(length)
end

functions["sleep"] = function()
    local seconds = get_args()
    if task and task.wait then
        task.wait(seconds)
    end
    push_args(0)
end

functions["nanosleep"] = function()
    local req_ptr, rem_ptr = get_args()
    local secs = buffer.readi32(memory, req_ptr)
    local nsecs = buffer.readi32(memory, req_ptr + 4)
    if task and task.wait then
        task.wait(secs + nsecs / 1e9)
    end
    if rem_ptr and rem_ptr ~= 0 then
        buffer.writei32(memory, rem_ptr, 0)
        buffer.writei32(memory, rem_ptr + 4, 0)
    end
    push_args(0)
end
