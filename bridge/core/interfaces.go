package core

import (
	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/repository"
)

type Configuration map[string]string

type BridgeImpl interface {
	// Name return the type of the bridge (e.g.: "github")
	Type() string

	// Configure handle the user interaction and return a key/value configuration
	// for future use
	Configure(repo repository.RepoCommon) (Configuration, error)

	// Importer return an Importer implementation if the import is supported
	Importer() Importer

	// Exporter return an Exporter implementation if the export is supported
	Exporter() Exporter
}

type Importer interface {
	ImportAll(repo *cache.RepoCache, conf Configuration) error
	Import(repo *cache.RepoCache, conf Configuration, id string) error
}

type Exporter interface {
	ExportAll(repo *cache.RepoCache, conf Configuration) error
	Export(repo *cache.RepoCache, conf Configuration, id string) error
}
