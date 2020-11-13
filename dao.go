package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var Dao *gorm.DB
var MapCache map[string][]Equipment
var ApplicationCache map[string]interface{}
var equipments []Equipment
var characters []Character

const (
	EquipmentDataFile   = "./equipment.json"
	ApplicationDataFile = "./application.json"
	CharacterDataFile   = "./character.json"
)

func init() {
	MapCache = make(map[string][]Equipment)
	ApplicationCache = make(map[string]interface{})
	var err error

	//文件读取
	filePtr, err := os.Open(EquipmentDataFile)
	if err != nil {
		log.Println("equipment文件打开失败", err.Error())
		return
	}
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&equipments)
	if err != nil {
		log.Println("读取equipment文件错误", err.Error())
	} else {
		log.Println("读取equipment文件成功")
	}

	//文件读取
	filePtr, err = os.Open(ApplicationDataFile)
	if err != nil {
		log.Println("application文件打开失败", err.Error())
		return
	}
	defer filePtr.Close()
	decoder = json.NewDecoder(filePtr)
	err = decoder.Decode(&ApplicationCache)
	if err != nil {
		log.Println("读取application文件错误", err.Error())
	} else {
		log.Println("读取application文件成功")
	}

	//角色信息读取
	filePtr, err = os.Open(CharacterDataFile)
	if err != nil {
		log.Println("character文件打开失败", err.Error())
		return
	}
	defer filePtr.Close()
	decoder = json.NewDecoder(filePtr)
	err = decoder.Decode(&characters)
	if err != nil {
		log.Println("读取character文件错误", err.Error())
	} else {
		log.Println("读取character文件成功")
	}

	// 数据库读取操作
	// 连接数据库
	//Dao, err = gorm.Open("mysql", "root:root@(127.0.0.1:3306)/python?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	panic(err)
	//}
	//err = Dao.Select("id,title,hot,id,map,url").Order("hot DESC").Find(&equipments, "enable = 1").Error

	if err == nil {
		for _, equipment := range equipments {
			var temp map[string]string //掉落地图
			if err := json.Unmarshal([]byte(equipment.Map), &temp); err == nil {
				for key, value := range temp {
					tV, _ := strconv.ParseInt(value[:len(value)-1], 10, 8)
					tK := strings.Replace(key, "\t", "", -1) //地图名
					var tempEquipment = Equipment{
						Id:       equipment.Id,
						Title:    equipment.Title,
						Hot:      equipment.Hot,
						Url:      equipment.Url,
						Priority: tV,
					}
					if MapCache[tK] == nil {
						MapCache[tK] = []Equipment{tempEquipment}
					} else {
						t := append(MapCache[tK], tempEquipment)
						MapCache[tK] = t
					}
				}
			}
		}
		//根据爆率排序
		for _, value := range MapCache {
			sort.Sort(EquipmentSlice(value))
		}
		// 初始化json文件
		//filePtr, err := os.Create(EquipmentDataFile)
		//if err != nil {
		//	log.Println("文件创建失败", err.Error())
		//	return
		//}
		//defer filePtr.Close()
		//encoder := json.NewEncoder(filePtr)
		//err = encoder.Encode(equipments)
		//if err != nil {
		//	log.Println("编码错误", err.Error())
		//} else {
		//	log.Println("编码成功")
		//}
		//filePtr, err = os.Create("./map.json")
		//if err != nil {
		//	log.Println("文件创建失败", err.Error())
		//	return
		//}
		//defer filePtr.Close()
		//encoder = json.NewEncoder(filePtr)
		//err = encoder.Encode(MapCache)
		//if err != nil {
		//	log.Println("编码错误", err.Error())
		//} else {
		//	log.Println("编码成功")
		//}
	}
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

func InitCharacterDataFile(_characters []Character) {
	filePtr, err := os.Create(CharacterDataFile)
	if err != nil {
		log.Println("character文件创建失败", err.Error())
		return
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(_characters)
	if err != nil {
		log.Println("character编码错误", err.Error())
	} else {
		log.Println("character编码成功")
	}
	copy(characters, _characters)
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
func UpdateJson() {
	UpdateApplicationJson()
	UpdateEquipmentJson()
}

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

func GetDao() *gorm.DB {
	return Dao
}

func GetCache() map[string][]Equipment {
	return MapCache
}
