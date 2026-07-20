<template>
  <div>
    <h2>Manajemen Aset</h2>
    <form @submit.prevent="createAsset" class="asset-form">
      <input v-model="assetCode" placeholder="Kode Aset" />
      <input v-model="assetName" placeholder="Nama Aset" />
      <select v-model="status">
        <option>Active</option>
        <option>Inactive</option>
      </select>
      <button type="submit">Tambah Aset</button>
    </form>

    <div class="search">
      <input v-model="q" placeholder="Cari aset..." @keyup.enter="search" />
      <button @click="search">Cari</button>
    </div>

    <ul>
      <li v-for="a in assets" :key="a.id">{{ a.asset_code || a.AssetCode }} - {{ a.asset_name || a.AssetName }} - {{ a.status }}</li>
    </ul>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../api'

const assetCode = ref('')
const assetName = ref('')
const status = ref('Active')
const q = ref('')
const assets = ref([])

async function createAsset() {
  try {
    await api.post('/assets', {
      asset_code: assetCode.value,
      asset_name: assetName.value,
      status: status.value,
    }, {
      headers: { 'X-User-Role': 'hod' }
    })
    assetCode.value = ''
    assetName.value = ''
    await search()
  } catch (e) {
    alert('gagal membuat aset: ' + (e?.response?.data || e.message))
  }
}

async function search() {
  if (!q.value) return
  try {
    const res = await api.get('/assets?q=' + encodeURIComponent(q.value))
    assets.value = res.data
  } catch (e) {
    alert('gagal mencari aset')
  }
}
</script>

<style scoped>
.asset-form { display:flex; gap:8px; margin-bottom:12px }
.search { margin:12px 0 }
</style>
