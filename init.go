package main

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	var err error

	// 判断数据库类型
	if config.Database.DbType == "mysql" {
		// 连接到 MySQL 数据库
		//user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
		//db, err = gorm.Open("mysql", &gorm.Config{config.Database.Url})
		// TODO
	} else if config.Database.DbType == "postgres" {
		// 连接到 PostgreSQL 数据库
		// "postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword"
		// TODO
	} else if config.Database.DbType == "sqlite" {
		db, err = gorm.Open(sqlite.Open(config.Database.Url), &gorm.Config{})
	}
	if err != nil {
		panic("failed to connect database")
		return db, err
	} else {
		// Migrate the schema
		_ = db.AutoMigrate(&User{})
		_ = db.AutoMigrate(&Dns{})
		return db, err
	}

}

func InitilizeRedis() (RedisDNS, error) {
	redisConnect2, err := RedisLogin(config.Redis.Server, config.Redis.Password)
	if err != nil {
		return RedisDNS{}, errors.New("failed to connect redis")
	} else {

		redisServer.c = redisConnect2
		return redisServer, nil
	}

}
