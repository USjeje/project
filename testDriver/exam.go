package main

import (
	"MyProje02/test01/testDriver/utils"
	"fmt"
	"strconv"
	"time"
)

var (
	//车道数
	LaneNum = 5
	//车道管道
	chLane = make(chan int8, LaneNum)

	//考试违规名单
	chVio = make(chan string, StuNum)
	//考生是否及格名单,true代表及格，false代表不及格
	passMap = make(map[string]bool, StuNum)
	//学生成绩
	scoreMap = make(map[string]int8, StuNum)
	//记录参加考试考生姓名,方便考生考完试查询成绩
	sliceExamNames = make([]string, 0)
)


//TakeExam 考试
func TakeExam(name string, examNum int) {
	//进入车道准备考试
	chLane<- 1

	fmt.Println(examNum, "号考生", name, "正在考试")

	//将姓名加到姓名切片中，方便考完试查询成绩
	sliceExamNames = append(sliceExamNames, name)

	//考生考试时间
	<-time.After(400 * time.Millisecond)

	//模拟10%的违规几率
	probabilityOfViolation := utils.GetRandomInt(0, 100)
	if probabilityOfViolation < 10 {
		//违规
		fmt.Println(name, "考试违规")
		//放入违规名单
		chVio<- name
		//成绩置为0
		scoreMap[name] = 0
		//标记为不及格
		passMap[name] = false
	} else {
		//不违规
		//随机得到成绩，不违规的有70%的及格率
		num := utils.GetRandomInt(0, 100)
		if num < 70 {
			//70%及格的部分,成绩在70到100之间随机
			scoreMap[name] = int8(utils.GetRandomInt(70, 101))
			//标记为及格
			passMap[name] = true
			//fmt.Println("num=", num, "score = ", scoreMap[name])
		} else {
			//30%的不及格率,成绩在0到69之间随机
			scoreMap[name] = int8(utils.GetRandomInt(0, 70))
			//标记为不及格
			passMap[name] = false
			//fmt.Println("num=", num, "score = ", scoreMap[name])
		}
	}

	//释放资源
	<-chLane
}

//Patrol 考官巡考
func Patrol() {
	//每三秒巡视一次
	ticker := time.NewTicker(2 * time.Second)
	//一直巡考直至主线程结束
	for {
		fmt.Println("监考老师在巡考")
		select {
		case name := <-chVio:
			fmt.Println(name, "考试违规！！！被请出考场")
		default :
			fmt.Println("考场秩序良好！！！")
		}
		<-ticker.C
	}
}

//QueryScore 二级缓存查询成绩
func QueryScore(name string) {
	//先尝试是够能从redis中读取到成绩
	score, err := utils.DoRedisCommand("get "+name, "int")
	//如果没有从redis中查到,就从Mysql中查询，并写入到redis
	if err != nil || score == nil{
		score, err = utils.QueryScoreFromMySql(name)
		//fmt.Println("进行mysql数据库查询，并将结果给redis, name,score", name, score)
		//将成绩写入到redis
		_, err = utils.DoRedisCommand("set "+name+" "+strconv.Itoa(score.(int)), "string")
		utils.HandleErr(err, `WriteScoreToRedis(name, score)`)
	} else {
		fmt.Println("query from Redis", name, ":", score)
	}
}

