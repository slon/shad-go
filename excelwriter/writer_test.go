package excelwriter_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"

	"gitlab.com/slon/shad-go/excelwriter"
)

const Sheet = "Sheet1"

type MyInt int

type Basic struct {
	I    int    `xlsx:"i"`
	UI   uint   `xlsx:"ui"`
	I16  int16  `xlsx:"i_16"`
	UI16 uint16 `xlsx:"ui_16"`
	I32  int32  `xlsx:"i_32"`
	UI32 uint32 `xlsx:"ui_32"`
	I64  int64  `xlsx:"i_64"`
	UI64 uint64 `xlsx:"ui_64"`

	Float  float32 `xlsx:"float"`
	Double float64 `xlsx:"double"`
	Bool   bool    `xlsx:"bool"`
	String string  `xlsx:"string"`

	MyInt MyInt `xlsx:"my_int"`
}

type Inner struct {
	Name string `xlsx:"name"`
}

type TextInt int

func (i TextInt) MarshalText() ([]byte, error) {
	return []byte("__" + fmt.Sprint(i) + "__"), nil
}

type TextStruct struct {
	V string
}

func (s TextStruct) MarshalText() (text []byte, err error) {
	return []byte("__" + s.V + "__"), nil
}

func TestWriter_basic(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet(Sheet)
	require.NoError(t, err)

	w := excelwriter.New(f)

	b := &Basic{
		I:    -1,
		UI:   1,
		I16:  -16,
		UI16: 16,
		I32:  -32,
		UI32: 32,
		I64:  -64,
		UI64: 64,

		Float:  -2.5,
		Double: 5.5,
		Bool:   true,
		String: "hello",

		MyInt: 23,
	}
	require.NoError(t, w.WriteRow(b))

	storeFile(t, f, "basic")

	rows := readAll(t, f, true)
	require.Equal(t, []string{"i", "ui", "i_16", "ui_16", "i_32", "ui_32", "i_64", "ui_64", "float", "double", "bool", "string", "my_int"}, rows[0])
	require.Equal(t, []string{"-1", "1", "-16", "16", "-32", "32", "-64", "64", "-2.5", "5.5", "1", "hello", "23"}, rows[1])
}

func TestWriter_map(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet(Sheet)
	require.NoError(t, err)

	w := excelwriter.New(f)

	require.NoError(t, w.WriteRow(map[string]any{
		"id":    29,
		"name":  "hello",
		"title": "yo",
	}))

	require.NoError(t, w.WriteRow(map[string]any{}))

	require.NoError(t, w.WriteRow(map[TextStruct]any{
		{"id"}:   18,
		{"name"}: "tm",
	}))

	storeFile(t, f, "map")

	rows := readAll(t, f, true)
	require.Equal(t, []string{"id", "name", "title", "__id__", "__name__"}, rows[0])
	require.Equal(t, []string{"29", "hello", "yo"}, rows[1])
	require.Equal(t, []string{"", "", "", "18", "tm"}, rows[2])
}

type E struct {
	ID      string `xlsx:"id"`
	NoTag   string
	Skipped int     `xlsx:"-"`
	private float32 `xlsx:"private"`

	TextInt    TextInt    `xlsx:"text-int"`
	TextStruct TextStruct `xlsx:"text-struct"`

	Struct    Inner  `xlsx:"inner"`
	StructPtr *Inner `xlsx:"inner-ptr"`

	Slice []int  `xlsx:"slice"`
	Array [2]int `xlsx:"arr"`

	Map   map[string]string     `xlsx:"map"`
	TMMap map[TextStruct]string `xlsx:"tm-map"`

	F  func()   `xlsx:"func"`
	Ch chan int `xlsx:"chan"`
}

func TestWriter_elaborate(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet("Sheet1")
	require.NoError(t, err)

	w := excelwriter.New(f)

	e := &E{
		ID:         "42",
		NoTag:      "no-tag",
		Skipped:    29,
		private:    2,
		TextInt:    77,
		TextStruct: TextStruct{V: "text-struct"},
		Struct:     Inner{Name: "inner"},
		StructPtr:  &Inner{Name: "inner-ptr"},
		Slice:      []int{1, 2, 3},
		Array:      [2]int{4, 5},
		Map:        map[string]string{"foo": "bar"},
		TMMap:      map[TextStruct]string{{"foo"}: "bar"},
		F:          func() { fmt.Println("hello") },
		Ch:         make(chan int, 10),
	}
	require.NoError(t, w.WriteRow(e))

	storeFile(t, f, "elaborate")

	rows := readAll(t, f, true)
	require.Equal(t, []string{"id", "notag", "private", "text-int", "text-struct", "inner", "inner-ptr", "slice", "arr", "map", "tm-map"}, rows[0])
	require.Equal(t, []string{"42", "no-tag", "2", "__77__", "__text-struct__", `{"Name":"inner"}`, `{"Name":"inner-ptr"}`, "[1,2,3]", "[4,5]", `{"foo":"bar"}`, `{"__foo__":"bar"}`}, rows[1])
}

func TestWriter_ptr(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet(Sheet)
	require.NoError(t, err)

	w := excelwriter.New(f)

	require.NoError(t, w.WriteRow(struct{ Name string }{"foo"}))
	require.NoError(t, w.WriteRow(&struct{ Name string }{"bar"}))
	require.NoError(t, w.WriteRow(map[string]string{"name": "baz"}))
	require.NoError(t, w.WriteRow(&map[string]string{"name": "qux"}))

	storeFile(t, f, "ptr")

	rows := readAll(t, f, true)
	require.Equal(t, []string{"name"}, rows[0])
	require.Equal(t, []string{"foo"}, rows[1])
	require.Equal(t, []string{"bar"}, rows[2])
	require.Equal(t, []string{"baz"}, rows[3])
	require.Equal(t, []string{"qux"}, rows[4])
}

func TestWriter_unsupported(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet("Sheet1")
	require.NoError(t, err)

	w := excelwriter.New(f)

	require.Error(t, w.WriteRow(2))
	require.Error(t, w.WriteRow(nil))
	var b *Basic
	require.Error(t, w.WriteRow(b))
}

type S struct {
	P  float64 `xlsx:"percent,numfmt:9"`
	D  int     `xlsx:"date,numfmt:15"`
	D2 int     `xlsx:",numfmt:15"`
	T  float64 `xlsx:"time,numfmt:22"`
	Y  int     `xlsx:"jpy,numfmt:194"`
}

func TestWriter_style(t *testing.T) {
	f := excelize.NewFile()
	_, err := f.NewSheet("Sheet1")
	require.NoError(t, err)

	w := excelwriter.New(f)

	s := &S{
		P:  0.5,
		D:  100,
		D2: 100,
		T:  100.5,
		Y:  200,
	}
	require.NoError(t, w.WriteRow(s))

	storeFile(t, f, "style")

	rows := readAll(t, f, false)
	require.Equal(t, []string{"percent", "date", "d2", "time", "jpy"}, rows[0])
	require.Equal(t, []string{"50%", "9-Apr-00", "9-Apr-00", "4/9/00 12:00", "¥200.00"}, rows[1])
}

func readAll(t *testing.T, f *excelize.File, raw bool) [][]string {
	t.Helper()

	var out [][]string

	rows, err := f.Rows(Sheet)
	require.NoError(t, err)

	var opts []excelize.Options
	if raw {
		opts = append(opts, excelize.Options{RawCellValue: true})
	}

	for i := 1; rows.Next(); i++ {
		row, err := rows.Columns(opts...)
		require.NoError(t, err)

		out = append(out, row)
	}

	return out
}

func storeFile(t *testing.T, f *excelize.File, prefix string) {
	t.Helper()

	out, err := os.CreateTemp("/tmp", prefix+"-*.xlsx")
	require.NoError(t, err)

	_, err = f.WriteTo(out)
	require.NoError(t, err)

	t.Logf("stored file in %s", out.Name())
}
