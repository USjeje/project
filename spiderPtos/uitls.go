package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/axgle/mahonia"
)

//字符格式转换

//ConvertToByte src为要转换的字符串，srcCode为待转换的编码格式，targetCode为要转换的编码格式
func ConvertToByte(src, srcCode, targetCode string) []byte {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tarCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tarCoder.Translate([]byte(srcResult), true)
	return cdata
}

//WriteResult 将[]string结果写入文件result/txt
func WriteResult(str []string) {
	filename := `E:\Learning\goproject\src\go_project\spider_img\test03\result.txt`

	//打开文件
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, v := range str {
		writer.WriteString(v + "\r\n")
	}

	//一定不要忘了刷新缓冲
	writer.Flush()
}


//HandleErr 处理错误
func HandleErr(err error, str string) {
	if err != nil {
		fmt.Println(err, str)
		//0代表成功，非0代表失败
		os.Exit(1)
	}
}

//GetHTML 返回URL的页面代码
func GetHTML(url string) string {
	resp, err := http.Get(url)
	HandleErr(err, `http.Get(url)`)
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	HandleErr(err, `ioutil.ReadAll(resp.Body)`)
	return string(bytes)
}

//GetImgPageURL 获取图片的url并返回切片
func GetImgPageURL(url string) []string {
	html := GetHTML(url)
	imgURLs := make([]string, 0)
	re := regexp.MustCompile(reImg)
	result := re.FindAllStringSubmatch(html, -1)
	for _, v := range result {
		imgURLs = append(imgURLs, v[1])
	}
	return imgURLs
}
