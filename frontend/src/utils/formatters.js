/**
 * Centralized formatting and date utilities for AsetKu Frontend
 */

/**
 * Returns today's local date string formatted as YYYY-MM-DD
 */
export function getTodayDateStr() {
  const d = new Date()
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

/**
 * Formats a date string or timestamp into Indonesian locale string (DD MMM YYYY, HH:mm)
 */
export function formatDate(dateStr) {
  if (!dateStr) return '—'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

/**
 * Formats a number with Indonesian thousand separators (e.g. 150.000)
 */
export function formatNumber(num) {
  return (num || 0).toLocaleString('id-ID')
}

/**
 * Maps role or department code to human-readable Indonesian label
 */
export function formatDepartmentLabel(roleOrDept) {
  if (!roleOrDept) return 'Staff Operasional'
  const map = {
    dept_akunting: 'Departement Akunting',
    dept_spa: 'Departement Spa',
    dept_sales: 'Department Sales',
    dept_hr: 'Department Human Resources',
    dept_fb_kitchen: 'Department Food Beverage Kitchen',
    dept_fb_service: 'Department Food Beverage Service',
    dept_housekeeping: 'Department House Keeping',
    dept_frontoffice: 'Department Front Office',
    admin: 'Administrator',
    hod: 'HOD Engineer',
    management: 'Supervisor Engineer',
    engineer: 'Staff Engineer',
    external: 'Staff Operasional'
  }
  return map[roleOrDept] || roleOrDept
}
