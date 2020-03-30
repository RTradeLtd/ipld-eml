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

func TestConverter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: "xapi.temporal.cloud:9090",
		Insecure:      true,
	})
	defer cl.Close()
	if err != nil {
		t.Fatal(err)
	}
	converter := NewConverter(ctx, cl)
	var files = []string{"./samples/sample1.eml", "./samples/sample2.eml", "./samples/sample3.eml", "./samples/sample4.eml", "./samples/sample5.eml", "./samples/sample6.eml", "./samples/sample7.eml", "./samples/sample8.eml"}
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: "localhost:9090",
		Insecure:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	converter := NewConverter(ctx, cl)
	var files = []string{"./samples/sample1.eml", "./samples/sample2.eml", "./samples/sample3.eml", "./samples/sample4.eml", "./samples/sample5.eml", "./samples/sample6.eml", "./samples/sample7.eml", "./samples/sample8.eml"}
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
	/*	for _, hash := range hashes {
			ch, err := converter.GetChunkedEmail(hash)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := cl.Blockstore(ctx, &pb.BlockstoreRequest{
				RequestType: pb.BSREQTYPE_BS_GET_STATS,
				Cids:        []string{hash},
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, block := range resp.GetBlocks() {
				totalSize += block.GetSize_()
			}
			for _, chash := range ch.GetParts() {
				resp, err := cl.Blockstore(ctx, &pb.BlockstoreRequest{
					RequestType: pb.BSREQTYPE_BS_GET_STATS,
					Cids:        []string{chash},
				})
				if err != nil {
					t.Fatal(err)
				}
				for _, block := range resp.GetBlocks() {
					totalSize += block.GetSize_()
				}
			}
		}
	*/
	size, err := converter.CalculateChunkedEmailSize(hashes...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total size of raw blocks: ", size)
}

func TestNonChunkSizeCalc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cl, err := client.NewClient(client.Opts{
		ListenAddress: "localhost:9090",
		Insecure:      true,
	})
	defer cl.Close()
	if err != nil {
		t.Fatal(err)
	}
	converter := NewConverter(ctx, cl)
	var files = []string{"./samples/sample1.eml", "./samples/sample2.eml", "./samples/sample3.eml", "./samples/sample4.eml", "./samples/sample5.eml", "./samples/sample6.eml", "./samples/sample7.eml", "./samples/sample8.eml"}
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
	size, err := converter.CalculateEmailSize(hashes...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("size: ", size)
}
