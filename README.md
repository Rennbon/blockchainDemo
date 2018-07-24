### 区块链demo练习
>demo中所有的privatekey和publickey等重要内容都存进了数据库，主要是为了便于调试
### 如果mongodb不太了解的话，可以直接docker安装
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


