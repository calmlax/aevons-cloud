import request from '@/utils/request'

export default {
  page:   (query: any) => request.get('/v1/log/login/log/page', { params: query }),
  delete: (ids: any[]) => request.delete(`/v1/log/login/log/${ids.join(',')}`),
  clear:  () => request.delete('/v1/log/login/log'),
}
