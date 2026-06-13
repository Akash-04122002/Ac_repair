<template>
  <v-app>
    <AppHeader :business="config.business" @book="openBooking()" />

    <v-main>
      <HeroSection :business="config.business" @book="openBooking()" />
      <HowItWorks />
      <ServicesSection :categories="config.categories" @book="openBooking" />
      <WhyUsSection :items="config.whyUs" :reviews="config.reviews" />
      <AreasSection :areas="config.areas" :area-label="config.business.area" />
      <FaqSection :faq="config.faq" />
      <ContactSection :business="config.business" @book="openBooking()" />
      <AppFooter :business="config.business" />
    </v-main>

    <FloatingActions :business="config.business" @emergency="openBooking(null, null, true)" />

    <BookingDialog
      v-model="bookingOpen"
      :categories="config.categories"
      :business="config.business"
      :prefill="prefill"
      :urgent="urgent"
    />
  </v-app>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { fetchConfig, FALLBACK_CONFIG } from './api'
import AppHeader from './components/AppHeader.vue'
import HeroSection from './components/HeroSection.vue'
import HowItWorks from './components/HowItWorks.vue'
import ServicesSection from './components/ServicesSection.vue'
import WhyUsSection from './components/WhyUsSection.vue'
import AreasSection from './components/AreasSection.vue'
import FaqSection from './components/FaqSection.vue'
import ContactSection from './components/ContactSection.vue'
import AppFooter from './components/AppFooter.vue'
import FloatingActions from './components/FloatingActions.vue'
import BookingDialog from './components/BookingDialog.vue'

const config = ref(FALLBACK_CONFIG)
const bookingOpen = ref(false)
const urgent = ref(false)
const prefill = ref({ categoryId: null, service: null })

onMounted(async () => {
  config.value = await fetchConfig()
})

function openBooking(categoryId = null, service = null, isUrgent = false) {
  prefill.value = { categoryId, service }
  urgent.value = isUrgent
  bookingOpen.value = true
}
</script>
