<template>
  <div class="posts-section">
    <div v-if="posts.length === 0" class="empty-state">
      <span class="empty-icon">📝</span>
      <h3>Aucun post</h3>
      <p>Soyez le premier à publier dans ce groupe</p>
    </div>
    <div v-else class="posts-list">
      <div v-for="post in posts" :key="post.id" class="post-card">
        <div class="post-header">
          <div class="author-avatar">
            <img v-if="post.author.profilePic" :src="post.author.profilePic" :alt="post.author.name" />
            <span v-else class="default-author">👤</span>
          </div>
          <div class="post-info">
            <h4>{{ post.author.firstName }} {{ post.author.lastName }}</h4>
            <span class="post-date">{{ formatDate(post.createdAt) }}</span>
          </div>
        </div>
        <div class="post-content">
          <p>{{ post.content }}</p>
          <img v-if="post.image" :src="post.image" :alt="post.content" class="post-image" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  posts: {
    type: Array,
    default: () => []
  }
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.post-card {
  background: rgba(15, 15, 23, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
}

.post-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-author {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
}

.post-info h4 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 0.875rem;
  font-weight: 600;
}

.post-date {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.75rem;
}

.post-content p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 16px 0;
  line-height: 1.6;
}

.post-image {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  border-radius: 12px;
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
</style>
