package release

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"example.com/nats-microservices-opd/shared"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
	_ "github.com/go-sql-driver/mysql"
)

const (
	Version = "0.1.0"
)

// Server is a component.
type Server struct {
	*shared.Component
}


// ListenReleaseEvents Listen to release events and update the temporary table
func (s *Server) ListenReleaseEvents() error {
	nc := s.NATS()
	nc.Subscribe("patient.release", func(msg *nats.Msg) {
		var req *shared.ReleaseEvent
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		
		log.Printf("New Patient Release Event received for PatientID %d\n",
			req.ID)

			// Insert data to the database
		db := s.DB()

		insForm, err := db.Prepare("INSERT INTO release_reports(id, time, next_state, post_medication, notes) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(req.ID, req.Time, req.NextState, req.PostMedication, req.Notes)

	})

	return nil

}

// HandlePendingView processes requests to view pending releases.
func (s *Server) HandlePendingView(w http.ResponseWriter, r *http.Request) {
	// Retrieve pending inspections from the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM release_reports")
    if err != nil {
        panic(err.Error())
    }

	type allReleases []shared.ReleaseEvent
	var releases = allReleases{}

    for selDB.Next() {
		var newRelease shared.ReleaseEvent
        var id int
        var time, next_state, post_medication, notes string
        err = selDB.Scan(&id, &time, &next_state, &post_medication, &notes)
        if err != nil {
            panic(err.Error())
        }
        newRelease.ID = id
		newRelease.Time = time
        newRelease.NextState = next_state
		newRelease.PostMedication = post_medication
		newRelease.Notes = notes
		releases = append(releases, newRelease)
    }

	fmt.Println(releases)
	json.NewEncoder(w).Encode(releases)
}

// HandleDischargeRecord processes patient discharge requests.
func (s *Server) HandleDischargeRecord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var discharge *shared.DischargeRequest
	err = json.Unmarshal(body, &discharge)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert data to the database
	db := s.DB()

	insForm, err := db.Prepare("INSERT INTO discharge_details(id, time, state, post_medication, notes, next_visit) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(discharge.ID, discharge.Time, discharge.State, discharge.PostMedication, discharge.Notes, discharge.NextVisit)

	// Remove the entry from pending release table if it exists
	removeData, err := db.Prepare("DELETE FROM release_reports WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	removeData.Exec(discharge.ID)
    
	// Send admission  request if required
	if discharge.State == "admission" {
		discharge.RequestID = nuid.Next()
		// Publish event to the NATS server
		nc := s.NATS()

		//var registration_event shared.RegistrationEvent
		admission_event := shared.AdmissionEvent{discharge.ID, discharge.Time, discharge.Notes}
		reg_event, err := json.Marshal(admission_event)

		if err != nil {
			log.Fatal(err)
			return
		}

		log.Printf("requestID:%s - Publishing inspection event with patientID %d\n", discharge.RequestID, discharge.ID)
		// Publishing the message to NATS Server
		nc.Publish("patient.admission", reg_event)
	}

	json.NewEncoder(w).Encode("Patient discharge recorded successfully")
}

func (s *Server) HandleHomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Release Service v%s\n", Version))
}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func (s *Server) ListenAndServe(addr string) error {

	// Start listening to patient registration events
	s.ListenReleaseEvents()

	r := mux.NewRouter()
	router := r.PathPrefix("/opd/release/").Subrouter()

	// Handle base path requests
	// GET /opd/inspection
	router.HandleFunc("/", s.HandleHomeLink)
	
	// View pending patient release requests
	// POST /opd/release/pending
	router.HandleFunc("/pending", s.HandlePendingView).Methods("GET")

	// Handle discharge requests
	// GET /opd/treatment/tests/{id}
	router.HandleFunc("/discharge", s.HandleDischargeRecord).Methods("POST")

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
