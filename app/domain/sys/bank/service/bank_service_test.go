package service

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func TestBankService_Create(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	service := NewBankService()
	ctx := xtest2.NewContext()

	list := getList()
	manyCmd := &command.BankCreateManyCommand{}
	manyCmd.Data = list
	manyCmd.CommandId = idutils.NewId()
	err := service.CreateMany(ctx, manyCmd)
	if err != nil {
		t.Error(err)
	}
}

func TestBankService_DeleteAll(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	service := NewBankService()
	ctx := xtest2.NewContext()

	delCmd := &command.BankDeleteCommand{}
	delCmd.CommandId = idutils.NewId()

	err := service.DeleteAll(ctx, delCmd)
	if err != nil {
		t.Error(err)
	}
}

func getList() []*model.Bank {
	list := []*model.Bank{
		{
			Name:     "中国人民银行",
			Keywords: "央行,人行,PBoC,",
		},
		{
			Name:     "国家开发银行",
			Keywords: "国开行,CDB,",
		},
		{
			Name:     "中国进出口银行",
			Keywords: "进出口行,Exim,",
		},
		{
			Name:     "中国农业发展银行",
			Keywords: "农发行,ADBC,",
		},
		{
			Name:     "中国工商银行",
			Keywords: "工行,ICBC,",
		},
		{
			Name:     "中国农业银行",
			Keywords: "农行,ABC,",
		},
		{
			Name:     "中国银行",
			Keywords: "中行,BOC,",
		},
		{
			Name:     "中国建设银行",
			Keywords: "建行,CCB,",
		},
		{
			Name:     "交通银行",
			Keywords: "交行,BCM,",
		},
		{
			Name:     "中国邮政储蓄银行",
			Keywords: "邮储银行,PSBC,邮政",
		},
		{
			Name:     "招商银行",
			Keywords: "招行,CMB,",
		},
		{
			Name:     "兴业银行",
			Keywords: "兴业,CIB,",
		},
		{
			Name:     "中信银行",
			Keywords: "中信,CITIC,",
		},
		{
			Name:     "中国光大银行",
			Keywords: "光大,CEB,",
		},
		{
			Name:     "平安银行",
			Keywords: "平安,PAB,",
		},
		{
			Name:     "上海浦东发展银行",
			Keywords: "浦发银行,SPDB,",
		},
		{
			Name:     "华夏银行",
			Keywords: "华夏,HXB,",
		},
		{
			Name:     "广发银行",
			Keywords: "广发,CGB,",
		},
		{
			Name:     "中国民生银行",
			Keywords: "民生,CMBC,",
		},
		{
			Name:     "北京银行",
			Keywords: "北京银行,",
		},
		{
			Name:     "上海银行",
			Keywords: "上行,BOSC,",
		},
		{
			Name:     "宁波银行",
			Keywords: "宁波银行,",
		},
		{
			Name:     "南京银行",
			Keywords: "南京银行,",
		},
		{
			Name:     "杭州银行",
			Keywords: "杭州,",
		},
		{
			Name:     "重庆银行",
			Keywords: "重庆银行,",
		},
		{
			Name:     "江苏银行",
			Keywords: "江苏银行,",
		},
		{
			Name:     "徽商银行",
			Keywords: "徽商银行,",
		},
		{
			Name:     "青岛银行",
			Keywords: "青岛银行,",
		},
		{
			Name:     "厦门银行",
			Keywords: "厦门银行,",
		},
		{
			Name:     "成都银行",
			Keywords: "成都银行,",
		},
		{
			Name:     "重庆农村商业银行",
			Keywords: "渝农商行,",
		},
		{
			Name:     "北京农村商业银行",
			Keywords: "北京农商行,北京农商",
		},
		{
			Name:     "上海农村商业银行",
			Keywords: "上海农商行,",
		},
		{
			Name:     "广东农村商业银行",
			Keywords: "广东农信,",
		},
		{
			Name:     "江苏农村商业银行",
			Keywords: "江苏农信,",
		},
		{
			Name:     "四川农村商业银行",
			Keywords: "四川农信,",
		},
		{
			Name:     "浙江网商银行",
			Keywords: "网商银行,MYbank,",
		},
		{
			Name:     "深圳微众银行",
			Keywords: "微众银行,WeBank,",
		},
		{
			Name:     "新网银行",
			Keywords: "新网银行,",
		},
		{
			Name:     "苏宁银行",
			Keywords: "苏宁银行,",
		},
		{
			Name:     "汇丰银行",
			Keywords: "汇丰,HSBC,",
		},
		{
			Name:     "渣打银行",
			Keywords: "渣打,Standard Chartered,",
		},
		{
			Name:     "花旗银行",
			Keywords: "花旗,Citibank,",
		},
		{
			Name:     "德意志银行",
			Keywords: "德意志,Deutsche Bank,",
		},
		{
			Name:     "东亚银行",
			Keywords: "东亚,BEA,",
		},
		{
			Name:     "瑞银集团",
			Keywords: "瑞银,UBS,",
		},
		{
			Name:     "日本三菱日联银行",
			Keywords: "三菱日联,MUFG,",
		},
		{
			Name:     "新加坡大华银行",
			Keywords: "大华银行,UOB,",
		},
		{
			Name:     "恒生银行",
			Keywords: "恒生,",
		},
		{
			Name:     "浙商银行",
			Keywords: "浙商,CZBANK,",
		},
		{
			Name:     "渤海银行",
			Keywords: "渤海,",
		},
		{
			Name:     "恒丰银行",
			Keywords: "恒丰,",
		},
		{
			Name:     "长沙银行",
			Keywords: "长沙银行,",
		},
		{
			Name:     "郑州银行",
			Keywords: "郑州银行,",
		},
		{
			Name:     "西安银行",
			Keywords: "西安银行,",
		},
		{
			Name:     "天津银行",
			Keywords: "天津银行,",
		},
		{
			Name:     "河北银行",
			Keywords: "河北银行,",
		},
		{
			Name:     "山西银行",
			Keywords: "山西银行,",
		},
		{
			Name:     "内蒙古银行",
			Keywords: "内蒙古银行,",
		},
		{
			Name:     "辽沈银行",
			Keywords: "辽沈银行,",
		},
		{
			Name:     "吉林银行",
			Keywords: "吉林银行,",
		},
		{
			Name:     "龙江银行",
			Keywords: "龙江银行,",
		},
		{
			Name:     "哈尔滨银行",
			Keywords: "哈行,",
		},
		{
			Name:     "锦州银行",
			Keywords: "锦州银行,",
		},
		{
			Name:     "盛京银行",
			Keywords: "盛京银行,",
		},
		{
			Name:     "大连银行",
			Keywords: "大连银行,",
		},
		{
			Name:     "鞍山银行",
			Keywords: "鞍山银行,",
		},
		{
			Name:     "抚顺银行",
			Keywords: "抚顺银行,",
		},
		{
			Name:     "丹东银行",
			Keywords: "丹东银行,",
		},
		{
			Name:     "本溪银行",
			Keywords: "本溪银行,",
		},
		{
			Name:     "营口银行",
			Keywords: "营口银行,",
		},
		{
			Name:     "阜新银行",
			Keywords: "阜新银行,",
		},
		{
			Name:     "辽阳银行",
			Keywords: "辽阳银行,",
		},
		{
			Name:     "葫芦岛银行",
			Keywords: "葫芦岛银行,",
		},
		{
			Name:     "盘锦银行",
			Keywords: "盘锦银行,",
		},
		{
			Name:     "朝阳银行",
			Keywords: "朝阳银行,",
		},
		{
			Name:     "铁岭银行",
			Keywords: "铁岭银行,",
		},
		{
			Name:     "白山银行",
			Keywords: "白山银行,",
		},
		{
			Name:     "长春发展农村商业银行",
			Keywords: "长春农商行,",
		},
		{
			Name:     "延边农村商业银行",
			Keywords: "延边农商行,",
		},
		{
			Name:     "四平农村商业银行",
			Keywords: "四平农商行,",
		},
		{
			Name:     "通化农村商业银行",
			Keywords: "通化农商行,",
		},
		{
			Name:     "松原农村商业银行",
			Keywords: "松原农商行,",
		},
		{
			Name:     "白城农村商业银行",
			Keywords: "白城农商行,",
		},
		{
			Name:     "黑龙江省农村信用社联合社",
			Keywords: "黑龙江农信,",
		},
		{
			Name:     "齐齐哈尔农村商业银行",
			Keywords: "齐齐哈尔农商行,",
		},
		{
			Name:     "牡丹江农村商业银行",
			Keywords: "牡丹江农商行,",
		},
		{
			Name:     "佳木斯农村商业银行",
			Keywords: "佳木斯农商行,",
		},
		{
			Name:     "大庆农村商业银行",
			Keywords: "大庆农商行,",
		},
		{
			Name:     "鸡西农村商业银行",
			Keywords: "鸡西农商行,",
		},
		{
			Name:     "双鸭山农村商业银行",
			Keywords: "双鸭山农商行,",
		},
		{
			Name:     "伊春农村商业银行",
			Keywords: "伊春农商行,",
		},
		{
			Name:     "七台河农村商业银行",
			Keywords: "七台河农商行,",
		},
		{
			Name:     "鹤岗农村商业银行",
			Keywords: "鹤岗农商行,",
		},
		{
			Name:     "黑河农村商业银行",
			Keywords: "黑河农商行,",
		},
		{
			Name:     "绥化农村商业银行",
			Keywords: "绥化农商行,",
		},
		{
			Name:     "大兴安岭农村商业银行",
			Keywords: "大兴安岭农商行,",
		},
		//
		{
			Name:     "支付宝",
			Keywords: "Alipay,支付宝,",
		},
		{
			Name:     "财付通",
			Keywords: "微信支付,WeChat Pay,Tenpay,",
		},
		{
			Name:     "银联商务",
			Keywords: "UnionPay Commerce,",
		},
		{
			Name:     "蚂蚁集团",
			Keywords: "Ant Group,",
		},
		{
			Name:     "腾讯金融科技",
			Keywords: "Tencent Fintech,",
		},
		{
			Name:     "京东科技",
			Keywords: "京东数科,JD Digits,",
		},
		{
			Name:     "度小满金融",
			Keywords: "百度金融,Du Xiaoman Financial,",
		},
		{
			Name:     "拉卡拉",
			Keywords: "Lakala,",
		},
		{
			Name:     "通联支付",
			Keywords: "Allinpay,",
		},
		{
			Name:     "汇付天下",
			Keywords: "Huifu,",
		},
		{
			Name:     "快钱",
			Keywords: "99Bill,",
		},
		{
			Name:     "易宝支付",
			Keywords: "YeePay,",
		},
		{
			Name:     "苏宁易付宝",
			Keywords: "苏宁支付,Suning Pay,",
		},
		{
			Name:     "美团支付",
			Keywords: "钱袋宝,Meituan Pay,",
		},
		{
			Name:     "翼支付",
			Keywords: "BestPay,",
		},
		{
			Name:     "和包支付",
			Keywords: "CMPay,",
		},
		{
			Name:     "招联消费金融",
			Keywords: "招联,",
		},
		{
			Name:     "马上消费金融",
			Keywords: "马上金融,",
		},
		{
			Name:     "捷信消费金融",
			Keywords: "Home Credit,",
		},
		{
			Name:     "蚂蚁消费金融",
			Keywords: "花呗,借呗,",
		},
		{
			Name:     "陆金所",
			Keywords: "Lufax,",
		},
		{
			Name:     "360数科",
			Keywords: "360金融,360 DigiTech,",
		},
		{
			Name:     "乐信",
			Keywords: "LexinFintech,",
		},
		{
			Name:     "信也科技",
			Keywords: "拍拍贷,FinVolution Group,",
		},
		{
			Name:     "众安在线财产保险",
			Keywords: "众安保险,ZhongAn Online,",
		},
		{
			Name:     "东方财富",
			Keywords: "East Money,",
		},
		{
			Name:     "同花顺",
			Keywords: "iFinD,",
		},
		{
			Name:     "富途控股",
			Keywords: "富途牛牛,Futu,",
		},
		{
			Name:     "老虎证券",
			Keywords: "Tiger Brokers,",
		},
		{
			Name:     "连连数字",
			Keywords: "连连支付,LianLian Pay,",
		},
		{
			Name:     "PingPong",
			Keywords: "PingPong,",
		},
		{
			Name:     "空中云汇",
			Keywords: "Airwallex,",
		},
		{
			Name:     "联动优势",
			Keywords: "UMPAY,",
		},
		{
			Name:     "宝付支付",
			Keywords: "Baofu Pay,",
		},
		{
			Name:     "杉德支付",
			Keywords: "SandPay,",
		},
		{
			Name:     "卡友支付",
			Keywords: "CardYou,",
		},
		{
			Name:     "付临门",
			Keywords: "Fulinfu,",
		},
		{
			Name:     "随行付",
			Keywords: "VBill,",
		},
		{
			Name:     "国通星驿",
			Keywords: "星驿付,",
		},
		{
			Name:     "中付支付",
			Keywords: "Zhongfu,",
		},
		{
			Name:     "嘉联支付",
			Keywords: "Jialian,",
		},
		{
			Name:     "现代金控",
			Keywords: "Modern Financial,",
		},
		{
			Name:     "易生支付",
			Keywords: "Esen,",
		},
		{
			Name:     "海科融通",
			Keywords: "Hicard,",
		},
		{
			Name:     "银盛支付",
			Keywords: "Yinsheng,",
		},
		{
			Name:     "平安付",
			Keywords: "Ping An Pay,",
		},
		{
			Name:     "沃支付",
			Keywords: "联通支付,Unicom Pay,",
		},
		{
			Name:     "网银在线",
			Keywords: "京东支付,JD Pay,",
		},
		{
			Name:     "中邮消费金融",
			Keywords: "中邮消金,",
		},
		{
			Name:     "兴业消费金融",
			Keywords: "兴业消金,",
		},
		{
			Name:     "平安消费金融",
			Keywords: "平安消金,",
		},
		{
			Name:     "海尔消费金融",
			Keywords: "海尔消金,",
		},
		{
			Name:     "小米消费金融",
			Keywords: "小米消金,",
		},
		{
			Name:     "阳光消费金融",
			Keywords: "阳光消金,",
		},
		{
			Name:     "中原消费金融",
			Keywords: "中原消金,",
		},
		{
			Name:     "湖北消费金融",
			Keywords: "湖北消金,",
		},
		{
			Name:     "锦程消费金融",
			Keywords: "锦程消金,",
		},
		{
			Name:     "盛银消费金融",
			Keywords: "盛银消金,",
		},
		{
			Name:     "哈银消费金融",
			Keywords: "哈银消金,",
		},
		{
			Name:     "尚诚消费金融",
			Keywords: "尚诚消金,",
		},
		{
			Name:     "金美信消费金融",
			Keywords: "金美信消金,",
		},
		{
			Name:     "唯品富邦消费金融",
			Keywords: "唯品富邦消金,",
		},
		{
			Name:     "晋商消费金融",
			Keywords: "晋商消金,",
		},
		{
			Name:     "蒙商消费金融",
			Keywords: "蒙商消金,",
		},
		{
			Name:     "幸福消费金融",
			Keywords: "幸福消金,",
		},
		{
			Name:     "首信易支付",
			Keywords: "PayEase,",
		},
		{
			Name:     "钱海",
			Keywords: "Oceanpayment,",
		},
		{
			Name:     "艾贝盈",
			Keywords: "iPayLinks,",
		},
		{
			Name:     "寻汇",
			Keywords: "Sunrate,",
		},
		{
			Name:     "易极付",
			Keywords: "Yiji Pay,",
		},
		{
			Name:     "环迅支付",
			Keywords: "IPS,",
		},
		{
			Name:     "圣亚云",
			Keywords: "Shengpay,",
		},
		{
			Name:     "裕福支付",
			Keywords: "Yufu,",
		},
		{
			Name:     "新生支付",
			Keywords: "Xinsheng,",
		},
		{
			Name:     "中金支付",
			Keywords: "CPCN,",
		},
		{
			Name:     "爱农驿站",
			Keywords: "inong,",
		},
		{
			Name:     "得仕股份",
			Keywords: "DeShi,",
		},
		{
			Name:     "快捷通",
			Keywords: "KJT,",
		},
		{
			Name:     "中铁银通",
			Keywords: "CRB,",
		},
		{
			Name:     "邦付宝",
			Keywords: "Bonpay,",
		},
		{
			Name:     "中汇支付",
			Keywords: "Zhonghui,",
		},
		{
			Name:     "智付电子支付",
			Keywords: "Dinpay,",
		},
		{
			Name:     "北京市政交通一卡通",
			Keywords: "北京一卡通,",
		},
		{
			Name:     "上海公共交通卡",
			Keywords: "上海交通卡,",
		},
		{
			Name:     "深圳通",
			Keywords: "Shenzhen Tong,",
		},
		{
			Name:     "岭南通",
			Keywords: "羊城通,Lingnan Pass,",
		},
		{
			Name:     "武汉城市一卡通",
			Keywords: "武汉通,",
		},
		{
			Name:     "重庆畅通卡",
			Keywords: "渝城通,",
		},
		{
			Name:     "天津城市一卡通",
			Keywords: "天津城市卡,",
		},
		{
			Name:     "苏州市民卡",
			Keywords: "苏州市民卡,",
		},
		{
			Name:     "南京市民卡",
			Keywords: "南京市民卡,",
		},
		{
			Name:     "杭州市民卡",
			Keywords: "杭州市民卡,",
		},
		{
			Name:     "厦门e通卡",
			Keywords: "厦门e通卡,",
		},
		{
			Name:     "青岛琴岛通卡",
			Keywords: "琴岛通,",
		},
		{
			Name:     "大连明珠卡",
			Keywords: "明珠卡,",
		},
		{
			Name:     "哈尔滨城市通",
			Keywords: "哈尔滨城市通,",
		},
		{
			Name:     "西安长安通",
			Keywords: "长安通,",
		},
		{
			Name:     "成都天府通",
			Keywords: "天府通,",
		},
		{
			Name:     "宁波甬城通",
			Keywords: "甬城通,",
		},
		{
			Name:     "福州榕城通",
			Keywords: "榕城通,",
		},
	}
	for i, c := range list {
		c.Id = idutils.NewId()
		c.Order = i
		c.Keywords = c.Keywords + c.Name + ","
	}
	return list
}
