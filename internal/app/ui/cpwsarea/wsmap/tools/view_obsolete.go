package tools

import (
	"sdmm/internal/app/ui/cpwsarea/wsmap/pmap/overlay"
	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/util"
)

// ToolViewObsolete allows viewing information about obsolete objects.
// Click on an obsolete object to see its original path and variables.
type ToolViewObsolete struct {
	tool
}

func (ToolViewObsolete) Name() string {
	return TNViewObsolete
}

func newViewObsolete() *ToolViewObsolete {
	return &ToolViewObsolete{}
}

func (t *ToolViewObsolete) process() {
	// Highlight obsolete objects when hovering
	if ed == nil || cs == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	// Check if this is an obsolete prefab
	obsConfig := dmmap.ObsoleteConfig{
		ObjectPath: ed.Prefs().Editor.ObsoleteObjectPath,
		TurfPath:   ed.Prefs().Editor.ObsoleteTurfPath,
		AreaPath:   ed.Prefs().Editor.ObsoleteAreaPath,
	}

	if dmmap.IsObsoletePrefab(instance.Prefab().Path(), obsConfig) {
		coord := instance.Coord()
		ed.OverlayPushTile(coord, overlay.ColorToolPickTileFill, overlay.ColorToolPickTileBorder)
	}
}

func (t *ToolViewObsolete) onStart(coord util.Point) {
	if ed == nil {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	// Check if this is an obsolete prefab
	obsConfig := dmmap.ObsoleteConfig{
		ObjectPath: ed.Prefs().Editor.ObsoleteObjectPath,
		TurfPath:   ed.Prefs().Editor.ObsoleteTurfPath,
		AreaPath:   ed.Prefs().Editor.ObsoleteAreaPath,
	}

	if !dmmap.IsObsoletePrefab(instance.Prefab().Path(), obsConfig) {
		return
	}

	// Get original path and vars from the obsolete object
	vars := instance.Prefab().Vars()
	originalPath, _ := vars.Value("original_path")
	originalVars, _ := vars.Value("original_vars")

	// Signal to show the info popup (will be handled by UI layer)
	// For now, just select the instance and log
	ed.InstanceSelect(instance)

	// Store info for popup display
	SetObsoleteInfo(originalPath, originalVars)
}

// ObsoleteInfo holds information about the currently viewed obsolete object.
type ObsoleteInfo struct {
	OriginalPath string
	OriginalVars string
	ShowPopup    bool
}

var currentObsoleteInfo ObsoleteInfo

// SetObsoleteInfo sets the info to display in the popup.
func SetObsoleteInfo(path, vars string) {
	currentObsoleteInfo = ObsoleteInfo{
		OriginalPath: path,
		OriginalVars: vars,
		ShowPopup:    true,
	}
}

// GetObsoleteInfo returns the current obsolete info and clears the popup flag.
func GetObsoleteInfo() ObsoleteInfo {
	info := currentObsoleteInfo
	currentObsoleteInfo.ShowPopup = false
	return info
}

// ClearObsoleteInfo clears the popup.
func ClearObsoleteInfo() {
	currentObsoleteInfo.ShowPopup = false
}
