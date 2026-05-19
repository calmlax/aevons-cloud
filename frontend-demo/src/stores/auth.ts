import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { TokenPair } from '../types/auth'

export const useAuthStore = defineStore('auth', () => {
  const stored = localStorage.getItem('third_part_token')
  const tokenPair = ref<TokenPair | null>(stored ? JSON.parse(stored) : null)

  const clientId = ref<string>(import.meta.env.VITE_OAUTH_CLIENT_ID || '')
  const clientSecret = ref<string>(import.meta.env.VITE_OAUTH_CLIENT_SECRET || '')

  watch(
    tokenPair,
    (val) => {
      if (val) localStorage.setItem('third_part_token', JSON.stringify(val))
      else localStorage.removeItem('third_part_token')
    },
    { deep: true },
  )

  const isLoggedIn = computed(() => !!tokenPair.value?.access_token)

  function setTokenPair(pair: TokenPair) {
    tokenPair.value = pair
  }

  function clearAuth() {
    tokenPair.value = null
  }

  return { tokenPair, clientId, clientSecret, isLoggedIn, setTokenPair, clearAuth }
})
