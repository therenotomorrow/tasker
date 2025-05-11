tasker
======

<p>
<a href="https://github.com/therenotomorrow/tasker/actions?query=workflow%3ACI+event%3Apush+branch%3Amaster" target="_blank">
    <img src="https://github.com/therenotomorrow/tasker/actions/workflows/ci.yml/badge.svg?event=push&branch=master" alt="CI">
</a>
<a href="https://codecov.io/gh/therenotomorrow/tasker" target="_blank">
    <img src="https://codecov.io/gh/therenotomorrow/tasker/graph/badge.svg?token=CGYK1Y72S2" alt="Coverage">
</a>
<a href="https://github.com/therenotomorrow/tasker/releases" target="_blank">
    <img src="https://img.shields.io/github/v/release/therenotomorrow/tasker" alt="Releases">
</a>
<a href="https://roadmap.sh/projects/task-tracker" target="_blank">
    <img src="https://img.shields.io/badge/project-task_tracker-blue" alt="Project">
</a>
</p>

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
