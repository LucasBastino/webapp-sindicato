package company

import "strings"

func buildCompanyFilters(baseQuery string, filters companyFilters) (string, []any) {
	args := []any{}
	query := baseQuery + " WHERE 1=1"

	var statusConditions []string
	if filters.statuses.ShowActive {
		statusConditions = append(statusConditions, " deleted_at IS NULL ")
	}
	if filters.statuses.ShowDeleted {
		statusConditions = append(statusConditions, " deleted_at IS NOT NULL ")
	}

	if len(statusConditions) > 0 {
		query += " AND ( " + strings.Join(statusConditions, " OR ") + " ) "
	} else{
		// si no hay filtros seleccionados no devuelvo nada
		query += " AND 1=0 "
	}

	if key := strings.TrimSpace(filters.searchKey); key != "" {
		like := "%" + key + "%"
		query += " AND (name LIKE ? OR company_number LIKE ?)"
		args = append(args, like, like)
	}
	
	return query, args
}
