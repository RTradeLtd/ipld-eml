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

To reliably estimate space savings, and performance there is a set of sample emails included in the repository in the `samples` directory. The root of the samples directory contains emails I've sent to myself as a initial test dataset, and an email I received from a newsletter. The `samples/generated` directory contains 5000 emails randomly generated with the `analysis` package. The samples contained here contain highly duplicated data. It is meant to showcase a best case space savings example.

Overview of the various samples:

`sample1.eml` is a basic email message with no attachments
`sample2.eml` is an email message with an attachment
`sample3.eml` is `sample2.eml` but forwarded to myself
`sample4.eml` is a few replies to `sample3.eml` and sending the same image back
`sample5.eml` is a few replies to `sample4.eml` with roughly 1.6MB in attachments/embedded files
`sample6.eml` is a reply to `sample5.eml` with CC+BCC, and more files
`sample7.eml` is a reply to `sample6.eml` but with samples 1 -> 6 attached
`sample8.eml` is an email i received from the golang weekly mailing list

# samples (generated)

The `generated` directory contains 5000 emails generated with the fake email generator in the `analysis` package. Each email has a randomly generated 720x720 image attached to it, as well as one emoji per paragraph, with a total of 100 paragraphs. 

The following command was used to generate the data:

```shell
$> ./eml-util gfe --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
```

Size on disk is about 144MB and size on IPFS is about 133MB which gives us about an 8% space savings on average

# space savings

The non-generated samples are intended to show-case a best case space savings about, when there is predominantly duplicated data. In the non-generated samples, the bulk of the data is composed of the same picture, meant to simulate a situation where the same photo (ex: cat picture) is sent to many different people. 

The generated samples are inteded to show-case the "average/worst-case" space savings, when deduplication is largely derived the nature of content-addressing because there will occasionally be emails that have a small number of chunks that are shared by other emails, as opposed to the exact same image being sent to multiple different people, which leads to deduplication savings as you only need to store the image once.

| Sample Set | IPLD Format | Number Of Emails | IPFS Size | Disk Size | Space Savings | Scenario |
|------------|-------------|------------------|-----------|-----------|---------------|----------|
| Real | Pure UnixFS | 8 | 1.93MB | 11MB | 578% | Best Case (lots of duplicated emails + images) |
| Generated | Pure UnixFS | 5000 | 133MB | 144MB | 8% | Worst Case (virtually no duplicated emails and images) | 

At face value the worst case savings of 8% might not seem like much. However if we extrapolate to larger data sizes even with 8% savings it makes a huge difference.

| Scenario | Disk Size | IPFS Size | Space Saved (no raid / raid-0) | Space Saved (raid-1)
|----------|-----------|-----------|-------------|------|
| Best | 20PB | 3.46PB | 16.54PB | 33.08PB
| Best | 20GB | 3.46GB | 16.54GB | 33.08GB
| Best | 20MB | 3.46MB | 16.54MB | 33.08MB
| Worst | 20PB | 18.4PB | 1.6PB | 3.2PB
| Worst | 20GB | 18.4GB | 1.6GB | 3.2GB
| Worst | 20MB | 18.4MB | 1.6MB | 3.2MB

Even at 20PB, saving 1.6PB amounts to significant real world financial savings, which when you're operating at that scale of storage is huge. Massive email stores, and archives aren't just taking 20PB and using a bunch of cheap Western Digital disks without any redundancy. They're using enterprise grade hard drives which in and off itself is expensive, but there also using things like RAID, zRAID, etc... which amplifies the space savings even more.

# cli usage

## fake email generation

```shell
$> eml-util generate-fake-emails --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
$> eml-util gen-fake-emails --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
$> eml-util gfe --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
```

## converts emails to ipld eml objects

```shell
$> eml-util --email.dir=samples/generated/10k convert
$> eml-util --email.dir=samples/generated/10k con
$> eml-util --email.dir=samples/generated/10k c
```
## Benchmarking

```shell
$> eml-util benchmark
$> eml-util bench
$> eml-util b
```