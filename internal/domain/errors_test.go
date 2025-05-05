package domain_test

import (
	"testing"

	"github.com/therenotomorrow/tasker/internal/domain"
)

func TestUnitError(t *testing.T) {
	t.Parallel()

	const err = domain.Error("stop")

	if got, want := err.Error(), "stop"; got != want {
		t.Errorf("Error() got = %v, want = %v", got, want)
	}
}
