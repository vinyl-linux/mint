# Mint Binary Specification

This document describes the binary format of Marshalled data.

## Binary representation

Mint data is **not** self documenting; data is represented as a stream of bytes and parsed by understanding the length of each field.

For fixed-length types this is easy; when a parser expects a 64 bit integer, for instance, it need only read the next 8 bytes and parse them according to the encoding below.

Unbounded types, such as strings, or slices, one of three things happens:

1. strings are prefixed with an int64 containing the length of the string
1. maps, and _most_ slices are prefixed with a uint32 containing the number of bytes to read
1. Some slices have a fixed length; a slice of values corresponding to the months of the year, for instance, will always be of length 12. In this case, no prefixed size need be written

## Types

### Fixed length

Mint supports the following fixed length types. These types are read from a stream of data directly.

| Name     | Description                                                                                              | Size (bytes) | Encoding      |
|----------|----------------------------------------------------------------------------------------------------------|--------------|---------------|
| datetime | int64 of nanoseconds since the [epoch](https://en.wikipedia.org/wiki/Unix_time) without TZ               | 8 Bytes      | Little Endian |
| uuid     | fixed length array of 16 bytes containing the components of a uuid                                       | 16 Bytes     | Little Endian |
| int16    | 16 bit integer, useful where numbers are known to be low                                                 | 2 Bytes      | Little Endian |
| uint16   | Unsigned 16 bit integer                                                                                  | 2 Bytes      | Little Endian |
| int32    | 32 bit integer                                                                                           | 4 Bytes      | Little Endian |
| uint32   | Unsigned 32 bit integer; useful for larger positive numbers where keeping binary size down is important  | 4 Bytes      | Little Endian |
| int64    | 64 bit integer; the recommended numeric type for modern architectures                                    | 8 Bytes      | Little Endian |
| float32  | 32 bit floating point number                                                                             | 4 Bytes      | Little Endian |
| float64  | 64 bit floating point number; recommended for floats unless architecture or space dictates               | 8 Bytes      | Little Endian |
| byte     | Single byte of information, useful for encoded data when coupled with a Slice type                       | 1 Byte       | Little Endian |
| bool     | Boolean value; either true or false                                                                      | 1 Byte       | Little Endian |

There are also the following type aliases:

1. `int8` -> `byte`

### Unbounded types

Mint supports the following types of data which require special care; either they're prefixed with the number of bytes to read, or have a specific length known at compile time

| Name   | Description                                                                                                                                                                                                                                                          | Effective max size   |
|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|
| string | Arbitrary lengthed strings of Bytes, prefixed with an `int64` of the length of the string in bytes                                                                                                                                                                   | 9.2 Exabytes         |
| slice  | Arbitrary lengthed list of *scalars*; slices of slices, or slices of maps, or similar are technically possible but discouraged due to the complexity of how lengths are stored.  Prefixed with a uint32 containing the number of elements in the slice               | 4.2 million elements |
| array  | An array is a fixed length list of data ([we use go's terminology for sequence types](https://go.dev/blog/slices-intro)) and so has no prefixed size. This can be very efficient for known sequence lengths, but a lot of empty/ nil fields are likewise inefficient | Unbounded            |
| map    | A map is an associative slice. It is serialised as a slice of `[k0, v0, k1, v1, ... kn, vn]`                                                                                                                                                                         | 2.1 million elements |
