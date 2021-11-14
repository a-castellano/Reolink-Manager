package webcam

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Webcam struct {
	IP       string
	User     string
	Password string
	token    string
}

type WebcamErrorResponse struct {
	Detail string `json:"detail"`
	Code   int    `json:"rspCode"`
}

type WebcamResponseToken struct {
	Detail string `json:"detail"`
	Code   int    `json:"rspCode"`
}

type WebcamResponseValue struct {
	Token WebcamResponseToken `json:"Token"`
}

type WebcamResponse struct {
	CMD          string              `json:"cmd"`
	Code         int                 `json:"code"`
	ErrorReponse WebcamErrorResponse `json:"error"`
	Value        WebcamResponseValue `json:"value"`
}

func (w Webcam) Connect(client http.Client) error {

	var webcamResponses []WebcamResponse

	dataString := fmt.Sprintf("{'cmd':'Login','action':0,'param':{'User':{'userName': '%s','password': '%s'}}}", w.User, w.Password)
	url := fmt.Sprintf("http://%s/cgi-bin/api.cgi?cmd=Login&token=null", w.IP)

	data := []byte(dataString)

	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("User-Agent", "github.com/a-castellano/Reolink-Manager")
	req.Header.Set("Content-Type", "application/json")

	resp, postErr := client.Do(req)
	if postErr != nil {
		return postErr
	}
	defer resp.Body.Close()

	body, readBodyErr := ioutil.ReadAll(resp.Body)
	if readBodyErr != nil {
		return readBodyErr
	}
	json.Unmarshal([]byte(body), &webcamResponses)

	webcamResponse := webcamResponses[0]

	if webcamResponse.Code != 0 {
		return errors.New("Login failed.")
	}

	return nil
}
