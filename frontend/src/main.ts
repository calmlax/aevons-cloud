import { createApp } from 'vue';
import ArcoVue from '@arco-design/web-vue';
import '@arco-design/web-vue/dist/arco.css';
import 'virtual:svg-icons-register';

import App from './App.vue';
import SvgIcon from './components/SvgIcon/index.vue';
import { registerGlobalDirectives } from './directives';
import store from './store';
import { useAuthStore } from './store/modules/auth';
import { i18n, initializeLocale } from './locale';
import { useThemeStore } from './store/modules/theme';
import router from './router';
import './assets/style/global.css';
import './assets/style/animations.css';
import { useDict } from '@/components/DictTag/dict' 
import DictTag from '@/components/DictTag/index.vue'


const bootstrap = async () => {
	initializeLocale();

	const app = createApp(App);
	app.use(store);

	const authStore = useAuthStore();
	authStore.initializeAuth();

	const themeStore = useThemeStore();
	themeStore.initializeTheme();

	registerGlobalDirectives(app);
	app.config.globalProperties.$useDict = useDict
	app.use(i18n).use(ArcoVue).use(router)
	app.component('SvgIcon', SvgIcon)
	app.component('DictTag', DictTag)
	app.mount('#app');
};

void bootstrap();