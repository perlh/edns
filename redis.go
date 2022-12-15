package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
)

type RedisDNS struct {
	c redis.Conn
}

// RedisLogin RedisConnect redis-11201.c54.ap-northeast-1-2.ec2.cloud.redislabs.com:11201,5n0YYY2ogDQYClgvffuRFcBHWLq6TpmL
func RedisLogin(serverURL string, password string) (c redis.Conn, err error) {
	// 连接redis
	c, err = redis.Dial("tcp", serverURL)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return c, errors.New("Connect to redis error")
	} else {
		//fmt.Println("Connect to redis ok")
	}

	// 密码鉴权
	_, err = c.Do("AUTH", password)
	if err != nil {
		fmt.Println("auth failed:", err)
		return c, errors.New("auth failed")
	} else {
		//fmt.Println("auth ok:")
	}
	return c, nil
}

func (r *RedisDNS) SetDns(dns Dns, second int) (err error) {
	var key string
	var value []byte
	switch dns.RecodeType {
	case "a":
		{
			jsonvalue, _ := json.Marshal(dns)
			key = strings.Join([]string{dns.HostRecode, dns.Domain}, ".")
			value = jsonvalue
			break
		}
	case "aaaa":
		{
			jsonvalue, _ := json.Marshal(dns)
			key = strings.Join([]string{dns.HostRecode, dns.Domain}, ".")
			value = jsonvalue
			break
		}
	case "cname":
		{
			jsonvalue, _ := json.Marshal(dns)
			key = strings.Join([]string{dns.HostRecode, dns.Domain}, ".")
			value = jsonvalue
			break
		}

	}

	// 就算存在，也会覆盖掉之前的值
	// 写入数据
	_, err = r.c.Do("SETNX", key, value)
	if err != nil {
		fmt.Println("redis set failed:", err)
		return err
	} else {
		//fmt.Println("redis set ok")
	}
	// second小于0表示永久生效
	if second > 0 {
		// 设置过期时间为6秒
		ret, err := r.c.Do("EXPIRE", "blog.hsm.cool", second)
		if ret == int64(1) {
			//fmt.Println("success")
		} else {
			return err
		}
	}
	return nil
}
func (r *RedisDNS) isKeyExit(key string) (bool, error) {
	// 判断key是否存在
	is_key_exit, err := redis.Bool(r.c.Do("EXISTS", key))
	if err != nil {
		fmt.Println("error:", err)
		return false, err
	} else {
		fmt.Printf("exists or not: %v \n", is_key_exit)
		return is_key_exit, nil
	}
}

func (r *RedisDNS) DeleteDns(dns Dns) (ok bool) {
	domain, err := dnsServer.GetDomain(dns)
	if err != nil {
		return false
	}
	_, err = r.c.Do("DEL", domain)
	if err != nil {
		fmt.Println("redis delelte failed:", err)
		return false
	}
	return true
}

func (r *RedisDNS) ReadDnsByDomain(domain string) Dns {
	//log.Println("key:", domain)
	value, err := redis.Bytes(r.c.Do("GET", domain))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		//fmt.Printf("Get: %v \n", value)
	}
	var dns11 Dns
	err2 := json.Unmarshal(value, &dns11)
	if err2 != nil {
		fmt.Println(err)
	}
	//log.Println("value:", dns11.RecodeValue)
	return dns11
}

func (r *RedisDNS) ReadDns(dns Dns) Dns {
	var key string
	switch dns.RecodeType {
	case "a":
		{
			key = strings.Join([]string{dns.HostRecode, dns.Domain}, ".")
			break
		}
	case "aaaa":
		{
			key = strings.Join([]string{dns.HostRecode, dns.Domain}, ".")
			break
		}
	case "cname":
		{
			key = strings.Join([]string{dns.HostRecode, dns.Domain}, ".")
			break
		}
	default:
		{
			key = "google.com"
			break
		}
	}
	value, err := redis.Bytes(r.c.Do("GET", key))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		//fmt.Printf("Get: %v \n", value)
	}
	var dns11 Dns
	err2 := json.Unmarshal(value, &dns11)
	if err2 != nil {
		fmt.Println(err)
	}
	return dns11
}
