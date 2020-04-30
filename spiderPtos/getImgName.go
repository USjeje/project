package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var (
	//reImg 图片正则
	//http正则包括http和https
	//reImg = `<img[\s\S]+?src="(http[\s\S]+?)"`
	//上面的[\s\S]包括换行符
	reImg = `<img.+?src="(http.+?)"`
	//基于alt属性为图片命名
	reImgName = `<img.+?alt="(.+?)"`
	//在同一个img下抓url和alt, 由于不知url和alt相对位置,所以先抓url后在抓alt
	reImgWithAlt = `<img.+?src="(https?.+?)".*?>` //获得图片不完全
	//有的图片不是src而是其他的
	reImgInfo = `<img.+?(https?:.+?)".*?>`
	//从img标签中提取alt内容
	reAlt = `<img.+?alt="(.+?)"`

	//从img的URL中提取原始的图片名
	/* GetNameFormURL
	eg : http://cms-bucket.ws.126.net/2020/0409/efccae1cj00q8h9lw002pc0008w005lc.jpg?
	*/
	reImgNameFormURL = `/\w+.((jpg)|(jpeg)|(png)|(gif)|(bmp)|(webp)|(swf)|(ico))`
)

//GetRandomInt 获得不同时间的随机数,生成[start, end)之间的随机数
func GetRandomInt(start, end int) int {
	//不能让时间纳秒重合
	randomMT.Lock()
	<-time.After(1 * time.Nanosecond) //效率不高，可以换成NewTimer(d.C)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := start + r.Intn(end-start)
	randomMT.Unlock()
	return ret
}

//GetRandomName 获得随机文件名
func GetRandomName() string {
	timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	randomNum := strconv.Itoa(GetRandomInt(1000, 10000)) //1000-9999 四位随机数
	return timestamp + "_" + randomNum
}

//使用alt属性来给文件命名
//<img ... alt="第二排比对手宽敞 奔驰SUV配七座满足家用" src="..."...>
//有一个问题是得到的alt比实际由http|https下载下来的图片个数多，原因是reImgName的正则表达式中没有限制http
func spiderImgName(url string) {
	html := GetHTML(url)
	imgName := make([]string, 0)

	//将gbk转为utf-8
	bytes := ConvertToByte(html, "gbk", "utf-8")
	re := regexp.MustCompile(reImgName)
	result := re.FindAllStringSubmatch(string(bytes), -1)

	fmt.Println("从alt属性获得图片名字个数为： ", len(result))
	for _, v := range result {
		imgName = append(imgName, v[1])
		fmt.Println(v[1])
	}
}

//先抓url后抓alt
func spiderImgNameWithAlt(url string) {
	html := GetHTML(url)
	imgURLWithAlt := make([]string, 0)

	//将gbk转为utf-8
	bytes := ConvertToByte(html, "gbk", "utf-8")
	re := regexp.MustCompile(reImgWithAlt)
	result := re.FindAllStringSubmatch(string(bytes), -1)
	//fmt.Println(result)

	fmt.Println("从alt属性获得图片名字个数为： ", len(result))
	for _, v := range result {
		imgURLWithAlt = append(imgURLWithAlt, v[0])
		//fmt.Println(i, v[0], "\n", GetImgNameFormTag(v[0]))
	}

	//将结果写入result.txt文件
	WriteResult(imgURLWithAlt)
}

//GetImgTagInfo 得到img标签信息, 饭返回[]map[url]alt
func GetImgTagInfo(url string) []map[string]string {
	//初始化切片
	imgInfos := make([]map[string]string, 0)
	//获取页面
	html := GetHTML(url)
	//转换格式
	//bytes := ConvertToByte(html, "gbk", "utf-8")
	//正则匹配
	re := regexp.MustCompile(reImgInfo)
	result := re.FindAllStringSubmatch(string(bytes), -1)
	//赋值
	for _, v := range result {
		//初始化map
		imgInfo := make(map[string]string)
		url := v[1]
		filename := GetImgNameFormTag(v[0])
		imgInfo["url"] = url
		//文件名加三位随机数是防止alt属性一样
		imgInfo["filename"] = filename + strconv.Itoa(GetRandomInt(100, 1000))
		imgInfos = append(imgInfos, imgInfo)
		// fmt.Println(i, imgInfo["url"], imgInfo["filename"])
		// fmt.Println("------------------------------------------------------------")
	}
	return imgInfos
}

//GetNameFormURL 从URL里面提取原始的图片名
func GetNameFormURL(url string) string {
	html := GetHTML(url)
	re := regexp.MustCompile(reImg)
	result := re.FindAllStringSubmatch(html, -1)
	for _, result := range result {
		imgURL := result[1]
		fmt.Println(imgURL)
	}
	return ""
}

//GetImgNameFormTag 从<img>标签中提取alt, 有alt就用alt,没有就用时间戳加随机数
func GetImgNameFormTag(imgTag string) string {
	//尝试从imgTag中提取alt
	re := regexp.MustCompile(reAlt)
	result := re.FindAllStringSubmatch(imgTag, 1)
	if len(result) > 0 {
		return result[0][1]
	}
	return GetRandomName()
}
