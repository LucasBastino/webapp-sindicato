package member

import "strings"

func buildMemberFilters(baseQuery string, filters memberFilters) (string, []any) {
	args := []any{}
	query := baseQuery + " WHERE 1=1"

	if filters.searchKey != "" {
		like := filters.searchKey + "%"
		query += " AND (M.name LIKE ? OR M.last_name LIKE ? OR M.dni LIKE ?)"
		args = append(args, like, like, like)
	}

	var statusConditions []string
	if filters.statuses.ShowActive {
		statusConditions = append(statusConditions, " C.deleted_at IS NULL AND M.deleted_at IS NULL")
	}
	if filters.statuses.ShowInactive {
		statusConditions = append(statusConditions, " C.deleted_at IS NOT NULL AND M.deleted_at IS NULL")
	}
	if filters.statuses.ShowDeleted {
		statusConditions = append(statusConditions, " M.deleted_at IS NOT NULL")
	}

	if len(statusConditions) > 0{
		query += " AND ( " + strings.Join(statusConditions, " OR ") + " )" 
	} else{
		// si no hay filtros seleccionados no devuelvo nada
		query += " AND 1=0"
	}

	if filters.companyID != nil {
		query += " AND M.id_company = ?"
		args = append(args, *filters.companyID)
	}


	return query, args
}
