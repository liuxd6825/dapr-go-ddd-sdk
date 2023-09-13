package randomutils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"math"
	"math/rand"
	"time"

	crypto "crypto/rand"
	"math/big"
)

var lastName = []string{
	"赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋",
	"沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张", "孔", "曹", "严", "华", "金", "魏",
	"陶", "姜", "戚", "谢", "邹", "喻", "柏", "水", "窦", "章", "云", "苏", "潘", "葛", "奚", "范", "彭",
	"郎", "鲁", "韦", "昌", "马", "苗", "凤", "花", "方", "任", "袁", "柳", "鲍", "史", "唐", "费", "薛",
	"雷", "贺", "倪", "汤", "滕", "殷", "罗", "毕", "郝", "安", "常", "傅", "卞", "齐", "元", "顾", "孟",
	"平", "黄", "穆", "萧", "尹", "姚", "邵", "湛", "汪", "祁", "毛", "狄", "米", "伏", "成", "戴", "谈",
	"宋", "茅", "庞", "熊", "纪", "舒", "屈", "项", "祝", "董", "梁", "杜", "阮", "蓝", "闵", "季", "贾",
	"路", "娄", "江", "童", "颜", "郭", "梅", "盛", "林", "钟", "徐", "邱", "骆", "高", "夏", "蔡", "田",
	"樊", "胡", "凌", "霍", "虞", "万", "支", "柯", "管", "卢", "莫", "柯", "房", "裘", "缪", "解", "应",
	"宗", "丁", "宣", "邓", "单", "杭", "洪", "包", "诸", "左", "石", "崔", "吉", "龚", "程", "嵇", "邢",
	"裴", "陆", "荣", "翁", "荀", "于", "惠", "甄", "曲", "封", "储", "仲", "伊", "宁", "仇", "甘", "武",
	"符", "刘", "景", "詹", "龙", "叶", "幸", "司", "黎", "溥", "印", "怀", "蒲", "邰", "从", "索", "赖",
	"卓", "屠", "池", "乔", "胥", "闻", "莘", "党", "翟", "谭", "贡", "劳", "逄", "姬", "申", "扶", "堵",
	"冉", "宰", "雍", "桑", "寿", "通", "燕", "浦", "尚", "农", "温", "别", "庄", "晏", "柴", "瞿", "阎",
	"连", "习", "容", "向", "古", "易", "廖", "庾", "终", "步", "都", "耿", "满", "弘", "匡", "国", "文",
	"寇", "广", "禄", "阙", "东", "欧", "利", "师", "巩", "聂", "关", "荆", "司马", "上官", "欧阳", "夏侯",
	"诸葛", "闻人", "东方", "赫连", "皇甫", "尉迟", "公羊", "澹台", "公冶", "宗政", "濮阳", "淳于", "单于",
	"太叔", "申屠", "公孙", "仲孙", "轩辕", "令狐", "徐离", "宇文", "长孙", "慕容", "司徒", "司空"}
var firstName = []string{
	"伟", "刚", "勇", "毅", "俊", "峰", "强", "军", "平", "保", "东", "文", "辉", "力", "明", "永", "健", "世", "广", "志", "义",
	"兴", "良", "海", "山", "仁", "波", "宁", "贵", "福", "生", "龙", "元", "全", "国", "胜", "学", "祥", "才", "发", "武", "新",
	"利", "清", "飞", "彬", "富", "顺", "信", "子", "杰", "涛", "昌", "成", "康", "星", "光", "天", "达", "安", "岩", "中", "茂",
	"进", "林", "有", "坚", "和", "彪", "博", "诚", "先", "敬", "震", "振", "壮", "会", "思", "群", "豪", "心", "邦", "承", "乐",
	"绍", "功", "松", "善", "厚", "庆", "磊", "民", "友", "裕", "河", "哲", "江", "超", "浩", "亮", "政", "谦", "亨", "奇", "固",
	"之", "轮", "翰", "朗", "伯", "宏", "言", "若", "鸣", "朋", "斌", "梁", "栋", "维", "启", "克", "伦", "翔", "旭", "鹏", "泽",
	"晨", "辰", "士", "以", "建", "家", "致", "树", "炎", "德", "行", "时", "泰", "盛", "雄", "琛", "钧", "冠", "策", "腾", "楠",
	"榕", "风", "航", "弘", "秀", "娟", "英", "华", "慧", "巧", "美", "娜", "静", "淑", "惠", "珠", "翠", "雅", "芝", "玉", "萍",
	"红", "娥", "玲", "芬", "芳", "燕", "彩", "春", "菊", "兰", "凤", "洁", "梅", "琳", "素", "云", "莲", "真", "环", "雪", "荣",
	"爱", "妹", "霞", "香", "月", "莺", "媛", "艳", "瑞", "凡", "佳", "嘉", "琼", "勤", "珍", "贞", "莉", "桂", "娣", "叶", "璧",
	"璐", "娅", "琦", "晶", "妍", "茜", "秋", "珊", "莎", "锦", "黛", "青", "倩", "婷", "姣", "婉", "娴", "瑾", "颖", "露", "瑶",
	"怡", "婵", "雁", "蓓", "纨", "仪", "荷", "丹", "蓉", "眉", "君", "琴", "蕊", "薇", "菁", "梦", "岚", "苑", "婕", "馨", "瑗",
	"琰", "韵", "融", "园", "艺", "咏", "卿", "聪", "澜", "纯", "毓", "悦", "昭", "冰", "爽", "琬", "茗", "羽", "希", "欣", "飘",
	"育", "滢", "馥", "筠", "柔", "竹", "霭", "凝", "晓", "欢", "霄", "枫", "芸", "菲", "寒", "伊", "亚", "宜", "可", "姬", "舒",
	"影", "荔", "枝", "丽", "阳", "妮", "宝", "贝", "初", "程", "梵", "罡", "恒", "鸿", "桦", "骅", "剑", "娇", "纪", "宽", "苛",
	"灵", "玛", "媚", "琪", "晴", "容", "睿", "烁", "堂", "唯", "威", "韦", "雯", "苇", "萱", "阅", "彦", "宇", "雨", "洋", "忠",
	"宗", "曼", "紫", "逸", "贤", "蝶", "菡", "绿", "蓝", "儿", "翠", "烟", "小", "轩"}
var lastNameLen = len(lastName)
var firstNameLen = len(firstName)

func StringNumber(size int) string {
	rand.Seed(time.Now().UnixNano())
	return randomString(size, 0)
}

func StringLower(size int) string {
	rand.Seed(time.Now().UnixNano())
	return randomString(size, 1)
}

func StringUpper(size int) string {
	rand.Seed(time.Now().UnixNano())
	return randomString(size, 2)
}

func String(size int) string {
	rand.Seed(time.Now().UnixNano())
	return randomString(size, 2)
}

func Int() int {
	v := RangeRand(0, 1000)
	return int(v)
}

func IntMax(max int) int {
	v := RangeRand(0, int64(max-1))
	return int(v)
}

func Int64() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(100000)
}

func Int64Max(max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max)
}

func Uint32() uint32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

func Uint64() uint64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint64()
}

func Float32() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32()
}

func PFloat32() *float32 {
	rand.Seed(time.Now().UnixNano())
	v := Float32()
	return &v
}

func Float64() float64 {
	rand.Seed(time.Now().Unix())
	v := float64(RangeRand(0, 1000)) + float64(RangeRand(0, 100))/100
	return v
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := crypto.Int(crypto.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	}
	result, _ := crypto.Int(crypto.Reader, big.NewInt(max-min+1))
	return min + result.Int64()
}

func PFloat64() *float64 {
	rand.Seed(time.Now().UnixNano())
	v := Float64()
	return &v
}

// Email 随机生成英文字符串
func Email() string {
	rand.Seed(time.Now().UnixNano())
	name := StringLower(8)
	return fmt.Sprintf("%v@163.com", name)
}

// UUID 随机生成ID
func UUID() string {
	return idutils.NewId()
}

// PUUID 随机生成ID
func PUUID() *string {
	v := UUID()
	return &v
}

// NewId 新建ID
func NewId() string {
	return UUID()
}

// IpAddr 随机生成IP地址
func IpAddr() string {
	rand.Seed(time.Now().UnixNano())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

// PIpAddr 随机生成IP地址
func PIpAddr() *string {
	v := IpAddr()
	return &v
}

// NameCN 随机生成中文名称
func NameCN() string {
	rand.Seed(time.Now().UnixNano())     //设置随机数种子
	var first string                     //名
	for i := 0; i <= rand.Intn(1); i++ { //随机产生2位或者3位的名
		first = fmt.Sprint(firstName[rand.Intn(firstNameLen-1)])
	}
	//返回姓名
	return fmt.Sprintf("%s%s", fmt.Sprint(lastName[rand.Intn(lastNameLen-1)]), first)
}

// PNameCN 随机生成ID
func PNameCN() *string {
	v := NameCN()
	return &v
}

// NameEN 随机生成ID
func NameEN() string {
	return StringLower(6)
}

// PNameEN 随机生成ID
func PNameEN() *string {
	v := NameEN()
	return &v
}

// DateString 随机生成日期字符串
func DateString() string {
	rand.Seed(time.Now().UnixNano())
	max := 2099
	min := 2000
	year := rand.Intn(max-min) + min
	month := rand.Intn(12) + 1
	daysInMonth := 31

	switch month {
	case 2:
		if year%100 != 0 && year%4 == 0 || year%400 == 0 {
			daysInMonth = 29
		} else {
			daysInMonth = 28
		}
	case 4, 6, 9, 11:
		daysInMonth = 30
	}

	day := rand.Intn(daysInMonth) + 1
	return fmt.Sprintf("%v-%v-%v", year, month, day)
}

// PDateString 随机生成日期字符串
func PDateString() *string {
	v := DateString()
	return &v
}

// TimeString 随机生成时间字符串
func TimeString() string {
	return Time().Format("2006-01-02 15:04:05")
}

// PTimeString 随机生成时间字符串
func PTimeString() *string {
	value := TimeString()
	return &value
}

func Date() time.Time {
	value := Time()
	date := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(date, 0)
}

func PDate() *time.Time {
	value := Date()
	return &value
}

func Time() time.Time {
	min := time.Date(2000, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2050, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func PTime() *time.Time {
	value := Time()
	return &value
}

func Now() time.Time {
	return time.Now()
}

func PNow() *time.Time {
	value := Now()
	return &value
}

func Boolean() bool {
	return Uint32()%2 == 1
}

func PBoolean() *bool {
	p := Uint32()%2 == 1
	return &p
}

func JsonTime() types.JSONTime {
	return types.JSONTime(Time())
}

func PJsonTime() *types.JSONTime {
	value := JsonTime()
	return &value
}

func JsonDate() types.JSONDate {
	return types.JSONDate(Date())
}

func PJsonDate() *types.JSONDate {
	value := JsonDate()
	return &value
}

//
//  randomString
//  @Description: 随机生成字符串
//  @param size size 随机码的位数
//  @param kind 0:纯数字/1:小写字母/2:大写字母/3:数字、大小写字母
//  @return string
//
func randomString(size int, kind int) string {
	ikind, kinds, rsbytes := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		rsbytes[i] = uint8(base + rand.Intn(scope))
	}
	return string(rsbytes)
}

// Options 参数
/*type Options struct {
	Length  *int
	MinYear *int
	MaxYear *int
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Merge(opts ...*Options) *Options {
	for _, i := range opts {
		if i == nil {
			continue
		}
		if i.Length != nil {
			o.Length = i.Length
		}
		if i.MaxYear != nil {
			o.MaxYear = i.MaxYear
		}
		if i.MinYear != nil {
			o.MinYear = i.MinYear
		}
	}
	return o
}

func (o *Options) GetLength() int {
	if o.Length != nil {
		return *o.Length
	}
	return 20
}

func (o *Options) SetLength(len int) *Options {
	i := len
	o.Length = &i
	return o
}

func (o *Options) GetMaxYear() int {
	if o.Length != nil {
		return *o.MaxYear
	}
	return 3000
}

func (o *Options) SetMaxYear(len int) *Options {
	v := len
	o.MaxYear = &v
	return o
}

func (o *Options) GetMinYear() int {
	if o.Length != nil {
		return *o.MinYear
	}
	return 1900
}

func (o *Options) SetMinYear(len int) *Options {
	v := len
	o.MinYear = &v
	return o
}
*/
