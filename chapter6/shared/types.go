package shared

// RegistrationRequest is the request to register a patient.
type RegistrationRequest struct {
	// Full Name of the patient.
	FullName string `json:"full_name,omitempty"`

	// Address of the patient.
	Address string `json:"address,omitempty"`

	// National Identification Number of the patient.
	ID int `json:"id"`

	// Sexual orientation 
	Sex string `json:"sex,omitempty"`

	// Sexual orientation 
	Email string `json:"email,omitempty"`	

	// Sexual orientation 
	Phone int `json:"phone,omitempty"`	

	// Other details
	Remarks string `json:"remarks,omitempty"`

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

type RegistrationEvent struct {
	// ID of the patient
	ID int `json:"id"`

	// Token of the patient
	Token uint64 `json:"token"`
}

// RegistrationRequest is the request to register a patient.
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

// TestRequest is the request to do a test on a patient.
type TestRequest struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time the inspection was done.
	Time string `json:"time,omitempty"`

	// Observations from the inspection.
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

// MedicationRequest is the request to report a medication instance on a patient.
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

// MedicationRequest is the request to report a medication instance on a patient.
type DischargeRequest struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time the inspection was done.
	Time string `json:"time,omitempty"`

	// Details of the dose of medication.
	State string `json:"test_name,omitempty"`	

	// Details of the dose of medication.
	PostMedication string `json:"post_medication,omitempty"`

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// Next visit(if any) date 
	NextVisit string `json:"next_visit,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// AdmissionEvent is the request to admit a patient.
type AdmissionEvent struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Medication schedule 
	Time string `json:"time,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	
}

// RegistrationRequest is the request to register a patient.
type ReleaseEvent struct {
	// National Identification Number of the patient.
	ID int `json:"id"`

	// Time of the release event
	Time string `json:"time"`

	// Medication schedule 
	NextState string `json:"next_state,omitempty"`

	// Tests to be carried out 
	PostMedication string `json:"post_medication,omitempty"`	

	// Special notes 
	Notes string `json:"notes,omitempty"`	

	// RequestID is the ID from the request.
	RequestID string `json:"request_id,omitempty"`
}

// RegistrationRequest is the request to register a patient.
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

