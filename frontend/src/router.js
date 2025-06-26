import { createRouter, createWebHistory } from 'vue-router';

import Login from './pages/Login.vue';
import Register from './pages/Register.vue';
import Profile from './pages/Profile.vue';
import Posts from './pages/Posts.vue';
import Groups from './pages/Groups.vue';
import Notifications from './pages/Notifications.vue';
import Chat from './pages/Chat.vue';
import Home from './pages/Home.vue';

const routes = [
  { path: '/', name: 'Home', component: Home, meta: { requiresAuth: true } },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { 
    path: '/profile', 
    name: 'Profile', 
    component: Profile,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/posts', 
    name: 'Posts', 
    component: Posts,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/groups', 
    name: 'Groups', 
    component: Groups,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/notifications', 
    name: 'Notifications', 
    component: Notifications,
    meta: { requiresAuth: true } 
  },
  { 
    path: '/chat', 
    name: 'Chat', 
    component: Chat,
    meta: { requiresAuth: true } 
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const publicPages = ['/login', '/register'];
  const authRequired = !publicPages.includes(to.path);

  if (!authRequired) {
    return next();
  }

  // Vérifie la session côté backend
  try {
    const res = await fetch('http://localhost:8081/sessionActive', {
      method: 'GET',
      credentials: 'include'
    });
    const data = await res.json();
    if (res.ok && data.success) {
      next();
    } else {
      next('/login');
    }
  } catch (e) {
    next('/login');
  }
});

export default router;