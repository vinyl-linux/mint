<p align="center">
  <img src="docs/logo.png" />
</p>

# mint

[![Coverage Status](https://coveralls.io/repos/github/vinyl-linux/mint/badge.svg?branch=main)](https://coveralls.io/github/vinyl-linux/mint?branch=main)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_mint&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=vinyl-linux_mint)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_mint&metric=security_rating)](https://sonarcloud.io/dashboard?id=vinyl-linux_mint)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=vinyl-linux_mint&metric=sqale_index)](https://sonarcloud.io/dashboard?id=vinyl-linux_mint)

Mint, a contraction of Message Interchange, is a binary message format (similar to protobuf, thrift, etc.) designed to provide fast serialisation, a reasonably rich set of data types, and an easily extensible system of type validations. A deeper dive into what this means [may be here](docs/scheme.md)

The niche mint fills, over such illustrious alternatives, is simplicity; both on the wire and off.

Mint was designed to:

1. Provide a non-self-describing stream of binary data representing scalars, arrays, maps, and complex types to be as small and fast on the wire as possible
2. Provide a simple annotation schema for mint definitions (documents in our parlance) to drive _validation_ and _transformation_ functions thus shifting data quality assertions left
3. Provide a mechanism for fixed length arrays, thus avoiding extraneous field-length bytes
4. Provide a rich set of datatypes, avoiding the complexity tools like protobuf bring with imports, and complex nested types for something as simple as timestamps

Mint is also designed with code generation in mind; a mint document should describe data types in such a way that code may be generated.

## Sample document

```protobuf
type Location {
    +mint:doc:"Location points to a specific location"
    +mint:validate:string_not_empty
    string Location = 0;

    +custom:validate:valid_lat
    +mint:doc:"Latitude relates the latitude of the"
    +mint:doc:"described location"
    float32 Latitude = 1;

    +mint:doc:"Longitude relates the longitude of the"
    +custom:validate:valid_long
    +mint:doc:"described location"
    float32 Longitude = 2

    +mint:doc:"Tags contain an arbitrary list of tags for"
    +mint:doc:"labelling this location in some way"
    []string Tags = 3;

    +mint:doc:"Labels contains a map key/values representing"
    +mint:doc:"this location in some way"
    map<string,string> Labels = 4;

    +mint:doc:"ID is a UUID representing this location and is"
    +mint:doc:"ever unchanging (whereas the Location name may)"
    uuid ID = 5;
}

type WeatherForecast {
    +mint:doc:"ForecastedAt contains the datetime at which this forecast"
    +mint:doc:"was created"
    +mint:validate:date_in_past
    +mint:transform:date_in_utc
    datetime ForecastedAt = 0;

    +mint:doc:"Location contains a reference to the specified"
    +mint:doc:"location of this forecast"
    Location Location = 1;

    +mint:doc:"Temperature is a float pointing to the forecast"
    +mint:doc:"temperature"
    +custom:validate:seems_valid
    float32 Temperature = 3;

    +mint:doc:"CloudCoverage provides how cloudy it is in"
    +mint:doc:"oktas"
    +custom:validate:valid_okta
    int32 CloudCoverage = 2;

    +mint:doc:"Date this forecast is for"
    +mint:transform:date_in_utc
    datetime ForecastedFor = 4;

    +mint:doc:"WeatherKeys is a tuple that holds some arbitrary data that means..."
    +mint:doc:"something"
    [5]int16 WeatherKeys = 5;
}
```

Mint documents are, essentially, protobuf documents with a few distinctions:

1. They neither define nor power services, such as gRPC services
2. They are annotated with validations and transformations which are called when types are serialised to binary

These distinctions are explored further in [parser/](parser/)

## Installation

```bash
$ curl https://github.com/vinyl-linux/mint/releases/download/0.4.0/mint -o /usr/local/bin/mint
$ mint
mint, a contraction of Message Interchange, is a binary message format
(similar to protobuf, thrift, etc.) designed to provide fast serialisation,
a reasonably rich set of data types, and an easily extensible system of
type validations.

The mint command is used to manage mint documents,
including code generation, validation, and more

Usage:
  mint [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate go code from mint documents
  help        Help about any command
  validate    Validate the documents in a directory

Flags:
      --config string   config file (default is /home/user/.config/mint/mint.yaml)
  -h, --help            help for mint

Use "mint [command] --help" for more information about a command.
```

### Usage

The two main uses for the mint command are:

1. To validate mint documents against the reference parser implementation provided by this package; and
2. To generate go code from documents.

Both uses are served in much the same way:

```bash
$ mint validate path/to/mint-documents
$ mint generate path/to/mint-documents
```

Respectively.

Both commands come with various flags and options, which can be seen (respectively) with:

```bash
$ mint help validate
$ mint help generate
```
