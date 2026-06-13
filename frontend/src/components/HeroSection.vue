<template>
  <section class="hero">
    <v-container class="py-12 py-md-16">
      <v-row align="center">
        <v-col cols="12" md="7">
          <v-chip color="secondary" variant="flat" size="small" class="mb-4">
            <v-icon start icon="mdi-map-marker" />{{ business.area }} · {{ business.hours }}
          </v-chip>
          <h1 class="text-h3 text-md-h2 font-weight-bold text-white mb-4">
            Electrician & Plumber,<br />at your door in Chennai
          </h1>
          <p class="text-h6 font-weight-regular text-white mb-8" style="opacity: 0.92">
            {{ business.tagline || 'Trusted, verified professionals. Book in seconds on WhatsApp.' }}
          </p>

          <div class="d-flex flex-wrap ga-3">
            <v-btn size="x-large" color="secondary" prepend-icon="mdi-calendar-check" @click="$emit('book')">
              Book Now
            </v-btn>
            <v-btn
              size="x-large"
              variant="flat"
              color="white"
              :href="`tel:${business.phone}`"
              prepend-icon="mdi-phone"
            >
              Call Us
            </v-btn>
          </div>

          <div class="d-flex flex-wrap ga-6 mt-8 text-white">
            <div><div class="text-h5 font-weight-bold">4.8★</div><div class="text-caption">Avg. rating</div></div>
            <div><div class="text-h5 font-weight-bold">2 hrs</div><div class="text-caption">Avg. arrival</div></div>
            <div><div class="text-h5 font-weight-bold">100%</div><div class="text-caption">All Chennai</div></div>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
defineProps({ business: { type: Object, required: true } })
defineEmits(['book'])
</script>

<style scoped>
.hero {
  background: linear-gradient(135deg, #0D47A1 0%, #1565C0 55%, #1976D2 100%);
}
</style>
