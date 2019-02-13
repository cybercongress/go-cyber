package app

import "github.com/cybercongress/cyberd/x/rank"

type Options struct {
	// main options
	ComputeUnit rank.ComputeUnit
	AllowSearch bool

	// debug options
	FailBeforeHeight int64
}
