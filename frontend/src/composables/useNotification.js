import { ref } from 'vue'

const showToast = ref(false)
const toastMsg = ref('')
const toastType = ref('success')
let toastTimer = null

/**
 * Vue 3 Composable for UI Toast Notifications
 */
export function useNotification() {
  function notify(msg, type = 'success', duration = 4000) {
    if (toastTimer) clearTimeout(toastTimer)
    toastMsg.value = msg
    toastType.value = type
    showToast.value = true
    toastTimer = setTimeout(() => {
      showToast.value = false
    }, duration)
  }

  function dismiss() {
    if (toastTimer) clearTimeout(toastTimer)
    showToast.value = false
  }

  return {
    showToast,
    toastMsg,
    toastType,
    notify,
    dismiss
  }
}
