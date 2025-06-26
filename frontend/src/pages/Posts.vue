<template>
  <div v-if="isAuthenticated">
    <h1>Publications</h1>
    <p>Liste des publications à venir.</p>
  </div>
  <div v-else>
    <p>Redirecting to login...</p>
  </div>
</template>

<script>
export default {
  name: 'Posts',
  data() {
    return {
      isAuthenticated: false,
    };
  },
  created() {
    this.checkAuth();
  },
  methods: {
    checkAuth() {
      const token = localStorage.getItem('token');
      console.log("lo",token)
      if (token==="undefined") {
      console.log("jii",token)

        // Redirect to login if no token exists
        this.$router.push({ name: 'Login' });
      } else {
        this.isAuthenticated = true;
        
        // Optional: Validate token with backend
        // this.validateToken(token);
      }
    },
    // Optional: Add token validation with backend
    async validateToken(token) {
      try {
        const res = await fetch('http://localhost:8081/validate-token', {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!res.ok) throw new Error('Invalid token');
      } catch (err) {
        localStorage.removeItem('token');
        this.$router.push({ name: 'Login' });
      }
    }
  }
};
</script>

<style scoped>
</style>