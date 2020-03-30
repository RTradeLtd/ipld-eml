package ipldeml

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/gogo/protobuf/proto"
)

var (
	listenAddress string
)

func init() {
	listenAddress = os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "xapi.temporal.cloud:9090"
	}
}

func getSamples(t *testing.T, dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	fs := []string{}
	for _, fh := range files {
		if !fh.IsDir() {
			fs = append(fs, dir+"/"+fh.Name())
		}
	}
	return fs
}

func TestConverterGenerated(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: listenAddress,
		Insecure:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	converter := NewConverter(ctx, cl)
	files := getSamples(t, "samples/generated")
	// only test 1000 since in CI this is taking forever
	for i, file := range files {
		if i >= 1000 {
			break
		}
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
			email2Copy, err := converter.GetEmail(hash)
			if err != nil {
				t.Fatal(err)
			}
			chunkHash, err := converter.PutEmailChunked(email2)
			if err != nil {
				t.Fatal(err)
			}
			email3, err := converter.GetEmailChunked(chunkHash)
			if err != nil {
				t.Fatal(err)
			}
			// tests equality between dedicated ipld format, and unixfs format
			if !proto.Equal(email2Copy, email3) {
				t.Fatal("invalid email")
			}
		}()
	}
}

func TestConverter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: listenAddress,
		Insecure:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	converter := NewConverter(ctx, cl)
	files := getSamples(t, "samples")
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
			email2Copy, err := converter.GetEmail(hash)
			if err != nil {
				t.Fatal(err)
			}
			chunkHash, err := converter.PutEmailChunked(email2)
			if err != nil {
				t.Fatal(err)
			}
			email3, err := converter.GetEmailChunked(chunkHash)
			if err != nil {
				t.Fatal(err)
			}
			// tests equality between dedicated ipld format, and unixfs format
			if !proto.Equal(email2Copy, email3) {
				t.Fatal("invalid email")
			}
		}()
	}

}

func TestChunkSizeCalc(t *testing.T) {
	t.Skip("update not yet deplyoed to public temporalx instance")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: listenAddress,
		Insecure:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	converter := NewConverter(ctx, cl)
	files := getSamples(t, "samples")
	var hashes = make([]string, len(files))
	for i, file := range files {
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
		hash, err := converter.PutEmailChunked(email1)
		if err != nil {
			t.Fatal(err)
		}
		hashes[i] = hash
	}
	size, err := converter.CalculateChunkedEmailSize(hashes...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total size of raw blocks: ", size)
}

func TestNonChunkSizeCalc(t *testing.T) {
	t.Skip("update not yet deplyoed to public temporalx instance")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: listenAddress,
		Insecure:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	converter := NewConverter(ctx, cl)
	files := getSamples(t, "samples")
	var hashes []string
	var foundHashes = make(map[string]bool)
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
			if !foundHashes[hash] {
				foundHashes[hash] = true
				hashes = append(hashes, hash)
			}
		}()
	}
	size, err := converter.CalculateEmailSize(false, hashes...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("size: ", size)
}
