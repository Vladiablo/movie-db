package httphelpers

import (
	"strconv"
	"strings"
)

type Int32Slice []int32

func (s *Int32Slice) UnmarshalParam(src string) error {
	var vals []int32
	uniqVals := make(map[int32]struct{}, 2)

	for strVal := range strings.SplitSeq(src, ",") {
		strVal := strings.TrimSpace(strVal)
		if len(strVal) == 0 {
			continue
		}

		parsedVal, err := strconv.ParseInt(strVal, 10, 32)
		if err != nil {
			return err
		}

		val := int32(parsedVal)
		if _, ok := uniqVals[val]; ok {
			continue
		}

		uniqVals[val] = struct{}{}
		vals = append(vals, val)
	}

	*s = vals

	return nil
}
