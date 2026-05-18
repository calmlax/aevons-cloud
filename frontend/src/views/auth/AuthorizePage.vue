<template>
  <div class="authorize-page">
    <div class="authorize-card">

      <div v-if="loading" class="center">
        <a-spin tip="加载中..." />
      </div>

      <a-result v-else-if="error" status="error" :title="error">
        <template #extra>
          <a-button @click="goBack">返回</a-button>
        </template>
      </a-result>

      <!-- 仅保留授权确认步骤 -->
      <template v-else>
        <div class="app-info">
          <img v-if="info.logo_uri" :src="info.logo_uri" class="app-logo" />
          <div v-else class="app-logo-placeholder">{{ info.client_name?.charAt(0) }}</div>
          <h2 class="app-name">{{ info.client_name || info.client_id }}</h2>
          <p class="app-desc">该应用申请访问你的账号</p>
        </div>

        <div v-if="scopeList.length" class="scope-box">
          <p class="scope-title">申请以下权限：</p>
          <a-checkbox-group v-model="selectedScopes" direction="vertical" style="width:100%">
            <a-checkbox v-for="s in scopeList" :key="s" :value="s" style="margin-bottom:8px">
              <div class="scope-item">
                <span class="scope-name">{{ s }}</span>
                <span class="scope-desc">{{ scopeLabels[s] ?? '' }}</span>
              </div>
            </a-checkbox>
          </a-checkbox-group>
        </div>

        <div class="action-row">
          <a-button long status="danger" style="margin-bottom:8px" @click="handleDeny">拒绝</a-button>
          <a-button type="primary" long :loading="submitting" @click="handleApprove">同意授权</a-button>
        </div>
        <div class="back-link">
          <a-link @click="handleSwitchAccount">切换账号</a-link>
        </div>
      </template>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { getAuthorizeInfo, approveAuthorize } from '@/api/auth'
import { useAuthStore } from '@/store/modules/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loading = ref(true)
const error = ref('')
const submitting = ref(false)

// 当前授权用的 access token（已登录用户的 token）
const currentToken = ref('')

const info = ref({
  client_id: '',
  client_name: '',
  logo_uri: '',
  scope: '',
  redirect_uri: '',
  state: '',
  autoapprove: false,
})

const scopeList = computed(() =>
  info.value.scope ? info.value.scope.split(',').map((s) => s.trim()).filter(Boolean) : [],
)
const selectedScopes = ref<string[]>([])

const scopeLabels: Record<string, string> = {
  openid: '获取你的唯一标识',
  profile: '获取你的基本信息（昵称、头像等）',
  email: '获取你的邮箱地址',
  phone: '获取你的手机号',
}

const clientId = route.query.client_id as string ?? ''
const redirectUri = route.query.redirect_uri as string ?? ''

onMounted(async () => {
  if (!clientId) {
    error.value = '缺少 client_id 参数'
    loading.value = false
    return
  }

  // 获取 URL 参数检查登录状态。如果没有 Token 则强定向前去登录页。
  const accessToken = authStore.accessToken

  if (!accessToken) {
    // 平滑跳转到独立登录页并附带返回地址
    router.replace({ path: '/login', query: { redirect: route.fullPath } })
    return
  }

  currentToken.value = accessToken

  // 获取展示信息
  try {
    const data = await getAuthorizeInfo(clientId, redirectUri)
    info.value = data
    selectedScopes.value = data.scope
      ? data.scope.split(',').map((s) => s.trim()).filter(Boolean)
      : []

    // autoapprove=1：静默授权直接跳转
    if (data.autoapprove) {
      await doApprove(currentToken.value, selectedScopes.value)
    }
  } catch (err: unknown) {
    error.value = extractMsg(err)
  } finally {
    loading.value = false
  }
})

async function handleApprove() {
  if (!currentToken.value) {
    Message.error('长时间未操作，请重新登录')
    return handleSwitchAccount()
  }
  submitting.value = true
  try {
    await doApprove(currentToken.value, selectedScopes.value)
  } catch (err: unknown) {
    Message.error(extractMsg(err))
  } finally {
    submitting.value = false
  }
}

async function doApprove(token: string, scopes: string[]) {
  const res = await approveAuthorize({
    state: info.value.state,
    access_token: token,
    scopes,
  })
  window.location.href = res.redirect_uri
}

function handleDeny() {
  const base = info.value.redirect_uri || redirectUri
  if (base) {
    const sep = base.includes('?') ? '&' : '?'
    window.location.href = base + sep + 'error=access_denied&state=' + info.value.state
  } else {
    goBack()
  }
}

function handleSwitchAccount() {
  authStore.logout() // 清除当前凭证
  router.replace({ path: '/login', query: { redirect: route.fullPath } })
}

function goBack() {
  window.history.back()
}

function extractMsg(err: unknown): string {
  if (err && typeof err === 'object') {
    const e = err as Record<string, unknown>
    if (e.response && typeof e.response === 'object') {
      const res = e.response as Record<string, unknown>
      if (res.data && typeof res.data === 'object') {
        const data = res.data as Record<string, unknown>
        if (typeof data.message === 'string') return data.message
      }
    }
    if (typeof e.message === 'string') return e.message
  }
  return '操作失败'
}
</script>

<style scoped>
.authorize-page {
  min-height: 100vh;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}
.authorize-card {
  background: #fff;
  border-radius: 8px;
  padding: 40px 36px;
  width: 400px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.center {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}
.app-info {
  text-align: center;
  margin-bottom: 24px;
}
.app-logo {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  object-fit: cover;
}
.app-logo-placeholder {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  background: #165dff;
  color: #fff;
  font-size: 28px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}
.app-name {
  margin: 12px 0 4px;
  font-size: 18px;
  font-weight: 600;
}
.app-desc {
  color: #666;
  font-size: 13px;
  margin: 0;
}
.scope-box {
  background: #f7f8fa;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 20px;
}
.scope-title {
  font-size: 13px;
  color: #666;
  margin: 0 0 12px;
}
.scope-item {
  display: flex;
  flex-direction: column;
}
.scope-name {
  font-weight: 500;
  font-size: 14px;
}
.scope-desc {
  font-size: 12px;
  color: #999;
}
.action-row {
  display: flex;
  flex-direction: column;
}
.back-link {
  text-align: center;
  margin-top: 12px;
  font-size: 13px;
}
</style>
