local core = require("apisix.core")

local plugin_name = "client-resource-auth"

local schema = {
    type = "object",
    properties = {
        client_id = { type = "string", minLength = 1 },
        resources = {
            type = "array",
            items = { type = "string", minLength = 1 },
        },
        allow_all = {
            type = "boolean",
            default = true,
        },
    },
}

local _M = {
    version = 0.2,
    priority = 2500,
    name = plugin_name,
    schema = schema,
}

local function respond(status, message, extra)
    local body = {
        code = status,
        message = message,
    }

    if extra then
        body.data = extra
    end

    return core.response.exit(status, body)
end

local function has_permission(permissions, expected, allow_all)
    for _, permission in ipairs(permissions or {}) do
        if permission == expected then
            return true
        end

        if allow_all and permission == "*" then
            return true
        end
    end

    return false
end

function _M.check_schema(conf)
    return core.schema.check(schema, conf)
end

function _M.access(conf, ctx)
    if not conf.client_id and not conf.resources then
        return
    end

    local user = ctx.aevons_login_user
    if not user then
        return respond(401, "authorization.user_context.missing")
    end

    if conf.client_id and user.client_id ~= conf.client_id then
        return respond(403, "authorization.client.denied", {
            expected_client_id = conf.client_id,
            actual_client_id = user.client_id,
        })
    end

    local resources = conf.resources or {}
    if #resources == 0 then
        return
    end

    local permissions = user.permissions or {}
    local allow_all = conf.allow_all ~= false

    for _, resource in ipairs(resources) do
        if not has_permission(permissions, resource, allow_all) then
            return respond(403, "authorization.permission.denied", {
                required_permission = resource,
            })
        end
    end
end

return _M
