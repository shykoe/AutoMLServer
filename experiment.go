package main

import (
	"os"
	"os/exec"
	log "github.com/sirupsen/logrus"

)
type experment struct {
	write4out *os.File
	read4nni *os.File
	write4nni *os.File
	read4out *os.File
	tunerType string
	tunerArgs string
	input chan interface{}
	output chan interface{}
}
func (e *experment) listen(){
	for ; ;  {
		out := make([]byte, 10000)
		_, err := e.read4out.Read(out)
		if err != nil{
			log.Warn("read error", err)
			continue
		}
		err, cmdData := convert(out)
		if err != nil{
			log.Warn("convert error", err)
			continue
		}
		e.output <- cmdData
	}
}
func (e *experment) run(){
	cmd := exec.Command("python", "-m", "nni", "--tuner_class_name", e.tunerType, "--tuner_args", e.tunerArgs)
	cmd.Env = os.Environ()
	r1, w1, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	r2, w2, err := os.Pipe()
	if err != nil{
		panic(err)
	}
	pip :=make([]*os.File,2)
	pip[0] = r1
	pip[1] = w2
	cmd.ExtraFiles = pip
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	e.read4nni = r1
	e.write4out = w1
	e.write4nni = w2
	e.read4out = r2
	go cmd.Run()

}
func (e *experment) send(data IpcData){
	log.Info("send", data.cedType, data.cmdContent)
	byteContent := data.decode()
	_, err := e.write4out.Write(byteContent)
	if err != nil {
		panic(err)
	}
}