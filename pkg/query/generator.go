package generator

import (
	"strings"
)

func DynamicUpdateStatement(column []string, json map[string]interface{}) string {
	var columns []string

	for row, dataColumn := range column {
		value := json[dataColumn].(string)
		totalRow := len(column)
		if value != "" {
			if row+1 == totalRow {
				columns = append(columns, dataColumn+" = :"+dataColumn+" ")
			} else {
				columns = append(columns, dataColumn+" = :"+dataColumn+", ")
			}
		}
	}
	queryColumn := strings.Join(columns, "")
	return queryColumn
}
