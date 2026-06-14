<template>
  <section id="contact" class="py-14 contact-bg">
    <v-container>
      <v-row align="center" justify="center">
        <v-col cols="12" md="6" class="text-center text-md-left">
          <v-chip color="white" variant="flat" size="small" class="mb-4 font-weight-bold" style="color:#0D47A1">
            GET IN TOUCH
          </v-chip>
          <h2 class="text-h4 text-md-h3 font-weight-black text-white mb-3">
            Ready to get it fixed?
          </h2>
          <p class="text-h6 font-weight-regular text-white mb-8" style="opacity:.88">
            Book now or call us — we're available Mon–Sun, 8 AM–8 PM.
          </p>
          <div class="d-flex flex-wrap ga-3 justify-center justify-md-start">
            <v-btn
              size="x-large"
              color="secondary"
              variant="flat"
              prepend-icon="mdi-calendar-check"
              class="contact-cta"
              @click="$emit('book')"
            >Book an appointment</v-btn>
            <v-btn
              size="x-large"
              variant="outlined"
              color="white"
              :href="`tel:${business.phone}`"
              prepend-icon="mdi-phone"
              style="border-color:rgba(255,255,255,.5)!important"
            >Call now</v-btn>
          </div>
        </v-col>

        <v-col cols="12" md="5" offset-md="1">
          <v-card class="contact-card pa-2" rounded="xl" elevation="0">
            <v-list lines="two" bg-color="transparent">
              <v-list-item
                v-for="item in contactItems"
                :key="item.subtitle"
                :prepend-icon="item.icon"
                :title="item.title"
                :subtitle="item.subtitle"
                :href="item.href"
                :target="item.target"
                :base-color="item.color"
                rounded="lg"
                class="mb-1"
              />
            </v-list>
            <div class="pa-3 pt-1">
              <v-btn
                block
                color="#25D366"
                variant="flat"
                size="large"
                rounded="lg"
                prepend-icon="mdi-whatsapp"
                :href="waLink"
                target="_blank"
              >
                Chat on WhatsApp
              </v-btn>
            </div>
          </v-card>
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

const contactItems = computed(() => [
  { icon: 'mdi-phone', title: props.business.phone, subtitle: 'Tap to call us directly', href: `tel:${props.business.phone}`, color: 'primary' },
  { icon: 'mdi-whatsapp', title: props.business.phone, subtitle: 'Message or book on WhatsApp', href: waLink.value, target: '_blank', color: '#25D366' },
  { icon: 'mdi-clock-outline', title: props.business.hours, subtitle: 'Working hours', color: 'secondary' },
  { icon: 'mdi-map-marker', title: props.business.area, subtitle: 'Service coverage', color: 'error' },
])
</script>

<style scoped>
.contact-bg {
  background: linear-gradient(135deg, #0a1f4e 0%, #0D47A1 55%, #1565C0 100%);
}
.contact-cta {
  box-shadow: 0 4px 20px rgba(251,140,0,.5);
}
.contact-card {
  background: rgba(255,255,255,.97) !important;
}
</style>
