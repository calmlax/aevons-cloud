import request, { download } from '@/utils/request'

export default {
    page: (query: any) => request.get('/v1/sys/user/page', { params: query }),
    list: (query: any) => request.get('/v1/sys/user/list', { params: query }),
    getById: (id: string) => request.get(`/v1/sys/user/${id}`),
    getRelations: (id: string) => request.get(`/v1/sys/user/${id}/relations`),
    add: (data: any) => request.post('/v1/sys/user', data),
    update: (id: string, data: any) => request.put(`/v1/sys/user/${id}`, data),
    updateStatus: (id: string | number, status: number) => request.put(`/v1/sys/user/${id}/status`, { status }),
    resetPassword: (id: string | number, password: string) => request.put(`/v1/sys/user/${id}/reset-password`, { password }),
    delete: (ids: any) => request.delete(`/v1/sys/user/${ids}`),
    download: (query: any, filename: string) => download('/v1/sys/user/export', { ...query }, `${filename}.xlsx`, {}),
}
