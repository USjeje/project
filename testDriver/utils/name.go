package utils

import (
	"io/ioutil"
	"strings"
)

var (
	//姓名由surname:姓, seniority:辈分, 名:lastName构成
	surnameFilePath   = `E:\Learning\goproject\src\MyProje02\test01\testDriver\name\surnames.txt`
	seniorityFilePath = `E:\Learning\goproject\src\MyProje02\test01\testDriver\name\seniority.txt`
	lastNamePath      = `E:\Learning\goproject\src\MyProje02\test01\testDriver\name\lastName.txt`
)

//initSurnames 从姓氏文件中读取姓到切片
func initSurname() []string {
	bytes, err := ioutil.ReadFile(surnameFilePath)
	HandleErr(err, `ioutil.ReadFile(surnameFilePath)`)
	//一定要注意，空格可能是"\n"或者为"\r\n",可以通过fmt.Printf("%q", target)来查看到底是哪一种
	surnameStr := strings.Replace(string(bytes), "\r\n", "", -1)
	return strings.Split(surnameStr, ",")
}

//initSeniority 从辈分文件中读取辈分到切片
func initSeniority() []string {
	bytes, err := ioutil.ReadFile(seniorityFilePath)
	HandleErr(err, `ioutil.ReadFile(seniorityFilePath)`)
	seniorityStr := strings.Replace(string(bytes), "\r\n", "", -1)
	return strings.Split(seniorityStr, "、")
}

//initLastName 从名文件中读取名到切片
func initLastName() []string {
	bytes, err := ioutil.ReadFile(lastNamePath)
	HandleErr(err, "ioutil.ReadFile(lastNamePath)")
	lastNameStr := strings.Replace(string(bytes), "\r\n", "", -1)
	//fmt.Printf("%q\n", strings.Split(lastNameStr, ""))
	return strings.Split(lastNameStr, "、")
}

//GetGetRandomName 返回一个随机的姓名
func GetRandomName() string {
	surnames := initSurname()
	senioritys := initSeniority()
	lastNames := initLastName()

	surname := surnames[GetRandomInt(0, len(surnames))]
	seniority := senioritys[GetRandomInt(0, len(senioritys))]
	lastName := lastNames[GetRandomInt(0, len(lastNames))]
	//name := surname + seniority + lastName
	//fmt.Println("surname", surname, "seniority", seniority, "lastName", lastName)
	return surname + seniority + lastName
}

//InitName 初始化考生姓名
func InitName(chName chan string, max int) {

	for i := 0; i < max; i++ {
		chName<- GetRandomName()
	}
	//关闭管道
	close(chName)

}
