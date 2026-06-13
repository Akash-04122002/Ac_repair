# QuickFix Chennai — Electrical & Plumbing Booking Site

A fast, mobile-first one-page site whose purpose is **booking appointments via WhatsApp**.
Customer picks a service + preferred time → taps a button → a pre-filled booking message
opens in WhatsApp and lands in your chat. You reply to confirm the slot.

- **Frontend:** Vue 3 + Vuetify 3 (Vite)
- **Backend:** Go (standard library only) — serves the site + a small config/booking API
- **No database** — content lives in `backend/config.json`; bookings travel over WhatsApp
- **Ships as one binary** (the built site is embedded into the Go executable)

---

## Edit your content (no rebuild needed)

Everything you'll want to change lives in **`backend/config.json`**:

- `business` — name, tagline, WhatsApp number, phone, hours, area
- `categories` → `services` — your services and **price ranges**
- `whyUs`, `reviews`, `areas`, `faq` — trust content

The running server reads `config.json` from disk on every request, so after deployment you
can edit it and refresh — **no rebuild required**. (A copy is also embedded as a fallback.)

> WhatsApp number is set in `config.json` → `business.whatsapp` (digits only, with country
> code, e.g. `919043050841`).

---

## Run locally

### Option A — production mode (single binary, what you deploy)
```bash
# 1. Build the frontend (outputs into backend/dist)
cd frontend
npm install
npm run build

# 2. Build & run the Go server (serves everything)
cd ../backend
go build -o quickfix.exe .
./quickfix.exe                 # http://localhost:8080
# If port 8080 is busy:  PORT=8099 ./quickfix.exe
```
Open the printed URL. Fill the booking form → it opens WhatsApp with the message pre-filled.

### Option B — dev mode (hot reload while editing the UI)
```bash
# terminal 1 — API
cd backend && go run .          # serves /api on :8080

# terminal 2 — Vite dev server (proxies /api to :8080)
cd frontend && npm run dev      # http://localhost:5173
```

---

## Deploy free (single binary)

The whole app is **one self-contained binary**, so any free host that runs a Go web service works.

### Recommended: Render (free)
1. Push this project to a GitHub repo.
2. On [render.com](https://render.com) → **New → Web Service** → connect the repo.
3. Settings:
   - **Root directory:** `backend`
   - **Build command:**
     ```
     cd ../frontend && npm install && npm run build && cd ../backend && go build -o quickfix .
     ```
   - **Start command:** `./quickfix`
   - Render provides the `PORT` env var automatically — the server already reads it.
4. Deploy → you get a free `https://<your-app>.onrender.com` URL.

> **Cold start:** free tiers sleep after inactivity, so the first visit after idle takes a
> few seconds to wake. Booking still works even during wake-up — the frontend has a local
> fallback that builds the WhatsApp link client-side.

Other free options that work the same way: **Fly.io**, **Koyeb**, **Google Cloud Run** (free tier).

---

## How a booking flows

```
Customer → booking stepper (category → service → date/time → details → review)
        → "Book on WhatsApp"
        → POST /api/booking  (Go validates + logs to bookings.json + returns wa.me URL)
        → browser opens wa.me/<number>?text=...   (pre-filled message)
        → customer taps SEND → booking arrives in YOUR WhatsApp
        → you reply to confirm the time slot
```

## API
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/config` | Returns `config.json` (site content) |
| POST | `/api/booking` | Validates, logs, returns `{ waUrl, message }` |
| GET | `/*` | Serves the embedded Vue site (SPA) |

## Known limitations (by design — no database)
- **Time slots are *preferred*, not locked** — you confirm/adjust over WhatsApp.
- **No automatic reminders/confirmations** (needs the paid WhatsApp Business API).
- **`bookings.json` is best-effort** — free hosts wipe disk on restart; your WhatsApp chat is
  the durable record.
- **No payments, login, or live tracking** — out of scope (all need a real database).
