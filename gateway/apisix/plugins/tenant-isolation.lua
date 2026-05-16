local plugin_name = "tenant-isolation"

local _M = {
    version = 0.1,
    priority = 2300,
    name = plugin_name,
    schema = {
        type = "object",
        properties = {}
    }
}

function _M.check_schema(_)
    return true
end

function _M.access(conf, ctx)
    -- 第一阶段占位：
    -- 后续在这里识别租户并注入租户上下文。
end

return _M
