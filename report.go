package main

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type Reporter interface {
	Append(ps *PageStats)
	Render() error
}

func structToMap(ps *PageStats) map[string]string {
	values := make(map[string]string)
	s := reflect.ValueOf(ps).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		values[typeOfT.Field(i).Name] = v
	}
	return values
}

func mapValues(m map[string]string, fields []string) []string {
	row := make([]string, 0, len(fields))
	for _, k := range fields {
		if v, ok := m[k]; ok {
			row = append(row, v)
		}
	}
	return row
}

type TableReporter struct {
	table  *tablewriter.Table
	fields []string
}

func StatHeaders() []string {
	statType := reflect.TypeOf(PageStats{})
	header := make([]string, 0, statType.NumField())
	for i := 0; i < statType.NumField(); i++ {
		field := statType.Field(i)
		header = append(header, field.Name)
	}
	return header
}

func NewTableReporter() *TableReporter {
	tr := &TableReporter{
		table:  tablewriter.NewWriter(os.Stdout),
		fields: StatHeaders(),
	}

	tr.table.SetHeader(tr.fields)
	return tr
}

func NewTSVReporter() *TableReporter {
	tr := NewTableReporter()

	tr.table.SetAutoWrapText(false)
	tr.table.SetAutoFormatHeaders(true)
	tr.table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tr.table.SetAlignment(tablewriter.ALIGN_LEFT)
	tr.table.SetCenterSeparator("")
	tr.table.SetColumnSeparator("")
	tr.table.SetRowSeparator("")
	tr.table.SetHeaderLine(false)
	tr.table.SetBorder(false)
	tr.table.SetTablePadding("\t") // pad with tabs
	tr.table.SetNoWhiteSpace(true)

	return tr
}

func (tr *TableReporter) Append(ps *PageStats) {
	m := structToMap(ps)
	row := mapValues(m, tr.fields)
	tr.table.Append(row)
}

func (tr *TableReporter) Render() error {
	tr.table.Render()
	return nil
}

type CSVReporter struct {
	w      *csv.Writer
	fields []string
}

func NewCSVReporter() *CSVReporter {
	cr := &CSVReporter{
		w:      csv.NewWriter(os.Stdout),
		fields: StatHeaders(),
	}
	cr.w.Write(cr.fields)
	return cr
}

func (cr *CSVReporter) Append(ps *PageStats) {
	m := structToMap(ps)
	row := mapValues(m, cr.fields)
	cr.w.Write(row)
}

func (cr *CSVReporter) Render() error {
	cr.w.Flush()
	return nil
}
