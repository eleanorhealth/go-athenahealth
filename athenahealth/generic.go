package athenahealth

type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ErrorMessageResponse struct {
	Message string `json:"errormessage"`
	Success bool   `json:"success"`
}
