package services
import (
    "encoding/csv"
    "os"
)
func ReadFlatFileColumns(filePath string, delimiter rune) ([]string, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    reader := csv.NewReader(f)
    reader.Comma = delimiter
    record, err := reader.Read()
    return record, err
}
func PreviewFlatFile(filePath string, delimiter rune, limit int) ([][]string, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    reader := csv.NewReader(f)
    reader.Comma = delimiter
    records := [][]string{}
    for i := 0; i <= limit; i++ {
        record, err := reader.Read()
        if err != nil {
            break
        }
        records = append(records, record)
    }
    return records, nil
}
