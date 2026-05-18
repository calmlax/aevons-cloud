<script setup lang="ts">
import { computed } from 'vue';
import { PRESET_COLORS } from './types';
import type { DictOption } from './types';

const props = withDefaults(defineProps<{
  /** 字典选项列表（后端 DictDataDTO 数组） */
  options?: DictOption[];
  /** 当前值（对应 dictValue） */
  value: string | number | null | undefined;
  /** 找不到匹配项时的回退文字 */
  fallback?: string;
  /** 标签尺寸 */
  size?: 'small' | 'medium' | 'large';
  /** 是否显示圆点指示器 */
  dot?: boolean;
}>(), {
  options: () => [],
  fallback: undefined,
  size: 'medium',
  dot: false,
});

const safeOptions = computed(() => Array.isArray(props.options) ? props.options : []);

// 匹配逻辑：根据 dictValue 匹配（对应后端字段）
const matched = computed<DictOption | undefined>(() =>
  safeOptions.value.find((opt) => String(opt.dictValue ?? opt.value) === String(props.value))
);

// 颜色样式（读取 tagType）
const colorStyle = computed(() => {
  if (!matched.value) {
    return PRESET_COLORS.default;
  }
  const colorKey = matched.value.tagType ?? matched.value.color ?? 'default';
  return PRESET_COLORS[colorKey] ?? PRESET_COLORS.default;
});

// 显示文本（读取 text 字段）
const label = computed(() => {
  if (matched.value) return matched.value.label ?? matched.value.text ?? matched.value.value ?? matched.value.dictValue;
  if (props.fallback !== undefined) return props.fallback;
  return props.value != null ? String(props.value) : '—';
});

// 自定义类名（支持 tagClass）
const customClass = computed(() => matched.value?.tagClass || '');
</script>

<template>
  <span
    class="dict-tag"
    :class="[`dict-tag--${size}`, customClass]"
    :style="{
      background: colorStyle.bg,
      color: colorStyle.text,
      borderColor: colorStyle.border,
    }"
  >
    <!-- 图标支持 -->
    <span v-if="matched?.icon" class="dict-tag-icon">{{ matched.icon }}</span>
    <span v-if="dot" class="dict-tag-dot" :style="{ background: colorStyle.text }" />
    {{ label }}
  </span>
</template>

<style scoped>
.dict-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  border: 1px solid;
  border-radius: 6px;
  font-weight: 500;
  white-space: nowrap;
  line-height: 1;
  transition:
    background var(--dur-fast) var(--ease-out),
    color var(--dur-fast) var(--ease-out);
}

.dict-tag--small {
  padding: 2px 7px;
  font-size: 11px;
  border-radius: 5px;
}

.dict-tag--medium {
  padding: 3px 9px;
  font-size: 12px;
}

.dict-tag--large {
  padding: 5px 12px;
  font-size: 13px;
  border-radius: 8px;
}

.dict-tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dict-tag-icon {
  font-size: 12px;
  flex-shrink: 0;
}
</style>
