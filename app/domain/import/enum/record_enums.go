package enum

type MasterType string

const (
	MasterTypeRecord   MasterType = "银行流水"
	MasterTypeContract            = "合同"
	MasterTypeProduct             = "产品"
	MasterTypeHuman               = "人员"
	MasterTypeCompany             = "公司"
	MasterTypeAccount             = "账号"
)

func (t MasterType) Name() string {
	return string(t)
}

type SourceType int

const (
	SourceTypeScan   SourceType = 1 // 扫描导入
	SourceTypeDb     SourceType = 2 // 数据库导入
	SourceTypeExcel  SourceType = 3 // excel导入
	SourceTypeCustom SourceType = 4 // 自定义
)

func (e SourceType) String() string {
	res := "UNKNOWN"
	switch e {
	case SourceTypeDb:
		res = "Db"
	case SourceTypeScan:
		res = "Scan"
	case SourceTypeExcel:
		res = "Excel"
	case SourceTypeCustom:
		res = "Custom"
	default:
		res = "UNKNOWN"
	}
	return res
}

func (e SourceType) Title() string {
	res := "UNKNOWN"
	switch e {
	case SourceTypeDb:
		res = "数据库"
	case SourceTypeScan:
		res = "扫描"
	case SourceTypeExcel:
		res = "Excel"
	case SourceTypeCustom:
		res = "自定义"
	default:
		return "UNKNOWN"
	}
	return res
}

type CurrencyType int

const (
	CurrencyCNY CurrencyType = iota // 人民币
	CurrencyFRF                     // 法国法郎
	CurrencyHKD                     // 港元
	CurrencyCHF                     // 瑞士法郎
	CurrencyUSD                     // 美元
	CurrencyCAD                     // 加拿大元
	CurrencyGBP                     // 英镑
	CurrencyNLG                     // 荷兰盾
	CurrencyTHB                     // 泰铢
	CurrencyJPY                     // 日元
	CurrencyEUR                     // 欧元
	CurrencySUR                     // 俄罗斯卢布
	CurrencySGD                     // 新加坡元
	CurrencyKRW                     // 韩国元
	CurrencyDEM                     // 德国马克
)

func (e CurrencyType) String() string {
	res := ""
	switch e {
	case CurrencyCNY:
		res = "CNY"
	case CurrencyFRF:
		res = "FRF"
	case CurrencyHKD:
		res = "HKD"
	case CurrencyCHF:
		res = "CHF"
	case CurrencyUSD:
		res = "USD"
	case CurrencyCAD:
		res = "CAD"
	case CurrencyGBP:
		res = "GBP"
	case CurrencyNLG:
		res = "NLG"
	case CurrencyDEM:
		res = "DEM"
	case CurrencyTHB:
		res = "THB"
	case CurrencyJPY:
		res = "JPY"
	case CurrencyEUR:
		res = "EUR"
	case CurrencySUR:
		res = "SUR"
	case CurrencySGD:
		res = "SGD"
	case CurrencyKRW:
		res = "KRW"
	default:
		res = ""
	}
	return res
}

func (e CurrencyType) Title() string {
	res := "UNKNOWN"
	switch e {
	case CurrencyCNY:
		res = "人民币"
	case CurrencyFRF:
		res = "法国法郎"
	case CurrencyHKD:
		res = "港元"
	case CurrencyCHF:
		res = "瑞士法郎"
	case CurrencyUSD:
		res = "美元"
	case CurrencyCAD:
		res = "加拿大元"
	case CurrencyGBP:
		res = "英镑"
	case CurrencyNLG:
		res = "荷兰盾"
	case CurrencyDEM:
		res = "德国马克"
	case CurrencyTHB:
		res = "泰铢"
	case CurrencyJPY:
		return "日元"
	case CurrencyEUR:
		return "欧元"
	case CurrencySUR:
		return "俄罗斯卢布"
	case CurrencySGD:
		return "新加坡元"
	case CurrencyKRW:
		return "韩国元"
	default:
		return "UNKNOWN"
	}
	return res
}
