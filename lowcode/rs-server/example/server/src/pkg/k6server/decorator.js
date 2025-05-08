import { server } from "../k6/server";
export var get = function (path, schema) {
    var fn = function (target, key, descriptor) {
        server.get({ path: path, handle: target[key], schema: schema });
    };
    return fn;
};
export var put = function (path, schema) {
    var fn = function (target, key, descriptor) {
        server.put({ path: path, handle: target[key], schema: schema });
    };
    return fn;
};
export var post = function (path, schema) {
    var fn = function (target, key, descriptor) {
        server.post({ path: path, handle: target[key], schema: schema });
    };
    return fn;
};
export var del = function (path, schema) {
    var fn = function (target, key, descriptor) {
        server.delete({ path: path, handle: target[key], schema: schema });
    };
    return fn;
};
