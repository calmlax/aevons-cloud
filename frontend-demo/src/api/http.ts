import axios from 'axios'
import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios'

const http: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
})

http.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  const { useAuthStore } = await import('../stores/auth')
  const authStore = useAuthStore()
  const hasAuth = Object.keys(config.headers ?? {}).some(
    (k) => k.toLowerCase() === 'authorization',
  )
  if (authStore.tokenPair?.access_token && !hasAuth) {
    config.headers['Authorization'] = `Bearer ${authStore.tokenPair.access_token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    // Unwrap unified envelope { code, message, data }
    if (response.data && typeof response.data === 'object' && 'data' in response.data) {
      response.data = response.data.data
    }
    return response
  },
  (error) => Promise.reject(error),
)

export default http
