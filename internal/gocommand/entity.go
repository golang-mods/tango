package gocommand

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/serrors"
)

var _ Command = (*Entity)(nil)

type Entity struct{}

func NewEntity() Command { return &Entity{} }

func (Entity) Install(
	directory string,
	pkgs []string,
	env []string,
	logDownloading func(string),
) error {
	arguments := make([]string, 1, len(pkgs)+1)
	arguments[0] = "install"
	arguments = append(arguments, pkgs...)

	command := exec.Command("go", arguments...)
	command.Env = append(append(os.Environ(), env...), "GOBIN="+directory)

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()

	if err := command.Start(); err != nil {
		return err
	}

	var errorBuffer errorBuffer
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		const label = "go: downloading "
		text := scanner.Text()
		if strings.HasPrefix(text, label) {
			logDownloading(text[len(label):])
		} else {
			fmt.Fprintln(&errorBuffer, text)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := command.Wait(); err != nil {
		return errors.Join(err, errorBuffer.error())
	}

	return nil
}

var ErrBadPath = errors.New("bad path")

func (Entity) Versions(pkg string) (*Versions, error) {
	for current := pkg; current != "."; current = path.Dir(current) {
		internal, err := versions(current)
		if err != nil {
			if errors.Is(err, errNotFound) {
				continue
			} else {
				return nil, err
			}
		}

		if len(internal.Versions) > 0 {
			versions, err := toVersions(internal.Versions)
			if err != nil {
				return nil, err
			}
			return &Versions{Path: internal.Path, Versions: versions}, nil
		}
	}

	return nil, serrors.Format("%w", ErrBadPath)("path", pkg)
}

type internalVersions struct {
	Path     string
	Versions []string
}

func versions(pkg string) (*internalVersions, error) {
	var errorBuffer versionsErrorBuffer

	command := exec.Command("go", "list", "-json", "-m", "-versions", pkg)
	command.Stderr = &errorBuffer

	reader, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if err := command.Start(); err != nil {
		return nil, err
	}

	var versions internalVersions
	decodeErr := json.NewDecoder(reader).Decode(&versions)

	if err := errors.Join(command.Wait(), errorBuffer.error()); err != nil {
		return nil, err
	}
	if decodeErr != nil {
		return nil, decodeErr
	}

	return &versions, nil
}

func toVersions(versions []string) ([]*semver.Version, error) {
	return exerrors.Map(versions, func(version string, _ int) (*semver.Version, error) {
		return semver.NewVersion(version)
	})
}

var ErrParse = errors.New("parse error")

func (Entity) BinaryVersion(name string) (*BinaryVersion, error) {
	var errorBuffer errorBuffer

	command := exec.Command("go", "version", "-m", name)
	command.Stderr = &errorBuffer

	reader, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if err := command.Start(); err != nil {
		return nil, err
	}

	var version BinaryVersion

	scanner := bufio.NewScanner(reader)
	if scanner.Scan() {
		const separator = ": go"

		text := scanner.Text()
		index := strings.LastIndex(text, separator)
		if index == -1 {
			return nil, serrors.Format("%w", ErrParse)("text", text)
		}

		goVersionText := text[index+len(separator):]
		goVersion, err := semver.NewVersion(goVersionText)
		if err != nil {
			return nil, serrors.Format("%w", err)("version", goVersionText)
		}

		version.File = text[:index]
		version.GoVersion = goVersion
	}

	for scanner.Scan() {
		switch parts := strings.Split(scanner.Text(), "\t")[1:]; parts[0] {
		case "path":
			version.Path = parts[1]
		case "mod":
			module, err := sliceToModule(parts[1:])
			if err != nil {
				return nil, err
			}
			version.Module = module
		case "dep":
			dependence, err := sliceToModule(parts[1:])
			if err != nil {
				return nil, err
			}
			version.Dependencies = append(version.Dependencies, dependence)
		case "build":
			version.Builds = append(version.Builds, parts[1])
		}
	}

	if err := errors.Join(scanner.Err(), command.Wait(), errorBuffer.error()); err != nil {
		return nil, err
	}

	slices.Sort(version.Builds)

	return &version, nil
}

func sliceToModule(parts []string) (*Module, error) {
	version, err := semver.NewVersion(parts[1])
	if err != nil {
		return nil, serrors.Format("%w", err)("version", parts[1])
	}

	return &Module{
		Name:    parts[0],
		Version: version,
		Hash:    parts[2],
	}, nil
}
