import request, { download } from '@/utils/request'

/**
 * 角色信息表 Model
 *
 * @author 
 * @date 2026-04-18
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
        return request.get('/sys/v1/role/page',{params: query})
    },

    /**
     * 查询列表
     * @param query 查询参数
     */
    list: (query: any) => {
        return request.get('/sys/v1/role/list',{params: query})
    },

    /**
     * 查询详细
     * @param query 查询参数
     */
    getById: (id: string) => {
        return request.get(`/sys/v1/role/${id}`)
    },

    /**
     * 新增
     * @param data 新增数据
     */
    add: (data: any) => {
        return request.post('/sys/v1/role',data)
    },

    /**
     * 修改
     * @param data 修改数据
     */
    update: (id: string,data: any) => {
        return request.put(`/sys/v1/role/${id}`,data)
    },

    /**
     * 删除
     * @param ids 删除Ids
     */
    delete: (ids: any) => {
        return request.delete(`/sys/v1/role/${ids}`)
    },

    /**
     * 获取角色已关联的菜单ID列表
     */
    getMenuIds: (id: string) => {
        return request.get(`/sys/v1/role/${id}/menu`)
    },

    /**
     * 获取角色已关联的部门ID列表
     */
    getDeptIds: (id: string) => {
        return request.get(`/sys/v1/role/${id}/dept`)
    },

    /**
     * 导出Excel
     * @param filename Excel文件名称
     * @param query 查询参数
     */
    download: (query: any, filename: string) => {
        download("/sys/v1/role/export", { ...query }, `${filename}.xlsx`,{});
    },

}
