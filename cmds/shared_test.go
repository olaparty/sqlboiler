package cmds

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestOutHandler(t *testing.T) {
	buf := &bytes.Buffer{}

	saveTestHarnessStdout := testHarnessStdout
	testHarnessStdout = buf
	defer func() {
		testHarnessStdout = saveTestHarnessStdout
	}()

	data := tplData{
		Table: "patrick",
	}

	templateOutputs := [][]byte{[]byte("hello world"), []byte("patrick's dreams")}

	if err := outHandler("", templateOutputs, &data); err != nil {
		t.Error(err)
	}

	if out := buf.String(); out != "hello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}
}

type NopWriteCloser struct {
	io.Writer
}

func (NopWriteCloser) Close() error {
	return nil
}

func nopCloser(w io.Writer) io.WriteCloser {
	return NopWriteCloser{w}
}

func TestOutHandlerFiles(t *testing.T) {
	saveTestHarnessFileOpen := testHarnessFileOpen
	defer func() {
		testHarnessFileOpen = saveTestHarnessFileOpen
	}()

	file := &bytes.Buffer{}
	testHarnessFileOpen = func(path string) (io.WriteCloser, error) {
		return nopCloser(file), nil
	}

	data := tplData{
		Table: "patrick",
	}

	templateOutputs := [][]byte{[]byte("hello world"), []byte("patrick's dreams")}

	if err := outHandler("folder", templateOutputs, &data); err != nil {
		t.Error(err)
	}

	if out := file.String(); out != "hello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}
}

func TestBuildImportString(t *testing.T) {
}

func TestCombineImports(t *testing.T) {
	a := imports{
		standard:   []string{"fmt"},
		thirdparty: []string{"github.com/pobri19/sqlboiler", "gopkg.in/guregu/null.v3"},
	}
	b := imports{
		standard:   []string{"os"},
		thirdparty: []string{"github.com/pobri19/sqlboiler"},
	}

	c := combineImports(a, b)

	if c.standard[0] != "fmt" && c.standard[1] != "os" {
		t.Errorf("Wanted: fmt, os got: %#v", c.standard)
	}
	if c.thirdparty[0] != "github.com/pobri19/sqlboiler" && c.thirdparty[1] != "gopkg.in/guregu/null.v3" {
		t.Errorf("Wanted: github.com/pobri19/sqlboiler, gopkg.in/guregu/null.v3 got: %#v", c.thirdparty)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	hasDups := func(possible []string) error {
		for i := 0; i < len(possible)-1; i++ {
			for j := i + 1; j < len(possible); j++ {
				if possible[i] == possible[j] {
					return fmt.Errorf("found duplicate: %s [%d] [%d]", possible[i], i, j)
				}
			}
		}

		return nil
	}

	if len(removeDuplicates([]string{})) != 0 {
		t.Error("It should have returned an empty slice")
	}

	oneItem := []string{"patrick"}
	slice := removeDuplicates(oneItem)
	if ln := len(slice); ln != 1 {
		t.Error("Length was wrong:", ln)
	} else if oneItem[0] != slice[0] {
		t.Errorf("Slices differ: %#v %#v", oneItem, slice)
	}

	slice = removeDuplicates([]string{"hello", "patrick", "hello"})
	if ln := len(slice); ln != 2 {
		t.Error("Length was wrong:", ln)
	}
	if err := hasDups(slice); err != nil {
		t.Error(err)
	}

	slice = removeDuplicates([]string{"five", "patrick", "hello", "hello", "patrick", "hello", "hello"})
	if ln := len(slice); ln != 3 {
		t.Error("Length was wrong:", ln)
	}
	if err := hasDups(slice); err != nil {
		t.Error(err)
	}
}

func TestCombineStringSlices(t *testing.T) {
	var a, b []string
	slice := combineStringSlices(a, b)
	if ln := len(slice); ln != 0 {
		t.Error("Len was wrong:", ln)
	}

	a = []string{"1", "2"}
	slice = combineStringSlices(a, b)
	if ln := len(slice); ln != 2 {
		t.Error("Len was wrong:", ln)
	} else if slice[0] != a[0] || slice[1] != a[1] {
		t.Errorf("Slice mismatch: %#v %#v", a, slice)
	}

	b = a
	a = nil
	slice = combineStringSlices(a, b)
	if ln := len(slice); ln != 2 {
		t.Error("Len was wrong:", ln)
	} else if slice[0] != b[0] || slice[1] != b[1] {
		t.Errorf("Slice mismatch: %#v %#v", b, slice)
	}

	a = b
	b = []string{"3", "4"}
	slice = combineStringSlices(a, b)
	if ln := len(slice); ln != 4 {
		t.Error("Len was wrong:", ln)
	} else if slice[0] != a[0] || slice[1] != a[1] || slice[2] != b[0] || slice[3] != b[1] {
		t.Errorf("Slice mismatch: %#v + %#v != #%v", a, b, slice)
	}
}
