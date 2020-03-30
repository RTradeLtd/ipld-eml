package analysis

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	outdir := "outdir"
	t.Cleanup(func() {
		os.RemoveAll(outdir)
	})
	if err := GenerateMessages(outdir, 10, 10, 10); err != nil {
		t.Fatal(err)
	}
	files, err := ioutil.ReadDir(outdir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 10 {
		fmt.Println(len(files))
		t.Fatal("bad number of files")
	}
}
