package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// The built Vue site is embedded so the whole app ships as one binary.
// `all:` keeps files that start with "_" (Vite sometimes emits them).
//
//go:embed all:dist
var distFS embed.FS

// config.json is embedded as a fallback; at runtime we prefer the file on disk
// so the owner can edit content without rebuilding.
//
//go:embed config.json
var embeddedConfig []byte

type booking struct {
	Category string `json:"category"`
	Service  string `json:"service"`
	Price    string `json:"price"`
	Date     string `json:"date"`
	Slot     string `json:"slot"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Area     string `json:"area"`
	Coords   string `json:"coords"`
	MapsLink string `json:"mapsLink"`
	Notes    string `json:"notes"`
	Urgent   bool   `json:"urgent"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/config", handleConfig)
	mux.HandleFunc("/api/booking", handleBooking)
	mux.Handle("/", spaHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("QuickFix Chennai running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// GET /api/config — prefer config.json on disk, fall back to embedded copy.
func handleConfig(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		data = embeddedConfig
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(data)
}

// POST /api/booking — validate, append to bookings.json (best-effort), return wa.me URL.
func handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var b booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(b.Name) == "" || strings.TrimSpace(b.Phone) == "" || strings.TrimSpace(b.Area) == "" {
		http.Error(w, "name, phone and area are required", http.StatusUnprocessableEntity)
		return
	}

	logBooking(b) // best-effort; ignores errors (disk may be read-only/ephemeral)
	waNumber := whatsappNumber()
	msg := buildMessage(b)
	waURL := "https://wa.me/" + waNumber + "?text=" + url.QueryEscape(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"waUrl": waURL, "message": msg})
}

func buildMessage(b booking) string {
	var sb strings.Builder
	if b.Urgent {
		sb.WriteString("🚨 URGENT Booking\n\n")
	} else {
		sb.WriteString("🛠️ New Booking\n\n")
	}
	if b.Category != "" {
		sb.WriteString("Service: " + b.Category)
		if b.Service != "" {
			sb.WriteString(" – " + b.Service)
		}
		sb.WriteString("\n")
	}
	if b.Price != "" {
		sb.WriteString("Est. price: " + b.Price + "\n")
	}
	if b.Date != "" || b.Slot != "" {
		sb.WriteString("Preferred: " + strings.TrimPrefix(b.Date+", "+b.Slot, ", ") + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString("Name: " + b.Name + "\n")
	sb.WriteString("Phone: " + b.Phone + "\n")
	sb.WriteString("Area: " + b.Area + "\n")
	if b.Coords != "" {
		sb.WriteString("📍 GPS: " + b.Coords + "\n")
	}
	if b.MapsLink != "" {
		sb.WriteString("🗺️ Map: " + b.MapsLink + "\n")
	}
	if b.Notes != "" {
		sb.WriteString("Notes: " + b.Notes + "\n")
	}
	sb.WriteString("\n(Sent from website — please confirm the time slot.)")
	return sb.String()
}

// Append one booking as a JSON line to bookings.json. Best-effort log only.
func logBooking(b booking) {
	f, err := os.OpenFile("bookings.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	rec := map[string]any{"at": time.Now().Format(time.RFC3339), "booking": b}
	line, _ := json.Marshal(rec)
	f.Write(append(line, '\n'))
}

// Pull the WhatsApp number out of the active config so it stays in one place.
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

// Serve the embedded SPA, falling back to index.html for unknown paths.
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
			// Not a real file → serve index.html (SPA routing).
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
