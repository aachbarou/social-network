<script setup>
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div id="app">
    <nav>
      <router-link to="/">Accueil</router-link> |
      <router-link to="/login">Connexion</router-link> |
      <router-link to="/register">Inscription</router-link> |
      <router-link to="/profile">Profil</router-link> |
      <router-link to="/posts">Publications</router-link> |
      <router-link to="/groups">Groupes</router-link> |
      <router-link to="/notifications">Notifications</router-link> |
      <router-link to="/chat">Chat</router-link>
    </nav>
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style>
nav {
  background: #f0f0f0;
  padding: 1em;
}
nav a {
  margin: 0 0.5em;
  text-decoration: none;
  color: #333;
}
nav a.router-link-exact-active {
  font-weight: bold;
  color: #007bff;
}
</style>
