package cli

import (
	"strconv"
	"text/template"

	"github.com/therenotomorrow/tasker/internal/domain"
)

const namespace = "tasker"

func compileTemplates() *template.Template {
	templates := template.New(namespace)

	for name, body := range map[int]string{
		notEnoughArgsTpl:      notEnoughArgsBody,
		unknownCommandTpl:     unknownCommandBody,
		invalidTaskIDTpl:      invalidTaskIDBody,
		invalidDescriptionTpl: invalidDescriptionBody,
		invalidStatusTpl:      invalidStatusBody,
		taskNotFoundTpl:       taskNotFoundBody,
		unexpectedErrorTpl:    unexpectedErrorBody,
		taskAlreadyDoneTpl:    taskAlreadyDoneBody,
		taskListIsEmptyTpl:    taskListIsEmptyBody,

		addTaskTpl:    addTaskBody,
		updateTaskTpl: updateTaskBody,
		deleteTaskTpl: deleteTaskBody,
		markTaskTpl:   markTaskBody,
		listTaskTpl:   listTaskBody,

		helpTpl: helpBody,
	} {
		_, _ = templates.New(strconv.Itoa(name)).Parse(body)
	}

	return templates
}

const (
	notEnoughArgsTpl = 1 + iota<<1
	unknownCommandTpl
	invalidTaskIDTpl
	invalidDescriptionTpl
	invalidStatusTpl
	taskNotFoundTpl
	unexpectedErrorTpl
	taskAlreadyDoneTpl
	taskListIsEmptyTpl

	notEnoughArgsBody      = `error: not enough arguments for command "{{ .Command }}"`
	unknownCommandBody     = `error: unknown command "{{ .Command }}"`
	invalidTaskIDBody      = `error: invalid "id" parameter, must be positive integer`
	invalidDescriptionBody = `error: invalid "description" parameter, must be not empty`
	invalidStatusBody      = `error: invalid "status" parameter, must be one of {{ .Statuses }}`
	taskNotFoundBody       = `error: task (ID: {{ .TaskID }}) not found`
	unexpectedErrorBody    = `error: unexpected behaviour "{{ .Error }}"`
	taskAlreadyDoneBody    = `error: cannot change status for done task`
	taskListIsEmptyBody    = `error: task list is empty`
)

func (cli *Cli) errNotEnoughArgs(command string) int {
	_ = cli.template(notEnoughArgsTpl).Execute(cli.config.Output, map[string]string{"Command": command})

	return noArgs
}

func (cli *Cli) errUnknownCommand(command string) int {
	_ = cli.template(unknownCommandTpl).Execute(cli.config.Output, map[string]string{"Command": command})

	return failure
}

func (cli *Cli) errInvalidTaskID(id string) int {
	_ = cli.template(invalidTaskIDTpl).Execute(cli.config.Output, map[string]string{"TaskID": id})

	return invalid
}

func (cli *Cli) errInvalidDescription() int {
	_ = cli.template(invalidDescriptionTpl).Execute(cli.config.Output, nil)

	return invalid
}

func (cli *Cli) errInvalidStatus(statuses []domain.Status) int {
	_ = cli.template(invalidStatusTpl).Execute(cli.config.Output, map[string]any{"Statuses": statuses})

	return invalid
}

func (cli *Cli) errTaskNotFound(id string) int {
	_ = cli.template(taskNotFoundTpl).Execute(cli.config.Output, map[string]string{"TaskID": id})

	return failure
}

func (cli *Cli) errUnexpected(err error) int {
	_ = cli.template(unexpectedErrorTpl).Execute(cli.config.Output, map[string]string{"Error": err.Error()})

	return unknown
}

func (cli *Cli) errTaskAlreadyDone() int {
	_ = cli.template(taskAlreadyDoneTpl).Execute(cli.config.Output, nil)

	return failure
}

func (cli *Cli) errTaskListIsEmpty() int {
	_ = cli.template(taskListIsEmptyTpl).Execute(cli.config.Output, nil)

	return failure
}

const (
	addTaskTpl = 2 + iota<<1
	updateTaskTpl
	deleteTaskTpl
	markTaskTpl
	listTaskTpl

	addTaskBody    = `task added successfully (ID: {{ .TaskID }})`
	updateTaskBody = `task updated successfully`
	deleteTaskBody = `task deleted successfully`
	markTaskBody   = `task status changed successfully`
	listTaskBody   = `{{ range . }}
---- id: {{ .ID }}
description | {{ .Description }}
status      | {{ .Status }}
created at  | {{ .CreatedAt.Format "02 Jan 2006 15:04:05" }}
last update | {{ .LastUpdate }} ago
{{ end -}}`
)

func (cli *Cli) template(name int) *template.Template {
	return cli.templates.Lookup(strconv.Itoa(name))
}

const (
	helpTpl = iota

	helpBody = `manage tasks with ease from the command line:
 - tasker add "description"
      add a new task with the given description
 - tasker update <id> "new description"
      update the description of an existing task by its ID
 - tasker delete <id>
      delete the task with the specified ID
 - tasker mark <id> <status>
      set a new status for the task ("todo", "progress", or "done")
 - tasker work <id>
      shortcut to mark the task as "progress"
 - tasker done <id>
      shortcut to mark the task as "done"
 - tasker list [status]
      list all tasks, if a status is provided, only tasks with that status will be shown
 - tasker help
      show this help message and exit`
)
