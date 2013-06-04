package version_sort

import (
	"testing"
	"sort"
)

func TestLess(t *testing.T) {
	strs := []string{"0.0", "1.1", "1.2.3.2", "1.2.3.4", "2",}
	if ByVersion(strs).Less(1,0) {
		t.Errorf("%#v > %#v but Less said it's not.\n", strs[1], strs[0])
	}
	for i, iv := range(strs) {
		if ByVersion(strs).Less(i,i) {
			t.Errorf("%#v is not less than itself\n", iv)
		}
		for j, jv := range(strs) {
			if i < j {
				if ByVersion(strs).Less(j,i) {
					t.Errorf("%#v < %#v but Less said it's not\n", iv, jv)
				}
			} else if i > j {
				if ByVersion(strs).Less(i,j) {
					t.Errorf("%#v > %#v but Less said it's not\n", iv, jv)
				}
			}
		}
	}
}

func TestSort(t *testing.T) {
	strs := []string{
		"1.2.3.2",
		"1.1",
		"1.2.3.4",
		"2",
		"0.0",
	}
	sorted_strs := []string{
		strs[4],
		strs[1],
		strs[0],
		strs[2],
		strs[3],
	}
	fail_strs := []string{
		strs[3],
		strs[0],
		strs[4],
		strs[1],
		strs[2],
	}

	sort.Sort(ByVersion(strs))
	for k := range(strs) {
		if strs[k] != sorted_strs[k] {
			t.Errorf("%#v != %#v\n", strs, sorted_strs)
			break
		}
	}
	for k := range(strs) {
		if strs[k] == fail_strs[k] {
			t.Errorf("%#v != %#v\n", strs, fail_strs)
			break
		}
	}
}