<template>
  <div class="status-page">
    <v-container class="py-10" style="max-width:560px">

      <!-- loading -->
      <div v-if="loading" class="text-center py-16">
        <v-progress-circular indeterminate color="primary" size="48" class="mb-4" />
        <p class="text-medium-emphasis">Loading booking…</p>
      </div>

      <!-- not found -->
      <v-card v-else-if="!booking" rounded="xl" elevation="2" class="pa-8 text-center">
        <v-icon icon="mdi-alert-circle-outline" color="error" size="56" class="mb-4" />
        <h2 class="text-h6 font-weight-bold mb-2">Booking not found</h2>
        <p class="text-medium-emphasis mb-4">The link may be invalid or expired.</p>
        <v-btn color="primary" href="/">Go to homepage</v-btn>
      </v-card>

      <!-- booking found -->
      <template v-else>
        <!-- status banner -->
        <v-card :color="statusColor" variant="flat" rounded="xl" class="mb-5 pa-5">
          <div class="d-flex align-center ga-3">
            <v-icon :icon="statusIcon" color="white" size="40" />
            <div>
              <div class="text-h6 font-weight-bold text-white">{{ statusTitle }}</div>
              <div class="text-body-2 text-white" style="opacity:.88">{{ statusSub }}</div>
            </div>
          </div>
        </v-card>

        <!-- booking details -->
        <v-card rounded="xl" elevation="1" class="mb-5">
          <v-card-title class="text-subtitle-1 font-weight-bold pa-4 pb-2">Booking details</v-card-title>
          <v-list density="comfortable" class="px-2">
            <v-list-item prepend-icon="mdi-wrench"        :title="booking.category + (booking.service ? ' – ' + booking.service : '')" subtitle="Service" />
            <v-list-item prepend-icon="mdi-cash"          :title="booking.price || '—'"  subtitle="Est. price" />
            <v-list-item prepend-icon="mdi-calendar-clock" :title="confirmedTime"         subtitle="Confirmed time" />
            <v-list-item prepend-icon="mdi-account"       :title="booking.name"          :subtitle="booking.phone + ' · ' + booking.area" />
            <v-list-item v-if="booking.notes" prepend-icon="mdi-text" :title="booking.notes" subtitle="Notes" />
          </v-list>
        </v-card>

        <!-- proposed time (if owner or user proposed a change) -->
        <v-card v-if="hasProposal" rounded="xl" elevation="1" class="mb-5" color="amber-lighten-5" border="amber">
          <v-card-text class="pa-4">
            <div class="d-flex align-center ga-2 mb-1">
              <v-icon icon="mdi-clock-outline" color="amber-darken-3" />
              <span class="font-weight-bold text-amber-darken-3">Proposed new time</span>
            </div>
            <p class="text-body-1 font-weight-bold mb-0">{{ booking.proposedDate }}, {{ booking.proposedSlot }}</p>
          </v-card-text>
        </v-card>

        <!-- history timeline -->
        <v-card v-if="booking.history?.length" rounded="xl" elevation="1" class="mb-5">
          <v-card-title class="text-subtitle-1 font-weight-bold pa-4 pb-2">Timeline</v-card-title>
          <v-timeline side="end" density="compact" class="px-4 pb-4">
            <v-timeline-item
              v-for="(h, i) in booking.history"
              :key="i"
              :dot-color="h.actor === 'owner' ? 'primary' : 'secondary'"
              size="small"
            >
              <div class="text-body-2 font-weight-medium">{{ historyLabel(h) }}</div>
              <div class="text-caption text-medium-emphasis">{{ fmtTime(h.at) }} · {{ h.actor }}</div>
            </v-timeline-item>
          </v-timeline>
        </v-card>

        <!-- confirmed CTA -->
        <v-card v-if="booking.status === 'confirmed'" rounded="xl" color="success" variant="flat" class="pa-5 text-center">
          <v-icon icon="mdi-check-circle" color="white" size="48" class="mb-3" />
          <h3 class="text-h6 font-weight-bold text-white mb-1">You're all set!</h3>
          <p class="text-body-2 text-white mb-0" style="opacity:.9">
            Our technician will arrive at the confirmed time. You'll hear from us on WhatsApp.
          </p>
        </v-card>

        <!-- auto-refresh note -->
        <p v-else class="text-caption text-center text-medium-emphasis mt-4">
          <v-icon icon="mdi-refresh" size="14" class="mr-1" />This page refreshes automatically every 15 seconds.
        </p>
      </template>
    </v-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const booking = ref(null)
const loading = ref(true)
let timer = null

async function fetchBooking() {
  try {
    const res = await fetch(`/api/booking/${route.params.id}`)
    if (res.ok) booking.value = await res.json()
    else booking.value = null
  } catch {
    booking.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchBooking()
  timer = setInterval(fetchBooking, 15000)
})
onUnmounted(() => clearInterval(timer))

const statusColor = computed(() => {
  if (!booking.value) return 'grey'
  const map = { pending: 'primary', owner_proposed: 'amber-darken-2', user_proposed: 'orange-darken-2', confirmed: 'success' }
  return map[booking.value.status] || 'grey'
})
const statusIcon = computed(() => {
  if (!booking.value) return 'mdi-help'
  const map = { pending: 'mdi-clock-outline', owner_proposed: 'mdi-calendar-edit', user_proposed: 'mdi-calendar-edit', confirmed: 'mdi-check-circle' }
  return map[booking.value.status] || 'mdi-help'
})
const statusTitle = computed(() => {
  if (!booking.value) return ''
  const map = { pending: 'Waiting for confirmation', owner_proposed: 'New time proposed by us', user_proposed: 'Your time request sent', confirmed: 'Booking Confirmed ✅' }
  return map[booking.value.status] || ''
})
const statusSub = computed(() => {
  if (!booking.value) return ''
  const map = {
    pending: 'We received your booking and will confirm shortly.',
    owner_proposed: 'We suggested a new time — please check your WhatsApp.',
    user_proposed: 'We received your preferred time and will respond shortly.',
    confirmed: 'Your appointment is locked in. See you soon!',
  }
  return map[booking.value.status] || ''
})
const hasProposal = computed(() => booking.value?.proposedDate && booking.value.status !== 'confirmed')
const confirmedTime = computed(() => {
  if (!booking.value) return '—'
  if (booking.value.status === 'confirmed' && booking.value.proposedDate) {
    return `${booking.value.proposedDate}, ${booking.value.proposedSlot} (updated)`
  }
  return `${booking.value.date}, ${booking.value.slot}`
})

function historyLabel(h) {
  if (h.action === 'accepted') return h.actor === 'owner' ? 'Owner confirmed booking' : 'You accepted the time'
  if (h.action === 'proposed') return `${h.actor === 'owner' ? 'Owner' : 'You'} proposed: ${h.date}, ${h.slot}`
  return h.action
}
function fmtTime(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleString('en-IN', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.status-page { min-height: 100vh; background: #F0F4FA; padding-top: 16px; }
</style>
