import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons';
import AutoImport from 'unplugin-auto-import/vite'


export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  const apiProxyTarget = env.VITE_API_PROXY_TARGET || 'http://localhost:8021';

  return {
  plugins: [
    vue(),
    createSvgIconsPlugin({
      iconDirs: [resolve(__dirname, 'src/assets/icons')],
      symbolId: 'icon-[name]',
    }),
    AutoImport({
      imports: [
        'vue',
        'vue-router',
        {
          '@arco-design/web-vue': ['Message', 'Modal', 'Notification', 'Drawer', 'Spin'],
        },
      ],
      dts: true,
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
    // 强制所有包使用同一个 Vue 实例，避免 vue-echarts 等库出现多实例冲突
    dedupe: ['vue', '@vue/runtime-core'],
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
    },
  },
  };
});