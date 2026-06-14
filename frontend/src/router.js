import { createRouter, createWebHistory } from 'vue-router'
import BookingStatus from './pages/BookingStatus.vue'
import OwnerAction from './pages/OwnerAction.vue'
import UserReply from './pages/UserReply.vue'

// Lazy-load the heavy landing page
const Landing = () => import('./Landing.vue')

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',             component: Landing },
    { path: '/booking/:id', component: BookingStatus },
    { path: '/owner/:id',   component: OwnerAction },
    { path: '/reply/:id',   component: UserReply },
  ],
  scrollBehavior(to) {
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    return { top: 0 }
  },
})
