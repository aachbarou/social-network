<template>
  <div class="chat-container">
    <div class="chat-layout">
      <!-- Users/Friends List Section -->
      <div class="friends-list-section">
        <!-- Header -->
        <div class="header">
          <div class="header-icon-wrapper">
            <svg class="header-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path>
            </svg>
          </div>
          <div class="header-title">ChatApp</div>
        </div>

        <!-- Current User Profile -->
        <div v-if="currentUser" class="profile-card">
          <div class="profile-avatar-wrapper">
            <img 
                 :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)"
                 alt="Avatar" class="profile-avatar" style="max-width:120px; max-height:120px; object-fit:contain; border:2px solid #e879c6; background:#fff;" />
        
          </div>
          <div class="profile-name">{{ currentUser.nickname || `${currentUser.firstName} ${currentUser.lastName}` }}</div>
          <div class="profile-email">{{ currentUser.email }}</div>
        </div>

        <!-- Friends List -->
        <div class="friends-list-container">
          <div class="friends-list-header">
            <span class="friends-list-title">Friends</span>
            <span class="friends-list-count">{{ following.length }}</span>
          </div>
          <div class="friends-list custom-scrollbar">
            <button
              v-for="friend in following"
              :key="friend.id"
              @click="selectUser(friend)"
              class="friend-item"
              :class="{ 'active': selectedUser && selectedUser.id === friend.id }"
            >
              <div class="friend-avatar-wrapper">
                 <img 
                 :src="getFullImageUrl(friend.avatar || friend.imagePath)"
                 alt="Avatar" class="profile-avatar" style="max-width:120px; max-height:120px; object-fit:contain; border:2px solid #e879c6; background:#fff;" />
              </div>
              <div class="friend-name">{{ friend.nickname || friend.firstName }}</div>
            </button>
          </div>
        </div>
      </div>

      <!-- Chat Section -->
      <div class="chat-area-wrapper">
        <div class="chat-box">
          <!-- Chat Active State -->
          <div v-if="selectedUser" class="active-chat">
            <!-- Chat Header -->
            <div class="chat-header">
              <div class="chat-header-user">
                 <img 
                 :src="getFullImageUrl(selectedUser.avatar || selectedUser.imagePath)"
                 alt="Avatar" class="chat-header-avatar" style="max-width:120px; max-height:120px; object-fit:contain; border:2px solid #e879c6; background:#fff;" />
      
                <div class="chat-header-text">
                  <div class="chat-header-name-wrapper">
                    <span class="chat-header-name">{{ selectedUser.nickname || selectedUser.firstName }}</span>
                  </div>
                  <span class="chat-header-email">{{ selectedUser.email }}</span>
                </div>
              </div>
               <button class="reset-chat-btn" @click="resetChat">Reset Chat</button>
            </div>
            
            <!-- Messages -->
            <div ref="messageContainer" class="messages-container custom-scrollbar">
              <div class="messages-list">
                <div class="messages-grid">
                  <template v-for="(message, index) in messages[selectedUser.id]" :key="index">
                    <!-- Message from other user -->
                    <div v-if="message.sender !== currentUser.id" class="message-wrapper received">
                      <div class="message-layout">
                        <div class="message-avatar-wrapper">
                          <img 
                 :src="getFullImageUrl(selectedUser.avatar || selectedUser.imagePath)"
                 alt="Avatar" class="chat-header-avatar" style="max-width:120px; max-height:120px; object-fit:contain; border:2px solid #e879c6; background:#fff;" />
      
                        </div>
                        <div class="message-bubble received">
                          <div>{{ message.text }}</div>
                        </div>
                      </div>
                    </div>
                    <!-- Message from current user -->
                    <div v-else class="message-wrapper sent">
                      <div class="message-layout sent">
                        <div class="message-avatar-wrapper">
                          <img 
                 :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)"
                 alt="Avatar" class="chat-header-avatar" style="max-width:120px; max-height:120px; object-fit:contain; border:2px solid #e879c6; background:#fff;" />
      
                        </div>
                        <div class="message-bubble sent">
                          <div>{{ message.text }}</div>
                        </div>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </div>

            <!-- Message Input -->
            <div class="message-input-section">
              <div class="input-wrapper">
                <div class="input-relative-wrapper">
                  <input
                    type="text"
                    v-model="newMessage"
                    @keyup.enter="sendMessage"
                    class="message-input"
                    placeholder="Type your message..."
                  />
                </div>
              </div>
              <div class="send-button-wrapper">
                <button @click="sendMessage" class="send-button">
                  <span>Send</span>
                  <span class="send-button-icon-wrapper">
                    <svg class="send-button-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
                    </svg>
                  </span>
                </button>
              </div>
            </div>
          </div>
          <!-- Chat Empty State -->
          <div v-else class="empty-chat-container">
            <div class="empty-chat-icon-wrapper">
              <svg class="empty-chat-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
            </div>
            <h2 class="empty-chat-title">Select a conversation</h2>
            <p class="empty-chat-text">Choose one of your friends from the list on the left to start chatting.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue';
import { useUserStore } from '../stores/userStore';
import { useFollowStore } from '../stores/followStore';
import { storeToRefs } from 'pinia';

const userStore = useUserStore();
const followStore = useFollowStore();

const { user: currentUser } = storeToRefs(userStore);
const { following } = storeToRefs(followStore);

const selectedUser = ref(null);
const newMessage = ref('');
const messageContainer = ref(null);
const messages = reactive({});

//  Get  avatar 
//user[currentUser].avatar =" https://images.app.goo.gl/mU2UBreWdeTh1Emj9"
function getFullImageUrl(path) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `http://localhost:8081/${path.replace(/^\/+/,'')}`
}

/**
 * Scrolls the message container to the most recent message.
 */
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
    }
  });
};

/**
 * Sets the currently active user for the chat window.
 * @param {object} user - The user object to select.
 */
const selectUser = (user) => {
  selectedUser.value = user;
  if (!messages[user.id]) {
    messages[user.id] = [];
  }
  scrollToBottom();
};

/**
 * Adds the new message to the conversation.
 */
const sendMessage = () => {
  if (newMessage.value.trim() === '' || !selectedUser.value) {
    return;
  }

  messages[selectedUser.value.id].push({
    sender: currentUser.value.id,
    text: newMessage.value.trim(),
  });

  newMessage.value = '';
  scrollToBottom();
};

/**
 * Clears all messages for the currently selected chat.
 */
const resetChat = () => {
    if (selectedUser.value) {
        messages[selectedUser.value.id] = [];
    }
}

/**
 * On component mount, fetches user data and their following list.
 */
onMounted(async () => {
  await userStore.reloadUserAfterRefresh();
  // Download  the  avatar  
  if (currentUser.value && currentUser.value.id) {
    await followStore.fetchFollowing(currentUser.value.id);

    if (following.value.length > 0) {
      selectUser(following.value[0]);
    }
  }
});
</script>

<style scoped>
:root {
  --interactive-gradient: linear-gradient(135deg, #c0dd1b 0%, #73db1e 100%);
  --bg-primary: #111118; /* Was #0a0a0f */
  --bg-secondary: rgba(35, 36, 50, 0.8); /* Was rgba(25, 26, 38, 0.8) */
  --bg-tertiary: rgba(55, 58, 78, 0.9); /* Was rgba(40, 42, 58, 0.9) */
  --border-color: rgba(255, 255, 255, 0.15); /* Was 0.1 */
  --text-primary: #c90505;
  --text-secondary: rgba(255, 255, 255, 0.85); /* Was 0.7 */
  --text-muted: rgba(255, 255, 255, 0.6); /* Was 0.5 */
}

.chat-container {
  display: flex;
  height: 100vh;
  color: var(--text-secondary);
  background-color: var(--bg-primary);
}

.chat-layout {
  display: flex;
  flex-direction: row;
  height: 100%;
  width: 100%;
  overflow-x: hidden;
}

/* --- Sidebar --- */
.friends-list-section {
  display: flex;
  flex-direction: column;
  padding: 1.5rem 0.5rem 1.5rem 1.5rem;
  width: 20rem;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  flex-shrink: 0;
  backdrop-filter: blur(10px);
}

.header {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 3rem;
  width: 100%;
  flex-shrink: 0;
  margin-bottom: 1rem;
}

.header-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.75rem;
  color: var(--text-primary);
  background: var(--interactive-gradient);
  height: 2.5rem;
  width: 2.5rem;
}
.header-icon {
  width: 1.5rem;
  height: 1.5rem;
}

.header-title {
  margin-left: 0.75rem;
  font-weight: 700;
  font-size: 1.5rem;
  color: var(--text-primary);
}

/* --- Profile Card --- */
.profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  background: rgba(0,0,0,0.2);
  border: 1px solid var(--border-color);
  margin-top: 1rem;
  width: 100%;
  padding: 1.5rem 1rem;
  border-radius: 1rem;
  flex-shrink: 0;
}
.profile-avatar-wrapper {
  height: 5rem;
  width: 5rem;
  border-radius: 50%;
  border: 2px solid var(--border-color);
  overflow: hidden;
}
.profile-avatar {
  height: 100%;
  width: 100%;
  object-fit: cover;
}
.profile-name {
  font-size: 1rem;
  font-weight: 600;
  margin-top: 0.75rem;
  color: var(--text-primary);
}
.profile-email {
  font-size: 0.875rem;
  color: var(--text-muted);
}

/* --- Friends List --- */
.friends-list-container {
  display: flex;
  flex-direction: column;
  margin-top: 1.5rem;
  flex-grow: 1;
  min-height: 0;
}
.friends-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
  padding: 0 0.5rem;
  color: var(--text-muted);
}
.friends-list-title {
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.friends-list-count {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(82, 79, 79, 0.1);
  height: 1.25rem;
  min-width: 1.25rem;
  padding: 0 0.25rem;
  border-radius: 0.625rem;
  font-size: 0.75rem;
  font-weight: bold;
  color: var(--text-secondary);
}
.friends-list {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  margin-top: 0.5rem;
  flex-grow: 1;
  overflow-y: auto;
  padding-right: 0.5rem;
}
.friend-item {
  display: flex;
  align-items: center;
  border-radius: 0.75rem;
  padding: 0.75rem;
  width: 100%;
  text-align: left;
  border: 1px solid transparent;
  background: none;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}
.friend-item:hover {
  background-color: rgba(255,255,255,0.05);
  border-color: rgba(255,255,255,0.1);
}
.friend-item.active {
  color: white;
  background: var(--interactive-gradient);
}
.friend-avatar-wrapper {
  height: 2.5rem;
  width: 2.5rem;
  border-radius: 50%;
  flex-shrink: 0;
  overflow: hidden;
}
.friend-avatar {
  height: 100%;
  width: 100%;
  object-fit: cover;
}
.friend-name {
  margin-left: 0.75rem;
  font-size: 1rem;
  font-weight: 500;
  color: inherit;
}

/* --- Chat Section --- */
.chat-area-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1 1 0%;
  height: 100%;
  padding: 1.5rem;
}
.chat-box {
  display: flex;
  flex-direction: column;
  flex: 1 1 0%;
  border-radius: 1.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  height: 100%;
  padding: 1rem;
  min-height: 0;
  backdrop-filter: blur(10px);
}

/* --- Active Chat --- */
.active-chat {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}
.chat-header-user {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.chat-header-avatar {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  object-fit: cover;
}
.chat-header-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
}
.chat-header-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}
.chat-header-email {
  font-size: 0.875rem;
  color: var(--text-muted);
}
.reset-chat-btn {
  padding: 0.5rem 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255, 107, 107, 0.3);
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
  color: white;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}
.reset-chat-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.2);
}

/* --- Messages --- */
.messages-container {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow-y: auto;
  margin-bottom: 1rem;
  min-height: 0;
}
.messages-list {
  margin-top: auto;
}
.messages-grid {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 1rem 0;
  padding: 1rem 0.5rem;
}
.message-wrapper {
  display: flex;
  max-width: 80%;
}
.message-wrapper.received {
  grid-column: 1 / span 10;
  justify-content: flex-start;
}
.message-wrapper.sent {
  grid-column: 4 / span 10;
  justify-content: flex-end;
}
.message-layout {
  display: flex;
  align-items: flex-end;
  flex-direction: row;
  gap: 0.75rem;
}
.message-layout.sent {
  flex-direction: row-reverse;
}
.message-avatar-wrapper {
  height: 2.5rem;
  width: 2.5rem;
  border-radius: 50%;
  flex-shrink: 0;
  overflow: hidden;
}
.message-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.message-bubble {
  position: relative;
  font-size: 1rem;
  padding: 0.75rem 1rem;
  border-radius: 1.25rem;
  color: var(--text-secondary);
}
.message-bubble.received {
  background-color: var(--bg-tertiary);
  border-top-left-radius: 0.25rem;
}
.message-bubble.sent {
  background: var(--interactive-gradient);
  color: var(--text-primary);
  border-top-right-radius: 0.25rem;
}

/* --- Message Input --- */
.message-input-section {
  display: flex;
  align-items: center;
  height: 4rem;
  width: 100%;
  padding: 0;
  flex-shrink: 0;
}
.input-wrapper {
  flex-grow: 1;
  position: relative;
}
.message-input {
  width: 100%;
  background-color: rgba(0,0,0,0.2);
  border: 1px solid var(--border-color);
  border-radius: 1rem;
  padding-left: 1rem;
  padding-right: 1rem;
  height: 3.5rem;
  color: var(--text-primary);
  font-size: 1rem;
  transition: all 0.2s;
}
.message-input:focus {
  outline: none;
  border-color: #e879c6;
  background: var(--bg-tertiary);
  box-shadow: 0 0 0 2px rgba(232, 121, 198, 0.5);
}
.send-button-wrapper {
  margin-left: 1rem;
}
.send-button {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--interactive-gradient);
  border-radius: 1rem;
  color: rgb(240, 219, 219);
  height: 3.5rem;
  width: 3.5rem;
  flex-shrink: 0;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}
.send-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.2);
}
.send-button span {
  display: block ;
}
.send-button-icon-wrapper {
  margin-left: 0;
}
.send-button-icon {
  width: 1.5rem;
  height: 1.5rem;
  transform: rotate(45deg);
  margin-left: 2px;
  margin-top: -2px;
}

/* --- Empty State --- */
.empty-chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  align-items: center;
  justify-content: center;
  text-align: center;
}
.empty-chat-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--text-primary);
  background: var(--interactive-gradient);
  height: 5rem;
  width: 5rem;
}
.empty-chat-icon {
  width: 2.5rem;
  height: 2.5rem;
}
.empty-chat-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-top: 1.5rem;
  color: var(--text-primary);
}
.empty-chat-text {
  color: var(--text-muted);
  margin-top: 0.5rem;
}

/* --- Custom Scrollbar --- */
.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15); /* Was 0.1 */
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.25); /* Was 0.2 */
}
</style>
