<template>
  <section class="hero">
    <!-- background photo with dark overlay -->
    <div class="hero-bg" />

    <v-container class="hero-content py-16 py-md-24">
      <v-row align="center" style="min-height: 520px">
        <v-col cols="12" md="7">
          <v-chip color="white" variant="flat" size="small" class="mb-5 chip-pill" style="color:#1565C0; font-weight:700;">
            <v-icon start icon="mdi-map-marker" color="#1565C0" size="16" />
            {{ business.area }} &nbsp;·&nbsp; {{ business.hours }}
          </v-chip>

          <h1 class="text-h3 text-md-h2 font-weight-black text-white mb-4 hero-title">
            Fast Electrician &amp;<br />Plumber in Chennai
          </h1>
          <p class="text-h6 font-weight-regular text-white mb-8 hero-sub">
            Verified professionals at your door — booked in seconds on WhatsApp.
          </p>

          <div class="d-flex flex-wrap ga-3 mb-10">
            <v-btn
              size="x-large"
              color="secondary"
              class="hero-cta"
              prepend-icon="mdi-calendar-check"
              href="#services"
            >Book Now — It's Free</v-btn>
            <v-btn
              size="x-large"
              variant="outlined"
              color="white"
              :href="`tel:${business.phone}`"
              prepend-icon="mdi-phone"
              class="hero-call"
            >{{ business.phone }}</v-btn>
          </div>

          <!-- stats row -->
          <div class="d-flex flex-wrap ga-6">
            <div v-for="s in stats" :key="s.label" class="stat-item text-white">
              <div class="text-h5 font-weight-black">{{ s.value }}</div>
              <div class="text-caption" style="opacity:.8">{{ s.label }}</div>
            </div>
          </div>
        </v-col>

        <!-- trust badge card -->
        <v-col cols="12" md="5" class="d-none d-md-flex justify-end">
          <v-card class="trust-card pa-5" rounded="xl" elevation="8">
            <div class="text-subtitle-2 text-medium-emphasis mb-4 font-weight-bold text-uppercase" style="letter-spacing:1px">Why customers trust us</div>
            <v-list density="compact" bg-color="transparent">
              <v-list-item
                v-for="t in trust"
                :key="t.text"
                :prepend-icon="t.icon"
                :title="t.text"
                :base-color="t.color"
                class="px-0"
              />
            </v-list>
            <v-divider class="my-4" />
            <div class="d-flex align-center ga-2">
              <v-rating :model-value="4.8" color="amber" density="compact" size="small" readonly half-increments />
              <span class="text-body-2 font-weight-bold">4.8 / 5</span>
              <span class="text-caption text-medium-emphasis">· 200+ reviews</span>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
defineProps({ business: { type: Object, required: true } })
defineEmits(['book'])

const stats = [
  { value: '2 hrs', label: 'Avg. arrival' },
  { value: '24 x 7', label: 'Available' },
  { value: '100+', label: 'Jobs done' },
  { value: '90%', label: 'All Chennai' },
]
const trust = [
  { icon: 'mdi-shield-check', text: 'Background-verified technicians', color: 'success' },
  { icon: 'mdi-cash-check', text: 'Transparent pricing, no surprises', color: 'primary' },
  { icon: 'mdi-clock-fast', text: 'Same-day service available', color: 'secondary' },
  { icon: 'mdi-medal', text: 'Service warranty on every job', color: 'warning' },
]
</script>

<style scoped>
.hero {
  position: relative;
  overflow: hidden;
  background: #0a2550;
}
.hero-bg {
  position: absolute;
  inset: 0;
  background-image: url('https://images.unsplash.com/photo-1621905251918-48416bd8575a?w=1400&q=60&auto=format&fit=crop');
  background-size: cover;
  background-position: center;
  filter: brightness(0.32);
}
.hero-content {
  position: relative;
  z-index: 1;
}
.hero-title {
  line-height: 1.12;
  text-shadow: 0 2px 12px rgba(0,0,0,.4);
}
.hero-sub {
  opacity: .93;
  max-width: 520px;
}
.hero-cta {
  font-size: 1rem !important;
  box-shadow: 0 4px 20px rgba(251,140,0,.5);
}
.hero-call {
  border-color: rgba(255,255,255,.6) !important;
}
.chip-pill {
  border-radius: 999px !important;
}
.stat-item {
  min-width: 72px;
}
.trust-card {
  background: rgba(255,255,255,0.97) !important;
  max-width: 340px;
  width: 100%;
}
</style>
