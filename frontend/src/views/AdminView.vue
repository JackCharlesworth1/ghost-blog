<template>
  <div class="admin-container">
    <div v-if="!adminStore.authenticated" class="login-wrapper">
      <div class="login-box">
        <h2 class="login-title">Admin Login</h2>
        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label for="username">Username</label>
            <input
              id="username"
              v-model="username"
              type="text"
              required
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="form-input"
            />
          </div>
          <div v-if="loginError" class="error-message">
            {{ loginError }}
          </div>
          <button type="submit" :disabled="loading" class="login-button">
            {{ loading ? 'Logging in...' : 'Login' }}
          </button>
        </form>
      </div>
    </div>

    <div v-else class="admin-panel">
      <div class="admin-header">
        <h1>Admin Panel</h1>
        <div class="admin-actions">
          <span class="post-count">Total Posts: {{ adminStore.total }}</span>
          <button @click="handleLogout" class="logout-button">Logout</button>
        </div>
      </div>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <div v-if="loading" class="loading">Loading posts...</div>

      <div v-else-if="adminStore.posts.length === 0" class="empty">
        No posts found.
      </div>

      <div v-else class="posts-table">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Title</th>
              <th>Country</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="post in adminStore.posts" :key="post.id">
              <td>{{ post.id }}</td>
              <td class="title-cell">
                <div class="title-content">
                  <strong>{{ post.title }}</strong>
                  <p class="content-preview">{{ truncate(post.content, 100) }}</p>
                </div>
              </td>
              <td>{{ getCountryFlag(post.country) }} {{ post.country }}</td>
              <td>{{ formatDate(post.created_at) }}</td>
              <td>
                <button
                  @click="handleDelete(post.id)"
                  :disabled="deleting === post.id"
                  class="delete-button"
                >
                  {{ deleting === post.id ? 'Deleting...' : 'Delete' }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAdminStore } from '../stores/admin'

const adminStore = useAdminStore()

const username = ref('')
const password = ref('')
const loginError = ref(null)
const deleting = ref(null)

const loading = computed(() => adminStore.loading)
const error = computed(() => adminStore.error)

const handleLogin = async () => {
  loginError.value = null
  adminStore.setCredentials(username.value, password.value)
  await adminStore.fetchAllPosts()

  if (adminStore.error) {
    loginError.value = adminStore.error
  }
}

const handleLogout = () => {
  adminStore.logout()
  username.value = ''
  password.value = ''
}

const handleDelete = async (postId) => {
  if (!confirm('Are you sure you want to delete this post?')) {
    return
  }

  deleting.value = postId
  try {
    await adminStore.deletePost(postId)
  } catch (err) {
    alert('Failed to delete post')
  } finally {
    deleting.value = null
  }
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString()
}

const truncate = (text, length) => {
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}

const getCountryFlag = (countryName) => {
  // Map country names to ISO 3166-1 alpha-2 codes
  const countryToCode = {
    'United States': 'US',
    'United Kingdom': 'GB',
    'Canada': 'CA',
    'Australia': 'AU',
    'Germany': 'DE',
    'France': 'FR',
    'Italy': 'IT',
    'Spain': 'ES',
    'Japan': 'JP',
    'China': 'CN',
    'India': 'IN',
    'Brazil': 'BR',
    'Mexico': 'MX',
    'Russia': 'RU',
    'South Korea': 'KR',
    'Netherlands': 'NL',
    'Sweden': 'SE',
    'Norway': 'NO',
    'Denmark': 'DK',
    'Finland': 'FI',
    'Poland': 'PL',
    'Belgium': 'BE',
    'Austria': 'AT',
    'Switzerland': 'CH',
    'Ireland': 'IE',
    'Portugal': 'PT',
    'Greece': 'GR',
    'Czech Republic': 'CZ',
    'Romania': 'RO',
    'Hungary': 'HU',
    'New Zealand': 'NZ',
    'Singapore': 'SG',
    'Thailand': 'TH',
    'Vietnam': 'VN',
    'Philippines': 'PH',
    'Indonesia': 'ID',
    'Malaysia': 'MY',
    'Turkey': 'TR',
    'Israel': 'IL',
    'Saudi Arabia': 'SA',
    'United Arab Emirates': 'AE',
    'South Africa': 'ZA',
    'Egypt': 'EG',
    'Argentina': 'AR',
    'Chile': 'CL',
    'Colombia': 'CO',
    'Peru': 'PE',
    'Venezuela': 'VE',
    'Ukraine': 'UA',
    'Bulgaria': 'BG',
    'Croatia': 'HR',
    'Slovakia': 'SK',
    'Slovenia': 'SI',
    'Lithuania': 'LT',
    'Latvia': 'LV',
    'Estonia': 'EE',
    'Iceland': 'IS',
    'Luxembourg': 'LU',
    'Malta': 'MT',
    'Cyprus': 'CY',
    'Pakistan': 'PK',
    'Bangladesh': 'BD',
    'Sri Lanka': 'LK',
    'Hong Kong': 'HK',
    'Taiwan': 'TW',
    'Local': '🏠',
    'Unknown': '🌍'
  }

  const code = countryToCode[countryName]

  // Special handling for non-country codes
  if (code === '🏠' || code === '🌍') {
    return code
  }

  if (!code) {
    return '🌍' // Default globe emoji for unmapped countries
  }

  // Convert country code to flag emoji
  // Flag emojis are created using regional indicator symbols
  // A = 0x1F1E6, B = 0x1F1E7, ..., Z = 0x1F1FF
  const codePoints = [...code].map(char => 0x1F1E6 + char.charCodeAt(0) - 65)
  return String.fromCodePoint(...codePoints)
}

onMounted(() => {
  if (adminStore.authenticated) {
    adminStore.fetchAllPosts()
  }
})
</script>

<style scoped>
.admin-container {
  min-height: 100vh;
  background-color: #1a1a1a;
}

.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 2rem;
}

.login-box {
  background-color: #2a2a2a;
  border-radius: 8px;
  padding: 2rem;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
}

.login-title {
  color: #e0e0e0;
  margin-bottom: 1.5rem;
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  color: #b0b0b0;
  font-weight: 500;
}

.form-input {
  background-color: #1a1a1a;
  border: 1px solid #3a3a3a;
  color: #e0e0e0;
  padding: 0.75rem;
  border-radius: 4px;
  transition: border-color 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #5a5a5a;
}

.error-message {
  background-color: #4a2d2d;
  color: #ff9090;
  padding: 0.75rem;
  border-radius: 4px;
  text-align: center;
}

.login-button {
  background-color: #4a4a4a;
  color: #e0e0e0;
  border: none;
  padding: 0.75rem;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-button:hover:not(:disabled) {
  background-color: #5a5a5a;
}

.login-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.admin-panel {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #3a3a3a;
}

.admin-header h1 {
  color: #e0e0e0;
  margin: 0;
}

.admin-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.post-count {
  color: #b0b0b0;
  font-size: 1rem;
}

.logout-button {
  background-color: #4a2d2d;
  color: #ff9090;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.logout-button:hover {
  background-color: #5a3d3d;
}

.loading,
.empty {
  text-align: center;
  color: #b0b0b0;
  padding: 2rem;
  font-size: 1.25rem;
}

.posts-table {
  background-color: #2a2a2a;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background-color: #3a3a3a;
}

th {
  color: #e0e0e0;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
}

tbody tr {
  border-bottom: 1px solid #3a3a3a;
  transition: background-color 0.2s;
}

tbody tr:hover {
  background-color: #333333;
}

tbody tr:last-child {
  border-bottom: none;
}

td {
  color: #d0d0d0;
  padding: 1rem;
  vertical-align: top;
}

.title-cell {
  max-width: 500px;
}

.title-content strong {
  display: block;
  margin-bottom: 0.5rem;
  color: #e0e0e0;
}

.content-preview {
  color: #b0b0b0;
  font-size: 0.875rem;
  margin: 0;
  line-height: 1.5;
}

.delete-button {
  background-color: #4a2d2d;
  color: #ff9090;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  white-space: nowrap;
}

.delete-button:hover:not(:disabled) {
  background-color: #5a3d3d;
}

.delete-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
