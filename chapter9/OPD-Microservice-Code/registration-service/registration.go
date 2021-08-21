package registration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"
	"strconv"

	"example.com/nats-microservices-opd/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nuid"
)

const (
	Version = "0.1.0"
)

// Server is a component.
type Server struct {
	*shared.Component
}

var ops uint64

// generateTokenNumber Generate token number for patient
func generateTokenNumber(start uint64) uint64 {
  if start > 0 {
	ops = start
	return ops
  }
  atomic.AddUint64(&ops, 1)
  return ops
}

// HandleTokenReset processes token reset requests.
func (s *Server) HandleTokenReset(w http.ResponseWriter, r *http.Request) {
	resetID, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	generateTokenNumber(resetID)
	json.NewEncoder(w).Encode("Token reset successful")
}

// HandleToken processes token generation requests for registered patients.
func (s *Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	token := generateTokenNumber(0)
	patientID, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println("Token %d generated for user %d", token, patientID)
	// Publish event to the NATS server
	nc := s.NATS()

	registration_event := shared.RegistrationEvent{patientID, token}
	reg_event, err := json.Marshal(registration_event)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("tokenID:%d - Publishing registration event with patientID %d\n", token, patientID)
	// Publishing the message to NATS Server
	nc.Publish("patient.register", reg_event)
	json.NewEncoder(w).Encode(registration_event)


}


// HandleRegister processes patient registration requests.
func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var registration *shared.RegistrationRequest
	err = json.Unmarshal(body, &registration)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert data to the database
	db := s.DB()

	insForm, err := db.Prepare("INSERT INTO patient_details(id, full_name, address, sex, phone, remarks) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(registration.ID, registration.FullName, registration.Address, registration.Sex, registration.Phone, registration.Remarks)

	// Tag the request with an ID for tracing in the logs.
	registration.RequestID = nuid.Next()
	fmt.Println(registration)

	// Publish event to the NATS server
	nc := s.NATS()

	//var registration_event shared.RegistrationEvent
	tokenNo := generateTokenNumber(0)
	registration_event := shared.RegistrationEvent{registration.ID, tokenNo}
	reg_event, err := json.Marshal(registration_event)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("requestID:%s - Publishing registration event with patientID %d\n", registration.RequestID, registration.ID)
	// Publishing the message to NATS Server
	nc.Publish("patient.register", reg_event)

	json.NewEncoder(w).Encode(registration_event)
}

// HandleUpdate processes requests to update patient details.
func (s *Server) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	//patientID := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request *shared.RegistrationRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db := s.DB()

	insForm, err := db.Prepare("UPDATE patient_details SET full_name=?, address=?, sex=?, phone=?, remarks=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(request.FullName, request.Address, request.Sex, request.Phone, request.Remarks, request.ID)

	json.NewEncoder(w).Encode("Record for Patient updated sucessfully")
}

// HandleView processes requests to view patient data.
func (s *Server) HandleView(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["id"]
	// Insert data to the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM patient_details WHERE ID=?", patientID)
    if err != nil {
        panic(err.Error())
    }

    registration := shared.RegistrationRequest{}
    for selDB.Next() {
        var id, phone int
        var full_name, address, sex, remarks string
        err = selDB.Scan(&id, &full_name, &address, &sex, &phone, &remarks)
        if err != nil {
            panic(err.Error())
        }
        registration.ID = id
        registration.FullName = full_name
        registration.Address = address
		registration.Sex = sex
		registration.Phone = phone
		registration.Remarks = remarks
    }

	fmt.Println(registration)
	json.NewEncoder(w).Encode(registration)
}

func (s *Server) HandleHomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Registration Service v%s\n", Version))
}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func (s *Server) ListenAndServe(addr string) error {

	r := mux.NewRouter()
	router := r.PathPrefix("/opd/patient/").Subrouter()

	// Handle base path requests
	// GET /opd/patient
	router.HandleFunc("/", s.HandleHomeLink)
	// Handle registration requests
	// POST /opd/patient/register
	router.HandleFunc("/register", s.HandleRegister).Methods("POST")

	// Handle update requests
	// PUT /opd/patient/update
	router.HandleFunc("/update", s.HandleUpdate).Methods("PUT")

	// Handle view requests
	// GET /opd/patient/view/{id}
	router.HandleFunc("/view/{id}", s.HandleView).Methods("GET")

	// Handle token requests
	// GET /opd/patient/token
	router.HandleFunc("/token/{id}", s.HandleToken).Methods("GET")

	// Handle token reset requests
	// GET /opd/patient/token/reset/{id}
	router.HandleFunc("/token/reset/{id}", s.HandleTokenReset).Methods("GET")

	//router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	//log.Fatal(http.ListenAndServe(":8080", router))

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

	go srv.Serve(l)

	return nil
}
