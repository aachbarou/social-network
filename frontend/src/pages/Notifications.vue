<template>
  <div class="notifications-page">
    <div class="notifications-header">
      <h1>Notifications</h1>
      <div class="header-actions">
        <button
          v-if="notificationStore.unreadCount > 0"
          @click="notificationStore.markAllAsRead()"
          class="mark-all-read-btn"
        >
          Tout marquer comme lu ({{ notificationStore.unreadCount }})
        </button>
      </div>
    </div>

    <!-- Error -->
    <div v-if="notificationStore.error" class="error-banner">
      <div class="error-content">
        <span>{{ notificationStore.error }}</span>
        <button @click="notificationStore.dismissError()" class="dismiss-btn">×</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="notificationStore.loading" class="loading">
      <div class="loading-spinner"></div>
      <span>Chargement des notifications...</span>
    </div>

    <!-- Empty State -->
    <div v-else-if="notificationStore.notifications.length === 0" class="empty-state">
      Vous n'avez aucune notification.
    </div>

    <!-- Notifications List -->
    <div v-else>
      <!-- Follow Requests -->
      <div class="notification-section" v-if="notificationStore.followRequests.length > 0">
        <h2>Demandes de suivi</h2>
        <div class="notification-list">
          <div
            v-for="notif in notificationStore.followRequests"
            :key="notif.id"
            class="notification-card follow-request"
          >
            <div class="notification-header">
              <div class="notification-avatar">
                <img v-if="notif.user?.avatar" :src="notif.user.avatar" alt="Avatar" />
                <div v-else class="default-avatar">
                  {{ notif.user?.firstName?.[0] || 'U' }}
                </div>
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong>
                  souhaite vous suivre
                </div>
                <div class="notification-time">
                  {{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}
                </div>
              </div>
            </div>
            <div class="notification-actions">
              <button
                @click="notificationStore.respondToFollowRequest(notif.id, 'ACCEPT')"
                class="accept-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-sm"></span>
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Accepter' }}
              </button>
              <button
                @click="notificationStore.respondToFollowRequest(notif.id, 'DECLINE')"
                class="decline-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Refuser' }}
              </button>
              <button
                @click="dismissWithConfirmation(notif.id, 'follow')"
                class="dismiss-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                Supprimer
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Group Invites -->
      <div class="notification-section" v-if="notificationStore.groupInvites.length > 0">
        <h2>Invitations de groupe</h2>
        <p class="section-info">Utilisez "Rejoindre" ou "Refuser" pour répondre aux invitations.</p>
        <div class="notification-list">
          <div
            v-for="notif in notificationStore.groupInvites"
            :key="notif.id"
            class="notification-card group-invite"
          >
            <div class="notification-header">
              <div class="notification-avatar group">
                <UserGroupIcon class="icon-sm" />
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong>
                  vous invite à rejoindre le groupe
                  <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                </div>
                <div class="notification-subtitle" v-if="notif.group?.description">
                  {{ notif.group.description }}
                </div>
                <div class="notification-time">
                  {{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}
                </div>
              </div>
            </div>
            <div class="notification-actions">
              <button
                @click="notificationStore.respondToGroupInvite(notif.id, 'ACCEPT')"
                class="accept-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-sm"></span>
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Rejoindre' }}
              </button>
              <button
                @click="notificationStore.respondToGroupInvite(notif.id, 'DECLINE')"
                class="decline-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Refuser' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Group Requests -->
      <div class="notification-section" v-if="notificationStore.groupRequests.length > 0">
        <h2>Demandes d'adhésion aux groupes</h2>
        <p class="section-info">Utilisez "Accepter" ou "Refuser" pour répondre aux demandes.</p>
        <div class="notification-list">
          <div
            v-for="notif in notificationStore.groupRequests"
            :key="notif.id"
            class="notification-card group-request"
          >
            <div class="notification-header">
              <div class="notification-avatar group">
                <UserGroupIcon class="icon-sm" />
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong>
                  demande à rejoindre votre groupe
                  <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                </div>
                <div class="notification-time">
                  {{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}
                </div>
              </div>
            </div>
            <div class="notification-actions">
              <button
                @click="notificationStore.respondToGroupRequest(notif.id, 'ACCEPT', notif.group?.id)"
                class="accept-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-sm"></span>
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Accepter' }}
              </button>
              <button
                @click="notificationStore.respondToGroupRequest(notif.id, 'DECLINE', notif.group?.id)"
                class="decline-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Refuser' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Event Notifications -->
      <div class="notification-section" v-if="notificationStore.eventNotifications.length > 0">
        <h2>Événements de groupe</h2>
        <div class="notification-list">
          <div
            v-for="notif in notificationStore.eventNotifications"
            :key="notif.id"
            class="notification-card event-notification"
            :class="{ 'unread': !notif.read }"
            @click="notificationStore.markAsReadOnView(notif.id)"
          >
            <div class="notification-header">
              <div class="notification-avatar event">
                <CalendarIcon class="icon-sm" />
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  Nouvel événement <strong>{{ notif.event?.title || 'Événement' }}</strong> dans le groupe
                  <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                  <span v-if="!notif.read" class="unread-dot"></span>
                </div>
                <div class="notification-subtitle" v-if="notif.event?.description">
                  {{ notif.event.description }}
                </div>
                <div class="notification-date" v-if="notif.event?.date">
                  <CalendarDaysIcon class="icon-xs" />
                  {{ new Date(notif.event.date).toLocaleDateString() }}
                </div>
                <div class="notification-time">
                  {{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}
                </div>
              </div>
            </div>
            <div class="notification-actions">
              <button
                @click.stop="dismissWithConfirmation(notif.id, 'event')"
                class="dismiss-btn"
              >
                Supprimer
              </button>
              <button
                @click.stop="notificationStore.markAsRead(notif.id)"
                class="mark-read-btn"
                v-if="!notif.read"
              >
                Marquer comme lu
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notificationStore'
import { UserGroupIcon, CalendarIcon, CalendarDaysIcon } from '@heroicons/vue/24/outline'
import wsService from '../services/websocketService'

const router = useRouter()
const notificationStore = useNotificationStore()
let removeWebSocketListener = null

function dismissWithConfirmation(notificationId, type) {
  let message = ''
  switch (type) {
    case 'follow':
      message = 'Supprimer cette demande de suivi ? La personne pourra soumettre une nouvelle demande.'
      break
    case 'event':
      message = 'Supprimer cette notification d\'événement ?'
      break
    default:
      message = 'Supprimer cette notification ?'
  }
  if (confirm(message)) {
    notificationStore.dismissNotification(notificationId)
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
  removeWebSocketListener = wsService.addListener('notificationsPage', (data) => {
    if (data.action === 'notification') {
      notificationStore.fetchNotifications()
    }
  })
})

onUnmounted(() => {
  if (removeWebSocketListener) removeWebSocketListener()
})
</script>
<style scoped>
/* --- Global Card Style Reuse --- */
.notification-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.notification-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
}

/* --- Enhanced Default Avatar --- */
.default-avatar {
  background: linear-gradient(135deg, #3b3c5c, #4b4e7c);
  color: white;
  font-size: 20px;
  font-weight: bold;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* --- Enhanced Section Titles --- */
.notification-section h2 {
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

/* --- Empty State Styling --- */
.empty-state {
  color: #aaa;
  background: #1e1e2e;
  border: 1px dashed #444;
  border-radius: 12px;
  padding: 2rem;
  margin: 2rem 0;
  font-style: italic;
  font-size: 15px;
}

/* --- Error Animation (Optional) --- */
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-4px); }
  40%, 80% { transform: translateX(4px); }
}

.error-banner {
  animation: shake 0.3s ease-in-out;
}

/* --- New Responsive Utility --- */
@media (max-width: 600px) {
  .notifications-page {
    padding: 1rem;
  }

  .notification-card {
    padding: 1rem;
  }

  .notification-title {
    font-size: 14px;
  }

  .notification-subtitle,
  .notification-date,
  .notification-time {
    font-size: 12px;
  }

  .notification-actions {
    flex-direction: column;
    align-items: flex-end;
  }

  .notification-actions button {
    width: 100%;
  }
}

</style>