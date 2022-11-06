package webcam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Webcam struct {
	IP        string
	User      string
	Password  string
	token     string
	leaseTime int
}

type WebcamErrorResponse struct {
	Detail string `json:"detail"`
	Code   int    `json:"rspCode"`
}

type WebcamResponseToken struct {
	LeaseTime int    `json:"leaseTime"`
	Token     string `json:"name"`
}

type WebcamResponseValue struct {
	Token WebcamResponseToken `json:"Token"`
	State int                 `json:"state"`
}
type WebcamMotionResponseValue struct {
	State int `json:"state"`
}

type WebcamResponse struct {
	CMD          string              `json:"cmd"`
	Code         int                 `json:"code"`
	ErrorReponse WebcamErrorResponse `json:"error"`
	Value        WebcamResponseValue `json:"value"`
}

type WebcamMotionResponse struct {
	CMD          string                    `json:"cmd"`
	Code         int                       `json:"code"`
	ErrorReponse WebcamErrorResponse       `json:"error"`
	Value        WebcamMotionResponseValue `json:"value"`
}

// getToken Return webcam current login token
func (w Webcam) getToken() string {
	return w.token
}

// getToken Return webcam current login token
func (w Webcam) makeLoginRequest(client http.Client, url string, dataString string) (WebcamResponse, error) {
	var webcamResponses []WebcamResponse

	data := []byte(dataString)

	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if reqErr != nil {
		return WebcamResponse{}, reqErr
	}
	req.Header.Set("User-Agent", "github.com/a-castellano/Reolink-Manager")
	req.Header.Set("Content-Type", "application/json")

	resp, postErr := client.Do(req)
	if postErr != nil {
		return WebcamResponse{}, postErr
	}
	defer resp.Body.Close()

	body, readBodyErr := ioutil.ReadAll(resp.Body)
	if readBodyErr != nil {
		return WebcamResponse{}, readBodyErr
	}
	unmarshalError := json.Unmarshal([]byte(body), &webcamResponses)
	if unmarshalError != nil {
		return WebcamResponse{}, unmarshalError
	}

	webcamResponse := webcamResponses[0]

	return webcamResponse, nil
}

// Connect tries to login into webcam interface
func (w *Webcam) Connect(client http.Client) error {

	dataString := fmt.Sprintf("[{\"cmd\":\"Login\",\"action\":0,\"param\":{\"User\":{\"userName\":\"%s\",\"password\":\"%s\"}}}]", w.User, w.Password)
	url := fmt.Sprintf("http://%s/cgi-bin/api.cgi?cmd=Login&token=null", w.IP)

	webcamResponse, reponseErr := w.makeLoginRequest(client, url, dataString)

	//	fmt.Println(reponseErr)
	if reponseErr != nil {
		return reponseErr
	}

	if webcamResponse.Code != 0 {
		return errors.New("Login failed.")
	} else {
		w.token = webcamResponse.Value.Token.Token
		now := time.Now()
		w.leaseTime = int(now.Unix()) + webcamResponse.Value.Token.LeaseTime
	}

	return nil
}

// Checks if lease time has expired
func (w *Webcam) expiredToken() bool {

	var expired bool = false
	now := time.Now()
	nowSeconds := int(now.Unix())
	if w.leaseTime-nowSeconds < 10 {
		expired = true
	}
	return expired
}

// Retrieves motion sensor value
func (w *Webcam) MotionDetected(client http.Client) (bool, error) {

	var motionDetected bool = false
	if w.expiredToken() {
		w.Connect(client)
	}
	url := fmt.Sprintf("http://%s/cgi-bin/api.cgi?cmd=GetMdState&token=%s", w.IP, w.token)
	webcamResponse, reponseErr := w.makeLoginRequest(client, url, "")
	if reponseErr != nil {
		return motionDetected, reponseErr
	}
	if webcamResponse.Code != 0 {
		return motionDetected, errors.New("Motion detection Failed.")
	}
	if webcamResponse.Value.State != 0 {
		motionDetected = true
	}
	return motionDetected, nil
}

// Reboot webcam
func (w Webcam) Reboot(client http.Client) error {

	if w.getToken() == "" {
		return errors.New("Connect must be performed before rebooting the webcam.")
	}

	dataString := fmt.Sprintf("[{\"cmd\":\"Reboot\",\"action\":0,\"param\":{}}]")
	url := fmt.Sprintf("http://%s/cgi-bin/api.cgi?cmd=Reboot&token=%s", w.IP, w.token)

	webcamResponse, reponseErr := w.makeLoginRequest(client, url, dataString)

	if reponseErr != nil {
		return reponseErr
	}

	if webcamResponse.Code != 0 {
		errorString := fmt.Sprintf("Error rebooting webcam, %s.", webcamResponse.ErrorReponse.Detail)
		return errors.New(errorString)
	}

	return nil
}
