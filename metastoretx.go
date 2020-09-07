package metastore

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

// MetaTx handles operations of meta informations in transcation.
type MetaTx struct {
	TableName string

	Ctx context.Context
	Tx  *sql.Tx
}

// OptionalLock lock meta data row with given metaKey.
// The lock is placed at modifyAt field. Data row may not exist.
func (n *MetaTx) OptionalLock(metaKey string) (ok bool, modifyAt int64, err error) {
	ok = false
	if err = n.Tx.QueryRowContext(n.Ctx, sqlTxStmtFetchMetaModifyTimeWithLock(n.TableName), metaKey).Scan(&modifyAt); nil != err {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return
	}
	return true, modifyAt, nil
}

// RequiredLock lock meta data row with given metaKey.
// The lock is placed at modifyAt field. Data row must exists.
func (n *MetaTx) RequiredLock(metaKey string) (modifyAt int64, err error) {
	err = n.Tx.QueryRowContext(n.Ctx, sqlTxStmtFetchMetaModifyTimeWithLock(n.TableName), metaKey).Scan(&modifyAt)
	return
}

// Unlock update modification time of dat row with given metaKey.
// The transcation lock must be release with Tx.Commit() method.
func (n *MetaTx) Unlock(metaKey string) (err error) {
	modifyAt := time.Now().Unix()
	_, err = n.Tx.ExecContext(n.Ctx, sqlTxStmtUpdateMetaRowModifyTime(n.TableName), modifyAt, metaKey)
	return
}

func (n *MetaTx) store(metaKey, metaValue string) (err error) {
	modifyAt := time.Now().Unix()
	_, err = n.Tx.ExecContext(
		n.Ctx,
		sqlStmtStoreMetaValue(n.TableName),
		metaKey, metaValue, modifyAt, metaValue, modifyAt)
	return
}

// StoreInt64 put int64 value into store.
func (n *MetaTx) StoreInt64(metaKey string, value int64) (err error) {
	textValue := strconv.FormatInt(value, 10)
	return n.store(metaKey, textValue)
}
