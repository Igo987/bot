package api

import (
	"context"
	"encoding/json"
	"github/Igo87/crypt/config"
	"github/Igo87/crypt/pkg/logger"
	"github/Igo87/crypt/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock.go -package=mocks
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Run(ctx context.Context, h http.Handler)
}

type handler struct {
	srv *service.Service
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveGet(w, r)
}

func (h *handler) serveGet(w http.ResponseWriter, r *http.Request) {
	mux := mux.NewRouter()
	mux.HandleFunc("/api/coins", h.handleCoinsGet).Methods(http.MethodGet)
	mux.HandleFunc("/api/coins/{name}", h.handleCoinsByNameGet).Methods(http.MethodGet)
	mux.ServeHTTP(w, r)
}

func (h *handler) Run(ctx context.Context) {
	addr := config.Cfg.GetPort()
	srv := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.LogStart().Error("Failed to start server: %v", err)
		}
	}()
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Cfg.Waiting)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.LogStart().Error("Failed to shutdown server: %v", err)
	}
}

func (h *handler) handleCoinsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := h.srv.GetDataByToday(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.encodeResponse(w, http.StatusOK, data)
}

func (h *handler) handleCoinsByNameGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	data, err := h.srv.GetDataByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if data == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	h.encodeResponse(w, http.StatusOK, data)
}

// encodeResponse encodes the provided data into JSON format and writes it to the http.ResponseWriter.
//
// Parameters:
// - w: http.ResponseWriter - the writer to write the response to.
// - status: int - the HTTP status code to set in the response.
// - data: interface{} - the data to be encoded into JSON.
//
// Return: None.
func (h *handler) encodeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	flusher, _ := w.(http.Flusher)
	flusher.Flush()
}

// NewHandler returns a new instance of the Handler interface.
func NewHandler(srv *service.Service) *handler {
	return &handler{srv: srv}
}
