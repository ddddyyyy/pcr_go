package main

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/semicircle/gozhszht"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	EquipmentDataFile   = "./equipment.json"
	ApplicationDataFile = "./application.json"
	CharacterDataFile   = "./character.json"
)

var MapCache map[string][]Equipment
var ApplicationCache map[string]interface{}
var equipments []Equipment
var characters []Character

func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func t2s(text string) string {
	return gozhszht.ToSimple(text)
	//cc, err := gocc.New("t2s")
	//if err != nil {
	//	log.Println(err)
	//	return ""
	//}
	//res, err := cc.Convert(text)
	//if err != nil {
	//	log.Println(err)
	//	return ""
	//}
	//return res
}

//根据excel同步角色信息
func SyncCharacterDataFromExcel() {
	characters := map[int]interface{}{
		1000: []string{"未知角色", "未知キャラ", "Unknown"},
		1001: []string{"日和", "ヒヨリ", "Hiyori", "日和莉", "猫拳", "🐱👊"},
		1002: []string{"优衣", "ユイ", "Yui", "种田", "普田", "由衣", "结衣", "ue", "↗↘↗↘"},
		1003: []string{"怜", "レイ", "Rei", "剑圣", "普怜", "伶"},
		1004: []string{"禊", "ミソギ", "Misogi", "未奏希", "炸弹", "炸弹人", "💣"},
		1005: []string{"茉莉", "マツリ", "Matsuri", "跳跳虎", "老虎", "虎", "🐅"},
		1006: []string{"茜里", "アカリ", "Akari", "妹法", "妹妹法"},
		1007: []string{"宫子", "ミヤコ", "Miyako", "布丁", "布", "🍮"},
		1008: []string{"雪", "ユキ", "Yuki", "小雪", "镜子", "镜法", "伪娘", "男孩子", "男孩纸", "雪哥"},
		1009: []string{"杏奈", "アンナ", "Anna", "中二", "煤气罐"},
		1010: []string{"真步", "マホ", "Maho", "狐狸", "真扎", "咕噜灵波", "真布", "🦊"},
		1011: []string{"璃乃", "リノ", "Rino", "妹弓"},
		1012: []string{"初音", "ハツネ", "Hatsune", "hego", "星法", "星星法", "⭐法", "睡法"},
		1013: []string{"七七香", "ナナカ", "Nanaka", "娜娜卡", "77k", "77香"},
		1014: []string{"霞", "カスミ", "侦探", "杜宾犬", "驴", "驴子", "🔍"},
		1015: []string{"美里", "ミサト", "Misato", "圣母"},
		1016: []string{"铃奈", "スズナ", "Suzuna", "暴击弓", "暴弓", "爆击弓", "爆弓", "政委"},
		1017: []string{"香织", "カオリ", "Kaori", "琉球犬", "狗子", "狗", "狗拳", "🐶", "🐕", "🐶👊🏻", "🐶👊"},
		1018: []string{"伊绪", "イオ", "Io", "老师", "魅魔"},

		1020: []string{"美美", "ミミ", "Mimi", "兔子", "兔兔", "兔剑", "萝卜霸断剑", "人参霸断剑", "天兔霸断剑", "🐇", "🐰"},
		1021: []string{"胡桃", "クルミ", "Kurumi", "铃铛", "🔔"},
		1022: []string{"依里", "ヨリ", "Yori", "姐法", "姐姐法"},
		1023: []string{"绫音", "アヤネ", "Ayane", "熊锤", "🐻🔨", "🐻"},

		1025: []string{"铃莓", "スズメ", "Suzume", "女仆", "妹抖"},
		1026: []string{"铃", "リン", "Rin", "松鼠", "🐿", "🐿️"},
		1027: []string{"惠理子", "エリコ", "Eriko", "病娇"},
		1028: []string{"咲恋", "サレン", "Saren", "充电宝", "青梅竹马", "幼驯染", "院长", "园长", "🔋", "普电"},
		1029: []string{"望", "ノゾミ", "Nozomi", "偶像", "小望", "🎤"},
		1030: []string{"妮诺", "扇子", "ニノン", "Ninon", "妮侬"},
		1031: []string{"忍", "シノブ", "Shinobu", "普忍", "鬼父", "💀"},
		1032: []string{"秋乃", "アキノ", "Akino", "哈哈剑"},
		1033: []string{"真阳", "マヒル", "Mahiru", "奶牛", "🐄", "🐮", "真☀"},
		1034: []string{"优花梨", "ユカリ", "Yukari", "由加莉", "黄骑", "酒鬼", "奶骑", "圣骑", "🍺", "🍺👻"},

		1036: []string{"镜华", "キョウカ", "Kyouka", "小仓唯", "xcw", "小苍唯", "8岁", "八岁", "喷水萝", "八岁喷水萝", "8岁喷水萝"},
		1037: []string{"智", "トモ", "Tomo", "卜毛"},
		1038: []string{"栞", "シオリ", "Shiori", "tp弓", "小栞", "白虎弓", "白虎妹"},

		1040: []string{"碧", "アオイ", "Aoi", "香菜", "香菜弓", "绿毛弓", "毒弓", "绿帽弓", "绿帽"},

		1042: []string{"千歌", "チカ", "Chika", "绿毛奶"},
		1043: []string{"真琴", "マコト", "Makoto", "狼", "🐺", "月月", "朋", "狼姐"},
		1044: []string{"伊莉亚", "イリヤ", "Iriya", "伊利亚", "伊莉雅", "伊利雅", "yly", "吸血鬼", "那个女人"},
		1045: []string{"空花", "クウカ", "Kuuka", "抖m", "抖"},
		1046: []string{"珠希", "タマキ", "Tamaki", "猫剑", "🐱剑", "🐱🗡️"},
		1047: []string{"纯", "ジュン", "Jun", "黑骑", "saber"},
		1048: []string{"美冬", "ミフユ", "Mifuyu", "子龙", "赵子龙", "泳装美冬"},
		1049: []string{"静流", "シズル", "Shizuru", "姐姐"},
		1050: []string{"美咲", "ミサキ", "Misaki", "大眼", "👀", "👁️", "👁"},
		1051: []string{"深月", "ミツキ", "Mitsuki", "眼罩", "抖s", "医生"},
		1052: []string{"莉玛", "リマ", "Rima", "Lima", "草泥马", "羊驼", "🦙", "🐐"},
		1053: []string{"莫妮卡", "モニカ", "Monika", "毛二力"},
		1054: []string{"纺希", "ツムギ", "Tsumugi", "裁缝", "蜘蛛侠", "🕷️", "🕸️"},
		1055: []string{"步未", "アユミ", "Ayumi", "步美", "路人", "路人妹"},
		1056: []string{"流夏", "ルカ", "Ruka", "大姐", "大姐头", "儿力", "luka", "刘夏"},
		1057: []string{"吉塔", "ジータ", "Jiita", "姬塔", "团长", "吉他", "🎸", "骑空士", "qks"},
		1058: []string{"贪吃佩可", "ペコリーヌ", "Pecoriinu", "佩可莉姆", "吃货", "佩可", "公主", "饭团", "🍙"},
		1059: []string{"可可萝", "コッコロ", "Kokkoro", "可可罗", "妈", "普白"},
		1060: []string{"凯留", "キャル", "Kyaru", "凯露", "百地希留耶", "希留耶", "Kiruya", "黑猫", "臭鼬", "普黑", "接头霸王", "街头霸王"},
		1061: []string{"矛依未", "ムイミ", "Muimi", "诺维姆", "Noemu", "夏娜", "511", "无意义", "天楼霸断剑", "姆咪", "母咪"},

		1063: []string{"亚里莎", "アリサ", "Arisa", "鸭梨瞎", "瞎子", "亚里沙", "鸭梨傻", "亚丽莎", "亚莉莎", "瞎子弓", "🍐🦐", "yls"},

		1065: []string{"嘉夜", "カヤ", "Kaya", "憨憨龙", "龙拳", "🐲👊🏻", "🐉👊🏻", "接龙笨比", "鬼道嘉夜"},
		1066: []string{"祈梨", "イノリ", "Inori", "梨老八", "李老八", "龙锤", "🐲🔨"},
		1067: []string{"穗希", "ホマレ", "Homare"},
		1068: []string{"拉比林斯达", "ラビリスタ", "Rabirisuta", "迷宫女王", "模索路晶", "模索路", "晶"},
		1069: []string{"真那", "マナ", "Mana", "霸瞳皇帝", "千里真那", "千里", "霸瞳", "霸铜"},
		1070: []string{"似似花", "ネネカ", "Neneka", "变貌大妃", "现士实似似花", "現士実似々花", "現士実", "现士实", "nnk", "448", "捏捏卡", "变貌", "大妃"},
		1071: []string{"克莉丝提娜", "クリスティーナ", "Kurisutiina", "誓约女君", "克莉丝提娜·摩根", "Christina", "Cristina", "克总", "女帝", "克", "摩根"},
		1072: []string{"可萝爹", "長老", "Chourou", "岳父", "爷爷"},
		1073: []string{"拉基拉基", "ラジニカーント", "Rajinigaanto", "跳跃王", "Rajiraji", "Lajilaji", "垃圾垃圾", "教授"},

		1075: []string{"贪吃佩可(夏日)", "ペコリーヌ(サマー)", "Pekoriinu(Summer)", "佩可莉姆(夏日)", "水吃", "水饭", "水吃货", "水佩可", "水公主", "水饭团", "水🍙", "泳吃", "泳饭", "泳吃货", "泳佩可", "泳公主", "泳饭团", "泳🍙", "泳装吃货", "泳装公主", "泳装饭团", "泳装🍙", "佩可(夏日)", "🥡", "👙🍙", "泼妇"},
		1076: []string{"可可萝(夏日)", "コッコロ(サマー)", "Kokkoro(Summer)", "水白", "水妈", "水可", "水可可", "水可可萝", "水可可罗", "泳装妈", "泳装可可萝", "泳装可可罗"},
		1077: []string{"铃莓(夏日)", "スズメ(サマー)", "Suzume(Summer)", "水女仆", "水妹抖"},
		1078: []string{"凯留(夏日)", "キャル(サマー)", "Kyaru(Summer)", "凯露(夏日)", "水黑", "水黑猫", "水臭鼬", "泳装黑猫", "泳装臭鼬", "潶", "溴", "💧黑"},
		1079: []string{"珠希(夏日)", "タマキ(サマー)", "Tamaki(Summer)", "水猫剑", "水猫", "泳猫", "渵", "💧🐱🗡️", "水🐱🗡️"},
		1080: []string{"美冬(夏日)", "ミフユ(サマー)", "Mifuyu(Summer)", "水子龙", "水美冬"},
		1081: []string{"忍(万圣节)", "シノブ(ハロウィン)", "Shinobu(Halloween)", "万圣忍", "瓜忍", "🎃忍", "🎃💀"},
		1082: []string{"宫子(万圣节)", "ミヤコ(ハロウィン)", "Miyako(Halloween)", "万圣宫子", "万圣布丁", "狼丁", "狼布丁", "万圣🍮", "🐺🍮", "🎃🍮", "👻🍮"},
		1083: []string{"美咲(万圣节)", "ミサキ(ハロウィン)", "Misaki(Halloween)", "万圣美咲", "万圣大眼", "瓜眼", "🎃眼", "🎃👀", "🎃👁️", "🎃👁"},
		1084: []string{"千歌(圣诞节)", "チカ(クリスマス)", "Chika(Xmas)", "圣诞千歌", "圣千", "蛋鸽", "🎄💰🎶", "🎄千🎶", "🎄1000🎶"},
		1085: []string{"胡桃(圣诞节)", "クルミ(クリスマス)", "Kurumi(Xmas)", "圣诞胡桃", "圣诞铃铛"},
		1086: []string{"绫音(圣诞节)", "アヤネ(クリスマス)", "Ayane(Xmas)", "圣诞熊锤", "蛋锤", "圣锤", "🎄🐻🔨", "🎄🐻", "圣诞熊槌"},
		1087: []string{"日和(新年)", "ヒヨリ(ニューイヤー)", "Hiyori(NewYear)", "新年日和", "春猫", "👘🐱"},
		1088: []string{"优衣(新年)", "ユイ(ニューイヤー)", "Yui(NewYear)", "新年优衣", "春田", "新年由衣"},
		1089: []string{"怜(新年)", "レイ(ニューイヤー)", "Rei(NewYear)", "春剑", "春怜", "春伶", "新春剑圣", "新年怜", "新年剑圣"},
		1090: []string{"惠理子(情人节)", "エリコ(バレンタイン)", "Eriko(Valentine)", "情人节病娇", "恋病", "情病", "恋病娇", "情病娇"},
		1091: []string{"静流(情人节)", "シズル(バレンタイン)", "Shizuru(Valentine)", "情人节静流", "情姐", "情人节姐姐"},
		1092: []string{"安", "アン", "An", "胖安", "55kg"},
		1093: []string{"露", "ルゥ", "Ruu", "逃课女王"},
		1094: []string{"古蕾娅", "グレア", "Gurea", "龙姬", "古雷娅", "古蕾亚", "古雷亚", "古蕾雅", "🐲🐔", "🐉🐔"},
		1095: []string{"空花(大江户)", "クウカ(オーエド)", "Kuuka(Ooedo)", "江户空花", "江户抖m", "江m", "花m", "江花"},
		1096: []string{"妮诺(大江户)", "ニノン(オーエド)", "Ninon(Ooedo)", "江户扇子", "忍扇"},
		1097: []string{"雷姆", "レム", "Remu", "蕾姆"},
		1098: []string{"拉姆", "ラム", "Ramu"},
		1099: []string{"爱蜜莉雅", "エミリア", "Emiria", "艾米莉亚", "emt", "爱蜜莉亚", "EMT"},
		1100: []string{"铃奈(夏日)", "スズナ(サマー)", "Suzuna(Summer)", "瀑击弓", "水爆", "水爆弓", "水暴", "瀑", "水暴弓", "瀑弓", "泳装暴弓", "泳装爆弓"},
		1101: []string{"伊绪(夏日)", "イオ(サマー)", "Io(Summer)", "水魅魔", "水老师", "泳装魅魔", "泳装老师"},
		1102: []string{"美咲(夏日)", "ミサキ(サマー)", "Misaki(Summer)", "水大眼", "泳装大眼"},
		1103: []string{"咲恋(夏日)", "サレン(サマー)", "Saren(Summer)", "水电", "泳装充电宝", "泳装咲恋", "水着咲恋", "水电站", "水电宝", "水充", "👙🔋"},
		1104: []string{"真琴(夏日)", "マコト(サマー)", "Makoto(Summer)", "水狼", "浪", "水🐺", "泳狼", "泳月", "泳月月", "泳朋", "水月", "水月月", "水朋", "👙🐺"},
		1105: []string{"香织(夏日)", "カオリ(サマー)", "Kaori(Summer)", "水狗", "泃", "水🐶", "水🐕", "泳狗"},
		1106: []string{"真步(夏日)", "マホ(サマー)", "泳装真步", "Maho(Summer)", "水狐狸", "水狐", "水壶", "水真步", "水maho", "氵🦊", "水🦊", "💧🦊"},
		1107: []string{"碧(插班生)", "アオイ(編入生)", "Aoi(Hennyuusei)", "生菜", "插班碧", "学院碧"},
		1108: []string{"克萝依", "クロエ", "Kuroe", "华哥", "黑江", "黑江花子", "花子"},
		1109: []string{"琪爱儿", "チエル", "Chieru", "切露", "茄露", "茄噜", "切噜"},
		1110: []string{"优妮", "ユニ", "Yuni", "真行寺由仁", "由仁", "优尼", "u2", "优妮辈先", "辈先", "书记", "uni", "先辈", "仙贝", "油腻", "优妮先辈", "学姐", "18岁黑丝学姐"},
		1111: []string{"镜华(万圣节)", "キョウカ(ハロウィン)", "Kyouka(Halloween)", "万圣镜华", "万圣小仓唯", "万圣xcw", "猫仓唯", "黑猫仓唯", "mcw", "猫唯", "猫仓", "喵唯"},
		1112: []string{"禊(万圣节)", "ミソギ(ハロウィン)", "Misogi(Halloween)", "万圣禊", "万圣炸弹人", "瓜炸弹人", "万圣炸弹", "万圣炸", "瓜炸", "南瓜炸", "🎃💣"},
		1113: []string{"美美(万圣节)", "ミミ(ハロウィン)", "Mimi(Halloween)", "万圣兔", "万圣兔子", "万圣兔兔", "绷带兔", "绷带兔子", "万圣美美", "绷带美美", "万圣🐰", "绷带🐰", "🎃🐰", "万圣🐇", "绷带🐇", "🎃🐇"},
		1114: []string{"露娜", "ルナ", "Runa", "Luna", "露仓唯", "露cw"},
		1115: []string{"克莉丝提娜(圣诞节)", "クリスティーナ(クリスマス)", "Kurisutiina(Xmas)", "Christina(Xmas)", "Cristina(Xmas)", "圣诞克", "圣诞克总", "圣诞女帝", "蛋克", "圣克", "必胜客"},
		1116: []string{"望(圣诞节)", "ノゾミ(クリスマス)", "Nozomi(Xmas)", "圣诞望", "圣诞偶像", "蛋偶像", "蛋望"},
		1117: []string{"伊莉亚(圣诞节)", "イリヤ(クリスマス)", "Iriya(Xmas)", "圣诞伊莉亚", "圣诞伊利亚", "圣诞伊莉雅", "圣诞伊利雅", "圣诞yly", "圣诞吸血鬼", "圣伊", "圣yly"},

		1119: []string{"可可萝(新年)", "コッコロ(ニューイヤー)", "Kokkoro(NewYear)", "春妈", "春可可", "春白", "新年妈", "新可"},
		1120: []string{"凯留(新年)", "キャル(ニューイヤー)", "Kyaru(NewYear)", "凯露(新年)", "春凯留", "春黑猫", "春黑", "春臭鼬", "新年凯留", "新年黑猫", "新年臭鼬", "唯一神"},
		1121: []string{"铃莓(新年)", "スズメ(ニューイヤー)", "Suzume(NewYear)", "春铃莓", "春女仆", "春妹抖", "新年铃莓", "新年女仆", "新年妹抖"},
		1122: []string{"霞(魔法少女)", "カスミ(マジカル)", "魔法少女霞", "魔法侦探", "魔法杜宾犬", "魔法驴", "魔法驴子", "魔驴", "魔法霞", "魔法少驴"},
		1123: []string{"栞(魔法少女)", "シオリ(マジカル)", "Shiori(MagiGirl)", "魔法少女栞", "魔法tp弓", "魔法小栞", "魔法白虎弓", "魔法白虎妹", "魔法白虎", "魔栞", "魔法栞"},
		1124: []string{"卯月(NGs)", "ウヅキ(デレマス)", "Udsuki(DEREM@S)", "卯月", "卵用", "Udsuki(DEREMAS)", "岛村卯月"},
		1125: []string{"凛(偶像大师)", "凛(NGs)", "リン(デレマス)", "Rin(DEREM@S)", "凛", "Rin(DEREMAS)", "涩谷凛", "西部凛"},
		1126: []string{"未央(NGs)", "ミオ(デレマス)", "Mio(DEREM@S)", "未央", "Mio(DEREMAS)", "本田未央"},
		1127: []string{"铃(游侠)", "リン(レンジャー)", "铃(自卫队)", "Rin(Ranger)", "骑兵松鼠", "游侠松鼠", "游骑兵松鼠", "护林员松鼠", "护林松鼠", "游侠🐿️", "武松"},
		1128: []string{"真阳(游侠)", "マヒル(レンジャー)", "真阳(自卫队)", "Mahiru(Ranger)", "骑兵奶牛", "游侠奶牛", "游骑兵奶牛", "护林员奶牛", "护林奶牛", "游侠🐄", "游侠🐮", "牛叉"},
		1129: []string{"璃乃(奇幻)", "リノ(ワンダー)", "Rino(Wonder)", "璃乃(奇境)", "璃乃(仙境)", "爽弓", "爱丽丝弓", "爱弓", "兔弓", "奇境妹弓", "奇幻妹弓", "奇幻璃乃", "仙境妹弓", "白丝妹弓", "爱丽丝妹弓"},
		1130: []string{"步未(奇幻)", "アユミ(ワンダー)", "Ayumi(Wonder)", "步未(奇境)", "步未(仙境)", "路人兔", "兔人妹", "爱丽丝路人", "奇境路人", "奇幻路人", "奇幻步未", "仙境路人"},
		1131: []string{"流夏(夏日)", "ルカ(サマー)", "Ruka(Summer)", "泳装流夏", "水流夏", "泳装刘夏", "水刘夏", "泳装大姐", "泳装大姐头", "水大姐", "水大姐头", "水儿力", "泳装儿力", "水流"},
		1132: []string{"杏奈(夏日)", "アンナ(サマー)", "Anna(Summer)", "泳装中二", "泳装煤气罐", "水中二", "水煤气罐", "冲", "冲二"},
		1133: []string{"七七香(夏日)", "ナナカ(サマー)", "Nanaka(Summer)", "泳装娜娜卡", "泳装77k", "泳装77香", "水娜娜卡", "水77k", "水77香", "水七七香", "泳装七七香"},
		1134: []string{"泳装初音", "初音(夏日)", "ハツネ(サマー)", "Hatsune(Summer)", "水星", "海星", "水hego", "水星法", "泳装星法", "水⭐法", "水睡法", "湦"},
		1135: []string{"美里(夏日)", "泳装美里", "ミサト(サマー)", "Misato(Summer)", "水母", "泳装圣母", "水圣母"},
		1136: []string{"纯(夏日)", "ジュン(サマー)", "Jun(Summer)", "泳装黑骑", "水黑骑", "泳装纯", "水纯", "小次郎"},
		1137: []string{"茜里(天使)", "アカリ(エンジェル)", "Akari(Angel)", "天使妹法", "天使茜里", "丘比特妹法"},
		1138: []string{"依里(天使)", "ヨリ(エンジェル)", "Yori(Angel)", "天使姐法", "天使依里", "丘比特姐法"},
		1139: []string{"纺希(万圣节)", "ツムギ(ハロウィン)", "Tsumugi(Halloween)", "万圣裁缝", "万圣蜘蛛侠", "🎃🕷️", "🎃🕸️", "万裁", "瓜裁", "鬼裁", "鬼才"},
		1140: []string{"怜(万圣节)", "レイ(ハロウィン)", "Rei(Halloween)", "万圣剑圣", "万剑", "瓜剑", "瓜怜", "万圣怜"},
		1141: []string{"茉莉(万圣节)", "マツリ(ハロウィン)", "Matsuri(Halloween)", "万圣跳跳虎", "万圣老虎", "瓜虎", "🎃🐅"},

		//# =================================== #
		1701: []string{"环奈", "カンナ", "Kanna", "桥本环奈", "红二力", "毛大力", "毛小力", "毛六力", "可大萝", "大可萝", "缝合怪"},

		//# =================================== #

		1802: []string{"公主优依", "公优", "优衣(公主)", "ユイ(プリンセス)", "Yui(Princess)", "公主优衣", "公主yui", "公主种田", "公主田", "公主ue", "掉毛优衣", "掉毛yui", "掉毛ue", "掉毛", "飞翼优衣", "飞翼ue", "飞翼", "飞翼高达", "飞田"},

		1804: []string{"贪吃佩可(公主)", "ペコリーヌ(プリンセス)", "Pekoriinu(Princess)", "公主吃", "公主饭", "公主吃货", "公主佩可", "公主饭团", "公主🍙", "命运高达", "高达", "命运公主", "高达公主", "命吃", "春哥高达", "🤖🍙", "🤖"},
		1805: []string{"可可萝(公主)", "コッコロ（プリンセス）", "Kokkoro(Princess)", "公主妈", "月光妈", "蝶妈", "蝴蝶妈", "月光蝶妈", "公主可", "公主可萝", "公主可可萝", "月光可", "月光可萝", "月光可可萝", "蝶可", "蝶可萝", "蝶可可萝"},

		//# =================================== #
		1900: []string{"爱梅斯", "アメス", "Amesu", "菲欧", "フィオ", "Fio"},

		1907: []string{"大古", "タイゴ", "Taigo", "大吾", "鬼道大吾"},
		1908: []string{"花凛", "カリン", "Karin", "绿毛恶魔"},
		1909: []string{"涅比亚", "ネビア", "Nevia", "Nebia"},
		1910: []string{"真崎", "マサキ", "Masaki"},
		1911: []string{"米涅尔β", "ミネルβ", "MineruBeta", "米涅尔", "ミネル", "Mineru"},

		1914: []string{"豪绅", "ゴウシン", "Goushin"},
		1915: []string{"克里吉塔", "クレジック", "Kurejikku"},
		1916: []string{"基洛", "キイロ", "Kiiro"},
		1917: []string{"善", "ゼーン", "Seen"},
		1918: []string{"兰法", "ランファ", "Ranfa"},
		1919: []string{"阿佐尔德", "アンゾールド", "Anzoorudo"},
		1920: []string{"美空", "ミソラ", "Misora"},

		//# =================================== #
		4031: []string{"骷髅", "髑髏", "Dokuro", "骷髅老爹", "老爹"},

		9000: []string{"祐树", "ユウキ", "Yuuki", "骑士", "骑士君"},
	}

	log.SetPrefix("[PCR]")
	log.SetFlags(log.Ldate | log.Lshortfile)
	xlsx, err := excelize.OpenFile("./RANK.xlsx")
	if err != nil {
		log.Println(err)
		return
	}
	rows := xlsx.GetRows(xlsx.GetSheetMap()[1])
	c := 0

	cache := make(map[int]bool)

	var characterList []Character

	for i, row := range rows {
		distance := xlsx.GetCellValue(xlsx.GetSheetMap()[1], "C"+strconv.Itoa(i))
		if isNum(distance) {
			name := strings.Replace(t2s(row[3]), "\n", " ", -1)
			if len(name) == 0 {
				continue
			}
			oc := c
			for _, realName := range strings.Split(name, " ") {
				if len(realName) == 0 {
					continue
				}
				for k, v := range characters {
					_, ok := cache[k]
					if ok {
						continue
					}
					vs := v.([]string)
					for _, s := range vs {
						if strings.Compare(s, realName) == 0 {
							c += 1
							cache[k] = true
							//去除重复的名字
							names := strings.Split(name, " ")
							res1 := ""
							res2 := strings.Join(vs, "")
							for _, n := range names {
								if !strings.Contains(res2, n) {
									res1 += n
								}
							}
							//过滤非中文的昵称
							//只取两个
							res2 = ""
							rc := 0
							var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5].*$")
							for _, n := range vs {
								if hzRegexp.MatchString(n) {
									res2 += n
									rc += 1
								}
								if rc == 2 {
									break
								}
							}
							name = res1 + res2
							_distance, _ := strconv.ParseFloat(distance, 64)
							characterList = append(characterList, Character{
								Id:             int64(k),
								Name:           name,
								RealName:       names[0],
								Distance:       int64(_distance),
								Rank:           row[6],
								Star:           row[7],
								UnionEquipment: row[8],
								Inform:         strings.Replace(row[9], "\n", " ", -1),
							})
							break
						}
					}
					if oc != c {
						break
					}
				}
				if oc != c {
					break
				}
			}
		}
	}
	InitCharacterDataFile(characterList)
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

func InitEquipmentDataInfoFromJsonFile() {
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
	}
}

func SyncEquipmentDataCompareWithDataBaseAndJson() {
	var es []Equipment
	compile := regexp.MustCompile(`\d.*\.png$`)
	Dao.Select("id,title,hot,id,map,url").Order("hot DESC").Find(&es, "enable = 1")
	for _, e := range es {
		for i := 0; i < len(equipments); i++ {
			if equipments[i].Id == e.Id {
				equipments[i].Map = e.Map
				e.Hot = equipments[i].Hot
				_id := compile.FindStringSubmatch(e.Url)[0]
				e.Id, _ = strconv.ParseInt(_id[:len(_id)-4], 10, 64)
				equipments[i].Id = e.Id
				Dao.Model(&e).Update("hot", e.Hot)
			}
		}
	}
	UpdateEquipmentJson()
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
