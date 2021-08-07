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

func main() {
	// Initialize loggers
	initLoggers()

	// Initialize tracing
	initTracing()

	// Create a connection to the NATS server over TLS with the RootCA
	nc, err := nats.Connect("tls://nats.example.com:4222", nats.Name("inspection-service-1"),
							nats.RootCAs("../ca.pem"),nats.UserInfo("inspection_service", "ins123"))
		
	if err != nil {
		ErrorLogger.Println(err)
		log.Fatal(err)
	}
	//defer nc.Close()

	InfoLogger.Println("Listening on [patient.register] subject")
	//fmt.Println("Listening on [patient.register] subject")

	// Subscribe
	nc.Subscribe("patient.register", func(msg *nats.Msg) {
		TraceLogger.Printf("Received on [%s]: '%s' \n", msg.Subject, string(msg.Data))
		//fmt.Printf("Received on [%s]: '%s' \n", msg.Subject, string(msg.Data))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		ErrorLogger.Println(err)
		log.Fatal(err)
	}

	// Start the health check HTTP service
	ListenAndServe("0.0.0.0:9000")

	runtime.Goexit()

}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func ListenAndServe(addr string) error {

	r := mux.NewRouter()
	router := r.PathPrefix("/subscriber/").Subrouter()
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
    file, err := os.OpenFile("subscriber.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
