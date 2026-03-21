--- <stdlib.h> helpers

functions["malloc"] = function()
    local size = get_args()
    local dest = malloc(size)

    push_args(dest)
end
