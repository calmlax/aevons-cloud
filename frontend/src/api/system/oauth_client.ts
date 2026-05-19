import request, { download } from '@/utils/request'

/**
 * 终端应用 Model
 *
 * @author 
 * @date 2026-04-16
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
export default {

    /**
     * 分页查询
     * @param query 查询参数
     */
    page: (query: any) => {
        return request.get('/sys/v1/oauth/client/page',{params: query})
    },

    /**
     * 查询列表
     * @param query 查询参数
     */
    list: (query: any) => {
        return request.get('/sys/v1/oauth/client/list',{params: query})
    },

    /**
     * 查询详细
     * @param query 查询参数
     */
    getById: (id: string) => {
        return request.get(`/sys/v1/oauth/client/${id}`)
    },

    /**
     * 新增
     * @param data 新增数据
     */
    add: (data: any) => {
        return request.post('/sys/v1/oauth/client',data)
    },

    /**
     * 修改
     * @param data 修改数据
     */
    update: (id: string,data: any) => {
        return request.put(`/sys/v1/oauth/client/${id}`,data)
    },

    /**
     * 删除
     * @param ids 删除Ids
     */
    delete: (ids: any) => {
        return request.delete(`/sys/v1/oauth/client/${ids}`)
    },

    /**
     * 导出Excel
     * @param filename Excel文件名称
     * @param query 查询参数
     */
    download: (query: any, filename: string) => {
        download("/sys/v1/oauth/client/export", { ...query }, `${filename}.xlsx`,{});
    },

    /**
     * 刷新网关 oauth_client 资源缓存
     */
    refreshGatewayCache: () => {
        return request.post('/sys/v1/oauth/client/refresh-cache')
    },

}
