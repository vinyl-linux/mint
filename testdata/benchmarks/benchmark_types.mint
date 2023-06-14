type Benchmarker {
	uuid ID = 0;
	string ShortString = 1;
	string LongString = 2;
	[]string ManyShortStrings = 3;
	[]string ManyLongStrings = 4;
	int64 SomeNumber = 5;
}
