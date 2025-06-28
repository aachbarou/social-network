<template>
  <div class="group-actions">
    <!-- Join button for non-members -->
    <button 
      v-if="canJoinGroup"
      @click="handleJoin" 
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
      v-if="isMember || isAdmin"
      @click="$emit('chat')" 
      class="btn-primary"
    >
      <span class="icon">💬</span>
      Chat
    </button>
    
    <!-- Leave button for members (not admins) -->
    <button 
      v-if="canLeaveGroup"
      @click="handleLeave" 
      class="btn-icon btn-danger"
      :disabled="isLeaving"
      title="Quitter le groupe"
    >
      <span 
        class="icon" 
        :class="{ 'animate-spin': isLeaving }"
      >
        {{ isLeaving ? '⏳' : '🚪' }}
      </span>
    </button>
    
    <button @click="$router.go(-1)" class="btn-secondary">
      <span class="icon">←</span>
      Retour
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps({
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

const emit = defineEmits(['join', 'leave', 'chat'])

const isJoining = ref(false)
const isLeaving = ref(false)

const canJoinGroup = computed(() => {
  return props.group && 
         (props.group.privacy === 'public' || props.group.privacy === 'private') &&
         !props.isMember && 
         !props.isAdmin
})

const canLeaveGroup = computed(() => {
  return props.group && 
         (props.isMember || props.isAdmin)
})

const handleJoin = async () => {
  if (isJoining.value) return
  
  isJoining.value = true
  
  try {
    // Show animation for minimum time
    await new Promise(resolve => setTimeout(resolve, 1200))
    
    // After animation, emit the join event
    emit('join')
  } finally {
    isJoining.value = false
  }
}

const handleLeave = async () => {
  if (isLeaving.value) return
  
  isLeaving.value = true
  
  try {
    // Show animation for minimum time
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // After animation, emit the leave event
    emit('leave')
  } finally {
    isLeaving.value = false
  }
}
</script>

<style scoped>
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

.btn-icon {
  width: 48px;
  height: 48px;
  padding: 0;
  border: none;
  border-radius: 50%;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
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

.btn-danger {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
}

.btn-primary:hover,
.btn-secondary:hover,
.btn-icon:hover {
  transform: translateY(-1px);
}

.btn-danger:hover {
  transform: translateY(-1px);
  background: linear-gradient(135deg, #ff5252 0%, #d63031 100%);
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.icon {
  font-size: 0.875rem;
}
</style>
