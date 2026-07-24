<template>
  <div v-if="show" class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-card" :style="{ maxWidth: maxWidth }">
      <div class="modal-header">
        <h3>{{ title }}</h3>
        <button class="close-btn" @click="$emit('close')">✕</button>
      </div>
      <div class="modal-body">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: { type: Boolean, default: false },
  title: { type: String, default: 'Dialog' },
  maxWidth: { type: String, default: '550px' }
})
defineEmits(['close'])
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  padding: 16px;
}

.modal-card {
  width: 100%;
  background: #ffffff;
  border-radius: 2px !important;
  border: 1px solid #cbd5e1;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  animation: popIn 0.15s ease-out;
}

@keyframes popIn {
  from { opacity: 0; transform: scale(0.98); }
  to { opacity: 1; transform: scale(1); }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  background: #0f172a;
  color: #ffffff;
  border-bottom: 2px solid #d97706;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.05rem;
  color: #ffffff;
  font-weight: 700;
  letter-spacing: 0.3px;
}

.close-btn {
  background: transparent;
  border: none;
  font-size: 1.1rem;
  color: #cbd5e1;
  cursor: pointer;
  padding: 2px 8px;
  border-radius: 2px !important;
}

.close-btn:hover {
  background: #1e293b;
  color: #ffffff;
}

.modal-body {
  padding: 20px;
  max-height: 80vh;
  overflow-y: auto;
}
</style>
