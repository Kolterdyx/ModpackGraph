package loaders

import (
	"encoding/json"
)

type FabricV1Entrypoint struct {
	Adapter string `json:"adapter"`
	Value   string `json:"value"`
}

func (fv1e *FabricV1Entrypoint) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		fv1e.Value = asString
		fv1e.Adapter = "default"
		return nil
	}
	var asObject struct {
		Adapter string `json:"adapter"`
		Value   string `json:"value"`
	}
	if err := json.Unmarshal(data, &asObject); err != nil {
		return err
	}
	fv1e.Adapter = asObject.Adapter
	fv1e.Value = asObject.Value
	return nil
}

type FabricV1JarEntry struct {
	Path string `json:"path"`
}

type FabricV1MixinEntry struct {
	Config      string              `json:"config"`
	Environment FabricV1Environment `json:"environment"`
}

func (fvmx *FabricV1MixinEntry) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		fvmx.Config = asString
		fvmx.Environment = FabricV1EnvironmentUniversal
		return nil
	}
	var asObject struct {
		Config      string              `json:"config"`
		Environment FabricV1Environment `json:"environment"`
	}
	if err := json.Unmarshal(data, &asObject); err != nil {
		return err
	}
	fvmx.Config = asObject.Config
	fvmx.Environment = asObject.Environment
	return nil
}

type FabricV1Environment string

const (
	FabricV1EnvironmentClient    FabricV1Environment = "client"
	FabricV1EnvironmentServer    FabricV1Environment = "server"
	FabricV1EnvironmentUniversal FabricV1Environment = "*"
)

var AllFabricV1Environments = []FabricV1Environment{
	FabricV1EnvironmentClient,
	FabricV1EnvironmentServer,
	FabricV1EnvironmentUniversal,
}

func (fve *FabricV1Environment) UnmarshalJSON(data []byte) error {
	var envStr string
	if err := json.Unmarshal(data, &envStr); err != nil {
		return err
	}
	switch env := FabricV1Environment(envStr); env {
	case FabricV1EnvironmentClient, FabricV1EnvironmentServer, FabricV1EnvironmentUniversal:
		*fve = env
	default:
		*fve = FabricV1EnvironmentUniversal
	}
	return nil
}

type FabricV1VersionMatcher []string

func (fvvm *FabricV1VersionMatcher) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*fvvm = FabricV1VersionMatcher{asString}
		return nil
	}
	var asArray []string
	if err := json.Unmarshal(data, &asArray); err != nil {
		return err
	}
	*fvvm = asArray
	return nil
}

// FabricAuthor represents an author entry
type FabricV1Person struct {
	Name    string            `json:"name"`
	Contact map[string]string `json:"contact"`
}

func (fa *FabricV1Person) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		fa.Name = asString
		fa.Contact = make(map[string]string)
		return nil
	}
	var asObject struct {
		Name    string            `json:"name"`
		Contact map[string]string `json:"contact"`
	}
	if err := json.Unmarshal(data, &asObject); err != nil {
		return err
	}
	fa.Name = asObject.Name
	fa.Contact = asObject.Contact
	return nil
}
