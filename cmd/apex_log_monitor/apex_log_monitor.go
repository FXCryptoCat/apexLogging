package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/william1034/apexLogging/cmd/common_flags"
	"github.com/william1034/apexLogging/internal/apex_monitor"
	"github.com/william1034/apexLogging/internal/influx"
	"github.com/william1034/apexLogging/internal/tick"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	. "time"
)

func main() {

	flags := common_flags.GetFlags()
	log.SetLevel(log.TraceLevel) //Use this until we have gotten past the install

	if *flags.Install {
		installAsService(flags)
		fmt.Println("Please run the following to complete the installation.")
		fmt.Println("sudo systemctl daemon-reload && ")
		fmt.Println("sudo systemctl enable apexMonitor.service &&")
		fmt.Println("sudo systemctl start apexMonitor.service &&")
		fmt.Println("sudo systemctl status apexMonitor.service ")
		fmt.Println()
		os.Exit(0)
	}

	configureLogs(flags.LogLevel, flags.LogFile)
	realTimeMonitor(flags)
}

func configureLogs(logLevel string, logFile string) {

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.SetLevel(level)

	if len(logFile) > 0 {
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal("Unable to create log file " + logFile)
		}
		log.SetOutput(f)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func installAsService(flags common_flags.CommonFlags) {
	if !*flags.Quiet {
		fmt.Println("You have elected to install this as a service. Make sure you are sudo. Press Y/y to continue")
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		if char != 'Y' && char != 'y' {
			fmt.Println("Exiting installation")
			os.Exit(0)
		}
	}

	saveServiceFile(flags)
	saveConfigFile(flags)
	copyBinary(flags)

}

func realTimeMonitor(flags common_flags.CommonFlags) {
	////The apex actually brings everything back. Slowly
	duration := Duration(*flags.Delay) * Second
	ticker := NewTicker(duration)
	quit := make(chan struct{})

	apexClient := apex_monitor.NewApexMonitor(flags.ApexIp, flags.ApexUserName, flags.ApexPassword)
	ic := influx.NewInfluxClient(flags.InfluxIp, flags.InfluxUser, flags.InfluxPassword)

	go func() {
		for {
			select {
			case <-ticker.C:
				nanos := Now().UnixNano()
				apexStatus, err := apexClient.GetStatus(nanos)

				if err != nil {
					fmt.Println(err)
					//Sleep for a few seconds before continuing....
					Sleep(Millisecond * 5000)

					//Try to reauth
					apexClient.ReAuth()
					continue
				}

				log.Trace("Fetching another record")
				tick := tick.GetTickLineFromApexStatus(*apexStatus)
				if *flags.DisableDataSource {
					log.Debug(tick)
				} else {
					ic.WriteRecord(tick)
				}

				//fmt.Print(tick)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit

}

func copyBinary(flags common_flags.CommonFlags) {
	src := "apexMonitor"
	dst := path.Join(flags.ExecDir, src)

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		fmt.Print("Source file does not exist", err)
		log.Fatal("Source file does not exist", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		fmt.Printf("%s is not a regular file %s\n", src, err)
		log.WithError(err).Fatalf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		fmt.Printf("Cannot open %s - %s", src, err)
		log.WithError(err).Fatalf("Cannot open %s", src)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Cannot create %s - %s", dst, err)
		log.WithError(err).Fatalf("Cannot create %s", dst)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		fmt.Printf("Cannot copy %s to  %s - %s", src, dst, err)
		log.WithError(err).Fatalf("Cannot copy %s to  %s", src, dst)
	}

	fixFilePermissions(dst, flags, 777)
}

func saveConfigFile(flags common_flags.CommonFlags) {
	configFileName := path.Join(flags.ExecDir, "apexMonitor.config.yaml")

	*flags.Install = false
	configYaml, _ := yaml.Marshal(flags)

	err := ioutil.WriteFile(configFileName, configYaml, 777)
	if err != nil {
		log.Fatalf("Unable to save file %s", configFileName)
	}

	fixFilePermissions(configFileName, flags, 777)

	fmt.Println("Saved " + configFileName)
}

func fixFilePermissions(configFileName string, flags common_flags.CommonFlags, perm os.FileMode) {
	err := os.Chown(configFileName, *flags.Uid, *flags.Gid)
	if err != nil {
		log.Errorf("Unable to group to gid: %d and user to uid: ", *flags.Uid, *flags.Gid)
		log.Fatal(err)
	}

	err = os.Chmod(configFileName, perm)
	if err != nil {
		log.Errorf("Unable change permissions to " + perm.String())
		log.Fatal(err)
	}
}

func saveServiceFile(flags common_flags.CommonFlags) {
	cfg := ini.Empty()
	unitSection, _ := cfg.NewSection("Unit")
	_, _ = unitSection.NewKey("Description", "APEX Influxdb Logger")
	_, _ = unitSection.NewKey("After", "network.target remote-fs.target nss-lookup.target")

	serviceSection, _ := cfg.NewSection("Service")
	_, _ = serviceSection.NewKey("Type", "simple")
	_, _ = serviceSection.NewKey("User", flags.OsUserName)

	_, _ = serviceSection.NewKey("WorkingDirectory", flags.ExecDir)

	execFile := path.Join(flags.ExecDir, "apexMonitor")
	parameters := `-config="apexMonitor.config.yaml"`
	_, _ = serviceSection.NewKey("ExecStart", execFile+" "+parameters)

	_, _ = serviceSection.NewKey("PrivateTmp", "true")
	_, _ = serviceSection.NewKey("LimitNOFILE", "infinity")
	_, _ = serviceSection.NewKey("KillMode", "mixed")
	_, _ = serviceSection.NewKey("Restart", "on-failure")
	_, _ = serviceSection.NewKey("RestartSec", "5s")

	installSection, _ := cfg.NewSection("Install")
	_, _ = installSection.NewKey("WantedBy", "multi-user.target")

	configFileName := path.Join("/etc/systemd/system", "apexMonitor.service")
	err := cfg.SaveTo(configFileName)
	if err != nil {
		fmt.Println("Unable to save config file: " + configFileName)
		log.Fatal("Unable to save config file: " + configFileName)

	}

	fmt.Println("Saved " + configFileName)
}
