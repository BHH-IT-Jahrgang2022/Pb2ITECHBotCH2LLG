# How to MongoDB

## Prequisits

- Have Mongo installed

## Insert Data

- Access Mongo Shell

```bash
mongo
```

- Select Database

```bash
use Chatty
```

- Insert Data to "categorie"

```bash
db._categorie_.insert({"_fieldname_": "_fieldvalue_", "_fieldname_": "_fieldvalue_"})
```

```bash
db.windowfly.insert({"keywords": ["Windowfly", "ab*st(u|ü|ue)rz], "response": "Bitte wärmen"})
```

- See data per categorie

```bash
db._categorie_.find()
```


- Alter a specific entry/ document

```bash
db._categorie_.find()

{"_id": ObjectID("abc123def"), "_field_": "_value_"}

db._categorie_.updateOne({_id: ObjectId("abc123def")}, {$set: {_field_: "new value"}})
```

## Setup Users & remote Access

- Create global admin user

```bash
use admin

db.createUser(
  {
    user: "admin",
    pwd: "password",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase" ]
  }
)
```

- Create additional API user

```bash
use myDB

db.createUser(
  {
    user: "apiUser",
    pwd: "apiPassword",
    roles: [ { role: "readWrite", db: "myDB" } ]
  }
)
```

- edit config

```/etc/mongod.conf```

```bash
net:
  port: 27017
  bindIp: 0.0.0.0  // or a better choice

...

security:
  authorization: "enabled"

```

- Finish

restart MongoDB and login again:

```bash
mongo -u user -p password --authenticationDatabase myDB
```

Note: myDB for the admin user is _admin_

