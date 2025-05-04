package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"tasker/internal/domain"
	"tasker/internal/usecases"
	"time"
)

func (cli *Cli) add(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("add")
	}

	description := args[0]
	task, err := cli.useCases.AddTask(ctx, description)

	switch {
	case errors.Is(err, domain.ErrEmptyDescription):
		return cli.errInvalidDescription()
	case err != nil:
		return cli.errUnexpected(err)
	}

	_ = cli.template(addTaskTpl).Execute(os.Stdout, map[string]uint64{"TaskID": task.ID})

	return success
}

func (cli *Cli) update(ctx context.Context, args []string) int {
	if len(args) < twoArgs {
		return cli.errNotEnoughArgs("update")
	}

	taskID, newDesc := args[0], args[1]
	_, err := cli.useCases.UpdateTask(ctx, taskID, newDesc)

	switch {
	case errors.Is(err, domain.ErrEmptyDescription):
		return cli.errInvalidDescription()
	case errors.Is(err, domain.ErrTaskNotFound):
		return cli.errTaskNotFound(taskID)
	case errors.Is(err, domain.ErrInvalidTaskID):
		return cli.errInvalidTaskID(taskID)
	case err != nil:
		return cli.errUnexpected(err)
	}

	_ = cli.template(updateTaskTpl).Execute(os.Stdout, nil)

	return success
}

func (cli *Cli) delete(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("delete")
	}

	taskID := args[0]
	err := cli.useCases.DeleteTask(ctx, taskID)

	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		return cli.errTaskNotFound(taskID)
	case errors.Is(err, domain.ErrInvalidTaskID):
		return cli.errInvalidTaskID(taskID)
	case err != nil:
		return cli.errUnexpected(err)
	}

	_ = cli.template(deleteTaskTpl).Execute(os.Stdout, nil)

	return success
}

func (cli *Cli) mark(ctx context.Context, args []string) int {
	if len(args) < twoArgs {
		return cli.errNotEnoughArgs("mark")
	}

	taskID, status := args[0], args[1]
	_, err := cli.useCases.MarkTask(ctx, taskID, status)

	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		return cli.errTaskNotFound(taskID)
	case errors.Is(err, domain.ErrInvalidTaskID):
		return cli.errInvalidTaskID(taskID)
	case errors.Is(err, domain.ErrTaskAlreadyDone):
		return cli.errTaskAlreadyDone()
	case errors.Is(err, domain.ErrInvalidStatus):
		return cli.errInvalidStatus()
	case err != nil:
		return cli.errUnexpected(err)
	}

	_ = cli.template(markTaskTpl).Execute(os.Stdout, nil)

	return success
}

func (cli *Cli) work(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("work")
	}

	args = append(args, string(domain.StatusProgress))

	return cli.mark(ctx, args)
}

func (cli *Cli) done(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("done")
	}

	args = append(args, string(domain.StatusDone))

	return cli.mark(ctx, args)
}

type listView struct {
	ID          uint64
	Description string
	Status      domain.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastUpdate  string
}

func lastUpdateString(t time.Time) string {
	var lastUpdate string

	const oneDay = 24

	switch duration := time.Since(t); {
	case duration.Hours() < 1:
		lastUpdate = fmt.Sprintf("%d minute(s)", int(duration.Minutes()))
	case duration.Hours() < oneDay:
		lastUpdate = fmt.Sprintf("%d hour(s)", int(duration.Hours()))
	default:
		lastUpdate = fmt.Sprintf("%d day(s)", int(duration.Hours()/oneDay))
	}

	return lastUpdate
}

func (cli *Cli) list(ctx context.Context, args []string) int {
	status := ""
	if len(args) > 0 {
		status = args[0]
	}

	list, err := cli.useCases.ListTasks(ctx, usecases.ListParams{Status: status})

	switch {
	case errors.Is(err, domain.ErrInvalidStatus):
		return cli.errInvalidStatus()
	case errors.Is(err, domain.ErrEmptyTasks):
		return cli.errTaskListIsEmpty()
	case err != nil:
		return cli.errUnexpected(err)
	}

	views := make([]listView, len(list))

	for idx, task := range list {
		views[idx] = listView{
			ID:          task.ID,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			LastUpdate:  lastUpdateString(task.UpdatedAt),
		}
	}

	slices.SortStableFunc(views, func(a, b listView) int {
		return -a.UpdatedAt.Compare(b.UpdatedAt)
	})

	_ = cli.template(listTaskTpl).Execute(os.Stdout, views)

	return success
}

func (cli *Cli) help() int {
	_ = cli.template(helpTpl).Execute(os.Stdout, nil)

	return success
}

func (cli *Cli) usage() int {
	_ = cli.help()

	return noArgs
}

func (cli *Cli) unknown(command string) int {
	return cli.errUnknownCommand(command)
}
