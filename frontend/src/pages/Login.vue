<template>
  <div class="login-bg">
    <div class="login-card">
      <div class="login-logo">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none"><circle cx="24" cy="24" r="24" fill="#1976d2"/><text x="50%" y="55%" text-anchor="middle" fill="#fff" font-size="22" font-family="Arial" dy=".3em">SN</text></svg>
      </div>
      <h2>Connexion à Social Network</h2>
      <form class="login-form">
        <label for="email">Email</label>
        <input id="email" type="email" placeholder="Votre email" required />
        <label for="password">Mot de passe</label>
        <input id="password" type="password" placeholder="Votre mot de passe" required />
        <button type="submit">Se connecter</button>
      </form>
      <div class="register-link">
        Pas de compte ? <router-link to="/register">Créer un compte</router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Login"
}
</script>

<style scoped>
.login-bg {
  min-height: 100vh;
  background: #f6f8fa;
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-card {
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(25, 118, 210, 0.07);
  padding: 36px 30px 24px 30px;
  width: 100%;
  max-width: 350px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.login-logo {
  margin-bottom: 16px;
}
h2 {
  margin-bottom: 20px;
  color: #1976d2;
  font-weight: 700;
  font-size: 1.25rem;
  text-align: center;
}
.login-form {
  width: 100%;
  display: flex;
  flex-direction: column;
}
.login-form label {
  margin-bottom: 5px;
  color: #333;
  font-size: 1rem;
  font-weight: 500;
}
.login-form input {
  padding: 10px 12px;
  margin-bottom: 15px;
  border-radius: 6px;
  border: 1.2px solid #cfd8dc;
  background: #f7fbff;
  font-size: 1rem;
  color: #222;
  transition: border 0.2s, background 0.2s;
  outline: none;
}
.login-form input:focus {
  border: 1.5px solid #1976d2;
  background: #e3f0ff;
}
button {
  width: 100%;
  padding: 11px;
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-weight: bold;
  font-size: 1.08rem;
  cursor: pointer;
  margin-top: 6px;
  box-shadow: 0 1px 4px rgba(25, 118, 210, 0.06);
  transition: background 0.2s;
}
button:hover {
  background: #1251a3;
}
.register-link {
  margin-top: 15px;
  text-align: center;
  font-size: 1rem;
}
.register-link a {
  color: #1976d2;
  text-decoration: underline;
  font-weight: 500;
}
</style>
