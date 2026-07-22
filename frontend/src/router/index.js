import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import AssetManagement from '../views/AssetManagement.vue'
import WorkOrder from '../views/WorkOrder.vue'
import PreventiveMaintenance from '../views/PreventiveMaintenance.vue'
import ActivityLog from '../views/ActivityLog.vue'
import UtilityMonitoring from '../views/UtilityMonitoring.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard },
  { path: '/assets', component: AssetManagement },
  { path: '/workorders', component: WorkOrder },
  { path: '/maintenance', component: PreventiveMaintenance },
  { path: '/activitylogs', component: ActivityLog },
  { path: '/utility', component: UtilityMonitoring },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
