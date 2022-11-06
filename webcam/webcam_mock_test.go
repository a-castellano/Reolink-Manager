package webcam

import (
	"net/http"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

type RoundTripperMockTwoRequests struct {
	Response       *http.Response
	SecondResponse *http.Response
	counter        int
	RespErr        error
}

func (rtm *RoundTripperMockTwoRequests) RoundTrip(*http.Request) (*http.Response, error) {
	if rtm.counter == 0 {
		rtm.counter++
		return rtm.Response, rtm.RespErr
	} else {
		return rtm.SecondResponse, rtm.RespErr
	}
}
