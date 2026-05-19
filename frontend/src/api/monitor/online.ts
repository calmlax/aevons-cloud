import request from '@/utils/request'

export default {
  list: () => request.get('/sys/v1/monitor/online'),
  forceLogout: (token: string) => request.delete(`/sys/v1/monitor/online/${token}`),
}
