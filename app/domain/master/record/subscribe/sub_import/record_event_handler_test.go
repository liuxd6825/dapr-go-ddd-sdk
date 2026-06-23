package sub_import

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	event2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

var jsonText = `{
	"commandId":"01K1822XWCV6TPKKY5XPQFXKTK",
	"occurredOn":"2012-09-10T19:20:00Z",
	"data":{
		"caseid":"1001",
		"docid":"F3YOfd110hzRCKuT7eZPSHL4k",
		"fileid":"FIcuNxpulmMopISDYhPdHyrAK",
		"filename":"10w.xlsx",
		"isadditems":false,
		"date":"2012-09-10T19:20:00Z",
		"items":[{
			"date":"2012-09-10T19:20:00Z"
		}]
	}
}`

func Test_Unmarshal(t *testing.T) {
	date := time.Now()
	ctx := xtest.NewContext()
	data := &event2.RecordImportMasterEventData{
		Items: []*field.RecordFields{
			{
				Date: &date,
			},
		},
	}
	ev := event2.NewRecordImportMasterEvent(ctx, "test", data)
	jsonByte, err := json.Marshal(ev)
	if err != nil {
		t.Error(err)
		return
	}
	jText := string(jsonByte)
	var e *event2.RecordImportMasterEvent
	err = json.Unmarshal([]byte(jText), &e)
	if err != nil {
		t.Error(err)
		return
	}
	d1 := ev.Data.Items[0].Date
	d2 := e.Data.Items[0].Date
	if !d1.Equal(*d2) {
		t.Error("error ", d1, d2)
	}
}
