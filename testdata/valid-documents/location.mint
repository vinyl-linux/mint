type Location {
    +mint:doc:"Location points to a specific location"
    +mint:validate:string_not_empty
    string Location = 0;

    +custom:validate:valid_lat
    +mint:doc:"Latitude relates the latitude of the"
    +mint:doc:"described location"
    float32 Latitude = 1;

    // Note how you can also drop random validators into
    // docs, for..... reasons I guess
    +mint:doc:"Longitude relates the longitude of the"
    +custom:validate:valid_long
    +mint:doc:"described location"
    float32 Longitude = 2
}
