package main

import "sort"

//Equipment 装备结构体
type Equipment struct {
	Id       int64  `json:"id" gorm:"primary_key"`
	Title    string `json:"title"`
	Kind     string `json:"kind"`
	Map      string `json:"map"`
	Enable   int    `json:"enable"`
	Hot      int64  `json:"hot"`
	Url      string `json:"url"`
	Priority int64  `json:"priority"`
	IsSelect bool   `json:"isSelect"`
}

//角色信息
type Character struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	RealName       string `json:"realName"`
	Distance       int64  `json:"distance"`
	Rank           string `json:"rank"`
	Star           string `json:"star"`
	UnionEquipment string `json:"unionEquipment"`
	Inform         string `json:"inform"`
}

// 请求的参数数组
type Param struct {
	List []string `json:"list"`
}

func (Equipment) TableName() string {
	return "equipment"
}

func (Character) TableName() string {
	return "character"
}

// map排序
type ValueSorter struct {
	Keys   []string
	Values []int64
}

// map按value排序
func NewValueSorter(m map[string]int64) *ValueSorter {
	vs := &ValueSorter{
		Keys:   make([]string, 0, len(m)),
		Values: make([]int64, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Values = append(vs.Values, v)
	}
	return vs
}
func (vs *ValueSorter) Len() int { return len(vs.Values) }

func (vs *ValueSorter) Sort() {
	sort.Sort(vs)
}
func (vs *ValueSorter) Swap(i, j int) {
	vs.Values[i], vs.Values[j] = vs.Values[j], vs.Values[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
func (vs *ValueSorter) Less(i, j int) bool { return vs.Values[i] > vs.Values[j] }

//根据爆率排序
type EquipmentSlice []Equipment

func (a EquipmentSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a EquipmentSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a EquipmentSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Priority < a[i].Priority
}

//根据热度排序
type EquipmentHotSlice []Equipment

func (a EquipmentHotSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a EquipmentHotSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a EquipmentHotSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Hot < a[i].Hot
}
