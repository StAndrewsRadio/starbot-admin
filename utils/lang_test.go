package utils_test

import (
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"testing"
)

func TestFieldsN(t *testing.T) {
	tables := []struct {
		name     string
		input    string
		n        int
		expected []string
	}{
		{"zero", "hello split this please", 0, []string{"hello split this please"}},
		{"long", "!register Monday   10PM @someone  Here's a long show name!", 4, []string{
			"!register", "Monday", "10PM", "@someone", "Here's a long show name!"}},
		{"short", "!register  Monday 10PM @someone    ShortShowName!", 4, []string{
			"!register", "Monday", "10PM", "@someone", "ShortShowName!"}},
		{"leading", " hey this is a thing", 4, []string{
			"hey", "this", "is", "a thing"}},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()

			actual := utils.FieldsN(table.input, table.n)
			if !utils.StringSliceEquals(actual, table.expected) {
				t.Errorf("got %q; expected %q where n=%d", actual, table.expected, table.n)
			}
		})
	}
}
