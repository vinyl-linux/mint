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

    +mint:doc:"Tags contain an arbitrary list of tags for"
    +mint:doc:"labelling this location in some way"
    []string Tags = 3;

    +mint:doc:"Labels contains a map key/values representing"
    +mint:doc:"this location in some way"
    map<string,string> Labels = 4;

    +mint:doc:"ID is a UUID representing this location and is"
    +mint:doc:"ever unchanging (whereas the Location name may)"
    uuid ID = 5;

    +mint:doc:"Type of location this is, such as home or whatever"
    LocationType Type = 6;
}

enum LocationType {
     Home
     Work
     School
     Other
}
