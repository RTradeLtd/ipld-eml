# ipld-eml

`ipld-eml` is an RFC-5322 compliant IPLD email object format. It allows taking emails and storing them as typed objects on IPFS.  Emails are converted into a protocol buffer object before being stored on IPFS. Currently there are two methods of storing on IPFS, as a UnixFS object, or as a dedicated IPLD object.

## unixfs

The workflow for unixfs is similar to the dedicated IPLD object, except we take the protocol buffer object, marshal it, and store it as a unixfs object

## dedicated ipld object

The workflow for this involves manually chunking the email protocol buffer object into chunks of slightly under 1MB in size. These chunks are then recorded in a wrapper object, which is then stored as a unixfs object. Because individual DAG objects can't be larger than 1MB in size, otherwise they will be unable to be transferred through the network, it is possible that storing the email chunk wrapper object will be larger than 1MB in size. As such, the unixfs object type allows us to conveniently not have to deal with the maximum size of the wrapper object.


# samples

## overview

`sample1.eml` is a basic email message with no attachments
`sample2.eml` is an email message with an attachment
`sample3.eml` is `sample2.eml` but forwarded to myself
`sample4.eml` is a few replies to `sample3.eml` and sending the same image back
`sample5.eml` is a few replies to `sample4.eml` with roughly 1.6MB in attachments/embedded files
`sample6.eml` is a reply to `sample5.eml` with CC+BCC, and more files
`sample7.eml` is a reply to `sample6.eml` but with samples 1 -> 6 attached
`sample8.eml` is an email i received from the golang weekly mailing list

## generated

The `generated` directory contains 5000 emails generated with the fake email generator in the `analysis` package. Each email has a randomly generated 720x720 image attached to it, as well as one emoji per paragraph.

The following command was used to generate the data:

```shell
$> ./eml-util gfe --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
```