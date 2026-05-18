import request from '@/utils/request'

export default {
  list: () => request.get('/v1/sys/monitor/online'),
  forceLogout: (token: string) => request.delete(`/v1/sys/monitor/online/${token}`),
}
