import { createRouter, createWebHistory } from 'vue-router'
import BrowseView from '../views/BrowseView.vue'
import SubmitView from '../views/SubmitView.vue'
import AdminView from '../views/AdminView.vue'

const routes = [
  {
    path: '/',
    redirect: '/browse'
  },
  {
    path: '/browse',
    name: 'Browse',
    component: BrowseView
  },
  {
    path: '/submit',
    name: 'Submit',
    component: SubmitView
  },
  {
    path: '/admin',
    name: 'Admin',
    component: AdminView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
