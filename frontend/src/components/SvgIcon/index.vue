<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  name: string;
  size?: string | number;
  color?: string;
}>();

const symbolId = computed(() => `#icon-${props.name}`);

const sizeVal = computed(() => {
  if (!props.size) return '1em';
  return typeof props.size === 'number' ? `${props.size}px` : props.size;
});
</script>

<template>
  <svg
    :style="{
      width: sizeVal,
      height: sizeVal,
      color: color ?? 'currentColor',
      verticalAlign: '-0.15em',
    }"
    aria-hidden="true"
    class="svg-icon"
  >
    <use :href="symbolId" :fill="color" />
  </svg>
</template>

<style scoped>
.svg-icon {
  display: inline-block;
  overflow: hidden;
  fill: currentColor;
  flex-shrink: 0;
}
</style>
