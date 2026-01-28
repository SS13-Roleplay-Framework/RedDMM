package tools

import (
	"strconv"

	"sdmm/internal/app/ui/cpwsarea/wsmap/pmap/overlay"
	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/util"

	"github.com/rs/zerolog/log"
)

// ToolReplaceObsolete allows replacing obsolete objects with new types.
// Select a type in the Environment tree, then click obsolete objects to replace them.
// Variables from the original object are transferred to the new object when possible.
type ToolReplaceObsolete struct {
	tool

	replacedCount int
}

func (ToolReplaceObsolete) Name() string {
	return TNReplaceObsolete
}

func newReplaceObsolete() *ToolReplaceObsolete {
	return &ToolReplaceObsolete{}
}

func (t *ToolReplaceObsolete) process() {
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
		ed.OverlayPushTile(coord, overlay.ColorToolReplaceAltTileFill, overlay.ColorToolReplaceAltTileBorder)
	}
}

func (t *ToolReplaceObsolete) onStart(coord util.Point) {
	t.onMove(coord)
}

func (t *ToolReplaceObsolete) onMove(coord util.Point) {
	if ed == nil {
		return
	}

	// Get the selected prefab to replace with
	selectedPrefab, ok := ed.SelectedPrefab()
	if !ok {
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

	// Get original variables from the obsolete object
	vars := instance.Prefab().Vars()
	originalVarsStr, _ := vars.Value("original_vars")
	originalVars := dmmap.ParseOriginalVars(originalVarsStr)

	// Create new prefab with transferred variables
	newPrefab := t.createReplacementPrefab(selectedPrefab, originalVars)

	// Replace the instance
	tile := ed.Dmm().GetTile(coord)
	tile.InstancesRemove(instance)
	tile.InstancesAdd(dmmap.PrefabStorage.Put(newPrefab))
	tile.InstancesRegenerate()

	ed.UpdateCanvasByCoords([]util.Point{coord})
	t.replacedCount++
}

func (t *ToolReplaceObsolete) onStop(coord util.Point) {
	if t.replacedCount > 0 {
		log.Print("replaced obsolete objects:", t.replacedCount)
		go ed.CommitChanges("Replace Obsolete")
		t.replacedCount = 0
	}
}

// createReplacementPrefab creates a new prefab transferring applicable variables.
func (t *ToolReplaceObsolete) createReplacementPrefab(basePrefab *dmmprefab.Prefab, originalVars map[string]string) *dmmprefab.Prefab {
	newVars := basePrefab.Vars()

	// Transfer variables that exist on the new type
	for key, value := range originalVars {
		// Special handling for name - transfer from original_name if present
		if key == "name" {
			newVars = dmvars.Set(newVars, "name", value)
			continue
		}

		// For other vars, check if they exist on the base prefab's type
		// We always transfer the value - the sanitize on save will clean up if needed
		if key != "original_path" && key != "original_vars" {
			newVars = dmvars.Set(newVars, key, value)
		}
	}

	return dmmprefab.New(dmmprefab.IdNone, basePrefab.Path(), newVars)
}

func (t *ToolReplaceObsolete) OnDeselect() {
	t.replacedCount = 0
}
