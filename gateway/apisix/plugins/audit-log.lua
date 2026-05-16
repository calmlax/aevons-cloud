local plugin_name = "audit-log"

local _M = {
    version = 0.1,
    priority = 2200,
    name = plugin_name,
    schema = {
        type = "object",
        properties = {}
    }
}

function _M.check_schema(_)
    return true
end

function _M.log(conf, ctx)
    -- 第一阶段占位：
    -- 后续在这里接入审计日志和安全事件上报。
end

return _M
