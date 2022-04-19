package liquo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/wawandco/liquo/internal/log"
)

// ChangeSet with SQL and Rollback instructions.
type ChangeSet struct {
	ID          string   `xml:"id,attr"`
	Author      string   `xml:"author,attr"`
	SQL         []string `xml:"sql"`
	RollbackSQL string   `xml:"rollback"`
}

// Execute a changeset takes the SQL part of the changeset and runs it.
func (cs ChangeSet) Execute(conn *pgx.Conn, file string) error {
	var err error
	ctx := context.Background()

	row := conn.QueryRow(context.Background(), `SELECT orderexecuted FROM databasechangelog ORDER BY dateexecuted desc`)
	var order int
	if err = row.Scan(&order); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	var count int
	row = conn.QueryRow(ctx, `SELECT count(*) FROM databasechangelog WHERE id = $1`, cs.ID)
	if err = row.Scan(&count); err == nil && count > 0 {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error checking if changeset %v has already been executed:%w", cs.ID, err)
	}

	_, err = conn.Exec(ctx, cs.sql())
	if err != nil {
		return err
	}

	insertStmt := `
		INSERT
		INTO databasechangelog (id, author, filename, dateexecuted, orderexecuted,exectype)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err = conn.Exec(ctx, insertStmt, cs.ID, cs.Author, file, time.Now(), order+1, "EXECUTED")
	if err != nil {
		return err
	}

	log.Infof("Executed `%v`.", cs.ID)
	order++

	return nil
}

// Rollback the changeset runs the Rollback section of the
// changeset.
func (cs ChangeSet) Rollback(conn *pgx.Conn) error {
	log.Infof("Rolling back %v. \n", cs.ID)
	_, err := conn.Exec(context.Background(), cs.RollbackSQL)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `DELETE FROM databasechangelog WHERE id = $1`, cs.ID)
	if err != nil {
		return err
	}

	return nil
}

// sql concats the sql statements on the SQL array of the
// changeset.
func (cs ChangeSet) sql() string {
	return strings.Join(cs.SQL, "\n")
}
