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
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: capitalize;
}

.active { background: #dcfce7; color: #166534; }
.maintenance, .in-progress { background: #fef3c7; color: #92400e; }
.damaged, .emergency, .high { background: #fee2e2; color: #991b1b; }
.open { background: #e0f2fe; color: #075985; }
.under-review { background: #f3e8ff; color: #6b21a8; }
.completed, .closed, .finish { background: #d1fae5; color: #065f46; }
.cancelled, .retired { background: #f1f5f9; color: #475569; }
.reserved { background: #ffedd5; color: #9a3412; }
.medium { background: #e0e7ff; color: #3730a3; }
.low { background: #f1f5f9; color: #334155; }
</style>
