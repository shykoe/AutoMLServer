package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	Initialize = "IN"
	RequestTrialJobs = "GE"
	ReportMetricData = "ME"
	UpdateSearchSpace = "SS"
	ImportData = "FD"
	AddCustomizedTrialJob = "AD"
	TrialEnd = "EN"
	Terminate = "TE"
	Ping = "PI"

	Initialized = "ID"
	NewTrialJob = "TR"
	SendTrialJobParameter = "SP"
	NoMoreTrialJobs = "NO"
	KillTrialJob = "KI"
)
type IpcData struct{
	cedType    string
	cmdLength  int
	cmdContent string
	ipcRemain  string
}
func (d *IpcData) decode() []byte{
	d.cmdLength = len(d.cmdContent)
	content := fmt.Sprintf("%s%06d%s",d.cedType, d.cmdLength, d.cmdContent)
	log.Info("cmd decode ", content)
	return []byte(content)
}
func convert(data []byte) (error, *IpcData){
	if len(data) < 8 {
		return errors.New("Invalid data! "), nil
	}
	commandType := string(data[:2])
	commandLength,err := strconv.Atoi(string(data[2:8]))
	if err != nil{
		return errors.New("Invalid data! "), nil
	}
	if len(data) < commandLength +8 {
		return errors.New("Invalid data! "), nil
	}
	command := string(data[8:commandLength+8])
	remain := string(data[commandLength+8:])
	ipcData := &IpcData{
		commandType,
		commandLength,
		command,
		remain,
	}
	return nil, ipcData
}