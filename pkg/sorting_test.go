package pkg

import (
	"testing"
)

func TestParseSortString(t *testing.T) {
	parsedResult := ParseSortString("name:dsc, foo:asc,bar:asc")

	if parsedResult == nil {
		t.Error("expected non-nil parsedResult")
	}

	if len(parsedResult) != 3 {
		t.Errorf("expected 2 results, got: %d", len(parsedResult))
	}

	expectedResults := []struct {
		field string
		dir   SortDirection
	}{
		{field: "NAME", dir: Descending},
		{field: "FOO", dir: Ascending},
		{field: "BAR", dir: Ascending},
	}

	for i, er := range expectedResults {
		if parsedResult[i].Field != er.field || parsedResult[i].Direction != er.dir {
			t.Errorf("expected at idx: %d field: %s dir: %d, got field: %s dir: %d", i,
				er.field, er.dir, parsedResult[i].Field, parsedResult[i].Direction)
		}
	}
}
