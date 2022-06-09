package validators

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
