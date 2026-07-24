<template>
  <div class="login-screen">
    <div class="login-panel">
      <div class="login-brand">
        <div class="brand-badge">🏢</div>
        <div>
          <p class="brand-label">AsetKu Hotel</p>
          <p class="brand-subtitle">Work Order & Asset Management System</p>
        </div>
      </div>

      <h1>Selamat Datang</h1>
      <p>Masukkan username dan password Anda dari database hotel.</p>

      <form @submit.prevent="login" class="login-form">
        <label>
          <span>Username</span>
          <input 
            type="text" 
            v-model="username" 
            placeholder="Masukkan username database"
            :disabled="isLoading"
            required
          />
        </label>

        <label>
          <span>Password</span>
          <div class="password-input-wrapper">
            <input 
              :type="showPassword ? 'text' : 'password'" 
              v-model="password" 
              placeholder="Masukkan password"
              :disabled="isLoading"
              required
            />
            <button 
              type="button" 
              class="eye-toggle-btn" 
              @click="showPassword = !showPassword"
              :title="showPassword ? 'Sembunyikan Password' : 'Tampilkan Password'"
              tabindex="-1"
            >
              {{ showPassword ? '🙈' : '👁️' }}
            </button>
          </div>
        </label>

        <button type="submit" class="submit-login-btn" :disabled="isLoading">
          {{ isLoading ? 'Memverifikasi Data...' : 'Masuk Ke Sistem' }}
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
const showPassword = ref(false)
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

    const data = res.data?.data || {}
    const token = data.token
    const role = data.role || 'external'
    const name = data.name || username.value

    if (!token) {
      throw new Error('Token tidak ditemukan dari server.')
    }

    // Save token and user details in sessionStorage (forces re-login on new sessions)
    sessionStorage.setItem('token', token)
    sessionStorage.setItem('user_role', role)
    sessionStorage.setItem('user_name', name)
    sessionStorage.setItem('username', data.username || username.value)

    // Sync to localStorage for current active session
    localStorage.setItem('token', token)
    localStorage.setItem('user_role', role)
    localStorage.setItem('user_name', name)
    localStorage.setItem('username', data.username || username.value)

    router.push('/dashboard')
  } catch (err) {
    console.error('Login error:', err)
    error.value = err.response?.data?.message || err.message || 'Username atau password salah. Silakan periksa kembali kredensial database Anda.'
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
  padding: 24px 16px;
  background: radial-gradient(circle at 50% 20%, #1e293b 0%, #090d16 100%);
  position: relative;
  overflow: hidden;
}

.login-screen::before {
  content: '';
  position: absolute;
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.15) 0%, rgba(0,0,0,0) 70%);
  top: 10%;
  left: 30%;
  border-radius: 50%;
  pointer-events: none;
}

.login-panel {
  width: min(420px, 100%);
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 28px;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.5), 0 0 1px rgba(255, 255, 255, 0.15);
  padding: 36px 32px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #ffffff;
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-brand {
  display: flex;
  gap: 14px;
  align-items: center;
  margin-bottom: 24px;
}

.brand-badge {
  width: 50px;
  height: 50px;
  border-radius: 16px;
  background: linear-gradient(135deg, #0284c7, #2563eb);
  color: #ffffff;
  display: grid;
  place-items: center;
  font-size: 1.5rem;
  box-shadow: 0 4px 14px rgba(2, 132, 199, 0.4);
}

.brand-label {
  margin: 0;
  font-weight: 800;
  font-size: 1.15rem;
  color: #ffffff;
}

.brand-subtitle {
  margin: 2px 0 0;
  color: #94a3b8;
  font-size: 0.85rem;
}

h1 {
  margin: 0 0 6px;
  font-size: 1.85rem;
  color: #ffffff;
  font-weight: 800;
}

p {
  margin: 0 0 24px;
  color: #94a3b8;
  font-size: 0.9rem;
}

.login-form {
  display: grid;
  gap: 18px;
}

.login-form label {
  display: grid;
  gap: 8px;
  font-weight: 600;
  font-size: 0.88rem;
  color: #cbd5e1;
}

.login-form input {
  width: 100%;
  border: 1px solid #334155;
  border-radius: 14px;
  padding: 13px 16px;
  font-size: 0.95rem;
  background: #1e293b;
  color: #ffffff;
  outline: none;
  transition: all 0.2s ease;
}

.login-form input:focus {
  border-color: #38bdf8;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
}

.password-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input-wrapper input {
  padding-right: 48px;
}

.eye-toggle-btn {
  position: absolute;
  right: 12px;
  background: transparent;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.eye-toggle-btn:hover {
  opacity: 1;
}

.submit-login-btn {
  display: inline-block;
  width: 100%;
  padding: 14px 0;
  border-radius: 14px;
  font-size: 1rem;
  font-weight: 700;
  color: #ffffff;
  background: linear-gradient(135deg, #0284c7, #2563eb);
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 8px;
  box-shadow: 0 4px 16px rgba(37, 99, 235, 0.35);
}

.submit-login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 22px rgba(37, 99, 235, 0.45);
}

.submit-login-btn:disabled {
  background: #475569;
  cursor: not-allowed;
  transform: none;
}

.error-message {
  margin: 0;
  color: #f87171;
  font-size: 0.85rem;
  line-height: 1.4;
}
</style>