package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/therenotomorrow/tasker/internal/domain"
	"github.com/therenotomorrow/tasker/internal/usecases"
)

func (cli *Cli) Add(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("add")
	}

	description := args[0]
	task, err := cli.use.AddTask(ctx, description)

	switch {
	case errors.Is(err, domain.ErrEmptyDescription):
		return cli.errInvalidDescription()
	case err != nil:
		return cli.errUnexpected(err)
	}

	_ = cli.template(addTaskTpl).Execute(os.Stdout, map[string]uint64{"TaskID": task.ID})

	return success
}

func (cli *Cli) Update(ctx context.Context, args []string) int {
	if len(args) < twoArgs {
		return cli.errNotEnoughArgs("update")
	}

	taskID, newDesc := args[0], args[1]
	_, err := cli.use.UpdateTask(ctx, taskID, newDesc)

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

func (cli *Cli) Delete(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("delete")
	}

	taskID := args[0]
	err := cli.use.DeleteTask(ctx, taskID)

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

func (cli *Cli) Mark(ctx context.Context, args []string) int {
	if len(args) < twoArgs {
		return cli.errNotEnoughArgs("mark")
	}

	taskID, status := args[0], args[1]
	_, err := cli.use.MarkTask(ctx, taskID, status)

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

func (cli *Cli) Work(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("work")
	}

	args = append(args, string(domain.StatusProgress))

	return cli.Mark(ctx, args)
}

func (cli *Cli) Done(ctx context.Context, args []string) int {
	if len(args) < oneArg {
		return cli.errNotEnoughArgs("done")
	}

	args = append(args, string(domain.StatusDone))

	return cli.Mark(ctx, args)
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

func (cli *Cli) List(ctx context.Context, args []string) int {
	status := ""
	if len(args) > 0 {
		status = args[0]
	}

	list, err := cli.use.ListTasks(ctx, usecases.ListParams{Status: status})

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
		return b.UpdatedAt.Compare(a.UpdatedAt)
	})

	_ = cli.template(listTaskTpl).Execute(os.Stdout, views)

	return success
}

func (cli *Cli) Help() int {
	_ = cli.template(helpTpl).Execute(os.Stdout, nil)

	return success
}

func (cli *Cli) Usage() int {
	_ = cli.Help()

	return noArgs
}

func (cli *Cli) Unknown(command string) int {
	return cli.errUnknownCommand(command)
}
