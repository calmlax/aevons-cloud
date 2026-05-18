<script setup lang="ts">
import { reactive, ref, watchEffect, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import JSEncrypt from 'jsencrypt';

import { useAuthStore } from '@/store/modules/auth';
import { getPublicKey } from '@/api/auth';
import {
  passkeyLoginBegin,
  passkeyLoginFinish,
  parseRequestOptions,
} from '@/api/passkey';

const router = useRouter();
const route = useRoute();
const { t } = useI18n();
const authStore = useAuthStore();
const { login } = authStore;
const submitting = ref(false);
const activeTab = ref('password');
const sendingCode = ref(false);
const countdown = ref(0);
const passkeyLoading = ref(false);
let timer: ReturnType<typeof setInterval> | null = null;
const form = reactive({
  email: '',
  password: '',
  code: '',
});

const startCountdown = () => {
  countdown.value = 60;
  timer = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      clearInterval(timer!);
      timer = null;
    }
  }, 1000);
};

const handleSendCode = async () => {
  if (!form.email) {
    Message.warning(t('auth.loginIncomplete'));
    return;
  }
  sendingCode.value = true;
  try {
    await authStore.sendRealEmailCode(form.email, 'login');
    Message.success(t('auth.codeSentSuccess'));
    startCountdown();
  } catch (error: any) {
    if (!error?.config) {
      Message.error(error?.message || t('auth.codeSendFailed'));
    }
  } finally {
    sendingCode.value = false;
  }
};

onUnmounted(() => {
  if (timer) clearInterval(timer);
});

watchEffect(() => {
  if (typeof route.query.email === 'string') {
    form.email = route.query.email;
  }
});

const resolveRedirect = () => {
  return typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard';
};

const onSubmit = async () => {
  if (!form.email) {
    Message.warning(t('auth.loginIncomplete'));
    return;
  }
  
  if (activeTab.value === 'password' && !form.password) {
    Message.warning(t('auth.loginIncomplete'));
    return;
  }

  if (activeTab.value === 'email' && !form.code) {
    Message.warning(t('auth.loginIncomplete'));
    return;
  }

  submitting.value = true;

  try {
    if (activeTab.value === 'password') {
      const pubKeyResp = await getPublicKey();
      const encryptor = new JSEncrypt();
      encryptor.setPublicKey(pubKeyResp.public_key);
      const encryptedPassword = encryptor.encrypt(form.password);

      if (!encryptedPassword) {
        throw new Error('加密失败，请重试');
      }

      await login({
        grant_type: 'password',
        email: form.email,
        password: '',
        encryptedPassword,
        key_id: pubKeyResp.key_id,
      });
    } else {
      await login({
        grant_type: 'email',
        email: form.email,
        password: '',
        code: form.code,
      });
    }

    Message.success(t('auth.loginSuccess'));
    await router.replace(resolveRedirect());
  } catch (error: any) {
    if (!error?.config) {
      Message.error(error?.message || '操作失败');
    }
  } finally {
    submitting.value = false;
  }
};

// ── Passkey 登录 ──────────────────────────────────────────────────────────────

const isPasskeySupported = typeof window !== 'undefined' &&
  typeof window.PublicKeyCredential !== 'undefined';

const onPasskeyLogin = async () => {
  if (!isPasskeySupported) {
    Message.warning(t('auth.passkeyNotSupported'));
    return;
  }

  passkeyLoading.value = true;
  try {
    // 1. 获取 challenge
    const { options: optionsJSON, session_key } = await passkeyLoginBegin();

    // 2. 调用浏览器 API
    const requestOptions = parseRequestOptions(optionsJSON);
    const credential = await navigator.credentials.get({
      publicKey: requestOptions,
    }) as PublicKeyCredential | null;

    if (!credential) {
      Message.warning(t('auth.passkeyCancelled'));
      return;
    }

    // 3. 发送给后端验证，获取 token
    const tokenPair = await passkeyLoginFinish(session_key, credential);

    // 4. 存储 token，获取用户信息
    authStore.setToken(tokenPair.access_token);
    await authStore.refreshUserInfo();

    Message.success(t('auth.loginSuccess'));
    await router.replace(resolveRedirect());
  } catch (error: any) {
    if (error?.name === 'NotAllowedError') {
      // 用户取消或超时，不弹错误
      return;
    }
    if (!error?.config) {
      Message.error(error?.message || t('auth.passkeyFailed'));
    }
  } finally {
    passkeyLoading.value = false;
  }
};
</script>

<template>
  <div class="auth-form-card panel-card">
    <div class="auth-form-header">
      <p class="eyebrow">{{ t('auth.loginEyebrow') }}</p>
      <h2>{{ t('auth.loginTitle') }}</h2>
      <p>{{ t('auth.loginDescription') }}</p>
    </div>

    <a-tabs v-model:active-key="activeTab" class="login-tabs">
      <a-tab-pane key="password" :title="t('auth.passwordLogin')" />
      <a-tab-pane key="email" :title="t('auth.emailCodeLogin')" />
      <a-tab-pane key="passkey" title="Passkey" />
    </a-tabs>

    <!-- 密码 / 邮箱验证码登录 -->
    <a-form v-if="activeTab !== 'passkey'" layout="vertical">
      <a-form-item field="email" :label="t('auth.email')">
        <a-input v-model="form.email" size="large" allow-clear :placeholder="t('auth.emailPlaceholder')" @keyup.enter="onSubmit" />
      </a-form-item>

      <a-form-item v-if="activeTab === 'password'" field="password" :label="t('auth.password')">
        <a-input-password
          v-model="form.password"
          size="large"
          allow-clear
          :placeholder="t('auth.passwordPlaceholder')"
          @keyup.enter="onSubmit"
        />
      </a-form-item>

      <a-form-item v-if="activeTab === 'email'" field="code" :label="t('auth.emailCode')">
        <div style="display: flex; gap: 8px; width: 100%;">
          <a-input v-model="form.code" size="large" allow-clear :placeholder="t('auth.emailCodePlaceholder')" style="flex: 1" @keyup.enter="onSubmit" />
          <a-button size="large" :disabled="countdown > 0 || sendingCode" :loading="sendingCode" @click="handleSendCode">
            {{ countdown > 0 ? t('auth.resendCodeIn', { seconds: countdown }) : t('auth.sendCode') }}
          </a-button>
        </div>
      </a-form-item>

      <a-button type="primary" long size="large" :loading="submitting" @click="onSubmit">
        {{ t('auth.loginAction') }}
      </a-button>
    </a-form>

    <!-- Passkey 登录 -->
    <div v-else class="passkey-panel">
      <div class="passkey-icon">🔑</div>
      <p class="passkey-desc">{{ t('auth.passkeyDesc') }}</p>
      <a-button
        type="primary"
        long
        size="large"
        :loading="passkeyLoading"
        :disabled="!isPasskeySupported"
        @click="onPasskeyLogin"
      >
        {{ isPasskeySupported ? t('auth.passkeyLogin') : t('auth.passkeyNotSupported') }}
      </a-button>
    </div>

    <div class="auth-links-row">
      <router-link to="/auth/register">{{ t('auth.goRegister') }}</router-link>
      <router-link to="/auth/forgot-password">{{ t('auth.goForgotPassword') }}</router-link>
    </div>
  </div>
</template>

<style scoped>
.passkey-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 24px 0 8px;
}
.passkey-icon {
  font-size: 48px;
  line-height: 1;
}
.passkey-desc {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-3);
  text-align: center;
  max-width: 280px;
}
</style>