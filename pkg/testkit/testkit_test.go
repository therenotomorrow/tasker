package testkit_test

import (
	"testing"

	"github.com/therenotomorrow/tasker/pkg/testkit"
)

func TestUnitConstError(t *testing.T) {
	t.Parallel()

	const err = testkit.ConstError("const")

	if got, want := err.Error(), "const"; got != want {
		t.Errorf("ConstError() got = %v, want = %v", got, want)
	}
}
