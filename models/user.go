package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// User 用户表 model 定义
type User struct {
	gorm.Model
	UserName	string	`gorm:"column:name;type:varchar(100);unique_index;default:null"`
	Password	string	`gorm:"column:password;type:varchar(100);default:null"`
	Email		string	`gorm:"column:email;type:varchar(100);unique_index;default:null"`
	Avatar		string	`gorm:"column:avatar;type:varchar(100);default:null"`
	Status		string	`sql:"type:ENUM('ENABLE', 'DISABLE')"`
}

// Insert 新增用户
func (user *User) Insert() (userID uint, err error) {

	result := DB.Create(&user)
	userID = user.ID
	if result.Error != nil {
		err = result.Error
	}
	return
}

// FindOne 查询用户详情
func (user *User) FindOne(condition map[string]interface{}) (*User, error) {
	var userInfo User
	r := DB.Select("id, name, email, avatar, password")
		re :=r.Where(condition)
		fmt.Println(userInfo,222222)
		result := re.First(&userInfo)
	fmt.Println(r,999999,re,8888888,result,7777777,userInfo,6666)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	if userInfo.ID > 0 {
		return &userInfo, nil
	}
	return nil, nil
}

// FindAll 获取用户列表
func (user *User) FindAll(pageNum int, pageSize int, condition interface{}) (users []User, err error) {

	result := DB.Offset(pageNum).Limit(pageSize).Select("id", "name", "email").Where(condition).Find(&users)
	err = result.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

// UpdateOne 修改用户
func (user *User) UpdateOne(userID uint, data map[string]interface{}) (*User, error) {
	err := DB.Model(&User{}).Where("id = ?", userID).Updates(data).Error
	if err != nil {
		return nil, err
	}
	var updUser User
	err = DB.Select([]string{"id", "name", "email", "avatar"}).First(&updUser, userID).Error
	if err != nil {
		return nil, err
	}
	return &updUser, nil
}

// DeleteOne 删除用户
func (user *User) DeleteOne(userID uint) (delUser User, err error) {
	if err = DB.Select([]string{"id"}).First(&user, userID).Error; err != nil {
		return
	}

	if err = DB.Delete(&user).Error; err != nil {
		return
	}
	delUser = *user
	return
}
