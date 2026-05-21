import request, { requestDownload, requestUpload } from '@/utils/request'
import { saveAs } from 'file-saver'

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
    exportFile: async (query: any, filename: string) => {
        const blob = await requestDownload('/sys/v1/user/export', { ...query })
        saveAs(new Blob([blob]), `${filename}.xlsx`)
    },
    downloadTemplate: async (filename: string) => {
        const blob = await requestDownload('/sys/v1/user/import/template')
        saveAs(new Blob([blob]), `${filename}.xlsx`)
    },
    importExcel: (formData: FormData) => requestUpload('/sys/v1/user/import', formData),
}
