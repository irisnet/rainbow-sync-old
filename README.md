# rainbow-sync
A daemon that synchronizes IRIS hub data for the Rainbow wallet backend

# Structure

- `conf`: config of project
- `block`: parse asset detail and tx function module
- `model`: mongodb script to create database
- `task`: main logic of sync-server, sync data from blockChain and write to database
- `db`: database model
- `msgs`: tx msgs model
- `lib`: cdc and client pool functions
- `utils`: common functions
- `main.go`: bootstrap project

# SetUp
## Database
Use Mongodb  to store IRIS hub data

# Build And Run

- Build: `make all`
- Run: `make run`
- Cross compilation: `make build-linux`

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
-e "DB_PASSWD=password" -e "DB_DATABASE=db_name"  \&
-e "SER_BC_FULL_NODES=tcp://localhost:26657,..." rainbow-sync
```


## Environment Params

| param | type | default |description | example |
| :--- | :--- | :--- | :---: | :---: |
| DB_ADDR | string | "" | db addr | 127.0.0.1:27017,127.0.0.2:27017... |
| DB_USER | string | "" | db user | user |
| DB_PASSWD | string | "" |db passwd  | password |
| DB_DATABASE | string | "" |database name  | db_name |
| SER_BC_FULL_NODES | string | tcp://localhost:26657 | iris full node rpc url | tcp://localhost:26657, tcp://127.0.0.2:26657 |
| WORKER_NUM_EXECUTE_TASK | string | 30 | number of threads executing synchronization TX task | 30 |
| WORKER_MAX_SLEEP_TIME | string | 120 | the maximum time (in seconds) that synchronization TX threads are allowed to be out of work | 120 |
| BLOCK_NUM_PER_WORKER_HANDLE | string | 50 | number of blocks per sync TX task | 50 |
| BEHIND_BLOCK_NUM | string | 0 | wait block num to handle tx | 0 |

- Remarks
  - synchronizes  block chain data from  specify block height(such as:17908 current time:1576208532)
  
     At first,stop the rainbow-sync and create the task. Run:
  ```bash
     ﻿﻿db.sync_iris_task.insert({'start_height':NumberLong(17908),'end_height':NumberLong(0),'current_height':NumberLong(0),'status':'unhandled','last_update_time':NumberLong(1576208532)})
  ```
  Then,start the rainbow-sync.