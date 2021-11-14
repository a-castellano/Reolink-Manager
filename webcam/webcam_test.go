package webcam

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestConnectFailed(t *testing.T) {

	var errMessage = "Login failed."

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
	[
   {
      "cmd" : "Login",
      "code" : 1,
      "error" : {
         "detail" : "login failed",
         "rspCode" : -7
      }
   }
]
	`))}}}

	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass"}
	err := webcam.Connect(client)

	if err == nil {
		t.Errorf("Connect should fail.")
	}

	if errMessage != err.Error() {
		t.Errorf("Error message should be '%s', not '%s'.", errMessage, err)
	}

}
