package testhttp

import "encoding/json"

type TestRequest struct {
	base *TestRequest
	data map[string]interface{}
}

// Create a new end point test.
func New() *TestRequest {
	return &TestRequest{
		data: make(map[string]interface{}),
	}
}

// Set a base test.
func (test *TestRequest) SetBase(base *TestRequest) *TestRequest {
	test.base = base
	return test
}

// Returns base test
func (test *TestRequest) GetBase() *TestRequest {
	return test.base
}

// Creates a new instance based on current request.
func (test *TestRequest) GiveBirth() *TestRequest {
	new := New()
	new.SetBase(test)
	return new
}

// Set URL to request
func (req *TestRequest) SetRoute(url string) *TestRequest {
	req.data["route"] = url
	return req
}

// Returns route url
func (req *TestRequest) GetRoute() (url string, ok bool) {
	urlInterface, ok := req.data["route"]
	if !ok {
		return "", ok
	}
	url, ok = urlInterface.(string)
	return
}

// Set header data
func (req *TestRequest) SetHeaders(headers map[string]string) *TestRequest {
	req.data["headers"] = headers
	return req
}

// Returns header data
func (req *TestRequest) GetHeader() (headers map[string]string, ok bool) {
	objInterface, ok := req.data["headers"]
	if !ok {
		return map[string]string{}, ok
	}
	headers, ok = objInterface.(map[string]string)
	return
}

// Set body (JSON)
func (req *TestRequest) SetBodyJSON(data map[string]interface{}) *TestRequest {
	req.data["body"] = data
	return req
}

// Set body
func (req *TestRequest) SetBody(body string) *TestRequest {
	req.data["body"] = body
	return req
}

// Get body
func (req *TestRequest) GetBody(string) (body string, ok bool) {

	// Check if data exists
	objInterface, ok := req.data["body"]
	if !ok {
		return "", ok
	}

	// If data is string
	body, ok = objInterface.(string)
	if ok {
		return body, ok
	}

	// If data is JSON
	bodyJSON, ok := objInterface.(map[string]interface{})
	if !ok {
		// its not a json??
		return "", false
	}

	// Convert the map to a JSON
	bodyByte, err := json.Marshal(bodyJSON)
	if err != nil {
		return "", false
	}

	return string(bodyByte), ok

}

// Execute the test. Returns a new clean test object with the base of current test.
func (test *TestRequest) Execute(check func()) (newTest *TestRequest) {

	// Returns a new clean test object with the base of current test.
	new := New()
	new.SetBase(test.base)
	return new

}
