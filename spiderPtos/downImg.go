package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//DownImg 同步下载图片
func DownImg(imgURL, filename string) {
	var flag bool

	resp, err := http.Get(imgURL)
	HandleErr(err, `http.Get(imgURL)`)
	defer resp.Body.Close()
	imgBytes, err := ioutil.ReadAll(resp.Body)
	HandleErr(err, `ioutil.ReadAll(resp.Body)`)
	//filename 文件路径加文件名
	filePath := imgPath + filename + imgSuffix

	//写入文件
	err = ioutil.WriteFile(filePath, imgBytes, 0666)
	//有时alt的属性不能当文件名,所以换个名字再试试
	if err != nil {
		fmt.Println(`ioutil.WriteFile error : `, err)
		//fmt.Println(imgURL)
		filename = GetRandomName()
		filePath = imgPath + filename + imgSuffix
		err1 := ioutil.WriteFile(filePath, imgBytes, 0666)
		if err1 != nil {
			fmt.Println("下载再次失败，Error : ", err1)
			flag = true
		}
	}
	if !flag {
		fmt.Println("Download success : ", filename)
	}
}

//DownImgAsync 异步下载图片
func DownImgAsync(url, filename string) {
	downloadWG.Add(1)
	//开启协程
	go func() {
		chDownImg <- 1
		DownImg(url, filename)
		<-chDownImg
		downloadWG.Done()
	}()
	downloadWG.Wait()
}
