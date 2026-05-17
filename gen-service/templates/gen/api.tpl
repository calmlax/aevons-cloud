import request, { download } from '@/utils/request'

/**
 * {{.Comment}} Model
 *
 * @author {{.Author}}
 * @date {{.Date}}
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
        return request.get('{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}/page',{params: query})
    },

    /**
     * 查询列表
     * @param query 查询参数
     */
    list: (query: any) => {
        return request.get('{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}/list',{params: query})
    },

    /**
     * 查询详细
     * @param query 查询参数
     */
    getById: (id: string) => {
        return request.get(`{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}/${id}`)
    },

    /**
     * 新增
     * @param data 新增数据
     */
    add: (data: any) => {
        return request.post('{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}',data)
    },

    /**
     * 修改
     * @param data 修改数据
     */
    update: (id: string,data: any) => {
        return request.put(`{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}/${id}`,data)
    },

    /**
     * 删除
     * @param ids 删除Ids
     */
    delete: (ids: any) => {
        return request.delete(`{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}/${ids}`)
    },

    /**
     * 导出Excel
     * @param filename Excel文件名称
     * @param query 查询参数
     */
    download: (query: any, filename: string) => {
        download("{{.ApiPrefix}}{{.ModuleName}}/{{.Router}}/export", { ...query }, `${filename}.xlsx`,{});
    },

}