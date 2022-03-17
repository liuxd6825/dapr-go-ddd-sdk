package ddd

import "context"

type QueryAction func(ctx context.Context) error

func QueryLog(ctx context.Context, action QueryAction) error {
	err := action(ctx)
	return err
}

func CommandLog(ctx context.Context, action QueryAction) error {
	err := action(ctx)
	return err
}
