import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import 'bootstrap-icons/font/bootstrap-icons.min.css'
import "preline/preline";
import router from "@/router/index.ts";

const app = createApp(App)
app.use(router)
    .mount('#app')
