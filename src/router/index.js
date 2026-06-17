import { createRouter, createWebHashHistory } from 'vue-router'
import Home     from '../views/Home.vue'
import Preview  from '../views/Preview.vue'
import Settings from '../views/Settings.vue'
import Help     from '../views/Help.vue'
import About    from '../views/About.vue'

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/',         component: Home,     name: 'home',     meta: { title: 'Generate Report' } },
    { path: '/preview',  component: Preview,  name: 'preview',  meta: { title: 'Preview' } },
    { path: '/settings', component: Settings, name: 'settings', meta: { title: 'Settings' } },
    { path: '/help',     component: Help,     name: 'help',     meta: { title: 'Help & FAQs' } },
    { path: '/about',    component: About,    name: 'about',    meta: { title: 'About' } },
  ],
})
