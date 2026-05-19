import request from '@/utils/request'

export default {
  page:   (query: any) => request.get('/log/v1/login/log/page', { params: query }),
  delete: (ids: any[]) => request.delete(`/log/v1/login/log/${ids.join(',')}`),
  clear:  () => request.delete('/log/v1/login/log'),
}
