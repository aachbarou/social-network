<template>
  <div class="register-bg">
    <div class="register-card">
      <h2>Créer un compte</h2>
      <form class="register-form" @submit.prevent="handleRegister">
        <label for="email">Email</label>
        <input id="email" v-model="form.email" type="email" placeholder="Votre email" required />
        <label for="password">Mot de passe</label>
        <input id="password" v-model="form.password" type="password" placeholder="Votre mot de passe" required />
        <label for="firstname">Prénom</label>
        <input id="firstname" v-model="form.firstname" type="text" placeholder="Votre prénom" required />
        <label for="lastname">Nom</label>
        <input id="lastname" v-model="form.lastname" type="text" placeholder="Votre nom" required />
        <label for="dob">Date de naissance</label>
        <input id="dob" v-model="form.dateofbirth" type="date" required />
        <label for="avatar">Avatar (optionnel)</label>
        <input id="avatar" type="file" @change="onFileChange" />
        <label for="nickname">Pseudo (optionnel)</label>
        <input id="nickname" v-model="form.nickname" type="text" placeholder="Votre pseudo" />
        <label for="about">À propos de moi (optionnel)</label>
        <textarea id="about" v-model="form.aboutme" placeholder="Quelques mots sur vous..."></textarea>
        <button type="submit">S'inscrire</button>
      </form>
      <div v-if="successMsg" class="success-msg">{{ successMsg }}</div>
      <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
      <div class="login-link">
        Déjà un compte ? <router-link to="/login">Se connecter</router-link>
      </div>
    </div>
  </div>
</template>


<script>
export default {
  name: "Register",
  data() {
    return {
      form: {
        email: '',
        password: '',
        firstname: '',
        lastname: '',
        dateofbirth: '',
        nickname: '',
        aboutme: '',
        avatar: null
      },
      successMsg: '',
      errorMsg: ''
    }
  },
  methods: {
    check(txt) {
alert(txt)
    },
    onFileChange(e) {
      this.form.avatar = e.target.files[0];
    },
    async handleRegister() {
      this.successMsg = '';
      this.errorMsg = '';
      try {
        const formData = new FormData();
        for (const key in this.form) {
          if (this.form[key]) {
            formData.append(key, this.form[key]);
          }
        }
        const res = await fetch('http://localhost:8081/register', {
          method: 'POST',
          body: formData
        });
        if (res.ok) {
          this.successMsg = 'Inscription réussie ! Vous pouvez vous connecter.';
          this.form = { email: '', password: '', firstname: '', lastname: '', dateofbirth: '', nickname: '', aboutme: '', avatar: null };
        } else {
          const data = await res.json();
          this.errorMsg = data.error || "Erreur lors de l'inscription.";
        }
      } catch (err) {
        this.errorMsg = "Erreur réseau ou serveur.";
      }
    }
  }
}
</script>

<style scoped>
.register-bg {
  min-height: 100vh;
  background: #f6f8fa;
  display: flex;
  align-items: center;
  justify-content: center;
}
.register-card {
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(25, 118, 210, 0.07);
  padding: 32px 28px 18px 28px;
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
h2 {
  margin-bottom: 18px;
  color: #1976d2;
  font-weight: 700;
  font-size: 1.25rem;
  text-align: center;
}
.register-form {
  width: 100%;
  display: flex;
  flex-direction: column;
}
.register-form label {
  margin-bottom: 5px;
  color: #333;
  font-size: 1rem;
  font-weight: 500;
}
.register-form input,
.register-form textarea {
  padding: 10px 12px;
  margin-bottom: 13px;
  border-radius: 6px;
  border: 1.2px solid #cfd8dc;
  background: #f7fbff;
  font-size: 1rem;
  color: #222;
  transition: border 0.2s, background 0.2s;
  outline: none;
}
.register-form input:focus,
.register-form textarea:focus {
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
.login-link {
  margin-top: 15px;
  text-align: center;
  font-size: 1rem;
}
.login-link a {
  color: #1976d2;
  text-decoration: underline;
  font-weight: 500;
}
.success-msg {
  color: #388e3c;
  background: #e8f5e9;
  border: 1px solid #c8e6c9;
  padding: 8px 12px;
  border-radius: 6px;
  margin: 10px 0 0 0;
  text-align: center;
}
.error-msg {
  color: #d32f2f;
  background: #ffebee;
  border: 1px solid #ffcdd2;
  padding: 8px 12px;
  border-radius: 6px;
  margin: 10px 0 0 0;
  text-align: center;
}
</style>