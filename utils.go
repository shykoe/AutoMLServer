package main

import (
	"archive/tar"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)
func parseStatus(status string) string{
	switch status {
	case "running":
		return RUNNING
	case "success":
		return SUCCESS
	case "terminated":
		return USERCANCEL
	case "failed":
		return ERROR
	case "waiting":
		return READY
	}
	return ""
}
func getOptimizeMetric(OptimizeParam string, m *metric) (float32, error){
	switch OptimizeParam {
	case "loss":
		return m.loss,nil
	case "auc":
		return m.auc,nil
	case 	"predictAvg":
		return m.predictAvg,nil
	case 	"realAvg":
		return m.realAvg,nil
	case 	"copc":
		return m.copc,nil
	default:
		return 0,errors.New("Error in getOptimizeMetric!\n")
	}
}
func IsValidOptimizeParam(OptimizeParam string) bool {
	switch OptimizeParam {
	case
		"loss",
		"auc",
		"predictAvg",
		"realAvg",
		"copc":
		return true
	}
	return false
}
func initDB() error{
	var err error
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",USERNAME,PASSWORD,NETWORK,SERVER,PORT,DATABASE)
	DB, err = sql.Open("mysql",dsn)
	if err!=nil{
		log.Fatal(err)
		return err
	}
	DB.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(1000)//设置最大连接数
	DB.SetMaxIdleConns(0) //设置闲置连接数
	return nil
}
func parseMetric(data string) *metric{
	r := *regexp.MustCompile(`metric_count_ = (\d+).+ver:'([\d|_]+)'.+data_info:'([\d|_]+)'.+loglosss=([1-9]\d*.\d*|0.\d*[1-9]\d*).+auc=([1-9]\d*.\d*|0.\d*[1-9]\d*).+predict_avg=([1-9]\d*.\d*|0.\d*[1-9]\d*).+real_avg=([1-9]\d*.\d*|0.\d*[1-9]\d*).+copc=([1-9]\d*.\d*|0.\d*[1-9]\d*)`)
	match := r.FindAllStringSubmatch(data, -1)
	if match == nil{
		return nil
	}else{
		count,err := strconv.Atoi(match[0][1])
		if err!=nil{
			return nil
		}
		ver := match[0][2]
		dataInfo := match[0][3]
		loss, err:= strconv.ParseFloat(match[0][4],32)
		if err!=nil{
			return nil
		}
		auc, err:= strconv.ParseFloat(match[0][5],32)
		if err!=nil{
			return nil
		}
		predictAvg, err:= strconv.ParseFloat(match[0][6],32)
		if err!=nil{
			return nil
		}
		realAvg, err:= strconv.ParseFloat(match[0][7],32)
		if err!=nil{
			return nil
		}
		copc, err:= strconv.ParseFloat(match[0][8],32)
		if err!=nil{
			return nil
		}
		m := &metric{
			count:count,
			ver:ver,
			dataInfo:dataInfo,
			loss:float32(loss) ,
			auc:float32(auc),
			predictAvg:float32(predictAvg),
			realAvg:float32(realAvg),
			copc:float32(copc),
		}
		return m
	}
}
func initConfig( configPath string) error{
	config := make(map[string] string)
	data, err := ioutil.ReadFile(configPath)
	err = yaml.Unmarshal(data,&config)
	if err != nil{
		log.Fatal(err)
		return err
	}
	USERNAME = config["username"]
	PASSWORD = config["password"]
	SERVER = config["server"]
	DATABASE = config["database"]
	PORT = config["port"]
	S3AK = config["accessKey"]
	S3SK = config["secretKey"]
	TMPPATH = config["tmpPath"]
	BOREURL = config["boreUrl"]
	BUCKET = config["bucket"]
	S3HOST = config["s3host"]
	SKEY = config["skey"]
	BORELOGURL = config["boreLogUrl"]
	BORESTATUSURL = config["boreStatusUrl"]
	DEV = config["dev"]
	return nil
}
func getBoreLog(appName string, containerName string, logType string, offset int, length int) (int, []string, string, error){
	req, err := http.NewRequest("GET",BORELOGURL,nil)
	if err!=nil{
		log.Error(err)
		return 0, nil, "", err
	}
	q := req.URL.Query()
	q.Add("appinstance_name", appName)
	q.Add("container_name", containerName)
	q.Add("log_type", logType)
	q.Add("offset", strconv.Itoa(offset) )
	q.Add("length", strconv.Itoa(length))
	q.Add("skey", SKEY)
	req.URL.RawQuery = q.Encode()
	resp, err := http.Get(req.URL.String())
	log.Info("getBoreLog Url: ", req.URL.String())
	//resp, err := http.DefaultClient.Do(req)
	var data = make(map[string] interface{})
	response, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(response,&data); err!=nil{
		log.Error(err)
		return 0, nil, "", err
	}
	content := data["data"].(map[string]interface{})
	slog := content["sLog"].(string)
	logStr := []rune(slog)
	var result  []string
	for ; ;  {
		aLen := runeSearch(logStr,"\n")
		if aLen == -1{
			break
		}
		offset += aLen+1
		result = append(result, string(logStr[:aLen+1]) )
		logStr = logStr[aLen+1:]

	}
	return offset, result, slog, nil
}
func getBoreStatus(appName string) ( map[string]string, error){
	req, err := http.NewRequest("GET",BORESTATUSURL,nil)
	if err!=nil{
		log.Error(err)
		return nil, err
	}
	q := req.URL.Query()
	q.Add("appinstance_name", appName)
	q.Add("skey", SKEY)
	req.URL.RawQuery = q.Encode()
	resp, err := http.Get(req.URL.String())
	var data = make(map[string] interface{})
	response, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(response,&data); err!=nil{
		log.Error(err)
		return nil, err
	}
	content := data["data"].(map[string]interface{})
	result := make(map[string]string)
	for k,v := range content{
		if reflect.TypeOf(v).String() == "string"{
			result[k] = v.(string)
		}else{
			v2i := strconv.FormatFloat(v.(float64), 'E', -1,64)
			result[k] = v2i
		}
	}
	return result, nil
}
func runeSearch(text []rune, what string) int  {
	whatRunes := []rune(what)

	for i := range text {
		found := true
		for j := range whatRunes {
			if text[i+j] != whatRunes[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}
func uploadS3(fileName string, bucket string , key string, acl string) error{
	if DEV == "1"{
		return nil
	}
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(S3AK, S3SK, ""),
		Endpoint:         aws.String(S3HOST),
		Region:           aws.String("default"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	}

	sess,err := session.NewSession(s3Config)
	if err!=nil{
		log.Error("NewSession err:", err)
		return err
	}
	//svc := s3.New(sess)

	file, err := os.Open(fileName)
	if err != nil {
		log.Error("Unable to open file", err)
		return err
	}
	uploader := s3manager.NewUploader(sess)

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL: aws.String(acl),
	})
	if err != nil {
		log.Error("Unable to upload ", fileName," to ", bucket, err)
		return err
	}
	return nil
}
func createTar(path string, target string) error {
	fw, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fw.Close()
	tw := tar.NewWriter(fw)
	defer func() {
		if err := tw.Close(); err != nil {
			log.Error(err)
		}
	}()
	return filepath.Walk(path, func(fileName string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		// 这里需要处理下 hdr 中的 Name，因为默认文件的名字是不带路径的，
		// 打包之后所有文件就会堆在一起，这样就破坏了原本的目录结果
		// 例如： 将原本 hdr.Name 的 syslog 替换程 log/syslog
		// 这个其实也很简单，回调函数的 fileName 字段给我们返回来的就是完整路径的 log/syslog
		// strings.TrimPrefix 将 fileName 的最左侧的 / 去掉，
		hdr.Name = strings.TrimPrefix(strings.TrimPrefix(fileName, path), string(filepath.Separator))
		if hdr.Name == ""{
			return nil
		}
		// 写入文件信息
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// 判断下文件是否是标准文件，如果不是就不处理了，
		// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		fr, err := os.Open(fileName)
		defer fr.Close()
		if err != nil {
			return err
		}

		// copy 文件数据到 tw
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}

		// 记录下过程，这个可以不记录，这个看需要，这样可以看到打包的过程
		//log.Printf("成功打包 %s ，共写入了 %d 字节的数据\n", fileName, n)

		return nil
	})
}
