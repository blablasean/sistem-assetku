import { ref, computed } from 'vue'

const userRole = ref(sessionStorage.getItem('user_role') || 'external')
const userName = ref(sessionStorage.getItem('user_name') || 'Guest User')

/**
 * Vue 3 Composable for Authentication & Role Permissions
 */
export function useAuth() {
  function syncAuth() {
    userRole.value = sessionStorage.getItem('user_role') || 'external'
    userName.value = sessionStorage.getItem('user_name') || 'Guest User'
  }

  const isAdmin = computed(() => userRole.value === 'admin')
  const isHOD = computed(() => userRole.value === 'hod')
  const isSupervisor = computed(() => userRole.value === 'management')
  const isStaffOnly = computed(() => userRole.value === 'external')

  const canManageAssets = computed(() => {
    const r = userRole.value
    return r === 'admin' || r === 'hod' || r === 'management'
  })

  const canMutate = computed(() => {
    const r = userRole.value
    return r === 'admin' || r === 'hod'
  })

  const canAssignWorker = computed(() => {
    const r = userRole.value
    return r === 'admin' || r === 'hod' || r === 'management'
  })

  const roleLabel = computed(() => {
    const map = {
      admin: 'Administrator',
      hod: 'HOD Engineer',
      management: 'Supervisor Engineer',
      engineer: 'Staff Engineer',
      external: 'Staff Operasional'
    }
    return map[userRole.value] || 'User'
  })

  return {
    userRole,
    userName,
    syncAuth,
    isAdmin,
    isHOD,
    isSupervisor,
    isStaffOnly,
    canManageAssets,
    canMutate,
    canAssignWorker,
    roleLabel
  }
}
