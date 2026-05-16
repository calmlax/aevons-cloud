local plugin_name = "jwt-enterprise-auth"

local _M = {
    version = 0.1,
    priority = 2400,
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
    -- 后续在这里实现 JWT 校验、用户上下文注入和角色上下文构建。
end

return _M
