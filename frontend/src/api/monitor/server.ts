import request from '@/utils/request'

export default {
  getInfo: () => request.get('/v1/sys/monitor/server'),
}
