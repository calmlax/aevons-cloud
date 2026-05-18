<script setup lang="ts">
import { computed } from 'vue';
import { IconLayers, IconStorage, IconRobot, IconMobile } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const products = computed(() => [
  {
    id: 1,
    icon: IconLayers,
    name: t('portal.productsPage.items.0.title'),
    category: t('portal.productsPage.items.0.badge'),
    description: t('portal.productsPage.items.0.desc')
  },
  {
    id: 2,
    icon: IconStorage,
    name: t('portal.productsPage.items.1.title'),
    category: t('portal.productsPage.items.1.badge'),
    description: t('portal.productsPage.items.1.desc')
  },
  {
    id: 3,
    icon: IconRobot,
    name: t('portal.productsPage.items.2.title'),
    category: t('portal.productsPage.items.2.badge'),
    description: t('portal.productsPage.items.2.desc')
  },
  {
    id: 4,
    icon: IconMobile,
    name: t('portal.productsPage.items.3.title'),
    category: t('portal.productsPage.items.3.badge'),
    description: t('portal.productsPage.items.3.desc')
  }
]);
</script>

<template>
  <div class="products-wrapper">
    <div class="page-header">
      <h1 class="title">{{ t('portal.productsPage.title1') }} <span class="highlight">{{ t('portal.productsPage.title2') }}</span></h1>
      <p class="subtitle">{{ t('portal.productsPage.subtitle') }}</p>
    </div>

    <div class="product-gallery">
      <div v-for="item in products" :key="item.id" class="product-item">
        <div class="product-icon-box">
          <component :is="item.icon" class="main-icon" />
        </div>
        <div class="product-content">
          <div class="meta-row">
            <span class="category-tag">{{ item.category }}</span>
          </div>
          <h2 class="product-name">{{ item.name }}</h2>
          <p class="product-desc">{{ item.description }}</p>
          <div class="product-action">
            <a class="learn-more">{{ t('portal.productsPage.learnMore') }} <span class="arrow">→</span></a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.products-wrapper {
  padding: 80px 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 80px;
  animation: fadeInDown 0.8s ease-out;
}

.title {
  font-size: 56px;
  font-weight: 900;
  color: var(--color-text-1);
  letter-spacing: -1.5px;
  margin-bottom: 24px;
}

.highlight {
  background: linear-gradient(135deg, rgb(var(--primary-5)), #6B8E23);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.subtitle {
  font-size: 20px;
  color: var(--color-text-3);
  max-width: 600px;
  margin: 0 auto;
}

/* 艺术画廊网格 */
.product-gallery {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
  gap: 40px;
}

.product-item {
  display: flex;
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 20px;
  padding: 32px;
  transition: all 0.4s cubic-bezier(0.2, 0.8, 0.2, 1);
  position: relative;
  overflow: hidden;
  animation: fadeUp 0.8s ease backwards;
}

/* 交错动画 */
.product-item:nth-child(1) { animation-delay: 0.1s; }
.product-item:nth-child(2) { animation-delay: 0.2s; }
.product-item:nth-child(3) { animation-delay: 0.3s; }
.product-item:nth-child(4) { animation-delay: 0.4s; }

.product-item:hover {
  transform: translateY(-8px) scale(1.02);
  border-color: rgb(var(--primary-6));
  box-shadow: 0 24px 48px rgba(0,0,0,0.06);
}

.product-item::after {
  content: '';
  position: absolute;
  top: -50px; left: -50px;
  width: 100px; height: 100px;
  background: radial-gradient(circle, rgba(var(--primary-6), 0.1) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.5s ease;
}
.product-item:hover::after {
  opacity: 1;
}

.product-icon-box {
  flex-shrink: 0;
  width: 72px;
  height: 72px;
  background: var(--color-fill-2);
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 28px;
  transition: all 0.4s ease;
}
.product-item:hover .product-icon-box {
  background: rgb(var(--primary-6));
  transform: rotate(-5deg);
}

.main-icon {
  font-size: 36px;
  color: var(--color-text-2);
  transition: all 0.4s ease;
}
.product-item:hover .main-icon {
  color: #fff;
}

.meta-row {
  margin-bottom: 8px;
}
.category-tag {
  font-size: 13px;
  font-weight: 600;
  color: rgb(var(--primary-6));
  background: var(--color-primary-light-1);
  padding: 4px 10px;
  border-radius: 6px;
  text-transform: uppercase;
}

.product-name {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-1);
  margin-bottom: 12px;
}

.product-desc {
  font-size: 15px;
  color: var(--color-text-3);
  line-height: 1.6;
  margin-bottom: 24px;
}

.product-action {
  margin-top: auto;
}
.learn-more {
  font-weight: 600;
  color: var(--color-text-2);
  font-size: 15px;
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  transition: color 0.3s ease;
}
.learn-more .arrow {
  margin-left: 6px;
  transition: transform 0.3s ease;
}
.product-item:hover .learn-more {
  color: rgb(var(--primary-6));
}
.product-item:hover .learn-more .arrow {
  transform: translateX(6px);
}

@keyframes fadeInDown {
  from { opacity: 0; transform: translateY(-30px); }
  to { opacity: 1; transform: translateY(0); }
}
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 768px) {
  .product-item {
    flex-direction: column;
  }
  .product-icon-box {
    margin-bottom: 24px;
  }
}
</style>
