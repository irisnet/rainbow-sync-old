// create database and user
// use rainbow-server
// db.createUser(
//     {
//         user:"iris",
//         pwd:"irispassword",
//         roles:[{role:"root",db:"admin"}]
//     }
// )
// create table
db.createCollection("sync_iris_asset_detail");
db.createCollection("sync_iris_block");
db.createCollection("sync_iris_task");
db.createCollection("sync_iris_tx");


// create index
// db.sync_iris_task.createIndex({"status": 1}, {"background": true});
// db.sync_iris_tx.createIndex({"to": -1, "height": -1});
// db.sync_iris_asset_detail.createIndex({"to": -1, "height": -1});
// db.sync_iris_asset_detail.createIndex({"to": -1, "subject": -1});
// db.sync_iris_block.createIndex({"height": -1}, {"unique": true});
// db.sync_iris_task.createIndex({"start_height": 1, "end_height": 1}, {"unique": true});
//
//
// db.sync_iris_tx.createIndex({'from': 1}, {'background': true});
// db.sync_iris_tx.createIndex({'initiator': 1}, {'background': true});
// db.sync_iris_tx.createIndex({"type": 1}, {"background": true});
/*
 * remove collection data
 */
// db.sync_iris_asset_detail.deleteMany({});
// db.sync_block.deleteMany({});
// db.sync_task.deleteMany({});