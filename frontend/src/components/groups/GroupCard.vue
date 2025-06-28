<template>
  <div class="group-card">
    <!-- Group Image -->
    <div class="group-image">
      <img 
        v-if="group.image" 
        :src="group.image" 
        :alt="group.name"
        @error="handleImageError"
      />
      <div v-else class="default-image">
        <span class="group-icon">👥</span>
      </div>
    </div>

    <!-- Group Info -->
    <div class="group-info">
      <h3 class="group-name">{{ group.name }}</h3>
      <p class="group-description">{{ group.description }}</p>
      
      <div class="group-meta">
        <span class="member-count">
          <span class="icon">👤</span>
          {{ group.memberCount || 0 }} membres
        </span>
        <span class="group-privacy" :class="group.privacy">
          <span class="icon">{{ privacyIcon }}</span>
          {{ privacyLabel }}
        </span>
      </div>
    </div>

    <!-- Actions -->
    <div class="group-actions">
      <button 
        v-if="group.privacy === 'public' || isMember"
        @click="$emit('view', group)"
        class="btn-secondary"
      >
        <span class="icon">👁️</span>
        Voir
      </button>

      <!-- Member actions -->
      <template v-if="isMember">
        <button 
          @click="$emit('chat', group.id)"
          class="btn-primary"
        >
          <span class="icon">💬</span>
          Chat
        </button>
      </template>

      <!-- Non-member actions -->
      <template v-else>
        <button 
          v-if="!group.requestPending"
          @click="handleJoinRequest"
          class="btn-primary"
          :disabled="props.loading"
        >
          <template v-if="props.loading">
            <span class="icon animate-spin">⏳</span>
            {{ group.privacy === 'public' ? 'Rejoindre...' : 'Envoi...' }}
          </template>
          <template v-else-if="group.privacy === 'public'">
            <span class="icon">➕</span>
            Rejoindre
          </template>
          <template v-else>
            <span class="icon">✉️</span>
            Demander
          </template>
        </button>
        
        <div v-else class="request-pending">
          <span class="icon">⏰</span>
          Demande envoyée
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  group: {
    type: Object,
    required: true
  },
  isMember: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['view', 'chat', 'join-request'])

const privacyIcon = computed(() => {
  return props.group.privacy === 'public' ? '🌐' : '🔒'
})

const privacyLabel = computed(() => {
  return props.group.privacy === 'public' ? 'Public' : 'Privé'
})

const handleImageError = (event) => {
  event.target.style.display = 'none'
  event.target.nextElementSibling.style.display = 'flex'
}

const handleJoinRequest = () => {
  emit('join-request', props.group.id)
}
</script>

<style scoped>
.group-card {
  background: rgba(15, 15, 23, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.group-card:hover {
  transform: translateY(-4px);
  border-color: rgba(232, 121, 198, 0.3);
  box-shadow: 0 10px 40px rgba(232, 121, 198, 0.2);
}

.group-image {
  width: 100%;
  height: 160px;
  position: relative;
  overflow: hidden;
}

.group-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-image {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.group-icon {
  font-size: 2rem;
  opacity: 0.8;
}

.group-info {
  padding: 20px;
  flex: 1;
}

.group-name {
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.group-description {
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 16px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.group-meta {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.member-count,
.group-privacy {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.group-privacy.public {
  color: rgba(120, 199, 255, 0.8);
}

.group-privacy.private {
  color: rgba(232, 121, 198, 0.8);
}

.group-actions {
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  gap: 12px;
}

.btn-primary,
.btn-secondary {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 0.875rem;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.request-pending {
  flex: 1;
  padding: 10px 16px;
  border-radius: 12px;
  background: rgba(255, 193, 7, 0.2);
  color: #ffc107;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  border: 1px solid rgba(255, 193, 7, 0.3);
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Responsive */
@media (max-width: 768px) {
  .group-actions {
    flex-direction: column;
  }
  
  .btn-primary,
  .btn-secondary {
    flex: none;
  }
}
</style>
