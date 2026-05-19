import request from '@/utils/request'

export default {
  // ── 任务 CRUD ──────────────────────────────────────────
  page: (query: any) => request.get('/job/v1/page', { params: query }),
  getById: (id: string) => request.get(`/job/v1/${id}`),
  add: (data: any) => request.post('/job/v1', data),
  update: (id: string, data: any) => request.put(`/job/v1/${id}`, data),
  delete: (ids: string) => request.delete(`/job/v1/${ids}`),

  // ── 调度操作 ───────────────────────────────────────────
  trigger: (id: string) => request.post(`/job/v1/${id}/trigger`),
  changeStatus: (id: string, status: number) => request.put(`/job/v1/${id}/status`, { status }),

  // ── 执行日志 ───────────────────────────────────────────
  logPage: (query: any) => request.get('/job/v1/log/page', { params: query }),
  logDelete: (ids: string) => request.delete(`/job/v1/log/${ids}`),
}
