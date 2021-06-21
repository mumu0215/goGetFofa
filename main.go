package main

import (
	"flag"
	"fmt"
	"github.com/fofapro/fofa-go/fofa"
	jsoniter "github.com/json-iterator/go"
	"os"
)

var(
	searchStr=flag.String("s","","fofa search string")
	outFile=flag.String("f","out.txt","output file")
	json=jsoniter.ConfigCompatibleWithStandardLibrary
)
type result struct {
	Domain  string `json:"domain,omitempty"`
	Host    string `json:"host,omitempty"`
	IP      string `json:"ip,omitempty"`
	Port    string `json:"port,omitempty"`
	Title   string `json:"title,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
}
type response struct {
	Error interface{} `json:"error"`
	Mode interface{}	`json:"mode"`
	Page interface{} `json:"page"`
	Query interface{} `json:"query"`
	Results [][]string `json:"results"`
	Size interface{}	`json:"size"`
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

func main() {
	flag.Parse()
	if len(*searchStr)==0{
		fmt.Println("please input search string")
		os.Exit(1)
	}
	email:=""
	apiKey:=""
	clt := fofa.NewFofaClient([]byte(email), []byte(apiKey))
	if clt == nil {
		fmt.Printf("create fofa client\n")
		return
	}
	result1, err := clt.QueryAsJSON(1, []byte(*searchStr))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	//fmt.Printf("%s\n", result1)
	fmt.Println(parseResult(result1))
}