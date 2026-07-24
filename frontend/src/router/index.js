import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import AssetManagement from '../views/AssetManagement.vue'
import WorkOrder from '../views/WorkOrder.vue'
import PreventiveMaintenance from '../views/PreventiveMaintenance.vue'
import ActivityLog from '../views/ActivityLog.vue'
import UserManagement from '../views/UserManagement.vue'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: Login, meta: { public: true } },
  { path: '/dashboard', component: Dashboard },
  { path: '/assets', component: AssetManagement, meta: { roles: ['management', 'hod', 'engineer', 'admin'] } },
  { path: '/workorders', component: WorkOrder },
  { path: '/maintenance', component: PreventiveMaintenance, meta: { roles: ['management', 'hod', 'engineer', 'admin'] } },
  { path: '/activitylogs', component: ActivityLog },
  { path: '/users', component: UserManagement, meta: { roles: ['admin'] } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard checking sessionStorage for session re-login requirement
router.beforeEach((to, from, next) => {
  const token = sessionStorage.getItem('token') || localStorage.getItem('token')
  const userRole = sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external'

  if (!to.meta.public && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else if (to.meta.roles && !to.meta.roles.includes(userRole)) {
    alert(`⛔ Akses Ditolak!\nRole "${userRole.toUpperCase()}" tidak memiliki wewenang untuk mengakses halaman ini. Anda dialihkan ke Work Order.`)
    next('/workorders')
  } else {
    next()
  }
})

export default router
