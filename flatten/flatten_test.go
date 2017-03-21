package flatten

import (
	"testing"
)

func TestGenerateFilename(t *testing.T) {
	type generateFilenameTestCase struct {
		n string
		i int
		l int
	}
	test_cases := []generateFilenameTestCase{
		{"hello", 2, 3},
		{"myFile", 24, 5},
		{"movie.mp4", 1, 1},
		{"song_23.flac", 99, 3},
	}
	expected_results := []string{
		"hello_002",
		"myFile_00024",
		"movie_1.mp4",
		"song_23_099.flac",
	}
	for i, tc := range test_cases {
		res := generateFilename(tc.n, tc.i, tc.l)
		if res != expected_results[i] {
			t.Errorf("generateFilename: expected %v, got %v", expected_results[i], res)
		}
	}

}
