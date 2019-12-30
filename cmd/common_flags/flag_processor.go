package common_flags

import (
	"flag"
	"os"
)

type CommonFlags struct {
	Delay int
	LogLevel string
	ApexIp string

	//So digging this out it fun. It usually doesn't change.
	//Use chrome. Sign into the classic dashboard, not fusion.
	//Open the settings
	//Go to advanced
	//Go to Site Settings
	//Go to Cookies and Site Data
	//Hard to find but "See all cookies and site data"
	//Search for your apex's IP
	//You want the connect.sid, the Content is what you put in your cookie (Usually)
	//-- I just went thru this and the cookie value I'm using is not what is in chrome.
	//TODO: Hmmm if giving the apex username and password can we get a new cookie when the old one expires
	ApexCookie string

	InfluxUser string
	InfluxPassword string
	InfluxIp string


}

func GetFlags() CommonFlags{
	delayPtr := flag.Int("delay", 30, `delay seconds between requests to apex`)
	logLevel := flag.String("loglevel", "warn", `trace, debug, info, warn, error, fatal`)
	apexIp := flag.String("apexIp", "", `The IP address for your APEX`)
	apexCookie := flag.String("apexCookie", "", `The cookie used to authenticate your requests`)
	influxUser := flag.String("influxUser", "", `InfluxDb User Name`)
	influxPassword := flag.String("influxPassword", "", `InfluxDb Password`)
	influxIp := flag.String("influxIp", "", `InfluxDb IP`)

	flag.Parse()

	if len(*apexIp)  == 0 || len(*apexCookie)==0 || len(*influxUser)==0 || len(*influxPassword) == 0 || len(*influxIp)==0 {
		flag.Usage()
		os.Exit(1)
	}

	//Todo: Add validation
	if delayPtr == nil {
		*delayPtr = 60
	}
	if logLevel == nil {
		*logLevel = "warn"
	}

	commonFlags := CommonFlags{
		Delay:          *delayPtr,
		LogLevel:       *logLevel,
		ApexIp:         *apexIp,
		ApexCookie:     *apexCookie,
		InfluxUser:     *influxUser,
		InfluxPassword: *influxPassword,
		InfluxIp:       *influxIp,
	}

	return commonFlags
}