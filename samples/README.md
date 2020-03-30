# overview

`sample1.eml` is a basic email message with no attachments
`sample2.eml` is an email message with an attachment
`sample3.eml` is `sample2.eml` but forwarded to myself
`sample4.eml` is a few replies to `sample3.eml` and sending the same image back
`sample5.eml` is a few replies to `sample4.eml` with roughly 1.6MB in attachments/embedded files
`sample6.eml` is a reply to `sample5.eml` with CC+BCC, and more files
`sample7.eml` is a reply to `sample6.eml` but with samples 1 -> 6 attached
`sample8.eml` is an email i received from the golang weekly mailing list

# generated

The `generated` directory contains 5000 emails generated with the fake email generator in the `analysis` package. Each email has a randomly generated 720x720 image attached to it, as well as one emoji per paragraph.

The following command was used to generate the data:

```shell
$> ./eml-util gfe --paragraph.count 100 --email.count 5000 --emoji.count 100 --outdir=samples/generated
```