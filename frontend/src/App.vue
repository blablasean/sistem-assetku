<template>
  <div class="app-container">
    <aside class="sidebar">
      <h2>AsetKu</h2>
      <nav>
        <ul>
          <li><a href="/src/views/Login.vue">Dashboard</a></li>
          <li><a href="/src/views/AssetManagement.vue">Manajemen Aset</a></li>
          <li><a href="/src/views/WorkOrder.vue">Work Order</a></li>
          <li><a href="/src/views/UtilityMonitoring.vue">Utility Monitoring</a></li>
        </ul>
      </nav>
    </aside>

    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
/* CSS Dasar untuk membagi layar menjadi Sidebar & Konten Utama */
.app-container {
  display: flex;
  min-height: 100vh;
  font-family: sans-serif;
}
.sidebar {
  width: 250px;
  background-color: #1e293b;
  color: white;
  padding: 20px;
}
.sidebar ul {
  list-style: none;
  padding: 0;
}
.sidebar li a {
  color: #cbd5e1;
  text-decoration: none;
  display: block;
  padding: 10px 0;
}
.main-content {
  flex: 1;
  background-color: #f8fafc;
  padding: 20px;
}
</style>