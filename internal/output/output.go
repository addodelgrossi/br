package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Options struct {
	JSON    bool
	CSV     bool
	Quiet   bool
	NoColor bool
}

type Field struct {
	Key   string
	Label string
	Value string
}

type Column struct {
	Key   string
	Label string
}

func PrintObject(w io.Writer, opts Options, fields []Field, raw any, quietValue string) error {
	switch {
	case opts.JSON:
		return printJSON(w, raw)
	case opts.CSV:
		writer := csv.NewWriter(w)
		header := make([]string, 0, len(fields))
		values := make([]string, 0, len(fields))
		for _, field := range fields {
			header = append(header, field.Key)
			values = append(values, field.Value)
		}
		if err := writer.Write(header); err != nil {
			return err
		}
		if err := writer.Write(values); err != nil {
			return err
		}
		writer.Flush()
		return writer.Error()
	case opts.Quiet:
		_, err := fmt.Fprintln(w, quietValue)
		return err
	default:
		width := 0
		for _, field := range fields {
			if len(field.Label) > width {
				width = len(field.Label)
			}
		}
		for _, field := range fields {
			if field.Value == "" {
				continue
			}
			if _, err := fmt.Fprintf(w, "%-*s  %s\n", width, field.Label, field.Value); err != nil {
				return err
			}
		}
		return nil
	}
}

func PrintList(w io.Writer, opts Options, columns []Column, rows []map[string]string, raw any, quietKey string) error {
	switch {
	case opts.JSON:
		return printJSON(w, raw)
	case opts.CSV:
		writer := csv.NewWriter(w)
		header := make([]string, 0, len(columns))
		for _, col := range columns {
			header = append(header, col.Key)
		}
		if err := writer.Write(header); err != nil {
			return err
		}
		for _, row := range rows {
			values := make([]string, 0, len(columns))
			for _, col := range columns {
				values = append(values, row[col.Key])
			}
			if err := writer.Write(values); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	case opts.Quiet:
		for _, row := range rows {
			value := row[quietKey]
			if value == "" && len(columns) > 0 {
				value = row[columns[0].Key]
			}
			if _, err := fmt.Fprintln(w, value); err != nil {
				return err
			}
		}
		return nil
	default:
		if len(rows) == 0 {
			_, err := fmt.Fprintln(w, "sem resultados")
			return err
		}

		widths := make(map[string]int, len(columns))
		for _, col := range columns {
			widths[col.Key] = len(col.Label)
		}
		for _, row := range rows {
			for _, col := range columns {
				if len(row[col.Key]) > widths[col.Key] {
					widths[col.Key] = len(row[col.Key])
				}
			}
		}

		var header strings.Builder
		for i, col := range columns {
			if i > 0 {
				header.WriteString("  ")
			}
			header.WriteString(fmt.Sprintf("%-*s", widths[col.Key], col.Label))
		}
		if _, err := fmt.Fprintln(w, header.String()); err != nil {
			return err
		}

		for _, row := range rows {
			var line strings.Builder
			for i, col := range columns {
				if i > 0 {
					line.WriteString("  ")
				}
				line.WriteString(fmt.Sprintf("%-*s", widths[col.Key], row[col.Key]))
			}
			if _, err := fmt.Fprintln(w, line.String()); err != nil {
				return err
			}
		}
		return nil
	}
}

func printJSON(w io.Writer, raw any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(raw)
}
