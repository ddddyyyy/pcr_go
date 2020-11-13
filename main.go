package main

import (
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

func main() {

	app := iris.New()

	app.Use(CounterHandler)

	app.HandleDir("/", "./dist")

	crs := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"}, //允许通过的主机名称
		AllowCredentials:   true,
		Debug:              true,
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: true,
	})

	app.Use(crs)

	v1 := app.Party("/").AllowMethods(iris.MethodGet, iris.MethodPost, iris.MethodOptions)
	{
		v1.Get("get", getEquipmentList)
		v1.Post("send", getEquipmentMap)
		v1.Post("sendRI", getEquipmentMapInfo)
		v1.Get("cache", getMapCache)
		v1.Get("updateHeartHot", updateHeartHot)
		v1.Get("updateEquipments", updateEquipmentData)
		v1.Get("characters", getCharacterListInfo)
		v1.Get("updateCharacters", updateCharacterData)
		v1.Get("application", getApplicationInfo)
	}

	//监听 HTTP/1.x & HTTP/2 客户端在  localhost 端口号8080 设置字符集
	_ = app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))
}

//得到装备列表
func getEquipmentList(ctx iris.Context) {
	_, _ = ctx.JSON(GetList())
}

//确保定时任务每次json文件都有数据更新
func updateHeartHot(ctx iris.Context) {
	_, _ = ctx.JSON(UpdateHeartHot())
}

func getMapCache(ctx iris.Context) {
	_, _ = ctx.JSON(GetMapCache())
}

//从文件重新读取装备信息
func updateEquipmentData(ctx iris.Context) {
	InitEquipmentDataInfoFromJsonFile()
	_, _ = ctx.JSON(equipments)
}

func updateCharacterData(ctx iris.Context) {
	SyncCharacterDataFromExcel()
	_, _ = ctx.JSON(characters)
}

func getCharacterListInfo(ctx iris.Context) {
	_, _ = ctx.JSON(characters)
}

func getApplicationInfo(ctx iris.Context) {
	_, _ = ctx.JSON(ApplicationCache)
}

//得到装备掉落地图信息
func getEquipmentMap(ctx iris.Context) {
	var param Param
	err := ctx.ReadJSON(&param)
	if err == nil {
		cache := make(map[string]int)
		length := len(param.List)
		for index, k := range param.List {
			cache[k] = length - index
		}
		sorts := getSortedKey(param, cache)
		//写入数据
		_, _ = ctx.JSON(sorts.Keys)
		UpdateEquipmentJson()
	}
}

type mapInfo struct {
	Title      string      `json:"title"`
	Equipments []Equipment `json:"equipments"`
}

func getEquipmentMapInfo(ctx iris.Context) {
	var param Param
	err := ctx.ReadJSON(&param)
	if err == nil {
		cache := make(map[string]int)
		length := len(param.List)
		for index, k := range param.List {
			cache[k] = length - index
		}
		sorts := getSortedKey(param, cache)
		var res []mapInfo
		for _, key := range sorts.Keys {
			info := mapInfo{
				Title:      key,
				Equipments: GetCache()[key],
			}
			var i int
			for i = 0; i < len(info.Equipments); i += 1 {
				info.Equipments[i].IsSelect = cache[info.Equipments[i].Title] != 0
			}
			res = append(res, info)
		}
		//写入数据
		_, _ = ctx.JSON(res)
		UpdateEquipmentJson()
	} else {
		fmt.Println(err)
	}
}

func getSortedKey(param Param, cache map[string]int) *ValueSorter {
	equipments := GetListByName(param.List)
	if equipments != nil && len(equipments) > 0 {
		var result = map[string]int64{} // 存放结果的map，包含地图名和权值
		for _, equipment := range equipments {
			UpdateHot(equipment.Id)    //更新热度
			var temp map[string]string //掉落地图
			if err := json.Unmarshal([]byte(equipment.Map), &temp); err == nil {
				for key, value := range temp {
					tV, _ := strconv.ParseInt(value[:len(value)-1], 10, 64)
					tV = tV * int64(cache[equipment.Title])
					tK := strings.Replace(key, "\t", "", -1) //地图名
					if strings.Contains(tK, "36") {
						continue
					}
					if strings.Contains(tK, "VH") || strings.Contains(tK, "HARD") {
						//优先度低于普通地图
						if result[tK] == 0 {
							result[tK] = tV - 10000 //初始化
						} else {
							result[tK] += tV //增加
						}
						continue
					}
					if result[tK] == 0 {
						result[tK] = tV //初始化
					} else {
						result[tK] += tV //增加
					}
				}
			} else {
				Error.Println(err)
			}
		}
		fmt.Println(result)
		//按照value排序
		sorts := NewValueSorter(result)
		sorts.Sort()
		return sorts
	}
	return nil
}
