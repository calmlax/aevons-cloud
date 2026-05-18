import request from '@/utils/request'

export default {
  list: () => request.get('/v1/sys/monitor/cache'),
  detail: (key: string) => request.get('/v1/sys/monitor/cache/detail', { params: { key } }),
  delete: (keys: string[]) => request.delete('/v1/sys/monitor/cache', { data: { keys } }),
  deleteByPrefix: (prefix: string) => request.delete('/v1/sys/monitor/cache/prefix', { params: { prefix } }),
}
