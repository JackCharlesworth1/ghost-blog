<template>
  <div class="submit-container">
    <div class="submit-form-wrapper">
      <h2 class="form-title">Share Your Anonymous Story</h2>

      <div v-if="success" class="success-message">
        Post created successfully!
        <router-link to="/browse" class="browse-link">Browse posts</router-link>
      </div>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <form @submit.prevent="handleSubmit" class="submit-form">
        <div class="form-group">
          <label for="title">Title</label>
          <input
            id="title"
            v-model="title"
            type="text"
            placeholder="Enter your post title..."
            maxlength="255"
            required
            class="form-input"
          />
          <span class="char-count">{{ title.length }}/255</span>
        </div>

        <div class="form-group">
          <label for="content">Content</label>
          <textarea
            id="content"
            v-model="content"
            placeholder="Write your story..."
            rows="10"
            required
            class="form-textarea"
          ></textarea>
        </div>

        <button
          type="submit"
          :disabled="loading || !title || !content"
          class="submit-button"
        >
          {{ loading ? 'Submitting...' : 'Submit Post' }}
        </button>

        <p class="disclaimer">
          Your post is anonymous. Your country will be detected from your IP address.
        </p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { usePostsStore } from '../stores/posts'

const postsStore = usePostsStore()

const title = ref('')
const content = ref('')
const loading = ref(false)
const error = ref(null)
const success = ref(false)

const handleSubmit = async () => {
  error.value = null
  success.value = false
  loading.value = true

  try {
    await postsStore.createPost({
      title: title.value,
      content: content.value
    })

    success.value = true
    title.value = ''
    content.value = ''
  } catch (err) {
    error.value = postsStore.error
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.submit-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  min-height: calc(100vh - 80px);
}

.submit-form-wrapper {
  background-color: #2a2a2a;
  border-radius: 8px;
  padding: 2rem;
  max-width: 600px;
  width: 100%;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
}

.form-title {
  color: #e0e0e0;
  margin-bottom: 1.5rem;
  font-size: 1.75rem;
  text-align: center;
}

.success-message {
  background-color: #2d4a2d;
  color: #90ee90;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  text-align: center;
}

.browse-link {
  color: #90ee90;
  text-decoration: underline;
  font-weight: bold;
}

.error-message {
  background-color: #4a2d2d;
  color: #ff9090;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  text-align: center;
}

.submit-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  position: relative;
}

label {
  color: #b0b0b0;
  font-weight: 500;
}

.form-input,
.form-textarea {
  background-color: #1a1a1a;
  border: 1px solid #3a3a3a;
  color: #e0e0e0;
  padding: 0.75rem;
  border-radius: 4px;
  font-family: inherit;
  transition: border-color 0.3s;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #5a5a5a;
}

.form-textarea {
  resize: vertical;
  min-height: 150px;
}

.char-count {
  color: #7a7a7a;
  font-size: 0.875rem;
  text-align: right;
}

.submit-button {
  background-color: #4a4a4a;
  color: #e0e0e0;
  border: none;
  padding: 1rem;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submit-button:hover:not(:disabled) {
  background-color: #5a5a5a;
}

.submit-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.disclaimer {
  color: #7a7a7a;
  font-size: 0.875rem;
  text-align: center;
  margin: 0;
}
</style>
