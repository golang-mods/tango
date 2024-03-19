package manifest

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/golang-mods/serrors"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
)

type Manifest struct {
	Tools []ManifestTool
}

type ManifestTool struct {
	ToolPath
	Version Constraints
	Tool
}

type externalManifest struct {
	Tools map[string]any `toml:"tools,omitempty"`
}

type externalEncodeManifest struct {
	Tools any `toml:"tools,omitempty"`
}

type externalManifestTool struct {
	Version Constraints `toml:"version"`
	Tool
}

var (
	ErrRequreViersionFiled     = errors.New("require version field")
	ErrEnvFieldMustStringSlice = errors.New("env field must string slice")
	ErrUnknownToolType         = errors.New("unknown tool type")
	ErrDuplicateBinaryName     = errors.New("duplicate binary name")
)

func EncodeManifest(writer io.Writer, manifest *Manifest) error {
	var external externalEncodeManifest

	external.Tools = createExternalManifestTools(manifest.Tools)

	return toml.NewEncoder(writer).Encode(&external)
}

func createExternalManifestTools(tools []ManifestTool) any {
	fields := make([]reflect.StructField, len(tools))
	values := make([]reflect.Value, len(tools))

	for i, tool := range tools {
		var reflectType reflect.Type
		var reflectValue reflect.Value
		if len(tool.Env) > 0 {
			value := externalManifestTool{Version: tool.Version, Tool: Tool{Env: tool.Env}}
			reflectType = reflect.TypeOf(value)
			reflectValue = reflect.ValueOf(value)
		} else {
			reflectType = reflect.TypeOf(tool.Version)
			reflectValue = reflect.ValueOf(tool.Version)
		}

		fields[i] = reflect.StructField{
			Name: fmt.Sprintf("V%d", i),
			Type: reflectType,
			Tag:  reflect.StructTag(fmt.Sprintf(`toml:"%s,inline"`, tool.Path)),
		}
		values[i] = reflectValue
	}

	value := reflect.New(reflect.StructOf(fields)).Elem()
	for i, v := range values {
		value.Field(i).Set(v)
	}

	return value.Addr().Interface()
}

func DecodeManifest(reader io.Reader) (*Manifest, error) {
	var external externalManifest
	if err := toml.NewDecoder(reader).Decode(&external); err != nil {
		return nil, err
	}

	var manifest Manifest
	if tools, err := createManifestTools(external.Tools); err != nil {
		return nil, err
	} else {
		manifest.Tools = tools
	}

	if duplicates := lo.FindDuplicatesBy(manifest.Tools, func(tool ManifestTool) string {
		return tool.BinaryName()
	}); len(duplicates) > 0 {
		paths := lo.Map(duplicates, func(tool ManifestTool, _ int) string { return tool.Path })
		return nil, serrors.Format("%w", ErrDuplicateBinaryName)("paths", paths)
	}

	if err := SortNamers(manifest.Tools); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func createManifestTools(externalTools map[string]any) ([]ManifestTool, error) {
	tools := make([]ManifestTool, len(externalTools))

	i := 0
	errs := lo.MapToSlice(externalTools, func(path string, fields any) error {
		tool := &tools[i]
		tool.Path = path
		i++

		switch fields := fields.(type) {
		case string:
			if err := tool.Version.UnmarshalText([]byte(fields)); err != nil {
				return err
			}

		case map[string]any:
			if version, ok := fields["version"].(string); ok {
				if err := tool.Version.UnmarshalText([]byte(version)); err != nil {
					return err
				}
			} else {
				return serrors.Format("%w", ErrRequreViersionFiled)("path", path)
			}

			if env, ok := fields["env"]; ok {
				if env, ok := env.([]any); ok {
					if env, ok := lo.FromAnySlice[string](env); ok {
						slices.Sort(env)
						tool.Env = env
					} else {
						return ErrEnvFieldMustStringSlice
					}
				} else {
					return ErrEnvFieldMustStringSlice
				}
			}

		default:
			return serrors.Format("%w", ErrUnknownToolType)("tool", fields)
		}

		return nil
	})

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return tools, nil
}
