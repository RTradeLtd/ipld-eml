# ipld-eml

`ipld-eml` is an RFC-5322 compliant IPLD email object format. It allows taking emails and storing them as typed objects on IPFS.  Emails are converted into a protocol buffer object before being stored on IPFS. Currently there are two methods of storing on IPFS, as a UnixFS object, or as a dedicated IPLD object.

## unixfs

The workflow for unixfs is similar to the dedicated IPLD object, except we take the protocol buffer object, marshal it, and store it as a unixfs object

## dedicated ipld object

The workflow for this involves manually chunking the email protocol buffer object into chunks of slightly under 1MB in size. These chunks are then recorded in a wrapper object, which is then stored as a unixfs object. Because individual DAG objects can't be larger than 1MB in size, otherwise they will be unable to be transferred through the network, it is possible that storing the email chunk wrapper object will be larger than 1MB in size. As such, the unixfs object type allows us to conveniently not have to deal with the maximum size of the wrapper object.
