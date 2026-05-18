<script setup lang="ts">
import { ref, computed } from 'vue';
import { IconSearch } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n';

const searchKeyword = ref('');
const activeKey = ref(['1']);
const { t } = useI18n();

const faqs = computed(() => [
  {
    id: '1',
    question: t('portal.helpPage.faqs.0.q'),
    answer: t('portal.helpPage.faqs.0.a')
  },
  {
    id: '2',
    question: t('portal.helpPage.faqs.1.q'),
    answer: t('portal.helpPage.faqs.1.a')
  },
  {
    id: '3',
    question: t('portal.helpPage.faqs.2.q'),
    answer: t('portal.helpPage.faqs.2.a')
  },
  {
    id: '4',
    question: t('portal.helpPage.faqs.3.q'),
    answer: t('portal.helpPage.faqs.3.a')
  }
]);
</script>

<template>
  <div class="help-wrapper">
    <div class="help-header">
      <h1 class="help-title">{{ t('portal.helpPage.title') }}</h1>
      <div class="search-box">
        <a-input 
          v-model="searchKeyword" 
          size="large" 
          :placeholder="t('portal.helpPage.searchPlaceholder')" 
          allow-clear
          class="custom-search-input"
        >
          <template #prefix>
            <IconSearch />
          </template>
        </a-input>
      </div>
    </div>

    <div class="faq-container">
      <h2 class="faq-section-title">{{ t('portal.helpPage.faqTitle') }}</h2>
      
      <a-collapse v-model:active-key="activeKey" :bordered="false" class="custom-collapse">
        <a-collapse-item 
          v-for="item in faqs" 
          :key="item.id" 
          :header="item.question"
          class="custom-collapse-item"
        >
          <p class="faq-answer">{{ item.answer }}</p>
        </a-collapse-item>
      </a-collapse>
    </div>
    
    <div class="contact-support">
      <!-- <p>{{ t('portal.helpPage.contactText') }}</p> -->
      <!-- <a-button type="outline" size="large">{{ t('portal.helpPage.contactBtn') }}</a-button> -->
    </div>
  </div>
</template>

<style scoped>
.help-wrapper {
  background-color: var(--color-bg-1);
  min-height: calc(100vh - 64px);
}

.help-header {
  padding: 100px 24px;
  background: var(--color-fill-1);
  text-align: center;
  border-bottom: 1px solid var(--color-border-1);
}

.help-title {
  font-size: 42px;
  font-weight: 800;
  color: var(--color-text-1);
  margin-bottom: 40px;
  letter-spacing: -1px;
}

.search-box {
  max-width: 680px;
  margin: 0 auto;
}

:deep(.custom-search-input) {
  border-radius: 12px;
  padding: 8px 16px;
  background-color: var(--color-bg-2);
  border: 1px solid transparent;
  box-shadow: 0 8px 24px rgba(0,0,0,0.05);
  transition: all 0.3s ease;
}
:deep(.custom-search-input:focus-within) {
  border-color: rgb(var(--primary-6));
  box-shadow: 0 8px 24px rgba(var(--primary-6), 0.15);
  background-color: var(--color-bg-1);
}
:deep(.custom-search-input .arco-input-prefix) {
  font-size: 20px;
  color: var(--color-text-3);
  margin-right: 12px;
}
:deep(.custom-search-input .arco-input) {
  font-size: 16px;
}

.faq-container {
  max-width: 800px;
  margin: 80px auto;
  padding: 0 24px;
}

.faq-section-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 32px;
  color: var(--color-text-1);
}

:deep(.custom-collapse) {
  background: transparent;
}
:deep(.custom-collapse-item) {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-1);
  border-radius: 16px !important;
  margin-bottom: 16px;
  overflow: hidden;
  transition: all 0.3s ease;
}
:deep(.custom-collapse-item:hover) {
  border-color: var(--color-border-3);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}
:deep(.custom-collapse-item > .arco-collapse-item-header) {
  border: none;
  padding: 24px;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-1);
}
:deep(.custom-collapse-item > .arco-collapse-item-content) {
  background-color: var(--color-bg-2);
  padding: 0 24px 24px;
}

.faq-answer {
  font-size: 15px;
  line-height: 1.8;
  color: var(--color-text-3);
  margin: 0;
}

.contact-support {
  text-align: center;
  padding: 60px 24px 120px;
}
.contact-support p {
  font-size: 18px;
  color: var(--color-text-2);
  margin-bottom: 24px;
}
</style>
