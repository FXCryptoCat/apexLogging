package influx

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/winterssy/sreq"
	"net/url"
	"strings"
	"time"
)

type DataSource interface {
	WriteRecord(tickline string)
	WriteRecords(ticklines []string)
}

//Used for posting messages directly to an influx db.
//TODO: Use interface to allow us to write to different dbs
type InfluxClient struct {
	influxHost       string
	influxUser       string
	influxPassword   string
}

func NewInfluxClient(host string, user string, password string) *InfluxClient {
	i := InfluxClient{influxHost: host, influxUser: user, influxPassword: password}

	timeout := 10 * time.Second
	sreq.SetTLSClientConfig( &tls.Config{InsecureSkipVerify: true})
	sreq.SetTimeout(timeout)

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
	log.Trace("InfluxClient::postData")

	uri := url.URL{
		Scheme:     "https",
		Path:       "write",
		User:       nil,
		Host:       client.influxHost,
		ForceQuery: false,
	}
	q := uri.Query()
	q.Set("db", "mydb") //<<---- Is this automatically created. If not do it.
	q.Set("u", client.influxUser)
	q.Set("p", client.influxPassword)
	uri.RawQuery = q.Encode()



	resp, err := sreq.
		Post(uri.String(),
			sreq.WithContent([]byte(tickData))).Text()
	fmt.Println(resp)
	if err != nil {
		log.Fatal("Unable to transmit to influxdb.", err)
	}

	//the sreq client must automatically close the response. Using go/http even forcing
	//the body to close was not reducing the file count on my raspberry.
	//use 'lsof -i PID | wc -l' to see how many connections are open.
	//defer resp.Body.Close()

}
