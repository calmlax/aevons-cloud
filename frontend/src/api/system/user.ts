import request, { download } from '@/utils/request'

export default {
    page: (query: any) => request.get('/sys/v1/user/page', { params: query }),
    list: (query: any) => request.get('/sys/v1/user/list', { params: query }),
    getById: (id: string) => request.get(`/sys/v1/user/${id}`),
    getRelations: (id: string) => request.get(`/sys/v1/user/${id}/relations`),
    add: (data: any) => request.post('/sys/v1/user', data),
    update: (id: string, data: any) => request.put(`/sys/v1/user/${id}`, data),
    updateStatus: (id: string | number, status: number) => request.put(`/sys/v1/user/${id}/status`, { status }),
    resetPassword: (id: string | number, password: string) => request.put(`/sys/v1/user/${id}/reset-password`, { password }),
    delete: (ids: any) => request.delete(`/sys/v1/user/${ids}`),
    download: (query: any, filename: string) => download('/sys/v1/user/export', { ...query }, `${filename}.xlsx`, {}),
}
