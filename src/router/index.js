import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Preview from '../views/Preview.vue'
import Settings from '../views/Settings.vue'

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/',         component: Home,     name: 'home' },
    { path: '/preview',  component: Preview,  name: 'preview' },
    { path: '/settings', component: Settings, name: 'settings' },
  ],
})
