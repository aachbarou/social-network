<template>
  <div class="chat-container">
    <div class="chat-layout">
      <div class="friends-list-section">
        <div class="header">
          <div class="header-icon-wrapper">
            <svg class="header-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path></svg>
          </div>
          <div class="header-title">ChatApp</div>
        </div>

        <div v-if="currentUser" class="profile-card">
          <div class="profile-avatar-wrapper">
            <img :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)" alt="Avatar" class="profile-avatar" />
          </div>
          <div class="profile-name">{{ currentUser.nickname || `${currentUser.firstName} ${currentUser.lastName}` }}</div>
          <div class="profile-email">{{ currentUser.email }}</div>
        </div>

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
                 <img :src="getFullImageUrl(friend.avatar || friend.imagePath)" alt="Avatar" class="friend-avatar" />
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
                 <img :src="getFullImageUrl(selectedUser.avatar || selectedUser.imagePath)" alt="Avatar" class="chat-header-avatar" />
                <div class="chat-header-text">
                  <span class="chat-header-name">{{ selectedUser.nickname || selectedUser.firstName }}</span>
                  <span class="chat-header-email">{{ selectedUser.email }}</span>
                </div>
              </div>
            </div>
            
            <div ref="messageContainer" class="messages-container custom-scrollbar">
              <div class="messages-list">
                <div class="messages-grid">
                  <template v-for="(message, index) in chatMessages" :key="message.id || index">
                    <div v-if="message.senderId !== currentUser.id" class="message-wrapper received">
                      <div class="message-layout">
                        <div class="message-avatar-wrapper">
                          <img :src="getFullImageUrl(selectedUser.avatar || selectedUser.imagePath)" alt="Avatar" class="message-avatar" />
                        </div>
                        <div class="message-bubble received"><div>{{ message.content || message.text }}</div></div>
                      </div>
                    </div>
                    <div v-else class="message-wrapper sent">
                      <div class="message-layout sent">
                        <div class="message-avatar-wrapper">
                          <img :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)" alt="Avatar" class="message-avatar" />
                        </div>
                        <div class="message-bubble sent"><div>{{ message.content || message.text }}</div></div>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </div>

            <div class="message-input-section">
                <div class="emojis-container">
                    <button @click="toggleEmojis" class="emoji-toggle-btn">😊</button>
                    <div v-if="showEmojis" class="emojis-picker">
                        <span v-for="emoji in emojiList" :key="emoji" @click="addEmoji(emoji)" class="emoji-item">{{ emoji }}</span>
                    </div>
                </div>
                <input
                    type="text"
                    v-model="newMessage"
                    @keyup.enter="sendMessage"
                    class="message-input"
                    placeholder="Type your message..."
                />
                <button @click="sendMessage" class="send-button">
                    <span class="send-button-icon-wrapper">
                        <svg class="send-button-icon" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"></path></svg>
                    </span>
                </button>
            </div>
          </div>
          <div v-else class="empty-chat-container">
            <div class="empty-chat-icon-wrapper">
              <svg class="empty-chat-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
            </div>
            <h2 class="empty-chat-title">Select a conversation</h2>
            <p class="empty-chat-text">Choose a friend from the list to start chatting.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { useUserStore } from '../stores/userStore';
import { useFollowStore } from '../stores/followStore';
import { useChatStore } from '../stores/chatStore';
import { storeToRefs } from 'pinia';
import wsService from '../services/websocketService.js';

const userStore = useUserStore();
const followStore = useFollowStore();
const chatStore = useChatStore();

const { user: currentUser } = storeToRefs(userStore);
const { following } = storeToRefs(followStore);

const selectedUser = ref(null);
const newMessage = ref('');
const messageContainer = ref(null);
const showEmojis = ref(false);
const removeWsListener = ref(null);
const historicalMessages = ref([]);

const chatMessages = computed(() => {
  if (!selectedUser.value || !currentUser.value) return [];
  const newMsgs = chatStore.getMessages(currentUser.value.id)(selectedUser.value.id, 'PERSON');
  return [...historicalMessages.value, ...newMsgs];
});

const emojiList = ['😊', '😂', '❤️', '👍', '🤔', '🎉', '😢', '😠'];

function getFullImageUrl(path) {
  if (!path) return 'https://placehold.co/120x120/e879c6/white?text=User';
  if (path.startsWith('http')) return path;
  try {
    return new URL(path.replace(/^\/+/, ''), 'http://localhost:8081').href;
  } catch (e) {
    console.error("Invalid image path:", path);
    return 'https://placehold.co/120x120/e879c6/white?text=Error';
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
    }
  });
};

const fetchPreviousMessages = async (user) => {
  if (!user) {
    historicalMessages.value = [];
    return;
  }
  try {
    const response = await fetch("http://localhost:8081/messages", {
        method: "POST",
        credentials: "include",
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            receiverId: user.id,
            type: 'PERSON'
        })
    });
    const data = await response.json();
    historicalMessages.value = data.chatMessage || [];
  } catch (error) {
    console.error("Failed to fetch previous messages:", error);
    historicalMessages.value = [];
  }
};

const selectUser = async (user) => {
  selectedUser.value = user;
  
  chatStore.removeUnreadMessages({ receiverId: user.id, type: 'PERSON' });

  if (currentUser.value) {
    const currentUserId = currentUser.value.id;
    chatStore.newChatMessages = chatStore.newChatMessages.filter(msg => 
        !((msg.senderId === currentUserId && msg.receiverId === user.id) || 
          (msg.senderId === user.id && msg.receiverId === currentUserId))
    );
  }
  
  await fetchPreviousMessages(user); 
  
  scrollToBottom();
};

const handleIncomingMessage = (data) => {
    if (data && data.action === 'notification' && data.notification) {
        const message = data.notification;
        console.log("Unwrapped WebSocket message:", message);
        chatStore.addNewChatMessage(message);
        if (selectedUser.value && (message.senderId === selectedUser.value.id || message.receiverId === selectedUser.value.id)) {
            scrollToBottom();
        }
    } else {
        console.warn("Received WebSocket message in an unexpected format:", data);
    }
};

const sendMessage = async () => {
  const messageContent = newMessage.value.trim();
  if (messageContent === '' || !selectedUser.value || !currentUser.value) {
    return;
  }

  const msgObj = {
    receiverId: selectedUser.value.id,
    content: messageContent,
    type: 'PERSON',
    senderId: currentUser.value.id,
    sender: { 
        id: currentUser.value.id,
        nickname: currentUser.value.nickname 
    }
  };
  try {
    const response = await fetch("http://localhost:8081/newMessage", {
        body: JSON.stringify(msgObj),
        method: "POST",
        credentials: "include",
        headers: { 'Content-Type': 'application/json' }
    });
    const data = await response.json();

    if (data.type === "Success") {
        chatStore.addNewChatMessage({ ...msgObj, sender: { nickname: currentUser.value.nickname } });
        newMessage.value = '';
        showEmojis.value = false;
        scrollToBottom();
    } else {
        console.error("Failed to send message:", data.message);
    }
  } catch (error) {
    console.error("Error sending message:", error);
  }
  const sent = wsService.send(msgObj);

  if (sent) {
    newMessage.value = '';
    showEmojis.value = false;
    scrollToBottom(); 
  } else {
    console.error("WebSocket is not connected. Message not sent.");
    alert("You are not connected to the chat server.");
  }
};

const toggleEmojis = () => {
  showEmojis.value = !showEmojis.value;
};

const addEmoji = (emoji) => {
  newMessage.value += emoji;
};

onMounted(async () => {
  try {
    await userStore.reloadUserAfterRefresh();
    if (currentUser.value && currentUser.value.id) {
      await followStore.fetchFollowing(currentUser.value.id);
      removeWsListener.value = wsService.addListener('chat-component', handleIncomingMessage);
    }
  } catch (error) {
    console.error("Error during component mount:", error);
  }
});

onUnmounted(() => {
  if (removeWsListener.value) {
    removeWsListener.value();
  }
});

</script>

<style scoped>
.chat-container { 
  display: flex; 
  height: 100vh; 
  color: #14171a; 
  background-color: #f7f9fc; 
  font-family: 'Inter', sans-serif; 
}
.chat-layout { display: flex; flex-direction: row; height: 100%; width: 100%; overflow: hidden; }

.friends-list-section { 
  display: flex; 
  flex-direction: column; 
  padding: 1.5rem; 
  width: 22rem; 
  background: #ffffff; 
  border-right: 1px solid #e1e8ed; 
  flex-shrink: 0; 
}
.header { display: flex; align-items: center; justify-content: center; height: 3rem; width: 100%; flex-shrink: 0; margin-bottom: 1rem; }
.header-icon-wrapper { 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  border-radius: 0.75rem; 
  color: #ffffff; 
  background: #1d9bf0; 
  height: 2.5rem; 
  width: 2.5rem; 
}
.header-icon { width: 1.5rem; height: 1.5rem; }
.header-title { margin-left: 0.75rem; font-weight: 700; font-size: 1.5rem; color: #14171a; }

.profile-card { 
  display: flex; 
  flex-direction: column; 
  align-items: center; 
  background: #f7f9fc; 
  border: 1px solid #e1e8ed; 
  margin-top: 1rem; 
  width: 100%; 
  padding: 1.5rem 1rem; 
  border-radius: 1rem; 
  flex-shrink: 0; 
}
.profile-avatar-wrapper { height: 5rem; width: 5rem; border-radius: 50%; border: 2px solid #e1e8ed; overflow: hidden; }
.profile-avatar { height: 100%; width: 100%; object-fit: cover; }
.profile-name { font-size: 1rem; font-weight: 600; margin-top: 0.75rem; color: #14171a; }
.profile-email { font-size: 0.875rem; color: #657786; }

.friends-list-container { display: flex; flex-direction: column; margin-top: 1.5rem; flex-grow: 1; min-height: 0; }
.friends-list-header { display: flex; justify-content: space-between; align-items: center; font-size: 0.875rem; padding: 0 0.5rem; color: #657786; }
.friends-list-title { font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
.friends-list-count { 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background-color: #eef2f7; 
  height: 1.25rem; 
  min-width: 1.25rem; 
  padding: 0 0.35rem; 
  border-radius: 0.625rem; 
  font-size: 0.75rem; 
  font-weight: bold; 
  color: #657786; 
}
.friends-list { display: flex; flex-direction: column; gap: 0.25rem; margin-top: 0.5rem; flex-grow: 1; overflow-y: auto; padding-right: 0.5rem; }
.friend-item { display: flex; align-items: center; border-radius: 0.75rem; padding: 0.75rem; width: 100%; text-align: left; border: 1px solid transparent; background: none; cursor: pointer; transition: all 0.2s ease; position: relative; color: #14171a; }
.friend-item:hover { background-color: #f7f9fc; }
.friend-item.active { color: #ffffff; background: #1d9bf0; font-weight: 600; }
.friend-avatar-wrapper { height: 2.5rem; width: 2.5rem; border-radius: 50%; flex-shrink: 0; overflow: hidden; margin-right: 0.75rem; }
.friend-avatar { height: 100%; width: 100%; object-fit: cover; }
.friend-name { font-size: 1rem; font-weight: 500; color: inherit; }
.friend-item.active .friend-name { font-weight: 600; }

.chat-area-wrapper { display: flex; flex-direction: column; flex: 1 1 0%; height: 100%; padding: 1.5rem; }
.chat-box { display: flex; flex-direction: column; flex: 1 1 0%; border-radius: 1.5rem; background: #ffffff; border: 1px solid #e1e8ed; height: 100%; min-height: 0; }
.active-chat { display: flex; flex-direction: column; height: 100%; }
.chat-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-bottom: 1px solid #e1e8ed; flex-shrink: 0; }
.chat-header-user { display: flex; align-items: center; gap: 1rem; }
.chat-header-avatar { width: 3rem; height: 3rem; border-radius: 50%; object-fit: cover; }
.chat-header-text { display: flex; flex-direction: column; line-height: 1.25; }
.chat-header-name { font-size: 1.25rem; font-weight: 600; color: #14171a; }
.chat-header-email { font-size: 0.875rem; color: #657786; }

.messages-container { display: flex; flex-direction: column; flex-grow: 1; overflow-y: auto; padding: 1rem; min-height: 0; }
.messages-list { margin-top: auto; }
.messages-grid { display: flex; flex-direction: column; gap: 1rem; }
.message-wrapper { display: flex; max-width: 80%; }
.message-wrapper.received { align-self: flex-start; }
.message-wrapper.sent { align-self: flex-end; }
.message-layout { display: flex; align-items: flex-end; flex-direction: row; gap: 0.75rem; }
.message-layout.sent { flex-direction: row-reverse; }
.message-avatar-wrapper { height: 2.5rem; width: 2.5rem; border-radius: 50%; flex-shrink: 0; overflow: hidden; }
.message-avatar { width: 100%; height: 100%; object-fit: cover; }
.message-bubble { position: relative; font-size: 1rem; padding: 0.75rem 1rem; border-radius: 1.25rem; color: #14171a; word-break: break-word; }
.message-bubble.received { background-color: #eef2f7; border-top-left-radius: 0.25rem; }
.message-bubble.sent { background: #1d9bf0; color: #ffffff; border-top-right-radius: 0.25rem; }

.message-input-section { display: flex; align-items: center; padding: 1rem 1.5rem; border-top: 1px solid #e1e8ed; flex-shrink: 0; gap: 1rem; position: relative; background: #ffffff; }
.emojis-container { position: relative; }
.emoji-toggle-btn { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: #657786; padding: 0.5rem; }
.emojis-picker { position: absolute; bottom: 100%; left: 0; background: #ffffff; border: 1px solid #e1e8ed; border-radius: 1rem; padding: 0.5rem; display: flex; flex-wrap: wrap; gap: 0.5rem; width: 200px; box-shadow: 0 -4px 20px rgba(0,0,0,0.1); z-index: 10; }
.emoji-item { cursor: pointer; font-size: 1.25rem; transition: transform 0.2s; }
.emoji-item:hover { transform: scale(1.2); }
.message-input { flex-grow: 1; background-color: #eef2f7; border: 1px solid #e1e8ed; border-radius: 1rem; padding: 0 1rem; height: 3.5rem; color: #14171a; font-size: 1rem; transition: all 0.2s; }
.message-input:focus { outline: none; border-color: #1d9bf0; background: #ffffff; box-shadow: 0 0 0 2px rgba(29, 155, 240, 0.3); }
.send-button { display: flex; align-items: center; justify-content: center; background: #1d9bf0; border-radius: 1rem; color: #ffffff; height: 3.5rem; width: 3.5rem; flex-shrink: 0; border: none; cursor: pointer; transition: all 0.2s ease; }
.send-button:hover { opacity: 0.9; }
.send-button-icon { width: 1.5rem; height: 1.5rem; }

.empty-chat-container { display: flex; flex-direction: column; height: 100%; align-items: center; justify-content: center; text-align: center; padding: 2rem; }
.empty-chat-icon-wrapper { display: flex; align-items: center; justify-content: center; border-radius: 50%; color: #ffffff; background: #1d9bf0; height: 5rem; width: 5rem; }
.empty-chat-icon { width: 2.5rem; height: 2.5rem; }
.empty-chat-title { font-size: 1.5rem; font-weight: 600; margin-top: 1.5rem; color: #14171a; }
.empty-chat-text { color: #657786; margin-top: 0.5rem; }
.custom-scrollbar::-webkit-scrollbar { width: 8px; }
.custom-scrollbar::-webkit-scrollbar-track { background: #f7f9fc; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #d1d9e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #b0b8c1; }
</style>