<template>
  <section id="areas" class="py-14 areas-bg">
    <v-container>
      <div class="text-center mb-10">
        <v-chip color="white" variant="flat" size="small" class="mb-3 font-weight-bold" style="color:#1565C0">
          <v-icon start icon="mdi-map-marker" color="#1565C0" />SERVICE AREAS
        </v-chip>
        <h2 class="text-h4 font-weight-black text-white mb-2">We cover all of Chennai</h2>
        <p class="text-subtitle-1 text-white mb-0" style="opacity:.8">
          Fast, reliable service across every major locality
        </p>
      </div>

      <div class="d-flex flex-wrap justify-center ga-3 mb-8">
        <div
          v-for="(a, i) in areas"
          :key="i"
          class="area-chip"
        >
          <v-icon icon="mdi-map-marker" size="15" class="mr-1" />{{ a }}
        </div>
      </div>

      <div class="text-center">
        <p class="text-white mb-4" style="opacity:.85">Not in the list? We likely cover your area too — just ask!</p>
        <v-btn
          color="secondary"
          size="large"
          variant="flat"
          prepend-icon="mdi-whatsapp"
          href="https://wa.me/919043050841"
          target="_blank"
        >Check your area on WhatsApp</v-btn>
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

<style scoped>
.areas-bg {
  background: linear-gradient(135deg, #0D47A1 0%, #1565C0 50%, #1976D2 100%);
}
.area-chip {
  display: inline-flex;
  align-items: center;
  background: rgba(255,255,255,.15);
  color: white;
  font-size: .875rem;
  font-weight: 600;
  padding: 8px 16px;
  border-radius: 999px;
  border: 1px solid rgba(255,255,255,.25);
  cursor: default;
  transition: background .2s, transform .15s;
  backdrop-filter: blur(4px);
}
.area-chip:hover {
  background: rgba(255,255,255,.28);
  transform: translateY(-2px);
}
</style>
