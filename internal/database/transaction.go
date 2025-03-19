package database

import (
	"context"
	"database/sql"
	"fmt"
)

func (pDB postgresDB) newTransaction() (postgresTx, error) {
	sqlTx, err := pDB.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return postgresTx{}, fmt.Errorf("new db transaction failed: %w", err)
	}

	return postgresTx{sqlTx}, nil
}

func (pTx postgresTx) commitOrRollbackOnError(err *error) {
	if *err != nil {
		pTx.Rollback()
	} else {
		pTx.Commit()
	}
}
