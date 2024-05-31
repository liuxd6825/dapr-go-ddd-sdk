var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
// @ts-ignore
import { db } from "../pkg/k6/db";
// @ts-ignore
import { context } from "../pkg/k6/common";
// @ts-ignore
import { newSchema } from "../pkg/k6/schema";
import { get, put, post } from "../pkg/k6server/decorator";
import { schemas as humanSchemas } from "./schema/human.schema";
var HumanService = /** @class */ (function () {
    function HumanService() {
        var _this = this;
        // @ts-ignore
        this.create_handler = function (rctx, cmd) {
            var data = cmd.data;
            rctx.executor().doCommand(function (ctx) {
                return _this.humans.create(ctx, data);
            }).setResponse();
        };
        // @ts-ignore
        this.update_handler = function (rctx, cmd) {
            rctx.executor().doCommand(function (ctx) {
                var err = _this.schema.validate(cmd.data);
                return err;
            }).setResponse();
        };
        // 按id查找数据
        // @ts-ignore
        this.findById_handler = function (rctx) {
            rctx.executor().doQueryOne(function (ctx) {
                return _this.humans.findById(ctx, rctx.getTenantId(), rctx.getId());
            }).setResponse();
        };
        // 分页查询数据
        // @ts-ignore
        this.findPaging_handler = function (rctx) {
            rctx.executor().doQuery(function (ctx) {
                var data = _this.humans.findPaging(ctx, rctx.getFindPaging());
                return { data: data, error: null };
            }).setResponse();
        };
        this.schema = newSchema(humanSchemas.default);
        this.createSchema = newSchema(humanSchemas.createCommand);
        var ctx = context.background();
        this.humans = db.model("human");
        this.humans.table(ctx, this.schema).create(ctx);
    }
    __decorate([
        post("/tenants/{tenantId}/human", humanSchemas.createCommand)
    ], HumanService.prototype, "create_handler", void 0);
    __decorate([
        put("/tenants/{tenantId}/human")
    ], HumanService.prototype, "update_handler", void 0);
    __decorate([
        get("/tenants/{tenantId}/human/{id}")
    ], HumanService.prototype, "findById_handler", void 0);
    __decorate([
        get("/tenants/{tenantId}/human")
    ], HumanService.prototype, "findPaging_handler", void 0);
    return HumanService;
}());
export { HumanService };
