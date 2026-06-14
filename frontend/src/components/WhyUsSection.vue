<template>
  <section id="why" class="py-14 why-bg">
    <v-container>
      <div class="text-center mb-10">
        <v-chip color="primary" variant="tonal" size="small" class="mb-3 font-weight-bold">WHY CHOOSE US</v-chip>
        <h2 class="text-h4 font-weight-black mb-2">Built on trust, every time</h2>
        <p class="text-subtitle-1 text-medium-emphasis">Thousands of Chennai homes serviced since 2019</p>
      </div>

      <!-- feature tiles -->
      <v-row class="mb-12">
        <v-col v-for="(item, i) in items" :key="i" cols="6" md="3">
          <div class="feature-tile text-center pa-5">
            <div class="feature-icon-wrap mb-3">
              <v-icon :icon="item.icon" size="32" color="white" />
            </div>
            <h3 class="text-subtitle-1 font-weight-bold mb-1">{{ item.title }}</h3>
            <p class="text-body-2 text-medium-emphasis">{{ item.text }}</p>
          </div>
        </v-col>
      </v-row>

      <!-- reviews -->
      <div class="text-center mb-6">
        <h3 class="text-h5 font-weight-black">What our customers say</h3>
        <div class="d-flex align-center justify-center ga-1 mt-1">
          <v-rating :model-value="4.8" color="amber" density="compact" size="small" readonly half-increments />
          <span class="text-body-2 font-weight-bold ml-1">4.8 / 5</span>
          <span class="text-caption text-medium-emphasis">· 200+ verified reviews</span>
        </div>
      </div>

      <v-row>
        <v-col v-for="(r, i) in reviews" :key="i" cols="12" sm="6" md="3">
          <v-card class="review-card h-100" rounded="xl" elevation="0">
            <v-card-text class="pa-5">
              <v-rating :model-value="r.rating" color="amber" density="compact" size="small" readonly half-increments class="mb-3" />
              <p class="text-body-2 review-text mb-4">"{{ r.text }}"</p>
              <div class="d-flex align-center ga-2">
                <v-avatar color="primary" size="36">
                  <span class="text-caption font-weight-bold text-white">{{ r.name[0] }}</span>
                </v-avatar>
                <div>
                  <div class="text-subtitle-2 font-weight-bold">{{ r.name }}</div>
                  <div class="text-caption text-medium-emphasis">
                    <v-icon icon="mdi-map-marker" size="12" />{{ r.area }}
                  </div>
                </div>
              </div>
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

<style scoped>
.why-bg { background: #fff; }

.feature-tile {
  border-radius: 20px;
  background: #F8FAFF;
  border: 1px solid rgba(21,101,192,.08);
  transition: transform .2s, box-shadow .2s;
  height: 100%;
}
.feature-tile:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(21,101,192,.12);
}
.feature-icon-wrap {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  background: linear-gradient(135deg, #1565C0, #1976D2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  box-shadow: 0 4px 16px rgba(21,101,192,.3);
}

.review-card {
  background: #F8FAFF;
  border: 1px solid rgba(21,101,192,.08);
  transition: transform .2s, box-shadow .2s;
}
.review-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,.10);
}
.review-text {
  line-height: 1.6;
  color: #444;
}
</style>
