package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	filepath    = "./avs3/test1.avs3"
	newfilepath = "./avs3/newAvs3_3.txt"
)

type Pos struct {
	pos    int
	header int //0 代表B0 、1代表B3、2代表B2、4代表B6
}

func main() {
	//str, _ := os.Getwd()
	//fmt.Println(str)
	read1()
}

func read1() {
	file, err := os.Open(newfilepath)
	if err != nil {
		fmt.Println("open file failed!, err:", err)
		return
	}

	defer func() {
		file.Close()
	}()
	// 循环读取文件
	var content []byte
	var tmp = make([]byte, 128)

	for {
		n, err := file.Read(tmp)
		if err == io.EOF {
			fmt.Println("文件读完了")
			break
		}
		if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		fmt.Println(string(tmp))
		content = append(content, tmp[:n]...)
	}
	fmt.Println(string(content))
}

func read2() {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("read file failed, err:", err)
		return
	}
	var left int
	var right int

	//按照video sequence来写
	for i := 3; i < len(content)-4; {
		if isSequenceHeader(content, i) {
			//处理序列头和帧内编码
			i++
			left = i
			for {
				if isInFrameHeader(content, i) {
					//找到帧间编码头
					right = i - 4
					i++
					break
				}
				i++
			}
			//数据赋值
			//changeValue(content,left,right)
			//break
			left = i
			for {
				if isUserInfoHeader(content, i) {

					right = i - 4
					i++
					break
				}
				i++
			}
			//数据赋值
			changeValue(content, left, right)
			left = i
		}
		if isBtwFrameHeader(content, i) {
			i++
			left = i
			for {
				if isUserInfoHeader(content, i) {

					right = i - 4
					i++
					break
				}
				i++
			}
			changeValue(content, left, right)
		}

		i++
	}

	err = ioutil.WriteFile(newfilepath, content, 0666)
	if err != nil {
		fmt.Println("write file failed, err:", err)
		return
	}
}

func read3() {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("read file failed, err:", err)
		return
	}

	err = ioutil.WriteFile(newfilepath, content, 0666)
	if err != nil {
		fmt.Println("write file failed, err:", err)
		return
	}
}

func changeValue(content []byte, left int, right int) {
	for i := left; i < right; i++ {
		content[i] = 0
	}
}

func isUserInfoHeader(content []byte, pos int) bool {
	if content[pos] == 178 && content[pos-1] == 1 && content[pos-2] == 0 && content[pos-3] == 0 {
		return true
	}
	return false
}

//帧间
func isBtwFrameHeader(content []byte, pos int) bool {
	if content[pos] == 182 && content[pos-1] == 1 && content[pos-2] == 0 && content[pos-3] == 0 {
		return true
	}
	return false
}

func isSequenceHeader(content []byte, pos int) bool {
	if content[pos] == 176 && content[pos-1] == 1 && content[pos-2] == 0 && content[pos-3] == 0 {
		return true
	}
	return false
}

//帧内编码头
func isInFrameHeader(content []byte, pos int) bool {
	if content[pos] == 179 && content[pos-1] == 1 && content[pos-2] == 0 && content[pos-3] == 0 {
		return true
	}
	return false
}
