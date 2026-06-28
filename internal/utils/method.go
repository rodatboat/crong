package utils

import "github.com/rodatboat/crong/internal/entities"

// ReqMethodToHTTPMethod converts a ReqMethod enum to its HTTP method string representation
func ReqMethodToHTTPMethod(method entities.ReqMethod) string {
	methods := map[entities.ReqMethod]string{
		entities.GET:     "GET",
		entities.HEAD:    "HEAD",
		entities.POST:    "POST",
		entities.PUT:     "PUT",
		entities.PATCH:   "PATCH",
		entities.DELETE:  "DELETE",
		entities.CONNECT: "CONNECT",
		entities.OPTIONS: "OPTIONS",
		entities.TRACE:   "TRACE",
	}

	if httpMethod, exists := methods[method]; exists {
		return httpMethod
	}
	return ""
}

// HTTPMethodToReqMethod converts an HTTP method string to a ReqMethod enum
func HTTPMethodToReqMethod(method string) (entities.ReqMethod, bool) {
	methods := map[string]entities.ReqMethod{
		"GET":     entities.GET,
		"HEAD":    entities.HEAD,
		"POST":    entities.POST,
		"PUT":     entities.PUT,
		"PATCH":   entities.PATCH,
		"DELETE":  entities.DELETE,
		"CONNECT": entities.CONNECT,
		"OPTIONS": entities.OPTIONS,
		"TRACE":   entities.TRACE,
	}

	if reqMethod, exists := methods[method]; exists {
		return reqMethod, true
	}
	return 0, false
}
