// Code generated by "stringer --type OutputType internal/reconciler/logpipeline/reconciler.go"; DO NOT EDIT.

package logpipeline

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OTel-0]
	_ = x[FluentBit-1]
}

const _OutputType_name = "OTelFluentBit"

var _OutputType_index = [...]uint8{0, 4, 13}

func (i OutputType) String() string {
	if i < 0 || i >= OutputType(len(_OutputType_index)-1) {
		return "OutputType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OutputType_name[_OutputType_index[i]:_OutputType_index[i+1]]
}
