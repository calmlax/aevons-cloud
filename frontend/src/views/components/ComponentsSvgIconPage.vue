<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import IconPicker from '@/components/IconPicker/index.vue';

const { t } = useI18n();

// 图标展示用的常用图标列表
const showcaseIcons = [
  'dashboard', 'user', 'setting', 'search', 'edit', 'delete',
  'download', 'upload', 'chart', 'form', 'list', 'table',
  'message', 'notice', 'lock', 'eye', 'star', 'check',
  'alert', 'info', 'question', 'bug', 'code', 'tool',
];

const sizes = [16, 20, 24, 32, 40];

// 表单示例
const form = reactive({
  name: '',
  icon: '',
  color: '#2563eb',
  description: '',
});
const formSubmitted = ref(false);

const onSubmit = () => {
  if (!form.name || !form.icon) return;
  formSubmitted.value = true;
};

const onReset = () => {
  form.name = '';
  form.icon = '';
  form.color = '#2563eb';
  form.description = '';
  formSubmitted.value = false;
};
</script>

<template>
  <div class="page-stack">
    <!-- 基础展示 -->
    <a-card class="panel-card" :bordered="false" :title="t('svgIcon.basicTitle')">
      <div class="icon-showcase-grid">
        <div
          v-for="name in showcaseIcons"
          :key="name"
          class="icon-showcase-item"
        >
          <SvgIcon :name="name" :size="24" />
          <span>{{ name }}</span>
        </div>
      </div>
    </a-card>

    <!-- 尺寸示例 -->
    <a-card class="panel-card" :bordered="false" :title="t('svgIcon.sizeTitle')">
      <div class="icon-size-row">
        <div v-for="size in sizes" :key="size" class="icon-size-item">
          <SvgIcon name="dashboard" :size="size" />
          <span>{{ size }}px</span>
        </div>
      </div>
    </a-card>

    <!-- 颜色示例 -->
    <a-card class="panel-card" :bordered="false" :title="t('svgIcon.colorTitle')">
      <div class="icon-color-row">
        <SvgIcon name="star" :size="28" color="#f59e0b" />
        <SvgIcon name="check" :size="28" color="#10b981" />
        <SvgIcon name="alert" :size="28" color="#ef4444" />
        <SvgIcon name="info" :size="28" color="#3b82f6" />
        <SvgIcon name="lock" :size="28" color="#8b5cf6" />
        <SvgIcon name="bug" :size="28" color="#f97316" />
      </div>
    </a-card>

    <!-- 表单示例 -->
    <a-card class="panel-card" :bordered="false" :title="t('svgIcon.formTitle')">
      <a-row :gutter="24">
        <a-col :xs="24" :md="14">
          <a-form :model="form" layout="vertical">
            <a-form-item :label="t('svgIcon.form.name')" required>
              <a-input
                v-model="form.name"
                :placeholder="t('svgIcon.form.namePlaceholder')"
                allow-clear
              />
            </a-form-item>

            <a-form-item :label="t('svgIcon.form.icon')" required>
              <IconPicker
                v-model="form.icon"
                :placeholder="t('svgIcon.form.iconPlaceholder')"
              />
            </a-form-item>

            <a-form-item :label="t('svgIcon.form.color')">
              <a-input v-model="form.color" placeholder="#2563eb" allow-clear>
                <template #prefix>
                  <span
                    class="color-dot"
                    :style="{ background: form.color || '#2563eb' }"
                  />
                </template>
              </a-input>
            </a-form-item>

            <a-form-item :label="t('svgIcon.form.description')">
              <a-textarea
                v-model="form.description"
                :placeholder="t('svgIcon.form.descriptionPlaceholder')"
                :max-length="100"
                show-word-limit
                :auto-size="{ minRows: 2, maxRows: 4 }"
              />
            </a-form-item>

            <a-form-item>
              <a-space>
                <a-button
                  type="primary"
                  :disabled="!form.name || !form.icon"
                  @click="onSubmit"
                >
                  {{ t('svgIcon.form.submit') }}
                </a-button>
                <a-button @click="onReset">{{ t('svgIcon.form.reset') }}</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-col>

        <!-- 预览卡片 -->
        <a-col :xs="24" :md="10">
          <div class="icon-form-preview">
            <p class="eyebrow">{{ t('svgIcon.form.preview') }}</p>
            <div class="icon-preview-card">
              <div
                class="icon-preview-badge"
                :style="{ color: form.color || '#2563eb' }"
              >
                <SvgIcon
                  v-if="form.icon"
                  :name="form.icon"
                  :size="36"
                  :color="form.color || '#2563eb'"
                />
                <span v-else class="icon-preview-empty">?</span>
              </div>
              <div class="icon-preview-copy">
                <strong>{{ form.name || t('svgIcon.form.namePlaceholder') }}</strong>
                <span>{{ form.icon || t('svgIcon.form.iconPlaceholder') }}</span>
                <p v-if="form.description">{{ form.description }}</p>
              </div>
            </div>

            <a-alert
              v-if="formSubmitted"
              type="success"
              :content="t('svgIcon.form.submitSuccess', { name: form.name, icon: form.icon })"
              style="margin-top: 12px"
            />
          </div>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<style scoped>
.icon-showcase-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 8px;
}

.icon-showcase-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 6px;
  border: 1px solid var(--panel-border);
  border-radius: 10px;
  cursor: default;
  transition: background 0.15s;
}

.icon-showcase-item:hover {
  background: rgba(37, 99, 235, 0.06);
}

.icon-showcase-item span {
  font-size: 11px;
  color: var(--text-subtle);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.icon-size-row {
  display: flex;
  align-items: flex-end;
  gap: 24px;
  flex-wrap: wrap;
}

.icon-size-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.icon-size-item span {
  font-size: 12px;
  color: var(--text-subtle);
}

.icon-color-row {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.icon-form-preview {
  height: 100%;
}

.icon-preview-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  border: 1px solid var(--panel-border);
  border-radius: 14px;
  background: var(--panel-bg);
  margin-top: 8px;
}

.icon-preview-badge {
  display: grid;
  place-items: center;
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: rgba(37, 99, 235, 0.08);
  flex-shrink: 0;
}

.icon-preview-empty {
  font-size: 24px;
  color: var(--text-subtle);
}

.icon-preview-copy {
  min-width: 0;
}

.icon-preview-copy strong {
  display: block;
  font-size: 15px;
  color: var(--text-main);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.icon-preview-copy span {
  font-size: 12px;
  color: var(--text-subtle);
}

.icon-preview-copy p {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--text-subtle);
}

.color-dot {
  display: inline-block;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1px solid rgba(0, 0, 0, 0.1);
}
</style>
