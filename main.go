package main

import (
	"container/list"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
)

//func makemsg(cmdtype string , mgs string) []byte{
//	content := fmt.Sprintf("%s%06d%s",cmdtype, len(mgs), mgs)
//	//fmt.Print(content)
//	value := []byte(content)
//	return value
//}
//func decode(data []byte) (bool, string, string, string){
//	if len(data) < 8 {
//		return false, "", "", ""
//	}
//	commandType := string(data[:2])
//	commandLength,err := strconv.Atoi(string(data[2:8]))
//	if err != nil{
//		log.Printf("error", err)
//		return false, "", "", ""
//	}
//	if len(data) < commandLength +8 {
//		return false, "", "", ""
//	}
//	command := string(data[8:commandLength+8])
//	remain := string(data[commandLength+8:])
//	return true, commandType, command, remain
//}
//func newChild(f1 *os.File, f2 *os.File){
//	cmd := exec.Command("python", "-m", "nni", "--tuner_class_name", "TPE", "--tuner_args", "{\"optimize_mode\":\"maximize\"}")
//	cmd.Env = os.Environ()
//	pip :=make([]*os.File,2)
//	pip[0] = f1
//	pip[1] = f2
//	cmd.ExtraFiles = pip
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stdout
//	err := cmd.Run()
//	if err!=nil {
//		log.Fatal(err)
//	}
//}
type AddExpJson struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	S3 string `form:"S3" json:"S3" xml:"S3" binding:"required"`
	TunerType string `form:"tunerType" json:"tunerType" xml:"tunerType" binding:"required"`
	TunerArgs string `form:"tunerArgs" json:"tunerArgs" xml:"tunerArgs" binding:"required"`
	SearchSpace string `form:"SearchSpace" json:"SearchSpace" xml:"SearchSpace" binding:"required"`
	Parallel int `form:"Parallel" json:"Parallel" xml:"Parallel" `
}
func (a *AddExpJson) toString() string{
	str := fmt.Sprintf("User: %s\nS3: %s\nTunerType: %s\nTunerArgs: %s\nSearchSpace: %s\nParallel: %d", a.User, a.S3, a.TunerType ,a.TunerArgs, a.SearchSpace, a.Parallel)
	return str
}
func main() {
	var exp = make(map[string] *experment )
	r := gin.Default()
	r.POST("/addExp", func(context *gin.Context) {
		var json AddExpJson
		//data,err :=context.GetRawData(); if err== nil{
		//	panic(err)
		//}
		//log.Info(string(data))
		if err := context.ShouldBindJSON(&json); err !=nil{
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ids := fmt.Sprintf("%s_%d",json.User ,rand.Int31())
		log.Info("get json ", json.toString())
		var parallelNum int
		if json.Parallel!=0{
			parallelNum = json.Parallel
		}else{
			parallelNum = 1
		}
		newExp :=  &experment{
			tunerType:json.TunerType,
			tunerArgs:json.TunerArgs,
			ll : list.New(),
			trials : list.New(),
			output: make(chan struct{}),
			searchSpace:json.SearchSpace,
			parallel:parallelNum,
			expId:ids,
			S3:json.S3,
		}
		go newExp.run()
		exp[ids] = newExp
		context.JSON(http.StatusOK, gin.H{"status": "success", "id":ids})

	})
	r.Run(":8989")



	//tunerType := "TPE"
	//tunerArgs := "{\"optimize_mode\":\"maximize\"}"
	//newExp :=  &experment{
	//	tunerType:tunerType,
	//	tunerArgs:tunerArgs,
	//	ll : list.New(),
	//	output: make(chan struct{}),
	//}
	//initExp := IpcData{
	//	cedType:Initialize,
	//	cmdContent:"{\"learning_rate\": {\"_type\": \"choice\", \"_value\": [0.0001, 0.001, 0.01, 0.1]}}",
	//
	//}
	//go newExp.run()
	//
	//getJob := IpcData{
	//	cedType:RequestTrialJobs,
	//	cmdContent:"1",
	//}
	//newExp.send(initExp)
	//newExp.send(getJob)
	//for ; ;  {
	//	select {
	//	case <- newExp.wait():
	//		for{
	//			data := newExp.get()
	//			if data!= nil{
	//				log.Info("asdasd!!",data.cedType, data.cmdContent)
	//			}
	//			break
	//		}
	//
	//	}
	//}


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