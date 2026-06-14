<template>
  <v-app-bar color="surface" :elevation="scrolled ? 3 : 0" scroll-behavior="elevate" class="header-bar">
    <!-- logo -->
    <v-app-bar-title>
      <div class="d-flex align-center ga-2">
        <div class="logo-icon-wrap">
          <v-icon icon="mdi-flash" color="white" size="18" />
        </div>
        <span>
          <span class="font-weight-black text-primary" style="font-size:1.1rem">QuickFix</span>
          <span class="font-weight-bold" style="font-size:1.1rem"> Chennai</span>
        </span>
      </div>
    </v-app-bar-title>

    <!-- desktop nav -->
    <template #append>
      <div class="d-none d-md-flex align-center ga-1 mr-2">
        <v-btn variant="text" size="small" href="#services">Services</v-btn>
        <v-btn variant="text" size="small" href="#how">How it works</v-btn>
        <v-btn variant="text" size="small" href="#areas">Areas</v-btn>
        <v-btn variant="text" size="small" href="#faq">FAQ</v-btn>
      </div>
      <v-btn
        :href="`tel:${business.phone}`"
        variant="text"
        prepend-icon="mdi-phone"
        size="small"
        class="d-none d-sm-flex mr-1"
      >Call</v-btn>
      <v-btn
        color="secondary"
        variant="flat"
        prepend-icon="mdi-calendar-check"
        size="small"
        class="mr-2"
        @click="$emit('book')"
      >Book Now</v-btn>
    </template>
  </v-app-bar>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
defineProps({ business: { type: Object, required: true } })
defineEmits(['book'])

const scrolled = ref(false)
function onScroll() { scrolled.value = window.scrollY > 20 }
onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<style scoped>
.header-bar {
  backdrop-filter: blur(12px);
  background: rgba(255,255,255,0.95) !important;
}
.logo-icon-wrap {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, #1565C0, #1976D2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
</style>
