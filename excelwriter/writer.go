//go:build !solution

package excelwriter

import (
	"github.com/xuri/excelize/v2"
)

type Writer interface {
	// WriteRow appends struct fields or map values to excel sheet.
	//
	// Integers and floats are encoded as excel type Number.
	// Strings are encoded as excel type Text.
	// Booleans are encoded as excel type Logical.
	//
	// Pointer values are encoded as the value pointed to. A nil pointer is skipped.
	//
	// Values implementing encoding.TextMarshaler interface are encoded as excel Text.
	//
	// Interface values are encoded as the value contained in the interface. A nil interface is skipped.
	//
	// Channels and functions are skipped.
	//
	// Structs, maps, slices and arrays are marshaled into json and encoded as excel Text.
	//
	// Encoding of each struct field can be customized by format string
	// stored under the "xlsx" key in the field's tag.
	//
	//	 // Field appears in excel under column "my_field".
	//	 Field int `xlsx:"my_field"`
	//
	//	 // Field appears in excel under column "my_field" formatted with predeclared style 15 ("d-mmm-yy").
	//	 // Only applicable for integers and floats.
	//	 Field int `xlsx:"my_field,numfmt:15"`
	//
	//	// Field is ignored by this package
	//	Field int `xlsx:"-"`
	//
	// The first row is reserved for column names.
	// For structs column name must be either lowercase field name or the name from "xlsx" tag if present.
	// For maps column name is a map key.
	// If map key implements encoding.TextMarshaler then column name is string(key.MarshalText()).
	WriteRow(r any) error
}

func New(f *excelize.File) Writer {
	panic("implement me")
}
