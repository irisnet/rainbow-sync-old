# rainbow-sync
A daemon that synchronizes Game of Zone hub data for the Rainbow wallet backend


## Run
```bash
make all
nohup ./rainbow-sync > debug.log 2>&1 &
```

## Run with docker
You can run application with docker.
### Zone
- Build Rainbow-sync Image
```$xslt
docker build -t rainbow-sync .
```

- Run CosmosHub Application
   
```$xslt
docker run --name rainbow-sync-cosmos \&
-v /mnt/data/rainbow-sync/logs:/root/go/src/github.com/irisnet/rainbow-sync/logs \&
-e "DB_ADDR=127.0.0.1:27217" -e "DB_USER=user" \&
-e "DB_PASSWD=password" -e "DB_DATABASE=db_name" \&
-e "ZONE_NAME=cosmos" \&
-e "SER_BC_FULL_NODE_ZONE=tcp://localhost:26657,..." rainbow-sync:latest
```
- Run IrisHub Application
  
```$xslt
docker run --name rainbow-sync-iris \&
-v /mnt/data/rainbow-sync/logs:/root/go/src/github.com/irisnet/rainbow-sync/logs \&
-e "DB_ADDR=127.0.0.1:27217" -e "DB_USER=user" \&
-e "DB_PASSWD=password" -e "DB_DATABASE=db_name" \&
-e "ZONE_NAME=iris" \&
-e "SER_BC_FULL_NODE_ZONE=tcp://localhost:26657,..." rainbow-sync:latest
```


## environment params

| param | type | default |description | example |
| :--- | :--- | :--- | :---: | :---: |
| DB_ADDR | string | "" | db addr | 127.0.0.1:27017,127.0.0.2:27017... |
| DB_USER | string | "" | db user | user |
| DB_PASSWD | string | "" |db passwd  | password |
| DB_DATABASE | string | "" |database name  | db_name |
| ZONE_NAME | string | cosmos |zone name  | cosmos |
| SER_BC_FULL_NODE_ZONE | string | tcp://localhost:36657 |Zone full node rpc url  | tcp://localhost:36657, tcp://127.0.0.2:36657 |
| WORKER_NUM_EXECUTE_TASK_ZONE | string | 30 | 执行同步Zone的Tx任务的线程数 | 30 |
| WORKER_MAX_SLEEP_TIME_ZONE | string | 120 | 允许同步Zone的Tx线程处于不工作状态的最大时长（单位为：秒） | 120 |
| BLOCK_NUM_PER_WORKER_HANDLE_ZONE | string | 50 | 每个同步Zone的Tx任务所包含的Zone区块数 | 50 |


