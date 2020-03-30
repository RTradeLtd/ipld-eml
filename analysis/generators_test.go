package analysis

import "testing"

func TestGenerate(t *testing.T) {
	messages, err := GenerateMessages(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 10 {
		t.Fatal("bad number of messages")
	}
}
