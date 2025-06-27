import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  // State (reactive data)
  const user = ref(null)
  const isAuthenticated = ref(false)
  const loading = ref(false)
  const error = ref('')

  // Actions (functions to modify state)
  const login = async (credentials) => {
    loading.value = true
    error.value = ''
    
    try {
      const response = await fetch('http://localhost:8081/signin', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: credentials.email.trim(),
          password: credentials.password
        })
      })
      
      const data = await response.json()
      
      if (response.ok) {
        // Login successful - store user data
        user.value = data.user || { email: credentials.email }
        isAuthenticated.value = true
        error.value = ''
        return { success: true }
      } else {
        // Login failed
        error.value = data.error || "Erreur lors de la connexion"
        return { success: false, error: error.value }
      }
    } catch (err) {
      error.value = "Erreur réseau ou serveur"
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  const register = async (userData) => {
    loading.value = true
    error.value = ''
    
    try {
      // Check if userData is FormData or regular object
      const isFormData = userData instanceof FormData
      
      console.log('Registering user...', { isFormData })
      if (isFormData) {
        console.log('FormData contents:')
        for (let pair of userData.entries()) {
          console.log(pair[0], pair[1])
        }
      }
      
      const response = await fetch('http://localhost:8081/register', {
        method: 'POST',
        credentials: 'include',
        headers: isFormData ? {} : { 'Content-Type': 'application/json' },
        body: isFormData ? userData : JSON.stringify(userData)
      })
      
      console.log('Response status:', response.status)
      console.log('Response ok:', response.ok)
      
      const data = await response.json()
      console.log('Response data:', data)
      
      if (response.ok) {
        // Registration successful
        error.value = ''
        return { success: true, message: "Compte créé avec succès! Redirection vers la connexion..." }
      } else {
        error.value = data.error || "Erreur lors de l'inscription"
        return { success: false, error: error.value }
      }
    } catch (err) {
      console.error('Registration error details:', err)
      error.value = "Erreur réseau ou serveur: " + err.message
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    loading.value = true
    
    try {
      await fetch('http://localhost:8081/logout', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' }
      })
    } catch (err) {
      console.error('Logout request failed:', err)
    } finally {
      // Always clear user data regardless of request success
      user.value = null
      isAuthenticated.value = false
      error.value = ''
      loading.value = false
    }
  }

  const checkAuth = async () => {
    loading.value = true
    
    try {
      const response = await fetch('http://localhost:8081/sessionActive', {
        method: 'GET',
        credentials: 'include'
      })
      
      const data = await response.json()
      
      if (response.ok && data.type === 'Success') {
        isAuthenticated.value = true
        user.value = data.user || { email: 'user@example.com' }
        return true
      } else {
        isAuthenticated.value = false
        user.value = null
        return false
      }
    } catch (err) {
      isAuthenticated.value = false
      user.value = null
      return false
    } finally {
      loading.value = false
    }
  }

  const clearError = () => {
    error.value = ''
  }

  const setError = (message) => {
    error.value = message
  }

  // Return state and actions
  return {
    // State
    user,
    isAuthenticated,
    loading,
    error,
    
    // Actions
    login,
    register,
    logout,
    checkAuth,
    clearError,
    setError
  }
})
