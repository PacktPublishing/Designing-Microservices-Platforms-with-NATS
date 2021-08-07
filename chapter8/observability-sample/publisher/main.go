package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"net"
	"net/http"
	"time"
	"strconv"
	"github.com/nats-io/nats.go"
	"github.com/gorilla/mux"
	"github.com/nats-io/nuid"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    WarningLogger *log.Logger
    InfoLogger    *log.Logger
    ErrorLogger   *log.Logger
	TraceLogger   *log.Logger
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// RegistrationRequest contains data about the patient.
type RegistrationRequest struct {
	// Full Name of the patient.
	FullName string `json:"full_name,omitempty"`

	// Address of the patient.
	Address string `json:"address,omitempty"`

	// National Identification Number of the patient.
	ID int `json:"id"`

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

func main() {
	// Initialize loggers
	initLoggers()

	// Initialize Tracing
	initTracing()

	// Create a connection to the NATS server over TLS with the RootCA
	nc, err := nats.Connect("tls://nats.example.com:4222", nats.Name("registration-service"),
							nats.RootCAs("../ca.pem"), nats.UserInfo("registration_service", "reg123"))
	if err != nil {
		ErrorLogger.Println(err)
		log.Fatal(err)
	}
	//defer nc.Close()

	// Create event
	regEvent := RegistrationRequest{"Chanaka Fernando", "44 Seeduwa", 1111, nuid.Next()}
	reg_event, err := json.Marshal(regEvent)

	if err != nil {
		ErrorLogger.Println(err)
		log.Fatal(err)
	}

	InfoLogger.Printf("Publishing message with ID %s", regEvent.RequestID)

	// Publish a message on "patient.profile" subject
	subj, msg := "patient.register", reg_event
	nc.Publish(subj, msg)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		ErrorLogger.Println(err)
		log.Fatal(err)
	} else {
		TraceLogger.Printf("Published [%s] : '%s'\n", subj, msg)
		//log.Printf("Published [%s] : '%s'\n", subj, msg)
	}

	// Start the health check HTTP service
	ListenAndServe("0.0.0.0:9001")

	runtime.Goexit()
}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func ListenAndServe(addr string) error {

	r := mux.NewRouter()
	
	router := r.PathPrefix("/publisher/").Subrouter()
	router.Use(prometheusMiddleware)

	// Handle health check requests
	router.HandleFunc("/healthz", HandleHealthCheck)

	router.Path("/metrics").Handler(promhttp.Handler())

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	InfoLogger.Println("Starting the health check service on %s", addr)
	go srv.Serve(l)

	return nil
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Service available. \n"))
}

func initLoggers() {
    file, err := os.OpenFile("publisher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	TraceLogger = log.New(file, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func initTracing() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}