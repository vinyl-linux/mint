type Location {
    +mint:doc:"Location points to a specific location"
    string Location = 1;

    +custom:validate:valid_lat
    +mint:doc:"Latitude relates the latitude of the"
    +mint:doc:"described location"
    float Latitude = 2;

    // Note how you can also drop random validators into
    // docs, for..... reasons I guess
    +mint:doc:"Longitude relates the longitude of the"
    +custom:validate:valid_long
    +mint:doc:"described location"
    float Longitude = 3;
}
