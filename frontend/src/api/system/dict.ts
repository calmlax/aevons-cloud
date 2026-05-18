import request, { download } from '@/utils/request'

/**
 * 字典类型表 Model
 *
 * @author 
 * @date 2026-04-13
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
        return request.get('/v1/sys/dict/page',{params: query})
    },

    /**
     * 查询列表
     * @param query 查询参数
     */
    list: (query: any) => {
        return request.get('/v1/sys/dict/list',{params: query})
    },

    /**
     * 查询详细
     * @param query 查询参数
     */
    getById: (id: string) => {
        return request.get(`/v1/sys/dict/${id}`)
    },

    /**
     * 新增
     * @param data 新增数据
     */
    add: (data: any) => {
        return request.post('/v1/sys/dict',data)
    },

    /**
     * 修改
     * @param data 修改数据
     */
    update: (id: string,data: any) => {
        return request.put(`/v1/sys/dict/${id}`,data)
    },

    /**
     * 删除
     * @param ids 删除Ids
     */
    delete: (ids: any) => {
        return request.delete(`/v1/sys/dict/${ids}`)
    },

    /**
     * 导出Excel
     * @param filename Excel文件名称
     * @param query 查询参数
     */
    download: (query: any, filename: string) => {
        download("/v1/sys/dict/export", { ...query }, `${filename}.xlsx`,{});
    },

    getByDictType:(dictType: string) =>{
        return request.get(`/v1/sys/dict/type/${dictType}`)
    },

    refreshCache:() =>{
       return request.delete(`/v1/sys/dict/refresh-cache`)
    }

}