package main

import "time"

type Workplace struct {
	OID                 int    `gorm:"column:OID"`
	Name                string `gorm:"column:Name"`
	WorkplaceDivisionId int    `gorm:"column:WorkplaceDivisionID"`
	DeviceID            int    `gorm:"column:DeviceID"`
	Code                string `gorm:"column:Code"`
}

func (Workplace) TableName() string {
	return "workplace"
}

type WorkplacePort struct {
	OID          int    `gorm:"column:OID"`
	DevicePortID int    `gorm:"column:DevicePortID"`
	Type         string `gorm:"column:Type"`
}

func (WorkplacePort) TableName() string {
	return "workplace_port"
}

type DevicePort struct {
	OID      int `gorm:"column:OID"`
	DeviceID int `gorm:"column:DeviceID"`
}

func (DevicePort) TableName() string {
	return "device_port"
}

type Device struct {
	OID     int    `gorm:"column:OID"`
	IPAddress string `gorm:"column:IPAddress"`
}

func (Device) TableName() string {
	return "device"
}

type TerminalInputOrder struct {
	OID      int       `gorm:"column:OID"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE"`
	DeviceID int       `gorm:"column:DeviceID"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}
