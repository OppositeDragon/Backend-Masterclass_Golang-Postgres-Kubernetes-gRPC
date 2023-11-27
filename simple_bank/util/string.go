package util

import "database/sql"

func SqlNullStringToStringPtr(nullString sql.NullString) *string {
	if nullString.Valid && nullString.String != "" {
		return &nullString.String
	} else {
		return nil
	}
}

func SqlNullStringToString(nullString sql.NullString) string {
	result := SqlNullStringToStringPtr(nullString)
	if result == nil {
		return ""
	}
	return *result
}

func StringPtrToSqlNullString(str *string) sql.NullString {
	if str == nil || *str == "" {
		return sql.NullString{String: "", Valid: false}
	} else {
		return sql.NullString{String: *str, Valid: true}
	}
}

func StringToSqlNullString(str string) sql.NullString {
	return StringPtrToSqlNullString(&str)
}
