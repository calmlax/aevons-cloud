<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/store/modules/auth';
import { IconDashboard, IconThunderbolt, IconSafe, IconCode } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n';

const router = useRouter();
const authStore = useAuthStore();
const { t } = useI18n();

const handleActionBtn = (path: string) => {
  const routeUrl = router.resolve({ path });
  window.open(routeUrl.href, '_blank');
};

const features = computed(() => [
  {
    icon: IconThunderbolt,
    title: t('portal.features.items.0.title'),
    desc: t('portal.features.items.0.desc')
  },
  {
    icon: IconSafe,
    title: t('portal.features.items.1.title'),
    desc: t('portal.features.items.1.desc')
  },
  {
    icon: IconDashboard,
    title: t('portal.features.items.2.title'),
    desc: t('portal.features.items.2.desc')
  },
  {
    icon: IconCode,
    title: t('portal.features.items.3.title'),
    desc: t('portal.features.items.3.desc')
  }
]);
</script>

<template>
  <div class="home-wrapper">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-bg">
        <div class="blob blob-1"></div>
        <div class="blob blob-2"></div>
      </div>
      <div class="hero-content">
        <h1 class="hero-title">
          {{ t('portal.hero.title1') }}<br />
          <span class="gradient-text">{{ t('portal.hero.title2') }}</span>
        </h1>
        <p class="hero-subtitle">
          {{ t('portal.hero.subtitle') }}
        </p>
        <div class="hero-actions">
          <a-button type="primary" size="large" class="btn-primary" @click="router.push('/products')">
            {{ t('portal.hero.exploreBtn') }}
          </a-button>
          <a-button size="large" class="btn-secondary" @click="handleActionBtn(authStore.isAuthenticated ? '/dashboard' : '/login')">
            {{ authStore.isAuthenticated ? (t('portal.home') === 'Home' ? 'Console' : '控制台') : (t('portal.home') === 'Home' ? 'Login' : '登录控制台') }}
          </a-button>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features-section">
      <div class="section-header">
        <h2 class="section-title">{{ t('portal.features.title') }}</h2>
        <p class="section-desc">{{ t('portal.features.desc') }}</p>
      </div>
      
      <div class="features-grid">
        <div v-for="(feat, idx) in features" :key="idx" class="feature-card">
          <div class="feat-icon-wrapper">
            <component :is="feat.icon" class="feat-icon" />
          </div>
          <h3 class="feat-title">{{ feat.title }}</h3>
          <p class="feat-desc">{{ feat.desc }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* 全局设定 */
.home-wrapper {
  background-color: var(--color-bg-1);
}

/* --- Hero Section --- */
.hero-section {
  position: relative;
  min-height: 85vh;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  overflow: hidden;
  padding: 0 24px;
}

/* 背景光影矩阵效果 */
.hero-bg {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 0;
  overflow: hidden;
  opacity: 0.15;
}
:global(body[arco-theme="dark"]) .hero-bg {
  opacity: 0.25;
}
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
}
.blob-1 {
  width: 600px;
  height: 600px;
  background: rgb(var(--primary-6));
  top: -150px;
  left: -100px;
  animation: float 15s ease-in-out infinite alternate;
}
.blob-2 {
  width: 500px;
  height: 500px;
  background: var(--color-brand-light, #b3e5fc);
  bottom: -50px;
  right: -100px;
  animation: float 10s ease-in-out infinite alternate-reverse;
}

@keyframes float {
  0% { transform: translateY(0) scale(1); }
  100% { transform: translateY(-50px) scale(1.1); }
}

.hero-content {
  position: relative;
  z-index: 10;
  max-width: 900px;
  animation: fadeInUp 1s cubic-bezier(0.16, 1, 0.3, 1);
}

.hero-title {
  font-size: 72px;
  font-weight: 900;
  line-height: 1.1;
  letter-spacing: -2px;
  margin-bottom: 32px;
  color: var(--color-text-1);
}

.gradient-text {
  background: linear-gradient(135deg, rgb(var(--primary-6)) 0%, #00b4db 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.hero-subtitle {
  font-size: 22px;
  color: var(--color-text-3);
  margin: 0 auto 48px;
  max-width: 700px;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  justify-content: center;
  gap: 24px;
}

.btn-primary {
  height: 56px;
  padding: 0 40px;
  font-size: 18px;
  border-radius: 28px;
  box-shadow: 0 8px 24px rgba(var(--primary-6), 0.3);
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.btn-primary:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 32px rgba(var(--primary-6), 0.4);
}

.btn-secondary {
  height: 56px;
  padding: 0 40px;
  font-size: 18px;
  border-radius: 28px;
  background: var(--color-fill-2);
  color: var(--color-text-1);
  border: none;
  transition: all 0.3s ease;
}
.btn-secondary:hover {
  background: var(--color-fill-3);
  transform: translateY(-4px);
}

/* --- Features Section --- */
.features-section {
  padding: 120px 24px;
  max-width: 1248px;
  margin: 0 auto;
}

.section-header {
  text-align: center;
  margin-bottom: 80px;
}
.section-title {
  font-size: 40px;
  font-weight: 800;
  color: var(--color-text-1);
  margin-bottom: 16px;
  letter-spacing: -1px;
}
.section-desc {
  font-size: 18px;
  color: var(--color-text-3);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 32px;
}

.feature-card {
  background: var(--color-bg-2);
  border-radius: 24px;
  padding: 40px 32px;
  border: 1px solid var(--color-border-1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
  overflow: hidden;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 4px;
  background: linear-gradient(90deg, rgb(var(--primary-6)), transparent);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.4s ease;
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.08);
  border-color: var(--color-primary-light-3);
}

.feature-card:hover::before {
  transform: scaleX(1);
}

.feat-icon-wrapper {
  width: 56px;
  height: 56px;
  background: var(--color-primary-light-1);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  transition: all 0.3s ease;
}
.feature-card:hover .feat-icon-wrapper {
  background: rgb(var(--primary-6));
}

.feat-icon {
  font-size: 28px;
  color: rgb(var(--primary-6));
  transition: color 0.3s ease;
}
.feature-card:hover .feat-icon {
  color: #fff;
}

.feat-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--color-text-1);
  margin-bottom: 16px;
}

.feat-desc {
  font-size: 15px;
  color: var(--color-text-3);
  line-height: 1.6;
}

/* 动画定义 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式支持 */
@media (max-width: 768px) {
  .hero-title {
    font-size: 48px;
  }
  .hero-subtitle {
    font-size: 18px;
  }
  .hero-actions {
    flex-direction: column;
    padding: 0 24px;
  }
  .features-section {
    padding: 80px 24px;
  }
  .section-title {
    font-size: 32px;
  }
}
</style>
