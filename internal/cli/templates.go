package cli

import (
	"strconv"
	"text/template"
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
		_, _ = templates.New(strconv.Itoa(name)).Parse(body + "\n")
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

	notEnoughArgsBody      = `Error: not enough arguments for command "{{ .Command }}".`
	unknownCommandBody     = `Error: unknown command "{{ .Command }}".`
	invalidTaskIDBody      = `Error: invalid "id" parameter, must be positive integer.`
	invalidDescriptionBody = `Error: invalid "description" parameter, must be not empty.`
	invalidStatusBody      = `Error: invalid "status" parameter, must be one of available.`
	taskNotFoundBody       = `Error: task (ID: {{ .TaskID }}) not found.`
	unexpectedErrorBody    = `Error: unexpected behaviour "{{ .Error }}".`
	taskAlreadyDoneBody    = `Error: cannot change status for done task.`
	taskListIsEmptyBody    = `Error: task list is empty.`
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

func (cli *Cli) errInvalidStatus() int {
	_ = cli.template(invalidStatusTpl).Execute(cli.config.Output, nil)

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

	addTaskBody    = `Task added successfully (ID: {{ .TaskID }}).`
	updateTaskBody = `Task updated successfully.`
	deleteTaskBody = `Task deleted successfully.`
	markTaskBody   = `Task status changed successfully.`
	listTaskBody   = `{{- range $i, $e := . -}}{{- if $i }}
{{ end }}---- ID: {{ $e.ID }}
Description | {{ $e.Description }}
Status      | {{ $e.Status }}
Created At  | {{ $e.CreatedAt.Format "02 Jan 06 15:04:05" }}
Last Update | {{ $e.LastUpdate }} ago
{{ end -}}`
)

func (cli *Cli) template(name int) *template.Template {
	return cli.templates.Lookup(strconv.Itoa(name))
}

const (
	helpTpl = iota

	helpBody = `Manage your tasks with ease from the command line. Usage:
 - tasker add "description"
      Add a new task with the given description.
 - tasker update <id> "new description"
      Update the description of an existing task by its ID.
 - tasker delete <id>
      Delete the task with the specified ID.
 - tasker mark <id> <status>
      Set a new status for the task ("todo", "progress", or "done").
 - tasker work <id>
      Shortcut to mark the task as "progress".
 - tasker done <id>
      Shortcut to mark the task as "done".
 - tasker list [status]
      List all tasks. If a status is provided, only tasks with that status will be shown.
 - tasker help
      Show this help message.`
)
