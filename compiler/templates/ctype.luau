--- <ctype.h> helpers

local function sanitize_char(value)
    return string.char(bit32.band(value, 0xFF))
end

local function match_char(value, pattern)
    local ch = sanitize_char(value)
    return string.match(ch, pattern) ~= nil
end

local function classification(result)
    push_args(result and 1 or 0)
end

functions["isalpha"] = function()
    local value = get_args()
    classification(match_char(value, "%a"))
end

functions["isdigit"] = function()
    local value = get_args()
    classification(match_char(value, "%d"))
end

functions["isalnum"] = function()
    local value = get_args()
    classification(match_char(value, "%w"))
end

functions["isupper"] = function()
    local value = get_args()
    classification(match_char(value, "%u"))
end

functions["islower"] = function()
    local value = get_args()
    classification(match_char(value, "%l"))
end

functions["isspace"] = function()
    local value = get_args()
    classification(match_char(value, "%s"))
end

functions["ispunct"] = function()
    local value = get_args()
    classification(match_char(value, "%p"))
end

functions["isprint"] = function()
    local value = get_args()
    local byte = bit32.band(value, 0xFF)
    classification(byte >= 0x20 and byte <= 0x7E)
end

functions["toupper"] = function()
    local value = get_args()
    push_args(string.byte(string.upper(sanitize_char(value))))
end

functions["tolower"] = function()
    local value = get_args()
    push_args(string.byte(string.lower(sanitize_char(value))))
end
