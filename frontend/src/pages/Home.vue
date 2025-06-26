<template>
  <div class="home-page">
    <div class="welcome-section">
      <h1>Bienvenue sur Social Network</h1>
      <p>Connectez-vous avec vos amis et découvrez votre communauté</p>
      
      <div class="quick-stats">
        <div class="stat-card">
          <div class="stat-number">156</div>
          <div class="stat-label">Amis connectés</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">24</div>
          <div class="stat-label">Nouvelles publications</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">8</div>
          <div class="stat-label">Messages non lus</div>
        </div>
      </div>
      
      <div class="quick-actions">
        <button class="action-card" @click="$router.push('/posts')">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
            <polyline points="9,22 9,12 15,12 15,22"/>
          </svg>
          <span>Fil d'actualités</span>
        </button>
        
        <button class="action-card" @click="$router.push('/profile')">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          <span>Mon profil</span>
        </button>
        
        <button class="action-card" @click="$router.push('/chat')">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          <span>Messages</span>
        </button>
        
        <button class="action-card" @click="$router.push('/groups')">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </svg>
          <span>Groupes</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home',
  mounted() {
    // Auto-redirect to posts after a short delay for better UX
    setTimeout(() => {
      this.$router.push('/posts');
    }, 2000);
  }
};
</script>

<style scoped>
.home-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  text-align: center;
}

.welcome-section h1 {
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 30%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 2.5rem;
  font-weight: 900;
  margin: 0 0 15px 0;
  letter-spacing: -0.02em;
}

.welcome-section p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1.2rem;
  margin: 0 0 40px 0;
}

.quick-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 25px 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 40px rgba(232, 121, 198, 0.2);
}

.stat-number {
  font-size: 2rem;
  font-weight: 900;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 5px;
}

.stat-label {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.9rem;
  font-weight: 500;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.action-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 30px 20px;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  transition: all 0.4s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.action-card:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
  transform: translateY(-5px);
  box-shadow: 0 15px 40px rgba(232, 121, 198, 0.3);
}

.action-card svg {
  transition: all 0.3s ease;
}

.action-card:hover svg {
  transform: scale(1.1);
}

.action-card span {
  font-size: 1rem;
  font-weight: 600;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .home-page {
    padding: 20px 15px;
  }
  
  .welcome-section h1 {
    font-size: 2rem;
  }
  
  .welcome-section p {
    font-size: 1rem;
  }
  
  .quick-stats {
    grid-template-columns: 1fr;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
}
</style>
