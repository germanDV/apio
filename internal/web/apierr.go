package web

// ApiErr is an error returned by the API
type ApiErr struct {
	Msg  string
	Code int
}

// To satisfy the error interface
func (ae ApiErr) Error() string {
	return ae.Msg
}
