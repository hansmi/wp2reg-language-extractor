package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteTable(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []fileData
		want [][]string
	}{
		{
			name: "empty",
			want: [][]string{
				{"Index"},
			},
		},
		{
			name: "no rows",
			data: []fileData{
				{name: "first"},
				{name: "second"},
			},
			want: [][]string{
				{"Index", "first", "second"},
			},
		},
		{
			name: "one",
			data: []fileData{
				{
					name:    "first",
					entries: []string{"a", "b"},
				},
				{
					name:    "second",
					entries: []string{"x", "y", "z"},
				},
				{
					name:    "third",
					entries: []string{"t0", "t1", "t2"},
				},
			},
			want: [][]string{
				{"Index", "first", "second", "third"},
				{"0", "a", "x", "t0"},
				{"1", "b", "y", "t1"},
				{"2", "", "z", "t2"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := [][]string{}

			if err := writeTable(tc.data, func(row []string) error {
				got = append(got, row)
				return nil
			}); err != nil {
				t.Errorf("writeTable(%v) failed: %v", tc.data, err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Rows difference (-want +got):\n%s", diff)
			}
		})
	}
}
