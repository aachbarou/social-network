import { createRouter, createWebHistory } from 'vue-router';

import Login from './pages/Login.vue';
import Register from './pages/Register.vue';
import Profile from './pages/Profile.vue';
import Posts from './pages/Posts.vue';
import Groups from './pages/Groups.vue';
import Notifications from './pages/Notifications.vue';
import Chat from './pages/Chat.vue';

const routes = [
  { path: '/', name: 'Home', component: Login },
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

router.beforeEach((to, from, next) => {
  console.log("[Debug] Route:", to.path);
  console.log("[Debug] Requires Auth:", to.meta.requiresAuth);
  console.log("[Debug] Token exists:", !!localStorage.getItem("token"));

  if (to.meta.requiresAuth && !localStorage.getItem("token")) {
    console.log("[Debug] 🔴 Blocked - Redirecting to login");
    next({ name: "Login" }); // Force redirect
    return; // Stop execution
  }
  next(); // Proceed
});

export default router;