package core

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/pkg/errors"
)

var ErrImportNorSupported = errors.New("import is not supported")
var ErrExportNorSupported = errors.New("export is not supported")

var bridgeImpl map[string]reflect.Type

// Bridge is a wrapper around a BridgeImpl that will bind low-level
// implementation with utility code to provide high-level functions.
type Bridge struct {
	Name string
	impl BridgeImpl
	conf Configuration
}

// Register will register a new type of BridgeImpl
func Register(_type string, impl interface{}) {
	if bridgeImpl == nil {
		bridgeImpl = make(map[string]reflect.Type)
	}
	bridgeImpl[_type] = reflect.TypeOf(impl)
}

func LoadBridge(_type string, name string) (*Bridge, error) {
	implType, ok := bridgeImpl[_type]
	if !ok {
		return nil, fmt.Errorf("unknown bridge type %v", _type)
	}

	b := reflect.New(implType).Interface()

	deref := reflect.ValueOf(b).Elem().Interface()

	opp.Operations = append(opp.Operations, deref.(Operation))
}

func (b *Bridge) Configure(repo repository.RepoCommon) error {
	conf, err := b.impl.Configure(repo)
	if err != nil {
		return err
	}

	return b.storeConfig(repo, conf)
}

func (b *Bridge) storeConfig(repo repository.RepoCommon, conf Configuration) error {
	for key, val := range conf {
		storeKey := fmt.Sprintf("git-bug.%s.%s.%s", b.impl.Type(), b.Name, key)

		cmd := exec.Command("git", "config", "--replace-all", storeKey, val)
		cmd.Dir = repo.GetPath()

		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error while storing bridge configuration: %s", out)
		}
	}

	return nil
}

func (b Bridge) getConfig(repo repository.RepoCommon) (Configuration, error) {
	var err error
	if b.conf == nil {
		b.conf, err = b.loadConfig(repo)
		if err != nil {
			return nil, err
		}
	}

	return b.conf, nil
}

func (b Bridge) loadConfig(repo repository.RepoCommon) (Configuration, error) {
	key := fmt.Sprintf("git-bug.%s.%s", b.impl.Type(), b.Name)
	cmd := exec.Command("git", "config", "--get-regexp", key)
	cmd.Dir = repo.GetPath()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error while reading bridge configuration: %s", out)
	}

	lines := strings.Split(string(out), "\n")

	result := make(Configuration, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad bridge configuration: %s", line)
		}

		result[parts[0]] = parts[1]
	}

	return result, nil
}

func (b Bridge) ImportAll(repo *cache.RepoCache) error {
	importer := b.impl.Importer()
	if importer == nil {
		return ErrImportNorSupported
	}

	conf, err := b.getConfig(repo)
	if err != nil {
		return err
	}

	return b.impl.Importer().ImportAll(repo, conf)
}

func (b Bridge) Import(repo *cache.RepoCache, id string) error {
	importer := b.impl.Importer()
	if importer == nil {
		return ErrImportNorSupported
	}

	conf, err := b.getConfig(repo)
	if err != nil {
		return err
	}

	return b.impl.Importer().Import(repo, conf, id)
}

func (b Bridge) ExportAll(repo *cache.RepoCache) error {
	exporter := b.impl.Exporter()
	if exporter == nil {
		return ErrExportNorSupported
	}

	conf, err := b.getConfig(repo)
	if err != nil {
		return err
	}

	return b.impl.Exporter().ExportAll(repo, conf)
}

func (b Bridge) Export(repo *cache.RepoCache, id string) error {
	exporter := b.impl.Exporter()
	if exporter == nil {
		return ErrExportNorSupported
	}

	conf, err := b.getConfig(repo)
	if err != nil {
		return err
	}

	return b.impl.Exporter().Export(repo, conf, id)
}
