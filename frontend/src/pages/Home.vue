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

      </div>
      <div class="create-post-actions" v-if="showCreatePost">
        <!-- Show followers selection only for private posts -->
        <div v-if="selectedPrivacy === 'private'" class="followers-selector">
          <label>Sélectionnez les followers autorisés :</label>
          <div v-if="followers.length === 0">Aucun follower à afficher.</div>
          <div v-else class="followers-list">
            <label v-for="f in followers" :key="f.id" class="follower-checkbox">
              <input type="checkbox" :value="f.id" v-model="selectedFollowers" />
              {{ f.nickname || f.firstName || f.lastName || f.id }}
            </label>
          </div>
        </div>
        <!-- Info text for almost private posts -->
        <div v-if="selectedPrivacy === 'almost_private'" class="almost-private-info">
          <p>
            <Info :size="16" />
            Ce post sera visible par tous vos followers automatiquement.
          </p>
        </div>
        
        <div class="action-buttons">
          <!-- Privacy Toggle Button -->
          <button class="privacy-toggle-btn" @click="togglePrivacy" :title="getPrivacyTooltip()">
            <component :is="getPrivacyIconComponent()" :size="16" />
            <span class="privacy-text">{{ getPrivacyText() }}</span>
          </button>
          
          <label class="action-btn">
            <input type="file" ref="imageInput" @change="handleImageChange" accept="image/*" style="display: none" />
            <Image :size="18" />
            Photo
          </label>
          
          <button class="publish-btn" @click="publishPost">Publier</button>
        </div>
      </div>
      <div v-if="showCreatePost && selectedImage" class="image-preview-container">
        <img :src="imagePreview" class="image-preview" />
        <button class="remove-image-btn" @click="removeImage">×</button>
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
                <div class="post-time-container">
                  <component :is="getPostPrivacyIconComponent(post.visibility)" :size="14" class="privacy-icon-component" />
                  <span class="post-time">{{ post.timeAgo }}</span>
                </div>
              </div>
              <button class="post-menu-btn">
                <MoreHorizontal :size="20" />
              </button>
            </div>
            <div class="post-content">
              <p>{{ post.content }}</p>
              <div v-if="post.hasImage" class="post-image-container">
                <img v-if="post.imagePath" :src="getImageUrl(post.imagePath)" class="post-image" alt="Post image" @error="handleImageError" />
                <div v-else class="post-image-placeholder">
                  <Image :size="60" />
                  <span>Image</span>
                </div>
              </div>
              <pre v-if="false">Debug post: {{ JSON.stringify(post, null, 2) }}</pre>
            </div>
            <div class="post-actions">
              <button class="action-btn" @click="toggleComments(post.id)">
                <MessageCircle :size="18" />
                {{ post.commentsCount || 0 }} Commentaire{{ (post.commentsCount || 0) > 1 ? 's' : '' }}
              </button>
            </div>
            
            <!-- Comments Section -->
            <div v-if="activeCommentPostId === post.id" class="comments-section">
              <div v-if="post.comments && post.comments.length > 0" class="comments-list">
                <div v-for="comment in post.comments" :key="comment.id" class="comment">
                  <div class="comment-header">
                    <div class="user-avatar small">{{ comment.userInitials || (comment.author?.nickname ? comment.author.nickname.slice(0,2).toUpperCase() : 'U') }}</div>
                    <div class="comment-meta">
                      <h4 class="comment-user">
                        {{ comment.userName || comment.author?.nickname || 'Utilisateur' }}
                      </h4>
                      <span class="comment-time">
                        {{ isValidDate(comment.createdAt) ? formatDate(comment.createdAt) : '' }}
                      </span>
                    </div>
                  </div>
                  <div class="comment-content">
                    <p>{{ comment.content }}</p>
                    <img v-if="getCommentImage(comment)" :src="getImageUrl(getCommentImage(comment))" class="comment-image" @error="handleImageError" />
                  </div>
                </div>
              </div>
              <div class="comment-form">
                <div class="user-avatar small">VU</div>
                <input 
                  type="text" 
                  v-model="newCommentText" 
                  placeholder="Écrire un commentaire..." 
                  class="comment-input"
                  @keyup.enter="submitComment(post.id)"
                />
                <label class="comment-attach-btn">
                  <input type="file" @change="handleCommentImageChange" accept="image/*" style="display: none" />
                  <Image :size="18" />
                </label>
                <button class="comment-send-btn" @click="submitComment(post.id)">
                  <Send :size="18" />
                </button>
              </div>
              <div v-if="commentImagePreview" class="comment-image-preview-container">
                <img :src="commentImagePreview" class="comment-image-preview" />
                <button class="remove-image-btn" @click="removeCommentImage">×</button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import { useMainStore } from '../stores/postsStore';
import { useFollowStore } from '../stores/followStore';
import { useUserStore } from '../stores/userStore';
import { storeToRefs } from 'pinia';
import { ref, computed, watch } from 'vue';
import { Image, MoreHorizontal, MessageCircle, Send, Globe, Users, Lock, Info } from 'lucide-vue-next';

export default {
  name: 'Posts',
  components: {
    Image,
    MoreHorizontal,
    MessageCircle,
    Send,
    Globe,
    Users,
    Lock,
    Info
  },
  setup() {
    const mainStore = useMainStore();
    const followStore = useFollowStore();
    const userStore = useUserStore();
    const { posts } = storeToRefs(mainStore);
    const { followers } = storeToRefs(followStore);
    mainStore.fetchPosts(); // Fetch posts on component mount
    
    // For new post creation
    const newPostContent = ref('');
    const showCreatePost = ref(false);
    const selectedImage = ref(null);
    const imagePreview = ref('');
    const imageInput = ref(null);
    const selectedPrivacy = ref('public');
    const selectedFollowers = ref([]);
    
    // Fetch followers when privacy changes to PRIVATE or ALMOST_PRIVATE
    watch(selectedPrivacy, async (val) => {
      if ((val === 'private' || val === 'almost_private') && userStore.user?.id) {
        await followStore.fetchFollowers(userStore.user.id);
      }
    });
    
    // For comments
    const activeCommentPostId = ref(null);
    const newCommentText = ref('');
    const commentSelectedImage = ref(null);
    const commentImagePreview = ref('');
    
    const handleImageChange = (event) => {
      const file = event.target.files[0];
      if (file) {
        selectedImage.value = file;
        imagePreview.value = URL.createObjectURL(file);
      }
    };
    
    const removeImage = () => {
      selectedImage.value = null;
      imagePreview.value = '';
      if (imageInput.value) {
        imageInput.value.value = '';
      }
    };
    
    const handleCommentImageChange = (event) => {
      const file = event.target.files[0];
      if (file) {
        commentSelectedImage.value = file;
        commentImagePreview.value = URL.createObjectURL(file);
      }
    };
    
    const removeCommentImage = () => {
      commentSelectedImage.value = null;
      commentImagePreview.value = '';
    };
    
    const toggleComments = (postId) => {
      activeCommentPostId.value = activeCommentPostId.value === postId ? null : postId;
      newCommentText.value = '';
      removeCommentImage();
    };
    
    const getImageUrl = (path) => {
      if (!path) return '';
      
      // If the path already includes http:// or https://, use it as is
      if (path.startsWith('http://') || path.startsWith('https://')) {
        return path;
      }
      
      // Make sure path doesn't have double slashes when concatenating
      const cleanPath = path.startsWith('/') ? path.substring(1) : path;
      return `http://localhost:8081/${cleanPath}`;
    };
    
    const formatDate = (dateString) => {
      // Handle completely missing date - provide a friendly fallback
      if (!dateString || dateString === "" || dateString === null || dateString === undefined) {
        return 'Récemment'  // More user-friendly than "Date non disponible"
      }
      
      try {
        let date;
        
        // Handle Unix timestamp (both string and number)
        if (typeof dateString === 'string' && /^\d+$/.test(dateString)) {
          const timestamp = parseInt(dateString)
          // Try both seconds and milliseconds
          date = timestamp > 10000000000 ? new Date(timestamp) : new Date(timestamp * 1000)
        } else if (typeof dateString === 'number') {
          // Check if it's Unix timestamp (seconds) or JavaScript timestamp (milliseconds)
          date = dateString < 10000000000 ? new Date(dateString * 1000) : new Date(dateString)
        } else if (typeof dateString === 'string') {
          // Try parsing as ISO string
          date = new Date(dateString)
        } else {
          date = new Date(dateString)
        }
        
        if (isNaN(date.getTime())) {
          return 'Récemment'  // Fallback for invalid dates
        }
        
        const now = new Date()
        const diffMs = now - date
        const diffMins = Math.floor(diffMs / 60000)
        const diffHours = Math.floor(diffMins / 60)
        const diffDays = Math.floor(diffHours / 24)
        
        if (diffMins < 1) return 'À l\'instant'
        if (diffMins < 60) return `Il y a ${diffMins} min`
        if (diffHours < 24) return `Il y a ${diffHours}h`
        if (diffDays < 7) return `Il y a ${diffDays} jour${diffDays > 1 ? 's' : ''}`
        
        return date.toLocaleDateString('fr-FR', {
          day: 'numeric',
          month: 'short',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch (error) {
        return 'Récemment'  // Fallback for any errors
      }
    };
    
    const handleImageError = (event) => {
      console.error('Image failed to load:', event.target.src);
      // Fall back to placeholder if image fails to load
      event.target.style.display = 'none';
      event.target.parentElement.innerHTML += `
        <div class="post-image-placeholder">
          <svg width="60" height="60" viewBox="0 0 24 24" fill="currentColor">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21,15 16,10 5,21"/>
          </svg>
        </div>
        <span>Image non disponible</span>
      `;
      event.target.parentElement.appendChild(placeholder);
    };
    
    const submitComment = async (postId) => {
      if (!newCommentText.value && !commentSelectedImage.value) return;
      
      await mainStore.submitComment(postId, {
        body: newCommentText.value,
        image: commentSelectedImage.value
      });
      
      newCommentText.value = '';
      removeCommentImage();
    };
    
    // Privacy toggle functionality
    const togglePrivacy = () => {
      const privacyOptions = ['public', 'almost_private', 'private'];
      const currentIndex = privacyOptions.indexOf(selectedPrivacy.value);
      const nextIndex = (currentIndex + 1) % privacyOptions.length;
      selectedPrivacy.value = privacyOptions[nextIndex];
    };
    
    const getPrivacyIconComponent = () => {
      switch (selectedPrivacy.value) {
        case 'public': return 'Globe';
        case 'almost_private': return 'Users';
        case 'private': return 'Lock';
        default: return 'Globe';
      }
    };
    
    const getPrivacyText = () => {
      switch (selectedPrivacy.value) {
        case 'public': return 'Public';
        case 'almost_private': return 'Presque privé';
        case 'private': return 'Privé';
        default: return 'Public';
      }
    };
    
    const getPrivacyTooltip = () => {
      switch (selectedPrivacy.value) {
        case 'public': return 'Cliquez pour changer en Presque privé';
        case 'almost_private': return 'Cliquez pour changer en Privé';
        case 'private': return 'Cliquez pour changer en Public';
        default: return 'Cliquez pour changer la confidentialité';
      }
    };
    
    const publishPost = async () => {
      if (!newPostContent.value && !selectedImage.value) return;
      const formData = new FormData();
      formData.set('body', newPostContent.value);
      formData.set('privacy', selectedPrivacy.value.toUpperCase().replace('_', '_'));
      // Only send selected followers for private posts, not for almost_private
      if (selectedPrivacy.value === 'private') {
        formData.set('checkedfollowers', selectedFollowers.value.join(','));
      }
      if (selectedImage.value) {
        formData.append('image', selectedImage.value);
      }
      
      try {
        const response = await fetch('http://localhost:8081/newPost', {
          method: 'POST',
          credentials: 'include',
          body: formData
        });
        
        const result = await response.json();
        
        if (!response.ok) {
          console.error("Error creating post:", result);
          alert("Erreur lors de la création du post: " + (result.message || "Erreur inconnue"));
          return;
        }
        
        newPostContent.value = '';
        removeImage();
        showCreatePost.value = false;
                selectedPrivacy.value = 'public';
        selectedFollowers.value = [];
        await mainStore.fetchPosts();
      } catch (error) {
        console.error("Post creation failed:", error);
        alert("Erreur de connexion lors de la création du post");
      }
    };
    
    const formattedPosts = computed(() => {
      // Ensure posts.value is always an array
      const arr = Array.isArray(posts.value) ? posts.value : [];
      
      const formattedArr = arr
        .filter(post => post) // filter out nulls
        .map(post => {
          // Check for image path in various possible formats
          const imagePath = post.imagePath || post.image || (post.images && post.images.length > 0 ? post.images[0] : null);
          
          return {
            id: post.id,
            userName: post.author?.nickname || 'Utilisateur',
            userInitials: post.author?.nickname ? post.author.nickname.slice(0,2).toUpperCase() : '??',
            timeAgo: post.createdAt ? formatDate(post.createdAt) : '',
            content: post.content,
            hasImage: !!imagePath,
            imagePath: imagePath,
            likes: post.likes || 0,
            comments: post.comments || [],
            commentsCount: post.comments ? post.comments.length : 0,
            visibility: post.visibility || 'PUBLIC',
          };
        });
      
      return formattedArr;
    });
    // Helper to get privacy label
    function getPrivacyLabel(visibility) {
      if (visibility === 'PUBLIC') return 'Public';
      if (visibility === 'ALMOST_PRIVATE') return 'Presque privé';
      if (visibility === 'PRIVATE') return 'Privé';
      return '';
    }
    
    // Add this helper function inside setup()
    const getCommentImage = (comment) => {
      return comment.imagePath || comment.image || (comment.images && comment.images.length > 0 ? comment.images[0] : null);
    };
    
    function isValidDate(date) {
      return !isNaN(new Date(date).getTime());
    }
    
    const getPostPrivacyIconComponent = (visibility) => {
      switch (visibility) {
        case 'PUBLIC': return 'Globe';
        case 'ALMOST_PRIVATE': return 'Users';
        case 'PRIVATE': return 'Lock';
        default: return 'Globe';
      }
    };
    
    return {
      posts,
      mainStore,
      newPostContent,
      showCreatePost,
      selectedImage,
      imagePreview,
      imageInput,
      handleImageChange,
      removeImage,
      publishPost,
      formattedPosts,
      activeCommentPostId,
      newCommentText,
      commentSelectedImage,
      commentImagePreview,
      toggleComments,
      submitComment,
      handleCommentImageChange,
      removeCommentImage,
      getImageUrl,
      formatDate,
      handleImageError,
      getPrivacyLabel,
      getPostPrivacyIconComponent,
      getCommentImage,
      isValidDate,
      selectedPrivacy,
      followers,
      selectedFollowers,
      togglePrivacy,
      getPrivacyIconComponent,
      getPrivacyText,
      getPrivacyTooltip,
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
  flex-direction: column;
  gap: 15px;
}

.action-buttons {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 15px;
}

.privacy-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.privacy-toggle-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
  transform: translateY(-1px);
}

.privacy-icon {
  font-size: 1rem;
  line-height: 1;
}

.privacy-text {
  font-size: 0.85rem;
  font-weight: 500;
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

.post-time-container {
  display: flex;
  align-items: center;
  gap: 6px;
}

.post-time-container {
  display: flex;
  align-items: center;
  gap: 6px;
}

.privacy-icon-component {
  opacity: 1;
  color: rgba(255, 255, 255, 0.8);
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

.image-preview-container {
  position: relative;
  margin-top: 15px;
}

.image-preview {
  width: 100%;
  max-height: 250px;
  border-radius: 12px;
  object-fit: cover;
}

.remove-image-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.remove-image-btn:hover {
  background: rgba(255, 0, 0, 0.7);
  transform: scale(1.1);
}

.comments-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.comments-list {
  margin-bottom: 15px;
}

.comment {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 10px;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.user-avatar.small {
  width: 30px;
  height: 30px;
  font-size: 12px;
}

.comment-meta {
  flex: 1;
}

.comment-user {
  font-size: 0.85rem;
  margin: 0 0 2px 0;
  color: rgba(255, 255, 255, 0.9);
}

.comment-time {
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.5);
}

.comment-content {
  color: #fff;
}

.comment-content p {
  font-size: 0.9rem;
  margin: 0 0 8px 0;
}

.comment-image {
  max-width: 100%;
  max-height: 150px;
  border-radius: 8px;
  margin-top: 8px;
}

.comment-form {
  display: flex;
  align-items: center;
  gap: 10px;
}

.comment-input {
  flex: 1;
  padding: 8px 15px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: white;
  font-size: 0.9rem;
}

.comment-input:focus {
  outline: none;
  border-color: #e879c6;
  box-shadow: 0 0 0 2px rgba(232, 121, 198, 0.2);
}

.comment-attach-btn,
.comment-send-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.7);
  transition: all 0.2s ease;
}

.comment-attach-btn:hover,
.comment-send-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
}

.comment-send-btn {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border: none;
  color: #000000;
}

.comment-image-preview-container {
  margin-top: 10px;
  position: relative;
  width: 150px;
}

.followers-selector {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 16px;
}

.followers-selector label {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.9rem;
  font-weight: 500;
  margin-bottom: 12px;
  display: block;
}

.almost-private-info {
  background: rgba(52, 152, 219, 0.1);
  border: 1px solid rgba(52, 152, 219, 0.3);
  border-radius: 12px;
  padding: 16px;
}

.almost-private-info p {
  margin: 0;
  color: #5dade2;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 8px;
}

.followers-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.follower-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.07);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.8);
}

.follower-checkbox:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
}

.follower-checkbox input[type="checkbox"] {
  accent-color: #e879c6;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .posts-page {
    padding: 0 15px;
  }
  
  .action-buttons {
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .privacy-toggle-btn {
    order: -1;
    width: 100%;
    justify-content: center;
  }
  
  .publish-btn {
    width: 100%;
    justify-content: center;
  }
  
  .page-header h1 {
    font-size: 1.6rem;
  }
}
</style>