local core = require("apisix.core")
local cjson = require("cjson.safe")

local plugin_name = "audit-log"

local schema = {
    type = "object",
    properties = {
        include_headers = {
            type = "boolean",
            default = false,
        },
    },
}

local _M = {
    version = 0.2,
    priority = 2200,
    name = plugin_name,
    schema = schema,
}

function _M.check_schema(conf)
    return core.schema.check(schema, conf)
end

function _M.log(conf, ctx)
    local user = ctx.aevons_login_user or {}
    local record = {
        plugin = plugin_name,
        request_id = ctx.var.request_id,
        method = ctx.var.request_method,
        uri = ctx.var.uri,
        status = ngx.status,
        client_ip = core.request.get_remote_client_ip(ctx),
        user_id = user.user_id,
        username = user.username,
        client_id = user.client_id,
        upstream_addr = ctx.var.upstream_addr,
        service_name = ctx.service_name,
    }

    if conf.include_headers then
        record.headers = ngx.req.get_headers()
    end

    core.log.warn("audit_log ", cjson.encode(record))
end

return _M
