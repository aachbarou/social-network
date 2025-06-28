<template>
  <div class="events-section">
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
</template>

<script setup>
defineProps({
  events: {
    type: Array,
    default: () => []
  }
})

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
</script>

<style scoped>
.events-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.event-card {
  background: rgba(15, 15, 23, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
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

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
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

/* Responsive */
@media (max-width: 768px) {
  .event-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>
