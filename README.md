# ipld-eml

`ipld-eml` is an RFC-5322 compliand IPLD object format for storing email messages, in both a space efficient, and time efficient manner. TemporalX is used as the interface into IPFS. Emails are converted into a protocol buffer object, before being stored onto IPFS. There are currently two methods for storing the IPLD objects:

* Entirely as a UnixFS object
* Chunked into 1MB blocks, with all blocks wrapped in a single unixfs object.

This repository also includes a CLI tool enabling you to convert emails manually, or generate fake emails

# data format overview

## unixfs workflow

* Email is converted into protocol buffer object
* Protocol buffer object is saved onto IPFS as a unixfs object

## chunked workflow

* Email is converted into protocol buffer object
* Object is serialzed
* Chunks the serialized byte slice into slight under 1MB in size
* Store byte slice on IPFS as a block
* Create a protocol buffer "chunked email" object
* Store a map of `chunk number -> block hash`
* Store chunked email object on ipfs as a unixfs object (done to avoid possible isuses with store protocol buffer object directly being larger than 1MB)

The chunked method has a very minor overhead compared to the pure unixfs object, but enables more fine-grained distribution of chunks across nodes in the network

# samples

To reliably estimate space savings, and performance there is a set of sample emails included in the repository in the `samples` directory. The root of the samples directory contains emails I've sent to myself as a initial test dataset, and an email I received from a newsletter. The `samples/generated` directory contains 5000 emails randomly generated with the `analysis` package. 

The following files are from the root of the `samples` directory:

`sample1.eml` is a basic email message with no attachments
`sample2.eml` is an email message with an attachment
`sample3.eml` is `sample2.eml` but forwarded to myself
`sample4.eml` is a few replies to `sample3.eml` and sending the same image back
`sample5.eml` is a few replies to `sample4.eml` with roughly 1.6MB in attachments/embedded files
`sample6.eml` is a reply to `sample5.eml` with CC+BCC, and more files
`sample7.eml` is a reply to `sample6.eml` but with samples 1 -> 6 attached
`sample8.eml` is an email i received from the golang weekly mailing list

## generated

The `generated` directory contains 5000 emails generated with the fake email generator in the `analysis` package. Each email has a randomly generated 720x720 image attached to it, as well as one emoji per paragraph, with a total of 100 paragraphs.

The following command was used to generate the data:

```shell
$> ./eml-util gfe --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
```