package main

import (
	"container/list"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"path"
	"time"
)

var DB *sql.DB
var (
	USERNAME      string
	PASSWORD      string
	NETWORK       = "tcp"
	SERVER        string
	PORT          string
	DATABASE      string
	Count         int64
	S3HOST        string
	S3AK          string
	S3SK          string
	TMPPATH       string
	BOREURL       string
	BUCKET        string
	SKEY          string
	BORELOGURL    string
	BORESTATUSURL string
	DEV           string
)

type AddExpJson struct {
	User string `form:"user" json:"user" xml:"user"  binding:"required"`
	//WorkDir string `form:"workDir" json:"workDir" xml:"workDir" binding:"required"`
	TunerType   string `form:"tunerType" json:"tunerType" xml:"tunerType" binding:"required"`
	TunerArgs   string `form:"tunerArgs" json:"tunerArgs" xml:"tunerArgs" binding:"required"`
	SearchSpace string `form:"SearchSpace" json:"SearchSpace" xml:"SearchSpace" binding:"required"`
	Parallel    int    `form:"Parallel" json:"Parallel" xml:"Parallel" binding:"required"`
	MaxTrialNum int    `form:"MaxTrialNum" json:"MaxTrialNum" xml:"MaxTrialNum" binding:"required"`
	//BoreFile string `form:"BoreFile" json:"BoreFile" xml:"BoreFile" binding:"required"`
	OptimizeParam string `form:"OptimizeParam" json:"OptimizeParam" xml:"OptimizeParam" binding:"required"`
	S3Tar         string `form:"S3Tar" json:"S3Tar" xml:"S3Tar" binding:"required"`
	S3BoreFile    string `form:"S3BoreFile" json:"S3BoreFile" xml:"S3BoreFile" binding:"required"`
}

func (a *AddExpJson) toString() string {
	str := fmt.Sprintf("User: %s\nS3Tar: %s\nTunerType: %s\nTunerArgs: %s\nSearchSpace: %s\nParallel: %d", a.User, a.S3Tar, a.TunerType, a.TunerArgs, a.SearchSpace, a.Parallel)
	return str
}

func main() {
	port := flag.Int("port", 8989, "Server port")
	configFile := flag.String("config", "./config.yml", "config file")
	flag.Parse()
	log.SetReportCaller(true)
	err := initConfig(*configFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = initDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	//var exp = make(map[string]*experiment)
	r := gin.Default()
	r.POST("/addExp", func(context *gin.Context) {
		var json AddExpJson
		if err := context.ShouldBindJSON(&json); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rand.Seed(time.Now().UnixNano())
		ids := fmt.Sprintf("%s_%d", json.User, rand.Int31())
		log.Info("get json ", json.toString())
		var parallelNum int
		if json.Parallel != 0 {
			parallelNum = json.Parallel
		} else {
			parallelNum = 1
		}
		ok := IsValidOptimizeParam(json.OptimizeParam)
		if !ok {
			context.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Error in OptimizeParam!"})
			return
		}
		var boreFile string
		if DEV == "1"{
			tarFile := "/Users/shykoe/go/src/auto/mock/20191223/s3/test.tar"
			err = extractTar(tarFile, fmt.Sprintf("%s/%s", TMPPATH, ids))
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
				return
			}
			boreFile = "/Users/shykoe/go/src/auto/mock/20191223/bore.json"
		}else{
			fileName := path.Base(json.S3Tar)
			tarFile := fmt.Sprintf("%s/%s", TMPPATH, fileName)
			err = downloadS3(tarFile, BUCKET, json.S3Tar)
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
				return
			}
			//ex tar
			err = extractTar(tarFile, fmt.Sprintf("%s/%s", TMPPATH, ids))
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
				return
			}

			// download bore.json
			boreFileName := path.Base(json.S3BoreFile)
			boreFile = fmt.Sprintf("%s/%s", TMPPATH, boreFileName)
			err = downloadS3(boreFile, BUCKET, json.S3BoreFile)
			if err != nil {
				context.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
				return
			}
		}
		//download Tar


		newExp := &experiment{
			tunerType:     json.TunerType,
			tunerArgs:     json.TunerArgs,
			ll:            list.New(),
			trials:        list.New(),
			output:        make(chan struct{}),
			trialChan:     make(chan string, 100),
			searchSpace:   json.SearchSpace,
			parallel:      parallelNum,
			expName:       ids,
			workDir:       fmt.Sprintf("%s/%s", TMPPATH, ids),
			runner:        json.User,
			maxTrialNum:   json.MaxTrialNum,
			boreFile:      boreFile,
			optimizeParam: json.OptimizeParam,
		}
		go newExp.run()
		//exp[ids] = newExp
		context.JSON(http.StatusOK, gin.H{"status": "success", "id": ids})

	})

	r.GET("/expStatus/:id", func(context *gin.Context) {
		ids := context.Param("id")
		rows, err := DB.Query("select  status from `t_experiment_info` where experiment_name = ?", ids)
		if err!=nil{
			context.JSON(http.StatusOK, gin.H{"status": "error", "msg": "DB error!"})
		}
		var name string
		for rows.Next() {
			if err := rows.Scan(&name); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Fatal(err)
			}
		}
		context.JSON(http.StatusOK, gin.H{"status": "success", "msg": name})
	})
	r.Run(fmt.Sprintf(":%d", *port))

}
