package models

// LoaderType represents the mod loader type
type LoaderType string

const (
	LoaderTypeFabric      LoaderType = "fabric"
	LoaderTypeForgeModern LoaderType = "forge_modern"
	LoaderTypeForgeLegacy LoaderType = "forge_legacy"
	LoaderTypeNeoForge    LoaderType = "neoforge"
	LoaderTypeQuilt       LoaderType = "quilt"
)

// String returns the string representation of the loader type
func (l LoaderType) String() string {
	return string(l)
}
