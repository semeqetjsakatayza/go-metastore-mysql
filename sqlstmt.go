package metastore

func sqlStmtInitMetaValue(metaStoreTableName string) string {
	return "INSERT IGNORE INTO `" + (metaStoreTableName) + "` (`meta_key`, `meta_value`, `modify_at`)" +
		" VALUES (?, ?, ?)"
}

func sqlStmtFetchMetaValueWithKey(metaStoreTableName string) string {
	return "SELECT `meta_value`, `modify_at`" +
		" FROM `" + (metaStoreTableName) + "`" +
		" WHERE (`meta_key` = ?)"
}

func sqlStmtStoreMetaValue(metaStoreTableName string) string {
	return "INSERT INTO `" + (metaStoreTableName) + "` (`meta_key`, `meta_value`, `modify_at`)" +
		" VALUES (?, ?, ?)" +
		" ON DUPLICATE KEY UPDATE `meta_value` = ?, `modify_at` = ?"
}
