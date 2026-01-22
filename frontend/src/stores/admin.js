import { defineStore } from 'pinia'
import axios from 'axios'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    posts: [],
    total: 0,
    loading: false,
    error: null,
    authenticated: false,
    credentials: null
  }),

  actions: {
    setCredentials(username, password) {
      this.credentials = btoa(`${username}:${password}`)
      this.authenticated = true
    },

    logout() {
      this.credentials = null
      this.authenticated = false
      this.posts = []
      this.total = 0
    },

    async fetchAllPosts(limit = 50, offset = 0) {
      if (!this.credentials) {
        this.error = 'Not authenticated'
        return
      }

      this.loading = true
      this.error = null
      try {
        const response = await axios.get('/api/admin/posts', {
          params: { limit, offset },
          headers: {
            Authorization: `Basic ${this.credentials}`
          }
        })

        this.posts = response.data.posts || []
        this.total = response.data.total || 0
      } catch (error) {
        if (error.response?.status === 401) {
          this.error = 'Invalid credentials'
          this.logout()
        } else {
          this.error = error.response?.data?.error || 'Failed to fetch posts'
        }
        console.error('Error fetching admin posts:', error)
      } finally {
        this.loading = false
      }
    },

    async deletePost(postId) {
      if (!this.credentials) {
        this.error = 'Not authenticated'
        return
      }

      this.loading = true
      this.error = null
      try {
        await axios.delete(`/api/admin/posts/${postId}`, {
          headers: {
            Authorization: `Basic ${this.credentials}`
          }
        })

        // Remove from local state
        this.posts = this.posts.filter(post => post.id !== postId)
        this.total--
      } catch (error) {
        if (error.response?.status === 401) {
          this.error = 'Invalid credentials'
          this.logout()
        } else {
          this.error = error.response?.data?.error || 'Failed to delete post'
        }
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
