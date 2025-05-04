tasker
======

> One of [roadmap.sh](https://roadmap.sh/projects) project. This is my small hobby.

`tasker` is a project used to track and manage tasks. It is a simple command line interface to track:
 - what you need to do
 - what you have done
 - what you are currently working on

Development
-----------

```shell
# required version
go version # go1.24.2

# setup code quality tools
make code

# setup pre-commit and check integrity
git config --local core.hooksPath .githooks

# just use it :V <3
./bin/tasker help
```
