<template>
  <div class="reply-page">
    <v-container class="py-10" style="max-width:560px">

      <div v-if="loading" class="text-center py-16">
        <v-progress-circular indeterminate color="primary" size="48" class="mb-4" />
        <p class="text-medium-emphasis">Loading your booking…</p>
      </div>

      <v-card v-else-if="!booking" rounded="xl" elevation="2" class="pa-8 text-center">
        <v-icon icon="mdi-alert-circle-outline" color="error" size="56" class="mb-4" />
        <h2 class="text-h6 font-weight-bold mb-2">Booking not found</h2>
        <v-btn color="primary" href="/">Go to homepage</v-btn>
      </v-card>

      <!-- already confirmed -->
      <v-card v-else-if="booking.status === 'confirmed'" rounded="xl" elevation="2" class="pa-8 text-center">
        <v-icon icon="mdi-check-circle" color="success" size="64" class="mb-4" />
        <h2 class="text-h5 font-weight-bold text-success mb-2">Booking Confirmed!</h2>
        <p class="text-medium-emphasis mb-1">Your appointment is confirmed for:</p>
        <p class="text-body-1 font-weight-bold">{{ booking.date }}, {{ booking.slot }}</p>
        <p class="text-caption text-medium-emphasis mt-3">You'll hear from us on WhatsApp before the visit.</p>
        <v-btn color="primary" class="mt-4" href="/">Back to homepage</v-btn>
      </v-card>

      <template v-else>
        <!-- header -->
        <div class="text-center mb-6">
          <v-icon icon="mdi-calendar-clock" color="primary" size="48" class="mb-3" />
          <h1 class="text-h5 font-weight-black mb-1">Time Proposal</h1>
          <p class="text-body-2 text-medium-emphasis">Hi {{ booking.name }}, we suggested a new time for your booking.</p>
        </div>

        <!-- original vs proposed -->
        <v-card rounded="xl" elevation="1" class="mb-5">
          <v-card-text class="pa-4">
            <v-row dense>
              <v-col cols="6">
                <div class="time-box time-original">
                  <div class="text-caption font-weight-bold mb-1" style="opacity:.7">Your request</div>
                  <div class="text-body-2 font-weight-bold">{{ booking.date }}</div>
                  <div class="text-caption">{{ booking.slot }}</div>
                </div>
              </v-col>
              <v-col cols="6">
                <div class="time-box time-proposed">
                  <div class="text-caption font-weight-bold mb-1">Our proposal</div>
                  <div class="text-body-2 font-weight-bold">{{ booking.proposedDate }}</div>
                  <div class="text-caption">{{ booking.proposedSlot }}</div>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>

        <!-- service summary -->
        <v-card rounded="xl" elevation="1" class="mb-5">
          <v-list density="comfortable" class="px-2">
            <v-list-item prepend-icon="mdi-wrench" :title="booking.category + (booking.service ? ' – ' + booking.service : '')" subtitle="Service" />
            <v-list-item prepend-icon="mdi-map-marker" :title="booking.area" subtitle="Your area" />
          </v-list>
        </v-card>

        <!-- success state -->
        <v-card v-if="done" rounded="xl" :color="doneColor" variant="flat" class="pa-6 text-center mb-5">
          <v-icon :icon="doneIcon" color="white" size="48" class="mb-3" />
          <h3 class="text-h6 font-weight-bold text-white mb-2">{{ doneTitle }}</h3>
          <p class="text-body-2 text-white mb-4" style="opacity:.9">{{ doneSub }}</p>
          <v-btn v-if="whatsappUrl" color="white" :href="whatsappUrl" target="_blank" prepend-icon="mdi-whatsapp" variant="flat"
            style="color:#1565C0!important">
            {{ doneWaLabel }}
          </v-btn>
        </v-card>

        <!-- action buttons -->
        <template v-else>
          <!-- Accept proposed time -->
          <v-card rounded="xl" elevation="1" class="mb-4 action-card" color="success-lighten-5" @click="acceptTime">
            <v-card-text class="pa-5 d-flex align-center ga-4">
              <v-avatar color="success" size="52">
                <v-icon icon="mdi-check" color="white" size="28" />
              </v-avatar>
              <div>
                <div class="text-subtitle-1 font-weight-bold">Accept new time</div>
                <div class="text-body-2 text-medium-emphasis">{{ booking.proposedDate }}, {{ booking.proposedSlot }}</div>
              </div>
              <v-spacer />
              <v-icon icon="mdi-chevron-right" color="success" />
            </v-card-text>
          </v-card>

          <!-- Suggest a different time -->
          <v-card rounded="xl" elevation="1" class="mb-4 action-card" @click="showReschedule = !showReschedule">
            <v-card-text class="pa-5 d-flex align-center ga-4">
              <v-avatar color="secondary" size="52">
                <v-icon icon="mdi-calendar-edit" color="white" size="26" />
              </v-avatar>
              <div>
                <div class="text-subtitle-1 font-weight-bold">Suggest a different time</div>
                <div class="text-body-2 text-medium-emphasis">Send us your preferred date & slot</div>
              </div>
              <v-spacer />
              <v-icon :icon="showReschedule ? 'mdi-chevron-up' : 'mdi-chevron-down'" color="secondary" />
            </v-card-text>

            <v-expand-transition>
              <div v-if="showReschedule" class="px-5 pb-5" @click.stop>
                <v-text-field
                  v-model="newDate"
                  type="date"
                  label="Your preferred date"
                  :min="today"
                  variant="outlined"
                  density="compact"
                  prepend-inner-icon="mdi-calendar"
                  class="mb-3"
                />
                <p class="text-body-2 text-medium-emphasis mb-2">Your preferred time slot</p>
                <v-chip-group v-model="newSlot" selected-class="text-white bg-secondary" column class="mb-4">
                  <v-chip v-for="s in slots" :key="s" :value="s" variant="outlined" filter size="small">{{ s }}</v-chip>
                </v-chip-group>
                <v-btn
                  color="secondary"
                  variant="flat"
                  block
                  :disabled="!newDate || !newSlot"
                  :loading="sending"
                  prepend-icon="mdi-send"
                  @click="requestChange"
                >Send my preferred time</v-btn>
              </div>
            </v-expand-transition>
          </v-card>
        </template>
      </template>
    </v-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const booking = ref(null)
const loading = ref(true)
const showReschedule = ref(false)
const newDate = ref('')
const newSlot = ref('')
const sending = ref(false)
const done = ref(false)
const doneColor = ref('success')
const doneIcon = ref('mdi-check-circle')
const doneTitle = ref('')
const doneSub = ref('')
const doneWaLabel = ref('')
const whatsappUrl = ref('')

const today = new Date().toISOString().slice(0, 10)
const slots = ['8:00 – 10:00 AM', '10:00 – 12:00 PM', '12:00 – 2:00 PM', '2:00 – 4:00 PM', '4:00 – 6:00 PM', '6:00 – 8:00 PM']

async function fetchBooking() {
  try {
    const res = await fetch(`/api/booking/${route.params.id}`)
    booking.value = res.ok ? await res.json() : null
  } catch { booking.value = null }
  finally { loading.value = false }
}
onMounted(fetchBooking)

async function acceptTime() {
  const res = await fetch(`/api/booking/${route.params.id}/user-accept`, { method: 'POST' })
  const data = await res.json()
  whatsappUrl.value = data.waUrl
  doneColor.value = 'success'
  doneIcon.value = 'mdi-check-circle'
  doneTitle.value = 'Time accepted!'
  doneSub.value = 'Your booking is confirmed. We\'ll send a reminder on WhatsApp before the visit.'
  doneWaLabel.value = 'Notify owner on WhatsApp'
  done.value = true
}

async function requestChange() {
  sending.value = true
  try {
    const res = await fetch(`/api/booking/${route.params.id}/user-reschedule`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ date: newDate.value, slot: newSlot.value }),
    })
    const data = await res.json()
    whatsappUrl.value = data.waUrl
    doneColor.value = 'secondary'
    doneIcon.value = 'mdi-calendar-edit'
    doneTitle.value = 'Request sent!'
    doneSub.value = 'We received your preferred time and will get back to you on WhatsApp shortly.'
    doneWaLabel.value = 'Send request on WhatsApp'
    done.value = true
  } finally {
    sending.value = false
  }
}
</script>

<style scoped>
.reply-page { min-height: 100vh; background: #F0F4FA; padding-top: 16px; }
.action-card { cursor: pointer; transition: box-shadow .2s, transform .15s; }
.action-card:hover { transform: translateY(-2px); box-shadow: 0 6px 24px rgba(0,0,0,.10) !important; }
.time-box {
  border-radius: 12px;
  padding: 12px;
  text-align: center;
}
.time-original { background: #f0f0f0; }
.time-proposed { background: #E3F2FD; border: 2px solid #1565C0; }
</style>
