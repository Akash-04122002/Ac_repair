<template>
  <section id="faq" class="py-14" style="background:#fff">
    <v-container>
      <div class="text-center mb-10">
        <v-chip color="primary" variant="tonal" size="small" class="mb-3 font-weight-bold">FAQ</v-chip>
        <h2 class="text-h4 font-weight-black mb-2">Got questions? We've got answers</h2>
        <p class="text-subtitle-1 text-medium-emphasis">Everything you need to know before booking</p>
      </div>

      <v-row justify="center">
        <v-col cols="12" md="8">
          <div class="faq-list">
            <div
              v-for="(f, i) in faq"
              :key="i"
              class="faq-item"
              :class="{ open: openIdx === i }"
              @click="openIdx = openIdx === i ? null : i"
            >
              <div class="faq-question">
                <v-icon icon="mdi-help-circle-outline" color="primary" class="mr-3" size="22" />
                <span>{{ f.q }}</span>
                <v-icon
                  :icon="openIdx === i ? 'mdi-chevron-up' : 'mdi-chevron-down'"
                  class="ml-auto faq-chevron"
                  :class="{ rotated: openIdx === i }"
                />
              </div>
              <div class="faq-answer" :style="openIdx === i ? 'max-height:300px; opacity:1' : ''">
                <p class="text-body-2 text-medium-emphasis mb-0">{{ f.a }}</p>
              </div>
            </div>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </section>
</template>

<script setup>
import { ref } from 'vue'
defineProps({ faq: { type: Array, default: () => [] } })
const openIdx = ref(0)
</script>

<style scoped>
.faq-list { display: flex; flex-direction: column; gap: 12px; }

.faq-item {
  border-radius: 14px;
  border: 1.5px solid rgba(21,101,192,.12);
  background: #F8FAFF;
  overflow: hidden;
  cursor: pointer;
  transition: border-color .2s, box-shadow .2s;
}
.faq-item:hover { border-color: rgba(21,101,192,.3); box-shadow: 0 4px 16px rgba(21,101,192,.08); }
.faq-item.open { border-color: #1565C0; box-shadow: 0 4px 20px rgba(21,101,192,.15); background: white; }

.faq-question {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  font-weight: 600;
  font-size: .95rem;
  color: #1a1a2e;
  user-select: none;
}
.faq-chevron { color: #1565C0; transition: transform .25s; }
.faq-chevron.rotated { transform: rotate(180deg); }

.faq-answer {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
  padding: 0 20px;
  transition: max-height .3s ease, opacity .3s ease, padding .3s ease;
}
.faq-item.open .faq-answer {
  padding: 0 20px 18px 20px;
}
</style>
