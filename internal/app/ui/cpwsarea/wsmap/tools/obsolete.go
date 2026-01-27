package tools

import (
	"sdmm/internal/app/ui/dialog"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/util"
)

// ToolViewObsolete allows viewing information about obsolete objects on the map.
type ToolViewObsolete struct {
	tool
}

func newViewObsolete() *ToolViewObsolete {
	return &ToolViewObsolete{}
}

func (*ToolViewObsolete) Name() string {
	return TNViewObsolete
}

func (t *ToolViewObsolete) onStart(coord util.Point) {
	if ed == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	// Check if the instance is an obsolete object by checking its path
	path := instance.Prefab().Path()
	if isObsoletePath(path) {
		dialog.Open(dialog.NewObsoleteInfo(instance.Prefab().Vars()))
	}
}

// ToolReplaceObsolete allows replacing obsolete objects with a selected type.
type ToolReplaceObsolete struct {
	tool
}

func newReplaceObsolete() *ToolReplaceObsolete {
	return &ToolReplaceObsolete{}
}

func (*ToolReplaceObsolete) Name() string {
	return TNReplaceObsolete
}

func (t *ToolReplaceObsolete) onStart(coord util.Point) {
	if ed == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	// Check if the instance is an obsolete object
	path := instance.Prefab().Path()
	if !isObsoletePath(path) {
		return
	}

	// Get the selected prefab to replace with
	selectedPrefab, ok := ed.SelectedPrefab()
	if !ok {
		return
	}

	// Transfer valid vars from obsolete object to new object
	obsoleteVars := instance.Prefab().Vars()
	newVars := selectedPrefab.Vars() // Start with selected prefab defaults
	
	// Copy over any vars that aren't specific to the obsolete placeholder
	for _, k := range obsoleteVars.Keys() {
		// Skip internal/tracking vars
		if k == "original_type" || k == "original_name" || k == "original_vars" || 
		   k == "icon" || k == "icon_state" || k == "type" {
			continue
		}
		
		val := obsoleteVars.ValueV(k, "")
		newVars = dmvars.Set(newVars, k, val)
	}

	// Create new prefab with transferred vars
	newPrefab := dmmprefab.New(dmmprefab.IdNone, selectedPrefab.Path(), newVars)

	// Get tile to modify
	tile := ed.Dmm().GetTile(coord)
	if tile == nil {
		return
	}

	// Remove the obsolete instance and add the new one
	ed.InstanceDelete(instance)
	tile.InstancesAdd(newPrefab)
	tile.InstancesRegenerate()

	ed.UpdateCanvasByCoords([]util.Point{coord})
	ed.CommitChanges("Replace Obsolete: " + path)
}

// isObsoletePath checks if a path is one of the obsolete placeholder types
func isObsoletePath(path string) bool {
	return path == "/obj/obselete" || path == "/turf/obselete" || path == "/area/obselete"
}
