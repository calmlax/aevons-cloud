<script setup lang="ts">
import { useAppLocale } from './locale';

const { arcoLocale } = useAppLocale();
</script>

<template>
  <a-config-provider :locale="arcoLocale">
    <router-view />
  </a-config-provider>
</template>