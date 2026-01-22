<template>
  <div class="browse-container">
    <div v-if="loading && posts.length === 0" class="loading">
      Loading posts...
    </div>

    <div v-else-if="error && posts.length === 0" class="error">
      {{ error }}
      <button @click="loadPosts" class="retry-button">Retry</button>
    </div>

    <div v-else-if="posts.length === 0" class="empty">
      <div class="empty-posts">
        <div class="empty-message">
          <h3>No posts yet!</h3>
          <p>Be the first to share something.</p>
        </div>
        <router-link to="/submit" class="create-post-button">
          Create a Post
        </router-link>
      </div>
    </div>

    <div v-else class="posts-scroll-container">
      <div
        v-for="(post, index) in posts"
        :key="post.id"
        class="post-slide"
        :class="{ active: index === currentIndex }"
      >
        <div class="post-content">
          <div class="post-header">
            <h2 class="post-title">{{ post.title }}</h2>
            <span class="post-country">{{ post.country }}</span>
          </div>
          <div class="post-body">
            <p class="post-text">{{ post.content }}</p>
          </div>
          <div class="post-footer">
            <span class="post-date">{{ formatDate(post.created_at) }}</span>
          </div>
        </div>

        <div class="scroll-indicator">
          <div v-if="index > 0" class="scroll-hint up">↑</div>
          <div v-if="index < posts.length - 1 || hasMore" class="scroll-hint down">
            ↓
          </div>
        </div>
      </div>

      <!-- End of posts message -->
      <div
        v-if="!hasMore && currentIndex >= posts.length - 1"
        class="post-slide"
        :class="{ active: currentIndex === posts.length }"
      >
        <div class="end-of-posts">
          <div class="end-message">
            <h3>That's all for now!</h3>
            <p>You've reached the end of the posts.</p>
          </div>
          <router-link to="/submit" class="create-post-button">
            Create a Post
          </router-link>
        </div>
        <div class="scroll-indicator">
          <div class="scroll-hint up">↑</div>
        </div>
      </div>

      <div v-if="loading" class="loading-more">
        <div class="spinner"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { usePostsStore } from '../stores/posts'

const postsStore = usePostsStore()

const posts = computed(() => postsStore.posts)
const currentIndex = computed(() => postsStore.currentIndex)
const loading = computed(() => postsStore.loading)
const error = computed(() => postsStore.error)
const hasMore = computed(() => postsStore.hasMore)

let isScrolling = false

const loadPosts = async () => {
  await postsStore.fetchPosts(10)
}

const handleScroll = (e) => {
  if (isScrolling) return

  const delta = e.deltaY
  const threshold = 50

  if (Math.abs(delta) > threshold) {
    isScrolling = true

    if (delta > 0) {
      // Scroll down
      postsStore.nextPost()
    } else {
      // Scroll up
      postsStore.previousPost()
    }

    setTimeout(() => {
      isScrolling = false
    }, 800)
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`
  if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`
  if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`

  return date.toLocaleDateString()
}

onMounted(() => {
  postsStore.reset()
  loadPosts()
  window.addEventListener('wheel', handleScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('wheel', handleScroll)
})
</script>

<style scoped>
.browse-container {
  width: 100%;
  height: calc(100vh - 80px);
  overflow: hidden;
  position: relative;
  background-color: #1a1a1a;
}

.loading,
.error,
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #b0b0b0;
  font-size: 1.25rem;
  gap: 1rem;
}

.empty-posts {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2rem;
  background-color: #2a2a2a;
  border-radius: 12px;
  padding: 3rem;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
}

.empty-message {
  text-align: center;
}

.empty-message h3 {
  color: #e0e0e0;
  font-size: 1.75rem;
  margin: 0 0 0.5rem 0;
}

.empty-message p {
  color: #b0b0b0;
  font-size: 1.125rem;
  margin: 0;
}

.retry-button {
  background-color: #4a4a4a;
  color: #e0e0e0;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.retry-button:hover {
  background-color: #5a5a5a;
}

.submit-link {
  color: #90ee90;
  text-decoration: underline;
}

.posts-scroll-container {
  height: 100%;
  position: relative;
}

.post-slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transform: translateY(100%);
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  pointer-events: none;
}

.post-slide.active {
  opacity: 1;
  transform: translateY(0);
  pointer-events: auto;
}

.post-content {
  background-color: #2a2a2a;
  border-radius: 12px;
  padding: 3rem;
  max-width: 700px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
}

.post-header {
  margin-bottom: 2rem;
  border-bottom: 1px solid #3a3a3a;
  padding-bottom: 1rem;
}

.post-title {
  color: #e0e0e0;
  font-size: 2rem;
  margin: 0 0 0.5rem 0;
  word-wrap: break-word;
}

.post-country {
  display: inline-block;
  background-color: #3a3a3a;
  color: #b0b0b0;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
}

.post-body {
  margin-bottom: 2rem;
}

.post-text {
  color: #d0d0d0;
  font-size: 1.125rem;
  line-height: 1.8;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.post-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 1rem;
  border-top: 1px solid #3a3a3a;
}

.post-date {
  color: #7a7a7a;
  font-size: 0.875rem;
}

.scroll-indicator {
  position: fixed;
  right: 2rem;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.scroll-hint {
  color: #5a5a5a;
  font-size: 2rem;
  animation: bounce 2s infinite;
}

.scroll-hint.up {
  animation-direction: reverse;
}

@keyframes bounce {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.loading-more {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #3a3a3a;
  border-top-color: #e0e0e0;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Scrollbar styles for post content */
.post-content::-webkit-scrollbar {
  width: 8px;
}

.post-content::-webkit-scrollbar-track {
  background: #1a1a1a;
  border-radius: 4px;
}

.post-content::-webkit-scrollbar-thumb {
  background: #4a4a4a;
  border-radius: 4px;
}

.post-content::-webkit-scrollbar-thumb:hover {
  background: #5a5a5a;
}

.end-of-posts {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2rem;
  background-color: #2a2a2a;
  border-radius: 12px;
  padding: 3rem;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
}

.end-message {
  text-align: center;
}

.end-message h3 {
  color: #e0e0e0;
  font-size: 1.75rem;
  margin: 0 0 0.5rem 0;
}

.end-message p {
  color: #b0b0b0;
  font-size: 1.125rem;
  margin: 0;
}

.create-post-button {
  background-color: #90ee90;
  color: #1a1a1a;
  padding: 0.875rem 2rem;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  font-size: 1rem;
  transition: all 0.3s;
  box-shadow: 0 4px 8px rgba(144, 238, 144, 0.2);
}

.create-post-button:hover {
  background-color: #7cda7c;
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(144, 238, 144, 0.3);
}
</style>
