package metastore

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	mysqlerrors "github.com/semeqetjsakatayza/go-mysql-errors"
)

// MetaStore handles operations of meta informations.
type MetaStore struct {
	TableName string

	Ctx  context.Context
	Conn *sql.DB
}

// PrepareSchema perform revision check and upgrade of meta store table schema.
func (m *MetaStore) PrepareSchema() (schemaChanged bool, err error) {
	schemaMgmt := schemaManager{
		referenceTableName: m.TableName,
		ctx:                m.Ctx,
		conn:               m.Conn,
	}
	revRecords, err := schemaMgmt.fetchSchemaRevisionOfMetaStore()
	if nil != err {
		return
	}
	if len(revRecords) == 0 {
		revRecords = []*schemaRevisionOfMetaStore{{
			metaStoreTableName: m.TableName,
		}}
	}
	if schemaChanged, err = schemaMgmt.UpgradeSchemaOfMetaStore(revRecords); nil != err {
		return
	} else if schemaChanged {
		if revRecords, err = schemaMgmt.fetchSchemaRevisionOfMetaStore(); nil != err {
			return
		}
	}
	if !isMetaStoreSchemasUpToDate(revRecords) {
		err = fmt.Errorf("meta-store %s schema not up to date: %#v", m.TableName, revRecords)
	}
	return
}

// MigrateSchemaRevisionKeyGen1 update revision key from legacy integrated meta store.
func (m *MetaStore) MigrateSchemaRevisionKeyGen1(legacyRevisionKey string) (err error) {
	var ok bool
	if ok, _, _, err = m.fetch(legacyRevisionKey); !ok {
		return
	} else if nil != err {
		if mysqlerrors.IsTableNotExistError(err) {
			err = nil
		}
		return
	}
	_, err = m.Conn.Exec(sqlStmtMigrateLegacySchemaRevKeyGen1(m.TableName), makeMetaStoreRevKey(m.TableName), legacyRevisionKey)
	return
}

// Initialize given `metaKey` with given `metaValue`.
func (m *MetaStore) Initialize(metaKey, metaValue string) (err error) {
	modifyAt := time.Now().Unix()
	_, err = m.Conn.Exec(sqlStmtInitMetaValue(m.TableName), metaKey, metaValue, modifyAt)
	return
}

func (m *MetaStore) fetch(metaKey string) (ok bool, value string, modifyAt int64, err error) {
	ok = false
	if err = m.Conn.QueryRowContext(
		m.Ctx,
		sqlStmtFetchMetaValueWithKey(m.TableName),
		metaKey).Scan(&value, &modifyAt); nil != err {
		if err == sql.ErrNoRows {
			return false, "", 0, nil
		}
		return
	}
	ok = true
	return
}

func (m *MetaStore) store(metaKey, metaValue string) (err error) {
	modifyAt := time.Now().Unix()
	_, err = m.Conn.ExecContext(
		m.Ctx,
		sqlStmtStoreMetaValue(m.TableName),
		metaKey, metaValue, modifyAt, metaValue, modifyAt)
	return
}

// FetchBool get bool value from store.
func (m *MetaStore) FetchBool(metaKey string, defaultValue bool) (value bool, modifyAt int64, err error) {
	ok, textValue, modifyAt, err := m.fetch(metaKey)
	if nil != err {
		return
	}
	if ok {
		if "1" == textValue {
			value = true
		} else {
			value = false
		}
	} else {
		value = defaultValue
	}
	return
}

// StoreBool put bool value into store.
func (m *MetaStore) StoreBool(metaKey string, value bool) (err error) {
	var textValue string
	if value {
		textValue = "1"
	} else {
		textValue = "0"
	}
	return m.store(metaKey, textValue)
}

// FetchInt32 get int32 value from store.
func (m *MetaStore) FetchInt32(metaKey string, defaultValue int32) (value int32, modifyAt int64, err error) {
	ok, textValue, modifyAt, err := m.fetch(metaKey)
	if nil != err {
		return
	}
	if ok {
		var vRaw int64
		if vRaw, err = strconv.ParseInt(textValue, 10, 32); nil != err {
			value = defaultValue
		} else {
			value = int32(vRaw)
		}
	} else {
		value = defaultValue
	}
	return
}

// StoreInt32 put int32 value into store.
func (m *MetaStore) StoreInt32(metaKey string, value int32) (err error) {
	textValue := strconv.FormatInt(int64(value), 10)
	return m.store(metaKey, textValue)
}

// FetchInt64 get int64 value from store.
func (m *MetaStore) FetchInt64(metaKey string, defaultValue int64) (value, modifyAt int64, err error) {
	ok, textValue, modifyAt, err := m.fetch(metaKey)
	if nil != err {
		return
	}
	if ok {
		if value, err = strconv.ParseInt(textValue, 10, 64); nil != err {
			value = defaultValue
		}
	} else {
		value = defaultValue
	}
	return
}

// StoreInt64 put int64 value into store.
func (m *MetaStore) StoreInt64(metaKey string, value int64) (err error) {
	textValue := strconv.FormatInt(value, 10)
	return m.store(metaKey, textValue)
}

// FetchRevision get revision value from store.
// Return 0 if revision record not exists.
func (m *MetaStore) FetchRevision(metaKey string) (revValue int32, modifyAt int64, err error) {
	return m.FetchInt32(metaKey, 0)
}

// StoreRevision save revision record into store.
func (m *MetaStore) StoreRevision(metaKey string, revValue int32) (err error) {
	return m.StoreInt32(metaKey, revValue)
}
