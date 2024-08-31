package port

import (
	"context"
)

type (
	RecordReport       func(ctx context.Context, payload []byte) error
	CleanUpWithArchive func(ctx context.Context) error
)
