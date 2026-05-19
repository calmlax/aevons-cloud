<template>
  <div class="home-page">
    <div class="home-card">
      <div class="logo-area">
        <div class="logo">第三方应用</div>
      </div>

      <template v-if="authStore.isLoggedIn">
        <a-result status="success" title="已登录">
          <template #subtitle>
            <p>Access Token 已获取，登录成功。</p>
            <a-typography-text code style="word-break: break-all; font-size: 12px">
              {{ authStore.tokenPair?.access_token }}
            </a-typography-text>
          </template>
          <template #extra>
            <a-space>
              <a-button type="primary" @click="handleLoginAgain">重新授权</a-button>
              <a-button status="danger" @click="handleLogout">退出登录</a-button>
            </a-space>
          </template>
        </a-result>
      </template>

      <template v-else>
        <div class="intro">
          <h2>欢迎使用第三方应用</h2>
          <p>本应用通过 OAuth2 授权码模式接入授权中心，点击下方按钮开始登录。</p>
        </div>
        <a-button type="primary" long size="large" :loading="loading" @click="handleLogin">
          使用授权中心登录
        </a-button>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useAuthStore } from '../stores/auth'
import { exchangeCode } from '../api/auth'

const authStore = useAuthStore()
const loading = ref(false)

const AUTHORIZE_CENTER = import.meta.env.VITE_AUTHORIZE_CENTER as string
const REDIRECT_URI = import.meta.env.VITE_REDIRECT_URI as string

let authWindow: Window | null = null
let pollTimer: ReturnType<typeof setInterval> | null = null

function stopLoading() {
  loading.value = false
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  authWindow = null
}

function handleLogin() {
  const url = new URL(AUTHORIZE_CENTER)
  url.searchParams.set('client_id', authStore.clientId)
  url.searchParams.set('redirect_uri', REDIRECT_URI)
  url.searchParams.set('response_type', 'code')

  loading.value = true
  authWindow = window.open(url.toString(), 'oauth2_authorize', 'width=520,height=640')

  if (!authWindow) {
    // Popup blocked — fall back to redirect
    window.location.href = url.toString()
    return
  }

  // Poll for popup close (covers manual close and deny redirect)
  pollTimer = setInterval(() => {
    if (authWindow?.closed) {
      stopLoading()
    }
  }, 500)
}

async function handleMessage(event: MessageEvent) {
  if (event.origin !== window.location.origin) return
  if (event.data?.type !== 'oauth2_callback') return

  const { code, error } = event.data as { type: string; code?: string; error?: string; state?: string }

  // Deny or error from callback page
  if (error || !code) {
    stopLoading()
    if (error === 'access_denied') {
      Message.warning('用户拒绝了授权')
    } else {
      Message.error(error || '未收到授权码')
    }
    return
  }

  try {
    const tokenPair = await exchangeCode(authStore.clientId, authStore.clientSecret, code)
    authStore.setTokenPair(tokenPair)
    Message.success('登录成功')
  } catch (err: unknown) {
    Message.error(extractMsg(err))
  } finally {
    stopLoading()
  }
}

function handleLogout() {
  authStore.clearAuth()
  Message.info('已退出登录')
}

function handleLoginAgain() {
  authStore.clearAuth()
  handleLogin()
}

onMounted(() => window.addEventListener('message', handleMessage))
onUnmounted(() => {
  window.removeEventListener('message', handleMessage)
  if (pollTimer) clearInterval(pollTimer)
})

function extractMsg(err: unknown): string {
  if (err && typeof err === 'object') {
    const e = err as Record<string, unknown>
    if (e.response && typeof e.response === 'object') {
      const res = e.response as Record<string, unknown>
      if (res.data && typeof res.data === 'object') {
        const d = res.data as Record<string, unknown>
        if (typeof d.message === 'string') return d.message
      }
    }
    if (typeof e.message === 'string') return e.message
  }
  return '登录失败，请重试'
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}
.home-card {
  background: #fff;
  border-radius: 8px;
  padding: 48px 40px;
  width: 440px;
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.1);
}
.logo-area {
  text-align: center;
  margin-bottom: 32px;
}
.logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: #165dff;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  line-height: 1.2;
  text-align: center;
}
.intro {
  text-align: center;
  margin-bottom: 28px;
}
.intro h2 {
  margin: 0 0 8px;
  font-size: 20px;
  color: #1d2129;
}
.intro p {
  margin: 0;
  color: #86909c;
  font-size: 14px;
  line-height: 1.6;
}
</style>
