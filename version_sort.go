package version_sort

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

var _ = fmt.Println

type ByVersion []string

func (va ByVersion) Len() int {
	return len(va)
}

func (va ByVersion) Swap(i, j int) {
	va[i], va[j] = va[j], va[i]
}

// tokenize the string returning ints on one channel and strs on the other.
// these are meant to be compared to another parallel tokenizer
// if one channel closes first, it is less than the other unless the
// currently waiting token is a string.
// strings < nothing < int

var versionRegexp = regexp.MustCompile("[0-9]+|\\.|[^0-9.]+")

func tokenize(str string, retint chan int64, retstr chan string) {
	defer close(retint)
	defer close(retstr)
	for _, v := range versionRegexp.FindAllString(str, -1) {
		switch v[0] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i, _ := strconv.ParseUint(v, 10, 64)
			retint <- int64(i)
		case '.':
			// pass
		default:
			retstr <- v
		}
	}
}

func (va ByVersion) Less(i, j int) (ret bool) {
	iic := make(chan int64)
	isc := make(chan string)
	jic := make(chan int64)
	jsc := make(chan string)
	go tokenize(va[i], iic, isc)
	go tokenize(va[j], jic, jsc)

	// if i is a string and j is a string, compare them
	// if i is a string and j is an int, return true
	// if i is a string and j is empty, return true
	// if i is an int and j is a string, return false
	// if i is an int and j is an int, compare them
	// if i is an int and j is empty, return false
	// if i is empty and j is a string, return false
	// if i is empty and j is an int, return true
	// if i is empty and j is empty, return false

	for {
		select {
		case jiv, jiok := <-jic:
			select {
			case iiv, iiok := <-iic:
				if !iiok {
					if !jiok {
						return false
					}
					return true
				} else if !jiok {
					return false
				}
				if iiv < jiv {
					return true
				}
				if iiv > jiv {
					return false
				}
			case _, isok := <-isc:
				if !isok {
					if !jiok {
						return false
					}
					return true
				}
				return true
			}
		case jsv, jsok := <-jsc:
			select {
			case _, iiok := <-iic:
				if !iiok {
					if !jsok {
						return false
					}
					return false
				} else if !jsok {
					return false
				}
				return false
			case isv, isok := <-isc:
				if !isok {
					if !jsok {
						return false
					}
					return false
				} else if !jsok {
					return true
				}

				cmp := bytes.Compare([]byte(isv), []byte(jsv))
				if cmp < 0 {
					return true
				}
				if cmp > 0 {
					return false
				}
			}
		}
	}
	return false // arbitrary.  we can never reach this.
}
