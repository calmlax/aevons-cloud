import request from '@/utils/request'

export default {
  page:   (query: any) => request.get('/v1/log/oper/log/page', { params: query }),
  delete: (ids: any[]) => request.delete(`/v1/log/oper/log/${ids.join(',')}`),
  clear:  () => request.delete('/v1/log/oper/log'),
}
