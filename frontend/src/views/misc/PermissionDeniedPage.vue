<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

import { useAuthStore } from '@/store/modules/auth';
const authStore = useAuthStore();

const router = useRouter();
const { t } = useI18n();
const countdown = ref(5);
const targetPath = computed(() => authStore.getAuthHomePath());

let timer: number | undefined;

const navigateHome = async () => {
  await router.replace(targetPath.value);
};

const goBack = () => {
  router.back();
};

onMounted(() => {
  timer = window.setInterval(() => {
    countdown.value -= 1;
    if (countdown.value <= 0) {
      if (timer) {
        window.clearInterval(timer);
      }
      navigateHome();
    }
  }, 1000);
});

onUnmounted(() => {
  if (timer) {
    window.clearInterval(timer);
  }
});
</script>

<template>
  <div class="permission-page-wrapper">
    <div class="status-page panel-card permission-page">
      <a-result status="403" :title="t('permission.title')" :subtitle="t('permission.description')">
        <template #extra>
          <a-space direction="vertical" size="medium" fill>
            <div class="permission-countdown">{{ t('permission.countdown', { seconds: countdown }) }}</div>
            <a-space>
              <a-button type="primary" @click="() => router.push(authStore.getAuthHomePath())">{{ t('permission.backHome') }}</a-button>
              <a-button @click="goBack">{{ t('permission.goBack') }}</a-button>
            </a-space>
          </a-space>
        </template>
      </a-result>
    </div>
  </div>
</template>

<style scoped>
.permission-page-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
}

.permission-page {
  /* Removing inherited min-height overrides to let content define size, or keep standard status-page height */
  width: min(100%, 680px);
  text-align: center;
}
</style>