<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();
</script>

<template>
  <div class="panel-card query-aevo" style="padding: 24px">
    <a-typography-title :heading="2">带 Query 参数的页面示例</a-typography-title>
    <a-typography-paragraph>
      此页面是通过在 <code>menu.json</code> 中配置带有 <code>query</code> 参数动态生成的路由。
    </a-typography-paragraph>

    <a-card title="当前页面的 Query 参数" style="margin-top: 16px;">
      <a-descriptions 
        :data="Object.entries(route.query).map(([k, v]) => ({ label: k, value: String(v) }))" 
        bordered 
      />
      
      <a-alert type="info" style="margin-top: 16px;">
        您可以尝试在地址栏随意修改 / 追加 <code>?name=abc&id=999</code>，页面中的参数会自动响应更新。
      </a-alert>
    </a-card>
  </div>
</template>
