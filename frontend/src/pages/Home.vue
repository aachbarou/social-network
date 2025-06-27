<template>
  <div class="posts-page">
    <!-- Header -->
    <div class="page-header">
      <h1>Fil d'actualités</h1>
      <p>Découvrez les dernières publications de votre communauté</p>
    </div>

    <!-- Create Post Section -->
    <div class="create-post-card">
      <div class="create-post-header">
        <div class="user-avatar">VU</div>
        <input 
          type="text" 
          placeholder="Quoi de neuf aujourd'hui ?" 
          class="create-post-input"
          v-model="newPostContent"
          @click="showCreatePost = true"
        />
                     <div>showCreatePost: {{ showCreatePost }}</div>

      </div>
      <div class="create-post-actions" v-if="showCreatePost">
        <button class="action-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21,15 16,10 5,21"/>
          </svg>
          Photo
        </button>
        <button class="action-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="23 7 16 12 23 17 23 7"/>
            <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
          </svg>
          Vidéo
        </button>
        <button class="action-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 1v6m0 6v6m11-7h-6m-6 0H1"/>
          </svg>
          Sondage
        </button>
        <button class="publish-btn" @click="publishPost">Publier</button>
      </div>
    </div>

    <!-- Debug Print for showCreatePost -->
 

    <!-- Posts Feed -->
    <div class="posts-feed">
      <template v-for="post in formattedPosts" :key="post && post.id">
        <div v-if="post">
          <div class="post-card">
            <div class="post-header">
              <div class="user-avatar">{{ post.userInitials }}</div>
              <div class="post-meta">
                <h3 class="user-name">{{ post.userName }}</h3>
                <span class="post-time">{{ post.timeAgo }}</span>
              </div>
              <button class="post-menu-btn">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                  <circle cx="12" cy="12" r="1"/>
                  <circle cx="19" cy="12" r="1"/>
                  <circle cx="5" cy="12" r="1"/>
                </svg>
              </button>
            </div>
            <div class="post-content">
              <p>{{ post.content }}</p>
              <div v-if="post.hasImage" class="post-image-placeholder">
                <svg width="60" height="60" viewBox="0 0 24 24" fill="currentColor">
                  <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                  <circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21,15 16,10 5,21"/>
                </svg>
                <span>Image</span>
              </div>
            </div>
            <div class="post-actions">
              <button class="action-btn">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
                </svg>
                {{ post.likes }}
              </button>
              <button class="action-btn">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
                </svg>
                {{ post.comments }}
              </button>
              <button class="action-btn">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
                  <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
                </svg>
                Partager
              </button>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import { useMainStore } from '../stores/postsStore';
import { storeToRefs } from 'pinia';
import { ref, computed } from 'vue';

export default {
  name: 'Posts',
  setup() {
    const mainStore = useMainStore();
    const { posts } = storeToRefs(mainStore);
    mainStore.fetchPosts(); // Fetch posts on component mount
    
    // For new post creation
    const newPostContent = ref('');
    const showCreatePost = ref(false);
    const publishPost = async () => {
      const formData = new FormData();
      formData.set('body', newPostContent.value);
      // Add more fields as needed (privacy, groupId, etc.)
      await fetch('http://localhost:8081/newPost', {
        method: 'POST',
        credentials: 'include',
        body: formData
      });
      newPostContent.value = '';
      showCreatePost.value = false;
      await mainStore.fetchPosts();
    };
    
    const formattedPosts = computed(() => {
      // Ensure posts.value is always an array
      const arr = Array.isArray(posts.value) ? posts.value : [];
      return arr
        .filter(post => post) // filter out nulls
        .map(post => ({

          id: post.id,
          userName: post.author?.nickname || 'Utilisateur',
          userInitials: post.author?.nickname ? post.author.nickname.slice(0,2).toUpperCase() : '??',
          timeAgo: post.createdAt ? new Date(post.createdAt).toLocaleString() : '',
          content: post.content,
          hasImage: !!post.imagePath,
          likes: post.likes || 0,
          comments: post.comments ? post.comments.length : 0
        }));
    });
    
    return {
      posts,
      mainStore,
      newPostContent,
      showCreatePost,
      publishPost,
      formattedPosts
    };
  },
  data() {
    return {
      showCreatePost: false
    }
  }
}
</script>

<style scoped>
.posts-page {
  max-width: 600px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-header h1 {
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 30%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 2rem;
  font-weight: 900;
  margin: 0 0 10px 0;
  letter-spacing: -0.02em;
}

.page-header p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
  margin: 0;
}

.create-post-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  margin-bottom: 30px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.create-post-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid rgba(232, 121, 198, 0.3);
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #000;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.create-post-input {
  flex: 1;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 25px;
  color: #ffffff;
  font-size: 0.95rem;
  transition: all 0.3s ease;
}

.create-post-input:focus {
  outline: none;
  border-color: #e879c6;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.2);
}

.create-post-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.create-post-actions {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 15px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
  transform: translateY(-1px);
}

.publish-btn {
  padding: 8px 20px;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border: none;
  border-radius: 12px;
  color: #000000;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-left: auto;
}

.publish-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(232, 121, 198, 0.4);
}

.posts-feed {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.post-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.post-meta {
  flex: 1;
}

.user-name {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.95rem;
  font-weight: 600;
  margin: 0 0 2px 0;
}

.post-time {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
}

.post-menu-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  padding: 5px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.post-menu-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

.post-content p {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin: 0 0 15px 0;
}

.post-image {
  width: 100%;
  border-radius: 12px;
  margin-bottom: 15px;
  object-fit: cover;
}

.post-image-placeholder {
  width: 100%;
  height: 200px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.4);
  gap: 10px;
}

.post-image-placeholder span {
  font-size: 0.9rem;
  font-weight: 500;
}

.post-actions {
  display: flex;
  align-items: center;
  gap: 20px;
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.post-actions .action-btn {
  background: none;
  border: none;
  padding: 8px 12px;
  font-size: 0.85rem;
}

.post-actions .action-btn:hover {
  background: rgba(255, 255, 255, 0.05);
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .posts-page {
    padding: 0 15px;
  }
  
  .create-post-actions {
    flex-direction: column;
    align-items: stretch;
  }
  
  .publish-btn {
    margin-left: 0;
    align-self: flex-end;
  }
  
  .page-header h1 {
    font-size: 1.6rem;
  }
}
</style>