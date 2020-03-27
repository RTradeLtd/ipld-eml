# ipld-eml

This repository is a utility to enable taking `eml` files (a common format for storing email messages), and turns them into IPLD objects. It uses `parsemail` to read the eml file, and turn it into a protocol buffer object, with attachments and embeded files stored separately as a unixfs object. The protocol buffer object is then marshalled and stored as a unixfs object as well, to take advantage of its chunking functionalities.

Supports `RFC-5322`.

# overview

One of the really cool things about this is that there is massive savings when it comes to storing multiple emails. For example both `sample2.eml`  and `sample3.eml` reference the same jpeg file. With traditional systems, and as the emails are stored on disk in the eml files, the jpeg image is stored twice. By putting this on IPFS, due to the way we treat embedded files and attachments, we only have to pay the storage cost once.