<template>
  <div class="debug-panel">
    <h3>🔧 Debug Information</h3>
    
    <div class="debug-section">
      <h4>📊 Notification Store State</h4>
      <div class="debug-data">
        <pre>{{ JSON.stringify(debugData, null, 2) }}</pre>
      </div>
    </div>

    <div class="debug-section">
      <h4>🧪 Create Test Notification</h4>
      <div class="test-controls">
        <select v-model="testNotifType" class="debug-select">
          <option value="FOLLOW">Follow Request</option>
          <option value="GROUP_INVITE">Group Invitation</option>
          <option value="GROUP_REQUEST">Group Request</option>
          <option value="EVENT">Event Notification</option>
        </select>
        <button @click="createTestNotification" class="debug-btn">Create Test Notification</button>
      </div>
    </div>

    <div class="debug-section">
      <h4>🌐 API Test</h4>
      <div class="api-controls">
        <button @click="testNotificationAPI" class="debug-btn">Test Notification API</button>
        <button @click="testUserAPI" class="debug-btn">Test User API</button>
      </div>
      <div v-if="apiResponse" class="api-response">
        <h5>API Response:</h5>
        <pre>{{ JSON.stringify(apiResponse, null, 2) }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useNotificationStore } from '../stores/notificationStore'
import { useUserStore } from '../stores/userStore'

const notificationStore = useNotificationStore()
const userStore = useUserStore()
const testNotifType = ref('FOLLOW')
const apiResponse = ref(null)

const debugData = computed(() => ({
  notifications: {
    total: notificationStore.notifications.length,
    list: notificationStore.notifications,
    followRequests: notificationStore.followRequests.length,
    groupInvites: notificationStore.groupInvites.length,
    groupRequests: notificationStore.groupRequests.length,
    eventNotifications: notificationStore.eventNotifications.length
  },
  loading: notificationStore.loading,
  error: notificationStore.error,
  unreadCount: notificationStore.unreadCount,
  hasNewNotifications: notificationStore.hasNewNotifications,
  processingRequests: Array.from(notificationStore.processingRequests),
  currentUser: userStore.user
}))

const createTestNotification = () => {
  const mockNotification = {
    id: `test-${Date.now()}`,
    type: testNotifType.value,
    createdAt: new Date().toISOString(),
    read: false,
    user: {
      id: 'test-user',
      firstName: 'Test',
      lastName: 'User',
      nickname: 'testuser',
      avatar: null
    },
    group: testNotifType.value.includes('GROUP') ? {
      id: 'test-group',
      name: 'Test Group',
      description: 'A test group for debugging'
    } : null,
    event: testNotifType.value === 'EVENT' ? {
      id: 'test-event',
      title: 'Test Event',
      description: 'A test event for debugging',
      date: new Date().toISOString()
    } : null
  }
  
  console.log('🧪 Creating test notification:', mockNotification)
  notificationStore.handleNewNotification(mockNotification)
}

const testNotificationAPI = async () => {
  try {
    console.log('🌐 Testing notification API...')
    const response = await fetch('http://localhost:8081/notifications', {
      credentials: 'include'
    })
    const data = await response.json()
    apiResponse.value = {
      status: response.status,
      data: data
    }
    console.log('📬 API Response:', apiResponse.value)
  } catch (error) {
    apiResponse.value = {
      error: error.message
    }
    console.error('❌ API Error:', error)
  }
}

const testUserAPI = async () => {
  try {
    console.log('🌐 Testing user API...')
    const response = await fetch('http://localhost:8081/user', {
      credentials: 'include'
    })
    const data = await response.json()
    apiResponse.value = {
      status: response.status,
      data: data
    }
    console.log('👤 User API Response:', apiResponse.value)
  } catch (error) {
    apiResponse.value = {
      error: error.message
    }
    console.error('❌ User API Error:', error)
  }
}

onMounted(() => {
  console.log('🔧 Debug panel mounted')
})
</script>

<style scoped>
.debug-panel {
  background: rgba(239, 68, 68, 0.1);
  border: 2px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  padding: 1rem;
  margin: 1rem 0;
  font-family: 'Courier New', monospace;
}

.debug-section {
  margin-bottom: 1.5rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  padding: 1rem;
}

.debug-section h4 {
  color: #ef4444;
  margin: 0 0 1rem 0;
  font-size: 1rem;
}

.debug-data {
  background: rgba(0, 0, 0, 0.4);
  border-radius: 4px;
  padding: 1rem;
  max-height: 300px;
  overflow-y: auto;
}

.debug-data pre {
  color: #10b981;
  font-size: 12px;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.test-controls, .api-controls {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.debug-select {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 4px;
  color: white;
  padding: 0.5rem;
}

.debug-btn {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.5);
  border-radius: 4px;
  color: #ef4444;
  padding: 0.5rem 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.debug-btn:hover {
  background: rgba(239, 68, 68, 0.3);
  border-color: rgba(239, 68, 68, 0.7);
}

.api-response {
  margin-top: 1rem;
  background: rgba(0, 0, 0, 0.4);
  border-radius: 4px;
  padding: 1rem;
}

.api-response h5 {
  color: #10b981;
  margin: 0 0 0.5rem 0;
}

.api-response pre {
  color: #78dbff;
  font-size: 12px;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
