package service

import (
	_ "embed"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/pkg/mxgraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

const tenantId = "test"
const caseId = "1001"

//go:embed test_file/333.drawio
var linkUpdateTargetIsEmptyDrawio string

//go:embed test_file/diff1/1_add_object.json
var addObjectJson string

//go:embed test_file/diff1/2_update_label.json
var updateLabelJson string

//go:embed test_file/diff1/3_all.json
var allJson string

var drawId = "D001"
var neo4jDBKey = ""

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

func Test_AddObject(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(newEnv())
		diff := mxgraph.NewFileDiff(addObjectJson)
		if diff == nil {
			return errors.New("diff is nil")
		}

		graphService := NewGraphService()
		saveBatch := graphService.GetSaveBatch(newSaveFileCommand("", diff))
		t.Log(saveBatch)

		nodeDao := dao.NewGraphDao(neo4jDBKey)
		nodeDao.BatchSave(ctx, saveBatch, "D001")
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_UpdateLabel(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(newEnv())
		diff := mxgraph.NewFileDiff(updateLabelJson)
		if diff == nil {
			t.Fatal("diff is nil")
		}

		graphService := NewGraphService()
		saveBatch := graphService.GetSaveBatch(newSaveFileCommand("", diff))
		t.Log(saveBatch)

		nodeDao := dao.NewGraphDao("")
		nodeDao.BatchSave(ctx, saveBatch, drawId)

		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func Test_AllJson(t *testing.T) {
	gp.Try(func() error {
		ctx := xtest.NewContext()
		env.SetEnv(newEnv())

		diff := mxgraph.NewFileDiff(allJson)
		if diff == nil {
			t.Fatal("diff is nil")
		}

		graphService := NewGraphService()
		saveBatch := graphService.GetSaveBatch(newSaveFileCommand("", diff))
		t.Log(saveBatch)

		nodeDao := dao.NewGraphDao("")
		nodeDao.BatchSave(ctx, saveBatch, drawId)
		return nil
	}).Catch(func(e error) {
		t.Error(e)
	})
}

func newEnv() *env.Env {
	envCfg := xtest.NewEnvConfigNeo4j()
	env.SetEnv(envCfg)
	return envCfg
}

func newSaveFileCommand(xml string, diff *mxgraph.FileDiff) command.IDrawSaveCommand {
	cmd := &command.DrawSaveCommand{}
	cmd.CommandId = idutils.NewId()
	cmd.Data = command.DrawSaveCommandData{
		CaseId:   caseId,
		DrawId:   drawId,
		FileName: "testFile",
		XML:      xml,
		Diff:     diff,
	}
	return cmd
}
