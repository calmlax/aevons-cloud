import request from '@/utils/request'

export default {
  getInfo: () => request.get('/sys/v1/monitor/server'),
}
