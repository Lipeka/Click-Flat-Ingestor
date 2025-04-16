package services
import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"

	"clickhouse-ingestor/models"
)
func ConnectToClickHouse(req models.IngestionRequest) (*sql.DB, error) {
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s?secure=false",
		req.Username,
		req.JWTToken,
		req.Host,
		req.Port,
		req.Database,
	)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func ListClickHouseTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
func GetTableSchema(db *sql.DB, table string) ([]string, error) {
	query := fmt.Sprintf("DESCRIBE TABLE %s", table)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []string
	for rows.Next() {
		var name, typeCol, defaultType, defaultExpr, comment, codec, ttl string
		if err := rows.Scan(&name, &typeCol, &defaultType, &defaultExpr, &comment, &codec, &ttl); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, nil
}
