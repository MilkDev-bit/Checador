import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'
import 'izitoast/dist/css/iziToast.min.css'
import 'izimodal/css/iziModal.min.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
