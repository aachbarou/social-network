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
      const res = await fetch('http://localhost:8081/notifications', { credentials: 'include' })
      const data = await res.json()
      // Filtrer les notifications de type FOLLOW
      followRequests.value = (res.ok && data.notifications)
        ? data.notifications.filter(n => n.type === 'FOLLOW')
        : []

        
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
      const result = await res.json()
      
      // If follow was successful, trigger posts refresh to show new posts
      if (res.ok && result.type === 'Success') {
        // Import and use posts store to refresh the feed
        const { useMainStore } = await import('./postsStore')
        const postsStore = useMainStore()
        await postsStore.fetchPosts()
      }
      
      return result      
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
     const result = await res.json()
      
      // If unfollow was successful, trigger posts refresh
      if (res.ok && result.type === 'Success') {
        // Import and use posts store to refresh the feed
        const { useMainStore } = await import('./postsStore')
        const postsStore = useMainStore()
        await postsStore.fetchPosts()
      }
      
      return result    } catch (e) {
      return { error: 'Erreur réseau' }
    }
  }

  // Demander à suivre (pour profil privé)
  async function requestFollow(userId) {
    const result = await follow(userId)
    // Si la demande a été envoyée avec succès, on ajoute une demande en attente localement
    if (!isFollowRequested(userId)) {
      followRequests.value.push({
        type: 'FOLLOW',
        content: userId,
        // Ajoute d'autres champs si besoin (id, sender, etc.)
      })
    }
    return result
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
      const result = await res.json()
      
      // If follow request was accepted, refresh posts to show new posts
      if (res.ok && result.type === 'Success' && response === 'ACCEPT') {
        // Import and use posts store to refresh the feed
        const { useMainStore } = await import('./postsStore')
        const postsStore = useMainStore()
        await postsStore.fetchPosts()
      }
      
      return result    } catch (e) {
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
    // Cherche une notification de type FOLLOW pour cet utilisateur
    return followRequests.value.some(req => {
      // Selon la structure, req.content ou req.sender ou req.targetId peut contenir l'id
      // On vérifie les deux cas courants
      return (req.content === userId || req.sender === userId || req.userId === userId)
    })
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
