# Heading Code

* `keep-empty-line`

```go
package metastore

```

# Initialize Meta Value

Meta value manager.

* `builder`: `sqlStmtInitMetaValue`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ``` INTO `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
INSERT IGNORE INTO `x_meta_store` (`meta_key`, `meta_value`, `modify_at`)
VALUES (?, ?, ?);
```

# Fetch Meta Value with Meta Key

Meta value manager.

* `builder`: `sqlStmtFetchMetaValueWithKey`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ``` FROM `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
SELECT `meta_value`, `modify_at`
FROM `x_meta_store`
WHERE (`meta_key` = ?)
```

# Store Meta Value

Meta value manager.

* `builder`: `sqlStmtStoreMetaValue`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ``` INTO `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
INSERT INTO `x_meta_store` (`meta_key`, `meta_value`, `modify_at`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE `meta_value` = ?, `modify_at` = ?
```

# Migrate Legacy Schema Revision Key from Gen-1

Rename schema revision key from Gen-1

* `builder`: `sqlStmtMigrateLegacySchemaRevKeyGen1`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ``` UPDATE `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
UPDATE `x_meta_store`
SET `meta_key` = ? WHERE (`meta_key` = ?);
```

# Lock Meta Row with Fetching Modification Time

Transcation lock. Data row must exists before lock.

* `builder`: `sqlTxStmtFetchMetaModifyTimeWithLock`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ``` FROM `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
SELECT `modify_at`
FROM `x_meta_store`
WHERE (`meta_key` = ?) FOR UPDATE;
```

# Unlock Meta Row with Updating Modification Time

Release transcation lock.

* `builder`: `sqlTxStmtUpdateMetaRowModifyTime`, `metaStoreTableName string`
* `strip-spaces`
* `replace`:
  - ``` UPDATE `(x_meta_store)` ```
  - `$1`
  - ``` metaStoreTableName ```

```sql
UPDATE `x_meta_store`
SET `modify_at` = ? WHERE `meta_key`= ?;
```
