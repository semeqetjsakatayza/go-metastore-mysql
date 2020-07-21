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
