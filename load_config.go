package main

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

func InitConfig() (err error) {
	v := viper.New()
	v.SetConfigName("edns.config.json") // 配置文件名为 "config"
	v.AddConfigPath(".")                // 查找当前目录下的配置文件
	v.AddConfigPath("/etc/dns/")        // 查找 /etc/appname/ 目录下的配置文件
	v.AddConfigPath("/etc/")            // 查找 /etc/ 目录下的配置文件
	v.SetConfigType("json")
	err = v.ReadInConfig() // 读取配置文件
	if err != nil {
		// handle error
		return errors.New("读取配置文件错误！")
	}
	config.DnsPort = v.GetInt("server.port")
	if config.DnsPort == 0 {
		return errors.New("server.port 读取错误！")
	}
	config.DnsListenAddress = v.GetString("server.listenAddress")
	if config.DnsListenAddress == "" {
		return errors.New("server.listenAddress 读取错误！")
	}
	config.ApiPort = v.GetInt("api.port")
	if config.ApiPort == 0 {
		return errors.New("api.port 读取错误！")
	}
	config.ApiListenAddress = v.GetString("api.listenAddress")

	if config.ApiListenAddress == "" {
		return errors.New("api.listenAddress 读取错误！")
	}
	config.user.Email = v.GetString("superuser.email")
	if config.user.Email == "" {
		return errors.New("superuser.email 读取错误！")
	}
	config.user.Passwd = v.GetString("superuser.password")
	if config.user.Passwd == "" {
		return errors.New("superuser.password 读取错误！")
	}
	config.user.Role = 1
	config.user.CreateTime = time.Now().Unix()
	config.Database.DbType = v.GetString("db.type")
	if config.Database.DbType == "" {
		return errors.New("db.type 读取错误！")
	}
	config.Database.Url = v.GetString("db.url")
	if config.Database.Url == "" {
		return errors.New("db.url 读取错误！")
	}

	config.Redis.Server = v.GetString("redis.server")
	if config.Redis.Server == "" {
		return errors.New("redis.server 读取错误！")
	}

	config.Redis.Password = v.GetString("redis.password")
	if config.Redis.Server == "" {
		return errors.New("redis.password 读取错误！")
	}
	config.Debug.Token = v.GetBool("debug.token")
	//if config.Debug.Token ==  {
	//	return errors.New("redis.password 读取错误！")
	//}
	return nil
}
