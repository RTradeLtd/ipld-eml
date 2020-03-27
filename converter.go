package ipldeml

import (
	"bytes"
	"context"
	"io"

	"github.com/DusanKasan/parsemail"
	xpb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/go-temporalx-sdk/client"
	"github.com/RTradeLtd/ipld-eml/pb"
)

// Converter takes eml files and converting them to an ipfs friendly version
type Converter struct {
	ctx     context.Context
	xclient *client.Client
}

// NewConverter instantiates our new converter
func NewConverter(ctx context.Context, xclient *client.Client) *Converter {
	return &Converter{
		ctx:     ctx,
		xclient: xclient,
	}
}

// GetEmail is a helper function to retrieve an email object
// from ipfs, and return its protocol buffer type
func (c *Converter) GetEmail(hash string) (*pb.Email, error) {
	resp, err := c.xclient.DownloadFile(c.ctx, &xpb.DownloadRequest{
		Hash: hash,
	}, false)
	if err != nil {
		return nil, err
	}
	email := new(pb.Email)
	if err := email.Unmarshal(resp.Bytes()); err != nil {
		return nil, err
	}
	// normalize time values
	email.Date = email.Date.UTC()
	return email, nil
}

// PutEmail is a helper function to store an email objecto n ipfs
func (c *Converter) PutEmail(email *pb.Email) (string, error) {
	data, err := email.Marshal()
	if err != nil {
		return "", err
	}
	resp, err := c.xclient.UploadFile(c.ctx, bytes.NewReader(data), 0, nil, false)
	if err != nil {
		return "", err
	}
	return resp.GetHash(), nil
}

// Convert takes a reader for an eml file, and returns the ipfs hash
func (c *Converter) Convert(reader io.Reader) (*pb.Email, error) {
	eml, err := parsemail.Parse(reader)
	if err != nil {
		return nil, err
	}
	email := &pb.Email{
		Headers: pb.Header{
			Values: make(map[string]*pb.Headers, len(eml.Header)),
		},
		Attachments:   make([]pb.Attachment, len(eml.Attachments)),
		EmbeddedFiles: make([]pb.EmbeddedFile, len(eml.EmbeddedFiles)),
	}
	// set header
	for k, v := range eml.Header {
		email.Headers.Values[k] = &pb.Headers{Values: v}
	}
	// set subject
	email.Subject = eml.Subject
	// set the addresses
	addrs := pb.Addresses{
		From:    make([]pb.Address, len(eml.From)),
		ReplyTo: make([]pb.Address, len(eml.ReplyTo)),
		To:      make([]pb.Address, len(eml.To)),
		Cc:      make([]pb.Address, len(eml.Cc)),
		Bcc:     make([]pb.Address, len(eml.Bcc)),
	}
	if eml.Sender != nil {
		addrs.Sender = &pb.Address{
			Name:    eml.Sender.Name,
			Address: eml.Sender.Address,
		}
	}
	for i, v := range eml.From {
		addrs.From[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	for i, v := range eml.ReplyTo {
		addrs.ReplyTo[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	for i, v := range eml.To {
		addrs.To[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	for i, v := range eml.Cc {
		addrs.Cc[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	for i, v := range eml.Bcc {
		addrs.Bcc[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	email.Addresses = addrs
	email.Date = eml.Date.UTC()
	email.MessageID = eml.MessageID
	email.InReplyTo = eml.InReplyTo
	email.References = eml.References
	var resent = &pb.Resent{
		ResentFrom:      make([]pb.Address, len(eml.ResentFrom)),
		ResentTo:        make([]pb.Address, len(eml.ResentTo)),
		ResentCc:        make([]pb.Address, len(eml.ResentCc)),
		ResentBcc:       make([]pb.Address, len(eml.ResentBcc)),
		ResentMessageId: eml.ResentMessageID,
	}
	for i, v := range eml.ResentFrom {
		resent.ResentFrom[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	if eml.ResentSender != nil {
		resent.ResentSender = &pb.Address{
			Name:    eml.ResentSender.Name,
			Address: eml.ResentSender.Address,
		}
	}
	for i, v := range eml.ResentTo {
		resent.ResentTo[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	resent.ResentDate = eml.ResentDate
	for i, v := range eml.ResentCc {
		resent.ResentCc[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	for i, v := range eml.ResentBcc {
		resent.ResentBcc[i] = pb.Address{
			Name:    v.Name,
			Address: v.Address,
		}
	}
	email.Resent = resent
	email.HtmlBody = eml.HTMLBody
	email.TextBody = eml.TextBody
	for i, attach := range eml.Attachments {
		// file size 0 == no progress eports
		resp, err := c.xclient.UploadFile(c.ctx, attach.Data, 0, nil, false)
		if err != nil {
			return nil, err
		}
		email.Attachments[i] = pb.Attachment{
			FileName:    attach.Filename,
			ContentType: attach.ContentType,
			DataHash:    resp.GetHash(),
		}
	}
	for i, embed := range eml.EmbeddedFiles {
		resp, err := c.xclient.UploadFile(c.ctx, embed.Data, 0, nil, false)
		if err != nil {
			return nil, err
		}
		email.EmbeddedFiles[i] = pb.EmbeddedFile{
			ContentId:   embed.CID,
			ContentType: embed.ContentType,
			DataHash:    resp.GetHash(),
		}
	}
	return email, nil
}
