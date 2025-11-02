package mysqlpermissionadapter

import "database/sql"

// rowsToEntities converts *sql.Rows to [][]any preserving the column order.
func rowsToEntities(rows *sql.Rows) ([][]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	entities := make([][]any, 0)

	for rows.Next() {
		entity := make([]any, len(cols))
		dest := make([]any, len(cols))
		for i := range dest {
			dest[i] = &entity[i]
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
