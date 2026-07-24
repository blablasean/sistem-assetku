<template>
  <span :class="['status-chip', normalizedStatus]">
    {{ status }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    default: 'Active'
  }
})

const normalizedStatus = computed(() => {
  const s = props.status ? props.status.toLowerCase().replace(/\s+/g, '-') : 'active'
  return s
})
</script>

<style scoped>
.status-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 2px !important;
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  border: 1px solid transparent;
}

.active { background: #dcfce7; color: #166534; border-color: #bbf7d0; }
.maintenance, .in-progress { background: #fef3c7; color: #92400e; border-color: #fde68a; }
.damaged, .emergency, .high { background: #fee2e2; color: #991b1b; border-color: #fca5a5; }
.open { background: #e0f2fe; color: #075985; border-color: #bae6fd; }
.under-review { background: #f3e8ff; color: #6b21a8; border-color: #e9d5ff; }
.completed, .closed, .finish { background: #d1fae5; color: #065f46; border-color: #a7f3d0; }
.cancelled, .retired { background: #f1f5f9; color: #475569; border-color: #cbd5e1; }
.reserved { background: #ffedd5; color: #9a3412; border-color: #fed7aa; }
.medium { background: #e0e7ff; color: #3730a3; border-color: #c7d2fe; }
.low { background: #f1f5f9; color: #334155; border-color: #e2e8f0; }
</style>
