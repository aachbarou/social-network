<template>
  <div class="notification-test-panel">
    <div class="test-panel-header">
      <h3>🧪 Notification System Test Panel</h3>
      <p>Use this panel to test different notification types and features</p>
    </div>

    <div class="test-sections">
      <!-- Connection Status -->
      <div class="test-section">
        <h4>🔗 WebSocket Connection</h4>
        <div class="status-indicator">
          <div class="status-dot" :class="{ connected: wsConnected, disconnected: !wsConnected }"></div>
          <span>{{ wsConnected ? 'Connected' : 'Disconnected' }}</span>
          <button @click="toggleWebSocket" class="test-btn">
            {{ wsConnected ? 'Disconnect' : 'Reconnect' }}
          </button>
        </div>
      </div>

      <!-- Notification Stats -->
      <div class="test-section">
        <h4>📊 Notification Statistics</h4>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-number">{{ notificationStore.notifications.length }}</div>
            <div class="stat-label">Total Notifications</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">{{ notificationStore.unreadCount }}</div>
            <div class="stat-label">Unread</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">{{ notificationStore.followRequests.length }}</div>
            <div class="stat-label">Follow Requests</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">{{ notificationStore.groupInvites.length }}</div>
            <div class="stat-label">Group Invites</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">{{ notificationStore.groupRequests.length }}</div>
            <div class="stat-label">Group Requests</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">{{ notificationStore.eventNotifications.length }}</div>
            <div class="stat-label">Events</div>
          </div>
        </div>
      </div>

      <!-- Test Actions -->
      <div class="test-section">
        <h4>⚡ Test Actions</h4>
        <div class="test-buttons">
          <button @click="refreshNotifications" class="test-btn primary" :disabled="notificationStore.loading">
            <span v-if="notificationStore.loading" class="loading-spinner-sm"></span>
            {{ notificationStore.loading ? 'Refreshing...' : 'Refresh Notifications' }}
          </button>
          <button @click="markAllRead" class="test-btn secondary" :disabled="notificationStore.unreadCount === 0">
            Mark All as Read ({{ notificationStore.unreadCount }})
          </button>
          <button @click="simulateNotification" class="test-btn accent">
            Simulate New Notification
          </button>
          <button @click="clearError" class="test-btn error" v-if="notificationStore.error">
            Clear Error
          </button>
        </div>
      </div>

      <!-- Processing Status -->
      <div class="test-section" v-if="notificationStore.processingRequests.size > 0">
        <h4>⏳ Processing Requests</h4>
        <div class="processing-list">
          <div v-for="requestId in Array.from(notificationStore.processingRequests)" :key="requestId" class="processing-item">
            <div class="loading-spinner-sm"></div>
            <span>Processing request: {{ requestId.substring(0, 8) }}...</span>
          </div>
        </div>
      </div>

      <!-- Error Display -->
      <div class="test-section error-section" v-if="notificationStore.error">
        <h4>❌ Current Error</h4>
        <div class="error-message">{{ notificationStore.error }}</div>
      </div>

      <!-- Testing Guide -->
      <div class="test-section">
        <h4>📖 Testing Guide</h4>
        <div class="testing-guide">
          <div class="test-scenario">
            <h5>🤝 Follow Request Testing</h5>
            <ol>
              <li>Open two different browsers or incognito windows</li>
              <li>Log in as different users in each window</li>
              <li>Make one user's profile private in settings</li>
              <li>Try to follow the private user from the other account</li>
              <li>Check real-time notification in the private user's window</li>
            </ol>
          </div>

          <div class="test-scenario">
            <h5>👥 Group Invitation Testing</h5>
            <ol>
              <li>Log in as two different users</li>
              <li>Create a group with one user</li>
              <li>Invite the other user to the group</li>
              <li>Check notification in the invited user's window</li>
              <li>Test accept/decline functionality</li>
            </ol>
          </div>

          <div class="test-scenario">
            <h5>📝 Group Request Testing</h5>
            <ol>
              <li>Create a private group with one user</li>
              <li>Request to join the group from another user</li>
              <li>Check notification in the group admin's window</li>
              <li>Test approve/decline functionality</li>
            </ol>
          </div>

          <div class="test-scenario">
            <h5>📅 Event Notification Testing</h5>
            <ol>
              <li>Have two users join the same group</li>
              <li>Create an event in the group with one user</li>
              <li>Check event notification in the other member's window</li>
              <li>Test real-time delivery and notification actions</li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useNotificationStore } from '../stores/notificationStore'
import wsService from '../services/websocketService'

const notificationStore = useNotificationStore()
const wsConnected = ref(false)

let wsStatusInterval = null

// Check WebSocket connection status
const checkWebSocketStatus = () => {
  wsConnected.value = wsService.isConnected()
}

const toggleWebSocket = () => {
  if (wsConnected.value) {
    wsService.disconnect()
  } else {
    wsService.connect()
  }
}

const refreshNotifications = () => {
  notificationStore.fetchNotifications()
}

const markAllRead = () => {
  notificationStore.markAllAsRead()
}

const clearError = () => {
  notificationStore.dismissError()
}

const simulateNotification = () => {
  // Simulate a new notification for testing
  const mockNotification = {
    id: `test-${Date.now()}`,
    type: 'FOLLOW',
    user: {
      firstName: 'Test',
      lastName: 'User',
      nickname: 'testuser'
    },
    createdAt: new Date().toISOString(),
    read: false
  }
  notificationStore.handleNewNotification(mockNotification)
}

onMounted(() => {
  checkWebSocketStatus()
  wsStatusInterval = setInterval(checkWebSocketStatus, 1000)
})

onUnmounted(() => {
  if (wsStatusInterval) {
    clearInterval(wsStatusInterval)
  }
})
</script>

<style scoped>
.notification-test-panel {
  background: rgba(15, 15, 23, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1.5rem;
  margin: 1rem 0;
  backdrop-filter: blur(10px);
}

.test-panel-header {
  text-align: center;
  margin-bottom: 2rem;
}

.test-panel-header h3 {
  color: #e879c6;
  margin: 0;
  font-size: 1.5rem;
}

.test-panel-header p {
  color: #888;
  margin: 0.5rem 0 0 0;
}

.test-sections {
  display: grid;
  gap: 1.5rem;
}

.test-section {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  padding: 1rem;
}

.test-section h4 {
  color: #78dbff;
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.status-dot.connected {
  background: #10b981;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
}

.status-dot.disconnected {
  background: #ef4444;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.5);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 1rem;
}

.stat-card {
  background: rgba(120, 219, 255, 0.1);
  border: 1px solid rgba(120, 219, 255, 0.2);
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #78dbff;
}

.stat-label {
  color: #888;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.test-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.test-btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.test-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.test-btn.primary {
  background: linear-gradient(90deg, #78dbff 0%, #7879c6 100%);
  color: white;
}

.test-btn.secondary {
  background: #23233a;
  color: #78dbff;
  border: 1px solid rgba(120, 219, 255, 0.3);
}

.test-btn.accent {
  background: linear-gradient(90deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.test-btn.error {
  background: #ef4444;
  color: white;
}

.test-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  opacity: 0.9;
}

.processing-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.processing-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #888;
}

.error-section {
  border: 1px solid rgba(239, 68, 68, 0.3);
  background: rgba(239, 68, 68, 0.1);
}

.error-message {
  color: #ef4444;
  font-weight: 500;
}

.testing-guide {
  display: grid;
  gap: 1rem;
}

.test-scenario {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 6px;
  padding: 1rem;
}

.test-scenario h5 {
  color: #e879c6;
  margin: 0 0 0.5rem 0;
}

.test-scenario ol {
  color: #d1d5db;
  margin: 0;
  padding-left: 1.5rem;
}

.test-scenario li {
  margin-bottom: 0.25rem;
}

.loading-spinner-sm {
  border: 2px solid rgba(120, 219, 255, 0.1);
  border-radius: 50%;
  border-top: 2px solid #78dbff;
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
