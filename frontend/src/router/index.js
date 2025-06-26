import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/pages/Login.vue'
import Posts from '@/pages/Posts.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login },
  { path: '/posts', name: 'Posts', component: Posts, meta: { requiresAuth: true } },
  { path: '/', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
