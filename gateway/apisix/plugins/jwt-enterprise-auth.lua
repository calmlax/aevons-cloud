local core = require("apisix.core")
local http = require("resty.http")
local cjson = require("cjson.safe")

local plugin_name = "jwt-enterprise-auth"

local default_excluded_paths = {
    "/api/v1/auth/login",
    "/api/v1/auth/refresh",
    "/api/v1/auth/email/code",
    "/api/v1/auth/register",
    "/api/v1/auth/reset-password",
    "/api/v1/auth/public-key",
    "/api/v1/auth/authorize",
    "/api/v1/auth/callback",
    "/api/v1/auth/passkey/login/begin",
    "/api/v1/auth/passkey/login/finish",
}

local schema = {
    type = "object",
    properties = {
        auth_service_url = {
            type = "string",
            minLength = 1,
            default = "http://127.0.0.1:10701",
        },
        user_info_path = {
            type = "string",
            minLength = 1,
            default = "/api/v1/auth/user",
        },
        timeout = {
            type = "integer",
            minimum = 100,
            default = 2000,
        },
        excluded_paths = {
            type = "array",
            items = { type = "string", minLength = 1 },
        },
    },
}

local _M = {
    version = 0.2,
    priority = 2600,
    name = plugin_name,
    schema = schema,
}

local function respond(status, message)
    return core.response.exit(status, {
        code = status,
        message = message,
    })
end

local function is_excluded(path, patterns)
    if not patterns then
        return false
    end

    for _, pattern in ipairs(patterns) do
        if pattern == path then
            return true
        end

        if #pattern >= 2 and string.sub(pattern, -2) == "/*" then
            local prefix = string.sub(pattern, 1, #pattern - 1)
            if string.sub(path, 1, #prefix) == prefix then
                return true
            end
        end
    end

    return false
end

local function normalize_array(value)
    if type(value) == "table" then
        return value
    end
    return {}
end

function _M.check_schema(conf)
    return core.schema.check(schema, conf)
end

function _M.access(conf, ctx)
    local path = ctx.var.uri or ""
    local excluded_paths = conf.excluded_paths or default_excluded_paths
    if is_excluded(path, excluded_paths) then
        return
    end

    local auth_header = ngx.req.get_headers()["Authorization"]
    if not auth_header or string.sub(auth_header, 1, 7) ~= "Bearer " then
        return respond(401, "authorization.token.missing")
    end

    local token = string.sub(auth_header, 8)
    if token == "" then
        return respond(401, "authorization.token.missing")
    end

    local httpc = http.new()
    httpc:set_timeout(conf.timeout or 2000)

    local res, err = httpc:request_uri((conf.auth_service_url or "http://127.0.0.1:10701") ..
        (conf.user_info_path or "/api/v1/auth/user"), {
        method = "GET",
        headers = {
            Authorization = auth_header,
            ["X-Request-Id"] = ctx.var.request_id,
        },
        ssl_verify = false,
    })

    if not res then
        core.log.error("auth service request failed: ", err)
        return respond(503, "authorization.auth_service.unavailable")
    end

    if res.status ~= 200 then
        core.log.warn("auth service rejected token, status: ", res.status, ", body: ", res.body)
        if res.status >= 500 then
            return respond(503, "authorization.auth_service.unavailable")
        end
        return respond(401, "authorization.token.invalid")
    end

    local payload, decode_err = cjson.decode(res.body)
    if not payload then
        core.log.error("decode auth service response failed: ", decode_err)
        return respond(503, "authorization.auth_service.invalid_response")
    end

    local user = payload.data
    if payload.code ~= 0 or type(user) ~= "table" then
        core.log.warn("auth service returned invalid payload: ", res.body)
        return respond(401, "authorization.token.invalid")
    end

    user.roles = normalize_array(user.roles)
    user.depts = normalize_array(user.depts)
    user.permissions = normalize_array(user.permissions)

    ctx.aevons_login_user = user
    ctx.aevons_permissions = user.permissions

    ngx.req.set_header("X-User-Id", tostring(user.user_id or ""))
    ngx.req.set_header("X-Username", user.username or "")
    ngx.req.set_header("X-Nickname", user.nickname or "")
    ngx.req.set_header("X-Client-Id", user.client_id or "")
    ngx.req.set_header("X-User-Permissions", cjson.encode(user.permissions))
    ngx.req.set_header("X-User-Roles", cjson.encode(user.roles))
    ngx.req.set_header("X-User-Depts", cjson.encode(user.depts))
end

return _M
