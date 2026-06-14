<template>
  <section id="services" class="py-14 services-bg">
    <v-container>
      <div class="text-center mb-10">
        <v-chip color="primary" variant="tonal" size="small" class="mb-3 font-weight-bold">OUR SERVICES</v-chip>
        <h2 class="text-h4 text-md-h3 font-weight-black mb-2">What can we fix for you?</h2>
        <p class="text-subtitle-1 text-medium-emphasis mx-auto" style="max-width:520px">
          Transparent price ranges · final estimate confirmed before work starts
        </p>
      </div>

      <!-- category tabs with photos -->
      <div class="cat-tabs mb-8 d-flex justify-center ga-3 flex-wrap">
        <div
          v-for="c in categories"
          :key="c.id"
          class="cat-tab"
          :class="{ active: tab === c.id }"
          @click="tab = c.id"
        >
          <div class="cat-tab-img-wrap">
            <img :src="c.image" :alt="c.name" class="cat-tab-img" loading="lazy" />
            <div class="cat-tab-overlay" />
          </div>
          <v-icon :icon="c.icon" class="cat-tab-icon" />
          <span class="cat-tab-label">{{ c.name }}</span>
        </div>
      </div>

      <!-- service cards -->
      <v-window v-model="tab">
        <v-window-item v-for="c in categories" :key="c.id" :value="c.id">
          <v-row>
            <v-col v-for="(s, i) in c.services" :key="i" cols="12" sm="6" md="4">
              <v-card class="service-card h-100 d-flex flex-column" rounded="xl" elevation="0">
                <!-- photo -->
                <div class="service-img-wrap">
                  <img :src="s.image" :alt="s.name" class="service-img" loading="lazy" />
                  <div class="service-img-overlay" />
                  <v-chip color="secondary" size="small" class="price-chip font-weight-bold">
                    {{ s.price }}
                  </v-chip>
                </div>

                <v-card-item class="pt-3 pb-1">
                  <template #prepend>
                    <v-avatar color="primary" variant="tonal" size="36">
                      <v-icon :icon="s.icon || c.icon" size="20" />
                    </v-avatar>
                  </template>
                  <v-card-title class="text-body-1 font-weight-bold text-wrap">{{ s.name }}</v-card-title>
                </v-card-item>

                <v-spacer />

                <v-card-actions class="px-4 pb-4">
                  <v-btn
                    block
                    color="primary"
                    variant="flat"
                    rounded="lg"
                    prepend-icon="mdi-whatsapp"
                    @click="$emit('book', c.id, s.name)"
                  >Book this service</v-btn>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </v-window-item>
      </v-window>
    </v-container>
  </section>
</template>

<script setup>
import { ref, watch } from 'vue'
const props = defineProps({ categories: { type: Array, default: () => [] } })
defineEmits(['book'])
const tab = ref(null)
watch(() => props.categories, (c) => { if (c.length && !tab.value) tab.value = c[0].id }, { immediate: true })
</script>

<style scoped>
.services-bg {
  background: #F0F4FA;
}

/* category selector tabs */
.cat-tab {
  position: relative;
  width: 140px;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  border: 3px solid transparent;
  transition: border-color .2s, transform .2s, box-shadow .2s;
  box-shadow: 0 2px 8px rgba(0,0,0,.08);
}
.cat-tab:hover { transform: translateY(-3px); box-shadow: 0 6px 20px rgba(0,0,0,.13); }
.cat-tab.active { border-color: #1565C0; box-shadow: 0 4px 20px rgba(21,101,192,.3); }

.cat-tab-img-wrap { position: relative; height: 88px; }
.cat-tab-img { width: 100%; height: 100%; object-fit: cover; display: block; }
.cat-tab-overlay {
  position: absolute; inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,.6) 0%, rgba(0,0,0,.1) 100%);
}
.cat-tab-icon {
  position: absolute;
  top: 50%; left: 50%;
  transform: translate(-50%, -60%);
  color: white !important;
  font-size: 28px !important;
}
.cat-tab-label {
  display: block;
  text-align: center;
  font-size: .8rem;
  font-weight: 700;
  padding: 6px 4px;
  background: white;
  color: #333;
}
.cat-tab.active .cat-tab-label { color: #1565C0; }

/* service photo card */
.service-card {
  background: white;
  transition: transform .2s, box-shadow .2s;
  border: 1px solid rgba(0,0,0,.06);
}
.service-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0,0,0,.12) !important;
}

.service-img-wrap {
  position: relative;
  height: 180px;
  overflow: hidden;
  border-radius: 12px 12px 0 0;
}
.service-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform .4s ease;
}
.service-card:hover .service-img { transform: scale(1.06); }
.service-img-overlay {
  position: absolute; inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,.3) 0%, transparent 60%);
}
.price-chip {
  position: absolute;
  bottom: 10px;
  right: 10px;
  font-size: .75rem !important;
}
</style>
