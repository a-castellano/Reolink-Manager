package webcam

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
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

func TestConnectSucceded(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "Login",
      "code" : 0,
      "value" : {
         "Token" : {
            "leaseTime" : 3600,
            "name" : "fef39ed8155f884"
         }
      }
   }
]
	`))}}}

	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass"}
	err := webcam.Connect(client)

	if err != nil {
		t.Errorf("Connect shouldn't fail.")
	}

	token := webcam.getToken()
	if token != "fef39ed8155f884" {
		t.Errorf("Token should be 'fef39ed8155f884', not '%s'.", token)
	}

	expiredToken := webcam.expiredToken()
	if expiredToken {
		t.Errorf("Token shouldn't be expired.")
	}

}

func TestRebootWithoutToken(t *testing.T) {

	var errMessage = "Connect must be performed before rebooting the webcam."
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "Reboot",
      "code" : 1,
      "error" : {
         "detail" : "please login first",
         "rspCode" : -6
      }
   }
]
	`))}}}

	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass"}
	err := webcam.Reboot(client)

	if err == nil {
		t.Errorf("Connect should fail.")
	}

	if err.Error() != errMessage {
		t.Errorf("Reboot error should be '%s', not '%s'.", errMessage, err.Error())
	}
}

func TestRebootFailedToken(t *testing.T) {

	var errMessage = "Error rebooting webcam, please login first."
	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "Reboot",
      "code" : 1,
      "error" : {
         "detail" : "please login first",
         "rspCode" : -6
      }
   }
]
	`))}}}

	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken"}
	err := webcam.Reboot(client)

	if err == nil {
		t.Errorf("Connect should fail.")
	} else {

		if err.Error() != errMessage {
			t.Errorf("Reboot error should be '%s', not '%s'.", errMessage, err.Error())
		}
	}
}

func TestReboot(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "Reboot",
      "code" : 0,
      "value" : {
         "rspCode" : 200
      }
   }
]
	`))}}}

	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken"}
	err := webcam.Reboot(client)

	if err != nil {
		t.Errorf("Connect shouldn't fail.")
	}
}

func TestNonExpiredToken(t *testing.T) {

	now := time.Now()
	nowSeconds := int(now.Unix()) + 500
	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken", leaseTime: nowSeconds}
	expire := webcam.expiredToken()

	if expire != false {
		t.Errorf("Token has not expired yet, expire should be false but it is true.")
	}
}

func TestExpiredToken(t *testing.T) {

	now := time.Now()
	nowSeconds := int(now.Unix()) - 500
	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken", leaseTime: nowSeconds}
	expire := webcam.expiredToken()

	if expire != true {
		t.Errorf("Token has expired, expire should be true, not false.")
	}
}

func TestMotionSensorDetection(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "GetMdState",
      "code" : 0,
      "value" : {
         "state" : 1
      }
   }
]
	`))}}}
	now := time.Now()
	nowSeconds := int(now.Unix()) + 500
	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken", leaseTime: nowSeconds}
	motion, err := webcam.MotionDetected(client)

	if err != nil {
		t.Errorf("MotionDetected shouldn't fail.")
	}

	if motion != true {
		t.Errorf("motion should be true, not false.")
	}
}

func TestMotionSensorNoDetection(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "GetMdState",
      "code" : 0,
      "value" : {
         "state" : 0
      }
   }
]
	`))}}}
	now := time.Now()
	nowSeconds := int(now.Unix()) + 500
	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken", leaseTime: nowSeconds}
	motion, err := webcam.MotionDetected(client)

	if err != nil {
		t.Errorf("MotionDetected shouldn't fail.")
	}

	if motion != false {
		t.Errorf("motion should be false, not true.")
	}
}

func TestMotionSensorErrorCode(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "GetMdState",
      "code" : 1,
      "value" : {
         "state" : 0
      }
   }
]
	`))}}}
	now := time.Now()
	nowSeconds := int(now.Unix()) + 500
	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken", leaseTime: nowSeconds}
	_, err := webcam.MotionDetected(client)

	if err == nil {
		t.Errorf("MotionDetected should fail.")
	}

}

func TestMotionSenssorReLogin(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMockTwoRequests{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "Login",
      "code" : 0,
      "value" : {
         "Token" : {
            "leaseTime" : 3600,
            "name" : "fef39ed8155f884"
         }
      }
   }
]
	`))},
		SecondResponse: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
[
   {
      "cmd" : "GetMdState",
      "code" : 0,
      "value" : {
         "state" : 0
      }
   }
]
	`))},
	}}
	now := time.Now()
	nowSeconds := int(now.Unix()) - 500
	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass", token: "testtoken", leaseTime: nowSeconds}
	_, err := webcam.MotionDetected(client)

	if err != nil {
		t.Errorf("MotionDetected shouldn't fail.")
	}

}

func TestConnectFailedBadJson(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
	[
   {s
   }
]
	`))}}}

	webcam := Webcam{IP: "10.10.0.1", User: "user", Password: "pass"}
	err := webcam.Connect(client)

	if err == nil {
		t.Errorf("Connect should fail.")
	}

}
