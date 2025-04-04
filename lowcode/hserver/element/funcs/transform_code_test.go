package funcs

import "testing"

func Test_Transform(t *testing.T) {
	tsCode := `
	// @go-runtime
	let pkg:PKG;

	(function (self: HumanService, ctx: any, pkg: PKG) {
		self.schema = pkg.schema.loadFile("/definition/db/human.json");
		self.dao = pkg.mongo.newDao("default", "human");
		self.dao.table(ctx, self.schema).create(ctx);
		self.feign = require("./human-feign.js").default;
	})()
		`

	code, err := TransformCode(tsCode)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(code))
}
