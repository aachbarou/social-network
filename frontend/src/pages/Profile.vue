<template>
  <div class="profile-page">
    <!-- Bouton pour changer le statut public/privé si c'est son propre profil -->
    <div v-if="isOwnProfile">
      <button @click="askPrivacyChange(isPrivate ? 'public' : 'private')">
        Rendre profil {{ isPrivate ? 'public' : 'privé' }}
      </button>
    </div>
    <!-- Boutons d'action pour les autres profils -->
    <div v-else>
      <button v-if="isFollowing" @click="unfollow">Unfollow</button>
      <button v-else-if="isPrivate && !isFollowRequested" @click="follow">Demander à suivre</button>
      <button v-else-if="!isFollowing && isPublic && !isFollowRequested" @click="follow">Suivre</button>
      <span v-if="isFollowRequested">Demande envoyée</span>
    </div>
    <!-- Affichage conditionnel des infos -->
    <div v-if="isOwnProfile || isPublic || isFollower || isFollowing" class="profile-infos">
      <h1 v-if="isOwnProfile">Mon profil</h1>
      <h1 v-else>Profil utilisateur</h1>
      <div v-if="loading">
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
            <p v-if="user.currentUser !== undefined && isOwnProfile"><strong>Vous êtes connecté</strong></p>
          </div>
        </div>
        <div class="profile-followers">
          <h3>Followers ({{ profileFollowers.length }})</h3>
          <ul>
            <li v-for="f in profileFollowers" :key="f.id">{{ f.nickname || f.firstName }}</li>
          </ul>
          <h3>Abonnements ({{ profileFollowing.length }})</h3>
          <ul>
            <li v-for="f in profileFollowing" :key="f.id">{{ f.nickname || f.firstName }}</li>
          </ul>
        </div>
        <div class="profile-posts">
          <h3>Posts ({{ posts.length }})</h3>
          <div v-if="posts.length === 0">Aucun post pour le moment.</div>
          <div v-for="post in posts" :key="post.id" class="profile-post-card">
            <div class="profile-post-content">{{ post.content }}</div>
            <img v-if="post.image" :src="post.image" alt="Image du post" class="profile-post-image" />
          </div>
        </div>
      </div>
      <div v-else>
        <p class="error-msg">Profil introuvable.</p>
      </div>
    </div>
    <div v-else>
      <div class="profile-info">
        <h2>{{ user?.nickname || (user?.firstName + ' ' + user?.lastName) || 'Profil utilisateur' }}</h2>
      </div>
      <div class="profile-private-msg">Ce profil est privé.</div>
    </div>
    <!-- Pop-up de confirmation -->
    <div v-if="showPrivacyConfirm" class="popup-confirm">
      <div>Confirmer le changement de statut du profil ?</div>
      <button @click="confirmPrivacyChange">Confirmer</button>
      <button @click="cancelPrivacyChange">Annuler</button>
    </div>
    <!-- Gestion des demandes de suivi -->
    <div v-if="isOwnProfile && isPrivate">
      <h3>Demandes de suivi en attente</h3>
      <div v-if="followStore.followRequests.length === 0">Aucune demande.</div>
      <ul>
        <li v-for="req in followStore.followRequests" :key="req.id" style="margin-bottom: 0.5rem;">
          <span>{{ req.nickname || req.firstName || 'Utilisateur' }}</span>
          <button @click="respondToRequest(req.requestId || req.id, 'ACCEPT')">Accepter</button>
          <button @click="respondToRequest(req.requestId || req.id, 'DECLINE')">Refuser</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useFollowStore } from '../stores/followStore'

const route = useRoute()
const userStore = useUserStore()
const followStore = useFollowStore()
const user = ref(null)
const profileFollowers = ref([])
const profileFollowing = ref([])
const following = followStore.following // Utilise directement le ref du store Pinia
const posts = ref([])
const loading = ref(true)
const isOwnProfile = computed(() => !route.params.id)
const isPublic = computed(() => user.value?.status?.toLowerCase() !== 'private')
const isPrivate = computed(() => user.value?.status?.toLowerCase() === 'private')
const showPrivacyConfirm = ref(false)
const pendingPrivacy = ref(null)
const isFollower = computed(() => !!user.value?.follower)
// const isFollowing = computed(() => {
//   return user.value?.id && Array.isArray(following.value)
//     ? following.value.some(f => f.id === user.value.id)
//     : false
// })
const isFollowing = computed(() => !!user.value?.following)
console.log(isFollowing.value, 'isFollowing')
const isFollowRequested = computed(() => user.value?.id ? followStore.isFollowRequested(user.value.id) : false)

function getFullImageUrl(path) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:8081/${path.replace(/^\/+/,'')}`
}

const fetchUserPosts = async (userId) => {
  const res = await fetch(`http://localhost:8081/userPosts?id=${userId}`, { credentials: 'include' })
  const data = await res.json()
  if (res.ok && data.posts) posts.value = data.posts
}

const askPrivacyChange = (newPrivacy) => {
  pendingPrivacy.value = newPrivacy
  showPrivacyConfirm.value = true
}
const confirmPrivacyChange = async () => {
  if (!user.value) return
  const newStatus = pendingPrivacy.value === 'public' ? 'public' : 'private'
  await userStore.changeStatus(newStatus)
  showPrivacyConfirm.value = false
  pendingPrivacy.value = null
  await loadProfile()
}
const cancelPrivacyChange = () => {
  showPrivacyConfirm.value = false
  pendingPrivacy.value = null
}
function follow() {
  if (!user.value?.id) return
  if (isPublic.value) {
    followStore.follow(user.value.id)
      .then(async () => {
        await loadProfile() // Recharge le profil pour mettre à jour le bouton
      })
      .catch(() => {
        // Optionnel : afficher une erreur
      })
  } else {
    followStore.requestFollow(user.value.id)
    // Optionnel : afficher une notification ou un message
  }
}

function unfollow() {
  if (!user.value?.id) return
  followStore.unfollow(user.value.id)
    .then(async () => {
      await loadProfile() // Recharge le profil pour mettre à jour le bouton
    })
}

async function loadProfile() {
  loading.value = true
  let userId
  if (route.params.id) {
    // Profil d'un autre utilisateur
    userId = route.params.id
    const res = await fetch(`http://localhost:8081/userData?userId=${userId}`, { credentials: 'include' })
    const data = await res.json()
    user.value = (data.users && data.users.length > 0) ? data.users[0] : null
    // Affichage : followers/following du profil visité
    await followStore.fetchFollowers(userId)
    profileFollowers.value = [...followStore.followers]
    await followStore.fetchFollowing(userId)
    profileFollowing.value = [...followStore.following]
    // Pour le bouton : always fetch following du user courant
    if (userStore.user?.id) {
      await followStore.fetchFollowing(userStore.user.id)
      console.log('following (utilisateur courant):', following.value)
    }
  } else {
    // Mon propre profil
    await userStore.reloadUserAfterRefresh()
    await userStore.fetchUserData()
    user.value = userStore.user
    userId = user.value?.id
    await followStore.fetchFollowers(userId)
    profileFollowers.value = [...followStore.followers]
    await followStore.fetchFollowing(userId)
    profileFollowing.value = [...followStore.following]
    // Recharge aussi la liste des following du user courant (lui-même)
    await followStore.fetchFollowing(userId)
    console.log('following (utilisateur courant):', following.value)
  }
  if (userId) {
    await fetchUserPosts(userId)
  } else {
    profileFollowers.value = []
    profileFollowing.value = []
    posts.value = []
  }
  loading.value = false
}

onMounted(loadProfile)
watch(() => route.params.id, loadProfile)
onMounted(async () => {
  await loadProfile()
  if (isOwnProfile.value && isPrivate.value) {
    await followStore.fetchFollowRequests()
  }
})

watch([isOwnProfile, isPrivate], async ([own, priv]) => {
  if (own && priv) {
    await followStore.fetchFollowRequests()
  }
})

async function respondToRequest(requestId, response) {
  await followStore.respondToFollowRequest(requestId, response)
  await followStore.fetchFollowRequests()
}
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
.profile-private-msg {
  color: #888;
  font-style: italic;
  margin-top: 1em;
}
.popup-confirm {
  position: fixed;
  top: 30%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  border: 1px solid #ccc;
  padding: 2em;
  z-index: 1000;
}
</style>
