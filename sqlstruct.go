package metastore

import (
	"context"
	"database/sql"
	"fmt"

	mysqlerrors "github.com/semeqetjsakatayza/go-mysql-errors"
)

func makeMetaStoreRevKey(tableName string) string {
	return metaKeyMetaStoreSchemaRev + ":" + tableName
}

func sqlCreateMetaStore(metaStoreTableName string) string {
	return "CREATE TABLE `" + (metaStoreTableName) + "` (" +
		"`meta_key` varchar(128) CHARACTER SET ascii NOT NULL COMMENT 'Key of meta information'," +
		"`meta_value` text NOT NULL COMMENT 'Value of meta information'," +
		"`modify_at` bigint(20) NOT NULL DEFAULT '0' COMMENT 'Timestamp of last update'," +
		"PRIMARY KEY (`meta_key`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='" + (metaStoreTableName) + " storage table'"
}

func makeSQLMigrateMetaStoreToRev2(metaStoreTableName string) string {
	return "ALTER TABLE `" + (metaStoreTableName) + "`" +
		" CHANGE COLUMN `meta_key` `meta_key` VARCHAR(128) CHARACTER SET 'ascii' NOT NULL" +
		" COMMENT 'Key of meta information'"
}

// ** SQL schema external filter

const metaKeyMetaStoreSchemaRev = "meta-store.schema"

const currentMetaStoreSchemaRev = 2

type schemaRevision struct {
	MetaStore []*schemaRevisionOfMetaStore
}

func (rev *schemaRevision) IsUpToDate() bool {
	for _, revRecord := range rev.MetaStore {
		if currentMetaStoreSchemaRev != revRecord.currentRev {
			return false
		}
	}
	return true
}

type schemaManager struct {
	metaStoreTableName string
	ctx                context.Context
	conn               *sql.DB
}

func (m *schemaManager) FetchSchemaRevision() (schemaRev *schemaRevision, err error) {
	schemaRev = &schemaRevision{}
	if schemaRev.MetaStore, err = m.fetchSchemaRevisionOfMetaStore(); nil != err {
		return nil, err
	}
	return schemaRev, nil
}

type schemaRevisionOfMetaStore struct {
	currentRev         int32
	metaStoreTableName string
}

func (m *schemaManager) fetchSchemaRevisionOfMetaStore() (revisionRecords []*schemaRevisionOfMetaStore, err error) {
	metaStore := MetaStore{
		TableName: m.metaStoreTableName,
		Ctx:       m.ctx,
		Conn:      m.conn,
	}
	rev, _, err := metaStore.FetchRevision(makeMetaStoreRevKey(m.metaStoreTableName))
	if nil != err {
		if mysqlerrors.IsTableNotExistError(err) {
			err = nil
		}
		return
	}
	if 0 == rev {
		return
	}
	revisionRecords = []*schemaRevisionOfMetaStore{{
		currentRev:         rev,
		metaStoreTableName: m.metaStoreTableName,
	}}
	return
}

func (m *schemaManager) execMetaStoreSchemaModification(sqlStmt string, metaStoreTableName string, targetRev int32) (err error) {
	if _, err = m.conn.ExecContext(m.ctx, sqlStmt); nil != err {
		return
	}
	metaStore := MetaStore{
		TableName: m.metaStoreTableName,
		Ctx:       m.ctx,
		Conn:      m.conn,
	}
	err = metaStore.StoreRevision(makeMetaStoreRevKey(metaStoreTableName), targetRev)
	return
}

func (m *schemaManager) upgradeSchemaMetaStore(currentRev int32, metaStoreTableName string) (schemaChanged bool, err error) {
	switch currentRev {
	case currentMetaStoreSchemaRev:
		return false, nil
	case 0:
		if err = m.execMetaStoreSchemaModification(sqlCreateMetaStore(metaStoreTableName), metaStoreTableName, currentMetaStoreSchemaRev); nil == err {
			return true, nil
		}
	case 1:
		if err = m.execMetaStoreSchemaModification(makeSQLMigrateMetaStoreToRev2(metaStoreTableName), metaStoreTableName, 2); nil == err {
			schemaChanged = true
		}
		return
	default:
		err = fmt.Errorf("unknown meta-store schema revision: %d", currentRev)
	}
	return
}

func (m *schemaManager) UpgradeSchemaOfMetaStore(revisionRecords []*schemaRevisionOfMetaStore) (schemaChanged bool, err error) {
	for _, revRec := range revisionRecords {
		if changed, err := m.upgradeSchemaMetaStore(revRec.currentRev, revRec.metaStoreTableName); nil != err {
			return schemaChanged, fmt.Errorf("upgrade MetaStore failed (%#v): %#v", revRec, err)
		} else if changed {
			schemaChanged = true
		}
	}
	return schemaChanged, nil
}

// ** Generated code for 1 table entries
