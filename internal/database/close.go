package database

func (pDB postgresDB) Close() error {
	return pDB.db.Close()
}
