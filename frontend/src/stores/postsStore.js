import { defineStore } from 'pinia'

export const usePostsStore = defineStore('posts', {
  state: () => ({
    posts: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetchAllPosts() {
      this.loading = true
      this.error = null
      try {
        const res = await fetch('http://localhost:8081/allPosts', { credentials: 'include' })
        const data = await res.json()
        if (res.ok && data.posts) {
          this.posts = data.posts
        } else {
          this.error = data.message || 'Erreur lors du chargement des posts'
        }
      } catch (e) {
        this.error = 'Erreur réseau'
      } finally {
        this.loading = false
      }
    },
    async createPost(postData) {
      this.loading = true
      this.error = null
      try {
        const res = await fetch('http://localhost:8081/createPost', {
          method: 'POST',
          credentials: 'include',
          body: postData // FormData pour image
        })
        const data = await res.json()
        if (res.ok && data.type === 'Success') {
          await this.fetchAllPosts()
        } else {
          this.error = data.message || 'Erreur lors de la création du post'
        }
      } catch (e) {
        this.error = 'Erreur réseau'
      } finally {
        this.loading = false
      }
    },
    // Ajoute d'autres actions si besoin (fetchUserPosts, addComment, etc.)
  },
})
