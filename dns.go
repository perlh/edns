package main

import (
	"errors"
	"log"
)

type Dns struct {
	Id             int `gorm:"primaryKey;AUTO_INCREMENT"`
	Domain         string
	HostRecode     string
	RecodeType     string
	RecodeValue    string
	TTL            int
	LastOptionTime int64
	UserID         int
}

func (d *Dns) IsNotFull(dns Dns) (ok bool) {
	if dns.HostRecode == "" || dns.RecodeValue == "" || dns.Domain == "" || dns.LastOptionTime == 0 {
		return true
	}
	//log.Println(dns.RecodeType)
	if dns.RecodeType == "" {
		return true
	}
	return false
}

func (d *Dns) Add(dns Dns) (ok bool) {
	// 判断需要添加到dns信息是否为空
	//log.Println(dns)
	if !d.IsNotFull(dns) {
		// 如果存在，那么会执行if下面的代码
		if d.IsExit(dns) {
			//log.Println("存在")
			// 发现存在该域名，判断之前的域名是否为本用户所属
			var oldDns Dns
			d.SearchByFullDomain(&oldDns, dns.Domain, dns.HostRecode)
			if oldDns.UserID == dns.UserID {
				dnsServer.DeleteById(oldDns.Id)
			} else {
				return false
			}
		}
		// 添加到redis
		//log.Println("redis")
		err := redisServer.SetDns(dns, -1)
		if err != nil {
			log.Println("添加到dns错误")
			return false
		}
		//log.Println("test")
		db.Create(&dns)
		return true
	}

	//log.Println("dsafas ")

	return false

}

//
func (d *Dns) SearchByDomain(dnss *[]Dns, domain string) (ok bool) {
	db.Model(&Dns{}).Where("domain = ?", domain).Find(dnss)
	if len(*dnss) == 0 {
		return false
	}
	return true
}

func (d *Dns) SearchByFullDomain(dns *Dns, domain string, hostRecode string) (ok bool) {
	db.Model(&Dns{}).Where("domain = ?", domain).Where("host_recode", hostRecode).Find(dns)
	return true
}

func (d *Dns) SearchDb(dnss *[]Dns, dns Dns) (ok bool) {
	//var dnss []Dns
	tx := db.Model(&Dns{}).Where("domain = ?", dns.Domain).Where("host_recode = ?", dns.HostRecode)
	tx.Where("recode_type", dns.RecodeType).Where("recode_value", dns.RecodeValue).Find(&dnss)
	if len(*dnss) == 0 {
		return false
	}
	return true
}

func (d *Dns) GenDbAll(dnss *[]Dns) (ok bool) {
	db.Model(&Dns{}).Find(&dnss)
	return true
}

func (d *Dns) SearchById(dns *Dns, id int) (ok bool) {
	db.Model(&Dns{}).Where("id=?", id).Find(&dns)
	return true
}

func (d *Dns) DeleteById(id int) (ok bool) {
	if id <= 0 {
		return false
	}
	var dns Dns
	d.SearchById(&dns, id)
	//log.Println(dns)
	if redisServer.DeleteDns(dns) {
		db.Model(&Dns{}).Where("id = ?", dns.Id).Delete(&Dns{})
		//log.Println("删除成功")
		return true
	}
	return false

}

// DeleteDb 删除数据库
func (d *Dns) Delete(dns Dns) (ok bool) {
	if redisServer.DeleteDns(dns) {
		db.Model(&Dns{}).Where("id = ?", dns.Id).Delete(&Dns{})
		return true
	}
	return false
}

// DeleteDbByUID DeleteDb 删除数据库
func (d *Dns) DeleteByUID(userID string) (ok bool) {

	db.Model(&Dns{}).Where("user_id = ?", userID).Delete(&Dns{})
	//log.Println(row)
	return true
}

func (d *Dns) IsExit(dns Dns) (ok bool) {
	var dnss []Dns
	if d.SearchByDomain(&dnss, dns.Domain) {
		//log.Println(dnss)
		for _, value := range dnss {
			if value.Domain == dns.Domain && value.HostRecode == dns.HostRecode {
				return true
			}
		}
	}
	return false
}

func (d *Dns) GetDomain(dns Dns) (domain string, err error) {
	switch dns.RecodeType {
	case "a":
		{
			domain = dns.HostRecode + "." + dns.Domain
			return domain, nil
		}
	case "data":
		{
			domain = dns.HostRecode + "." + dns.Domain
			return domain, nil
			//break
		}
	case "cname":
		{
			domain = dns.HostRecode + "." + dns.Domain
			return domain, nil
		}
	default:

		break

	}

	return "", errors.New("error")

}
