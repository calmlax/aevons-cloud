UI编写原则
1.支持简体中文和英文多语言
2.支持arco-design-vue 的官方浅色和官方深色模式切换
3.支持PCWeb，手机h5, 平板
4.尽量使用加载骨架，除非不合理
5.目录层级要规范，如下
    web/
    ├── public/
    ├── src/
    │   ├── api/
    │   ├── assets/
    │   ├── components/
    │   ├── config/
    │   ├── directives/
    │   ├── layout/
    │   ├── locale/
    │   ├── mock/
    │   ├── router/
    │   ├── store/
    │   ├── types/
    │   ├── utils/
    │   ├── views/
    │   ├── App.vue
    │   └── main.ts
    ├── .env
    ├── .env.development
    ├── .env.production
    ├── .env.staging
    ├── .gitignore
6.用户体验ui效果动效要突出重视
7.页面要简洁干净
8.页面代码注释简体中文
9.代码要简洁优雅，性能合理
10.弹框要支持手机模式
11.vue3 + arco-design-vue 前端
12.web/src/assets/icons 存放 svg图标
