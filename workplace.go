package main

import (
	"bytes"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func runWorkplace(workplace Workplace) {
	logInfo(workplace.Name, "Workplace active, started running")
	workplaceSync.Lock()
	runningWorkplaces = append(runningWorkplaces, workplace)
	workplaceSync.Unlock()
	workplaceIsActive := true
	for workplaceIsActive && serviceRunning {
		logInfo(workplace.Name, "Workplace main loop started")
		timer := time.Now()
		terminalDeviceId := updateTerminalDeviceIdFor(workplace)
		logInfo(workplace.Name, "Actual terminal device id: "+strconv.Itoa(terminalDeviceId))
		zapsiDeviceId := updateZapsiDeviceIdFor(workplace)
		logInfo(workplace.Name, "Actual zapsi device id: "+strconv.Itoa(zapsiDeviceId))
		terminalHasOpenOrder := checkOpenTerminalInputOrderFor(terminalDeviceId, workplace)
		logInfo(workplace.Name, "Terminal has open terminal input order: "+strconv.FormatBool(terminalHasOpenOrder))
		if terminalHasOpenOrder {
			zapsiRelayEnabled, deviceIpAddress := checkForZapsiRelay(zapsiDeviceId, workplace)
			if !zapsiRelayEnabled {
				logInfo(workplace.Name, "Zapsi relay not enabled, enabling")
				enableZapsiRelay(deviceIpAddress, workplace)
			}
		}
		logInfo(workplace.Name, "Main loop ended in "+time.Since(timer).String())
		if time.Since(timer) < (downloadInSeconds * time.Second) {
			sleepTime := downloadInSeconds*time.Second - time.Since(timer)
			logInfo(workplace.Name, "Sleeping for "+sleepTime.String())
			time.Sleep(sleepTime)
		}
		workplaceIsActive = checkActive(workplace)
	}
	removeWorkplaceFromRunningWorkplaces(workplace)
	logInfo(workplace.Name, "Workplace not active, stopped running")
}

func enableZapsiRelay(deviceIpAddress string, workplace Workplace) {
	logInfo(workplace.Name, "Enabling device relay for " + deviceIpAddress)
	conn, err := net.Dial("tcp", deviceIpAddress+":80")
	if err != nil {
	}
	defer conn.Close()
	fmt.Fprintf(conn, "SET /Rele1")
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, conn)
}



func checkForZapsiRelay(zapsiDeviceId int, workplace Workplace) (bool, string) {
	logInfo(workplace.Name, "Checking for open relay")
	db, err := gorm.Open(mysql.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(workplace.Name, "Problem opening database: "+err.Error())
		return false, ""
	}

	var device Device
	db.Where("OID = ?", zapsiDeviceId).Find(&device)
	deviceHasOpenRelay := checkDeviceWithIp(device.IPAddress)
	logInfo(workplace.Name, "Relay open: "+strconv.FormatBool(deviceHasOpenRelay))
	return deviceHasOpenRelay, device.IPAddress
}

func checkDeviceWithIp(deviceIpAddress string) bool {
	conn, err := net.Dial("tcp", deviceIpAddress+":80")
	if err != nil {
	}
	defer conn.Close()
	fmt.Fprintf(conn, "GET /IO")
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, conn)
	result := buf.String()
	for i, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {
		if i == 5 {
			if string([]rune(line)[0]) == "1" {
				return true
			}
		}
	}
	return false
}

func checkOpenTerminalInputOrderFor(terminalDeviceId int, workplace Workplace) bool {
	logInfo(workplace.Name, "Checking open terminal input order, please wait for 3 seconds...")
	db, err := gorm.Open(mysql.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(workplace.Name, "Problem opening database: "+err.Error())
		return false
	}

	var terminalInputOrderFirst TerminalInputOrder
	db.Where("DTE is null").Where("DeviceID = ?", terminalDeviceId).Find(&terminalInputOrderFirst)
	time.Sleep(3 * time.Second)
	var terminalInputOrderSecond TerminalInputOrder
	db.Where("DTE is null").Where("DeviceID = ?", terminalDeviceId).Find(&terminalInputOrderSecond)
	if terminalInputOrderFirst.OID == 0 && terminalInputOrderSecond.OID == 0 {
		return false
	}
	return true
}

func updateZapsiDeviceIdFor(workplace Workplace) int {
	logInfo(workplace.Name, "Updating zapsi device id number")
	db, err := gorm.Open(mysql.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(workplace.Name, "Problem opening database: "+err.Error())
		return 0
	}
	var workplacePort WorkplacePort
	db.Where("WorkplaceID = ?", workplace.OID).Where("Type = 'current'").Find(&workplacePort)
	logInfo(workplace.Name, "Workplace port id: "+strconv.Itoa(workplacePort.OID))
	var deviceport DevicePort
	db.Where("OID = ?", workplacePort.DevicePortID).Find(&deviceport)
	logInfo(workplace.Name, "Zapsi port id: "+strconv.Itoa(deviceport.OID))
	return deviceport.DeviceID
}

func updateTerminalDeviceIdFor(workplace Workplace) int {
	logInfo(workplace.Name, "Updating terminal device id number")
	db, err := gorm.Open(mysql.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(workplace.Name, "Problem opening database: "+err.Error())
		return 0
	}
	var workplaceRefresh Workplace
	db.Where("OID = ?", workplace.OID).Find(&workplaceRefresh)
	return workplaceRefresh.DeviceID
}

func readActiveWorkplaces(reference string) {
	logInfo("MAIN", "Reading active workplaces")
	timer := time.Now()
	db, err := gorm.Open(mysql.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError(reference, "Problem opening database: "+err.Error())
		activeWorkplaces = nil
		return
	}
	db.Find(&activeWorkplaces)
	logInfo("MAIN", "Active workplaces read in "+time.Since(timer).String())
}

func checkActive(workplace Workplace) bool {
	for _, activeWorkplace := range activeWorkplaces {
		if activeWorkplace.Name == workplace.Name {
			logInfo(workplace.Name, "Workplace still active")
			return true
		}
	}
	logInfo(workplace.Name, "Workplace not active")
	return false
}

func removeWorkplaceFromRunningWorkplaces(workplace Workplace) {
	workplaceSync.Lock()
	for idx, runningWorkplace := range runningWorkplaces {
		if workplace.Name == runningWorkplace.Name {
			runningWorkplaces = append(runningWorkplaces[0:idx], runningWorkplaces[idx+1:]...)
		}
	}
	workplaceSync.Unlock()
}

func checkWorkplaceInRunningWorkplaces(workplace Workplace) bool {
	for _, runningWorkplace := range runningWorkplaces {
		if runningWorkplace.Name == workplace.Name {
			return true
		}
	}
	return false
}
