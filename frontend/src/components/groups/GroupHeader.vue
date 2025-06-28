<template>
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
    
    <GroupActions 
      :group="group"
      :is-member="isMember"
      :is-admin="isAdmin"
      @join="$emit('join')"
      @leave="$emit('leave')"
      @chat="$emit('chat')"
    />
  </div>
</template>

<script setup>
import GroupActions from './GroupActions.vue'

defineProps({
  group: {
    type: Object,
    required: true
  },
  isMember: {
    type: Boolean,
    default: false
  },
  isAdmin: {
    type: Boolean,
    default: false
  }
})

defineEmits(['join', 'leave', 'chat'])
</script>

<style scoped>
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

.icon {
  font-size: 0.875rem;
}

/* Responsive */
@media (max-width: 768px) {
  .group-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .group-info {
    flex-direction: column;
    text-align: center;
  }
}
</style>
