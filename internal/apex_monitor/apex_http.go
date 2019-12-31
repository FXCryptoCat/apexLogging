package apex_monitor

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type apexAuthResponse struct {
	ConnectSid string `json:"connect.sid"`
}

type apexAuthRequest struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	RememberMe string `json:"remember_me"`
}

func apexHttpRequestWithCookieAuth(url string, apexCookie string) (b []byte, err error) {
	log.Trace("requestWithCookies")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Failed request with cookies.", err)
		return
	}

	//Close the connection when done
	req.Close = true
	req.AddCookie(&http.Cookie{Name: "connect.sid", Value: apexCookie})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.WithField("Status", resp.StatusCode).WithField("Response", resp.Body).Warn("Request with cookies status code not 200.", err)
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return nil, err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes, nil
}
