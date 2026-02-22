package helperdb

import "github.com/jackc/pgx/v5/pgtype"

func SafeString(value *string) pgtype.Text {
	return pgtype.Text{String: *value, Valid: true}
}
