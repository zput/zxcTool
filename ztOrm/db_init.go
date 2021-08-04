package ztOrm

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"

	"log"
	"time"
)

type DbEngineer interface {
	initMongo() error
	initPostgres() error
	initMysql() error
	initRedis() error

	GetRedisEngine() (*redis.Client, error)
	GetPostgresEngine() (*xorm.Engine, error)
	GetMongoEngine() (*mgo.Session, error)
	GetMysqlEngine() (*xorm.Engine, error)
	InitTable(tableMap map[string]interface{}) error
}

func NewDbEngine(conf config.Configer) DbEngineer {
	return &DbEngine{conf: conf}
}

type DbEngine struct {
	conf           config.Configer
	postgresEngine *xorm.Engine
	mysqlEngine    *xorm.Engine
	redisClient    *redis.Client
	mgoSess        *mgo.Session
}

func (engine *DbEngine) InitTable(tableMap map[string]interface{}) error {
	if engine.postgresEngine == nil {
		if err := engine.initPostgres(); err != nil {
			return err
		}
	}

	for _, value := range tableMap {
		if err := engine.postgresEngine.Sync(value); err != nil {
			return err
		}
	}

	return nil
}

func (engine *DbEngine) initMongo() error {
	urls := engine.conf.String("mongo::urls")
	if urls == "" {
		panic("mongo::urls is invalid")
	}

	mgoDB := engine.conf.String("mongo::db")
	if mgoDB == "" {
		mgoDB = "test"
	}

	mgoSess, err := mgo.Dial(urls)
	if err != nil {
		return fmt.Errorf("initMongo.Err:%w. ", err)
	}

	engine.mgoSess = mgoSess
	DefaultMgoSess = mgoSess

	whetherDebug := engine.conf.String("mongo::debug")
	if whetherDebug == "true" {
		mgo.SetDebug(true)           // 设置DEBUG模式
		mgo.SetLogger(new(MongoLog)) // 设置日志.
	}

	return nil
}

// 实现 mongo.Logger 的接口
type MongoLog struct {
}

func (MongoLog) Output(calldepth int, s string) error {
	log.SetFlags(log.Lshortfile)
	return log.Output(calldepth, s)
}

func (engine *DbEngine) initPostgres() error {
	user := engine.conf.String("postgres::user")
	//user = "postgres"
	passWd := engine.conf.String("postgres::password")
	//passWd = "9dwit1234"

	host := engine.conf.String("postgres::host")
	//host = "127.0.0.1"

	port := engine.conf.DefaultInt("postgres::port", 3306)
	//port = 7000

	db := engine.conf.String("postgres::db")
	//db = "scratch"

	maxIdle := engine.conf.DefaultInt("postgres::maxidle", 3)
	maxOpen := engine.conf.DefaultInt("postgres::maxopen", 20)
	debug := engine.conf.DefaultBool("postgres::debug", true)

	source := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, passWd, db)

	logs.Info("db source: %s", source)

	postgresEngine, err := xorm.NewEngine("postgres", source)
	if err != nil {
		return fmt.Errorf("initPostgres.Err:%w. ", err)
	}

	engine.postgresEngine = postgresEngine
	DefaultPostgresEngine = engine.postgresEngine

	engine.postgresEngine.SetMaxIdleConns(maxIdle)
	engine.postgresEngine.SetMaxOpenConns(maxOpen)
	engine.postgresEngine.ShowSQL(debug)

	if err = postgresEngine.Ping(); err != nil {
		return fmt.Errorf("initPostgres.Err:%v ", err)
	}

	return nil
}

func (engine *DbEngine) initMysql() error {
	return nil
}

func (engine *DbEngine) initRedis() error {
	addr := engine.conf.String("redis::addr")
	//passWd := engine.conf.String("redis::passwd")
	poolSize := engine.conf.DefaultInt("redis::poolsize", 10)

	if addr == "" {
		return fmt.Errorf("initRedis.addrs is not found. ")
	}

	opt := redis.Options{
		Addr:         addr,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
		PoolSize:     poolSize,
	}

	engine.redisClient = redis.NewClient(&opt)
	DefaultRedis = engine.redisClient

	return nil
}

func (engine *DbEngine) GetRedisEngine() (*redis.Client, error) {
	if engine.redisClient != nil {
		return engine.redisClient, nil
	}

	if err := engine.initRedis(); err != nil {
		return nil, err
	}
	return engine.redisClient, nil
}

func (engine *DbEngine) GetPostgresEngine() (*xorm.Engine, error) {
	if engine.postgresEngine != nil {
		return engine.postgresEngine, nil
	}

	if err := engine.initPostgres(); err != nil {
		return nil, err
	}
	return engine.postgresEngine, nil
}

func (engine *DbEngine) GetMongoEngine() (*mgo.Session, error) {
	if engine.mgoSess != nil {
		return engine.mgoSess, nil
	}

	if err := engine.initMongo(); err != nil {
		return nil, err
	}
	return engine.mgoSess, nil
}

func (engine *DbEngine) GetMysqlEngine() (*xorm.Engine, error) {
	if engine.mysqlEngine != nil {
		return engine.mysqlEngine, nil
	}

	if err := engine.initMysql(); err != nil {
		return nil, err
	}
	return engine.mysqlEngine, nil
}

func GetPostGresEngine() *xorm.Engine {
	return DefaultPostgresEngine
}
