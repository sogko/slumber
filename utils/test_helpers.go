package utils

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
)

// MapFromJSON is a test helper function that decodes recorded response body to
// a specific struct type
// Note: this functions panics on error. For test usage only, not for production.
func MapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(fmt.Sprintf("mapFromJSON(): Not a valid JSON body\n%v", string(data)))
	}
	return result.(map[string]interface{})
}

// DecodeResponseToType is a test helper function that decodes recorded response body to
// a specific struct type
// Note: this functions panics on error. For test usage only, not for production.
func DecodeResponseToType(recorder *httptest.ResponseRecorder, target interface{}) map[string]interface{} {
	decoder := json.NewDecoder(recorder.Body)
	err := decoder.Decode(target)
	if err != nil {
		panic(fmt.Sprintf("DecodeResponseToType(): Not a valid JSON body\n%v", recorder.Body.String()))
	}
	return nil
}
