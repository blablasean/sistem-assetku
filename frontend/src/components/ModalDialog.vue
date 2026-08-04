<template>
  <Teleport to="body">
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
  </Teleport>
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
  background: rgba(15, 23, 42, 0.55);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99999 !important;
  padding: 16px;
}

.modal-card {
  width: 100%;
  background: #ffffff;
  border-radius: 16px !important;
  border: 1px solid #e2e8f0;
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.25), 0 0 0 1px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  animation: popIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes popIn {
  from { opacity: 0; transform: scale(0.96) translateY(8px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #ffffff;
  color: #0f172a;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.08rem;
  color: #0f172a;
  font-weight: 800;
  letter-spacing: -0.02em;
}

.close-btn {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  font-size: 0.95rem;
  color: #64748b;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50% !important;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
  border-color: #cbd5e1;
  transform: scale(1.05);
}

.modal-body {
  padding: 20px;
  max-height: 82vh;
  overflow-y: auto;
  overflow-x: hidden;
}

@media (max-width: 640px) {
  .modal-backdrop {
    padding: 8px !important;
  }
  .modal-card {
    max-width: 100% !important;
    border-radius: 12px !important;
    overflow: hidden !important;
  }
  .modal-header {
    padding: 12px 14px !important;
  }
  .modal-header h3 {
    font-size: 0.95rem !important;
  }
  .modal-body {
    padding: 14px 10px !important;
    max-height: 85vh !important;
    overflow-x: hidden !important;
  }
}
</style>
