/**
 * Export Utility Module
 * Provides clean, lightweight, decoupled helper functions for Excel exporting and printing.
 */

/**
 * Export structured data array to an Excel (.xls) file cleanly.
 * @param {string} filename - The name of the file to download (e.g. 'Laporan_Aset.xls')
 * @param {Array<string>} headers - Table column headers
 * @param {Array<Array<any>>} rows - Table data rows
 */
export function exportToExcel(filename, headers, rows) {
  let tableHtml = '<table border="1"><thead><tr>'
  
  headers.forEach(h => {
    tableHtml += `<th style="background:#007aff;color:#ffffff;font-weight:bold;padding:8px;">${escapeHtml(h)}</th>`
  })
  tableHtml += '</tr></thead><tbody>'

  rows.forEach(row => {
    tableHtml += '<tr>'
    row.forEach(cell => {
      tableHtml += `<td style="padding:6px;">${escapeHtml(cell != null ? String(cell) : '')}</td>`
    })
    tableHtml += '</tr>'
  })
  tableHtml += '</tbody></table>'

  const template = `
    <html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
    <head>
      <meta charset="utf-8">
      <!--[if gte mso 9]><xml><x:ExcelWorkbook><x:ExcelWorksheets><x:ExcelWorksheet>
      <x:Name>Laporan</x:Name>
      <x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions>
      </x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]-->
    </head>
    <body>${tableHtml}</body>
    </html>
  `

  const blob = new Blob(['\ufeff', template], { type: 'application/vnd.ms-excel;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename.endsWith('.xls') ? filename : `${filename}.xls`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Trigger native print dialog safely.
 */
export function triggerPrint() {
  window.print()
}

/**
 * Helper to escape HTML special characters.
 */
function escapeHtml(str) {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}
