<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { IconSearch, IconClose } from '@arco-design/web-vue/es/icon/index';

// 获取所有图标名称（只需要文件名，不需要内容）
const svgModules = import.meta.glob('@/assets/icons/*.svg', { eager: true });

const allIcons = Object.keys(svgModules).map((key) =>
  key.replace('/src/assets/icons/', '').replace('.svg', '')
);

const props = defineProps<{
  modelValue?: string;
  placeholder?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const visible = ref(false);
const keyword = ref('');

const filtered = computed(() => {
  const kw = keyword.value.trim().toLowerCase();
  return kw ? allIcons.filter((name) => name.includes(kw)) : allIcons;
});

const select = (name: string) => {
  emit('update:modelValue', name);
  visible.value = false;
  keyword.value = '';
};

const clear = (e: Event) => {
  e.stopPropagation();
  emit('update:modelValue', '');
};

// 关闭时清空搜索
watch(visible, (val) => {
  if (!val) keyword.value = '';
});
</script>

<template>
  <a-trigger
    v-model:popup-visible="visible"
    trigger="click"
    position="bl"
    :popup-offset="4"
    popup-container="body"
    :click-outside-to-close="true"
  >
    <!-- 触发器：输入框样式 -->
    <div class="icon-picker-trigger" :class="{ 'is-active': visible }">
      <span v-if="modelValue" class="icon-picker-preview">
        <SvgIcon :name="modelValue" :size="16" />
        <span class="icon-picker-name">{{ modelValue }}</span>
      </span>
      <span v-else class="icon-picker-placeholder">{{ placeholder ?? '选择图标' }}</span>
      <IconClose v-if="modelValue" class="icon-picker-clear" @click="clear" />
    </div>

    <!-- 弹出面板 -->
    <template #content>
      <div class="icon-picker-panel">
        <div class="icon-picker-search">
          <a-input
            v-model="keyword"
            placeholder="搜索图标"
            allow-clear
            size="small"
          >
            <template #prefix><IconSearch /></template>
          </a-input>
        </div>
        <div class="icon-picker-count">{{ filtered.length }} 个图标</div>
        <div class="icon-picker-grid">
          <button
            v-for="name in filtered"
            :key="name"
            type="button"
            :class="['icon-picker-item', { 'is-selected': modelValue === name }]"
            :title="name"
            @click="select(name)"
          >
            <SvgIcon :name="name" :size="20" />
            <span class="icon-picker-item-name">{{ name }}</span>
          </button>
        </div>
        <div v-if="!filtered.length" class="icon-picker-empty">没有匹配的图标</div>
      </div>
    </template>
  </a-trigger>
</template>

<style scoped>
.icon-picker-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 32px;
  padding: 0 12px;
  border: 1px solid var(--color-border-2, rgba(148, 163, 184, 0.4));
  border-radius: 6px;
  background: var(--color-fill-2, rgba(255, 255, 255, 0.8));
  cursor: pointer;
  transition: border-color 0.2s;
  user-select: none;
}

.icon-picker-trigger:hover,
.icon-picker-trigger.is-active {
  border-color: rgb(var(--primary-6, 37, 99, 235));
}

.icon-picker-preview {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.icon-picker-name {
  font-size: 13px;
  color: var(--text-main);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.icon-picker-placeholder {
  flex: 1;
  font-size: 13px;
  color: var(--text-subtle);
}

.icon-picker-clear {
  flex-shrink: 0;
  color: var(--text-subtle);
  font-size: 12px;
  transition: color 0.2s;
}

.icon-picker-clear:hover {
  color: var(--text-main);
}

/* 弹出面板 */
.icon-picker-panel {
  width: 320px;
  padding: 10px;
  border-radius: 12px;
  background: var(--panel-bg);
  border: 1px solid var(--panel-border);
  box-shadow: var(--card-shadow);
  backdrop-filter: blur(18px);
}

.icon-picker-search {
  margin-bottom: 8px;
}

.icon-picker-count {
  margin-bottom: 6px;
  font-size: 12px;
  color: var(--text-subtle);
}

.icon-picker-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 4px;
  max-height: 280px;
  overflow-y: auto;
}

.icon-picker-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  padding: 6px 4px;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  color: var(--text-main);
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}

.icon-picker-item:hover {
  background: var(--blue-soft);
  border-color: rgba(37, 99, 235, 0.2);
}

.icon-picker-item.is-selected {
  background: var(--blue-soft);
  border-color: rgba(37, 99, 235, 0.4);
}

.icon-picker-item-name {
  font-size: 10px;
  color: var(--text-subtle);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.icon-picker-empty {
  padding: 20px 0;
  text-align: center;
  font-size: 13px;
  color: var(--text-subtle);
}
</style>
