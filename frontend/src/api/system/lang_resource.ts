import request, { download } from '@/utils/request'

/**
 * 语言资源 Model
 *
 * @author 
 * @date 2026-04-19
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
        return request.get('/sys/v1/lang/resource/page',{params: query})
    },

    /**
     * 查询列表
     * @param query 查询参数
     */
    list: (query: any) => {
        return request.get('/sys/v1/lang/resource/list',{params: query})
    },

    /**
     * 查询详细
     * @param query 查询参数
     */
    getById: (id: string) => {
        return request.get(`/sys/v1/lang/resource/${id}`)
    },

    /**
     * 新增
     * @param data 新增数据
     */
    add: (data: any) => {
        return request.post('/sys/v1/lang/resource',data)
    },

    /**
     * 修改
     * @param data 修改数据
     */
    update: (id: string,data: any) => {
        return request.put(`/sys/v1/lang/resource/${id}`,data)
    },

    /**
     * 删除
     * @param ids 删除Ids
     */
    delete: (ids: any) => {
        return request.delete(`/sys/v1/lang/resource/${ids}`)
    },
    
    /** 获取命名空间下的所有 resourceKey */
    getKeysByNamespace: (namespace: string) =>
        request.get('/sys/v1/lang/resource/keys', { params: { namespace } }),

    /** 去重分页查询 resourceKey（支持 key/content 搜索） */
    pageKeys: (params: { namespace: string; resourceKey?: string; content?: string; pageNum?: number; pageSize?: number }) =>
        request.get('/sys/v1/lang/resource/keys/page', { params }),

    /** 获取某个 namespace+resourceKey 的所有语言翻译 */
    getTranslations: (namespace: string, resourceKey: string) =>
        request.get('/sys/v1/lang/resource/translations', { params: { namespace, resourceKey } }),

    /** 批量保存翻译（upsert） */
    saveTranslations: (data: { namespace: string; resourceKey: string; items: { langCode: string; content: string }[] }) =>
        request.post('/sys/v1/lang/resource/save-translations', data),

    /**
     * 导出Excel
     * @param filename Excel文件名称
     * @param query 查询参数
     */
    download: (query: any, filename: string) => {
        download("/sys/v1/lang/resource/export", { ...query }, `${filename}.xlsx`,{});
    },

}