<template>
  <!-- Floating WhatsApp chat button (always visible) -->
  <v-btn
    :href="waLink"
    target="_blank"
    icon="mdi-whatsapp"
    color="#25D366"
    size="large"
    elevation="6"
    class="fab-whatsapp"
  />

  <!-- Sticky bottom bar on mobile: quick Call + Emergency -->
  <div class="mobile-bar d-flex d-sm-none">
    <v-btn :href="`tel:${business.phone}`" color="primary" variant="flat" class="flex-1-1" prepend-icon="mdi-phone" rounded="0">
      Call
    </v-btn>
    <v-btn color="error" variant="flat" class="flex-1-1" prepend-icon="mdi-alarm-light" rounded="0" @click="$emit('emergency')">
      Emergency
    </v-btn>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ business: { type: Object, required: true } })
defineEmits(['emergency'])
const waLink = computed(
  () => `https://wa.me/${props.business.whatsapp}?text=${encodeURIComponent('Hi, I have a question about your service.')}`
)
</script>

<style scoped>
.fab-whatsapp {
  position: fixed;
  right: 16px;
  bottom: 80px;
  z-index: 1000;
}
@media (min-width: 600px) {
  .fab-whatsapp { bottom: 24px; }
}
.mobile-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
}
.flex-1-1 { flex: 1 1 0; }
</style>
