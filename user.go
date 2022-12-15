package main

import "log"

type User struct {
	Id         int `gorm:"primaryKey;AUTO_INCREMENT"`
	Email      string
	Passwd     string
	Role       int // Role为1表示管理员，0为普通用户
	CreateTime int64
}

func (u *User) Add(user User) (ok bool) {
	if user.Email == "" || user.Passwd == "" || user.CreateTime == 0 {
		// 不能为空
		return false
	}
	if u.IsExit(user.Email) {
		// 如果用户存在的话，false
		return false
	}
	//log.Println("test")
	db.Create(&user)
	return true
}

func (u *User) DeleteByEmail(email string) (ok bool) {
	//var user User
	ok = u.IsExit(email)
	if ok {
		// 删除DNS表

		// 删除用户表
		db.Model(&User{}).Where("email = ?", email).Delete(&User{})
		return true

	}
	return false
}

func (u *User) SearchByEmail(user *User, email string) (ok bool) {
	db.Model(&User{}).Where("email = ?", email).Find(user)

	if user.Id == 0 || user.Email == "" || user.Passwd == "" {
		//log.Println("用户不存在")
		return false
	}
	return true
}

func (u *User) IsExit(email string) (ok bool) {
	var user User
	db.Model(&User{}).Where("email = ?", email).Find(&user)
	if user.Id != 0 || user.Email != "" || user.Passwd != "" {
		return true
	}
	return false
}

func (u *User) Update(oldUser User, newUser User) (ok bool) {
	if u.IsExit(oldUser.Email) {
		if u.IsNull(newUser) {
			// 新用户不能为空
			log.Println("新用户不能为空")
			return false
		}
		if u.IsExit(newUser.Email) && oldUser.Email != newUser.Email {
			// 不能替换别的用户
			log.Println("不能替换别的用户")
			return false
		}
		var user User
		db.Model(&User{}).Where("email = ?", oldUser.Email).Find(&user)
		user.Passwd = newUser.Passwd
		user.Email = newUser.Email
		user.CreateTime = newUser.CreateTime
		db.Save(&user)
		return true
	}
	//用户不存在
	return false
}

func (u *User) GetAllUser(users *[]User) (ok bool) {

	db.Model(&User{}).Find(users)
	return true
}

func (u *User) IsNull(user User) (ok bool) {
	if user.Email == "" || user.Passwd == "" || user.CreateTime == 0 {
		return true
	}
	return false
}
