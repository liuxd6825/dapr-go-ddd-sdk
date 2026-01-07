# 开发

## 格式化
juicefs format --storage s3 --bucket http://192.168.120.224:9000/atp 
--access-key minioadmin \
--secret-key minioadmin \
redis://default:123456@192.168.120.224:6379/10 \
docfs

## 挂载
juicefs mount \
-d -o allow_other \
--metrics :9567 \
redis://default:123456@192.168.120.224:6379/10 \
/Users/lxd/Projects/duxm/h-master/file-store


 