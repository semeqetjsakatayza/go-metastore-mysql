# Heading Code

* `keep-empty-line`

```go
package metastore

import (
  "database/sql"
  "fmt"
  "context"

  mysqlerrors "github.com/semeqetjsakatayza/go-mysql-errors"
)

func makeMetaStoreRevKey(tableName string) string {
  return metaKeyMetaStoreSchemaRev+":"+tableName
}

```

# MetaStore (meta-store) r.2

* `builder`: `makeSQLCreate`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ```CREATE TABLE `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```
* `replace`:
  - ``` (x_meta_store) storage table ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
CREATE TABLE `x_meta_store` (
  `meta_key` varchar(128) CHARACTER SET ascii NOT NULL COMMENT 'Key of meta information',
  `meta_value` text NOT NULL COMMENT 'Value of meta information',
  `modify_at` bigint(20) NOT NULL DEFAULT '0' COMMENT 'Timestamp of last update',
  PRIMARY KEY (`meta_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='x_meta_store storage table';
```

## Routines

### fetch revision

```go
metaStore := MetaStore {
  TableName: m.referenceTableName,
  Ctx: m.ctx,
  Conn: m.conn,
}
rev, _, err := metaStore.FetchRevision(makeMetaStoreRevKey(m.referenceTableName))
if nil != err {
  if mysqlerrors.IsTableNotExistError(err) {
    err = nil
  }
  return
}
if 0 == rev {
  return
}
revisionRecords = []*schemaRevisionOfMetaStore{&schemaRevisionOfMetaStore{
  currentRev: rev,
	metaStoreTableName: m.referenceTableName,
},}
```

### update revision

```go
metaStore := MetaStore {
  TableName: m.referenceTableName,
  Ctx: m.ctx,
  Conn: m.conn,
}
err = metaStore.StoreRevision(makeMetaStoreRevKey(metaStoreTableName), targetRev)
```

## Migrations

* `strip-spaces`
* `replace`:
  - ```ALTER TABLE `(x_meta_store)`$```
  - `$1`
  - ``` metaStoreTableName ```

### To r.2

```sql
ALTER TABLE `x_meta_store`
CHANGE COLUMN `meta_key` `meta_key` VARCHAR(128) CHARACTER SET 'ascii' NOT NULL
COMMENT 'Key of meta information' ;
```
