<template>
  <div class="profile-page">
    <h1>Mon profil</h1>
    <div v-if="user">
      <div class="profile-header">
        <img v-if="user.avatar || user.imagePath" :src="user.avatar || user.imagePath" alt="Avatar" class="profile-avatar" />
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
        <button @click="toggleStatus">Rendre le profil {{ user.status === 'public' ? 'privé' : 'public' }}</button>
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
    </div>
    <div v-else>
      <p>Chargement du profil...</p>
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

const toggleStatus = async () => {
  if (!user.value) return
  const newStatus = user.value.status === 'public' ? 'private' : 'public'
  const res = await fetch('http://localhost:8081/changeStatus', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status: newStatus })
  })
  if (res.ok) {
    user.value.status = newStatus
  }
}

onMounted(async () => {
  if (!user.value) {
    await userStore.reloadUserAfterRefresh()
  }
  await userStore.fetchUserData()
  fetchFollowers()
  fetchFollowing()
})
</script>

<style scoped>
.profile-page {
  max-width: 600px;
  margin: 0 auto;
  padding: 2rem;
}
.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}
.profile-avatar {
  width: 90px;
  height: 90px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e879c6;
}
.profile-info h2 {
  margin: 0 0 0.5rem 0;
}
.profile-actions {
  margin: 1rem 0;
}
.profile-followers {
  margin-top: 2rem;
}
</style>
