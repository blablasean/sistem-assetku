<template>
  <div class="app-container">
    <aside class="sidebar">
      <h2>AsetKu</h2>
      <nav>
        <ul>
          <li><router-link to="/dashboard">Dashboard</router-link></li>
          <li><router-link to="/assets">Manajemen Aset</router-link></li>
          <li><router-link to="/workorders">Work Order</router-link></li>
          <li><router-link to="/utility">Utility Monitoring</router-link></li>
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
.app-container {
  display: flex;
  min-height: 100vh;
  font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}
.sidebar {
  width: 220px;
  background-color: #0f172a;
  color: white;
  padding: 20px;
}
.sidebar ul {
  list-style: none;
  padding: 0;
}
.sidebar li {
  margin-bottom: 12px;
}
.sidebar a {
  color: #cbd5e1;
  text-decoration: none;
}
.main-content {
  flex: 1;
  background-color: #f8fafc;
  padding: 24px;
}
</style>