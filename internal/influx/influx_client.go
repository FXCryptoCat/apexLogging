package influx

import (
	"bytes"
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

//Used for posting messages directly to an influx db.
//TODO: Use interface to allow us to write to different dbs
type InfluxClient struct {
	influxHost string
	influxUser string
	influxPassword string
}

func NewInfluxClient(host string, user string, password string) *InfluxClient {
	i := InfluxClient{influxHost:host,influxUser:user,influxPassword:password}
	return &i
}

//Writes a single tick record to influx
func (client *InfluxClient) WriteRecord(tickline string) {
	client.postData(tickline)
}

//Batch up multiple tick lines by joining them with a \n.
func (client *InfluxClient) WriteRecords(ticklines []string) {
	tickData := strings.Join(ticklines, "\n")
	client.postData(tickData)
}

func (client *InfluxClient) postData(tickData string) {
	logrus.Trace("InfluxClient::postData")
	//logrus.Debug(tickData)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	influxEndPoint := &http.Client{Transport: tr}

	uri := url.URL{
		Scheme:     "https",
		Path:       "write",
		User:       nil,
		Host:       client.influxHost,
		ForceQuery: false,
	}

	q := uri.Query()
	q.Set("db", "mydb")  //<<---- Is this automatically created. If not do it.
	q.Set("u", client.influxUser)
	q.Set("p", client.influxPassword)
	uri.RawQuery = q.Encode()

	resp, err := influxEndPoint.Post(uri.String(), "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(tickData)))
	if err != nil {
		logrus.Fatal("Unable to transmit to influxdb.", err)
	}
	defer resp.Body.Close()

}
