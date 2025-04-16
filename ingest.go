package services

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/cheggaaa/pb/v3"

	"clickhouse-ingestor/models"
)

func FlatFileToClickHouse(req models.IngestionRequest, db *sql.DB) (int, error) {
	f, err := os.Open(req.FlatFileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	delimiter := ','
	if len(req.Delimiter) > 0 {
		delimiter = rune(req.Delimiter[0])
	}

	reader := csv.NewReader(f)
	reader.Comma = delimiter

	records, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}

	if len(records) < 2 {
		return 0, fmt.Errorf("no data to ingest")
	}

	headers := records[0]
	dataRows := records[1:]

	selectedCols := headers
	if cols, ok := req.SelectedColumns[req.SelectedTables[0]]; ok && len(cols) > 0 {
		selectedCols = cols
	}

	colIndices := []int{}
	for _, col := range selectedCols {
		found := false
		for i, header := range headers {
			if header == col {
				colIndices = append(colIndices, i)
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf("column %s not found in CSV headers", col)
		}
	}

	placeholders := strings.Repeat("?,", len(selectedCols))
	placeholders = strings.TrimSuffix(placeholders, ",")

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		req.SelectedTables[0],
		strings.Join(selectedCols, ","),
		placeholders,
	))
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	bar := pb.StartNew(len(dataRows))

	for _, row := range dataRows {
		vals := make([]interface{}, len(colIndices))
		for i, idx := range colIndices {
			if idx >= len(row) {
				vals[i] = ""
			} else {
				vals[i] = row[idx]
			}
		}
		if _, err := stmt.Exec(vals...); err != nil {
			tx.Rollback()
			return 0, err
		}
		bar.Increment()
	}
	bar.Finish()

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(dataRows), nil
}

func ClickHouseToFlatFile(req models.IngestionRequest, db *sql.DB) (int, error) {
	f, err := os.Create(req.FlatFileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	delimiter := ','
	if len(req.Delimiter) > 0 {
		delimiter = rune(req.Delimiter[0])
	}

	writer := csv.NewWriter(f)
	writer.Comma = delimiter

	columns := "*"
	if cols, ok := req.SelectedColumns[req.SelectedTables[0]]; ok && len(cols) > 0 {
		columns = strings.Join(cols, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", columns, req.SelectedTables[0])
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	colNames, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	if err := writer.Write(colNames); err != nil {
		return 0, err
	}

	count := 0
	for rows.Next() {
		vals := make([]interface{}, len(colNames))
		ptrs := make([]interface{}, len(colNames))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return 0, err
		}
		record := make([]string, len(colNames))
		for i, v := range vals {
			record[i] = fmt.Sprint(v)
		}
		if err := writer.Write(record); err != nil {
			return 0, err
		}
		count++
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return 0, err
	}

	return count, nil
}

func buildJoinQuery(req models.IngestionRequest) (string, error) {
	if len(req.JoinConditions) == 0 {
		return "", fmt.Errorf("no join conditions provided")
	}
	base := fmt.Sprintf("SELECT * FROM %s", req.JoinConditions[0].LeftTable)
	for _, cond := range req.JoinConditions {
		base += fmt.Sprintf(" %s JOIN %s ON %s.%s = %s.%s",
			cond.JoinType,
			cond.RightTable,
			cond.LeftTable, cond.LeftColumn,
			cond.RightTable, cond.RightColumn,
		)
	}
	if req.OutputTable != "" {
		return fmt.Sprintf("INSERT INTO %s %s", req.OutputTable, base), nil
	}
	return base, nil
}
func ExecuteJoin(req models.IngestionRequest, db *sql.DB) (int, error) {
	query, err := buildJoinQuery(req)
	if err != nil {
		return 0, err
	}
	if _, err := db.Exec(query); err != nil {
		return 0, err
	}
	return 1, nil // success
}
