import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useFollowStore = defineStore('follow', () => {
  // Listes réactives
  const followers = ref([])
  const following = ref([])
  const followRequests = ref([])

  // Récupère la liste des followers
  async function fetchFollowers(userId) {
    try {
      const res = await fetch(`http://localhost:8081/followers?userId=${userId}`, { credentials: 'include' })
      const data = await res.json()
      followers.value = (res.ok && data.users) ? data.users : []
    } catch (e) {
      followers.value = []
    }
  }

  // Récupère la liste des following
  async function fetchFollowing(userId) {
    try {
      const res = await fetch(`http://localhost:8081/following?userId=${userId}`, { credentials: 'include' })
      const data = await res.json()
      following.value = (res.ok && data.users) ? data.users : []
    } catch (e) {
      following.value = []
    }
  }

  // Récupère les demandes de suivi en attente (pour profils privés)
  async function fetchFollowRequests() {
    try {
      const res = await fetch('http://localhost:8081/pendingFollowRequests', { credentials: 'include' })
      const data = await res.json()
      followRequests.value = (res.ok && data.requests) ? data.requests : []
    } catch (e) {
      followRequests.value = []
    }
  }

  // Suivre un utilisateur
  async function follow(userId) {
    try {
      const res = await fetch(`http://localhost:8081/follow?userId=${userId}`, {
        method: 'POST',
        credentials: 'include'
      })
      return await res.json()
    } catch (e) {
      return { error: 'Erreur réseau' }
    }
  }

  // Unfollow un utilisateur
  async function unfollow(userId) {
    try {
      const res = await fetch(`http://localhost:8081/unfollow?userId=${userId}`, {
        method: 'POST',
        credentials: 'include'
      })
      return await res.json()
    } catch (e) {
      return { error: 'Erreur réseau' }
    }
  }

  // Demander à suivre (pour profil privé)
  async function requestFollow(userId) {
    return follow(userId)
  }

  // Accepter ou refuser une demande de suivi (pour profils privés)
  async function respondToFollowRequest(requestId, response) {
    try {
      const res = await fetch('http://localhost:8081/responseFollowRequest', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ requestId, response })
      })
      return await res.json()
    } catch (e) {
      return { error: 'Erreur réseau' }
    }
  }

  // Vérifie si l'utilisateur courant suit un autre utilisateur
  function isFollowing(userId) {
    return following.value.some(f => f.id === userId)
  }

  // Vérifie si l'utilisateur courant est suivi par un autre utilisateur
  function isFollower(userId) {
    return followers.value.some(f => f.id === userId)
  }

  // Vérifie si une demande de suivi est en attente
  function isFollowRequested(userId) {
    return followRequests.value.some(r => r.id === userId)
  }

  return {
    followers,
    following,
    followRequests,
    fetchFollowers,
    fetchFollowing,
    fetchFollowRequests,
    follow,
    unfollow,
    requestFollow,
    respondToFollowRequest,
    isFollowing,
    isFollower,
    isFollowRequested
  }
})
