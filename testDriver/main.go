package main

/*
   考场签到，名字丢入管道
   只有5个车道，最多供5个人同时考试
   考生按签到顺序依次考试，给予考生10%的违规几率,违规者成绩置为0
   每三秒巡视一次，发现违规的清楚考场，否则输出考试时序良好
   没有违规的考试及格率为约70%
   所有考试看完考完后，生成考试记录
   当前目录下的成绩录入MySQL数据库，数据库允许一写多读
   成绩录入完毕通知考生，考生查阅自己的成绩
*/

import (
	"MyProje02/test01/testDriver/utils"
	"fmt"
	"sync"
	"time"
)

var (
	//wg 互斥锁
	wg sync.WaitGroup

	//StuNum 考生人数
	StuNum = 20
	//chName 考生姓名管道
	chName = make(chan string, StuNum)

)

func main() {

	//延迟关闭数据库
	defer func() {
		utils.CloseMysqlDB()
		utils.CloseRedis()
	}()


	//将名字放到chName管道中
	utils.InitName(chName, StuNum)

	//巡考
	go func () {
		Patrol()
	}()

	id := 0
	//考生并发考试
	for name := range chName {
		wg.Add(1)
		id++
		go func(name string, id int) {
			TakeExam(name, id)
			wg.Done()
		}(name, id)
	}

	//等待考生考试结束
	wg.Wait()

	//考生考试成绩为
	fmt.Println("\n本次考生成绩如下")
	for i,v := range scoreMap {
		fmt.Println(i, v)
	}

	//将成绩录入mysql
	wg.Add(1)
	go func() {
		utils.WriteScoreToMysql(scoreMap)
		wg.Done()
	}()

	fmt.Println("=================成绩录入mysql====================")
	//故意等待一段时间，确保写数据库抢到读写锁
	fmt.Println(sliceExamNames)
	time.Sleep(1*time.Second)


	//考生第一次查询自己的成绩
	for _,name := range sliceExamNames {
		wg.Add(1)
		//这里隐含函数必须把name传进去，不做的话，因为是协程，会导致一个名字被好几个协程当成参数传入QueryScore
		go func(name string) {
			QueryScore(name)
			wg.Done()
		}(name)
	}
	<-time.After(1*time.Second)
	//考生查询自己的成绩
	for _,name := range sliceExamNames {
		wg.Add(1)
		go func(name string) {
			QueryScore(name)
			wg.Done()
		}(name)
	}

	wg.Wait()

	fmt.Println("=================查询成绩结束====================")
	fmt.Println("END...")
}

