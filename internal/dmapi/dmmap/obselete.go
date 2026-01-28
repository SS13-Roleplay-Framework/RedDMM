package dmmap

import (
	"strconv"
	"strings"

	"sdmm/internal/dmapi/dm"
	"sdmm/internal/dmapi/dmenv"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"

	"github.com/rs/zerolog/log"
)

// ObsoleteConfig holds configuration for obsolete object replacement.
type ObsoleteConfig struct {
	ObjectPath string // Path to use for unknown /obj and /mob types
	TurfPath   string // Path to use for unknown /turf types
	AreaPath   string // Path to use for unknown /area types
}

// DefaultObsoleteConfig returns a default configuration (no replacement, just discard).
func DefaultObsoleteConfig() ObsoleteConfig {
	return ObsoleteConfig{}
}

// CreateObsoletePrefab creates an obsolete prefab from an unknown prefab.
// It stores the original path and variables in the obsolete object.
func CreateObsoletePrefab(dme *dmenv.Dme, originalPrefab *dmmprefab.Prefab, config ObsoleteConfig) *dmmprefab.Prefab {
	originalPath := originalPrefab.Path()

	// Determine which obsolete type to use based on the original path
	var obsoletePath string
	if dm.IsPath(originalPath, "/turf") {
		obsoletePath = config.TurfPath
	} else if dm.IsPath(originalPath, "/area") {
		obsoletePath = config.AreaPath
	} else {
		// /obj, /mob, and anything else go to object path
		obsoletePath = config.ObjectPath
	}

	// If no obsolete path configured, return nil (discard as before)
	if obsoletePath == "" {
		return nil
	}

	// Check if the obsolete path exists in the environment
	if _, ok := dme.Objects[obsoletePath]; !ok {
		log.Print("obsolete path not found in environment:", obsoletePath)
		return nil
	}

	// Serialize original variables to a string
	originalVars := serializeVars(originalPrefab.Vars())

	// Create the obsolete prefab with original data stored
	vars := &dmvars.MutableVariables{}
	vars.Put("original_path", strconv.Quote(originalPath))
	if originalVars != "" {
		vars.Put("original_vars", strconv.Quote(originalVars))
	}

	// Link to parent vars from environment
	envObj := dme.Objects[obsoletePath]
	immutableVars := vars.ToImmutable()
	immutableVars.LinkParent(envObj.Vars)

	return dmmprefab.New(dmmprefab.IdNone, obsoletePath, immutableVars)
}

// serializeVars converts variables to a semicolon-separated string for storage.
func serializeVars(vars *dmvars.Variables) string {
	if vars == nil {
		return ""
	}

	var parts []string
	for _, key := range vars.Iterate() {
		value, _ := vars.Value(key)
		parts = append(parts, key+"="+value)
	}
	return strings.Join(parts, ";")
}

// ParseOriginalVars parses the original_vars string back into key-value pairs.
func ParseOriginalVars(originalVars string) map[string]string {
	result := make(map[string]string)
	if originalVars == "" {
		return result
	}

	// Handle quoted string
	if strings.HasPrefix(originalVars, "\"") && strings.HasSuffix(originalVars, "\"") {
		originalVars = originalVars[1 : len(originalVars)-1]
	}

	pairs := strings.Split(originalVars, ";")
	for _, pair := range pairs {
		if idx := strings.Index(pair, "="); idx > 0 {
			key := pair[:idx]
			value := pair[idx+1:]
			result[key] = value
		}
	}
	return result
}

// IsObsoletePrefab checks if a prefab is an obsolete placeholder.
func IsObsoletePrefab(prefabPath string, config ObsoleteConfig) bool {
	return prefabPath == config.ObjectPath ||
		prefabPath == config.TurfPath ||
		prefabPath == config.AreaPath
}
