package util

import "database/sql"

func SqlNullStringToStringPtr(nullString sql.NullString) *string {
	if nullString.Valid && nullString.String != "" {
		return &nullString.String
	} else {
		return nil
	}
}

func StringPtrToSqlNullString(str *string) sql.NullString {
	if str == nil || *str == "" {
		return sql.NullString{String: "", Valid: false}
	} else {
		return sql.NullString{String: *str, Valid: true}
	}
}
