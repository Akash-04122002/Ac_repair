<template>
  <v-dialog v-model="open" max-width="640" scrollable>
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon :icon="urgent ? 'mdi-alarm-light' : 'mdi-calendar-check'" :color="urgent ? 'error' : 'primary'" class="mr-2" />
        {{ urgent ? 'Emergency booking' : 'Book an appointment' }}
        <v-spacer />
        <v-btn icon="mdi-close" variant="text" @click="open = false" />
      </v-card-title>

      <!-- progress dots -->
      <div class="px-4">
        <v-progress-linear :model-value="(step / 5) * 100" color="primary" height="6" rounded />
        <div class="text-caption text-medium-emphasis text-center mt-1">Step {{ step }} of 5</div>
      </div>

      <v-divider />

      <v-card-text style="min-height: 280px">
        <!-- Success state -->
        <div v-if="done" class="text-center py-8">
          <v-icon icon="mdi-whatsapp" color="#25D366" size="64" class="mb-4" />
          <h3 class="text-h6 font-weight-bold mb-2">Opening WhatsApp…</h3>
          <p class="text-body-2 text-medium-emphasis mb-4">
            Tap <strong>Send</strong> in WhatsApp to confirm your booking. We’ll reply to confirm your time slot.
          </p>
          <v-btn color="primary" variant="flat" :href="waUrl" target="_blank" prepend-icon="mdi-whatsapp">
            Open WhatsApp again
          </v-btn>
        </div>

        <v-window v-else v-model="step">
          <!-- Step 1: Category -->
          <v-window-item :value="1">
            <p class="text-subtitle-1 font-weight-medium mb-3">What do you need?</p>
            <v-row>
              <v-col v-for="c in categories" :key="c.id" cols="12" sm="4">
                <v-card
                  :color="form.categoryId === c.id ? 'primary' : undefined"
                  :variant="form.categoryId === c.id ? 'flat' : 'outlined'"
                  @click="selectCategory(c)"
                >
                  <v-card-text class="text-center">
                    <v-icon :icon="c.icon" size="36" class="mb-2" />
                    <div class="font-weight-bold">{{ c.name }}</div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>

          <!-- Step 2: Service -->
          <v-window-item :value="2">
            <p class="text-subtitle-1 font-weight-medium mb-3">Choose a service</p>
            <v-list v-if="activeCategory" lines="two" select-strategy="single-independent">
              <v-list-item
                v-for="(s, i) in activeCategory.services"
                :key="i"
                :active="form.service === s.name"
                color="primary"
                @click="selectService(s)"
              >
                <v-list-item-title class="font-weight-medium">{{ s.name }}</v-list-item-title>
                <v-list-item-subtitle class="text-secondary">{{ s.price }}</v-list-item-subtitle>
                <template #append>
                  <v-icon v-if="form.service === s.name" icon="mdi-check-circle" color="primary" />
                </template>
              </v-list-item>
            </v-list>
          </v-window-item>

          <!-- Step 3: Date + time -->
          <v-window-item :value="3">
            <p class="text-subtitle-1 font-weight-medium mb-3">Preferred date & time</p>
            <v-text-field
              v-model="form.date"
              type="date"
              label="Date"
              :min="today"
              variant="outlined"
              prepend-inner-icon="mdi-calendar"
            />
            <p class="text-body-2 text-medium-emphasis mb-2">Time slot</p>
            <v-chip-group v-model="form.slot" selected-class="text-white bg-primary" column>
              <v-chip v-for="s in slots" :key="s" :value="s" variant="outlined" filter>{{ s }}</v-chip>
            </v-chip-group>
          </v-window-item>

          <!-- Step 4: Details -->
          <v-window-item :value="4">
            <p class="text-subtitle-1 font-weight-medium mb-3">Your details</p>
            <v-text-field v-model="form.name" label="Full name" variant="outlined" prepend-inner-icon="mdi-account" :error-messages="errors.name" />
            <v-text-field v-model="form.phone" label="Phone number" variant="outlined" prepend-inner-icon="mdi-phone" inputmode="numeric" maxlength="10" :error-messages="errors.phone" />
            <v-text-field v-model="form.area" label="Area / locality in Chennai" variant="outlined" prepend-inner-icon="mdi-map-marker" :error-messages="errors.area" />
            <v-textarea v-model="form.notes" label="Notes (optional)" variant="outlined" rows="2" prepend-inner-icon="mdi-text" />
          </v-window-item>

          <!-- Step 5: Review -->
          <v-window-item :value="5">
            <p class="text-subtitle-1 font-weight-medium mb-3">Review your booking</p>
            <v-list density="comfortable">
              <v-list-item prepend-icon="mdi-wrench" :title="`${form.category} – ${form.service}`" subtitle="Service" />
              <v-list-item prepend-icon="mdi-cash" :title="form.price || '—'" subtitle="Estimated price" />
              <v-list-item prepend-icon="mdi-calendar-clock" :title="`${prettyDate} · ${form.slot}`" subtitle="Preferred time" />
              <v-list-item prepend-icon="mdi-account" :title="form.name" :subtitle="`${form.phone} · ${form.area}`" />
              <v-list-item v-if="form.notes" prepend-icon="mdi-text" :title="form.notes" subtitle="Notes" />
            </v-list>
            <v-alert type="info" variant="tonal" density="compact" class="mt-2">
              Your chosen time is a <strong>preferred slot</strong>. We’ll confirm or suggest the closest time on WhatsApp.
            </v-alert>
          </v-window-item>
        </v-window>
      </v-card-text>

      <v-divider v-if="!done" />

      <v-card-actions v-if="!done" class="pa-4">
        <v-btn v-if="step > 1" variant="text" prepend-icon="mdi-chevron-left" @click="step--">Back</v-btn>
        <v-spacer />
        <v-btn v-if="step < 5" color="primary" variant="flat" append-icon="mdi-chevron-right" :disabled="!canNext" @click="next">
          Next
        </v-btn>
        <v-btn v-else color="#25D366" variant="flat" prepend-icon="mdi-whatsapp" @click="submit">
          Book on WhatsApp
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { submitBooking } from '../api'

const props = defineProps({
  modelValue: Boolean,
  categories: { type: Array, default: () => [] },
  business: { type: Object, required: true },
  prefill: { type: Object, default: () => ({ categoryId: null, service: null }) },
  urgent: Boolean,
})
const emit = defineEmits(['update:modelValue'])

const open = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const step = ref(1)
const done = ref(false)
const waUrl = ref('')

const form = reactive({
  categoryId: null, category: '', service: '', price: '',
  date: '', slot: '', name: '', phone: '', area: '', notes: '',
})
const errors = reactive({ name: '', phone: '', area: '' })

const today = new Date().toISOString().slice(0, 10)
const slots = [
  '8:00 – 10:00 AM', '10:00 – 12:00 PM', '12:00 – 2:00 PM',
  '2:00 – 4:00 PM', '4:00 – 6:00 PM', '6:00 – 8:00 PM',
]

const activeCategory = computed(() => props.categories.find((c) => c.id === form.categoryId) || null)
const prettyDate = computed(() => {
  if (!form.date) return '—'
  const d = new Date(form.date + 'T00:00:00')
  return d.toLocaleDateString('en-IN', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' })
})

const canNext = computed(() => {
  if (step.value === 1) return !!form.categoryId
  if (step.value === 2) return !!form.service
  if (step.value === 3) return !!form.date && !!form.slot
  if (step.value === 4) return validateDetails(false)
  return true
})

function selectCategory(c) {
  form.categoryId = c.id
  form.category = c.name
  form.service = ''
  form.price = ''
}
function selectService(s) {
  form.service = s.name
  form.price = s.price
}

function validateDetails(showErrors = true) {
  const e = { name: '', phone: '', area: '' }
  if (!form.name.trim()) e.name = 'Please enter your name'
  if (!/^\d{10}$/.test(form.phone.trim())) e.phone = 'Enter a 10-digit phone number'
  if (!form.area.trim()) e.area = 'Please enter your area'
  if (showErrors) Object.assign(errors, e)
  return !e.name && !e.phone && !e.area
}

function next() {
  if (step.value === 4 && !validateDetails(true)) return
  if (step.value < 5) step.value++
}

async function submit() {
  if (!validateDetails(true)) { step.value = 4; return }
  waUrl.value = await submitBooking(
    { ...form, date: prettyDate.value, phone: '+91' + form.phone },
    props.business,
    { urgent: props.urgent }
  )
  done.value = true
  window.open(waUrl.value, '_blank')
}

// Reset / apply prefill each time the dialog opens.
watch(open, (isOpen) => {
  if (!isOpen) return
  done.value = false
  waUrl.value = ''
  Object.assign(errors, { name: '', phone: '', area: '' })
  Object.assign(form, {
    categoryId: null, category: '', service: '', price: '',
    date: '', slot: '', name: '', phone: '', area: '', notes: '',
  })
  step.value = 1

  const pre = props.prefill || {}
  if (pre.categoryId) {
    const c = props.categories.find((x) => x.id === pre.categoryId)
    if (c) {
      selectCategory(c)
      step.value = 2
      if (pre.service) {
        const s = c.services.find((x) => x.name === pre.service)
        if (s) { selectService(s); step.value = 3 }
      }
    }
  }
})
</script>
