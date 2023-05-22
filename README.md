# mint

Mint, a contraction of Message Interchange, is a binary message format (similar to protobuf, thrift, etc.) designed to provide fast serialisation, a reasonably rich set of data types, and an easily extensible system of type validations.

## Sample document

```protobuf
// WeatherForecast provides a forecast for the upcoming day
message WeatherForecast {
	// ForecastedAt contains the datetime at which this forecast
	// was created
	// +mint:validate:date_in_past
	// +mint:transform:date_in_utc
	datetime ForecastedAt = 1;
	
	// Location contains a reference to the specified
	// location of this forecast
	// +mint:validate:string_not_empty
	string Location = 2;

	// Temperature is a float pointing to the forecast
	// temperature
	// +custom:validate:seems_valid
	float Temperature = 3;

	// CloudCoverage provides how cloudy it is in
	// oktas
	// +mint:validate:int_range:0:8
	int CloudCoverage = 4;
}
```

Mint documents are, essentially, protobuf documents with a few distinctions:

1. They neither define nor power services, such as gRPC services
2. They are annotated with validations and transformations; including both out of the box
and custom provided.

## Generating Code from Mint Documents

```bash
$ mint doc.mint > doc.go
```

Or, for non-go projects:

```bash
$ mint --output-to rust do.mint > doc.rs
```
