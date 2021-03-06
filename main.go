package main

import (
	"flag"
	"fmt"
	"github.com/fofapro/fofa-go/fofa"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v2"
	"strings"
)

var(
	searchStr=flag.String("s","","fofa search string")
	json=jsoniter.ConfigCompatibleWithStandardLibrary
)
//type result struct {
//	Domain  string `json:"domain,omitempty"`
//	Host    string `json:"host,omitempty"`
//	IP      string `json:"ip,omitempty"`
//	Port    string `json:"port,omitempty"`
//	Title   string `json:"title,omitempty"`
//	Country string `json:"country,omitempty"`
//	City    string `json:"city,omitempty"`
//}
type config struct {
	Email string `yaml:"Email"`
	Api string `yaml:"Api"`
}
type response struct {
	Error interface{} `json:"error"`
	Mode interface{}	`json:"mode"`
	Page interface{} `json:"page"`
	Query interface{} `json:"query"`
	Results [][]string `json:"results"`
	Size interface{}	`json:"size"`
}

func parseConfig(fileName string) config{
	var temp config
	configData,err:=ioutil.ReadFile(fileName)
	if err!=nil{
		fmt.Println("fail to read config file")
		os.Exit(1)
	}
	err=yaml.Unmarshal(configData,&temp)
	if err!=nil{
		fmt.Println("fail to parse config file")
		os.Exit(1)
	}
	return temp
}
func parseResult(result1 []byte) [][]string {
	var temp response
	err:=json.UnmarshalFromString(string(result1),&temp)
	if err!=nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//fmt.Println(temp)
	return temp.Results
}
func fileOutPut(input [][]string){
	fileUrl,err:=os.OpenFile("url.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open url file")
		os.Exit(0)
	}
	var protocolSlice []string
	urlAll:=""
	for _,temp:=range input{
		tempStr:=""
		if len(temp[4])!=0{
			tempStr+=temp[0]+"\t"+temp[4]
			protocolSlice=append(protocolSlice,tempStr)
		}else {
			if strings.HasPrefix(temp[0],"http"){
				urlAll+=temp[0]+"\r\n"
			}else {
				urlAll+="http://"+temp[0]+"\r\n"
			}
		}
	}
	fileUrl.WriteString(strings.TrimSpace(urlAll))
	fileUrl.Close()
	fmt.Println(protocolSlice)
}
func main() {
	flag.Parse()
	if len(*searchStr)==0{
		fmt.Println("please input search string")
		os.Exit(1)
	}
	myConfig:=parseConfig("config.yaml")
	email:=myConfig.Email
	apiKey:=myConfig.Api
	clt := fofa.NewFofaClient([]byte(email), []byte(apiKey))
	if clt == nil {
		fmt.Printf("create fofa client\n")
		return
	}
	result1, err := clt.QueryAsJSON(1, []byte(*searchStr),[]byte("host,ip,port,title,protocol,country"))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	//fmt.Printf("%s\n", result1)
	//fmt.Println(parseResult(result1))
	fileOutPut(parseResult(result1))
}