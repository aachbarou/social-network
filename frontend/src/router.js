import { createRouter, createWebHistory } from 'vue-router';

import Login from './pages/Login.vue';
import Register from './pages/Register.vue';
import Profile from './pages/Profile.vue';
import Groups from './pages/Groups.vue';
import Notifications from './pages/Notifications.vue';
import Chat from './pages/Chat.vue';
import Home from './pages/Home.vue';
import MainLayout from './components/MainLayout.vue';

const routes = [
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Home', component: Home },
      { path: 'profile', name: 'Profile', component: Profile },
      { path: 'profile/:id', name: 'ProfileOther', component: Profile },
      { path: 'groups', name: 'Groups', component: Groups },
      { path: 'notifications', name: 'Notifications', component: Notifications },
      { path: 'chat', name: 'Chat', component: Chat }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const publicPages = ['/login', '/register'];
  const isPublicPage = publicPages.includes(to.path);

  // Check user authentication status
  try {
    const res = await fetch('http://localhost:8081/sessionActive', {
      method: 'GET',
      credentials: 'include'
    });
    const data = await res.json();
    const isAuthenticated = res.ok && data.type === 'Success';
    
    if (isPublicPage) {
      if (isAuthenticated) {
        // User is authenticated but trying to access login/register, redirect to home
        return next('/');
      } else {
        // User is not authenticated, allow access to login/register
        return next();
      }
    } else {
      // For protected routes, check authentication
      if (isAuthenticated) {
        // User is authenticated, allow access
        next();
      } else {
        // User is not authenticated, redirect to login
        next('/login');
      }
    }
  } catch (e) {
    if (isPublicPage) {
      // Network error on public page, allow access (user can still try to login)
      next();
    } else {
      // Network error on protected route, redirect to login
      next('/login');
    }
  }
});

export default router;