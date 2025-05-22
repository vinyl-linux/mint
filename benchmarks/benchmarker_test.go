package benchmarks

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
)

var (
	rando = int64(rand.Int())
	id    = uuid.Must(uuid.NewV4())

	str1     = makeString(1)
	str10    = makeString(10)
	str100   = makeString(100)
	str1000  = makeString(1_000)
	str10000 = makeString(10_000)

	strS1_10      = makeStringSlice(str1, 10)
	strS10_10     = makeStringSlice(str10, 10)
	strS10_100    = makeStringSlice(str10, 100)
	strS10_1000   = makeStringSlice(str10, 1000)
	strS100_1000  = makeStringSlice(str100, 10000)
	strS100_10000 = makeStringSlice(str100, 100000)
)

func BenchmarkMarshall(b *testing.B) {
	for _, bench := range []struct {
		name                              string
		shortString, longString           string
		shortStringSlice, longStringSlice []string
	}{
		{"small types", str1, str10, strS1_10, strS10_10},
		{"medium types", str10, str100, strS10_10, strS10_100},
		{"loadsa data", str100, str10000, strS10_1000, strS100_10000},
	} {
		b.Run(bench.name, func(b *testing.B) {
			bb := Benchmarker{
				ID:               id,
				SomeNumber:       rando,
				ShortString:      bench.shortString,
				LongString:       bench.longString,
				ManyShortStrings: bench.shortStringSlice,
				ManyLongStrings:  bench.longStringSlice,
			}

			for b.Loop() {
				buf := new(bytes.Buffer)

				err := bb.Marshall(buf)
				if err != nil {
					panic(err)
				}
			}
		})
	}
}

func makeStringSlice(s string, elems int) (out []string) {
	out = make([]string, elems)
	for idx := range out {
		out[idx] = s
	}

	return
}

func makeString(len int) string {
	sb := strings.Builder{}
	for i := 0; i < len; i++ {
		sb.WriteString("A")
	}

	return sb.String()
}
