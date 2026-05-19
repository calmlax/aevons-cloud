import request from '@/utils/request'

export default {
  list: () => request.get('/sys/v1/monitor/cache'),
  detail: (key: string) => request.get('/sys/v1/monitor/cache/detail', { params: { key } }),
  delete: (keys: string[]) => request.delete('/sys/v1/monitor/cache', { data: { keys } }),
  deleteByPrefix: (prefix: string) => request.delete('/sys/v1/monitor/cache/prefix', { params: { prefix } }),
}
