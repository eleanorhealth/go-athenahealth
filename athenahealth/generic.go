package athenahealth

type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ErrorMessageResponse struct {
	Message string `json:"errormessage"`
	Success bool   `json:"success"`
}

// StatusResponse
// https://docs.athenahealth.com/api/api-ref/appointment-booked#Appointment-Booked
// * This subroutine will return 1 on success, and will otherwise return an error message.
type StatusResponse struct {
	Status string `json:"status"`
}
