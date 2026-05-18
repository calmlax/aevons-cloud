import request, { download } from '@/utils/request'

/**
 * 参数配置表 Model
 *
 * @author 
 * @date 2026-04-17
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
        return request.get('/v1/sys/conf/page',{params: query})
    },

    /**
     * 查询列表
     * @param query 查询参数
     */
    list: (query: any) => {
        return request.get('/v1/sys/conf/list',{params: query})
    },

    /**
     * 查询详细
     * @param query 查询参数
     */
    getById: (id: string) => {
        return request.get(`/v1/sys/conf/${id}`)
    },

    /**
     * 新增
     * @param data 新增数据
     */
    add: (data: any) => {
        return request.post('/v1/sys/conf',data)
    },

    /**
     * 修改
     * @param data 修改数据
     */
    update: (id: string,data: any) => {
        return request.put(`/v1/sys/conf/${id}`,data)
    },

    /**
     * 删除
     * @param ids 删除Ids
     */
    delete: (ids: any) => {
        return request.delete(`/v1/sys/conf/${ids}`)
    },

    /**
     * 导出Excel
     * @param filename Excel文件名称
     * @param query 查询参数
     */
    download: (query: any, filename: string) => {
        download("/v1/sys/conf/export", { ...query }, `${filename}.xlsx`,{});
    },

}