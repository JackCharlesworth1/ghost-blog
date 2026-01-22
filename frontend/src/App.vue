<template>
  <div id="app">
    <nav class="navbar" v-if="!isAdminRoute">
      <div class="nav-container">
        <h1 class="logo">Ghost Blog</h1>
        <div class="nav-links">
          <router-link to="/browse" class="nav-link">Browse</router-link>
          <router-link to="/submit" class="nav-link">Submit</router-link>
        </div>
      </div>
    </nav>
    <main class="main-content" :class="{ 'admin-view': isAdminRoute }">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
</script>

<style scoped>
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #1a1a1a;
}

.navbar {
  background-color: #2a2a2a;
  border-bottom: 1px solid #3a3a3a;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #e0e0e0;
  margin: 0;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-link {
  color: #b0b0b0;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s;
  padding: 0.5rem 1rem;
  border-radius: 4px;
}

.nav-link:hover {
  color: #e0e0e0;
  background-color: #3a3a3a;
}

.nav-link.router-link-active {
  color: #ffffff;
  background-color: #4a4a4a;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.main-content.admin-view {
  background-color: #1a1a1a;
}
</style>
