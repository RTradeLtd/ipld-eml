package analysis

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	outdir := "outdir"
	t.Cleanup(func() {
		os.RemoveAll(outdir)
	})
	messages, err := GenerateMessages(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 10 {
		t.Fatal("bad number of messages")
	}
	if err := WritePartsToDisk(messages, outdir); err != nil {
		t.Fatal(err)
	}
}
