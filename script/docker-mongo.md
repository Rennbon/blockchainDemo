### mongodb配置
```
//docker启动容器
docker run --name bcmongo -p 27017:27017  -d mongo:latest 
//获取containerId
docker ps | grep bcmongo | cut -c 1-12
//进入容器Mongo
docker exec -it 'containerId' mongo
//show dbs
use blockChain
db.createCollection("account")
db.createCollection("tx")
```
