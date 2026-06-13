<template>
  <section id="contact" class="py-12" style="background: #0D47A1">
    <v-container>
      <v-row align="center">
        <v-col cols="12" md="7">
          <h2 class="text-h4 font-weight-bold text-white mb-4">Ready to book?</h2>
          <v-list bg-color="transparent" class="text-white">
            <v-list-item prepend-icon="mdi-phone" :title="business.phone" subtitle="Tap to call" :href="`tel:${business.phone}`" />
            <v-list-item prepend-icon="mdi-whatsapp" :title="business.phone" subtitle="Chat / book on WhatsApp" :href="waLink" target="_blank" />
            <v-list-item prepend-icon="mdi-clock-outline" :title="business.hours" subtitle="Working hours" />
            <v-list-item prepend-icon="mdi-map-marker" :title="business.area" subtitle="Service area" />
          </v-list>
        </v-col>
        <v-col cols="12" md="5" class="text-center text-md-right">
          <v-btn size="x-large" color="secondary" prepend-icon="mdi-calendar-check" @click="$emit('book')">
            Book an appointment
          </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ business: { type: Object, required: true } })
defineEmits(['book'])
const waLink = computed(() => `https://wa.me/${props.business.whatsapp}`)
</script>
