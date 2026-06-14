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
        <div v-if="done" class="text-center py-6">
          <v-icon icon="mdi-check-circle" color="success" size="64" class="mb-3" />
          <h3 class="text-h6 font-weight-bold mb-1">Booking Submitted!</h3>
          <p class="text-body-2 text-medium-emphasis mb-4">
            Tap <strong>Send</strong> in WhatsApp. We’ll confirm or propose a new time — you can track it below.
          </p>
          <v-btn color="#25D366" variant="flat" :href="waUrl" target="_blank" prepend-icon="mdi-whatsapp" block class="mb-3">
            Open WhatsApp &amp; Send
          </v-btn>
          <v-btn
            v-if="statusUrl"
            color="primary"
            variant="tonal"
            :href="statusUrl"
            target="_blank"
            prepend-icon="mdi-magnify"
            block
          >Track booking status</v-btn>
        </div>

        <div v-else>
          <!-- Step 1: Category -->
          <div v-if="step === 1">
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
          </div>

          <!-- Step 2: Service -->
          <div v-if="step === 2">
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
          </div>

          <!-- Step 3: Date + time -->
          <div v-if="step === 3">
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
          </div>

          <!-- Step 4: Details -->
          <div v-if="step === 4">
            <p class="text-subtitle-1 font-weight-medium mb-3">Your details</p>
            <v-text-field v-model="form.name" label="Full name" variant="outlined" prepend-inner-icon="mdi-account" :error-messages="errors.name" />
            <v-text-field v-model="form.phone" label="Phone number" variant="outlined" prepend-inner-icon="mdi-phone" inputmode="numeric" maxlength="10" :error-messages="errors.phone" />
            <v-text-field v-model="form.doorNo" label="Door / Flat No. & Building name" variant="outlined" prepend-inner-icon="mdi-home-city" :error-messages="errors.doorNo" />
            <v-combobox 
              v-model="form.area" 
              :items="areas" 
              label="Area / locality in Chennai" 
              variant="outlined" 
              prepend-inner-icon="mdi-map-marker" 
              :error-messages="errors.area"
            >
              <template #append-inner>
                <v-btn
                  icon="mdi-crosshairs-gps"
                  variant="text"
                  density="compact"
                  color="primary"
                  title="Use current location"
                  @click="getLocation"
                ></v-btn>
              </template>
            </v-combobox>
            <v-chip v-if="form.coords" color="success" variant="tonal" prepend-icon="mdi-map-marker-radius" class="mb-3">
              📍 GPS: {{ form.coords }}
            </v-chip>
            <v-textarea v-model="form.notes" label="Notes (optional)" variant="outlined" rows="2" prepend-inner-icon="mdi-text" />
          </div>

          <!-- Step 5: Review -->
          <div v-if="step === 5">
            <p class="text-subtitle-1 font-weight-medium mb-3">Review your booking</p>
            <v-list density="comfortable">
              <v-list-item prepend-icon="mdi-wrench" :title="`${form.category} – ${form.service}`" subtitle="Service" />
              <v-list-item prepend-icon="mdi-cash" :title="form.price || '—'" subtitle="Estimated price" />
              <v-list-item prepend-icon="mdi-calendar-clock" :title="`${prettyDate} · ${form.slot}`" subtitle="Preferred time" />
              <v-list-item prepend-icon="mdi-account" :title="form.name" :subtitle="`${form.phone} · ${form.area}`" />
              <v-list-item v-if="form.coords" prepend-icon="mdi-crosshairs-gps" :title="form.coords" subtitle="GPS coordinates" />
              <v-list-item v-if="form.notes" prepend-icon="mdi-text" :title="form.notes" subtitle="Notes" />
            </v-list>
            <v-alert type="info" variant="tonal" density="compact" class="mt-2">
              Your chosen time is a <strong>preferred slot</strong>. We’ll confirm or suggest the closest time on WhatsApp.
            </v-alert>
          </div>
        </div>
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
  areas: { type: Array, default: () => [] },
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
const statusUrl = ref('')
const loadingLocation = ref(false)
const locationNote = ref('')

const form = reactive({
  categoryId: null, category: '', service: '', price: '',
  date: '', slot: '', name: '', phone: '', doorNo: '', area: '', coords: '', notes: '',
})
const errors = reactive({ name: '', phone: '', doorNo: '', area: '' })

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
  const e = { name: '', phone: '', doorNo: '', area: '' }
  if (!form.name.trim()) e.name = 'Please enter your name'
  if (!/^\d{10}$/.test(form.phone.trim())) e.phone = 'Enter a 10-digit phone number'
  if (!form.doorNo.trim()) e.doorNo = 'Please enter your door/flat number'
  if (!form.area.trim()) e.area = 'Please enter your area'
  if (showErrors) Object.assign(errors, e)
  return !e.name && !e.phone && !e.doorNo && !e.area
}

function next() {
  if (step.value === 4 && !validateDetails(true)) return
  if (step.value < 5) step.value++
}

async function getLocation() {
  if (!navigator.geolocation) {
    errors.area = 'Geolocation is not supported by your browser'
    return
  }
  
  errors.area = ''
  errors.doorNo = ''
  locationNote.value = ''
  loadingLocation.value = true
  form.area = ''
  form.doorNo = ''
  form.coords = ''
  
  navigator.geolocation.getCurrentPosition(async (position) => {
    try {
      const { latitude, longitude } = position.coords
      form.coords = `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`
      let filled = false
      
      // --- Try Nominatim (OpenStreetMap) first ---
      try {
        const url = 'https://nominatim.openstreetmap.org/reverse?format=json' +
          '&lat=' + latitude + '&lon=' + longitude + '&zoom=21&addressdetails=1'
        const res = await fetch(url, {
          headers: { 'Accept-Language': 'en', 'User-Agent': 'QuickFixChennai/1.0' }
        })
        if (res.ok) {
          const data = await res.json()
          const addr = data.address || {}
          
          const doorParts = [
            addr.house_number || addr.building || '',
            addr.road || addr.street || addr.pedestrian || ''
          ].filter(Boolean)
          
          if (doorParts.length > 0) {
            form.doorNo = doorParts.join(', ')
          } else if (data.name) {
            form.doorNo = data.name
          }
          
          const areaParts = [
            addr.neighbourhood || addr.residential || '',
            addr.suburb || addr.city_district || addr.village || addr.town || '',
            addr.city || addr.county || '',
            addr.postcode || ''
          ].filter(Boolean)
          
          if (areaParts.length > 0) {
            form.area = areaParts.join(', ')
            filled = true
          }
        }
      } catch (e) {
        console.warn('Nominatim failed:', e)
      }
      
      // --- Backup: BigDataCloud (free, no key needed) ---
      if (!filled) {
        try {
          const url2 = 'https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=' +
            latitude + '&longitude=' + longitude + '&localityLanguage=en'
          const res2 = await fetch(url2)
          if (res2.ok) {
            const d = await res2.json()
            
            if (!form.doorNo) {
              const doorParts2 = [
                d.localityInfo?.administrative?.[0]?.name || '',
                d.principalSubdivision || ''
              ]
              // Use the locality for street info
              form.doorNo = d.locality || d.city || ''
            }
            
            const areaParts2 = [
              d.locality || '',
              d.city || '',
              d.postcode || ''
            ].filter(Boolean)
            
            // Remove duplicates
            const unique = [...new Set(areaParts2)]
            if (unique.length > 0) {
              form.area = unique.join(', ')
              filled = true
            }
          }
        } catch (e2) {
          console.warn('BigDataCloud also failed:', e2)
        }
      }
      
      if (!filled) {
        form.area = `Lat: ${latitude.toFixed(5)}, Lng: ${longitude.toFixed(5)}`
        locationNote.value = 'Could not fetch address. Please type your area manually.'
      }
      
      if (form.doorNo) errors.doorNo = ''
      if (form.area) errors.area = ''
      
      if (!form.doorNo || !form.doorNo.match(/^\d/)) {
        locationNote.value = 'Street auto-filled! Please type your door/flat number.'
      }
    } catch (err) {
      console.error(err)
      form.area = ''
      errors.area = 'Location fetch failed. Please type your area manually.'
    } finally {
      loadingLocation.value = false
    }
  }, (err) => {
    form.area = ''
    form.doorNo = ''
    form.coords = ''
    errors.area = 'Please allow location access to use this feature.'
    loadingLocation.value = false
  }, { 
    enableHighAccuracy: true, 
    timeout: 15000, 
    maximumAge: 0 
  })
}

async function submit() {
  if (!validateDetails(true)) { step.value = 4; return }
  const mapsLink = form.coords ? `https://maps.google.com/?q=${form.coords.replace(' ', '')}` : ''
  const result = await submitBooking(
    { ...form, date: prettyDate.value, phone: '+91' + form.phone, area: `${form.doorNo}, ${form.area}`, mapsLink },
    props.business,
    { urgent: props.urgent }
  )
  waUrl.value = result.waUrl
  statusUrl.value = result.statusUrl || ''
  done.value = true
  window.open(waUrl.value, '_blank')
}

// Reset / apply prefill each time the dialog opens.
watch(open, (isOpen) => {
  if (!isOpen) return
  done.value = false
  waUrl.value = ''
  Object.assign(errors, { name: '', phone: '', doorNo: '', area: '' })
  Object.assign(form, {
    categoryId: null, category: '', service: '', price: '',
    date: '', slot: '', name: '', phone: '', doorNo: '', area: '', coords: '', notes: '',
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
