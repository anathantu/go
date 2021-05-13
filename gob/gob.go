package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func newGobFile(){
	info:=map[string]string{
		"name":"this is github",
		"website":"wwww.github.com",
	}
	name:="demo.gob"
	file,_ := os.OpenFile(name,os.O_RDWR|os.O_CREATE,0777)
	defer file.Close()
	enc:=gob.NewEncoder(file)
	if err:=enc.Encode(info);err!=nil{
		fmt.Println(err)
	}
}

func readGobFile(){
	var M map[string]string
	file,_:=os.Open("demo.gob")
	d:=gob.NewDecoder(file)
	d.Decode(&M)
	fmt.Println(M)
}

func main(){
	newGobFile()
	readGobFile()
}
