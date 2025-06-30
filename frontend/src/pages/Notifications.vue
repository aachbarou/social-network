<template>
  <div>
    <h1>Notifications</h1>
    <div v-if="loading">Chargement...</div>
    <div v-else>
      <div v-if="followRequests.length === 0">Aucune demande de suivi.</div>
      <ul>
        <li v-for="notif in followRequests" :key="notif.id" style="margin-bottom: 0.5rem;">
          <span>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }} veut vous suivre.</span>
          <button @click="respondToRequest(notif.id, 'ACCEPT')">Accepter</button>
          <button @click="respondToRequest(notif.id, 'DECLINE')">Refuser</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const followRequests = ref([])
const loading = ref(true)
let ws = null

async function fetchNotifications() {
  loading.value = true
  try {
    const res = await fetch('http://localhost:8081/notifications', { credentials: 'include' })
    const data = await res.json()
    // Filtrer les notifications de type FOLLOW
    followRequests.value = (data.notifications || []).filter(n => n.type === 'FOLLOW')
  } catch (e) {
    followRequests.value = []
  } finally {
    loading.value = false
  }
}

async function respondToRequest(requestId, response) {
  await fetch('http://localhost:8081/responseFollowRequest', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ requestId, response })
  })
  await fetchNotifications()
}

function setupWebSocket() {
  // Remplace l'URL par celle de ton backend si besoin
  ws = new WebSocket('ws://localhost:8081/ws')
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      // Si c'est une notification, on rafraîchit la liste
      if (data.type === 'NOTIFICATION' || data.type === 'FOLLOW') {
        fetchNotifications()
      }
    } catch (e) {}
  }
  ws.onclose = () => {
    // Reconnexion simple après 2s si la connexion est perdue
    setTimeout(setupWebSocket, 2000)
  }
}

onMounted(() => {
  fetchNotifications()
  setupWebSocket()
})
onUnmounted(() => {
  if (ws) ws.close()
})
</script>
