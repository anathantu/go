package slice

import "fmt"

func createSlice(){
	//声明切片
	var s1 []int
	if s1==nil{
		fmt.Println("是空")
	}else{
		fmt.Println("不是空")
	}
	//用:=
	s2:=[]int{}
	//用make函数
	var s3=make([]int,0)
	fmt.Println(s1,s2,s3)
	//初始化赋值
	var s4 []int=make([]int,0,0)
	fmt.Println(s4)
	s5:=[]int{1,2,3,4}
	fmt.Println(s5)
	//从数组切片
	arr:=[5]int{1,2,3,4,5}
	var s6 []int
	s6=arr[1:4]
	fmt.Println(s6)
}

func rwSlice(){
	data:=[...]int{0,1,2,3,4,5}

	s:=data[2:4]
	s[0]+=100
	s[1]+=200

	fmt.Println(s)
	fmt.Println(data)
	/**
	[102 203]
	[0 1 102 203 4 5]
	 */
}