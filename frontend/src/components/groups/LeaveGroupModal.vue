<template>
  <div class="modal-overlay" @click="$emit('cancel')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>Quitter le groupe</h3>
      </div>
      <div class="modal-body">
        <p>Êtes-vous sûr de vouloir quitter le groupe <strong>{{ group?.name }}</strong> ?</p>
        <p class="warning-text">Cette action ne peut pas être annulée.</p>
      </div>
      <div class="modal-actions">
        <button @click="$emit('cancel')" class="btn-secondary">
          Annuler
        </button>
        <button @click="handleConfirm" class="btn-danger" :disabled="isLeaving">
          <span 
            class="icon" 
            :class="{ 'animate-spin': isLeaving }"
          >
            {{ isLeaving ? '⏳' : '🚪' }}
          </span>
          {{ isLeaving ? 'Quitter...' : 'Quitter' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  group: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const isLeaving = ref(false)

const handleConfirm = async () => {
  if (isLeaving.value) return
  
  isLeaving.value = true
  
  try {
    // Add minimum delay for UX
    const minDelay = new Promise(resolve => setTimeout(resolve, 1000))
    
    await Promise.all([
      emit('confirm'),
      minDelay
    ])
  } finally {
    isLeaving.value = false
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: rgba(15, 15, 23, 0.95);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0;
  max-width: 400px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.modal-header {
  padding: 24px 24px 0 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 20px;
}

.modal-header h3 {
  color: white;
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.modal-body {
  padding: 0 24px 20px 24px;
}

.modal-body p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 16px 0;
  line-height: 1.5;
}

.modal-body .warning-text {
  color: rgba(255, 107, 107, 0.8);
  font-size: 0.875rem;
  margin: 0;
}

.modal-actions {
  padding: 20px 24px 24px 24px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.btn-secondary,
.btn-danger {
  padding: 10px 20px;
  font-size: 0.875rem;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
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

.btn-secondary:hover,
.btn-danger:hover {
  transform: translateY(-1px);
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
