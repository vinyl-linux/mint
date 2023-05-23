# Mint Parser

This code is almost completely taken from https://github.com/alecthomas/participle/tree/master/_examples/protobuf
with a few modifications.

It is, therefore, licenced under the same terms.

My sincere and heartfelt thanks go to github user [alecthomas](https://github.com/alecthomas) for their work.

## Changes

The following changes have been made:

1. Mint does not support services, and so these are removed
1. Mint does not support embedded messages, and so these are removed
1. Mint has different base types, and these are reflected accordingly
1. Mint has annotations, which protobuf doesn't
