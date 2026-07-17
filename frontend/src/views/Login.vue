<template>
  <div class="login-container">
    <form @submit.prevent="login" class="login-form">
      <h2>Login ke AsetKu</h2>
      <div class="form-group">
        <label for="username">Username:</label>
        <input type="text" id="username" v-model="username" required />
      </div>
      <div class="form-group">
        <label for="password">Password:</label>
        <input type="password" id="password" v-model="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
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