# rainbow-sync
A daemon that synchronizes IRIS hub data for the Rainbow wallet backend


## Run
```bash make all
nohup ./rainbow-sync > debug.log 2>&1 &
```

## Run with docker
You can run application with docker.
### Image
- Build Rainbow-sync Image
```$xslt
docker build -t rainbow-sync .
```

- Run Application
```$xslt
docker run --name rainbow-sync \&
-v /mnt/data/rainbow-sync/logs:/root/go/src/github.com/irisnet/rainbow-sync/logs \&
-e "DB_ADDR=127.0.0.1:27217" -e "DB_USER=user" \&
-e "DB_PASSWD=password" -e "DB_DATABASE=db_name" \&
-e "IRIS_NETWORK=testnet" \&
-e "SER_BC_FULL_NODE=tcp://localhost:26657,..." rainbow-sync
```


## environment params

| param | type | default |description | example |
| :--- | :--- | :--- | :---: | :---: |
| DB_ADDR | string | "" | db addr | 127.0.0.1:27017,127.0.0.2:27017... |
| DB_USER | string | "" | db user | user |
| DB_PASSWD | string | "" |db passwd  | password |
| DB_DATABASE | string | "" |database name  | db_name |
| IRIS_NETWORK | string | "testnet" |irishub name  | testnet or mainnet |
| SER_BC_FULL_NODES | string | tcp://localhost:26657 | iris full node rpc url | tcp://localhost:26657, tcp://127.0.0.2:26657 |
| WORKER_NUM_EXECUTE_TASK | string | 30 | number of threads executing synchronization TX task | 30 |
| WORKER_MAX_SLEEP_TIME | string | 120 | the maximum time (in seconds) that synchronization TX threads are allowed to be out of work | 120 |
| BLOCK_NUM_PER_WORKER_HANDLE | string | 50 | number of blocks per sync TX task | 50 |


