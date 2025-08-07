package sub_import

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"testing"
	"time"
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
	ev := event.RecordImportMasterEvent{
		Event: events.Event[event.RecordImportMasterEventData]{
			EventId: "111",
			Data: event.RecordImportMasterEventData{
				Items: []*field.RecordFields{
					{
						Date: &date,
					},
				},
			},
		},
	}

	jsonByte, err := json.Marshal(ev)
	if err != nil {
		t.Error(err)
		return
	}
	jText := string(jsonByte)
	var e *event.RecordImportMasterEvent
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
