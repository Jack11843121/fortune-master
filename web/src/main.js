import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

// 每次成功加载后清除「旧缓存」标记，保证下次遇旧缓存时可再次自动刷新。
sessionStorage.removeItem('stale_reload')

createApp(App).use(router).mount('#app')