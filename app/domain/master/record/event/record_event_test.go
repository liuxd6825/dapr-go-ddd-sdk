package event

import (
	"encoding/json"
	"testing"
	"time"
)

var jsonText = `
{
  "_id": "6887061d1e87e987d040de34",
  "id": "C29F7MTVU7JXZKI70DVMTZLV",
  "app_id": "import",
  "created_time": "2025-07-28T13:09:49.325+08:00",
  "meta": null,
  "tenant_id": "test",
  "topic": "record-import-master-event",
  "data": {
    "commandId": "0002",
    "isValidOnly": false,
    "updateMask": null,
	"data": {
    "caseId": "1001",
    "docId": "F3YOfd110hzRCKuT7eZPSHL4k",
    "fileId": "FIcuNxpulmMopISDYhPdHyrAK",
    "fileName": "10w.xlsx",
    "isAddItems": true,
    "sheetName": "",
    "taskId": "taskId",
    "items": [
      {
        "acct": "1027600000094761",
        "acct_type": "",
        "amount": 41.59,
        "balance": 42178190.89,
        "bank_name": "招商银行",
        "caseid": "",
        "category": "",
        "ccy": "CNY",
        "date": "2025-07-28T13:09:49.325+08:00",
        "docid": "",
        "fileid": "",
        "filename": "",
        "id": "c7d282b0-95e6-41bf-bce1-225e78006ec9",
        "iden": "",
        "income": 41.59,
        "meta": null,
        "name": "北京玖富普惠信息技术有限公司",
        "notes": "",
        "opp_acct": "7029600130010717",
        "opp_acct_type": "",
        "opp_bank_name": "北京分行清算中心",
        "opp_category": "",
        "opp_iden": "",
        "opp_name": "银行卡跨行清算-银联代付",
        "payout": -41.59,
        "place": "",
        "row_num": 0,
        "serial": "",
        "sheetname": "",
        "summary": "订单20210822P2P8022C1FB2",
        "task_id": "",
        "type": ""
      }
    ]
  }
}
}

`

type Event struct {
	PID         string          `json:"_id"`
	ID          string          `json:"id"`
	AppID       string          `json:"app_id"`
	CreatedTime time.Time       `json:"created_time"`
	Data        json.RawMessage `json:"data"` // 使用 json.RawMessage 来接收原始的 JSON 数据
	Meta        json.RawMessage `json:"meta"`
	TenantID    string          `json:"tenant_id"`
	Topic       string          `json:"topic"`
}

func Test_JsonUnmarshal(t *testing.T) {
	var event Event
	//err := jsonutils.CustomJson.Unmarshal([]byte(jsonText), &event)
	err := json.Unmarshal([]byte(jsonText), &event)
	if err != nil {
		t.Error(err)
		return
	} else {
		t.Log(event)
	}

	var importEvent RecordImportMasterEvent
	err = json.Unmarshal(event.Data, &importEvent)
	if err != nil {
		t.Error(err)
		return
	} else {
		t.Log(event)
	}
}
