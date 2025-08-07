package restapi

import (
	"testing"
)

var params = "\t{\"_id\":\"6880d4920110ab34e1224890\",\"app_id\":\"import2\",\"created_time\":\"2025-07-23T20:24:50.047+08:00\",\"data\":{\"command\":{\"commandid\":\"0002\",\"data\":{\"caseid\":\"1001\",\"docid\":\"F3YOfd110hzRCKuT7eZPSHL4k\",\"fileid\":\"FIcuNxpulmMopISDYhPdHyrAK\",\"filename\":\"10w.xlsx\",\"isadditems\":true,\"items\":[{\"acct\":\"1027600000094762\",\"acct_type\":\"\",\"amount\":6600,\"balance\":957769213.91,\"bank_name\":\"招商银行\",\"caseid\":\"\",\"category\":\"\",\"ccy\":\"CNY\",\"date\":\"2021-08-31T00:00:00Z\",\"docid\":\"\",\"fileid\":\"\",\"filename\":\"\",\"id\":\"511a3f45-81d9-4363-9027-9e6dae373f4e\",\"iden\":\"\",\"income\":6600,\"meta\":null,\"name\":\"北京玖富普惠信息技术有限公司\",\"notes\":\"\",\"opp_acct\":\"1027600000094761\",\"opp_acct_type\":\"\",\"opp_bank_name\":\"北京分行知春支行\",\"opp_category\":\"\",\"opp_iden\":\"\",\"opp_name\":\"北京玖富普惠信息技术有限公司\",\"payout\":-6600,\"place\":\"\",\"row_num\":0,\"serial\":\"\",\"sheetname\":\"\",\"summary\":\"HOST202108314947407\",\"task_id\":\"\",\"type\":\"\"}],\"sheetname\":\"\",\"taskid\":\"taskId\"},\"isvalidonly\":false,\"updatemask\":null}},\"id\":\"TJBLYSC8DZZ5TTCBUAX3TPL3WLZPUS\",\"meta\":null,\"tenant_id\":\"test\",\"topic\":\"record-import-master-event\"}\n"

func Test_GetEventParams(t *testing.T) {

}
