<template>
  <div class="login-screen">
    <div class="login-panel">
      <div class="login-brand">
        <div class="brand-badge">A</div>
        <div>
          <p class="brand-label">AsetKu</p>
          <p class="brand-subtitle">Sistem manajemen aset perusahaan</p>
        </div>
      </div>

      <h1>Selamat datang</h1>
      <p>Masuk untuk melihat dashboard aset dan pekerjaan Anda.</p>

      <form @submit.prevent="login" class="login-form">
        <label>
          <span>Username</span>
          <input 
            type="text" 
            v-model="username" 
            placeholder="Masukkan username"
            :disabled="isLoading"
          />
        </label>
        <label>
          <span>Password</span>
          <input 
            type="password" 
            v-model="password" 
            placeholder="Masukkan password"
            :disabled="isLoading"
          />
        </label>

        <button type="submit" :disabled="isLoading">
          {{ isLoading ? 'Sedang masuk...' : 'Masuk' }}
        </button>

        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const username = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)
const router = useRouter()

async function login() {
  error.value = ''
  isLoading.value = true

  if (!username.value || !password.value) {
    error.value = 'Username dan password wajib diisi.'
    isLoading.value = false
    return
  }

  try {
    const res = await api.post('/auth/login', {
      username: username.value,
      password: password.value,
    })

    if (!res.data.status) {
      error.value = res.data.message || 'Login gagal'
      return
    }

    const token = res.data.data?.token
    if (!token) {
      error.value = 'Token tidak diterima dari server'
      return
    }

    // Store token in localStorage
    localStorage.setItem('token', token)
    localStorage.setItem('username', username.value)

    // Redirect to dashboard
    router.push('/dashboard')
  } catch (err) {
    console.error('Login error:', err)
    
    if (err.response?.status === 401) {
      error.value = 'Username atau password salah'
    } else if (err.response?.data?.message) {
      error.value = err.response.data.message
    } else if (err.code === 'ERR_NETWORK') {
      error.value = 'Tidak dapat terhubung ke server. Pastikan backend berjalan.'
    } else if (err.message === 'timeout of 5000ms exceeded') {
      error.value = 'Koneksi timeout. Server tidak merespons.'
    } else {
      error.value = err.message || 'Login gagal. Silakan coba lagi.'
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-screen {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  background: linear-gradient(180deg, #eef2ff 0%, #f8fafc 100%);
}

.login-panel {
  width: min(420px, 100%);
  background: #ffffff;
  border-radius: 28px;
  box-shadow: 0 24px 80px rgba(15, 23, 42, 0.12);
  padding: 32px;
  border: 1px solid #e5e7eb;
}

.login-brand {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 28px;
}

.brand-badge {
  width: 52px;
  height: 52px;
  border-radius: 18px;
  background: #2563eb;
  color: #ffffff;
  display: grid;
  place-items: center;
  font-weight: 700;
  font-size: 1.25rem;
}

.brand-label {
  margin: 0;
  font-weight: 700;
  font-size: 1rem;
}

.brand-subtitle {
  margin: 4px 0 0;
  color: #6b7280;
  font-size: 0.95rem;
}

h1 {
  margin: 0 0 10px;
  font-size: 2rem;
  line-height: 1.1;
}

p {
  margin: 0 0 24px;
  color: #4b5563;
}

.login-form {
  display: grid;
  gap: 18px;
}

.login-form label {
  display: grid;
  gap: 10px;
  font-weight: 600;
  color: #111827;
}

.login-form input {
  width: 100%;
  border: 1px solid #d1d5db;
  border-radius: 14px;
  padding: 14px 16px;
  font-size: 1rem;
  background: #f8fafc;
  outline: none;
}

.login-form input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.12);
}

/* Button styles */
button {
  display: inline-block;
  width: 100%;
  padding: 14px 0;
  border-radius: 14px;
  font-size: 1rem;
  font-weight: 700;
  color: #ffffff;
  background: #2563eb;
  border: none;
  cursor: pointer;
  transition: background 0.2s ease;
}

button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

button:not(:disabled):hover {
  background: #1d4ed8;
}

.error-message {
  margin: 0;
  color: #dc2626;
  font-size: 0.95rem;
}
</style>