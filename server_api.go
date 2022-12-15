package main

import (
	"edns/core"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// JsonResult json返回体
type JsonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func webApiServer() {
	http.HandleFunc("/register", httpRegister)
	http.HandleFunc("/get_dns", httpGetDns)
	http.HandleFunc("/register_dns", httpRegisterDns)
	http.HandleFunc("/dns_delete", httpDnsDelete)
	http.HandleFunc("/get_user_info", httpGetUserInfo)
	http.HandleFunc("/user_delete", httpDelete)
	addr := config.ApiListenAddress + ":" + strconv.Itoa(config.ApiPort)
	log.Println("listen http server: ", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getTime() string {
	currentTime := time.Now()
	result := currentTime.Unix()
	return strconv.Itoa(int(result >> 1))
}

func checkToken(token string) (User, bool) {

	// debug
	//return User{Email: "test", Passwd: "passwd", CreateTime: time.Now().Unix()}, true

	data := strings.Split(core.Base58Decoding(token), ":")
	if len(data) == 2 {
		var user User
		ok := userServer.SearchByEmail(&user, data[0])
		//log.Println("user:", user, err)
		var message string
		if ok {
			if config.Debug.Token {
				message = user.Email + user.Passwd
			} else {
				message = user.Email + user.Passwd + getTime()
			}
			//message := user.Email + user.Passwd + getTime()
			tokenMd5 := core.GenMd5(message)
			if tokenMd5 == data[1] {
				return user, true
			}
		}
		// 没有找到这个用户

	}
	return User{}, false
}

func encodeToken(user string, password string) (token string) {

	message := user + password + getTime()
	tokenMd5 := core.GenMd5(message)
	parseString := user + ":" + tokenMd5
	token = core.Base58Encoding(parseString)
	return token
}

func httpDelete(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		JsonResult
	}
	var response Response

	rootToken := r.PostFormValue("root_token")
	_, ok := checkToken(rootToken)
	if ok {
		email := r.PostFormValue("email")
		var user User
		ok := userServer.SearchByEmail(&user, email)
		if ok {

			ok = userServer.DeleteByEmail(email)
			if ok {
				response.Code = 200
				response.Msg = "删除用户成功"
				data, _ := json.Marshal(response)
				_, _ = w.Write(data)
				return
			} else {
				response.Code = 500
				response.Msg = "删除失败"
				data, _ := json.Marshal(response)
				//_, _ = w.Write(data)
				_, _ = w.Write(data)
				return
			}

		}
		response.Code = 500
		response.Msg = "删除失败,无此用户"
		data, _ := json.Marshal(response)
		_, _ = w.Write(data)
		return

	}
	response.Code = 500
	response.Msg = "认证失败"
	data, _ := json.Marshal(response)
	_, _ = w.Write(data)
	return

}

func httpDnsDelete(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		JsonResult
		Data Dns `json:"data"`
	}
	var response Response
	token := r.PostFormValue("token")
	user, ok := checkToken(token)
	if ok {

		dnsIdString := r.PostFormValue("dns_id")
		dnsId, err := strconv.Atoi(dnsIdString)
		if err == nil {
			var dns Dns
			dnsServer.SearchById(&dns, dnsId)
			// 智能删除属于自己创建的DNS记录
			if dns.UserID == user.Id {
				ok := dnsServer.Delete(dns)
				if ok {
					response.Code = 200
					response.Msg = "删除成功"
					response.Data = dns
					data, _ := json.Marshal(response)
					//_, _ = w.Write(data)
					_, _ = w.Write(data)
					return
				}
				response.Code = 500
				response.Msg = "删除失败"
				response.Data = dns
				data, _ := json.Marshal(response)
				_, _ = w.Write(data)
				return
			}

		}

	}

	response.Code = 500
	response.Msg = "认证失败"
	data, _ := json.Marshal(response)
	//_, _ = w.Write(data)
	_, _ = w.Write(data)
	return

}

func httpRegisterDns(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		JsonResult
		Data Dns `json:"data"`
	}
	var response Response
	token := r.PostFormValue("token")
	user, ok := checkToken(token)
	if ok {
		var dns Dns
		dns.HostRecode = r.PostFormValue("host_recode")
		dns.Domain = r.PostFormValue("domain")
		dns.RecodeType = r.PostFormValue("recode_type")
		dns.RecodeValue = r.PostFormValue("recode_value")
		dns.LastOptionTime = time.Now().Unix()
		dns.UserID = user.Id
		var err error
		dns.TTL, err = strconv.Atoi(r.PostFormValue("ttl"))
		if err == nil {
			ok = dnsServer.Add(dns)
			//log.Println("ok：", ok)
			if ok {
				//response. =
				response.Data = dns
				response.Code = 200
				response.Msg = "添加成功"
				data, _ := json.Marshal(response)
				_, _ = w.Write(data)
				return
			}

			response.Data = dns
			response.Code = 500
			response.Msg = "添加失败，是否存在相同"
			data, _ := json.Marshal(response)
			//_, _ = w.Write(data)
			_, _ = w.Write(data)
			return

		}
		response.Data = dns
		response.Code = 500
		response.Msg = "添加失败，请检查ttl和其他设置！"
		data, _ := json.Marshal(response)
		//_, _ = w.Write(data)
		_, _ = w.Write(data)
		return
	}

	response.Code = 500
	response.Msg = "认证失败"
	data, _ := json.Marshal(response)
	_, _ = w.Write(data)
	return

}

func httpGetUserInfo(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		JsonResult
		Data []User `json:"data"`
	}
	var response Response

	rootToken := r.PostFormValue("root_token")
	_, ok := checkToken(rootToken)

	if ok {

		var users []User
		if userServer.GetAllUser(&users) {
			response.Data = users
			response.Code = 200
			response.Msg = ""
			data, _ := json.Marshal(response)
			_, _ = w.Write(data)
			return
		}
		response.Code = 500
		response.Msg = "获取用户信息失败"
		data, _ := json.Marshal(response)
		_, _ = w.Write(data)
		return
	}
	response.Code = 500
	response.Msg = "认证失败"
	data, _ := json.Marshal(response)
	//_, _ = w.Write(data)
	_, _ = w.Write(data)
	return
}

func httpRegister(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		JsonResult
		User
	}
	var response Response

	rootToken := r.PostFormValue("root_token")
	_, ok := checkToken(rootToken)

	if ok {
		//return
		email := r.PostFormValue("email")
		passwd := r.PostFormValue("passwd")
		if passwd != "" || email != "" {
			if strings.Contains(email, ":") || strings.Contains(passwd, ":") {
				// 不能包含这个字符串
				response.Code = 500
				response.Msg = "参数不能包含':'字符"
				data, _ := json.Marshal(response)
				_, _ = w.Write(data)
				return
			}
			if !userServer.IsExit(email) {
				// 针对重复添加用户
				newUser := User{Email: email, Passwd: passwd, CreateTime: time.Now().Unix()}
				ok := userServer.Add(newUser)
				if ok {
					response.Code = 200
					response.Msg = "用户注册成功"
					//response.Token = token
					response.Email = email
					response.Passwd = passwd
					data, _ := json.Marshal(response)
					_, _ = w.Write(data)
					return
				}
			}
			response.Code = 500
			response.Msg = "用户已注册！"
			data, _ := json.Marshal(response)
			_, _ = w.Write(data)
			return
		}
		response.Code = 500
		response.Msg = "参数不能为空！"
		data, _ := json.Marshal(response)
		_, _ = w.Write(data)
		return
	}
	response.Code = 500
	response.Msg = "认证失败"
	data, _ := json.Marshal(response)
	//_, _ = w.Write(data)
	_, _ = w.Write(data)
	return

}

func httpGetDns(w http.ResponseWriter, r *http.Request) {
	// JsonResult json返回体
	type ResponseTable struct {
		Domain  string `json:"domain"`
		Value   string `json:"value"`
		DnsType string `json:"dns_type"`
	}

	type Response struct {
		JsonResult
		Data []Dns `json:"data"`
	}

	var response Response

	rootToken := r.PostFormValue("root_token")
	_, ok := checkToken(rootToken)
	if ok {
		var ddns []Dns
		dnsServer.GenDbAll(&ddns)

		response.Data = ddns
		response.Code = 200
		response.Msg = ""
		data, _ := json.Marshal(response)
		_, _ = w.Write(data)
		return
	}
	response.Code = 500
	response.Msg = "认证失败"
	data, _ := json.Marshal(response)
	//_, _ = w.Write(data)
	_, _ = w.Write(data)
	return

}

// 得到随机字符串
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImpr  https://colobu.com/2018/09/02/generate-random-string-in-Go/
func genRandDomain(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
