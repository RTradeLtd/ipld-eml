package ipldeml

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
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
	fh, err := os.Open("sample.eml")
	if err != nil {
		t.Fatal(err)
	}
	hash, err := converter.Convert(fh)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("hash: ", hash)
}
