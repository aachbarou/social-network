<template>
  <div class="main-layout">
    <!-- Navbar -->
    <nav class="navbar">
      <div class="navbar-left">
        <div class="navbar-logo">
          <img src="/src/assets/logo.png" alt="Social Network Logo" class="logo-image" />
          <span class="navbar-title">Social Network</span>
        </div>
      </div>
      
      <div class="navbar-center">
        <div class="search-bar">
          <MagnifyingGlassIcon class="search-icon" />
          <input
            type="text"
            placeholder="Rechercher..."
            v-model="searchQuery"
            @keyup.enter="searchProfiles"
            @input="onSearchInput"
          />
          <div v-if="showResults && searchResults.length > 0" class="search-results-dropdown">
            <div
              v-for="user in searchResults"
              :key="user.id"
              class="search-result-item"
              @click="goToProfile(user.id)"
            >
              <img :src="user.avatar ? ('/' + user.avatar) : '/src/assets/default-avatar.png'" class="search-result-avatar" />
              <span class="search-result-name">{{ user.nickname || (user.firstName + ' ' + user.lastName) }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="navbar-right">
        <div class="notification-wrapper">
          <button class="nav-icon-btn" @click="navigateTo('/notifications')">
            <BellIcon class="icon" :class="{ 'has-new': notificationStore.hasNewNotifications }" />
            <span v-if="notificationStore.unreadCount > 0" class="notification-badge">{{ notificationStore.unreadCount }}</span>
          </button>
          
          <!-- Notification Dropdown -->
          <div v-if="showNotificationDropdown" class="notification-dropdown">
            <div class="notification-header">
              <h3>Notifications</h3>
              <button 
                @click="notificationStore.markAllAsRead()" 
                v-if="notificationStore.unreadCount > 0" 
                class="mark-all-read"
              >
                Marquer tout lu ({{ notificationStore.unreadCount }})
              </button>
            </div>
            <div v-if="notificationStore.loading" class="notification-loading">
              <div class="loading-spinner-sm"></div>
              <span>Chargement...</span>
            </div>
            <div v-else-if="notificationStore.recentUnreadNotifications.length === 0" class="notification-empty">
              Aucune nouvelle notification
            </div>
            <div v-else class="notification-list">
              <!-- Use the new recentUnreadNotifications computed property -->
              <template v-for="notif in notificationStore.recentUnreadNotifications" :key="notif.id">
                <!-- Follow Requests -->
                <div v-if="notif.type === 'FOLLOW'" class="notification-item follow-request">
                  <div class="notification-avatar">
                    <img v-if="notif.user?.avatar" :src="notif.user.avatar" alt="Avatar" />
                    <div v-else class="default-avatar">{{ notif.user?.firstName?.[0] || 'U' }}</div>
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> souhaite vous suivre
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                    <div class="notification-actions">
                      <button 
                        @click="notificationStore.respondToFollowRequest(notif.id, 'ACCEPT')" 
                        class="accept-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-xs"></span>
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Accepter' }}
                      </button>
                      <button 
                        @click="notificationStore.respondToFollowRequest(notif.id, 'DECLINE')" 
                        class="decline-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Refuser' }}
                      </button>
                    </div>
                  </div>
                </div>
                
                <!-- Group Invites -->
                <div v-else-if="notif.type === 'GROUP_INVITE'" class="notification-item group-invite">
                  <div class="notification-avatar">
                    <UserGroupIcon class="icon-sm" />
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> vous invite à rejoindre le groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                    <div class="notification-actions">
                      <button 
                        @click="notificationStore.respondToGroupInvite(notif.id, 'ACCEPT')" 
                        class="accept-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-xs"></span>
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Rejoindre' }}
                      </button>
                      <button 
                        @click="notificationStore.respondToGroupInvite(notif.id, 'DECLINE')" 
                        class="decline-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Refuser' }}
                      </button>
                    </div>
                  </div>
                </div>
                
                <!-- Group Requests -->
                <div v-else-if="notif.type === 'GROUP_REQUEST'" class="notification-item group-request">
                  <div class="notification-avatar">
                    <UserGroupIcon class="icon-sm" />
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> demande à rejoindre votre groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                    <div class="notification-actions">
                      <button 
                        @click="notificationStore.respondToGroupRequest(notif.id, 'ACCEPT', notif.group?.id)" 
                        class="accept-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-xs"></span>
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Accepter' }}
                      </button>
                      <button 
                        @click="notificationStore.respondToGroupRequest(notif.id, 'DECLINE', notif.group?.id)" 
                        class="decline-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Refuser' }}
                      </button>
                    </div>
                  </div>
                </div>
                
                <!-- Event Notifications (automatically marked as read when viewed) -->
                <div v-else-if="notif.type === 'EVENT'" 
                     class="notification-item event-notification"
                     @click="notificationStore.markAsReadOnView(notif.id)">
                  <div class="notification-avatar">
                    <CalendarIcon class="icon-sm" />
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      Nouvel événement <strong>{{ notif.event?.title || 'Événement' }}</strong> dans le groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                  </div>
                  <!-- Visual indicator for unread -->
                  <div v-if="!notif.read" class="unread-indicator"></div>
                </div>
              </template>
            </div>
            <div class="notification-footer">
              <button @click="router.push('/notifications'); showNotificationDropdown = false">Voir toutes les notifications</button>
            </div>
          </div>
        </div>
        
        <button class="nav-icon-btn" @click="navigateTo('/chat')">
          <ChatBubbleLeftRightIcon class="icon" />
        </button>
        
        <div class="navbar-profile">
          <div class="profile-avatar">
            <span>{{ currentUser?.firstName?.[0] || 'U' }}</span>
          </div>
          <ChevronDownIcon class="dropdown-icon" />
        </div>
        
        <button class="nav-icon-btn logout-btn" @click="logout">
          <ArrowRightOnRectangleIcon class="icon" />
        </button>
      </div>
    </nav>

    <!-- Debug Panel for Notifications (Development only) -->
    <div v-if="true" class="debug-panel">
      <details>
        <summary>🐛 Notification Debug Info</summary>
        <div class="debug-content">
          <div class="debug-row">
            <strong>User Logged In:</strong> {{ debugInfo.userLoggedIn }}
          </div>
          <div class="debug-row">
            <strong>WebSocket:</strong> 
            <span :class="debugInfo.wsConnected ? 'status-connected' : 'status-disconnected'">
              {{ debugInfo.wsConnected ? 'Connected' : 'Disconnected' }}
            </span>
          </div>
          <div class="debug-row">
            <strong>Total Notifications:</strong> {{ debugInfo.totalNotifications }}
          </div>
          <div class="debug-row">
            <strong>Unread Count:</strong> {{ debugInfo.unreadCount }}
          </div>
          <div class="debug-row">
            <strong>Has New:</strong> {{ debugInfo.hasNewNotifications }}
          </div>
          <div class="debug-row">
            <strong>Has Any:</strong> {{ debugInfo.hasAnyNotifications }}
          </div>
          <div class="debug-row">
            <strong>Follow Requests:</strong> {{ debugInfo.followRequests }}
          </div>
          <div class="debug-row">
            <strong>Group Invites:</strong> {{ debugInfo.groupInvites }}
          </div>
          <div class="debug-row">
            <strong>Group Requests:</strong> {{ debugInfo.groupRequests }}
          </div>
          <div class="debug-row">
            <strong>Event Notifications:</strong> {{ debugInfo.eventNotifications }}
          </div>
          <button @click="testNotification" class="test-btn">🧪 Add Test Notification</button>
          <button @click="testWebSocketConnection" class="test-btn">🔌 Test WebSocket Connection</button>
          <button @click="forceWebSocketReconnect" class="test-btn">🔄 Force Reconnect</button>
          <button @click="simulateWebSocketMessage" class="test-btn">📡 Simulate WebSocket Message</button>
        </div>
      </details>
    </div>

    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-section">
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/' }"
          @click="navigateTo('/')"
        >
          <HomeIcon />
          <span>Accueil</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/profile' }"
          @click="navigateTo('/profile')"
        >
          <UserIcon />
          <span>Profil</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/groups' }"
          @click="navigateTo('/groups')"
        >
          <UserGroupIcon />
          <span>Groupes</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/chat' }"
          @click="navigateTo('/chat')"
        >
          <ChatBubbleLeftRightIcon />
          <span>Messages</span>
        </button>
      </div>
      
      <div class="sidebar-section">
        <div class="sidebar-section-title">Raccourcis</div>
        <button class="sidebar-item">
          <HeartIcon />
          <span>Favoris</span>
        </button>
        
        <button class="sidebar-item">
          <DocumentDuplicateIcon />
          <span>Pages</span>
        </button>
      </div>

      <div class="sidebar-footer">
        <button class="sidebar-item logout-btn" @click="logout">
          <ArrowRightOnRectangleIcon />
          <span>Déconnexion</span>
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="main-content">
      <div class="content-wrapper">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script>
import {
  MagnifyingGlassIcon,
  BellIcon,
  ChatBubbleLeftRightIcon,
  ChevronDownIcon,
  ArrowRightOnRectangleIcon,
  HomeIcon,
  UserIcon,
  UserGroupIcon,
  HeartIcon,
  DocumentDuplicateIcon,
  CalendarIcon
} from '@heroicons/vue/24/outline'
import { useUserStore } from '../stores/userStore'
import { useNotificationStore } from '../stores/notificationStore'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import wsService from '../services/websocketService'

export default {
  name: 'MainLayout',
  components: {
    MagnifyingGlassIcon,
    BellIcon,
    ChatBubbleLeftRightIcon,
    ChevronDownIcon,
    ArrowRightOnRectangleIcon,
    HomeIcon,
    UserIcon,
    UserGroupIcon,
    HeartIcon,
    DocumentDuplicateIcon,
    CalendarIcon
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const userStore = useUserStore()
    const notificationStore = useNotificationStore()
    const showResults = ref(false)
    const searchQuery = ref('')
    const searchResults = ref([])
    const showUserMenu = ref(false)
    const showNotificationDropdown = ref(false)
    
    let searchTimeout = null
    let removeWebSocketListener = null

    // Computed
    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const currentUser = computed(() => userStore.user)
    const currentRoute = computed(() => route.path)
    
    // Check if there are any notifications to display in dropdown
    const hasAnyNotifications = computed(() => {
      return notificationStore.followRequests.length > 0 ||
             notificationStore.groupInvites.length > 0 ||
             notificationStore.groupRequests.length > 0 ||
             notificationStore.eventNotifications.length > 0
    })

    // Debug info for notification system
    const debugInfo = computed(() => ({
      totalNotifications: notificationStore.notifications.length,
      unreadCount: notificationStore.unreadCount,
      hasNewNotifications: notificationStore.hasNewNotifications,
      hasAnyNotifications: hasAnyNotifications.value,
      followRequests: notificationStore.followRequests.length,
      groupInvites: notificationStore.groupInvites.length,
      groupRequests: notificationStore.groupRequests.length,
      eventNotifications: notificationStore.eventNotifications.length,
      wsConnected: wsService.getConnectionStatus().value,
      userLoggedIn: userStore.isLoggedIn
    }))

    // Methods
    async function navigateTo(path) {
      if (path === '/notifications') {
        showNotificationDropdown.value = !showNotificationDropdown.value
        if (showNotificationDropdown.value) {
          // Fetch fresh notifications when opening dropdown
          console.log('📂 Opening notification dropdown, fetching fresh notifications...')
          await notificationStore.fetchNotifications()
          
          // Auto-mark event notifications as read when they're viewed in dropdown
          // This simulates the behavior of modern social networks where viewing = read
          setTimeout(async () => {
            await notificationStore.markDropdownNotificationsAsRead()
          }, 1000) // Small delay to allow user to see the notifications first
        }
      } else {
        router.push(path)
      }
    }

    async function logout() {
      try {
        const response = await fetch('http://localhost:8081/logout', {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        // Always redirect to login page after logout attempt
        // Even if the logout request fails, we want to redirect the user
        router.push('/login');
        
        if (!response.ok) {
          console.warn('Logout request failed, but redirecting to login anyway');
        }
      } catch (error) {
        console.error('Error during logout:', error);
        // Still redirect to login page even if there's a network error
        router.push('/login');
      }
    }

    async function searchProfiles() {
      if (!searchQuery.value.trim()) {
        searchResults.value = [];
        showResults.value = false;
        return;
      }
      try {
        const response = await fetch('http://localhost:8081/allUsers', {
          credentials: 'include'
        });
        if (response.ok) {
          const data = await response.json();
          const users = data.users || [];
          const query = searchQuery.value.trim().toLowerCase();
          searchResults.value = users.filter(user => {
            return (
              (user.nickname && user.nickname.toLowerCase().includes(query)) ||
              (user.firstName && user.firstName.toLowerCase().includes(query)) ||
              (user.lastName && user.lastName.toLowerCase().includes(query)) ||
              (user.email && user.email.toLowerCase().includes(query))
            );
          });
          showResults.value = true;
        } else {
          searchResults.value = [];
          showResults.value = false;
        }
      } catch (e) {
        searchResults.value = [];
        showResults.value = false;
      }
    }

    function onSearchInput() {
      if (searchTimeout) clearTimeout(searchTimeout);
      if (!searchQuery.value.trim()) {
        searchResults.value = [];
        showResults.value = false;
        return;
      }
      // debounce: attend 300ms après la dernière frappe
      searchTimeout = setTimeout(() => {
        searchProfiles();
      }, 300);
    }

    function goToProfile(userId) {
      router.push(`/profile/${userId}`);
      showResults.value = false;
      searchQuery.value = '';
    }

    async function testNotification() {
      // Create a test notification
      const testNotif = {
        id: Date.now(),
        type: 'FOLLOW',
        message: 'Test notification',
        user: {
          id: 999,
          firstName: 'Test',
          lastName: 'User',
          nickname: 'testuser'
        },
        read: false,
        createdAt: new Date().toISOString()
      }
      
      console.log('🧪 Creating test notification:', testNotif)
      await notificationStore.handleNewNotification(testNotif)
    }

    function testWebSocketConnection() {
      console.log('🔌 Testing WebSocket connection...')
      const wsStatus = wsService.isConnected()
      console.log('WebSocket status:', {
        isConnected: wsStatus,
        readyState: wsStatus ? 'OPEN' : 'CLOSED'
      })
      
      if (wsStatus) {
        console.log('✅ WebSocket is connected and ready')
        // Don't send test messages, just confirm connection
      } else {
        console.log('❌ WebSocket is not connected')
      }
    }

    function forceWebSocketReconnect() {
      console.log('🔄 Forcing WebSocket reconnection...')
      wsService.disconnect()
      setTimeout(() => {
        wsService.connect()
      }, 1000)
    }

    function simulateWebSocketMessage() {
      console.log('📡 Simulating WebSocket notification message...')
      
      // Create a mock WebSocket message that matches backend format
      const mockMessage = {
        action: 'notification',
        notification: {
          id: Date.now(),
          type: 'FOLLOW',
          sender: '123',
          target: userStore.user?.id,
          content: '123',
          createdAt: new Date().toISOString(),
          read: false,
          user: {
            id: '123',
            firstName: 'Test',
            lastName: 'User',
            nickname: 'testuser'
          }
        }
      }
      
      // Simulate receiving the message through WebSocket service
      if (removeWebSocketListener) {
        console.log('🎯 Triggering WebSocket message handler...')
        // Find the mainLayout listener and call it directly
        const listeners = wsService.getListeners()
        listeners.forEach((callback, key) => {
          if (key === 'mainLayout') {
            callback(mockMessage)
          }
        })
      }
    }

    // Setup WebSocket for real-time notifications
    // WebSocket handler
    onMounted(async () => {
      console.log('🎯 MainLayout mounted, checking user auth status...')
      
      // First, check if user is authenticated by calling the backend
      try {
        console.log('🔍 Checking user session...')
        await userStore.reloadUserAfterRefresh()
        
        console.log('✅ User auth check complete:', {
          isLoggedIn: userStore.isLoggedIn,
          isAuthenticated: userStore.isAuthenticated,
          user: userStore.user
        })
        
        if (userStore.isLoggedIn) {
          console.log('✅ User is logged in, setting up notifications and WebSocket')
          await notificationStore.fetchNotifications()
          
          // Give a small delay to ensure the session is established
          setTimeout(() => {
            console.log('🔌 Setting up WebSocket listener...')
            removeWebSocketListener = wsService.addListener('mainLayout', async (data) => {
              console.log('🔔 MainLayout received WebSocket message:', data)
              // Handle notification messages
              if (data.action === 'notification') {
                console.log('📬 New notification received, updating store...')
                await notificationStore.handleNewNotification(data.notification)
                
                // Force reactivity update by accessing the computed properties
                // This ensures the template re-renders
                console.log('📊 Updated notification counts:', {
                  total: notificationStore.notifications.length,
                  unread: notificationStore.unreadCount,
                  hasNew: notificationStore.hasNewNotifications,
                  followRequests: notificationStore.followRequests.length,
                  groupInvites: notificationStore.groupInvites.length,
                  groupRequests: notificationStore.groupRequests.length,
                  eventNotifications: notificationStore.eventNotifications.length
                })
                
                // Force template update by accessing hasAnyNotifications
                console.log('📋 Has any notifications:', hasAnyNotifications.value)
              }
            })
          }, 1000) // 1 second delay to ensure session is ready
        } else {
          console.log('❌ User not logged in, skipping WebSocket setup')
        }
      } catch (error) {
        console.error('❌ Error checking user authentication:', error)
      }
    })

    onUnmounted(() => {
      // Clean up WebSocket listener
      if (removeWebSocketListener) {
        removeWebSocketListener()
      }
      if (searchTimeout) {
        clearTimeout(searchTimeout)
      }
    })

    return {
      router,
      currentRoute,
      searchQuery,
      searchResults,
      showResults,
      showUserMenu,
      showNotificationDropdown,
      userStore,
      notificationStore,
      hasAnyNotifications,
      debugInfo,
      navigateTo,
      logout,
      searchProfiles,
      onSearchInput,
      goToProfile,
      testNotification,
      testWebSocketConnection,
      forceWebSocketReconnect,
      simulateWebSocketMessage,
      isLoggedIn,
      currentUser
    }
  }
}
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  position: relative;
  z-index: 1;
}

/* Navbar Styles */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: rgba(15, 15, 23, 0.9);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 1000;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
}

.navbar-left {
  display: flex;
  align-items: center;
}

.navbar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.navbar-title {
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.navbar-center {
  flex: 1;
  max-width: 400px;
  margin: 0 40px;
}

.search-bar {
  position: relative;
  width: 100%;
}

.search-bar input {
  width: 100%;
  padding: 8px 15px 8px 45px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 25px;
  color: #ffffff;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.search-bar input:focus {
  outline: none;
  border-color: #e879c6;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.2);
}

.search-bar input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.search-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.4);
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.nav-icon-btn {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 8px;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-icon-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
  transform: translateY(-1px);
}

/* Sidebar Icon Styling */
.sidebar-item svg,
.sidebar-item > svg {
  width: 20px;
  height: 20px;
  color: rgba(255, 255, 255, 0.6);
  transition: color 0.3s ease;
}

.sidebar-item:hover svg,
.sidebar-item.active svg {
  color: #e879c6;
}

/* Notification Styling */
.notification-wrapper {
  position: relative;
}

.notification-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #e879c6;
  color: white;
  font-size: 10px;
  font-weight: bold;
  min-width: 18px;
  height: 18px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: notification-pulse 1.5s ease-in-out infinite;
}

@keyframes notification-pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

.icon.has-new {
  animation: icon-glow 1.5s ease-in-out infinite;
  color: #e879c6;
}

@keyframes icon-glow {
  0% { filter: drop-shadow(0 0 0px #e879c6); }
  50% { filter: drop-shadow(0 0 2px #e879c6); }
  100% { filter: drop-shadow(0 0 0px #e879c6); }
}

/* Different notification types styling */
.notification-item {
  border-bottom: 1px solid rgba(120, 119, 198, 0.05);
  transition: background 0.2s;
}

.notification-item:hover {
  background: #23233a;
}

.notification-item.follow-request {
  border-left: 3px solid #e879c6;
}

.notification-item.group-invite {
  border-left: 3px solid #78dbff;
}

.notification-item.group-request {
  border-left: 3px solid #7879c6;
}

.notification-item.event-notification {
  border-left: 3px solid #84cc16;
  position: relative;
  cursor: pointer;
}

.unread-indicator {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 8px;
  height: 8px;
  background: #e879c6;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.notification-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: #23233a;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.notification-avatar .default-avatar {
  color: white;
  font-size: 16px;
}

.icon-sm {
  width: 24px;
  height: 24px;
  color: white;
}

.icon-xs {
  width: 16px;
  height: 16px;
  display: inline-block;
  vertical-align: middle;
  margin-right: 4px;
}

.notification-content {
  flex: 1;
  margin-left: 10px;
}

.notification-text {
  font-size: 14px;
  color: #d1d5db;
  margin-bottom: 8px;
}

.notification-text strong {
  color: white;
}

.notification-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.notification-actions button {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.notification-actions button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.accept-btn {
  background: linear-gradient(90deg, #78dbff 0%, #7879c6 100%);
  color: white;
}

.accept-btn:hover {
  opacity: 0.9;
}

.decline-btn {
  background: #23233a;
  color: #d1d5db;
  border: 1px solid rgba(120, 119, 198, 0.3) !important;
}

.decline-btn:hover {
  background: #282846;
}

.notification-footer {
  padding: 10px;
  text-align: center;
  border-top: 1px solid rgba(120, 119, 198, 0.1);
}

.notification-footer button {
  background: none;
  border: none;
  color: #78dbff;
  font-size: 14px;
  cursor: pointer;
  padding: 5px 10px;
  transition: color 0.2s;
}

.notification-footer button:hover {
  color: #e879c6;
  text-decoration: underline;
}

/* Sidebar Styles */
.sidebar {
  position: fixed;
  top: 60px;
  left: 0;
  width: 250px;
  height: calc(100vh - 60px);
  background: rgba(15, 15, 23, 0.9);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 0;
  overflow-y: auto;
  z-index: 100;
  display: flex;
  flex-direction: column;
}

.sidebar-section {
  margin-bottom: 25px;
}

.sidebar-section-title {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 20px 10px 20px;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 20px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
}

.sidebar-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
  transform: translateX(5px);
}

.sidebar-item.active {
  background: linear-gradient(90deg, rgba(232, 121, 198, 0.2) 0%, transparent 100%);
  color: #e879c6;
  border-right: 3px solid #e879c6;
}

.sidebar-item.active svg {
  color: #e879c6;
}

.sidebar-footer {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1) !important;
  color: #ef4444 !important;
}

/* Main Content Styles */
.main-content {
  margin-left: 250px;
  margin-top: 60px;
  min-height: calc(100vh - 60px);
  position: relative;
  z-index: 1;
}

.content-wrapper {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

/* Icon Styling */
.search-icon,
.icon,
.dropdown-icon {
  width: 20px;
  height: 20px;
  color: rgba(255, 255, 255, 0.7);
  transition: color 0.3s ease;
}

.dropdown-icon {
  width: 16px;
  height: 16px;
}

/* Navbar Icons */
.nav-icon-btn:hover .icon {
  color: #e879c6;
}

/* Logo Styling */
.logo-image {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  animation: logoGlow 3s ease-in-out infinite;
  transition: all 0.3s ease;
  filter: drop-shadow(0 0 10px rgba(232, 121, 198, 0.3));
}

.logo-image:hover {
  transform: scale(1.1);
  filter: drop-shadow(0 0 20px rgba(120, 199, 255, 0.6)) 
          drop-shadow(0 0 30px rgba(232, 121, 198, 0.4));
}

@keyframes logoGlow {
  0%, 100% { 
    filter: drop-shadow(0 0 15px rgba(232, 121, 198, 0.4)) 
            drop-shadow(0 0 25px rgba(168, 85, 247, 0.2));
  }
  50% { 
    filter: drop-shadow(0 0 25px rgba(120, 199, 255, 0.5)) 
            drop-shadow(0 0 35px rgba(232, 121, 198, 0.3));
  }
}

/* Navbar Profile */
.navbar-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 20px;
  transition: background-color 0.3s ease;
}

.navbar-profile:hover {
  background: rgba(255, 255, 255, 0.05);
}

.navbar-profile .profile-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 12px;
}

/* Search Results Dropdown */
.search-results-dropdown {
  position: absolute;
  top: 110%;
  left: 0;
  width: 100%;
  background: rgba(15, 15, 23, 0.98);
  border-radius: 10px;
  box-shadow: 0 4px 24px rgba(0,0,0,0.25);
  z-index: 2000;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid rgba(255,255,255,0.08);
}

.search-result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.search-result-item:hover {
  background: rgba(232, 121, 198, 0.08);
}

.search-result-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  background: #eee;
}

.search-result-name {
  color: #fff;
  font-size: 1rem;
  font-weight: 500;
}

/* Add missing notification dropdown styles */
.notification-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 400px;
  max-height: 500px;
  background: rgba(15, 15, 23, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  z-index: 2000;
  overflow: hidden;
  backdrop-filter: blur(20px);
}

.notification-header {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.02);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notification-header h3 {
  margin: 0;
  color: white;
  font-size: 16px;
  font-weight: 600;
}

.mark-all-read {
  background: transparent;
  border: 1px solid rgba(120, 219, 255, 0.3);
  color: #78dbff;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.mark-all-read:hover {
  background: rgba(120, 219, 255, 0.1);
  border-color: rgba(120, 219, 255, 0.5);
}

.notification-loading,
.notification-empty {
  padding: 20px;
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.notification-time {
  font-size: 11px;
  color: #78dbff;
  margin-top: 4px;
  margin-bottom: 8px;
}

.loading-spinner-sm {
  border: 2px solid rgba(120, 219, 255, 0.1);
  border-radius: 50%;
  border-top: 2px solid #78dbff;
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

.loading-spinner-xs {
  border: 1px solid rgba(120, 219, 255, 0.1);
  border-radius: 50%;
  border-top: 1px solid #78dbff;
  width: 12px;
  height: 12px;
  animation: spin 1s linear infinite;
  display: inline-block;
  margin-right: 4px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.notification-list {
  max-height: 350px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: background 0.2s;
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.02);
}

/* Debug Panel Styles */
.debug-panel {
  position: fixed;
  top: 60px;
  right: 10px;
  background: rgba(15, 15, 23, 0.95);
  border: 1px solid rgba(120, 219, 255, 0.3);
  border-radius: 8px;
  padding: 0.5rem;
  z-index: 9999;
  font-size: 12px;
  max-width: 300px;
}

.debug-panel summary {
  cursor: pointer;
  color: #78dbff;
  font-weight: bold;
  padding: 0.25rem;
}

.debug-content {
  margin-top: 0.5rem;
  padding: 0.5rem 0;
}

.debug-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.25rem;
  color: #d1d5db;
}

.debug-row strong {
  color: #78dbff;
}

.status-connected {
  color: #10b981;
}

.status-disconnected {
  color: #ef4444;
}

.test-btn {
  background: linear-gradient(90deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
  font-size: 11px;
  cursor: pointer;
  margin-top: 0.5rem;
  width: 100%;
}

.test-btn:hover {
  opacity: 0.9;
}
</style>
