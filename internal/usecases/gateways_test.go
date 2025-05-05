package usecases_test

import (
	"testing"

	"github.com/therenotomorrow/tasker/internal/usecases"
)

func TestUnitStorageSaver(t *testing.T) {
	t.Parallel()

	var _ usecases.Saver = usecases.Storage(nil)
}

func TestUnitStorageUpdater(t *testing.T) {
	t.Parallel()

	var _ usecases.Updater = usecases.Storage(nil)
}

func TestUnitStorageDeleter(t *testing.T) {
	t.Parallel()

	var _ usecases.Deleter = usecases.Storage(nil)
}

func TestUnitStorageRetriever(t *testing.T) {
	t.Parallel()

	var _ usecases.Retriever = usecases.Storage(nil)
}

func TestUnitUse(t *testing.T) {
	t.Parallel()

	var _ usecases.Use = new(usecases.UseCases)
}
