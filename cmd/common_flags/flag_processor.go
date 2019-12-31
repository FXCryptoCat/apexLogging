package common_flags

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type CommonFlags struct {
	Delay             *int   `yaml:"delay,omitempty"`
	LogLevel          string `yaml:"log_level,omitempty"`
	ApexIp            string `yaml:"apex_ip,omitempty"`
	ApexCookie        string `yaml:"apex_cookie,omitempty"`
	InfluxUser        string `yaml:"influx_user,omitempty"`
	InfluxPassword    string `yaml:"influx_password,omitempty"`
	InfluxIp          string `yaml:"influx_ip,omitempty"`
	DisableDataSource *bool   `yaml:"disableDataSource,omitempty"`
}

func GetFlags() CommonFlags {
	delayPtr := flag.Int("delay", 30, `delay seconds between requests to apex`)
	logLevel := flag.String("loglevel", "warn", `trace, debug, info, warn, error, fatal`)
	apexIp := flag.String("apexIp", "", `The IP address for your APEX`)
	apexCookie := flag.String("apexCookie", "", `The cookie used to authenticate your requests`)
	influxUser := flag.String("influxUser", "", `InfluxDb User Name`)
	influxPassword := flag.String("influxPassword", "", `InfluxDb Password`)
	influxIp := flag.String("influxIp", "", `InfluxDb IP`)
	configFile := flag.String("configFile", "", `Config file. Overrides command line settings`)
	disableDataSource := flag.Bool("disableDataSource", false, `Disables writes to the data source`)

	flag.Parse()

	allParametersSpecified := len(*apexIp) > 0 || len(*apexCookie) > 0 || len(*influxUser) > 0 || len(*influxPassword) > 0 || len(*influxIp) > 0
	if len(*configFile) == 0 && !allParametersSpecified {
		log.Error("You must specify a config file or all of the command line arguments)")
		flag.Usage()
		os.Exit(1)
	}

	//Todo: Add validation
	//Set any commandline arguments and override them with the conifg file. They can be mixed
	//config file takes precedence.
	commonFlags := CommonFlags{
		Delay:          delayPtr,
		LogLevel:       *logLevel,
		ApexIp:         *apexIp,
		ApexCookie:     *apexCookie,
		InfluxUser:     *influxUser,
		InfluxPassword: *influxPassword,
		InfluxIp:       *influxIp,
		DisableDataSource: disableDataSource,
	}

	if len(*configFile) > 0 {
		loadConfigFile(*configFile, &commonFlags)
	}

	if commonFlags.Delay == nil {
		*delayPtr = 60
	}
	if len(commonFlags.LogLevel) == 0 {
		*logLevel = "warn"
	}
	//TODO: More validation

	return commonFlags
}

//Assumes it is in the local directory
func loadConfigFile(fileName string, commonFlags *CommonFlags) {
	yamlFile, err := ioutil.ReadFile("./" + fileName)
	if err != nil {
		log.Fatalf("yamlFile.Get %s err   #%v ", fileName, err)
	}

	var flags CommonFlags
	err = yaml.Unmarshal(yamlFile, &flags)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	//Combine the flags - one liners would be nice
	if flags.Delay != nil {
		commonFlags.Delay = flags.Delay
	}
	if len(flags.LogLevel) > 0 {
		commonFlags.LogLevel = flags.LogLevel
	}
	if len(flags.ApexIp) > 0 {
		commonFlags.ApexIp = flags.ApexIp
	}
	if len(flags.ApexCookie) > 0 {
		commonFlags.ApexCookie = flags.ApexCookie
	}
	if len(flags.InfluxUser) > 0 {
		commonFlags.InfluxUser = flags.InfluxUser
	}
	if len(flags.InfluxPassword) > 0 {
		commonFlags.InfluxPassword = flags.InfluxPassword
	}
	if len(flags.InfluxIp) > 0 {
		commonFlags.InfluxIp = flags.InfluxIp
	}
	if flags.DisableDataSource != nil {
		commonFlags.DisableDataSource = flags.DisableDataSource
	}


}
