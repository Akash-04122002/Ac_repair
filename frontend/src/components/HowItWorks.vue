<template>
  <section id="how" class="py-12">
    <v-container>
      <h2 class="text-h4 font-weight-bold text-center mb-2">How it works</h2>
      <p class="text-subtitle-1 text-medium-emphasis text-center mb-10">Booking takes under a minute</p>

      <v-row justify="center">
        <v-col v-for="(s, i) in steps" :key="i" cols="12" sm="4">
          <div class="text-center px-4">
            <v-avatar color="primary" size="64" class="mb-4">
              <v-icon :icon="s.icon" size="32" />
            </v-avatar>
            <div class="text-overline text-secondary">Step {{ i + 1 }}</div>
            <h3 class="text-h6 font-weight-bold mb-2">{{ s.title }}</h3>
            <p class="text-body-2 text-medium-emphasis">{{ s.text }}</p>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
const steps = [
  { icon: 'mdi-format-list-checks', title: 'Pick a service', text: 'Choose electrical, plumbing or an annual plan and your preferred time.' },
  { icon: 'mdi-whatsapp', title: 'Book on WhatsApp', text: 'Tap one button — your booking lands in our WhatsApp instantly.' },
  { icon: 'mdi-home-heart', title: 'We arrive', text: 'A verified professional reaches you and gets the job done.' },
]
</script>
