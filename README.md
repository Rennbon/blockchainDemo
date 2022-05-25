### 该库是学习库，完整开源库请转至boxwallet库 [一闪一闪亮晶晶，star一下好心情]
[boxwallet库](https://github.com/Rennbon/boxwallet)


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

- etc

## 维护及联系
- 微信号：`WB343688972` 或者 扫码，请标注`ETH`
- 欢迎 star 加我交流，会拉群
   
   <img src="https://github.com/Rennbon/pyhikvision/blob/master/doc/wechat.png" width="200px" >
