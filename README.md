# mint

Mint, a contraction of Message Interchange, is a binary message format (similar to protobuf, thrift, etc.) designed to provide fast serialisation, a reasonably rich set of data types, and an easily extensible system of type validations.

## Sample document

```protobuf
type WeatherForecast {
    +mint:doc:"ForecastedAt contains the datetime at which this forecast"
    +mint:doc:"was created"
    +mint:validate:date_in_past
    +mint:transform:date_in_utc
    datetime ForecastedAt = 0;

    +mint:doc:"Location contains a reference to the specified"
    +mint:doc:"location of this forecast"
    +mint:validate:string_not_empty
    Location Location = 1;

    +mint:doc:"Temperature is a float pointing to the forecast"
    +mint:doc:"temperature"
    +custom:validate:seems_valid
    float32 Temperature = 3;

    +mint:doc:"CloudCoverage provides how cloudy it is in"
    +mint:doc:"oktas"
    +custom:validate:valid_okta
    int32 CloudCoverage = 2;
}

type Location {
    +mint:doc:"Location points to a specific location"
    string Location = 0;

    +custom:validate:valid_lat
    +mint:doc:"Latitude relates the latitude of the"
    +mint:doc:"described location"
    float32 Latitude = 1;

    // Note how documentation strings don't have to be contiguous
    // for..... reasons I guess
    +mint:doc:"Longitude relates the longitude of the"
    +custom:validate:valid_long
    +mint:doc:"described location"
    float32 Longitude = 2
}
```

Mint documents are, essentially, protobuf documents with a few distinctions:

1. They neither define nor power services, such as gRPC services
2. They are annotated with validations and transformations which are called when types are serialised to binary

## Generating Code from Mint Documents

```bash
$ DIR=testdata/valid-documents mint
```
