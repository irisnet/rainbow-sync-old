# rainbow-sync
A daemon that synchronizes IRIS hub data for the Rainbow wallet backend


## Run
- Cosmos
```bash
cd service/cosmos && make all
nohup ./rainbow-sync > debug.log 2>&1 &
```
- Iris
```bash
cd service/iris && make all
nohup ./rainbow-sync > debug.log 2>&1 &
```

## Run with docker
You can run application with docker.
### Iris
- Build Rainbow-sync Image
```$xslt
cd service/iris && docker build -t rainbow-sync:dev01 .
```

- Run Application
```$xslt
docker run --name rainbow-sync \&
-v /mnt/data/rainbow-sync/logs:/root/go/src/github.com/irisnet/rainbow-sync/logs \&
-e "DB_ADDR=127.0.0.1:27217" -e "DB_USER=user" \&
-e "DB_PASSWD=password" -e "DB_DATABASE=db_name" \&
-e "IRIS_NETWORK=testnet" \&
-e "SER_BC_FULL_NODE=tcp://localhost:26657,..." rainbow-sync:dev01
```
### Cosmos
- Build Rainbow-sync Image
```$xslt
cd service/cosmos && docker build -t rainbow-sync:dev01 .
```
- Run Application
```$xslt
docker run --name rainbow-sync \&
-v /mnt/data/rainbow-sync/logs:/root/go/src/github.com/irisnet/rainbow-sync/logs \&
-e "DB_ADDR=127.0.0.1:27217" -e "DB_USER=user" \&
-e "DB_PASSWD=password" -e "DB_DATABASE=db_name" \&
-e "SER_BC_FULL_NODE_COSMOS=tcp://localhost:36657,..."  rainbow-sync:dev01
```


## environment params

| param | type | default |description | example |
| :--- | :--- | :--- | :---: | :---: |
| DB_ADDR | string | "" | db addr | 127.0.0.1:27017,127.0.0.2:27017... |
| DB_USER | string | "" | db user | user |
| DB_PASSWD | string | "" |db passwd  | password |
| DB_DATABASE | string | "" |database name  | db_name |
| IRIS_NETWORK | string | "testnet" |irishub name  | testnet or mainnet |
| SER_BC_FULL_NODE | string | tcp://localhost:26657 | iris full node rpc url | tcp://localhost:26657, tcp://127.0.0.2:26657 |
| WORKER_NUM_CREATE_TASK | string | 2 | 创建同步Iris的Tag任务的线程数 | 2 |
| WORKER_NUM_EXECUTE_TASK | string | 30 | 执行同步Iris的Tag任务的线程数 | 30 |
| WORKER_MAX_SLEEP_TIME | string | 120 | 允许同步Iris的Tag线程处于不工作状态的最大时长（单位为：秒） | 120 |
| BLOCK_NUM_PER_WORKER_HANDLE | string | 50 | 每个同步Iris的Tag任务所包含的Iris区块数 | 50 |
| SER_BC_FULL_NODE_COSMOS | string | tcp://localhost:36657 |cosmos full node rpc url  | tcp://localhost:36657, tcp://127.0.0.2:36657 |
| WORKER_NUM_CREATE_TASK_COSMOS | string | 2 | 创建同步Cosmos的Tx任务的线程数 | 2 |
| WORKER_NUM_EXECUTE_TASK_COSMOS | string | 30 | 执行同步Cosmos的Tx任务的线程数 | 30 |
| WORKER_MAX_SLEEP_TIME_COSMOS | string | 120 | 允许同步Cosmos的Tx线程处于不工作状态的最大时长（单位为：秒） | 120 |
| BLOCK_NUM_PER_WORKER_HANDLE_COSMOS | string | 50 | 每个同步Cosmos的Tx任务所包含的Cosmos区块数 | 50 |


