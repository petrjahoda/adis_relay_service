package main

import (
	"github.com/kardianos/service"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const version = "2021.1.2.4"
const serviceName = "Adis Relay Service"
const serviceDescription = "Work with Zapsi Relays"
const downloadInSeconds = 10

var config = "zapsi_uzivatel:zapsi@tcp(zapasidatabase:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
var serviceRunning = false

var (
	activeWorkplaces  []Workplace
	runningWorkplaces []Workplace
	workplaceSync     sync.Mutex
)

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	go p.run()
	serviceRunning = true
	return nil
}

func (p *program) Stop(service.Service) error {
	serviceRunning = false
	for len(runningWorkplaces) != 0 {
		logInfo("MAIN", serviceName+" ["+version+"] stopping...")
		time.Sleep(1 * time.Second)
	}
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	if runtime.GOOS == "windows" {
		config = "zapsi_uzivatel:zapsi@tcp(localopst:3306)/zapsi2?charset=utf8mb4&parseTime=True&loc=Local"
	}
	for {
		logInfo("MAIN", serviceName+" ["+version+"] running")
		start := time.Now()
		readActiveWorkplaces("MAIN")
		logInfo("MAIN", "Active workplaces: "+strconv.Itoa(len(activeWorkplaces))+", running workplaces: "+strconv.Itoa(len(runningWorkplaces)))
		for _, activeWorkplace := range activeWorkplaces {
			activeWorkplaceIsRunning := checkWorkplaceInRunningWorkplaces(activeWorkplace)
			if !activeWorkplaceIsRunning {
				go runWorkplace(activeWorkplace)
			}
		}

		if time.Since(start) < (downloadInSeconds * time.Second) {
			sleepTime := downloadInSeconds*time.Second - time.Since(start)
			logInfo("MAIN", "Sleeping for "+sleepTime.String())
			time.Sleep(sleepTime)
		}
	}
}
