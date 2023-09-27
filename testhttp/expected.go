package testhttp

type TestExpected struct {
	data map[string]interface{}
}

func (expected *TestExpected) SetStatusCode(code int) *TestExpected {
	expected.data["status"] = code
	return expected
}

func (expected *TestExpected) GetStatusCode() (status int, ok bool) {
	obj, ok := expected.data["status"]
	if !ok {
		return 0, false
	}
	status, ok = obj.(int)
	return
}

func (expected *TestExpected) AddHeader() {

}
