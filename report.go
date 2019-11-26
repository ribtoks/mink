package main

import (
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

type TableReporter struct {
	table *tablewriter.Table
}

func NewTableReporter() *TableReporter {
	tr := &TableReporter{
		table: tablewriter.NewWriter(os.Stdout),
	}

	statType := reflect.TypeOf(PageStats{})
	header := make([]string, 0, statType.NumField())
	for i := 0; i < statType.NumField(); i++ {
		field := statType.Field(i)
		header = append(header, field.Name)
	}
	tr.table.SetHeader(header)
	return tr
}

func (tr *TableReporter) Append(ps *PageStats) {
	m := structToMap(ps)
	row := make([]string, 0, len(m))
	for _, v := range m {
		row = append(row, v)
	}
	tr.table.Append(row)
}

func (tr *TableReporter) Render() {
	tr.table.Render()
}
