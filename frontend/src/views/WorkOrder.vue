<template>
  <div>
    <h2>Work Orders</h2>
    <form @submit.prevent="createWO" class="wo-form">
      <input v-model="assetId" placeholder="Asset ID" />
      <input v-model="description" placeholder="Deskripsi" />
      <button type="submit">Buat WO</button>
    </form>

    <div class="assign">
      <h3>Assign Engineer</h3>
      <input v-model.number="assignWoId" placeholder="WO ID" />
      <input v-model.number="assignEngineerId" placeholder="Engineer ID" />
      <button @click="assign">Assign</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../api'

const assetId = ref('')
const description = ref('')
const assignWoId = ref(0)
const assignEngineerId = ref(0)

async function createWO() {
  try {
    await api.post('/workorders', {
      asset_id: parseInt(assetId.value),
      description: description.value,
      requester_id: 1,
    }, {
      headers: { 'X-User-Role': 'hod' }
    })
    assetId.value = ''
    description.value = ''
    alert('Work order created')
  } catch (e) {
    alert('gagal membuat work order')
  }
}

async function assign() {
  try {
    await api.post('/workorders/assign', {
      wo_id: assignWoId.value,
      engineer_id: assignEngineerId.value,
    }, { headers: { 'X-User-Role': 'hod' } })
    alert('Engineer assigned')
  } catch (e) {
    alert('assign gagal')
  }
}
</script>

<style scoped>
.wo-form { display:flex; gap:8px; margin-bottom:16px }
.assign { margin-top:20px }
</style>
