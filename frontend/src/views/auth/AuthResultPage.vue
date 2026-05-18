<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

type ResultScene = 'register' | 'reset-password' | 'not-found';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();

const scene = computed<ResultScene>(() => {
  const rawScene = route.params.scene;
  if (rawScene === 'register' || rawScene === 'reset-password' || rawScene === 'not-found') {
    return rawScene;
  }

  return 'not-found';
});

const sceneConfig = computed(() => {
  if (scene.value === 'register') {
    return {
      status: 'success' as const,
      title: t('result.registerTitle'),
      subtitle: t('result.registerDescription', { email: route.query.email ?? '--' }),
      primaryLabel: t('auth.loginAction'),
      primaryAction: () => router.push({ path: '/login', query: { email: route.query.email ?? '' } }),
      secondaryLabel: t('result.backHome'),
      secondaryAction: () => router.push('/dashboard'),
    };
  }

  if (scene.value === 'reset-password') {
    return {
      status: 'success' as const,
      title: t('result.resetTitle'),
      subtitle: t('result.resetDescription', { email: route.query.email ?? '--' }),
      primaryLabel: t('auth.loginAction'),
      primaryAction: () => router.push('/login'),
      secondaryLabel: t('result.backHome'),
      secondaryAction: () => router.push('/dashboard'),
    };
  }

  return {
    status: '404' as const,
    title: t('result.notFoundTitle'),
    subtitle: t('result.notFoundDescription'),
    primaryLabel: t('result.backHome'),
    primaryAction: () => router.push('/dashboard'),
    secondaryLabel: t('auth.backToLogin'),
    secondaryAction: () => router.push('/login'),
  };
});
</script>

<template>
  <div class="status-page panel-card">
    <a-result :status="sceneConfig.status" :title="sceneConfig.title" :subtitle="sceneConfig.subtitle">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="sceneConfig.primaryAction()">{{ sceneConfig.primaryLabel }}</a-button>
          <a-button @click="sceneConfig.secondaryAction()">{{ sceneConfig.secondaryLabel }}</a-button>
        </a-space>
      </template>
    </a-result>
  </div>
</template>