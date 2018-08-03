### 区块链接入
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


### 目录结构介绍
- certs  秘钥相关
- coins  币种相关
- daemon 币种检测等后台进程
- database 数据库相关
- errors 错误
- script 脚本相关
- utils 工具
- wallets 钱包相关
### daemon 流程图
![daemon 流程图](https://github.com/Rennbon/blockchainDemo/raw/master/daemon/tx_daemon_flow_chart.jpg)
