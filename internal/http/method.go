package http

type RequestMethod string

const (
	Get     RequestMethod = "GET"
	Put     RequestMethod = "PUT"
	Post    RequestMethod = "POST"
	Patch   RequestMethod = "PATCH"
	Delete  RequestMethod = "DELETE"
	Head    RequestMethod = "HEAD"
	Options RequestMethod = "OPTIONS"
)
