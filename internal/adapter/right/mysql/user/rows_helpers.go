package mysqluseradapter

import "database/sql"

// rowsToEntities converts *sql.Rows to [][]any preserving column order.
func rowsToEntities(rows *sql.Rows) ([][]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	entities := make([][]any, 0)

	for rows.Next() {
		entity := make([]any, len(columns))
		dest := make([]any, len(columns))
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
