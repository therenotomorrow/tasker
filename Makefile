.PHONY: code test/smoke test/unit test/integration test/coverage

code:
	@"$(CURDIR)/scripts/code.sh"

test/smoke:
	@"$(CURDIR)/scripts/test.sh" smoke

test/unit:
	@"$(CURDIR)/scripts/test.sh" unit

test/integration:
	@"$(CURDIR)/scripts/test.sh" integration

test/coverage:
	@"$(CURDIR)/scripts/test.sh" coverage
