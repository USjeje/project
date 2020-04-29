package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
	"time"
)

var (
	//开放成绩查询读写锁
	dbMutex sync.RWMutex
	//全局数据库链接
	mysqlDB *sql.DB
	//redisPool连接池
	redisPool *redis.Pool

)

//mysql中driver_exam库中score表里的结构
type ExamScore struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Score int	`db:"score"`
}

//InitMysqlDB 得到数据库链接
func InitMysqlDB() {
	if mysqlDB == nil {
		//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
		db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/driving_exam")
		HandleErr(err, `sql.Open`)
		mysqlDB = db
		//defer db.Close()
	}
}

//InitRedisPool 初始化redis连接池
func initRedisPool(net, add string, maxActive, maxIdle int, idleTimeout time.Duration) (*redis.Pool){
	return &redis.Pool{
		MaxActive: maxActive,
		MaxIdle: maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(net, add)
		},
	}
}

//GetRedisConn 得到redis链接
func GetRedisConn() redis.Conn {
	if redisPool == nil {
		redisPool = initRedisPool("tcp","localhost:6379",100,10, 10*time.Second)
	}
	return redisPool.Get()
}

//CloseRedis 关闭redis链接池
func CloseRedis() {
	if redisPool != nil {
		redisPool.Close()
	}
}

//CloseMysqlDB 关闭数据库链接
func CloseMysqlDB() {
	if mysqlDB != nil {
		mysqlDB.Close()
	}
}

//WriteScoreToMysql 记录到MySQL
func WriteScoreToMysql(scoreMap map[string]int8) {
	//枷锁，写入期间不允许读
	dbMutex.Lock()

	//拿到数据库链接
	InitMysqlDB()

	for name, score := range scoreMap {
		_, err := mysqlDB.Exec("insert into score(name,score) values(?,?);", name, score)
		HandleErr(err, `db.Exec`)
		fmt.Println("插入mysql成功")
	}
	fmt.Println("=================成绩录入mysql完毕====================")

	//开放查询
	dbMutex.Unlock()
}

////QueryScore 二级缓存查询成绩
//func QueryScore(name string) {
//	//fmt.Println("QueryScore(name string).name=", name)
//	//先尝试是够能从redis中读取到成绩
//	score, err := QueryScoreFromRedis(name)
//	//如果没有从redis中查到,就从Mysql中查询，并写入到redis
//	if err != nil {
//		score, err = QueryScoreFromMySql(name)
//		//fmt.Println("进行mysql数据库查询，并将结果给redis, name,score", name, score)
//		//将成绩写入到redis
//		err = WriteScoreToRedis(name, score)
//		HandleErr(err, `WriteScoreToRedis(name, score)`)
//	} else {
//		fmt.Println("query from Redis", name, ":", score)
//	}
//}

//QueryScoreFromMySql 考生从MySql中查询自己的成绩
func QueryScoreFromMySql(name string) (int, error) {
	//写入期间不等进行数据库的读访问,读入期间不能写但可读
	dbMutex.RLock()

	//拿到数据库链接
	InitMysqlDB()

	ExamScores := make([]ExamScore, 0)
	//执行查询语句 Query返回一行或多行结果
	//如果不重名的话，返回一行数据
	//这里就不区分了，本来应该是按照考生id而不是姓名来进行一系列的操作
	rows, err := mysqlDB.Query("select * from score where name=?", name)
	//HandleErr(err, `db.Query`)
	if err != nil {
		fmt.Println(err, `db.Query`)
		//-1代表MySql查询出错
		return -1, err
	}
	//循环读取结果
	for rows.Next() {
		var examScore ExamScore
		//将每一行的结果都赋值到一个ExamScore对象中
		err := rows.Scan(&examScore.ID, &examScore.Name, &examScore.Score)
		//fmt.Println(examScore.ID, examScore.Name, examScore.Score)
		HandleErr(err, `rows.Scan`)
		//追加到切片中
		ExamScores = append(ExamScores, examScore)
	}
	fmt.Println("从Mysql查询结果为：", ExamScores)
	dbMutex.RUnlock()

	return ExamScores[0].Score, nil
}

//DoRedisCommand 通用的执行redis命令
//resultType : string, int, strings
func DoRedisCommand(cmd string, resultType string) (interface{}, error) {
	conn := GetRedisConn()
	//conn, err := redis.Dial("tcp", "localhost:6379")
	//HandleErr(err, `redis.Dial`)
	//defer conn.Close()

	strs := strings.Split(cmd, " ")
	args := make([]interface{}, 0)
	for _, arg := range strs[1:] {
		args = append(args, arg)
	}
	reply, err := conn.Do(strs[0], args...)
	HandleErr(err, `conn.Do(strs[0], args...)`)
	if reply != nil {
		//进行类型转变
		switch resultType {
			case "string" :
				return redis.String(reply, err)
			case "int" :
				return redis.Int(reply, err)
			case "strings" :
				return redis.Strings(reply, err)
			default:
				return redis.Strings(reply, err)
		}
	}
	return -2, errors.New("未能从redis中查到数据")

}

//QueryScoreFromRedis 从redis中获取考生成绩
func QueryScoreFromRedis(name string) (int, error) {
	//获取redis连接
	conn, err := redis.Dial("tcp", "localhost:6379")
	HandleErr(err, `redis.Dial`)
	defer conn.Close()

	//在redis中存string类型,查询不到的reply为nil，所以要进行判断
	reply, err := conn.Do("get", name)
	if reply != nil {
		//reply是空接口类型，成绩是int类型
		score, err := redis.Int(reply, err)
		if err != nil {
			fmt.Println(`conn.Do或者redis.Int ERR`, err)
			return -2, err
		}
		return score, nil
	} else {
		return -2, errors.New("未能从redis中查到数据")
	}
}

func WriteScoreToRedis(name string, score int) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	HandleErr(err, `redis.Dial`)
	defer conn.Close()

	//写入到redis的string类型
	conn.Do("set", name, score)
	HandleErr(err, `conn.Do("set", name, score)`)
	fmt.Println("插入redis成功")
	return err
}