package main

import (
	"sync"
)

//爬取分页图片

var (
	//MaxChan 最大个数
	MaxChan int = 10
	//chDownImg 下载图片的最大协程个数
	chDownImg = make(chan int8, MaxChan)
	//退出chDownImg协程
	//chDownComplete = make(chan bool, MaxChan)
	downloadWG sync.WaitGroup
	//时间锁
	randomMT sync.Mutex

	//图片下载路径
	imgPath = `E:\Learning\goproject\src\MyProje02\test01\spider_img\test05\imgs\`
	//图片后缀
	imgSuffix = ".jpg"
)




func main() {
	url := "https://www.duotoo.com/zt/rbmn/index.html"
	imginfos := GetImgTagInfo(url)

	for _, imginfoMap := range imginfos {
		DownImgAsync(imginfoMap["url"], imginfoMap["filename"])
	}
}




