<template>
  <section id="why" class="py-12">
    <v-container>
      <h2 class="text-h4 font-weight-bold text-center mb-8">Why choose us</h2>

      <v-row class="mb-8">
        <v-col v-for="(item, i) in items" :key="i" cols="6" md="3">
          <div class="text-center px-2">
            <v-icon :icon="item.icon" color="primary" size="40" class="mb-3" />
            <h3 class="text-subtitle-1 font-weight-bold">{{ item.title }}</h3>
            <p class="text-body-2 text-medium-emphasis">{{ item.text }}</p>
          </div>
        </v-col>
      </v-row>

      <h3 v-if="reviews.length" class="text-h5 font-weight-bold text-center mb-6">What customers say</h3>
      <v-row>
        <v-col v-for="(r, i) in reviews" :key="i" cols="12" sm="6" md="3">
          <v-card variant="tonal" color="primary" class="h-100">
            <v-card-text>
              <v-rating :model-value="r.rating" color="amber" density="compact" size="small" readonly half-increments />
              <p class="text-body-2 my-3">“{{ r.text }}”</p>
              <div class="text-subtitle-2 font-weight-bold">{{ r.name }}</div>
              <div class="text-caption text-medium-emphasis">{{ r.area }}</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
defineProps({
  items: { type: Array, default: () => [] },
  reviews: { type: Array, default: () => [] },
})
</script>
