<template>
  <div class="groups-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-content">
        <h1>Groupes</h1>
        <p>Découvrez et rejoignez des communautés</p>
      </div>
      <div class="header-actions">
        <button @click="showCreateModal = true" class="create-btn">
          <span class="icon">+</span>
          Créer un groupe
        </button>
        <button @click="loadData" class="refresh-btn" :disabled="loading">
          <span class="icon" :class="{ 'animate-spin': loading }">↻</span>
          Actualiser
        </button>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="tabs-container">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
      >
        <span class="icon">{{ tab.icon }}</span>
        {{ tab.label }}
        <span v-if="tab.count > 0" class="tab-count">{{ tab.count }}</span>
      </button>
    </div>

    <!-- Content -->
    <div class="content-section">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement des groupes...</p>
      </div>

      <!-- My Groups Tab -->
      <div v-else-if="activeTab === 'my'" class="groups-grid">
        <div v-if="myGroups.length === 0" class="empty-state">
          <span class="empty-icon">👥</span>
          <h3>Aucun groupe</h3>
          <p>Vous n'êtes membre d'aucun groupe pour le moment</p>
          <button @click="showCreateModal = true" class="btn-primary">
            Créer votre premier groupe
          </button>
        </div>
        <GroupCard 
          v-for="group in myGroups" 
          :key="group.id"
          :group="group"
          :is-member="true"
          @view="openGroup"
          @chat="openGroupChat"
        />
      </div>

      <!-- Discover Tab -->
      <div v-else-if="activeTab === 'discover'" class="groups-grid">
        <div v-if="publicGroups.length === 0" class="empty-state">
          <span class="empty-icon">🔍</span>
          <h3>Aucun groupe public</h3>
          <p>Il n'y a pas de groupes publics disponibles</p>
        </div>
        <GroupCard 
          v-for="group in publicGroups" 
          :key="group.id"
          :group="group"
          :is-member="group.member"
          @view="openGroup"
          @join-request="requestToJoin"
        />
      </div>

      <!-- Invitations Tab -->
      <div v-else-if="activeTab === 'invitations'" class="invitations-list">
        <div v-if="invitations.length === 0" class="empty-state">
          <span class="empty-icon">✉️</span>
          <h3>Aucune invitation</h3>
          <p>Vous n'avez pas d'invitations de groupe en attente</p>
        </div>
        <InvitationCard 
          v-for="invitation in invitations" 
          :key="invitation.id"
          :invitation="invitation"
          @respond="handleInvitationResponse"
        />
      </div>

      <!-- Requests Tab -->
      <div v-else-if="activeTab === 'requests'" class="requests-list">
        <div v-if="requests.length === 0" class="empty-state">
          <span class="empty-icon">📋</span>
          <h3>Aucune demande</h3>
          <p>Aucune demande d'adhésion en attente pour vos groupes</p>
        </div>
        <RequestCard 
          v-for="request in requests" 
          :key="request.id"
          :request="request"
          @respond="handleRequestResponse"
        />
      </div>
    </div>

    <!-- Create Group Modal -->
    <CreateGroupModal 
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="handleGroupCreated"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// Components
import GroupCard from '../components/groups/GroupCard.vue'
import InvitationCard from '../components/groups/InvitationCard.vue'
import RequestCard from '../components/groups/RequestCard.vue'
import CreateGroupModal from '../components/groups/CreateGroupModal.vue'

const router = useRouter()

// State
const activeTab = ref('my')
const showCreateModal = ref(false)
const loading = ref(false)

// Mock data - replace with real API calls
const myGroups = ref([])
const publicGroups = ref([])
const invitations = ref([])
const requests = ref([])

// Computed
const tabs = computed(() => [
  {
    id: 'my',
    label: 'Mes groupes',
    icon: '👥',
    count: myGroups.value.length
  },
  {
    id: 'discover',
    label: 'Découvrir',
    icon: '🔍',
    count: 0
  },
  {
    id: 'invitations',
    label: 'Invitations',
    icon: '✉️',
    count: invitations.value.length
  },
  {
    id: 'requests',
    label: 'Demandes',
    icon: '📋',
    count: requests.value.length
  }
])

// Methods
const loadData = async () => {
  loading.value = true
  try {
    // Load user groups
    const userGroupsRes = await fetch('http://localhost:8081/userGroups', {
      credentials: 'include'
    })
    if (userGroupsRes.ok) {
      const data = await userGroupsRes.json()
      myGroups.value = data.groups || []
    } else {
      console.error('Failed to load user groups:', userGroupsRes.status)
    }

    // Load all groups for discovery
    const allGroupsRes = await fetch('http://localhost:8081/allGroups', {
      credentials: 'include'
    })
    if (allGroupsRes.ok) {
      const data = await allGroupsRes.json()
      publicGroups.value = (data.groups || []).filter(g => !g.member)
    } else {
      console.error('Failed to load all groups:', allGroupsRes.status)
    }

    // Load invitations and requests from notifications
    await loadNotifications()
    
  } catch (error) {
    console.error('Error loading groups:', error)
  } finally {
    loading.value = false
  }
}

const loadNotifications = async () => {
  try {
    const notificationsRes = await fetch('http://localhost:8081/notifications', {
      credentials: 'include'
    })
    if (notificationsRes.ok) {
      const data = await notificationsRes.json()
      // Filter for group invitations
      invitations.value = (data.notifications || [])
        .filter(notif => notif.type === 'GROUP_INVITE')
        .map(notif => ({
          id: notif.id,
          group: notif.group || { name: 'Groupe', image: null, memberCount: 0 },
          inviter: notif.sender || { firstName: 'Utilisateur', lastName: '' },
          createdAt: notif.createdAt
        }))
      
      // Filter for group requests (if user is admin)
      requests.value = (data.notifications || [])
        .filter(notif => notif.type === 'GROUP_REQUEST')
        .map(notif => ({
          id: notif.id,
          user: notif.sender || { firstName: 'Utilisateur', lastName: '', profilePic: null },
          group: notif.group || { name: 'Groupe' },
          createdAt: notif.createdAt
        }))
    } else {
      console.error('Failed to load notifications:', notificationsRes.status)
      // Set empty arrays if notifications endpoint fails
      invitations.value = []
      requests.value = []
    }
  } catch (error) {
    console.error('Error loading notifications:', error)
    // Set empty arrays on error
    invitations.value = []
    requests.value = []
  }
}

const openGroup = (group) => {
  router.push(`/group/${group.id}`)
}

const requestToJoin = async (groupId) => {
  try {
    const response = await fetch('http://localhost:8081/newGroupRequest', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ groupId })
    })
    
    if (response.ok) {
      // Update UI to show request sent
      const group = publicGroups.value.find(g => g.id === groupId)
      if (group) {
        group.requestPending = true
      }
      
      // Show success feedback
      alert('Demande d\'adhésion envoyée avec succès!')
    } else {
      throw new Error('Failed to send join request')
    }
  } catch (error) {
    console.error('Error sending join request:', error)
    alert('Erreur lors de l\'envoi de la demande d\'adhésion')
  }
}

const handleGroupCreated = () => {
  showCreateModal.value = false
  // Reload data to show the new group
  loadData()
  activeTab.value = 'my'
}

const handleInvitationResponse = async (invitationId, action) => {
  try {
    const response = await fetch('http://localhost:8081/responseInviteRequest', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        notificationId: invitationId, 
        response: action === 'accept' ? 'ACCEPT' : 'DECLINE' 
      })
    })
    
    if (response.ok) {
      // Remove invitation from list
      invitations.value = invitations.value.filter(inv => inv.id !== invitationId)
      if (action === 'accept') {
        // Reload groups data
        loadData()
      }
    }
  } catch (error) {
    console.error('Error responding to invitation:', error)
  }
}

const handleRequestResponse = async (requestId, action) => {
  try {
    const response = await fetch('http://localhost:8081/responseGroupRequest', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        notificationId: requestId, 
        response: action === 'accept' ? 'ACCEPT' : 'DECLINE' 
      })
    })
    
    if (response.ok) {
      // Remove request from list
      requests.value = requests.value.filter(req => req.id !== requestId)
      if (action === 'accept') {
        // Reload groups data to update member counts
        loadData()
      }
    }
  } catch (error) {
    console.error('Error responding to request:', error)
  }
}

const openGroupChat = (groupId) => {
  // Navigate to chat page with group context
  console.log('Opening group chat for group:', groupId)
  // TODO: Implement group-specific chat or integrate with existing chat system
  // For now, just navigate to the general chat page
  // this.$router.push(`/chat?group=${groupId}`)
  
  // Simple alert for now to show functionality
  alert(`Chat du groupe ${groupId} - Fonctionnalité à implémenter`)
}

// Lifecycle
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.groups-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 30px;
  background: rgba(15, 15, 23, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.page-header > div:last-child {
  display: flex;
  gap: 12px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.header-content h1 {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.header-content p {
  color: rgba(255, 255, 255, 0.7);
  margin: 8px 0 0 0;
  font-size: 1.1rem;
}

.create-btn {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 15px;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(232, 121, 198, 0.4);
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 15px;
  padding: 12px 20px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  transform: translateY(-1px);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

.icon {
  font-size: 1rem;
  font-weight: bold;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.tabs-container {
  display: flex;
  gap: 4px;
  margin-bottom: 30px;
  background: rgba(15, 15, 23, 0.6);
  padding: 8px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.tab-btn {
  flex: 1;
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  padding: 12px 20px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 500;
  position: relative;
}

.tab-btn.active {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.tab-count {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 0.75rem;
  font-weight: 700;
  min-width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content-section {
  min-height: 400px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.loading-state p {
  color: rgba(255, 255, 255, 0.7);
  margin: 16px 0 0 0;
  font-size: 1.1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top: 3px solid #e879c6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.invitations-list,
.requests-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: rgba(15, 15, 23, 0.4);
  border-radius: 20px;
  border: 2px dashed rgba(255, 255, 255, 0.1);
  grid-column: 1 / -1;
}

.empty-state h3 {
  color: rgba(255, 255, 255, 0.8);
  margin: 16px 0 8px 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.empty-state p {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 24px 0;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 10px 20px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

/* Responsive */
@media (max-width: 768px) {
  .groups-page {
    padding: 15px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .tabs-container {
    flex-direction: column;
  }
  
  .groups-grid {
    grid-template-columns: 1fr;
  }
}
</style>
