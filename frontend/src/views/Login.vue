<template>
  <div class="login-screen">
    <div class="login-panel">
      <div class="login-brand">
        <div class="brand-badge">
          <img src="/assets/logo.png" alt="AsetKu Logo" class="login-logo-img" @error="logoFailed = true" v-if="!logoFailed" />
          <span v-else class="login-text-logo">A</span>
        </div>
        <div>
          <p class="brand-label">AssetKu</p>
          <p class="brand-subtitle">Work Order & Asset Management System</p>
        </div>
      </div>

      <h1>Selamat Datang</h1>
      <p>Masukkan username dan password Anda.</p>

      <form @submit.prevent="login" class="login-form">
        <label>
          <span>Username</span>
          <input 
            type="text" 
            v-model="username" 
            placeholder="Masukkan username"
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
              <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/><circle cx="12" cy="12" r="3"/></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9.88 9.88a3 3 0 1 0 4.24 4.24"/><path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"/><path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"/><line x1="2" x2="22" y1="2" y2="22"/></svg>
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
const logoFailed = ref(false)

function onLogoError(e) {
  if (e && e.target && e.target.src && e.target.src.includes('/assets/logo.png')) {
    e.target.src = '/logo.png'
  } else {
    logoFailed.value = true
  }
}
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

    // Save session in sessionStorage (session ends automatically when tab/browser is closed)
    localStorage.clear() // Ensure legacy persistent tokens are removed

    sessionStorage.setItem('token', token)
    sessionStorage.setItem('user_role', role)
    sessionStorage.setItem('user_name', name)
    sessionStorage.setItem('username', data.username || username.value)

    if (data.avatar) {
      sessionStorage.setItem('user_avatar', data.avatar)
    } else {
      sessionStorage.removeItem('user_avatar')
    }

    router.push('/dashboard')
  } catch (err) {
    console.error('Login error:', err)
    const apiErrorMsg = err.response?.data?.message || err.response?.data?.details || err.message
    error.value = apiErrorMsg || 'Username atau password salah. Silakan periksa kembali kredensial database Anda.'
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
  background: #f8fafc;
  position: relative;
  overflow: hidden;
}

.login-panel {
  width: min(400px, 100%);
  background: #ffffff;
  border-radius: 6px !important;
  box-shadow: 0 10px 30px -5px rgba(0, 0, 0, 0.08), 0 0 0 1px #e2e8f0;
  padding: 40px 36px;
  border: 1px solid #e2e8f0;
  color: #0f172a;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-brand {
  display: flex;
  gap: 14px;
  align-items: center;
  margin-bottom: 24px;
}

.brand-badge {
  width: 46px;
  height: 46px;
  border: 1px solid #e2e8f0;
  border-radius: 6px !important;
  background: #f8fafc;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
}

.login-logo-img {
  width: 39px;
  height: 39px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 20%;
  object-fit: contain;
}

.login-text-logo {
  font-weight: 900;
  font-size: 1.3rem;
  color: #ffffff;
}

.brand-label {
  margin: 0;
  font-weight: 800;
  font-size: 1.15rem;
  color: #0f172a;
  letter-spacing: -0.02em;
}

.brand-subtitle {
  margin: 2px 0 0;
  color: #64748b;
  font-size: 0.82rem;
  font-weight: 600;
}

h1 {
  margin: 0 0 6px;
  font-size: 1.75rem;
  color: #0f172a;
  font-weight: 800;
  letter-spacing: -0.02em;
}

p {
  margin: 0 0 24px;
  color: #64748b;
  font-size: 0.88rem;
}

.login-form {
  display: grid;
  gap: 18px;
}

.login-form label {
  display: grid;
  gap: 8px;
  font-weight: 700;
  font-size: 0.85rem;
  color: #0f172a;
}

.login-form input {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 4px !important;
  padding: 12px 14px;
  font-size: 0.95rem;
  background: #ffffff;
  color: #0f172a;
  outline: none;
  transition: all 0.15s ease;
}

.login-form input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
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
  font-size: 1.1rem;
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
  padding: 13px 0;
  border-radius: 4px !important;
  font-size: 0.95rem;
  font-weight: 700;
  color: #ffffff;
  background: #0f172a;
  border: 1px solid #0f172a;
  cursor: pointer;
  transition: all 0.15s ease;
  margin-top: 8px;
}

.submit-login-btn:hover {
  background: #1e293b;
}

.submit-login-btn:disabled {
  background: #94a3b8;
  border-color: #94a3b8;
  cursor: not-allowed;
}

.error-message {
  margin: 0;
  color: #dc2626;
  font-size: 0.85rem;
  line-height: 1.4;
  font-weight: 600;
}
</style>