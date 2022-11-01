package database

import (
	"strconv"
	"strings"
)

type setter interface {
	Set(src interface{}) error
}

type Entity interface {
	FieldMap() ([]string, []interface{})
	TableName() string
}

type Entities interface {
	Add() Entity
}

func AllNullEntity(e Entity) {
	_, fields := e.FieldMap()
	for _, field := range fields {
		f, ok := field.(setter)
		if ok {
			_ = f.Set(nil)
		}
	}
}

func GetFieldNames(e Entity) []string {
	fieldNames, _ := e.FieldMap()
	return fieldNames
}

func GetScanFields(e Entity, reqlist []string) []interface{} {
	allNames, allValues := e.FieldMap()
	// Allocate enough capacity for result slice
	n := len(allValues)
	if len(reqlist) < n {
		n = len(reqlist)
	}
	result := make([]interface{}, 0, n)
	for _, reqname := range reqlist {
		for i, name := range allNames {
			if name == reqname {
				result = append(result, allValues[i])
				break
			}
		}
	}
	return result
}

// GeneratePlaceholders returns a string of "$1, $2, ..., $n".
func GeneratePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	var builder strings.Builder
	sep := ", "
	for i := 1; i <= n; i++ {
		if i == n {
			sep = ""
		}
		builder.WriteString("$" + strconv.Itoa(i) + sep)
	}

	return builder.String()
}

// GetFieldNamesExcepts returns all field names from entity e excepts exceptedFieldNames.
func GetFieldNamesExcepts(e Entity, ignoredFieldNames []string) []string {
	numberIgnoredFieldNames := len(ignoredFieldNames)
	fieldNames, _ := e.FieldMap()
	if numberIgnoredFieldNames == 0 {
		return fieldNames
	}
	mapIgnoredFieldNames := make(map[string]bool)
	for _, exceptedFieldName := range ignoredFieldNames {
		mapIgnoredFieldNames[exceptedFieldName] = true
	}
	result := make([]string, 0, len(fieldNames)-numberIgnoredFieldNames)
	for _, fieldName := range fieldNames {
		if mapIgnoredFieldNames[fieldName] {
			continue
		}
		result = append(result, fieldName)
	}
	return result
}
