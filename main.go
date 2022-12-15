package main

import (
	"gorm.io/gorm"
	"log"
)

var (
	userServer  User
	dnsServer   Dns
	db          *gorm.DB
	redisServer RedisDNS
	config      Config
)

type DatabaseConfig struct {
	DbType string
	Url    string
}

type RedisConfig struct {
	Server   string
	Password string
}

// Config 运行前必须加载的参数
type Config struct {
	DnsPort          int
	DnsListenAddress string
	ApiPort          int
	ApiListenAddress string
	user             User
	Database         DatabaseConfig
	Redis            RedisConfig
	Debug            DebugConfig
}

type DebugConfig struct {
	Token bool
}

func main() {

	//log.Println(core.Base58Encoding("905008677@qq.com:d0c2aa7c83e6554a9ac0406ba57f71fb"))
	//log.Println(core.Base58Encoding("test@qq.com:a552e8af69beed1306d0896cbeb0b12f"))
	err := InitConfig()
	if err != nil {
		panic(err)
		return
	}
	log.Println("加载配置成功！")
	// 初始化数据库
	_, err = InitializeDB()
	if err != nil {
		panic(err)
		return
	}
	log.Println("初始化数据库成功！")
	var vailUser User
	ok := userServer.SearchByEmail(&vailUser, config.user.Email)

	if !ok {
		/*
			管理员用户不存在
			在这里添加用户
		*/
		ok = userServer.Add(config.user)
		if !ok {
			panic("初始化用户失败！")
			return
		}
	} else if vailUser.Passwd != config.user.Passwd || vailUser.Role != 1 {
		/*
			如果用户存在，但是跟现在配置的用户信息不一致，那么把之前的用户删掉，在添加新用户
		*/
		ok = userServer.DeleteByEmail(vailUser.Email)
		if !ok {
			panic("初始化用户失败")
			return
		}
		ok = userServer.Add(config.user)
		if !ok {
			panic("初始化用户失败！")
			return
		}
	}
	// 初始化Redis
	redisServer, err = InitilizeRedis()
	if err != nil {
		panic(err)
		return
	}
	log.Println("初始化Redis成功！")
	go webApiServer()
	serverDNS()

}
