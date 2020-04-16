package repository

import (
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// ParseCypherQueryResult parses a cypher query result record and store the result in target interface
//   - record: neo4j result record
//   - alias: the alias used in cypher query (e.g. m.title)
//   - target: target interface (e.g. models.Movie)
//     Target object should a "db" tab (e.g. `db:"title"`)
func ParseCypherQueryResult(record neo4j.Record, alias string, target interface{}) error {
	elem := reflect.ValueOf(target).Elem()

	for i := 0; i < elem.Type().NumField(); i++ {
		structField := elem.Type().Field(i)

		tag := structField.Tag.Get("db")
		fieldType := structField.Type
		fieldName := structField.Name

		if val, ok := record.Get(fmt.Sprintf("%s.%s", alias, tag)); ok {
			// Ignore nil values
			if val == nil {
				continue
			}
			field := elem.FieldByName(fieldName)
			if field.IsValid() {
				t := fieldType.String()
				switch t {
				case "string":
					field.SetString(val.(string))
				case "int64":
					field.SetInt(val.(int64))
				default:
					return fmt.Errorf("Invalid type: %s", t)
				}
			}
		}

	}

	return nil
}

// BoolPtr returns pointer to a boolean
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr returns pointer to a string
func StringPtr(str string) *string {
	return &str
}

// IntPtr returns pointer to an Int
func IntPtr(i int) *int {
	return &i
}

// PtrOrPtrEmptyString returns pointer to an Int
func PtrOrPtrEmptyString(ptr *string) *string {
	if ptr == nil {
		return StringPtr("")
	}
	return ptr
}

// Float64Ptr returns a pointer to the float64 value passed in.
func Float64Ptr(v float64) *float64 {
	return &v
}
