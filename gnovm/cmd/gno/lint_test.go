package main

import "testing"

func TestLintApp(t *testing.T) {
	tc := []testMainCase{
		{
			args:        []string{"lint"},
			errShouldBe: "flag: help requested",
		}, {
			args:                []string{"lint", "--set-exit-status=0", "../../tests/integ/run_main/"},
			stderrShouldContain: "./../../tests/integ/run_main: gno.mod file not found in current or any parent directory (code=1).",
		}, {
			args:                []string{"lint", "--set-exit-status=0", "../../tests/integ/undefined_variable_test/undefined_variables_test.gno"},
			stderrShouldContain: "undefined_variables_test.gno:6: name toto not declared (code=2)",
		}, {
			args:                []string{"lint", "--set-exit-status=0", "../../tests/integ/package_not_declared/main.gno"},
			stderrShouldContain: "main.gno:4: name fmt not declared (code=2).",
		}, {
			args: []string{"lint", "--set-exit-status=0", "../../tests/integ/minimalist_gnomod/"},
			// TODO: raise an error because there is a gno.mod, but no .gno files
		}, {
			args: []string{"lint", "--set-exit-status=0", "../../tests/integ/invalid_module_name/"},
			// TODO: raise an error because gno.mod is invalid
		}, {
			args:                []string{"lint", "--set-exit-status=0", "../../tests/integ/invalid_gno_file/"},
			stderrShouldContain: "../../tests/integ/invalid_gno_file/invalid.gno:1: expected 'package', found packag (code=2).",
		}, {
			args:                []string{"lint", "--set-exit-status=0", "../../tests/integ/typecheck_missing_return/"},
			stderrShouldContain: "../../tests/integ/typecheck_missing_return/main.gno:5:1: missing return (code=4).",
		},

		// TODO: 'gno mod' is valid?
		// TODO: are dependencies valid?
		// TODO: is gno source using unsafe/discouraged features?
		// TODO: check for imports of native libs from non _test.gno files
	}
	testMainCaseRun(t, tc)
}
