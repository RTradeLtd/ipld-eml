package ipldeml

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/gogo/protobuf/proto"
)

func TestConverter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: "xapi.temporal.cloud:9090",
		Insecure:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	converter := NewConverter(ctx, cl)
	var files = []string{"./samples/sample1.eml", "./samples/sample2.eml", "./samples/sample3.eml"}
	for _, file := range files {
		func() {
			fh, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			defer fh.Close()
			data, err := ioutil.ReadAll(fh)
			if err != nil {
				t.Fatal(err)
			}
			email1, err := converter.Convert(bytes.NewReader(append(data[0:0:0], data...)))
			if err != nil {
				t.Fatal(err)
			}
			hash, err := converter.PutEmail(email1)
			if err != nil {
				t.Fatal(err)
			}
			email2, err := converter.GetEmail(hash)
			if err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(email1, email2) {
				t.Fatal("not equal")
			}
		}()
	}
}
