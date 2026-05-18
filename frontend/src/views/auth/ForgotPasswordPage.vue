<script setup lang="ts">
import { reactive, ref, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';

import { useAuthStore } from '@/store/modules/auth';
const authStore = useAuthStore();

const router = useRouter();
const { t } = useI18n();
const submitting = ref(false);
const sendingCode = ref(false);
const countdown = ref(0);
const form = reactive({
  email: '',
  code: '',
  password: '',
  confirmPassword: '',
});

let timer: number | undefined;

const startCountdown = () => {
  countdown.value = 60;
  timer = window.setInterval(() => {
    countdown.value -= 1;
    if (countdown.value <= 0 && timer) {
      window.clearInterval(timer);
      timer = undefined;
    }
  }, 1000);
};

const onSendCode = async () => {
  if (!form.email) {
    Message.warning(t('auth.emailRequired'));
    return;
  }

  sendingCode.value = true;

  try {
    const code = await authStore.sendEmailCode(form.email, 'reset-password');
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

const onSubmit = async () => {
  if (!form.email || !form.code || !form.password || !form.confirmPassword) {
    Message.warning(t('auth.resetIncomplete'));
    return;
  }

  if (form.password !== form.confirmPassword) {
    Message.error(t('auth.passwordMismatch'));
    return;
  }

  submitting.value = true;

  try {
    await authStore.resetPassword({ email: form.email, password: form.password, code: form.code });
    await router.push({
      name: 'auth-result',
      params: { scene: 'reset-password' },
      query: { email: form.email },
    });
  } catch (error: any) {
    if (!error?.config) {
      Message.error(error?.message || t('auth.resetFailed'));
    }
  } finally {
    submitting.value = false;
  }
};

onUnmounted(() => {
  if (timer) {
    window.clearInterval(timer);
  }
});
</script>

<template>
  <div class="auth-form-card panel-card">
    <div class="auth-form-header">
      <p class="eyebrow">{{ t('auth.resetEyebrow') }}</p>
      <h2>{{ t('auth.resetTitle') }}</h2>
      <p>{{ t('auth.resetDescription') }}</p>
    </div>

    <a-form layout="vertical">
      <a-form-item field="email" :label="t('auth.email')">
        <a-input v-model="form.email" size="large" allow-clear :placeholder="t('auth.emailPlaceholder')" />
      </a-form-item>

      <a-form-item field="code" :label="t('auth.emailCode')">
        <a-input-group compact>
          <a-input v-model="form.code" size="large" allow-clear :placeholder="t('auth.emailCodePlaceholder')" />
          <a-button size="large" :disabled="countdown > 0" :loading="sendingCode" @click="onSendCode">
            {{ countdown > 0 ? t('auth.resendCodeIn', { seconds: countdown }) : t('auth.sendCode') }}
          </a-button>
        </a-input-group>
      </a-form-item>

      <a-form-item field="password" :label="t('auth.newPassword')">
        <a-input-password
          v-model="form.password"
          size="large"
          allow-clear
          :placeholder="t('auth.newPasswordPlaceholder')"
        />
      </a-form-item>

      <a-form-item field="confirmPassword" :label="t('auth.confirmPassword')">
        <a-input-password
          v-model="form.confirmPassword"
          size="large"
          allow-clear
          :placeholder="t('auth.confirmPasswordPlaceholder')"
        />
      </a-form-item>

      <a-button type="primary" long size="large" :loading="submitting" @click="onSubmit">
        {{ t('auth.resetAction') }}
      </a-button>
    </a-form>

    <div class="auth-links-row auth-links-row--single">
      <router-link to="/login">{{ t('auth.backToLogin') }}</router-link>
    </div>
  </div>
</template>