<template>
  <div class="chat-container">
    <div class="chat-layout">
      <div class="friends-list-section">
        <div class="header">
          <div class="header-icon-wrapper">
            <svg class="header-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path>
            </svg>
          </div>
          <div class="header-title">ChatApp</div>
        </div>

        <div v-if="currentUser" class="profile-card">
          <div class="profile-avatar-wrapper">
            <img :src="currentUser.avatar || 'imageUpload/default.svg'" alt="Avatar" class="profile-avatar">
          </div>
          <div class="profile-name">{{ currentUser.nickname || `${currentUser.firstName} ${currentUser.lastName}` }}</div>
          <div class="profile-email">{{ currentUser.email }}</div>
        </div>

        <div class="friends-list-container">
          <div class="friends-list-header">
            <span class="friends-list-title">Following</span>
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
                <img :src="friend.avatar || 'imageUpload/default.svg'" :alt="friend.nickname" class="friend-avatar">
              </div>
              <div class="friend-name">{{ friend.nickname || friend.firstName }}</div>
            </button>
          </div>
        </div>
      </div>

      <div class="chat-area-wrapper">
        <div class="chat-box">
          <div v-if="selectedUser" class="active-chat">
            <div class="chat-header">
              <div class="chat-header-user">
                <img :src="selectedUser.avatar || 'imageUpload/default.svg'" alt="Avatar" class="chat-header-avatar">
                <div class="chat-header-text">
                  <div class="chat-header-name-wrapper">
                    <span class="chat-header-name">{{ selectedUser.nickname || selectedUser.firstName }}</span>
                  </div>
                  <span class="chat-header-email">{{ selectedUser.email }}</span>
                </div>
              </div>
               <button class="reset-chat-btn" @click="resetChat">Reset Chat</button>
            </div>
            
            <div ref="messageContainer" class="messages-container custom-scrollbar">
              <div class="messages-list">
                <div class="messages-grid">
                  <template v-for="(message, index) in messages[selectedUser.id]" :key="index">
                    <div v-if="message.sender !== currentUser.id" class="message-wrapper received">
                      <div class="message-layout">
                        <div class="message-avatar-wrapper">
                          <img :src="selectedUser.avatar || 'imageUpload/default.svg'" alt="Avatar" class="message-avatar">
                        </div>
                        <div class="message-bubble received">
                          <div>{{ message.text }}</div>
                        </div>
                      </div>
                    </div>
                    <div v-else class="message-wrapper sent">
                      <div class="message-layout sent">
                        <div class="message-avatar-wrapper">
                          <img :src="currentUser.avatar || 'imageUpload/default.svg'" alt="Avatar" class="message-avatar">
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
const messageContainer = ref(null); // To control scroll position
const messages = reactive({}); // Object to hold message history for each user ID

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
  // Ensure the message array exists for the selected user
  if (!messages[user.id]) {
    messages[user.id] = [];
  }
  scrollToBottom();
};

/**
 * Adds the new message to the conversation and simulates a reply.
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

  if (currentUser.value && currentUser.value.id) {
    await followStore.fetchFollowing(currentUser.value.id);

    if (following.value.length > 0) {
      selectUser(following.value[0]);
    }
  }
});
</script>

<style scoped>
        .chat-container {
            display: flex;
            height: 100vh;
        }

        .chat-layout {
            display: flex;
            flex-direction: row;
            height: 100%;
            width: 100%;
            overflow-x: hidden;
            background-color: #f3f4f6;
        }

        .friends-list-section {
            display: flex;
            flex-direction: column;
            padding: 2rem 0.5rem 2rem 1.5rem;
            width: 16rem;
            background-color: #ffffff;
            flex-shrink: 0;
        }
        @media (min-width: 768px) {
            .friends-list-section {
                width: 20rem;
            }
        }

        .header {
            display: flex;
            align-items: center;
            justify-content: center;
            height: 3rem;
            width: 100%;
        }

        .header-icon-wrapper {
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 1rem;
            color: #4338ca;
            background-color: #e0e7ff;
            height: 2.5rem;
            width: 2.5rem;
        }
        .header-icon {
            width: 1.5rem;
            height: 1.5rem;
        }

        .header-title {
            margin-left: 0.5rem;
            font-weight: 700;
            font-size: 1.5rem;
        }

        .profile-card {
            display: flex;
            flex-direction: column;
            align-items: center;
            background-color: #e0e7ff;
            border: 1px solid #e5e7eb;
            margin-top: 1rem;
            width: 100%;
            padding: 1.5rem 1rem;
            border-radius: 0.5rem;
        }
        .profile-avatar-wrapper {
            height: 5rem;
            width: 5rem;
            border-radius: 9999px;
            border: 1px solid #d1d5db;
            overflow: hidden;
        }
        .profile-avatar {
            height: 100%;
            width: 100%;
            object-fit: cover;
        }
        .profile-name {
            font-size: 0.875rem;
            font-weight: 600;
            margin-top: 0.5rem;
        }
        .profile-email {
            font-size: 0.75rem;
            color: #6b7280;
        }

        .friends-list-container {
            display: flex;
            flex-direction: column;
            margin-top: 2rem;
            flex-grow: 1;
            min-height: 0;
        }
        .friends-list-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 0.75rem;
        }
        .friends-list-title {
            font-weight: 700;
        }
        .friends-list-count {
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: #d1d5db;
            height: 1rem;
            width: 1rem;
            border-radius: 9999px;
            font-size: 0.75rem;
        }
        .friends-list {
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
            margin-top: 1rem;
            margin-left: -0.5rem;
            margin-right: -0.5rem;
            flex-grow: 1;
            overflow-y: auto;
            padding-right: 0.5rem;
        }
        .friend-item {
            display: flex;
            align-items: center;
            border-radius: 0.75rem;
            padding: 0.5rem;
            width: 100%;
            text-align: left;
            border: none;
            background: none;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        .friend-item:hover {
            background-color: #f3f4f6;
        }
        .friend-item.active {
            background-color: #e0e7ff;
        }
        .friend-avatar-wrapper {
            display: flex;
            align-items: center;
            justify-content: center;
            height: 2rem;
            width: 2rem;
            background-color: #c7d2fe;
            border-radius: 9999px;
            flex-shrink: 0;
        }
        .friend-avatar {
            height: 100%;
            width: 100%;
            border-radius: 9999px;
            object-fit: cover;
        }
        .friend-name {
            margin-left: 0.5rem;
            font-size: 0.875rem;
            font-weight: 600;
        }

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
            flex-shrink: 0;
            border-radius: 1rem;
            background-color: #ffffff;
            height: 100%;
            padding: 1rem;
            min-height: 0;
        }

        .active-chat {
            display: flex;
            flex-direction: column;
            height: 100%;
        }
        .chat-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding-bottom: 0.75rem;
            border-bottom: 2px solid #e5e7eb;
            flex-shrink: 0;
        }
        .chat-header-user {
            display: flex;
            align-items: center;
            gap: 1rem;
        }
        .chat-header-avatar {
            width: 2.5rem;
            height: 2.5rem;
            border-radius: 9999px;
            object-fit: cover;
        }
        @media (min-width: 640px) {
            .chat-header-avatar {
                width: 3rem;
                height: 3rem;
            }
        }
        .chat-header-text {
            display: flex;
            flex-direction: column;
            line-height: 1.25;
        }
        .chat-header-name-wrapper {
            display: flex;
            align-items: center;
            margin-top: 0.25rem;
            font-size: 1.125rem;
        }
        .chat-header-name {
            color: #374151;
            margin-right: 0.75rem;
        }
        .chat-header-email {
            font-size: 0.875rem;
            color: #4b5563;
        }
        .reset-chat-btn {
            padding: 0.375rem 1rem;
            border-radius: 1rem;
            border: none;
            background-color: #ef4444;
            color: #ffffff;
            font-weight: 600;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        .reset-chat-btn:hover {
            background-color: #b91c1c;
        }

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
            gap: 0.5rem 0;
            padding: 0.75rem 0;
        }
        .message-wrapper {
            padding: 0.75rem;
            border-radius: 0.5rem;
        }
        .message-wrapper.received {
            grid-column: 1 / span 10;
        }
        .message-wrapper.sent {
            grid-column: 4 / span 10;
        }
        @media (min-width: 768px) {
            .message-wrapper.received {
                grid-column: 1 / span 7;
            }
            .message-wrapper.sent {
                grid-column: 7 / span 7;
            }
        }
        .message-layout {
            display: flex;
            align-items: center;
            flex-direction: row;
        }
        .message-layout.sent {
            justify-content: flex-start;
            flex-direction: row-reverse;
        }
        .message-avatar-wrapper {
            display: flex;
            align-items: center;
            justify-content: center;
            height: 2.5rem;
            width: 2.5rem;
            border-radius: 9999px;
            background-color: #4f46e5;
            flex-shrink: 0;
        }
        .message-avatar {
            width: 100%;
            height: 100%;
            border-radius: 9999px;
            object-fit: cover;
        }
        .message-bubble {
            position: relative;
            font-size: 0.875rem;
            padding: 0.5rem 1rem;
            box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1);
            border-radius: 0.75rem;
        }
        .message-bubble.received {
            margin-left: 0.75rem;
            background-color: #f3f4f6;
        }
        .message-bubble.sent {
            margin-right: 0.75rem;
            background-color: #e0e7ff;
        }

        .message-input-section {
            display: flex;
            align-items: center;
            height: 4rem;
            border-radius: 0.75rem;
            background-color: #ffffff;
            width: 100%;
            padding: 0 1rem;
            flex-shrink: 0;
        }
        .input-wrapper {
            flex-grow: 1;
        }
        .input-relative-wrapper {
            position: relative;
            width: 100%;
        }
        .message-input {
            width: 100%;
            border: 1px solid #d1d5db;
            border-radius: 0.75rem;
            padding-left: 1rem;
            height: 2.5rem;
        }
        .message-input:focus {
            outline: none;
            border-color: #a5b4fc;
        }
        .send-button-wrapper {
            margin-left: 1rem;
        }
        .send-button {
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: #4f46e5;
            border-radius: 0.75rem;
            color: #ffffff;
            padding: 0.5rem 1rem;
            flex-shrink: 0;
            border: none;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        .send-button:hover {
            background-color: #4338ca;
        }
        .send-button-icon-wrapper {
            margin-left: 0.5rem;
        }
        .send-button-icon {
            width: 1rem;
            height: 1rem;
            transform: rotate(45deg);
            margin-top: -1px;
        }

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
            border-radius: 1rem;
            color: #4338ca;
            background-color: #e0e7ff;
            height: 4rem;
            width: 4rem;
        }
        .empty-chat-icon {
            width: 2.5rem;
            height: 2.5rem;
        }
        .empty-chat-title {
            font-size: 1.5rem;
            font-weight: 600;
            margin-top: 1rem;
        }
        .empty-chat-text {
            color: #6b7280;
            margin-top: 0.5rem;
        }
        
        .custom-scrollbar::-webkit-scrollbar { width: 6px; }
        .custom-scrollbar::-webkit-scrollbar-track { background: #f1f5f9; border-radius: 10px; }
        .custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>
