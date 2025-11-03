package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/shyim/go-version"

	"github.com/onlishop/onlishop-cli/internal/validation"
)

type OnlishopBundle struct {
	path     string
	Composer onlishopBundleComposerJson
	config   *Config
}

func newOnlishopBundle(path string) (*OnlishopBundle, error) {
	composerJsonFile := fmt.Sprintf("%s/composer.json", path)
	if _, err := os.Stat(composerJsonFile); err != nil {
		return nil, err
	}

	jsonFile, err := os.ReadFile(composerJsonFile)
	if err != nil {
		return nil, fmt.Errorf("newOnlishopBundle: %v", err)
	}

	var composerJson onlishopBundleComposerJson
	err = json.Unmarshal(jsonFile, &composerJson)
	if err != nil {
		return nil, fmt.Errorf("newOnlishopBundle: %v", err)
	}

	if composerJson.Type != "onlishop-bundle" {
		return nil, fmt.Errorf("newOnlishopBundle: composer.json type is not onlishop-bundle")
	}

	if composerJson.Extra.BundleName == "" {
		return nil, fmt.Errorf("composer.json does not contain onlishop-bundle-name in extra")
	}

	cfg, err := readExtensionConfig(path)
	if err != nil {
		return nil, fmt.Errorf("newOnlishopBundle: %v", err)
	}

	extension := OnlishopBundle{
		Composer: composerJson,
		path:     path,
		config:   cfg,
	}

	return &extension, nil
}

type composerAutoload struct {
	Psr4 map[string]string `json:"psr-4"`
}

type onlishopBundleComposerJson struct {
	Name     string                          `json:"name"`
	Type     string                          `json:"type"`
	License  string                          `json:"license"`
	Version  string                          `json:"version"`
	Require  map[string]string               `json:"require"`
	Extra    onlishopBundleComposerJsonExtra `json:"extra"`
	Suggest  map[string]string               `json:"suggest"`
	Autoload composerAutoload                `json:"autoload"`
}

type onlishopBundleComposerJsonExtra struct {
	BundleName string `json:"onlishop-bundle-name"`
}

func (p OnlishopBundle) GetComposerName() (string, error) {
	return p.Composer.Name, nil
}

// GetRootDir returns the src directory of the bundle.
func (p OnlishopBundle) GetRootDir() string {
	return path.Join(p.path, "src")
}

func (p OnlishopBundle) GetSourceDirs() []string {
	var result []string

	for _, val := range p.Composer.Autoload.Psr4 {
		result = append(result, path.Join(p.path, val))
	}

	return result
}

// GetResourcesDir returns the resources directory of the onlishop bundle.
func (p OnlishopBundle) GetResourcesDir() string {
	return path.Join(p.GetRootDir(), "Resources")
}

func (p OnlishopBundle) GetResourcesDirs() []string {
	var result []string

	for _, val := range p.GetSourceDirs() {
		result = append(result, path.Join(val, "Resources"))
	}

	return result
}

func (p OnlishopBundle) GetName() (string, error) {
	return p.Composer.Extra.BundleName, nil
}

func (p OnlishopBundle) GetExtensionConfig() *Config {
	return p.config
}

func (p OnlishopBundle) GetOnlishopVersionConstraint() (*version.Constraints, error) {
	if p.config != nil && p.config.Build.OnlishopVersionConstraint != "" {
		constraint, err := version.NewConstraint(p.config.Build.OnlishopVersionConstraint)
		if err != nil {
			return nil, err
		}

		return &constraint, nil
	}

	onlishopConstraintString, ok := p.Composer.Require["onlishop/core"]

	if !ok {
		return nil, fmt.Errorf("require.onlishop/core is required")
	}

	onlishopConstraint, err := version.NewConstraint(onlishopConstraintString)
	if err != nil {
		return nil, err
	}

	return &onlishopConstraint, err
}

func (OnlishopBundle) GetType() string {
	return TypeOnlishopBundle
}

func (p OnlishopBundle) GetVersion() (*version.Version, error) {
	return version.NewVersion(p.Composer.Version)
}

func (p OnlishopBundle) GetChangelog() (*ExtensionChangelog, error) {
	return parseExtensionMarkdownChangelog(p)
}

func (p OnlishopBundle) GetLicense() (string, error) {
	return p.Composer.License, nil
}

func (p OnlishopBundle) GetPath() string {
	return p.path
}

func (p OnlishopBundle) GetIconPath() string {
	return ""
}

func (p OnlishopBundle) GetMetaData() *extensionMetadata {
	return &extensionMetadata{
		Label: extensionTranslated{
			Chinese: "FALLBACK",
			English: "FALLBACK",
		},
		Description: extensionTranslated{
			Chinese: "FALLBACK",
			English: "FALLBACK",
		},
	}
}

func (p OnlishopBundle) Validate(c context.Context, check validation.Check) {
	// OnlishopBundle validation is currently empty but signature updated to match interface
}
