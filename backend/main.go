package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed all:dist
var distFS embed.FS

//go:embed config.json
var embeddedConfig []byte

// ── Data types ────────────────────────────────────────────────────────────────

type BookingStatus string

const (
	StatusPending      BookingStatus = "pending"      // just submitted, owner not yet responded
	StatusOwnerProposed BookingStatus = "owner_proposed" // owner suggested a new time
	StatusUserProposed  BookingStatus = "user_proposed"  // user requested a different time
	StatusConfirmed    BookingStatus = "confirmed"    // both sides accepted
)

type Booking struct {
	ID        string        `json:"id"`
	At        string        `json:"at"`
	Status    BookingStatus `json:"status"`
	Category  string        `json:"category"`
	Service   string        `json:"service"`
	Price     string        `json:"price"`
	Date      string        `json:"date"`      // original requested date
	Slot      string        `json:"slot"`      // original requested slot
	ProposedDate string     `json:"proposedDate,omitempty"` // owner's counter-proposal
	ProposedSlot string     `json:"proposedSlot,omitempty"`
	Name      string        `json:"name"`
	Phone     string        `json:"phone"`
	Area      string        `json:"area"`
	Notes     string        `json:"notes,omitempty"`
	Urgent    bool          `json:"urgent"`
	History   []HistoryEntry `json:"history,omitempty"`
}

type HistoryEntry struct {
	At     string `json:"at"`
	Actor  string `json:"actor"` // "owner" | "user"
	Action string `json:"action"`
	Date   string `json:"date,omitempty"`
	Slot   string `json:"slot,omitempty"`
}

// ── In-memory store backed by bookings.json ───────────────────────────────────

var (
	mu       sync.RWMutex
	bookings = map[string]*Booking{}
)

func loadBookings() {
	data, err := os.ReadFile("bookings.json")
	if err != nil {
		return
	}
	// file is newline-delimited JSON objects
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var b Booking
		if json.Unmarshal([]byte(line), &b) == nil && b.ID != "" {
			bookings[b.ID] = &b
		}
	}
}

func saveBookings() {
	var sb strings.Builder
	for _, b := range bookings {
		line, _ := json.Marshal(b)
		sb.Write(line)
		sb.WriteByte('\n')
	}
	os.WriteFile("bookings.json", []byte(sb.String()), 0644)
}

// ── ID generation ─────────────────────────────────────────────────────────────

const idChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func newID() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = idChars[rand.Intn(len(idChars))]
	}
	return string(b)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func siteBase(r *http.Request) string {
	base := os.Getenv("SITE_URL")
	if base != "" {
		return strings.TrimRight(base, "/")
	}
	scheme := "https"
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") == "" {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}

func now() string { return time.Now().Format(time.RFC3339) }

func whatsappNumber() string {
	data, err := os.ReadFile("config.json")
	if err != nil {
		data = embeddedConfig
	}
	var cfg struct {
		Business struct {
			WhatsApp string `json:"whatsapp"`
		} `json:"business"`
	}
	if json.Unmarshal(data, &cfg) == nil && cfg.Business.WhatsApp != "" {
		return cfg.Business.WhatsApp
	}
	return "919043050841"
}

func json200(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func corsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// GET /api/config
func handleConfig(w http.ResponseWriter, r *http.Request) {
	corsHeaders(w)
	data, err := os.ReadFile("config.json")
	if err != nil {
		data = embeddedConfig
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(data)
}

// POST /api/booking  — create a new booking, return waUrl + bookingId
func handleCreateBooking(w http.ResponseWriter, r *http.Request) {
	corsHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Category string `json:"category"`
		Service  string `json:"service"`
		Price    string `json:"price"`
		Date     string `json:"date"`
		Slot     string `json:"slot"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Area     string `json:"area"`
		Notes    string `json:"notes"`
		Urgent   bool   `json:"urgent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Area) == "" {
		http.Error(w, "name, phone and area are required", http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	id := newID()
	for bookings[id] != nil {
		id = newID()
	}
	b := &Booking{
		ID:       id,
		At:       now(),
		Status:   StatusPending,
		Category: req.Category,
		Service:  req.Service,
		Price:    req.Price,
		Date:     req.Date,
		Slot:     req.Slot,
		Name:     req.Name,
		Phone:    req.Phone,
		Area:     req.Area,
		Notes:    req.Notes,
		Urgent:   req.Urgent,
	}
	bookings[id] = b
	saveBookings()
	mu.Unlock()

	base := siteBase(r)
	ownerLink := base + "/owner/" + id
	statusLink := base + "/booking/" + id

	// Build WhatsApp message to owner with Accept / Reschedule magic links
	var sb strings.Builder
	if req.Urgent {
		sb.WriteString("🚨 *URGENT Booking Request*\n\n")
	} else {
		sb.WriteString("🛠️ *New Booking Request*\n\n")
	}
	sb.WriteString("*Booking ID:* " + id + "\n")
	sb.WriteString("*Service:* " + req.Category)
	if req.Service != "" {
		sb.WriteString(" – " + req.Service)
	}
	sb.WriteString("\n")
	if req.Price != "" {
		sb.WriteString("*Price:* " + req.Price + "\n")
	}
	sb.WriteString("*Requested time:* " + req.Date + ", " + req.Slot + "\n\n")
	sb.WriteString("*Customer:* " + req.Name + "\n")
	sb.WriteString("*Phone:* " + req.Phone + "\n")
	sb.WriteString("*Area:* " + req.Area + "\n")
	if req.Notes != "" {
		sb.WriteString("*Notes:* " + req.Notes + "\n")
	}
	sb.WriteString("\n──────────────────\n")
	sb.WriteString("👇 *Open this link to Accept or Propose a new time:*\n")
	sb.WriteString(ownerLink + "\n")
	sb.WriteString("──────────────────\n")
	sb.WriteString("_(Customer can track status at: " + statusLink + ")_")

	waNumber := whatsappNumber()
	waURL := "https://wa.me/" + waNumber + "?text=" + url.QueryEscape(sb.String())

	json200(w, map[string]string{
		"bookingId":  id,
		"waUrl":      waURL,
		"statusUrl":  statusLink,
		"ownerUrl":   ownerLink,
	})
}

// GET /api/booking/:id — poll booking status
func handleGetBooking(w http.ResponseWriter, r *http.Request, id string) {
	corsHeaders(w)
	mu.RLock()
	b, ok := bookings[id]
	mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json200(w, b)
}

// POST /api/booking/:id/accept — owner accepts original time
func handleOwnerAccept(w http.ResponseWriter, r *http.Request, id string) {
	corsHeaders(w)
	mu.Lock()
	b, ok := bookings[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	b.Status = StatusConfirmed
	b.History = append(b.History, HistoryEntry{At: now(), Actor: "owner", Action: "accepted"})
	saveBookings()
	mu.Unlock()

	// Build WhatsApp message to customer confirming booking
	var sb strings.Builder
	sb.WriteString("✅ *Booking Confirmed!*\n\n")
	sb.WriteString("Hi " + b.Name + ", your booking is confirmed.\n\n")
	sb.WriteString("*Service:* " + b.Category)
	if b.Service != "" {
		sb.WriteString(" – " + b.Service)
	}
	sb.WriteString("\n")
	sb.WriteString("*Date & Time:* " + b.Date + ", " + b.Slot + "\n")
	sb.WriteString("*Area:* " + b.Area + "\n\n")
	sb.WriteString("Our technician will arrive at the confirmed time. Thank you for choosing QuickFix Chennai! 🙏")

	userPhone := strings.TrimPrefix(b.Phone, "+")
	waURL := "https://wa.me/" + userPhone + "?text=" + url.QueryEscape(sb.String())

	json200(w, map[string]string{
		"status": "confirmed",
		"waUrl":  waURL, // owner opens this to message the customer
	})
}

// POST /api/booking/:id/reschedule — owner proposes a new time
func handleOwnerReschedule(w http.ResponseWriter, r *http.Request, id string) {
	corsHeaders(w)
	var req struct {
		Date string `json:"date"`
		Slot string `json:"slot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Date == "" || req.Slot == "" {
		http.Error(w, "date and slot required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	b, ok := bookings[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	b.Status = StatusOwnerProposed
	b.ProposedDate = req.Date
	b.ProposedSlot = req.Slot
	b.History = append(b.History, HistoryEntry{At: now(), Actor: "owner", Action: "proposed", Date: req.Date, Slot: req.Slot})
	saveBookings()
	mu.Unlock()

	base := siteBase(r)
	replyLink := base + "/reply/" + id

	// Build WhatsApp message to customer with new proposed time
	var sb strings.Builder
	sb.WriteString("🕐 *Time Change Proposal*\n\n")
	sb.WriteString("Hi " + b.Name + "! We'd like to suggest a different time for your booking.\n\n")
	sb.WriteString("*Service:* " + b.Category)
	if b.Service != "" {
		sb.WriteString(" – " + b.Service)
	}
	sb.WriteString("\n")
	sb.WriteString("*Your requested time:* " + b.Date + ", " + b.Slot + "\n")
	sb.WriteString("*Our proposed time:* " + req.Date + ", " + req.Slot + "\n\n")
	sb.WriteString("👇 *Please tap the link below to Accept or Suggest a different time:*\n")
	sb.WriteString(replyLink)

	userPhone := strings.TrimPrefix(b.Phone, "+")
	waURL := "https://wa.me/" + userPhone + "?text=" + url.QueryEscape(sb.String())

	json200(w, map[string]string{
		"status":    "proposed",
		"waUrl":     waURL,
		"replyLink": replyLink,
	})
}

// POST /api/booking/:id/user-accept — user accepts owner's proposed time
func handleUserAccept(w http.ResponseWriter, r *http.Request, id string) {
	corsHeaders(w)
	mu.Lock()
	b, ok := bookings[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Promote proposed time to confirmed time
	if b.ProposedDate != "" {
		b.Date = b.ProposedDate
		b.Slot = b.ProposedSlot
	}
	b.Status = StatusConfirmed
	b.History = append(b.History, HistoryEntry{At: now(), Actor: "user", Action: "accepted"})
	saveBookings()
	mu.Unlock()

	// Build WhatsApp message to owner confirming user accepted
	var sb strings.Builder
	sb.WriteString("✅ *Customer Accepted the New Time!*\n\n")
	sb.WriteString("Booking ID: " + id + "\n")
	sb.WriteString("*Customer:* " + b.Name + " (" + b.Phone + ")\n")
	sb.WriteString("*Confirmed time:* " + b.Date + ", " + b.Slot + "\n")
	sb.WriteString("*Area:* " + b.Area)

	waNumber := whatsappNumber()
	waURL := "https://wa.me/" + waNumber + "?text=" + url.QueryEscape(sb.String())

	json200(w, map[string]string{"status": "confirmed", "waUrl": waURL})
}

// POST /api/booking/:id/user-reschedule — user proposes their own new time
func handleUserReschedule(w http.ResponseWriter, r *http.Request, id string) {
	corsHeaders(w)
	var req struct {
		Date string `json:"date"`
		Slot string `json:"slot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Date == "" || req.Slot == "" {
		http.Error(w, "date and slot required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	b, ok := bookings[id]
	if !ok {
		mu.Unlock()
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	b.Status = StatusUserProposed
	b.ProposedDate = req.Date
	b.ProposedSlot = req.Slot
	b.History = append(b.History, HistoryEntry{At: now(), Actor: "user", Action: "proposed", Date: req.Date, Slot: req.Slot})
	saveBookings()
	mu.Unlock()

	base := siteBase(r)
	ownerLink := base + "/owner/" + id

	// Build WhatsApp message to owner with user's counter-proposal
	var sb strings.Builder
	sb.WriteString("🕐 *Customer Requested a Different Time*\n\n")
	sb.WriteString("Booking ID: " + id + "\n")
	sb.WriteString("*Customer:* " + b.Name + "\n")
	sb.WriteString("*Service:* " + b.Category)
	if b.Service != "" {
		sb.WriteString(" – " + b.Service)
	}
	sb.WriteString("\n")
	sb.WriteString("*Customer's requested time:* " + req.Date + ", " + req.Slot + "\n\n")
	sb.WriteString("👇 *Tap to Accept or Propose another time:*\n")
	sb.WriteString(ownerLink)

	waNumber := whatsappNumber()
	waURL := "https://wa.me/" + waNumber + "?text=" + url.QueryEscape(sb.String())

	json200(w, map[string]string{
		"status":    "user_proposed",
		"waUrl":     waURL,
		"ownerLink": ownerLink,
	})
}

// ── Router ────────────────────────────────────────────────────────────────────

func apiRouter(w http.ResponseWriter, r *http.Request) {
	corsHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")

	// /api/config
	if path == "/api/config" {
		handleConfig(w, r)
		return
	}

	// /api/booking  (create)
	if path == "/api/booking" {
		handleCreateBooking(w, r)
		return
	}

	// /api/booking/:id/*
	if strings.HasPrefix(path, "/api/booking/") {
		rest := strings.TrimPrefix(path, "/api/booking/")
		parts := strings.SplitN(rest, "/", 2)
		id := parts[0]
		action := ""
		if len(parts) == 2 {
			action = parts[1]
		}

		switch {
		case action == "" && r.Method == http.MethodGet:
			handleGetBooking(w, r, id)
		case action == "accept" && r.Method == http.MethodPost:
			handleOwnerAccept(w, r, id)
		case action == "reschedule" && r.Method == http.MethodPost:
			handleOwnerReschedule(w, r, id)
		case action == "user-accept" && r.Method == http.MethodPost:
			handleUserAccept(w, r, id)
		case action == "user-reschedule" && r.Method == http.MethodPost:
			handleUserReschedule(w, r, id)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	loadBookings()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/", apiRouter)
	mux.Handle("/", spaHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("QuickFix Chennai running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// spaHandler serves the embedded Vue SPA, falling back to index.html for unknown paths.
func spaHandler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if _, err := fs.Stat(sub, path); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
