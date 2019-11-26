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

func structToMap(i interface{}) map[string]string {
	values := make(map[string]string)
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
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
		values[typ.Field(i).Name] = v
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
	return tr
}

func (tr *TableReporter) Append(ps *PageStats) {
	m := structToMap(ps)
	row := make([]string, len(m))
	for _, v := range m {
		row = append(row, v)
	}
	tr.table.Append(row)
}

func (tr *TableReporter) Render() {
	tr.table.Render()
}
