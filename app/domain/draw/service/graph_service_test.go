package service

import (
	_ "embed"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

const tenantId = "test"
const caseId = "1001"

//go:embed test_file/333.drawio
var linkUpdateTargetIsEmptyDrawio string

func Test_Diff3(t *testing.T) {
	drawFile, err := mxgraph.NewDrawioFile(linkUpdateTargetIsEmptyDrawio)
	if err != nil {
		t.Fatal(err)
	}
	element := drawFile.GetElement("kkcGR8n074Od3FpKltgG-5")
	if element == nil {
		t.Fatal("element is nil")
	}

}

//go:embed test_file/diff1/1_add_object.json
var addObjectJson string

func Test_AddObject(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())

		diff := mxgraph.NewFileDiff(addObjectJson)
		if diff == nil {
			return errors.New("diff is nil")
		}

		graphService := NewGraphService()
		saveBatch := graphService.GetSaveBatch("1001", diff)
		t.Log(saveBatch)

		nodeDao := dao.NewGraphDao()
		nodeDao.BatchSave(ctx, saveBatch, "D001")
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

//go:embed test_file/diff1/2_update_label.json
var updateLabelJson string

func Test_UpdateLabel(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())
		diff := mxgraph.NewFileDiff(updateLabelJson)
		if diff == nil {
			t.Fatal("diff is nil")
		}

		graphService := NewGraphService()
		saveBatch := graphService.GetSaveBatch("1001", diff)
		t.Log(saveBatch)
		nodeDao := dao.NewGraphDao()
		nodeDao.BatchSave(ctx, saveBatch, "D001")
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

//go:embed test_file/diff1/3_all.json
var allJson string

func Test_AllJson(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(xtest.NewEnvConfig_Neo4j())

		diff := mxgraph.NewFileDiff(allJson)
		if diff == nil {
			t.Fatal("diff is nil")
		}

		graphService := NewGraphService()
		saveBatch := graphService.GetSaveBatch("1001", diff)
		t.Log(saveBatch)

		nodeDao := dao.NewGraphDao()
		nodeDao.BatchSave(ctx, saveBatch, "D001")
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}
