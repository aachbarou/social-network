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
  { path: '/profile', name: 'Profile', component: Profile },
  { path: '/posts', name: 'Posts', component: Posts },
  { path: '/groups', name: 'Groups', component: Groups },
  { path: '/notifications', name: 'Notifications', component: Notifications },
  { path: '/chat', name: 'Chat', component: Chat },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;