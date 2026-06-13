import { buildBookingMessage, buildWaUrl } from './waMessage'

// Fallback content used if the API is unreachable (e.g. free-host cold start).
// Mirrors backend/config.json so the site is never broken.
export const FALLBACK_CONFIG = {
  business: {
    name: 'QuickFix Chennai — Electrical & Plumbing',
    whatsapp: '919043050841',
    phone: '+919043050841',
    hours: 'Mon–Sun, 8 AM – 8 PM',
    area: 'All over Chennai',
  },
  categories: [],
  reviews: [],
  areas: [],
  faq: [],
}

export async function fetchConfig() {
  try {
    const res = await fetch('/api/config', { cache: 'no-store' })
    if (!res.ok) throw new Error('bad status')
    return await res.json()
  } catch {
    return FALLBACK_CONFIG
  }
}

// Submits the booking. Tries the Go API (which validates + logs and returns the
// wa.me URL). If the server is asleep/unreachable, builds the same URL locally so
// booking ALWAYS works.
export async function submitBooking(booking, business, { urgent = false } = {}) {
  try {
    const res = await fetch('/api/booking', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...booking, urgent }),
    })
    if (res.ok) {
      const data = await res.json()
      if (data.waUrl) return data.waUrl
    }
  } catch {
    // fall through to local build
  }
  const msg = buildBookingMessage(booking, { urgent })
  return buildWaUrl(business.whatsapp, msg)
}
