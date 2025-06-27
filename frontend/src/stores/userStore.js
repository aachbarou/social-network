import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    isAuthenticated: false,
  }),
  actions: {
    setUser(user) {
      this.user = user
      this.isAuthenticated = !!user
    },
    logout() {
      this.user = null
      this.isAuthenticated = false
    },
    async fetchUserData() {
      if (!this.user || !this.user.id) return
      const res = await fetch(`http://localhost:8081/userData?userId=${this.user.id}`, { credentials: 'include' })
      const data = await res.json()
      if (res.ok && data.users && data.users.length > 0) {
        this.user = data.users[0]
      }
    },
    async reloadUserAfterRefresh() {
      // Récupère l'utilisateur minimal puis le profil complet
      const res = await fetch('http://localhost:8081/currentUser', { credentials: 'include' })
      const data = await res.json()
      if (res.ok && data.users && data.users.length > 0) {
        this.user = data.users[0]
        this.isAuthenticated = true
        await this.fetchUserData()
      } else {
        this.user = null
        this.isAuthenticated = false
      }
    }
  },
})
