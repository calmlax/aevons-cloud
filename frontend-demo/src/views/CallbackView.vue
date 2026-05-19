<template>
  <div class="callback-container">
    <a-spin v-if="processing" tip="正在处理授权回调..." />

    <a-result
      v-else-if="success"
      status="success"
      title="授权成功"
      sub-title="授权码已传回，窗口即将关闭..."
    />

    <a-result v-else status="error" title="授权失败" :sub-title="errorMsg">
      <template #extra>
        <a-button type="primary" @click="closeWindow">关闭窗口</a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { extractCallbackParams } from '../utils/auth'

const processing = ref(true)
const success = ref(false)
const errorMsg = ref('')

const closeWindow = () => window.close()

onMounted(() => {
  const { code, state } = extractCallbackParams(window.location.href)

  // Check for OAuth2 error response (e.g. user denied)
  const urlParams = new URLSearchParams(window.location.search)
  const oauthError = urlParams.get('error')
  if (oauthError) {
    processing.value = false
    errorMsg.value = oauthError === 'access_denied' ? '用户拒绝了授权' : oauthError
    // Notify opener so it can re-enable the login button
    if (window.opener) {
      window.opener.postMessage(
        { type: 'oauth2_callback', error: oauthError },
        window.location.origin,
      )
      setTimeout(() => window.close(), 1500)
    }
    return
  }

  if (!code) {
    processing.value = false
    errorMsg.value = '回调 URL 中缺少 code 参数'
    return
  }

  // Send code back to the opener window via postMessage
  if (window.opener) {
    window.opener.postMessage(
      { type: 'oauth2_callback', code, state },
      window.location.origin,
    )
    success.value = true
    processing.value = false
    setTimeout(() => window.close(), 1500)
  } else {
    // Opened directly (not as popup) — redirect home with code in query
    // The home page can pick it up if needed, but normally this shouldn't happen
    processing.value = false
    errorMsg.value = '无法找到父窗口，请关闭此页面后重试'
  }
})
</script>

<style scoped>
.callback-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}
</style>
