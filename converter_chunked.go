package ipldeml

import (
	"bytes"
	"errors"

	xpb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/ipld-eml/pb"
)

// contains converter function to deal with chunked messages

// GetEmailChunked is used to return an email from its chunked storage format
func (c *Converter) GetEmailChunked(hash string) (*pb.Email, error) {
	ep, err := c.GetChunkedEmail(hash)
	if err != nil {
		return nil, err
	}
	var (
		data []byte
		max  = len(ep.Parts)
	)
	for i := 0; i < max; i++ {
		resp, err := c.xclient.Dag(c.ctx, &xpb.DagRequest{
			RequestType: xpb.DAGREQTYPE_DAG_GET,
			Hash:        ep.Parts[int32(i)],
		})
		if err != nil {
			return nil, err
		}
		data = append(data, resp.GetRawData()...)
	}
	email := new(pb.Email)
	if err := email.Unmarshal(data); err != nil {
		return nil, err
	}
	return email, nil
}

// PutEmailChunked allows storing an email as a custom ipld dag object
// as opposed to a unixfs object type
func (c *Converter) PutEmailChunked(email *pb.Email) (string, error) {
	data, err := email.Marshal()
	if err != nil {
		return "", err
	}
	var dataSize = len(data)
	maxSize := (1024 * 1024 * 1024) - 1024
	if len(data) >= maxSize {
		return "", errors.New("do normal uplaod")
	}
	var (
		parts     = make(map[int32]string)
		lastChunk = 0
	)
	for i := 0; ; i++ {
		if lastChunk >= dataSize {
			break
		}
		barrier := lastChunk + maxSize
		if barrier > dataSize {
			barrier = dataSize
		}
		resp, err := c.xclient.Dag(c.ctx, &xpb.DagRequest{
			Data: data[lastChunk:barrier],
		})
		if err != nil {
			return "", err
		}
		lastChunk = barrier
		parts[int32(i)] = resp.GetHashes()[0]
	}
	ep := &pb.ChunkedEmail{
		Parts: parts,
	}
	epd, err := ep.Marshal()
	if err != nil {
		return "", err
	}
	resp, err := c.xclient.UploadFile(c.ctx, bytes.NewReader(epd), 0, nil, false)
	if err != nil {
		return "", err
	}
	return resp.GetHash(), nil
}

// GetChunkedEmail returns a ChunkedEmail object
func (c *Converter) GetChunkedEmail(hash string) (*pb.ChunkedEmail, error) {
	resp, err := c.xclient.DownloadFile(c.ctx, &xpb.DownloadRequest{Hash: hash}, false)
	if err != nil {
		return nil, err
	}
	ep := new(pb.ChunkedEmail)
	if err := ep.Unmarshal(resp.Bytes()); err != nil {
		return nil, err
	}
	return ep, nil
}

// CalculateChunkedEmailSize is used to calculate the size of chunked ipld eml objects
func (c *Converter) CalculateChunkedEmailSize(hashes ...string) (int64, error) {
	if len(hashes) == 0 {
		return 0, errors.New("no hashes provided")
	}
	var fileHashes = make(map[string]bool)
	var newHashes []string
	for _, hash := range hashes {
		chnk, err := c.GetChunkedEmail(hash)
		if err != nil {
			return 0, err
		}
		for _, chash := range chnk.GetParts() {
			if !fileHashes[chash] {
				fileHashes[chash] = true
				newHashes = append(newHashes, chash)
			}
		}
		em, err := c.GetEmailChunked(hash)
		if err != nil {
			return 0, err
		}
		for _, embed := range em.EmbeddedFiles {
			if !fileHashes[embed.DataHash] {
				fileHashes[embed.DataHash] = true
				newHashes = append(newHashes, embed.DataHash)
			}
		}
		for _, attach := range em.Attachments {
			if !fileHashes[attach.DataHash] {
				fileHashes[attach.DataHash] = true
				newHashes = append(newHashes, attach.DataHash)
			}
		}
	}
	var size int64
	hashes = append(hashes, newHashes...)
	for _, hash := range hashes {
		resp, err := c.xclient.Dag(c.ctx, &xpb.DagRequest{
			RequestType: xpb.DAGREQTYPE_DAG_STAT,
			Hash:        hash,
		})
		if err != nil {
			return 0, err
		}
		size += resp.GetNodeStats()[hash].GetCumulativeSize()
	}
	return size, nil
}
