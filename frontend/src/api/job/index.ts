import request from '@/utils/request'

export default {
  // ── 任务 CRUD ──────────────────────────────────────────
  page: (query: any) => request.get('/v1/job/page', { params: query }),
  getById: (id: string) => request.get(`/v1/job/${id}`),
  add: (data: any) => request.post('/v1/job', data),
  update: (id: string, data: any) => request.put(`/v1/job/${id}`, data),
  delete: (ids: string) => request.delete(`/v1/job/${ids}`),

  // ── 调度操作 ───────────────────────────────────────────
  trigger: (id: string) => request.post(`/v1/job/${id}/trigger`),
  changeStatus: (id: string, status: number) => request.put(`/v1/job/${id}/status`, { status }),

  // ── 执行日志 ───────────────────────────────────────────
  logPage: (query: any) => request.get('/v1/job/log/page', { params: query }),
  logDelete: (ids: string) => request.delete(`/v1/job/log/${ids}`),
}
