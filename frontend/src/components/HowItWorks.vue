<template>
  <section id="how" class="py-14" style="background:#fff">
    <v-container>
      <div class="text-center mb-10">
        <v-chip color="secondary" variant="tonal" size="small" class="mb-3 font-weight-bold">HOW IT WORKS</v-chip>
        <h2 class="text-h4 font-weight-black mb-2">Book in under a minute</h2>
        <p class="text-subtitle-1 text-medium-emphasis">No calls, no waiting — just WhatsApp</p>
      </div>

      <v-row justify="center" align="stretch">
        <v-col v-for="(s, i) in steps" :key="i" cols="12" sm="4">
          <div class="step-card text-center px-4 py-6 h-100">
            <!-- connector line (shown between steps on desktop) -->
            <div class="step-number-wrap mb-4">
              <div class="step-number">{{ i + 1 }}</div>
              <div v-if="i < steps.length - 1" class="step-connector d-none d-sm-block" />
            </div>
            <v-avatar :color="s.color" size="72" class="mb-4 step-avatar">
              <v-icon :icon="s.icon" size="36" color="white" />
            </v-avatar>
            <h3 class="text-h6 font-weight-bold mb-2">{{ s.title }}</h3>
            <p class="text-body-2 text-medium-emphasis">{{ s.text }}</p>
          </div>
        </v-col>
      </v-row>

      <div class="text-center mt-8">
        <v-btn color="secondary" size="large" variant="flat" prepend-icon="mdi-whatsapp" href="https://wa.me/919043050841" target="_blank">
          Chat with us on WhatsApp
        </v-btn>
      </div>
    </v-container>
  </section>
</template>

<script setup>
const steps = [
  {
    icon: 'mdi-format-list-checks',
    title: 'Pick a service',
    text: 'Choose from electrical, plumbing, or annual maintenance plans and select your preferred time slot.',
    color: 'primary',
  },
  {
    icon: 'mdi-whatsapp',
    title: 'Confirm on WhatsApp',
    text: 'Tap one button — your full booking details land in our WhatsApp instantly. No app needed.',
    color: '#25D366',
  },
  {
    icon: 'mdi-home-heart',
    title: 'We arrive & fix it',
    text: 'A background-verified professional reaches your door and gets the job done right.',
    color: 'secondary',
  },
]
</script>

<style scoped>
.step-card {
  position: relative;
  border-radius: 20px;
  background: #F8FAFF;
  border: 1px solid rgba(21,101,192,.08);
  transition: box-shadow .2s, transform .2s;
}
.step-card:hover {
  box-shadow: 0 8px 28px rgba(21,101,192,.12);
  transform: translateY(-4px);
}
.step-number-wrap {
  position: relative;
  display: inline-flex;
  align-items: center;
}
.step-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #1565C0;
  color: white;
  font-weight: 800;
  font-size: .9rem;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}
.step-avatar {
  box-shadow: 0 6px 20px rgba(0,0,0,.15);
}
</style>
