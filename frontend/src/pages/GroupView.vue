<template>
  <div class="group-view">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du groupe...</p>
    </div>

    <!-- Group Content -->
    <div v-else-if="group" class="group-content">
      <!-- Header -->
      <div class="group-header">
        <div class="group-info">
          <div class="group-avatar">
            <img v-if="group.image" :src="group.image" :alt="group.name" />
            <div v-else class="default-avatar">
              <span class="group-icon">👥</span>
            </div>
          </div>
          <div class="group-details">
            <h1>{{ group.name }}</h1>
            <p v-if="group.description">{{ group.description }}</p>
            <div class="group-meta">
              <span class="member-count">
                <span class="icon">👤</span>
                {{ group.memberCount || 0 }} membres
              </span>
              <span class="group-privacy">
                <span class="icon">{{ group.privacy === 'public' ? '🌐' : '🔒' }}</span>
                {{ group.privacy === 'public' ? 'Public' : 'Privé' }}
              </span>
            </div>
          </div>
        </div>
        <div class="group-actions">
          <!-- Join button for non-members of public groups -->
          <button 
            v-if="canJoinGroup"
            @click="joinGroup" 
            class="btn-primary"
            :disabled="isJoining"
          >
            <span 
              class="icon" 
              :class="{ 'animate-spin': isJoining }"
            >
              {{ isJoining ? '⏳' : '➕' }}
            </span>
            {{ isJoining ? 'Rejoindre...' : (group.privacy === 'public' ? 'Rejoindre le groupe' : 'Demander à rejoindre') }}
          </button>
          
          <!-- Chat button for members -->
          <button 
            v-if="userMembership.isMember || userMembership.isAdmin"
            @click="openChat" 
            class="btn-primary"
          >
            <span class="icon">💬</span>
            Chat
          </button>
          
          <button @click="$router.go(-1)" class="btn-secondary">
            <span class="icon">←</span>
            Retour
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
        </button>
      </div>

      <!-- Tab Content -->
      <div class="tab-content">
        <!-- Posts Tab -->
        <div v-if="activeTab === 'posts'" class="posts-section">
          <div v-if="posts.length === 0" class="empty-state">
            <span class="empty-icon">📝</span>
            <h3>Aucun post</h3>
            <p>Soyez le premier à publier dans ce groupe</p>
          </div>
          <div v-else class="posts-list">
            <div v-for="post in posts" :key="post.id" class="post-card">
              <div class="post-header">
                <div class="author-avatar">
                  <img v-if="post.author.profilePic" :src="post.author.profilePic" :alt="post.author.name" />
                  <span v-else class="default-author">👤</span>
                </div>
                <div class="post-info">
                  <h4>{{ post.author.firstName }} {{ post.author.lastName }}</h4>
                  <span class="post-date">{{ formatDate(post.createdAt) }}</span>
                </div>
              </div>
              <div class="post-content">
                <p>{{ post.content }}</p>
                <img v-if="post.image" :src="post.image" :alt="post.content" class="post-image" />
              </div>
            </div>
          </div>
        </div>

        <!-- Members Tab -->
        <div v-if="activeTab === 'members'" class="members-section">
          <div v-if="members.length === 0" class="empty-state">
            <span class="empty-icon">👥</span>
            <h3>Aucun membre</h3>
            <p>Ce groupe n'a pas encore de membres</p>
          </div>
          <div v-else class="members-grid">
            <div v-for="member in members" :key="member.id" class="member-card">
              <div class="member-avatar">
                <img v-if="member.profilePic" :src="member.profilePic" :alt="member.name" />
                <span v-else class="default-member">👤</span>
              </div>
              <div class="member-info">
                <h4>{{ member.firstName }} {{ member.lastName }} <span v-if="member.following" class="admin-crown">👑</span></h4>
              </div>
            </div>
          </div>
        </div>

        <!-- Events Tab -->
        <div v-if="activeTab === 'events'" class="events-section">
          <div v-if="events.length === 0" class="empty-state">
            <span class="empty-icon">📅</span>
            <h3>Aucun événement</h3>
            <p>Aucun événement prévu pour ce groupe</p>
          </div>
          <div v-else class="events-list">
            <div v-for="event in events" :key="event.id" class="event-card">
              <div class="event-date">
                <span class="day">{{ formatEventDay(event.date) }}</span>
                <span class="month">{{ formatEventMonth(event.date) }}</span>
              </div>
              <div class="event-info">
                <h4>{{ event.title }}</h4>
                <p>{{ event.description }}</p>
                <span class="event-time">
                  <span class="icon">🕐</span>
                  {{ formatEventTime(event.date) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else class="error-state">
      <span class="error-icon">❌</span>
      <h3>Groupe introuvable</h3>
      <p>Ce groupe n'existe pas ou vous n'avez pas accès</p>
      <button @click="$router.push('/groups')" class="btn-primary">
        Retour aux groupes
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// State
const loading = ref(true)
const group = ref(null)
const activeTab = ref('posts')
const posts = ref([])
const members = ref([])
const events = ref([])
const isJoining = ref(false)
const userMembership = ref({
  isMember: false,
  isAdmin: false
})

// Tabs
const tabs = computed(() => [
  { id: 'posts', label: 'Posts', icon: '📝' },
  { id: 'members', label: 'Membres', icon: '👥' },
  { id: 'events', label: 'Événements', icon: '📅' }
])

// Computed
const canJoinGroup = computed(() => {
  return group.value && 
         (group.value.privacy === 'public' || group.value.privacy === 'private') &&
         !userMembership.value.isMember && 
         !userMembership.value.isAdmin
})

// Methods
const loadGroup = async () => {
  loading.value = true
  try {
    const groupId = route.params.id
    
    // Load group details
    const groupRes = await fetch(`http://localhost:8081/groupInfo?groupId=${groupId}`, {
      credentials: 'include'
    })
    
    if (groupRes.ok) {
      const data = await groupRes.json()
      // Backend returns groups array, take the first one
      group.value = data.groups && data.groups.length > 0 ? data.groups[0] : null
      
      if (group.value) {
        // Set membership status from the group data itself
        userMembership.value.isMember = group.value.member || false
        userMembership.value.isAdmin = group.value.administrator || false
        
        // Load group content
        await loadPosts()
        await loadMembers()
        await loadEvents()
      }
    } else {
      group.value = null
    }
  } catch (error) {
    console.error('Error loading group:', error)
    group.value = null
  } finally {
    loading.value = false
  }
}

const loadPosts = async () => {
  try {
    const postsRes = await fetch(`http://localhost:8081/groupPosts?groupId=${route.params.id}`, {
      credentials: 'include'
    })
    if (postsRes.ok) {
      const data = await postsRes.json()
      posts.value = data.posts || []
    }
  } catch (error) {
    console.error('Error loading posts:', error)
    posts.value = []
  }
}

const loadMembers = async () => {
  try {
    const membersRes = await fetch(`http://localhost:8081/groupMembers?groupId=${route.params.id}`, {
      credentials: 'include'
    })
    if (membersRes.ok) {
      const data = await membersRes.json()
      members.value = data.users || []
    }
  } catch (error) {
    console.error('Error loading members:', error)
    members.value = []
  }
}

const loadEvents = async () => {
  try {
    const eventsRes = await fetch(`http://localhost:8081/groupEvents?groupId=${route.params.id}`, {
      credentials: 'include'
    })
    if (eventsRes.ok) {
      const data = await eventsRes.json()
      events.value = data.events || []
    }
  } catch (error) {
    console.error('Error loading events:', error)
    events.value = []
  }
}

const openChat = () => {
  alert(`Chat du groupe ${group.value.name} - Fonctionnalité à implémenter`)
}

const joinGroup = async () => {
  if (!group.value || isJoining.value) return
  
  console.log('Starting join process, setting isJoining to true')
  isJoining.value = true
  
  try {
    // Add minimum delay to ensure animation is visible
    const minDelay = new Promise(resolve => setTimeout(resolve, 1000))
    
    // For public groups, join directly. For private groups, send a request.
    const endpoint = group.value.privacy === 'public' 
      ? 'http://localhost:8081/joinPublicGroup'
      : 'http://localhost:8081/newGroupRequest'
    
    console.log('Making request to:', endpoint)
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({
        groupId: group.value.id
      })
    })

    // Wait for both the request and minimum delay
    await minDelay

    if (response.ok) {
      if (group.value.privacy === 'public') {
        // Successfully joined public group - update membership status and refresh
        userMembership.value.isMember = true
        group.value.member = true
        // Refresh the group data to get updated member count and members list
        await loadGroup()
      } else {
        // For private groups, the request was sent silently
        // No intrusive popup needed
      }
    }
  } catch (error) {
    console.error('Error joining group:', error)
  } finally {
    console.log('Join process finished, setting isJoining to false')
    isJoining.value = false
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatEventDay = (dateString) => {
  const date = new Date(dateString)
  return date.getDate()
}

const formatEventMonth = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', { month: 'short' })
}

const formatEventTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('fr-FR', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Lifecycle
onMounted(() => {
  loadGroup()
})
</script>

<style scoped>
.group-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
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

.animate-spin {
  animation: spin 1s linear infinite !important;
}

.error-icon,
.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.group-header {
  background: rgba(15, 15, 23, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 30px;
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-info {
  display: flex;
  gap: 20px;
  align-items: center;
}

.group-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.group-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.group-icon {
  font-size: 2.5rem;
  opacity: 0.8;
}

.group-details h1 {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
}

.group-details p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 12px 0;
}

.group-meta {
  display: flex;
  gap: 16px;
}

.member-count,
.group-privacy {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.group-actions {
  display: flex;
  gap: 12px;
}

.btn-primary,
.btn-secondary {
  padding: 12px 20px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-primary:hover,
.btn-secondary:hover {
  transform: translateY(-1px);
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
}

.tab-btn.active {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  transform: translateY(-1px);
}

.tab-content {
  min-height: 400px;
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
}

.empty-state h3 {
  color: rgba(255, 255, 255, 0.8);
  margin: 16px 0 8px 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.empty-state p {
  color: rgba(255, 255, 255, 0.5);
  margin: 0;
}

.posts-list,
.events-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.post-card,
.event-card {
  background: rgba(15, 15, 23, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
}

.post-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-author {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
}

.post-info h4 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 0.875rem;
  font-weight: 600;
}

.post-date {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.75rem;
}

.post-content p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 16px 0;
  line-height: 1.6;
}

.post-image {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  border-radius: 12px;
}

.members-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.member-card {
  background: rgba(15, 15, 23, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 16px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.member-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.member-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-member {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.member-info h4 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 1rem;
  font-weight: 600;
}

.admin-crown {
  color: #ffd700;
  font-size: 0.9rem;
  margin-left: 4px;
}

.member-role {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
}

.event-card {
  display: flex;
  gap: 20px;
  align-items: center;
}

.event-date {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border-radius: 12px;
  padding: 12px;
  text-align: center;
  color: white;
  flex-shrink: 0;
}

.event-date .day {
  display: block;
  font-size: 1.5rem;
  font-weight: 700;
}

.event-date .month {
  display: block;
  font-size: 0.875rem;
  opacity: 0.8;
}

.event-info h4 {
  color: white;
  margin: 0 0 8px 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.event-info p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 8px 0;
}

.event-time {
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
}

.icon {
  font-size: 0.875rem;
}

/* Responsive */
@media (max-width: 768px) {
  .group-view {
    padding: 15px;
  }
  
  .group-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .group-info {
    flex-direction: column;
    text-align: center;
  }
  
  .tabs-container {
    flex-direction: column;
  }
  
  .members-grid {
    grid-template-columns: 1fr;
  }
  
  .event-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>
