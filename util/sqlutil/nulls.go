package sqlutil

import (
	"database/sql"
	"time"
)

// ToSqlNullTime transforms a time.Time into a sql.NullTime. If the tm is empty it returns an empty struct.
func ToSqlNullTime(tm time.Time) sql.NullTime {
	if tm.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{Time: tm, Valid: true}
}
