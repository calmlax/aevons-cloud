<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
  <div class="about-wrapper">
    <div class="content-container">
      <div class="text-content">
        <h1 class="about-tagline">
          {{ t('portal.aboutPage.title1') }}<br/>
          <span>{{ t('portal.aboutPage.title2') }}</span>
        </h1>
        <div class="about-body">
          <p class="lead-text">
            {{ t('portal.aboutPage.lead') }}
          </p>
          <p>
            {{ t('portal.aboutPage.p1') }}
          </p>
          <p>
            {{ t('portal.aboutPage.p2') }}
          </p>
        </div>
        
        <div class="stats-row">
          <div class="stat-card">
            <h4 class="stat-num">98%</h4>
            <span class="stat-label">{{ t('portal.aboutPage.stats.code') }}</span>
          </div>
          <div class="stat-card">
            <h4 class="stat-num">10x</h4>
            <span class="stat-label">{{ t('portal.aboutPage.stats.time') }}</span>
          </div>
          <div class="stat-card">
            <h4 class="stat-num">∞</h4>
            <span class="stat-label">{{ t('portal.aboutPage.stats.scale') }}</span>
          </div>
        </div>
      </div>
      
      <div class="visual-content">
        <div class="glass-orb"></div>
        <div class="glass-pane main-pane"></div>
        <div class="glass-pane sub-pane"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.about-wrapper {
  padding: 100px 24px;
  background-color: var(--color-bg-1);
}

.content-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 80px;
}

.text-content {
  flex: 1;
  max-width: 600px;
  animation: slideInLeft 1s cubic-bezier(0.16, 1, 0.3, 1);
}

.about-tagline {
  font-size: 56px;
  font-weight: 900;
  line-height: 1.15;
  letter-spacing: -2px;
  color: var(--color-text-1);
  margin-bottom: 40px;
}
.about-tagline span {
  color: rgb(var(--primary-6));
}

.about-body {
  font-size: 16px;
  line-height: 1.8;
  color: var(--color-text-3);
  margin-bottom: 56px;
}
.lead-text {
  font-size: 20px;
  font-weight: 500;
  color: var(--color-text-2);
  margin-bottom: 24px;
}

.stats-row {
  display: flex;
  gap: 40px;
  border-top: 1px solid var(--color-border-2);
  padding-top: 32px;
}
.stat-card {
  display: flex;
  flex-direction: column;
}
.stat-num {
  font-size: 40px;
  font-weight: 800;
  margin: 0 0 8px 0;
  color: var(--color-text-1);
}
.stat-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-3);
  text-transform: uppercase;
}

.visual-content {
  flex: 1;
  position: relative;
  height: 600px;
  max-width: 500px;
  animation: slideInRight 1s cubic-bezier(0.16, 1, 0.3, 1) 0.2s backwards;
}

/* 玻璃态艺术品 */
.glass-orb {
  position: absolute;
  top: 10%; left: 0;
  width: 280px; height: 280px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgb(var(--primary-5)), #8E2DE2);
  filter: blur(40px);
  opacity: 0.6;
  animation: orbFloat 10s ease-in-out infinite alternate;
}
.glass-pane {
  position: absolute;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 24px;
  box-shadow: 0 24px 48px rgba(0,0,0,0.05);
}
:global(body[arco-theme="dark"]) .glass-pane {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
.main-pane {
  width: 320px; height: 440px;
  top: 50%; left: 50%;
  transform: translate(-50%, -50%) rotate(-5deg);
}
.sub-pane {
  width: 240px; height: 320px;
  top: 60%; left: 10%;
  transform: rotate(8deg);
  background: rgba(255, 255, 255, 0.2);
}

@keyframes orbFloat {
  0% { transform: translateY(0) scale(1); }
  100% { transform: translateY(-40px) scale(1.1); }
}

@keyframes slideInLeft {
  from { opacity: 0; transform: translateX(-40px); }
  to { opacity: 1; transform: translateX(0); }
}
@keyframes slideInRight {
  from { opacity: 0; transform: translateX(40px); }
  to { opacity: 1; transform: translateX(0); }
}

@media (max-width: 992px) {
  .content-container {
    flex-direction: column;
    text-align: center;
  }
  .stats-row {
    justify-content: center;
  }
  .visual-content {
    display: none;
  }
}
</style>
