<template>
  <div class="main-layout">
    <!-- Navbar -->
    <nav class="navbar">
      <div class="navbar-left">
        <div class="navbar-logo">
          <img src="/src/assets/logo.png" alt="Social Network Logo" class="logo-image" />
          <span class="navbar-title">Social Network</span>
        </div>
      </div>
      
      <div class="navbar-center">
        <div class="search-bar">
          <MagnifyingGlassIcon class="search-icon" />
          <input type="text" placeholder="Rechercher..." />
        </div>
      </div>
      
      <div class="navbar-right">
        <button class="nav-icon-btn" @click="navigateTo('/notifications')">
          <BellIcon class="icon" />
          <span class="notification-badge">3</span>
        </button>
        
        <button class="nav-icon-btn" @click="navigateTo('/chat')">
          <ChatBubbleLeftRightIcon class="icon" />
        </button>
        
        <div class="navbar-profile">
          <div class="profile-avatar">
            <span>JD</span>
          </div>
          <ChevronDownIcon class="dropdown-icon" />
        </div>
        
        <button class="nav-icon-btn logout-btn" @click="logout">
          <ArrowRightOnRectangleIcon class="icon" />
        </button>
      </div>
    </nav>

    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-section">
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/posts' }"
          @click="navigateTo('/posts')"
        >
          <HomeIcon />
          <span>Accueil</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/profile' }"
          @click="navigateTo('/profile')"
        >
          <UserIcon />
          <span>Profil</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/groups' }"
          @click="navigateTo('/groups')"
        >
          <UserGroupIcon />
          <span>Groupes</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/chat' }"
          @click="navigateTo('/chat')"
        >
          <ChatBubbleLeftRightIcon />
          <span>Messages</span>
        </button>
      </div>
      
      <div class="sidebar-section">
        <div class="sidebar-section-title">Raccourcis</div>
        <button class="sidebar-item">
          <HeartIcon />
          <span>Favoris</span>
        </button>
        
        <button class="sidebar-item">
          <DocumentDuplicateIcon />
          <span>Pages</span>
        </button>
      </div>

      <div class="sidebar-footer">
        <button class="sidebar-item logout-btn" @click="logout">
          <ArrowRightOnRectangleIcon />
          <span>Déconnexion</span>
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="main-content">
      <div class="content-wrapper">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script>
import {
  MagnifyingGlassIcon,
  BellIcon,
  ChatBubbleLeftRightIcon,
  ChevronDownIcon,
  ArrowRightOnRectangleIcon,
  HomeIcon,
  UserIcon,
  UserGroupIcon,
  HeartIcon,
  DocumentDuplicateIcon
} from '@heroicons/vue/24/outline'

export default {
  name: 'MainLayout',
  components: {
    MagnifyingGlassIcon,
    BellIcon,
    ChatBubbleLeftRightIcon,
    ChevronDownIcon,
    ArrowRightOnRectangleIcon,
    HomeIcon,
    UserIcon,
    UserGroupIcon,
    HeartIcon,
    DocumentDuplicateIcon
  },
  data() {
    return {
      currentRoute: this.$route.path
    }
  },
  watch: {
    '$route'(to) {
      this.currentRoute = to.path
    }
  },
  methods: {
    navigateTo(path) {
      this.$router.push(path)
    },
    async logout() {
      try {
        const response = await fetch('http://localhost:8081/logout', {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        // Always redirect to login page after logout attempt
        // Even if the logout request fails, we want to redirect the user
        this.$router.push('/login');
        
        if (!response.ok) {
          console.warn('Logout request failed, but redirecting to login anyway');
        }
      } catch (error) {
        console.error('Error during logout:', error);
        // Still redirect to login page even if there's a network error
        this.$router.push('/login');
      }
    }
  }
}
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  position: relative;
  z-index: 1;
}

/* Navbar Styles */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: rgba(15, 15, 23, 0.9);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 1000;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
}

.navbar-left {
  display: flex;
  align-items: center;
}

.navbar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.navbar-title {
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.navbar-center {
  flex: 1;
  max-width: 400px;
  margin: 0 40px;
}

.search-bar {
  position: relative;
  width: 100%;
}

.search-bar input {
  width: 100%;
  padding: 8px 15px 8px 45px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 25px;
  color: #ffffff;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.search-bar input:focus {
  outline: none;
  border-color: #e879c6;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.2);
}

.search-bar input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.search-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.4);
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.nav-icon-btn {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 8px;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-icon-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
  transform: translateY(-1px);
}

/* Sidebar Icon Styling */
.sidebar-item svg,
.sidebar-item > svg {
  width: 20px;
  height: 20px;
  color: rgba(255, 255, 255, 0.6);
  transition: color 0.3s ease;
}

.sidebar-item:hover svg,
.sidebar-item.active svg {
  color: #e879c6;
}

/* Special styling for notification badge */
.notification-badge {
  position: absolute;
  top: -2px;
  right: -2px;
  background: linear-gradient(135deg, #ff6b6b, #ee5a24);
  color: white;
  font-size: 10px;
  font-weight: bold;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(15, 15, 23, 0.9);
}

.profile-menu {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 5px 10px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.profile-menu:hover {
  background: rgba(255, 255, 255, 0.1);
}

.profile-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid rgba(232, 121, 198, 0.3);
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #000;
  font-weight: 600;
  font-size: 12px;
}

.profile-name {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.9rem;
  font-weight: 500;
}

/* Sidebar Styles */
.sidebar {
  position: fixed;
  top: 60px;
  left: 0;
  width: 250px;
  height: calc(100vh - 60px);
  background: rgba(15, 15, 23, 0.9);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 0;
  overflow-y: auto;
  z-index: 100;
  display: flex;
  flex-direction: column;
}

.sidebar-section {
  margin-bottom: 25px;
}

.sidebar-section-title {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 20px 10px 20px;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 20px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
}

.sidebar-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
  transform: translateX(5px);
}

.sidebar-item.active {
  background: linear-gradient(90deg, rgba(232, 121, 198, 0.2) 0%, transparent 100%);
  color: #e879c6;
  border-right: 3px solid #e879c6;
}

.sidebar-item.active svg {
  color: #e879c6;
}

.sidebar-footer {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1) !important;
  color: #ef4444 !important;
}

/* Main Content Styles */
.main-content {
  margin-left: 250px;
  margin-top: 60px;
  min-height: calc(100vh - 60px);
  position: relative;
  z-index: 1;
}

.content-wrapper {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

/* Icon Styling */
.search-icon,
.icon,
.dropdown-icon {
  width: 20px;
  height: 20px;
  color: rgba(255, 255, 255, 0.7);
  transition: color 0.3s ease;
}

.dropdown-icon {
  width: 16px;
  height: 16px;
}

/* Navbar Icons */
.nav-icon-btn:hover .icon {
  color: #e879c6;
}

/* Logo Styling */
.logo-image {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  animation: logoGlow 3s ease-in-out infinite;
  transition: all 0.3s ease;
  filter: drop-shadow(0 0 10px rgba(232, 121, 198, 0.3));
}

.logo-image:hover {
  transform: scale(1.1);
  filter: drop-shadow(0 0 20px rgba(120, 199, 255, 0.6)) 
          drop-shadow(0 0 30px rgba(232, 121, 198, 0.4));
}

@keyframes logoGlow {
  0%, 100% { 
    filter: drop-shadow(0 0 15px rgba(232, 121, 198, 0.4)) 
            drop-shadow(0 0 25px rgba(168, 85, 247, 0.2));
  }
  50% { 
    filter: drop-shadow(0 0 25px rgba(120, 199, 255, 0.5)) 
            drop-shadow(0 0 35px rgba(232, 121, 198, 0.3));
  }
}

/* Navbar Profile */
.navbar-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 20px;
  transition: background-color 0.3s ease;
}

.navbar-profile:hover {
  background: rgba(255, 255, 255, 0.05);
}

.navbar-profile .profile-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 12px;
}
</style>
