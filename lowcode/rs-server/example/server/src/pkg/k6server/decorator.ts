import {server} from "../k6/server";
import {Schema} from "../k6/schema";

export const get=(path:string, schema?:Schema):MethodDecorator=>{
    const fn: MethodDecorator=(target: any, key: string, descriptor: PropertyDescriptor)=>{
        server.get({ path:path, handle:target[key], schema})
    }
    return fn;
}


export const put=(path:string, schema?:Schema):MethodDecorator=>{
    const fn: MethodDecorator=(target: any, key: string, descriptor: PropertyDescriptor)=>{
        server.put({ path:path, handle:target[key], schema})
    }
    return fn;
}


export const post=(path:string, schema?:Schema):MethodDecorator=>{
    const fn: MethodDecorator=(target: any, key: string, descriptor: PropertyDescriptor)=>{
        server.post({path:path, handle:target[key], schema})
    }
    return fn;
}


export const del=(path:string, schema?:Schema):MethodDecorator=>{
    const fn: MethodDecorator=(target: any, key: string, descriptor: PropertyDescriptor)=>{
        server.delete({path:path, handle:target[key], schema})
    }
    return fn;
}
