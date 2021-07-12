package treatment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"example.com/nats-microservices-opd/shared"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

const (
	Version = "0.1.0"
)

// Server is a component.
type Server struct {
	*shared.Component
}

// ListenTreatmentEvents listens to events coming from inspection service
func (s *Server) ListenTreatmentEvents() error {
	nc := s.NATS()
	nc.Subscribe("patient.treatment", func(msg *nats.Msg) {
		var req *shared.InspectionEvent
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		
		log.Printf("New Patient Inspection Event received for PatientID %d\n",
			req.ID)

		// Insert data to the database
		db := s.DB()

		insForm, err := db.Prepare("INSERT INTO inspection_reports(id, medication, tests, notes) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(req.ID, req.Medication, req.Tests, req.Notes)

	})

	return nil

}

// HandlePendingView processes requests to view pending treatments.
func (s *Server) HandlePendingView(w http.ResponseWriter, r *http.Request) {
	// Retrieve pending inspections from the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM inspection_reports")
    if err != nil {
        panic(err.Error())
    }

	type allTreatments []shared.InspectionEvent
	var treatments = allTreatments{}

    for selDB.Next() {
		var newTreatment shared.InspectionEvent
        var id int
        var medication, tests, notes string
        err = selDB.Scan(&id, &medication, &tests, &notes)
        if err != nil {
            panic(err.Error())
        }
        newTreatment.ID = id
        newTreatment.Medication = medication
		newTreatment.Tests = tests
		newTreatment.Notes = notes
		treatments = append(treatments, newTreatment)
    }

	fmt.Println(treatments)
	json.NewEncoder(w).Encode(treatments)
}

// HandleRelease processes requests to initiate a patient release.
func (s *Server) HandleRelease(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var release *shared.ReleaseEvent
	err = json.Unmarshal(body, &release)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Tag the request with an ID for tracing in the logs.
	release.RequestID = nuid.Next()
	fmt.Println(release)

	// Publish event to the NATS server
	nc := s.NATS()
	
	release.RequestID = nuid.Next()
	release_event := shared.ReleaseEvent{release.ID, release.Time, release.NextState, release.PostMedication, release.Notes, release.RequestID}
	rel_event, err := json.Marshal(release_event)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("requestID:%s - Publishing inspection event with patientID %d\n", release.RequestID, release.ID)
	// Publishing the message to NATS Server
	nc.Publish("patient.release", rel_event)

	json.NewEncoder(w).Encode("Release event published")
}

// HandleTestRecord processes recording of tests related requests.
func (s *Server) HandleTestRecord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var test *shared.TestRequest
	err = json.Unmarshal(body, &test)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert data to the database
	db := s.DB()

	insForm, err := db.Prepare("INSERT INTO test_reports(id, time, test_name, results, status, notes) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(test.ID, test.Time, test.TestName, test.Results, test.Status, test.Notes)

	json.NewEncoder(w).Encode("Test recorded successfully")
}


// HandleMedicationRecord processes patient medication record requests.
func (s *Server) HandleMedicationRecord(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var medication *shared.MedicationRequest
	err = json.Unmarshal(body, &medication)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insert data to the database
	db := s.DB()

	insForm, err := db.Prepare("INSERT INTO medication_reports(id, time, dose, notes) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(medication.ID, medication.Time, medication.Dose, medication.Notes)

	// Remove the entry from pending medication table if it exists
	removeData, err := db.Prepare("DELETE FROM inspection_reports WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	removeData.Exec(medication.ID)

	json.NewEncoder(w).Encode("Record updated successfully")
}

// HandleTestView processes requests to view test data.
func (s *Server) HandleTestView(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["id"]
	// Insert data to the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM test_reports WHERE ID=?", patientID)
    if err != nil {
        panic(err.Error())
    }

	type allReports []shared.TestRequest
	var reports = allReports{}

    for selDB.Next() {
		var newReport shared.TestRequest
        var id int
        var time, test_name, results, status, notes string
        err = selDB.Scan(&id, &time, &test_name, &results, &status, &notes)
        if err != nil {
            panic(err.Error())
        }
        newReport.ID = id
        newReport.Time = time
        newReport.TestName = test_name
		newReport.Results = results
		newReport.Status = status
		newReport.Notes = notes
		reports = append(reports, newReport)
    }

	fmt.Println(reports)
	json.NewEncoder(w).Encode(reports)
    //defer db.Close()
}

// HandleHistoryView processes requests to view medication history data.
func (s *Server) HandleHistoryView(w http.ResponseWriter, r *http.Request) {
	patientID := mux.Vars(r)["id"]
	// Select data from the database
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM medication_reports WHERE ID=?", patientID)
    if err != nil {
        panic(err.Error())
    }

	type allMedications []shared.MedicationRequest
	var medications = allMedications{}

    for selDB.Next() {
		var newMedication shared.MedicationRequest
        var id int
        var time, dose, notes string
        err = selDB.Scan(&id, &time, &dose, &notes)
        if err != nil {
            panic(err.Error())
        }
        newMedication.ID = id
        newMedication.Time = time
        newMedication.Dose = dose
		newMedication.Notes = notes
		medications = append(medications, newMedication)
    }

	fmt.Println(medications)
	json.NewEncoder(w).Encode(medications)
}

func (s *Server) HandleHomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Treatment Service v%s\n", Version))
}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func (s *Server) ListenAndServe(addr string) error {

	// Start listening to patient registration events
	s.ListenTreatmentEvents()

	r := mux.NewRouter()
	router := r.PathPrefix("/opd/treatment/").Subrouter()

	// Handle base path requests
	// GET /opd/inspection
	router.HandleFunc("/", s.HandleHomeLink)
	
	// Handle medication record requests
	// POST /opd/treatment/medication/{id}
	router.HandleFunc("/medication", s.HandleMedicationRecord).Methods("POST")

	// Handle test result update requests
	// GET /opd/treatment/tests/{id}
	router.HandleFunc("/tests", s.HandleTestRecord).Methods("POST")

	// Handle patient release initialization requests
	// GET /opd/treatment/release
	router.HandleFunc("/release", s.HandleRelease).Methods("POST")

	// Handle test result view requests
	// GET /opd/treatment/tests/{id}
	router.HandleFunc("/tests/{id}", s.HandleTestView).Methods("GET")

	// Handle medication history view requests
	// GET /opd/treatment/history/{id}
	router.HandleFunc("/history/{id}", s.HandleHistoryView).Methods("GET")

	// Handle pending treatment list view requests
	// GET /opd/treatment/pending
	router.HandleFunc("/pending", s.HandlePendingView).Methods("GET")

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
