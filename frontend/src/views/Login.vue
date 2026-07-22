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
      <p>Masuk ke portal operasional hotel Anda.</p>

      <form @submit.prevent="login" class="login-form">
        <label>
          <span>Pilih Role Login / Preset Demo</span>
          <select v-model="selectedRolePreset" @change="applyRolePreset" class="role-select">
            <option value="management">👔 Management Engineer (Supervisor)</option>
            <option value="hod">🗂️ HOD Engineer (Head of Department)</option>
            <option value="engineer">🛠️ Staff Engineer (Teknisi Lapangan)</option>
            <option value="external">🛎️ External User (Staff Hotel / Departemen Lain)</option>
          </select>
        </label>

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
          {{ isLoading ? 'Sedang masuk...' : 'Masuk Ke Sistem' }}
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

const selectedRolePreset = ref('management')
const username = ref('admin')
const password = ref('password')
const error = ref('')
const isLoading = ref(false)
const router = useRouter()

function applyRolePreset() {
  const presets = {
    management: { user: 'mgr_eng', name: 'Budi (Supervisor)' },
    hod: { user: 'hod_eng', name: 'Pak Alex (HOD Engineer)' },
    engineer: { user: 'tech_deni', name: 'Deni (Staff Engineer)' },
    external: { user: 'staff_frontdesk', name: 'Rina (Staff Front Desk)' }
  }
  const p = presets[selectedRolePreset.value] || { user: 'admin', name: 'Admin Hotel' }
  username.value = p.user
  password.value = 'password'
}

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

    const token = res.data.data?.token || 'DEMO_JWT_TOKEN_' + Date.now()
    
    // Store token and user role in localStorage
    localStorage.setItem('token', token)
    localStorage.setItem('user_role', selectedRolePreset.value)
    localStorage.setItem('user_name', username.value.toUpperCase())

    router.push('/dashboard')
  } catch (err) {
    // Demo fallback for smooth offline testing
    localStorage.setItem('token', 'DEMO_JWT_TOKEN_' + Date.now())
    localStorage.setItem('user_role', selectedRolePreset.value)
    localStorage.setItem('user_name', username.value.toUpperCase())
    router.push('/dashboard')
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