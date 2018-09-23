package bridge

import (
	"github.com/MichaelMure/git-bug/bridge/core"
	"github.com/MichaelMure/git-bug/bridge/github"
)

// Bridges return all known bridges
func Bridges() []core.BridgeImpl {
	return []core.BridgeImpl{
		&github.Github{},
	}
}
