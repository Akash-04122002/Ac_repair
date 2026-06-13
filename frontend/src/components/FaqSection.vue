<template>
  <section id="faq" class="py-12">
    <v-container>
      <h2 class="text-h4 font-weight-bold text-center mb-8">Frequently asked questions</h2>
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-expansion-panels variant="accordion">
            <v-expansion-panel v-for="(f, i) in faq" :key="i" :title="f.q" :text="f.a" />
          </v-expansion-panels>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
defineProps({ faq: { type: Array, default: () => [] } })
</script>
