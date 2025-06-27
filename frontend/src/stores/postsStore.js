import { defineStore } from 'pinia';

export const useMainStore = defineStore('main', {
  state: () => ({
    count: 0,
    posts: [],
  }),
  actions: {
    increment() {
      this.count++;
    },
    async fetchPosts() {
      const res = await fetch('http://localhost:8081/allPosts', { credentials: 'include' });
      if (res.ok) {
        const data = await res.json();
        // If backend returns { type: 'Success', posts: [...] }, extract posts
        this.posts = Array.isArray(data.posts) ? data.posts : [];
        console.log('Posts data:', this.posts);
      }
    },
    async submitComment(postId, comment) {
      let commentData = new FormData();
      commentData.set('postid', postId);
      commentData.set('body', comment.body);
      commentData.set('image', comment.image);
      await fetch('http://localhost:8081/newComment', {
        method: 'POST',
        credentials: 'include',
        body: commentData
      });
      await this.fetchPosts();
    }
  }
});
