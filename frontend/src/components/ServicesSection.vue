<template>
  <section id="services" class="py-12" style="background: #EEF2F7">
    <v-container>
      <h2 class="text-h4 font-weight-bold text-center mb-2">Our services</h2>
      <p class="text-subtitle-1 text-medium-emphasis text-center mb-8">
        Transparent price ranges · final estimate confirmed before work starts
      </p>

      <v-tabs v-model="tab" color="primary" align-tabs="center" class="mb-6">
        <v-tab v-for="c in categories" :key="c.id" :value="c.id" :prepend-icon="c.icon">
          {{ c.name }}
        </v-tab>
      </v-tabs>

      <v-window v-model="tab">
        <v-window-item v-for="c in categories" :key="c.id" :value="c.id">
          <v-row>
            <v-col v-for="(s, i) in c.services" :key="i" cols="12" sm="6" md="4">
              <v-card variant="elevated" elevation="2" class="h-100 d-flex flex-column">
                <v-card-item>
                  <template #prepend>
                    <v-avatar color="primary" variant="tonal"><v-icon :icon="c.icon" /></v-avatar>
                  </template>
                  <v-card-title class="text-wrap text-body-1 font-weight-bold">{{ s.name }}</v-card-title>
                  <v-card-subtitle class="text-secondary font-weight-bold">{{ s.price }}</v-card-subtitle>
                </v-card-item>
                <v-spacer />
                <v-card-actions>
                  <v-btn block color="primary" variant="tonal" prepend-icon="mdi-whatsapp"
                    @click="$emit('book', c.id, s.name)">
                    Book this
                  </v-btn>
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
watch(
  () => props.categories,
  (c) => { if (c.length && !tab.value) tab.value = c[0].id },
  { immediate: true }
)
</script>
