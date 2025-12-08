## 格式化
juicefs format --storage s3 --bucket http://192.168.120.224:9000/atp --access-key minioadmin --secret-key minioadmin redis://default:123456@192.168.120.224:6379/10 docfs

## 192.168.120.224 test挂载
juicefs mount -d -o allow_other --metrics :9567 redis://default:123456@192.168.120.224:6379/10 /home/l/program/file-store

## Neo4j Docker 
sudo docker run -d \
--name test-neo4j \
-p 17474:7474 \
-p 17687:7687 \
-v $PWD/neo4j/data:/data \
-v $PWD/neo4j/logs:/logs \
-v $PWD/neo4j/conf:/conf \
-v $PWD/neo4j/plugins:/plugins \
--env NEO4J_AUTH=neo4j/12345678 \
neo4j:latest