<template>
  <div class="app-container">
    <!-- Section 1: Users List -->
    <div class="users-section">
      <h2 class="section-title">Users</h2>
      <ul class="user-list">
        <li v-for="user in users" :key="user.id" class="user-item">
          <span :class="['status-dot', { online: user.online }]"></span>
          <span class="user-name">{{ user.name }}</span>
        </li>
      </ul>
    </div>

    <!-- Section 2: Chat Box -->
    <div class="chat-section">
      <h2 class="section-title">Chat Box</h2>
      <div class="chat-messages" ref="chatMessagesContainer">
        <div v-for="(msg, idx) in messages" :key="idx" :class="['chat-message', { 'my-message': msg.user === 'You' }]">
          <span class="user">{{ msg.user }}:</span>
          <span class="text">{{ msg.text }}</span>
        </div>
      </div>
      <form @submit.prevent="sendMessage" class="chat-input">
        <input
          v-model="newMessage"
          type="text"
          placeholder="Type your message..."
          autocomplete="off"
        />
        <button type="submit">Send</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChatPage',
  data() {
    return {
      // Data for the Users section
      users: [
        { id: 1, name: 'Alice', online: true },
        { id: 2, name: 'Bob', online: false },
        { id: 3, name: 'Charlie', online: true },
        { id: 4, name: 'Diana', online: true },
        { id: 5, name: 'Edward', online: false },
        { id: 6, name: 'Fiona', online: true },
      ],
      // Data for the Chat Box section
      messages: [
        { user: 'Alice', text: 'Hello everyone! Ready to discuss the project?' },
        { user: 'Charlie', text: 'Hey Alice! I am. Just waiting for the others.' }
      ],
      newMessage: ''
    };
  },
  methods: {
    sendMessage() {
      if (this.newMessage.trim() !== '') {
        this.messages.push({
          user: 'You', // Assuming the current user is 'You'
          text: this.newMessage.trim()
        });
        this.newMessage = '';
        // Auto-scroll to the latest message
        this.$nextTick(() => {
          const container = this.$refs.chatMessagesContainer;
          if (container) {
            container.scrollTop = container.scrollHeight;
          }
        });
      }
    }
  }
};
</script>

<style scoped>
/* To make this component truly full-page, ensure the parent elements 
  (like html, body, and your Vue app's root div) have height: 100% and margin: 0.
*/
.app-container {
  display: flex;
  gap: 20px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  width: 100%;
  height: 100vh; /* Use full viewport height */
  padding: 20px;
  box-sizing: border-box; /* Include padding in total dimensions */
  background-color: #f0f2f5;
}

/* --- Section Base Styles --- */
.users-section, .chat-section {
  border: 1px solid #d1d5db;
  border-radius: 12px;
  padding: 20px;
  background: #ffffff;
  box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden; /* Prevent content from spilling out */
}

.section-title {
  margin-top: 0;
  margin-bottom: 15px;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 10px;
  color: #111827;
  font-size: 1.25rem;
}

/* --- Users Section Styles --- */
.users-section {
  flex: 0 0 250px; /* Base width 250px, doesn't grow */
}

.user-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto; /* Allow scrolling if user list is long */
}

.user-item {
  display: flex;
  align-items: center;
  padding: 10px 5px;
  border-bottom: 1px solid #f3f4f6;
  transition: background-color 0.2s;
}
.user-item:hover {
  background-color: #f9fafb;
}
.user-item:last-child {
  border-bottom: none;
}

.status-dot {
  height: 9px;
  width: 9px;
  border-radius: 50%;
  margin-right: 12px;
  background-color: #9ca3af; /* Offline by default */
}
.status-dot.online {
  background-color: #22c55e; /* Green for online */
}

.user-name {
  color: #374151;
  font-weight: 500;
}

/* --- Chat Section Styles --- */
.chat-section {
  flex: 1; /* Takes up the remaining horizontal space */
}

.chat-messages {
  flex-grow: 1; /* This is key: it makes the chat area fill available vertical space */
  overflow-y: auto; /* Allows scrolling for messages */
  margin-bottom: 15px;
  padding: 10px;
}

.chat-message {
  margin-bottom: 12px;
  max-width: 80%;
  word-wrap: break-word;
}
.chat-message .user {
  font-weight: 600;
  margin-right: 8px;
  display: block;
  margin-bottom: 4px;
  color: #4b5563;
}
.chat-message .text {
  background-color: #e5e7eb;
  padding: 8px 12px;
  border-radius: 18px;
  display: inline-block;
}

/* Style for messages sent by 'You' */
.chat-message.my-message {
  align-self: flex-end;
  text-align: right;
}
.chat-message.my-message .text {
  background-color: #3b82f6;
  color: #ffffff;
}

.chat-input {
  display: flex;
  gap: 10px;
}
.chat-input input {
  flex: 1;
  padding: 10px;
  border-radius: 20px;
  border: 1px solid #d1d5db;
  font-size: 1rem;
}
.chat-input input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px #bfdbfe;
}
.chat-input button {
  padding: 10px 20px;
  border-radius: 20px;
  border: none;
  background: #3b82f6;
  color: #fff;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s;
}
.chat-input button:hover {
  background: #2563eb;
}
</style>
