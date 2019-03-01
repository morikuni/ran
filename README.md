# ran

[![CircleCI](https://circleci.com/gh/morikuni/ran/tree/master.svg?style=shield)](https://circleci.com/gh/morikuni/ran/tree/master)
[![GoDoc](https://godoc.org/github.com/morikuni/ran?status.svg)](https://godoc.org/github.com/morikuni/ran)
[![Go Report Card](https://goreportcard.com/badge/github.com/morikuni/ran)](https://goreportcard.com/report/github.com/morikuni/ran)
[![codecov](https://codecov.io/gh/morikuni/ran/branch/master/graph/badge.svg)](https://codecov.io/gh/morikuni/ran)

ran is a task runner with the concept of event driven. 

## Install

```sh
$ go get github.com/morikuni/ran/cmd/ran
```

## Usage

ran's task file is written in YAML.

Here is an example of task file.

```yaml
env:
  GO111MODULE: on

commands:
  test:
    description: Run test
    tasks:
    - script: go test -v -race ./...

  install:
    description: Install command into your $GOBIN dir.
    tasks:
    - name: test
      call:
        command: test

    - script: go install github.com/morikuni/ran/cmd/ran
      when:
      - test.succeeded
```

There are 2 kind of top-level key `env` and `commands`.
You can define environment variables which is used in a execution time on `env` key.
You can also define commands on `commands` key.
In this example, there are 2 command, `test` and `install`.
These are used as the sub-commands of ran command.

Each commands have a list of tasks.
Task is defined as a bash script or it calls another command.
The `install` is defined by 2 tasks.
First one is named `test` and it calls another command `test`.
Second one has no name, but it has a `when` key.
You can define dependencies on `when` key.
In this example, the second `go install` script is executed only when the first task `test` was succeeded.

Let's execute the task file.
At the first, call `help` command.

```sh
$ ran help
Usage:
  ran [command]

Available Commands:
  help          Help about any command
  install       Install command into your $GOBIN dir.
  test          Run test

Flags:
  -f, --file string        ran definition file. (default "ran.yaml")
  -h, --help               help for ran
      --log-level string   log level. (debug, info, error, discard) (default "info")

Use "ran [command] --help" for more information about a command.
```

All commands in task file is printed.
It's very useful when you don't know much about the project yet.

At next, call `test` command.

```sh
$ ran test
> go test -v -race ./...
=== RUN   TestStdLogger
--- PASS: TestStdLogger (0.00s)
=== RUN   TestTaskRunner
=== RUN   TestTaskRunner/defer
=== RUN   TestTaskRunner/no_events
=== RUN   TestTaskRunner/success
=== RUN   TestTaskRunner/error
=== RUN   TestTaskRunner/call
--- PASS: TestTaskRunner (0.02s)
    --- PASS: TestTaskRunner/defer (0.00s)
    --- PASS: TestTaskRunner/no_events (0.00s)
    --- PASS: TestTaskRunner/success (0.01s)
    --- PASS: TestTaskRunner/error (0.01s)
    --- PASS: TestTaskRunner/call (0.00s)
=== RUN   Test_EventsToParams
--- PASS: Test_EventsToParams (0.00s)
PASS
ok      github.com/morikuni/ran (cached)
?       github.com/morikuni/ran/cmd/ran [no test files]
```

`go test` script is executed.
Next, `install` command.

```sh
$ ran install
> go test -v -race ./...
=== RUN   TestStdLogger
--- PASS: TestStdLogger (0.00s)
=== RUN   TestTaskRunner
=== RUN   TestTaskRunner/defer
=== RUN   TestTaskRunner/no_events
=== RUN   TestTaskRunner/success
=== RUN   TestTaskRunner/error
=== RUN   TestTaskRunner/call
--- PASS: TestTaskRunner (0.02s)
    --- PASS: TestTaskRunner/defer (0.00s)
    --- PASS: TestTaskRunner/no_events (0.00s)
    --- PASS: TestTaskRunner/success (0.01s)
    --- PASS: TestTaskRunner/error (0.01s)
    --- PASS: TestTaskRunner/call (0.00s)
=== RUN   Test_EventsToParams
--- PASS: Test_EventsToParams (0.00s)
PASS
ok      github.com/morikuni/ran (cached)
?       github.com/morikuni/ran/cmd/ran [no test files]
> go install github.com/morikuni/ran/cmd/ran
```

`go install` script is executed after `go test` script.

## Event

There are 4 types of the event in ran.

| event | when triggered |
| :-: | :- |
| started | on the script is started |
| finished | on the script is finished (does not depends on success or fail) |
| succeeded | on the script is succeeded |
| failed | on the script is failed |
