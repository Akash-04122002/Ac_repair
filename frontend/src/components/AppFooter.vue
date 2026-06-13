<template>
  <v-footer color="#08306B" class="text-white">
    <v-container>
      <v-row align="center">
        <v-col cols="12" md="6">
          <div class="font-weight-bold text-h6 mb-1">{{ business.name }}</div>
          <div class="text-body-2" style="opacity: 0.8">
            {{ business.area }} · {{ business.hours }}
          </div>
        </v-col>
        <v-col cols="12" md="6" class="text-md-right text-body-2" style="opacity: 0.8">
          {{ business.phone }} · © {{ year }} {{ business.name }}
        </v-col>
      </v-row>
    </v-container>
  </v-footer>
</template>

<script setup>
const props = defineProps({ business: { type: Object, required: true } })
const year = new Date().getFullYear()
</script>
