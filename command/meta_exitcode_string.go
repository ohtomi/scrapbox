// Code generated by "stringer -type ExitCode -output meta_exitcode_string.go meta.go"; DO NOT EDIT.

package command

import "fmt"

const _ExitCode_name = "ExitCodeOKExitCodeErrorExitCodeParseFlagsErrorExitCodeBadArgsExitCodeInvalidURLExitCodeProjectNotFoundExitCodePageNotFoundExitCodeFetchFailure"

var _ExitCode_index = [...]uint8{0, 10, 23, 46, 61, 79, 102, 122, 142}

func (i ExitCode) String() string {
	if i < 0 || i >= ExitCode(len(_ExitCode_index)-1) {
		return fmt.Sprintf("ExitCode(%d)", i)
	}
	return _ExitCode_name[_ExitCode_index[i]:_ExitCode_index[i+1]]
}
