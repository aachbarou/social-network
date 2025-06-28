import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useFollowStore = defineStore('follow', () => {
  const followRequests = ref([])
  const following = ref([])
  const followers = ref([])

  // Charger la liste des abonnements (following)
  async function fetchFollowing(userId) {
    const res = await fetch(`http://localhost:8081/following?userId=${userId}`, { credentials: 'include' })
    const data = await res.json()
    if (res.ok && data.users) following.value = data.users.map(u => u.id)
  }
  // Charger la liste des followers
  async function fetchFollowers(userId) {
    const res = await fetch(`http://localhost:8081/followers?userId=${userId}`, { credentials: 'include' })
    const data = await res.json()
    if (res.ok && data.users) followers.value = data.users.map(u => u.id)
  }

  // Demander à suivre ou suivre directement
  async function follow(userId) {
    await fetch('http://localhost:8081/follow', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userId })
    })
    await fetchFollowing(userId)
    await fetchFollowers(userId)
  }

  // Demander à suivre (pour profil privé)
  async function requestFollow(userId) {
    // Même endpoint que follow, le backend gère la logique
    await follow(userId)
  }

  // Unfollow
  async function unfollow(userId) {
    await fetch('http://localhost:8081/unfollow', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userId })
    })
    await fetchFollowing(userId)
    await fetchFollowers(userId)
  }

  // Accepter/refuser une demande reçue
  async function respondToRequest(userId, accept) {
    await fetch('http://localhost:8081/responseFollowRequest', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ userId, accept })
    })
  }

  // Getters (à adapter selon la logique de ton backend)
  function isFollowing(userId) {
    return following.value.includes(userId)
  }
  function isFollowRequested(userId) {
    return followRequests.value.some(r => r.userId === userId && r.status === 'pending')
  }

  return {
    followRequests,
    following,
    followers,
    follow,
    requestFollow,
    unfollow,
    respondToRequest,
    isFollowing,
    isFollowRequested,
    fetchFollowing,
    fetchFollowers
  }
})
