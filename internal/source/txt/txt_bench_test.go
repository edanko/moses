package txt

import (
	"path/filepath"
	"testing"
)

func BenchmarkParseFile(b *testing.B) {
	files, err := filepath.Glob("../../../testdata/txt/*.txt")
	if err != nil {
		b.Fatal(err)
	}
	if len(files) == 0 {
		b.Fatal("no fixtures")
	}

	b.ReportAllocs()
	var sink int
	for b.Loop() {
		for _, file := range files {
			nests, err := ParseFile(file)
			if err != nil {
				b.Fatal(err)
			}
			sink += len(nests)
		}
	}
	_ = sink
}
