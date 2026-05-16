local plugin_name = "ai-risk-control"

local _M = {
    version = 0.1,
    priority = 2100,
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
    -- 预留 AI 风控、异常流量识别与 AI WAF 扩展位点。
end

return _M
