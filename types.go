package models
type DataSourceType string
const (
    ClickHouse DataSourceType = "clickhouse"
    FlatFile   DataSourceType = "flatfile"
)
type JoinCondition struct {
    LeftTable   string `json:"left_table"`
    RightTable  string `json:"right_table"`
    LeftColumn  string `json:"left_column"`
    RightColumn string `json:"right_column"`
    JoinType    string `json:"join_type"` // INNER, LEFT, RIGHT etc.
}
type IngestionRequest struct {
    SourceType       DataSourceType        `json:"source_type"`
    Host             string                `json:"host"`
    Port             int                   `json:"port"`
    Database         string                `json:"database"`
    Username         string                `json:"username"`
    JWTToken         string                `json:"jwt_token"`
    FlatFileName     string                `json:"flat_file_name,omitempty"`
    Delimiter        string                `json:"delimiter,omitempty"`
    TargetDirection  string                `json:"target_direction"` // flatfile_to_clickhouse or clickhouse_to_flatfile
    SelectedTables   []string              `json:"selected_tables"`
    SelectedColumns  map[string][]string   `json:"selected_columns"` // table -> []columns
    JoinConditions   []JoinCondition       `json:"join_conditions,omitempty"`
    OutputTable      string                `json:"output_table,omitempty"`
}
