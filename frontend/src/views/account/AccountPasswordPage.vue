<script setup lang="ts">
import { ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';

import { useAuthStore } from '@/store/modules/auth';
const authStore = useAuthStore();

const currentPassword = ref('');
const nextPassword = ref('');
const confirmPassword = ref('');
const { t } = useI18n();

const resetForm = () => {
  currentPassword.value = '';
  nextPassword.value = '';
  confirmPassword.value = '';
};

const onSubmit = async () => {
  if (!currentPassword.value || !nextPassword.value || !confirmPassword.value) {
    Message.warning(t('password.incomplete'));
    return;
  }

  if (nextPassword.value !== confirmPassword.value) {
    Message.error(t('password.mismatch'));
    return;
  }

  try {
    await authStore.updatePassword({
      currentPassword: currentPassword.value,
      nextPassword: nextPassword.value,
    });
    Message.success(t('password.success'));
    resetForm();
  } catch (error) {
    const messageKey = error instanceof Error ? error.message : 'password.currentIncorrect';
    Message.error(t(messageKey));
  }
};
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false" :title="t('password.title')">
      <a-form layout="vertical" style="max-width: 460px">
        <a-form-item field="currentPassword" :label="t('password.current')">
          <a-input-password v-model="currentPassword" :placeholder="t('password.currentPlaceholder')" allow-clear />
        </a-form-item>

        <a-form-item field="nextPassword" :label="t('password.next')">
          <a-input-password v-model="nextPassword" :placeholder="t('password.nextPlaceholder')" allow-clear />
        </a-form-item>

        <a-form-item field="confirmPassword" :label="t('password.confirm')">
          <a-input-password v-model="confirmPassword" :placeholder="t('password.confirmPlaceholder')" allow-clear />
        </a-form-item>

        <a-space size="small">
          <a-button type="primary" @click="onSubmit">{{ t('password.save') }}</a-button>
          <a-button @click="resetForm">{{ t('password.cancel') }}</a-button>
        </a-space>
      </a-form>
    </a-card>
  </div>
</template>