--- <math.h> helpers
local function trunc(x: number): number
    return x > 0 and math.floor(x) or math.ceil(x)
end

local function round(x: number): number
    return math.floor(x + 0.5)
end

local function nearbyint(x: number): number
    return round(x)
end

local math_float_returns = {
    sin = {fn = math.sin, argc = 1}, sinf = {fn = math.sin, argc = 1},
    cos = {fn = math.cos, argc = 1}, cosf = {fn = math.cos, argc = 1},
    tan = {fn = math.tan, argc = 1}, tanf = {fn = math.tan, argc = 1},
    asin = {fn = math.asin, argc = 1}, acos = {fn = math.acos, argc = 1}, atan = {fn = math.atan, argc = 1}, atan2 = {fn = math.atan2, argc = 2},
    sinh = {fn = math.sinh, argc = 1}, cosh = {fn = math.cosh, argc = 1}, tanh = {fn = math.tanh, argc = 1},
    asinh = {fn = function(x) return math.log(x + math.sqrt(x*x + 1)) end, argc = 1},
    acosh = {fn = function(x) return math.log(x + math.sqrt(x*x - 1)) end, argc = 1},
    atanh = {fn = function(x) return 0.5 * math.log((1 + x) / (1 - x)) end, argc = 1},
    exp = {fn = math.exp, argc = 1}, expf = {fn = math.exp, argc = 1},
    log = {fn = math.log, argc = 1}, logf = {fn = math.log, argc = 1},
    log10 = {fn = math.log10, argc = 1}, log10f = {fn = math.log10, argc = 1},
    log2 = {fn = function(x) return math.log(x) / math.log(2) end, argc = 1},
    sqrt = {fn = math.sqrt, argc = 1}, sqrtf = {fn = math.sqrt, argc = 1},
    cbrt = {fn = function(x) return x^(1/3) end, argc = 1},
    pow = {fn = function(x,y) return x^y end, argc = 2},
    hypot = {fn = function(x,y) return math.sqrt(x^2 + y^2) end, argc = 2},
    ceil = {fn = math.ceil, argc = 1}, ceilf = {fn = math.ceil, argc = 1},
    floor = {fn = math.floor, argc = 1}, floorf = {fn = math.floor, argc = 1},
    trunc = {fn = trunc, argc = 1}, truncf = {fn = trunc, argc = 1},
    round = {fn = round, argc = 1}, roundf = {fn = round, argc = 1},
    nearbyint = {fn = nearbyint, argc = 1}, nearbyintf = {fn = nearbyint, argc = 1},
    fabs = {fn = math.abs, argc = 1}, fabsf = {fn = math.abs, argc = 1},
    copysign = {fn = function(x,y) return math.abs(x) * (y>=0 and 1 or -1) end, argc = 2},
    fmax = {fn = math.max, argc = 2}, fmaxf = {fn = math.max, argc = 2},
    fmin = {fn = math.min, argc = 2}, fminf = {fn = math.min, argc = 2},
    fma = {fn = function(x,y,z) return x*y + z end, argc = 3},
    fmaf = {fn = function(x,y,z) return x*y + z end, argc = 3},
    fmod = {fn = function(x,y) return x % y end, argc = 2},
}

local math_int_returns = {
    isnan = math.isnan,
    isinf = math.isinf,
    signbit = function(x) return x < 0 and 1 or 0 end,
}

for name, spec in pairs(math_float_returns) do
    functions[name] = function()
        local a1, a2, a3 = get_f_args()
        if spec.argc == 1 then
            push_f_args(spec.fn(a1))
        elseif spec.argc == 2 then
            push_f_args(spec.fn(a1, a2))
        else
            push_f_args(spec.fn(a1, a2, a3))
        end
    end
end

for name, func in pairs(math_int_returns) do
    functions[name] = function()
        push_args(func(get_f_args()))
    end
end
