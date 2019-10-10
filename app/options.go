package app

import (
	"github.com/cybercongress/cyberd/x/debug"
	"github.com/cybercongress/cyberd/x/rank"
)

type Options struct {
	ComputeUnit    rank.ComputeUnit
	AllowSearch    bool
	InvCheckPeriod uint

	Debug debug.Options
}
