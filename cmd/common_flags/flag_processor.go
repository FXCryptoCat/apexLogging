package common_flags

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

type CommonFlags struct {
	Delay             *int   `yaml:"delay,omitempty"`
	LogLevel          string `yaml:"log_level,omitempty"`
	ApexIp            string `yaml:"apex_ip,omitempty"`
	ApexUserName      string `yaml:"apex_username,omitempty"`
	ApexPassword      string `yaml:"apex_password,omitempty"`
	InfluxUser        string `yaml:"influx_user,omitempty"`
	InfluxPassword    string `yaml:"influx_password,omitempty"`
	InfluxIp          string `yaml:"influx_ip,omitempty"`
	DisableDataSource *bool  `yaml:"disableDataSource,omitempty"`
	Install           *bool  `yaml:"install,omitempty"`
	Quiet             *bool  `yaml:"quiet,omitempty"`
	LogFile           string `yaml:"log,omitempty"`
	Gid               *int   `yaml:"Gid,omitempty"`
	Uid               *int   `yaml:"Uid,omitempty"`
	ExecDir           string `yaml:"execDir,omitempty"`
	OsUserName        string `yaml:"userName,omitempty"`
	TickLogFile       string `yaml:"tickLogFile,omitempty"`
}

func GetFlags() CommonFlags {
	delayPtr := flag.Int("delay", 30, `delay seconds between requests to apex`)
	logLevel := flag.String("loglevel", "warn", `trace, debug, info, warn, error, fatal`)
	apexIp := flag.String("apexIp", "", `The IP address for your APEX`)
	apexUserName := flag.String("apexUserName", "", `The apex user name`)
	apexPassword := flag.String("apexPassword", "", `The apex password`)
	influxUser := flag.String("influxUser", "", `InfluxDb User Name`)
	influxPassword := flag.String("influxPassword", "", `InfluxDb Password`)
	influxIp := flag.String("influxIp", "", `InfluxDb IP`)
	configFile := flag.String("config", "", `Config file. Overrides command line settings`)
	disableDataSource := flag.Bool("disableDataSource", false, `Disables writes to the data source`)
	install := flag.Bool("install", false, `Install as a service. The command line arguments, or specified config file, will be used to populate the config file.`)
	quiet := flag.Bool("quiet", false, `Answers any prompts with a default value.`)
	logFile := flag.String("log", "", `Log file. If not specified stdout will be used`)
	execDir := flag.String("execDir", "/usr/bin", `Location where Install will copy exec and config file. Default /usr/bin`)
	osUserName := flag.String("osUserName", "", `User name to run service as, only used by install`)
	tickLogFile := flag.String("tickLogFile", "", `Where to store tick data locally (Optional)`)

	flag.Parse()

	allParametersSpecified := len(*apexIp) > 0 || len(*apexUserName) > 0 || len(*apexPassword) > 0 || len(*influxUser) > 0 || len(*influxPassword) > 0 || len(*influxIp) > 0
	if len(*configFile) == 0 && !allParametersSpecified {
		log.Error("You must specify a config file or all of the command line arguments)")
		flag.Usage()
		os.Exit(1)
	}

	//Todo: Add validation
	//Set any commandline arguments and override them with the conifg file. They can be mixed
	//config file takes precedence.
	commonFlags := CommonFlags{
		Delay:             delayPtr,
		LogLevel:          *logLevel,
		ApexIp:            *apexIp,
		ApexUserName:      *apexUserName,
		ApexPassword:      *apexPassword,
		InfluxUser:        *influxUser,
		InfluxPassword:    *influxPassword,
		InfluxIp:          *influxIp,
		DisableDataSource: disableDataSource,
		Install:           install,
		LogFile:           *logFile,
		Quiet:             quiet,
		ExecDir:           *execDir,
		OsUserName:        *osUserName,
		TickLogFile:       *tickLogFile,
	}

	if len(*configFile) > 0 {
		loadConfigFile(*configFile, &commonFlags)
	}
	setDefaults(&commonFlags)

	//TODO: More validation

	return commonFlags
}

//Create default configuration options where necessary. This is mostly for configuring
//the install process.
func setDefaults(commonFlags *CommonFlags) {
	if commonFlags.Delay == nil {
		*commonFlags.Delay = 60
	}
	if len(commonFlags.LogLevel) == 0 {
		commonFlags.LogLevel = "warn"
	}

	if commonFlags.Install != nil && *commonFlags.Install {
		//Get the user information for installation
		var u *user.User
		var err error
		if len(commonFlags.OsUserName) == 0 {
			u, err = user.Current()
		} else {
			u, err = user.Lookup(commonFlags.OsUserName)
		}
		if err != nil {
			fmt.Println("Unable to load information for user.", err)
			log.WithError(err).Fatal("Unable to load information for user.")
		}

		commonFlags.OsUserName = u.Name
		*commonFlags.Uid, _ = strconv.Atoi(u.Uid)
		*commonFlags.Gid, _ = strconv.Atoi(u.Gid)

		//Set the location of the executable
		if len(commonFlags.ExecDir) == 0 {
			commonFlags.ExecDir = `/usr/bin`
		}

	}
}

//Assumes the config file is in the same directory as the binary.
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

	//Combine the command line parameters with the config file parameters.
	//The config file parameters win.
	//NOTE: Should it be the other way around
	if flags.Delay != nil {
		commonFlags.Delay = flags.Delay
	}
	if len(flags.LogLevel) > 0 {
		commonFlags.LogLevel = flags.LogLevel
	}
	if len(flags.ApexIp) > 0 {
		commonFlags.ApexIp = flags.ApexIp
	}
	if len(flags.ApexUserName) > 0 {
		commonFlags.ApexUserName = flags.ApexUserName
	}
	if len(flags.ApexPassword) > 0 {
		commonFlags.ApexPassword = flags.ApexPassword
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
	if flags.Quiet != nil {
		commonFlags.Quiet = flags.Quiet
	}
	if len(flags.LogFile) > 0 {
		commonFlags.LogFile = flags.LogFile
	}
	if flags.Uid != nil && *flags.Uid != -1 {
		commonFlags.Uid = flags.Uid
	}
	if flags.Gid != nil && *flags.Gid != -1 {
		commonFlags.Gid = flags.Gid
	}
	if len(flags.ExecDir) > 0 {
		commonFlags.ExecDir = flags.ExecDir
	}
	if len(flags.OsUserName) > 0 {
		commonFlags.OsUserName = flags.OsUserName
	}
	if len(flags.TickLogFile) > 0 {
		commonFlags.TickLogFile = flags.TickLogFile
	}
}
