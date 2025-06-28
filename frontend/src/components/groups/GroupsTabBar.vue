<template>
  <div class="tabs-container">
    <button 
      v-for="tab in tabs" 
      :key="tab.id"
      @click="$emit('tab-change', tab.id)"
      :class="['tab-btn', { active: activeTab === tab.id }]"
    >
      <span class="icon">{{ tab.icon }}</span>
      {{ tab.label }}
      <span v-if="tab.count > 0" class="tab-count">{{ tab.count }}</span>
    </button>
  </div>
</template>

<script setup>
defineProps({
  activeTab: {
    type: String,
    required: true
  },
  tabs: {
    type: Array,
    required: true
  }
})

defineEmits(['tab-change'])
</script>

<style scoped>
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
  position: relative;
}

.tab-btn.active {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.tab-count {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  padding: 2px 6px;
  font-size: 0.75rem;
  font-weight: 700;
  min-width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon {
  font-size: 1rem;
}

/* Responsive */
@media (max-width: 768px) {
  .tabs-container {
    flex-direction: column;
  }
}
</style>
