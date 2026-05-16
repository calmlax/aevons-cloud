local plugin_name = "client-resource-auth"

local _M = {
    version = 0.1,
    priority = 2500,
    name = plugin_name,
    schema = {
        type = "object",
        properties = {
            client_id = { type = "string" },
            resources = {
                type = "array",
                items = { type = "string" }
            }
        }
    }
}

function _M.check_schema(_)
    return true
end

function _M.access(conf, ctx)
    -- 第一阶段占位：
    -- 后续在这里实现 client_id 识别、客户端状态校验、
    -- 资源匹配和 ALL 语义约束。
end

return _M
