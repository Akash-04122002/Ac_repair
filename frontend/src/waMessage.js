// Builds the pre-filled WhatsApp message text from a booking object.
// Kept in one place so the frontend fallback and the displayed summary stay identical.

export function buildBookingMessage(b, { urgent = false } = {}) {
  const lines = []
  lines.push(urgent ? '🚨 URGENT Booking' : '🛠️ New Booking')
  lines.push('')
  if (b.category) lines.push(`Service: ${b.category}${b.service ? ' – ' + b.service : ''}`)
  if (b.price) lines.push(`Est. price: ${b.price}`)
  if (b.date || b.slot) lines.push(`Preferred: ${[b.date, b.slot].filter(Boolean).join(', ')}`)
  lines.push('')
  if (b.name) lines.push(`Name: ${b.name}`)
  if (b.phone) lines.push(`Phone: ${b.phone}`)
  if (b.area) lines.push(`Area: ${b.area}`)
  if (b.notes) lines.push(`Notes: ${b.notes}`)
  lines.push('')
  lines.push('(Sent from website — please confirm the time slot.)')
  return lines.join('\n')
}

// Returns the wa.me URL that opens WhatsApp with the message pre-filled.
export function buildWaUrl(whatsappNumber, message) {
  return `https://wa.me/${whatsappNumber}?text=${encodeURIComponent(message)}`
}
