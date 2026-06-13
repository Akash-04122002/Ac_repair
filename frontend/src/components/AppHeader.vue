<template>
  <v-app-bar color="surface" elevation="2" scroll-behavior="elevate">
    <v-app-bar-title>
      <span class="font-weight-bold text-primary">QuickFix</span>
      <span class="font-weight-bold"> Chennai</span>
    </v-app-bar-title>

    <template #append>
      <v-btn
        :href="`tel:${business.phone}`"
        variant="text"
        prepend-icon="mdi-phone"
        class="d-none d-sm-flex"
      >
        Call
      </v-btn>
      <v-btn color="secondary" variant="flat" prepend-icon="mdi-calendar-check" @click="$emit('book')">
        Book Now
      </v-btn>
    </template>
  </v-app-bar>
</template>

<script setup>
defineProps({ business: { type: Object, required: true } })
defineEmits(['book'])
</script>
