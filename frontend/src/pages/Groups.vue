<template>
  <div class="groups-page">
    <!-- Header -->
    <GroupsHeader @create="showCreateModal = true" @refresh="loadData" :loading="loading" />

    <!-- Navigation Tabs -->
    <GroupsTabBar 
      :active-tab="activeTab" 
      :tabs="tabs"
      @tab-change="activeTab = $event"
    />

    <!-- Content -->
    <div class="content-section">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Chargement des groupes...</p>
      </div>

      <!-- My Groups Tab -->
      <div v-else-if="activeTab === 'my'" class="groups-grid">
        <div v-if="userGroups.length === 0" class="empty-state">
          <span class="empty-icon">👥</span>
          <h3>Aucun groupe</h3>
          <p>Vous n'êtes membre d'aucun groupe pour le moment</p>
          <button @click="showCreateModal = true" class="btn-primary">
            Créer votre premier groupe
          </button>
        </div>
        <GroupCard 
          v-for="group in userGroups" 
          :key="group.id"
          :group="group"
          :is-member="true"
          @view="openGroup"
          @chat="openGroupChat"
        />
      </div>

      <!-- Discover Tab -->
      <div v-else-if="activeTab === 'discover'" class="groups-grid">
        <div v-if="discoveryGroups.length === 0" class="empty-state">
          <span class="empty-icon">🔍</span>
          <h3>Aucun groupe public</h3>
          <p>Il n'y a pas de groupes publics disponibles</p>
        </div>
        <GroupCard 
          v-for="group in discoveryGroups" 
          :key="group.id"
          :group="group"
          :is-member="group.member"
          :loading="group.loading || false"
          @view="openGroup"
          @join-request="requestToJoin"
        />
      </div>

      <!-- Invitations Tab -->
      <div v-else-if="activeTab === 'invitations'" class="invitations-list">
        <div v-if="groupInvitations.length === 0" class="empty-state">
          <span class="empty-icon">✉️</span>
          <h3>Aucune invitation</h3>
          <p>Vous n'avez pas d'invitations de groupe en attente</p>
        </div>
        <InvitationCard 
          v-for="invitation in groupInvitations" 
          :key="invitation.id"
          :invitation="invitation"
          @respond="handleInvitationResponse"
        />
      </div>

      <!-- Requests Tab -->
      <div v-else-if="activeTab === 'requests'" class="requests-list">
        <div v-if="groupRequests.length === 0" class="empty-state">
          <span class="empty-icon">📋</span>
          <h3>Aucune demande</h3>
          <p>Aucune demande d'adhésion en attente pour vos groupes</p>
        </div>
        <RequestCard 
          v-for="request in groupRequests" 
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
import { ref, computed, onMounted, onActivated, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useGroupStore } from '../stores/groupStore'
import { storeToRefs } from 'pinia'

// Components
import GroupsHeader from '../components/groups/GroupsHeader.vue'
import GroupsTabBar from '../components/groups/GroupsTabBar.vue'
import GroupCard from '../components/groups/GroupCard.vue'
import InvitationCard from '../components/groups/InvitationCard.vue'
import RequestCard from '../components/groups/RequestCard.vue'
import CreateGroupModal from '../components/groups/CreateGroupModal.vue'

const router = useRouter()
const groupStore = useGroupStore()

// State
const activeTab = ref('my')
const showCreateModal = ref(false)

// Store refs - properly reactive
const {
  userGroups,
  publicGroups,
  groupInvitations,
  groupRequests,
  discoveryGroups,
  loading
} = storeToRefs(groupStore)

// Computed
const tabs = computed(() => [
  {
    id: 'my',
    label: 'Mes groupes',
    icon: '👥',
    count: userGroups.value.length
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
    count: groupInvitations.value.length
  },
  {
    id: 'requests',
    label: 'Demandes',
    icon: '📋',
    count: groupRequests.value.length
  }
])

// Methods
const loadData = async () => {
  console.log('Loading groups data...')
  try {
    await groupStore.loadGroupsData()
    console.log('Data loaded:', {
      userGroups: userGroups.value?.length || 0,
      publicGroups: publicGroups.value?.length || 0,
      invitations: groupInvitations.value?.length || 0,
      requests: groupRequests.value?.length || 0
    })
    // Force reactivity update
    await nextTick()
  } catch (error) {
    console.error('Error loading groups data:', error)
  }
}

const openGroup = (group) => {
  router.push(`/group/${group.id}`)
}

const requestToJoin = async (groupId) => {
  try {
    // Find the group to check its privacy
    const group = publicGroups.value.find(g => g.id === groupId)
    if (!group) return
    
    // Set loading state using store method
    groupStore.setGroupLoading(groupId, true)
    
    const isPublic = group.privacy === 'public'
    const result = await groupStore.joinGroup(groupId, isPublic)
    
    if (result.success) {
      if (result.isPublic) {
        // Successfully joined public group - add small delay to ensure backend operation completes
        await new Promise(resolve => setTimeout(resolve, 500))
        // Then redirect to group page
        router.push(`/group/${groupId}`)
      }
      // For private groups, the store already handles the UI update
    }
  } catch (error) {
    console.error('Error joining group:', error)
  } finally {
    // Reset loading state using store method
    groupStore.setGroupLoading(groupId, false)
  }
}

const handleGroupCreated = async () => {
  showCreateModal.value = false
  console.log('Group created, refreshing data...')
  // Reload data to show the new group
  await loadData()
  activeTab.value = 'my'
}

const handleInvitationResponse = async (invitationId, action) => {
  const success = await groupStore.respondToInvitation(invitationId, action)
  if (!success) {
    console.error('Failed to respond to invitation')
  }
}

const handleRequestResponse = async (requestId, action) => {
  const success = await groupStore.respondToRequest(requestId, action)
  if (!success) {
    console.error('Failed to respond to request')
  }
}

const openGroupChat = (groupId) => {
  // Navigate to chat page with group context
  console.log('Opening group chat for group:', groupId)
  // TODO: Implement group-specific chat or integrate with existing chat system
  // For now, just navigate to the general chat page
  // router.push(`/chat?group=${groupId}`)
  
  // Simple alert for now to show functionality
  alert(`Chat du groupe ${groupId} - Fonctionnalité à implémenter`)
}

// Lifecycle
onMounted(() => {
  console.log('Groups component mounted, loading data...')
  loadData()
})

// Reload data when navigating back to this page
onActivated(() => {
  console.log('Groups component activated, reloading data...')
  loadData()
})

// Watch for changes in store data
watch([userGroups, publicGroups], ([newUserGroups, newPublicGroups]) => {
  console.log('Store data changed:', {
    userGroups: newUserGroups?.length || 0,
    publicGroups: newPublicGroups?.length || 0
  })
}, { deep: true })
</script>

<style scoped>
.groups-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
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

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
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
  
  .groups-grid {
    grid-template-columns: 1fr;
  }
}
</style>
