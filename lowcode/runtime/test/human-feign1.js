"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Human = void 0;
// @ts-ignore

var baseUrl = "dapr://duxm-master-cmd-service/api/v1.0/human";
var Human = /** @class */ (function () {
    function Human() {
    }
    Human.prototype.findPaging = function (ctx, params) {
        var opts = {
            method: "get",
            url: "".concat(baseUrl),
            params: params_1.default.loadFile("./params/", ""),
        };
        return feign.get(ctx, opts, params);
    };
    Human.prototype.create = function (ctx, params) {
        var opts = {
            method: "post",
            url: "".concat(baseUrl),
            params: params_1.default.loadFile("./params/", ""),
        };
        return feign.post(ctx, opts, params);
    };
    return Human;
}());
var feign = new Human();
exports.default = feign;
