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
}
