package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func makemsg(cmdtype string , mgs string) []byte{
	content := fmt.Sprintf("%s%06d%s",cmdtype, len(mgs), mgs)
	//fmt.Print(content)
	value := []byte(content)
	return value
}
func decode(data []byte) (bool, string, string, string){
	if len(data) < 8 {
		return false, "", "", ""
	}
	commandType := string(data[:2])
	commandLength,err := strconv.Atoi(string(data[2:8]))
	if err != nil{
		log.Printf("error", err)
		return false, "", "", ""
	}
	if len(data) < commandLength +8 {
		return false, "", "", ""
	}
	command := string(data[8:commandLength+8])
	remain := string(data[commandLength+8:])
	return true, commandType, command, remain
}
func newChild(f1 *os.File, f2 *os.File){
	cmd := exec.Command("python", "-m", "nni", "--tuner_class_name", "TPE", "--tuner_args", "{\"optimize_mode\":\"maximize\"}")
	cmd.Env = os.Environ()
	pip :=make([]*os.File,2)
	pip[0] = f1
	pip[1] = f2
	cmd.ExtraFiles = pip
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err!=nil {
		log.Fatal(err)
	}
}

func main() {

	tunerType := "TPE"
	tunerArgs := "{\"optimize_mode\":\"maximize\"}"
	newExp :=  &experment{
		tunerType:tunerType,
		tunerArgs:tunerArgs,
	}
	newExp.run()

	//r1, w1, err := os.Pipe()
	//if err != nil {
	//	panic(err)
	//}
	//r2, w2, err := os.Pipe()
	//if err != nil {
	//	panic(err)
	//}
	//if err != nil {
	//	panic(err)
	//}
	//go newChild(r1,w2)
	//data := makemsg("IN", "{\"learning_rate\": {\"_type\": \"choice\", \"_value\": [0.0001, 0.001, 0.01, 0.1]}}")
	//_,err = w1.Write(data)
	//if err !=nil{
	//	panic(err)
	//}
	//data2 := makemsg("GE","1")
	//_,err = w1.Write(data2)
	//if err != nil{
	//	panic(err)
	//}
	//for ; ;  {
	//	out := make([]byte, 10000)
	//	_, err = r2.Read(out)
	//	isok, cmdtype, cmd, remain :=decode(out)
	//	if isok == false{
	//		continue
	//	}
	//	log.Print("asdasd!!",cmdtype, cmd, remain)
	//}

}