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
  name: 'Posts'
}
</script>

<style scoped>
h1 {
  color: #1976d2;
}
</style>