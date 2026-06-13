<template>
  <section id="areas" class="py-12" style="background: #EEF2F7">
    <v-container>
      <h2 class="text-h4 font-weight-bold text-center mb-2">We serve {{ areaLabel }}</h2>
      <p class="text-subtitle-1 text-medium-emphasis text-center mb-8">
        Fast service across every major Chennai locality
      </p>

      <div class="d-flex flex-wrap justify-center ga-2">
        <v-chip v-for="(a, i) in areas" :key="i" color="primary" variant="outlined" prepend-icon="mdi-map-marker">
          {{ a }}
        </v-chip>
      </div>
    </v-container>
  </section>
</template>

<script setup>
defineProps({
  areas: { type: Array, default: () => [] },
  areaLabel: { type: String, default: 'all over Chennai' },
})
</script>
