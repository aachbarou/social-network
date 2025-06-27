<template>
  <div class="profile-page">
    <h1>Mon profil</h1>
    <div v-if="userStore.loading">
      <p>Chargement du profil...</p>
    </div>
    <div v-else-if="user">
      <div class="profile-header">
        <div style="display: flex; align-items: center; gap: 1rem;">
          <img v-if="user.avatar || user.imagePath"
               :src="getFullImageUrl(user.avatar || user.imagePath)"
               alt="Avatar" class="profile-avatar" style="max-width:120px; max-height:120px; object-fit:contain; border:2px solid #e879c6; background:#fff;" />
        </div>
        <div class="profile-info">
          <h2>{{ user.nickname || (user.firstName + ' ' + user.lastName) }}</h2>
          <p v-if="user.firstName"><strong>Prénom :</strong> {{ user.firstName }}</p>
          <p v-if="user.lastName"><strong>Nom :</strong> {{ user.lastName }}</p>
          <p v-if="user.nickname"><strong>Nickname :</strong> {{ user.nickname }}</p>
          <p v-if="user.email"><strong>Email :</strong> {{ user.email }}</p>
          <p v-if="user.dateOfBirth"><strong>Date de naissance :</strong> {{ user.dateOfBirth }}</p>
          <p v-if="user.about"><strong>À propos :</strong> {{ user.about }}</p>
          <p v-if="user.status"><strong>Statut du profil :</strong> {{ user.status }}</p>
          <p v-if="user.currentUser !== undefined"><strong>Vous êtes connecté</strong></p>
        </div>
      </div>
      <div class="profile-actions">
        <button @click="toggleStatus">
          {{ user.status && user.status.toLowerCase() === 'public' ? 'Rendre le profil privé' : 'Rendre le profil public' }}
        </button>
      </div>
      <div class="profile-followers">
        <h3>Followers ({{ followers.length }})</h3>
        <ul>
          <li v-for="f in followers" :key="f.id">{{ f.nickname || f.firstName }}</li>
        </ul>
        <h3>Abonnements ({{ following.length }})</h3>
        <ul>
          <li v-for="f in following" :key="f.id">{{ f.nickname || f.firstName }}</li>
        </ul>
      </div>
      <div class="profile-posts">
        <h3>Mes posts ({{ posts.length }})</h3>
        <div v-if="posts.length === 0">Aucun post pour le moment.</div>
        <div v-for="post in posts" :key="post.id" class="profile-post-card">
          <div class="profile-post-content">{{ post.content }}</div>
          <img v-if="post.image" :src="post.image" alt="Image du post" class="profile-post-image" />
        </div>
      </div>
    </div>
    <div v-else-if="userStore.error">
      <p class="error-msg">{{ userStore.error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '../stores/userStore'

const userStore = useUserStore()
const user = computed(() => userStore.user)
const followers = ref([])
const following = ref([])
const posts = ref([])

const fetchFollowers = async () => {
  const res = await fetch('http://localhost:8081/followers', { credentials: 'include' })
  const data = await res.json()
  if (res.ok && data.users) followers.value = data.users
}
const fetchFollowing = async () => {
  const res = await fetch('http://localhost:8081/following', { credentials: 'include' })
  const data = await res.json()
  if (res.ok && data.users) following.value = data.users
}

const fetchUserPosts = async () => {
  if (!user.value || !user.value.id) return
  const res = await fetch(`http://localhost:8081/userPosts?id=${user.value.id}`, { credentials: 'include' })
  const data = await res.json()
  if (res.ok && data.posts) posts.value = data.posts
}

const toggleStatus = async () => {
  if (!user.value) return
  const newStatus = user.value.status && user.value.status.toLowerCase() === 'public' ? 'private' : 'public'
  await userStore.changeStatus(newStatus)
  await userStore.fetchUserData() // recharge le profil après changement
}

function getFullImageUrl(path) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  // Ajoute le host backend si le chemin est relatif
  return `http://localhost:8081/${path.replace(/^\/+/, '')}`
}

onMounted(async () => {
  if (!user.value) {
    await userStore.reloadUserAfterRefresh()
  }
  await userStore.fetchUserData()
  fetchFollowers()
  fetchFollowing()
  await fetchUserPosts()
})
</script>

<style scoped>
.profile-page {
  max-width: 500px;
  margin: 2rem auto;
  padding: 2rem 1.5rem;
  background: #181824;
  border-radius: 18px;
  box-shadow: 0 4px 24px rgba(120, 119, 198, 0.10);
}
.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 1.5rem;
}
.profile-avatar {
  width: 110px;
  height: 110px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #e879c6;
  background: #fff;
  box-shadow: 0 2px 8px rgba(120, 119, 198, 0.10);
}
.profile-info {
  text-align: center;
}
.profile-info h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  color: #e879c6;
}
.profile-info p {
  margin: 0.2rem 0;
  color: #d1d5db;
  font-size: 1rem;
}
.profile-actions {
  display: flex;
  justify-content: center;
  margin: 1.5rem 0 1rem 0;
}
.profile-actions button {
  background: linear-gradient(90deg, #e879c6 0%, #78dbff 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.profile-actions button:hover {
  background: linear-gradient(90deg, #78dbff 0%, #e879c6 100%);
}
.profile-followers {
  margin-top: 1.5rem;
  background: #23233a;
  border-radius: 10px;
  padding: 1rem;
}
.profile-followers h3 {
  margin-bottom: 0.5rem;
  color: #78dbff;
}
.profile-followers ul {
  display: flex;
  flex-wrap: wrap;
  gap: 0.7rem;
  list-style: none;
  padding: 0;
  margin: 0 0 1rem 0;
}
.profile-followers li {
  background: #282846;
  color: #fff;
  border-radius: 6px;
  padding: 0.3rem 0.8rem;
  font-size: 0.98rem;
}
.profile-posts {
  margin-top: 2rem;
}
.profile-posts h3 {
  color: #78dbff;
  margin-bottom: 0.7rem;
}
.profile-post-card {
  background: #23233a;
  border-radius: 12px;
  padding: 1rem;
  margin-bottom: 1rem;
  box-shadow: 0 2px 8px rgba(120, 119, 198, 0.08);
}
.profile-post-content {
  font-size: 1.08rem;
  margin-bottom: 0.5rem;
  color: #fff;
}
.profile-post-image {
  max-width: 100%;
  border-radius: 8px;
  margin-top: 0.5rem;
  background: #fff;
}
.error-msg {
  color: #e74c3c;
  font-weight: bold;
  margin-top: 2rem;
  text-align: center;
}
</style>
