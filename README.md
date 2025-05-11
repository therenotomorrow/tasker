tasker
======

> One of [roadmap.sh](https://roadmap.sh/projects) project. This is my small hobby.

Goal
----

`tasker` is a project used to track and manage tasks. It is a simple command line interface to track:
 - what you need to do
 - what you have done
 - what you are currently working on

System Requirements
-------------------

```shell
go version
# go version go1.24.x ...
```

Development
-----------

Download sources

```shell
PROJECT_ROOT=tasker
git clone https://github.com/therenotomorrow/tasker.git "$PROJECT_ROOT"
cd "$PROJECT_ROOT"
```

Taste it :heart:

```shell
# check code integrity
make code
# run application
./bin/tasker help
# use custom file
TASKER_FILE=custom.json ./bin/tasker help
```

Setup safe development

```shell
git config --local core.hooksPath .githooks
```

Testing
-------

Controls by [test.sh](./scripts/test.sh) or [Makefile](./Makefile) and contains:

```shell
# fast unit tests to be sure that no regression was 
make test/smoke
# same as test/smoke but with -race condition check
make test/unit
# integration tests that needed external resources, also with -race condition
make test/integration
# combines both (test/unit and test/integration) to create local coverage report in HTML
make test/coverage
```
