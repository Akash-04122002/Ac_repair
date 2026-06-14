<template>
  <div class="action-page">
    <v-container class="py-10" style="max-width:560px">

      <div v-if="loading" class="text-center py-16">
        <v-progress-circular indeterminate color="primary" size="48" class="mb-4" />
        <p class="text-medium-emphasis">Loading booking…</p>
      </div>

      <v-card v-else-if="!booking" rounded="xl" elevation="2" class="pa-8 text-center">
        <v-icon icon="mdi-alert-circle-outline" color="error" size="56" class="mb-4" />
        <h2 class="text-h6 font-weight-bold mb-2">Booking not found</h2>
        <v-btn color="primary" href="/">Go to homepage</v-btn>
      </v-card>

      <!-- already confirmed -->
      <v-card v-else-if="booking.status === 'confirmed'" rounded="xl" elevation="2" class="pa-8 text-center">
        <v-icon icon="mdi-check-circle" color="success" size="64" class="mb-4" />
        <h2 class="text-h5 font-weight-bold mb-2">Already confirmed!</h2>
        <p class="text-medium-emphasis">This booking is confirmed for {{ booking.date }}, {{ booking.slot }}.</p>
      </v-card>

      <template v-else>
        <!-- header -->
        <div class="d-flex align-center ga-2 mb-6">
          <div class="owner-badge">OWNER</div>
          <h1 class="text-h5 font-weight-black">Review Booking</h1>
        </div>

        <!-- booking summary -->
        <v-card rounded="xl" elevation="1" class="mb-5">
          <v-card-title class="pa-4 pb-2 text-subtitle-1 font-weight-bold">Booking #{{ booking.id }}</v-card-title>
          <v-list density="comfortable" class="px-2">
            <v-list-item prepend-icon="mdi-account"        :title="booking.name"   :subtitle="booking.phone + ' · ' + booking.area" />
            <v-list-item prepend-icon="mdi-wrench"         :title="booking.category + (booking.service ? ' – ' + booking.service : '')" subtitle="Service" />
            <v-list-item prepend-icon="mdi-calendar-clock" :title="booking.date + ', ' + booking.slot" subtitle="Requested time" />
            <v-list-item v-if="booking.notes" prepend-icon="mdi-text" :title="booking.notes" subtitle="Notes" />
          </v-list>
        </v-card>

        <!-- proposed time from user (if user proposed) -->
        <v-card v-if="booking.status === 'user_proposed'" rounded="xl" color="orange-lighten-5" border="orange" class="mb-5">
          <v-card-text class="pa-4">
            <div class="d-flex align-center ga-2 mb-1">
              <v-icon icon="mdi-account-clock" color="orange-darken-3" />
              <span class="font-weight-bold text-orange-darken-3">Customer's requested time</span>
            </div>
            <p class="text-body-1 font-weight-bold mb-0">{{ booking.proposedDate }}, {{ booking.proposedSlot }}</p>
          </v-card-text>
        </v-card>

        <!-- success state after action -->
        <v-card v-if="done" rounded="xl" :color="doneColor" variant="flat" class="pa-6 text-center mb-5">
          <v-icon :icon="doneIcon" color="white" size="48" class="mb-3" />
          <h3 class="text-h6 font-weight-bold text-white mb-2">{{ doneTitle }}</h3>
          <p class="text-body-2 text-white mb-4" style="opacity:.9">{{ doneSub }}</p>
          <v-btn color="white" :href="whatsappUrl" target="_blank" prepend-icon="mdi-whatsapp" variant="flat"
            style="color:#1565C0!important">
            {{ doneWaLabel }}
          </v-btn>
        </v-card>

        <!-- action panel -->
        <template v-else>
          <!-- Accept current time -->
          <v-card rounded="xl" elevation="1" class="mb-4 action-card" @click="acceptTime">
            <v-card-text class="pa-5 d-flex align-center ga-4">
              <v-avatar color="success" size="52">
                <v-icon icon="mdi-check" color="white" size="28" />
              </v-avatar>
              <div>
                <div class="text-subtitle-1 font-weight-bold">Accept this booking</div>
                <div class="text-body-2 text-medium-emphasis">
                  Confirm {{ booking.status === 'user_proposed' ? booking.proposedDate + ', ' + booking.proposedSlot : booking.date + ', ' + booking.slot }}
                </div>
              </div>
              <v-spacer />
              <v-icon icon="mdi-chevron-right" color="success" />
            </v-card-text>
          </v-card>

          <!-- Propose new time -->
          <v-card rounded="xl" elevation="1" class="mb-4 action-card" @click="showReschedule = !showReschedule">
            <v-card-text class="pa-5 d-flex align-center ga-4">
              <v-avatar color="primary" size="52">
                <v-icon icon="mdi-calendar-edit" color="white" size="26" />
              </v-avatar>
              <div>
                <div class="text-subtitle-1 font-weight-bold">Propose a new time</div>
                <div class="text-body-2 text-medium-emphasis">Suggest a different date or slot</div>
              </div>
              <v-spacer />
              <v-icon :icon="showReschedule ? 'mdi-chevron-up' : 'mdi-chevron-down'" color="primary" />
            </v-card-text>

            <v-expand-transition>
              <div v-if="showReschedule" class="px-5 pb-5" @click.stop>
                <v-text-field
                  v-model="newDate"
                  type="date"
                  label="New date"
                  :min="today"
                  variant="outlined"
                  density="compact"
                  prepend-inner-icon="mdi-calendar"
                  class="mb-3"
                />
                <p class="text-body-2 text-medium-emphasis mb-2">Time slot</p>
                <v-chip-group v-model="newSlot" selected-class="text-white bg-primary" column class="mb-4">
                  <v-chip v-for="s in slots" :key="s" :value="s" variant="outlined" filter size="small">{{ s }}</v-chip>
                </v-chip-group>
                <v-btn
                  color="primary"
                  variant="flat"
                  block
                  :disabled="!newDate || !newSlot"
                  :loading="sending"
                  prepend-icon="mdi-send"
                  @click="proposeTime"
                >Send new time to customer</v-btn>
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
  const actionUrl = booking.value.status === 'user_proposed'
    ? `/api/booking/${route.params.id}/user-accept`
    : `/api/booking/${route.params.id}/accept`
  const res = await fetch(actionUrl, { method: 'POST' })
  const data = await res.json()
  whatsappUrl.value = data.waUrl
  doneColor.value = 'success'
  doneIcon.value = 'mdi-check-circle'
  doneTitle.value = 'Booking confirmed!'
  doneSub.value = 'Tap the button below to send the confirmation to the customer on WhatsApp.'
  doneWaLabel.value = 'Send confirmation to customer'
  done.value = true
}

async function proposeTime() {
  sending.value = true
  try {
    const res = await fetch(`/api/booking/${route.params.id}/reschedule`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ date: newDate.value, slot: newSlot.value }),
    })
    const data = await res.json()
    whatsappUrl.value = data.waUrl
    doneColor.value = 'primary'
    doneIcon.value = 'mdi-calendar-edit'
    doneTitle.value = 'New time sent to customer!'
    doneSub.value = 'Tap below to send the proposed time to the customer over WhatsApp. They can accept or suggest another time.'
    doneWaLabel.value = 'Send proposed time to customer'
    done.value = true
  } finally {
    sending.value = false
  }
}
</script>

<style scoped>
.action-page { min-height: 100vh; background: #F0F4FA; padding-top: 16px; }
.owner-badge {
  background: #1565C0; color: white;
  font-size: .7rem; font-weight: 800;
  padding: 3px 10px; border-radius: 999px;
  letter-spacing: 1px;
}
.action-card { cursor: pointer; transition: box-shadow .2s, transform .15s; }
.action-card:hover { transform: translateY(-2px); box-shadow: 0 6px 24px rgba(0,0,0,.10) !important; }
</style>
