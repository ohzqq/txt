package sep

import (
	"strings"
)

type Splitter func(string) []string

func Split(str string, seps ...Sep) []string {
	var fields []string

	if len(seps) == 0 {
		seps = []Sep{Space}
	}

	if len(seps) == 1 {
		fields = append(fields, splitStrings(str, seps[0])...)
	} else {

		for _, sep := range seps[1:] {
			if len(fields) == 0 {
				fields = append(fields, splitStrings(str, sep)...)
				break
			}
			for _, s := range fields {
				fields = append(fields, splitStrings(s, sep)...)
			}
		}
	}

	return fields
}

func splitStrings(str string, seps ...Sep) []string {
	//var vals []string
	////for _, s := range str {
	//  vals = append(vals, strings.FieldsFunc(s, ShouldSplit(seps))...)
	//}

	// return vals
	return strings.FieldsFunc(str, ShouldSplit(seps))
}
