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
      try {
        const res = await fetch('http://localhost:8081/allPosts', { credentials: 'include' });
        if (res.ok) {
          const data = await res.json();
          // If backend returns { type: 'Success', posts: [...] }, extract posts
          this.posts = Array.isArray(data.posts) ? data.posts : [];
        } else {
          console.error('Failed to fetch posts:', await res.text());
        }
      } catch (error) {
        console.error('Error fetching posts:', error);
      }
    },
    async submitComment(postId, comment) {
      let commentData = new FormData();
      commentData.set('postid', postId);
      commentData.set('body', comment.body);
      if (comment.image) {
        commentData.append('image', comment.image);
      }
      await fetch('http://localhost:8081/newComment', {
        method: 'POST',
        credentials: 'include',
        body: commentData
      });
      await this.fetchPosts();
    },
    async submitPost(post) {
      let postData = new FormData();
      postData.set('body', post.body);
      postData.set('privacy', post.privacy);
      if (post.image) {
        postData.append('image', post.image);
      }
      await fetch('http://localhost:8081/newPost', {
        method: 'POST',
        credentials: 'include',
        body: postData
      });
      await this.fetchPosts();
    }
  }
});
