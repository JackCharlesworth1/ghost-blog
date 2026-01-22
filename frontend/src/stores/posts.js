import { defineStore } from 'pinia'
import axios from 'axios'

export const usePostsStore = defineStore('posts', {
  state: () => ({
    posts: [],
    currentIndex: 0,
    loading: false,
    error: null,
    hasMore: true
  }),

  getters: {
    currentPost: (state) => state.posts[state.currentIndex] || null
  },

  actions: {
    async fetchPosts(limit = 10) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get('/api/posts', {
          params: { limit, offset: this.posts.length }
        })

        if (!response.data || response.data.length === 0) {
          this.hasMore = false
        } else {
          this.posts.push(...response.data)
        }
      } catch (error) {
        // Only set error if we have no posts (initial load failure)
        // If we already have posts, just mark as no more posts available
        if (this.posts.length === 0) {
          this.error = error.response?.data?.error || 'Failed to fetch posts'
        } else {
          this.hasMore = false
        }
        console.error('Error fetching posts:', error)
      } finally {
        this.loading = false
      }
    },

    async createPost(postData) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.post('/api/posts', postData)
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create post'
        throw error
      } finally {
        this.loading = false
      }
    },

    nextPost() {
      // Allow scrolling to one past the last post to show end-of-posts message
      const maxIndex = this.hasMore ? this.posts.length - 1 : this.posts.length
      if (this.currentIndex < maxIndex) {
        this.currentIndex++
      }

      // Prefetch more posts when nearing the end
      if (this.currentIndex >= this.posts.length - 3 && this.hasMore && !this.loading) {
        this.fetchPosts()
      }
    },

    previousPost() {
      if (this.currentIndex > 0) {
        this.currentIndex--
      }
    },

    reset() {
      this.posts = []
      this.currentIndex = 0
      this.hasMore = true
      this.error = null
    }
  }
})
