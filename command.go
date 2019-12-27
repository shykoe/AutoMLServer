package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	Initialize            = "IN"
	RequestTrialJobs      = "GE"
	ReportMetricData      = "ME"
	UpdateSearchSpace     = "SS"
	ImportData            = "FD"
	AddCustomizedTrialJob = "AD"
	TrialEnd              = "EN"
	Terminate             = "TE"
	Ping                  = "PI"

	Initialized           = "ID"
	NewTrialJob           = "TR"
	SendTrialJobParameter = "SP"
	NoMoreTrialJobs       = "NO"
	KillTrialJob          = "KI"
)

type ReportData struct {
	ParameterId string  `json:"parameter_id"`
	TrialJobId  string  `json:"trial_job_id"`
	Type        string  `json:"type"`
	Sequence    int     `json:"sequence"`
	Value       float32 `json:"value"`
}

func (r *ReportData) toJsonStr() ([]byte, error) {
	return json.Marshal(r)
}

type IpcData struct {
	CedType    string
	CmdLength  int
	CmdContent string
	IpcRemain  string
}

func (d *IpcData) decode() []byte {
	d.CmdLength = len(d.CmdContent)
	content := fmt.Sprintf("%s%06d%s", d.CedType, d.CmdLength, d.CmdContent)
	log.Info("cmd decode ", content)
	return []byte(content)
}
func convert(data []byte) (error, *IpcData) {
	if len(data) < 8 {
		return errors.New("Invalid data! "), nil
	}
	commandType := string(data[:2])
	commandLength, err := strconv.Atoi(string(data[2:8]))
	if err != nil {
		return errors.New("Invalid data! "), nil
	}
	if len(data) < commandLength+8 {
		return errors.New("Invalid data! "), nil
	}
	command := string(data[8 : commandLength+8])
	remain := string(data[commandLength+8:])
	ipcData := &IpcData{
		commandType,
		commandLength,
		command,
		remain,
	}
	return nil, ipcData
}
