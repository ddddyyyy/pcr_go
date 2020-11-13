package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"sort"
	"strings"
)

var Dao *gorm.DB

func init() {

	InitEquipmentDataInfoFromJsonFile()

	// 数据库读取操作
	// 连接数据库
	//Dao, _ = gorm.Open("mysql", "root:root@(127.0.0.1:3306)/python?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	panic(err)
	//}
	//err = Dao.Select("id,title,hot,id,map,url").Order("hot DESC").Find(&equipments, "enable = 1").Error

}

//GetList 得到装备列表
func GetList() []Equipment {
	return equipments
	//var equipments []Equipment
	//if Dao.Select("id,title,hot,id,url").Order("hot DESC").Find(&equipments, "enable = 1").Error == nil {
	//	return equipments
	//} else {
	//	Warning.Println(Dao.Error)
	//}
	//
	//return nil
}

func GetListByName(keys []string) []Equipment {
	var res []Equipment
	for _, e := range equipments {
		for _, key := range keys {
			if strings.Compare(e.Title, key) == 0 {
				res = append(res, e)
				break
			}
		}
	}
	return res
	//Dao.Select("id,map,title").Where("title IN (?)", keys).Find(&equipments).Error
}

func GetMapCache() map[string][]Equipment {
	return MapCache
}

func UpdateHeartHot() string {
	for i := 0; i < len(equipments); i += 1 {
		if equipments[i].Id == 1 {
			if equipments[i].Hot != 0 {
				equipments[i].Hot = 0
			} else {
				equipments[i].Hot = -1
			}
			break
		}
	}
	UpdateEquipmentJson()
	return "success"
}

//UpdateHot 更新装备的热度
func UpdateHot(id int64) {
	//文件操作
	for i := 0; i < len(equipments); i += 1 {
		if equipments[i].Id == id {
			equipments[i].Hot++
			break
		}
	}
	//数据库操作
	//var obj Equipment
	//err := Dao.Find(&obj, "id = ?", id).Error
	//if err == nil {
	//	obj.Hot++
	//	Dao.Model(&obj).Update("hot", obj.Hot)
	//} else {
	//	Warning.Println(err)
	//}
}

//更新本地json文件
func UpdateEquipmentJson() {
	//根据热度排序
	sort.Sort(EquipmentHotSlice(equipments))

	//保存equipment
	filePtr, err := os.OpenFile(EquipmentDataFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println("equipment文件读取失败", err.Error())
		return
	}
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(equipments)
	if err != nil {
		log.Println("更新equipment文件错误", err.Error())
	} else {
		log.Println("更新equipment文件成功")
	}
	_ = filePtr.Close()

	//保存application
	filePtr, err = os.OpenFile(ApplicationDataFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println("application文件读取失败", err.Error())
		return
	}
	encoder = json.NewEncoder(filePtr)
	err = encoder.Encode(ApplicationCache)
	if err != nil {
		log.Println("更新application文件错误", err.Error())
	} else {
		log.Println("更新application文件成功")
	}
	_ = filePtr.Close()

}

func UpdateApplicationJson() {
	//保存application
	filePtr, err := os.OpenFile(ApplicationDataFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println("application文件读取失败", err.Error())
		return
	}
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(ApplicationCache)
	if err != nil {
		log.Println("更新application文件错误", err.Error())
	} else {
		log.Println("更新application文件成功")
	}
	_ = filePtr.Close()

}

//func GetDao() *gorm.DB {
//	return Dao
//}

func GetCache() map[string][]Equipment {
	return MapCache
}
