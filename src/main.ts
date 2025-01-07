import { createApp } from 'vue';
import App from './App.vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import * as Icons from '@element-plus/icons-vue';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import axios from 'axios'

const app = createApp(App);

// 注册所有图标
for (const [key, component] of Object.entries(Icons)) {
  app.component(key, component);
}

app.use(ElementPlus, { locale: zhCn });

axios.defaults.baseURL = process.env.VUE_APP_API_URL || 'http://localhost:80'

app.mount('#app');
