"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HumanFeign = void 0;
var baseUrl = "dapr://duxm-master-cmd-service/api/v1.0/human";
var HumanFeign = /** @class */ (function () {
    function HumanFeign() {
    }
    HumanFeign.prototype.findPaging = function (ctx, params) {
        console.log("HumanFeign.findPaging()");
    };
    HumanFeign.prototype.create = function (ctx, params) {
        console.log("HumanFeign.create()");
    };
    return HumanFeign;
}());
exports.HumanFeign = HumanFeign;
exports.default = new HumanFeign();
