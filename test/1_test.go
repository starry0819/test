package test

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"testing"
)

const(
	jsonStr = "{ \"token\": \" \"}"
)


func TestJson(t *testing.T) {
	var mes Send
	//var mes1 = Send{
	//	UserNumber: "100010",
	//}

	//bytes, _ := json.Marshal(mes1)
	//fmt.Println(string(bytes[:]))
	e := json.Unmarshal([]byte(jsonStr), &mes)
	s := mes.UserNumber + "_" + mes.Token
	fmt.Println(s)
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println("11111")
	fmt.Print(mes.Token == " ")
}

func TestMD5(t *testing.T) {
	//sum := md5.Sum([]byte("123456"))
	hash := md5.New()
	hash.Write([]byte("123456"))
	fmt.Println(hex.EncodeToString(hash.Sum(nil)))
}

func TestFlagString(t *testing.T) {
	var addr = flag.String("addr", "localhost:8011", "http service address")
	fmt.Println(*addr)
}

func TestSplit(t *testing.T) {
	str := "_a_"
	split := strings.Split(str, "_")
	for i, v := range split  {
		fmt.Println(fmt.Sprintf("index:%v, value: %v", i, v))
	}
}

func TestNotInit(t *testing.T) {
	s := &Send{UserNumber:"2"}
	c := clients{
		"a": {
			&Send{UserNumber: "1",
			}: true,
		},
		"b": {
			//s: true,
		},
	}
	fmt.Println(c["b"][s])
}

type clients map[string]map[*Send]bool

type Send struct {
	UserNumber string `json:"user_number"`
	Token string `json:"token"`
}