package shared

// RegistrationRequest contains data about the patient.
type RegistrationRequest struct {
	// Full Name of the patient.
	FullName string `json:"full_name,omitempty"`

	// Address of the patient.
	Address string `json:"address,omitempty"`

	// National Identification Number of the patient.
	ID int `json:"id"`

	// Sexual orientation 
	Sex string `json:"sex,omitempty"`

	// Email address 
	Email string `json:"email,omitempty"`	

	// Phone number 
	Phone int `json:"phone,omitempty"`	

	// Other details
	Remarks string `json:"remarks,omitempty"`

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// RegistrationEvent contains the details for a given registration instance
type RegistrationEvent struct {
	// ID of the patient
	ID int `json:"id"`

	// Token of the patient
	Token uint64 `json:"token"`
}

// InspectionRequest contains data related to patient inspection.
type InspectionRequest struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time the inspection was done.
	Time string `json:"time,omitempty"`

	// Observations from the inspection.
	Observations string `json:"observations,omitempty"`

	// Medication schedule 
	Medication string `json:"medication,omitempty"`

	// Tests to be carried out 
	Tests string `json:"tests,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// TestRequest contains data related to tests carried out for patients.
type TestRequest struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time the inspection was done.
	Time string `json:"time,omitempty"`

	// Name of the test.
	TestName string `json:"test_name,omitempty"`

	// Test results 
	Results string `json:"results,omitempty"`

	// Status of the test 
	Status string `json:"status,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// MedicationRequest contains data related to medication of a patient.
type MedicationRequest struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time the inspection was done.
	Time string `json:"time,omitempty"`

	// Details of the dose of medication.
	Dose string `json:"test_name,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// DischargeRequest contains details of patient discharge.
type DischargeRequest struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time the inspection was done.
	Time string `json:"time,omitempty"`

	// State of the discharge.
	State string `json:"test_name,omitempty"`	

	// Details of the medication after release.
	PostMedication string `json:"post_medication,omitempty"`

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// Next visit(if any) date 
	NextVisit string `json:"next_visit,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// AdmissionEvent contains data on patient admission to the hospital.
type AdmissionEvent struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time of the admission request 
	Time string `json:"time,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	
}

// ReleaseEvent contains data on the patient release.
type ReleaseEvent struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time of the release event
	Time string `json:"time"`

	// NextState of the patient 
	NextState string `json:"next_state,omitempty"`

	// Medication after release 
	PostMedication string `json:"post_medication,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// InspectionEvent contains data on inspection activities.
type InspectionEvent struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Medication schedule 
	Medication string `json:"medication,omitempty"`

	// Tests to be carried out 
	Tests string `json:"tests,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	
}

