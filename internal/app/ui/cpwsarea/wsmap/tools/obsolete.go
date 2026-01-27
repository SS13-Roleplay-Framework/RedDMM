package tools

import (
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
		// TODO: Open obsolete info window with the instance data
		// For now, just log that we clicked an obsolete object
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

	// TODO: Transfer vars from obsolete object to new object
	// For now, just replace the instance with the selected prefab
	tile := ed.Dmm().GetTile(coord)
	if tile == nil {
		return
	}

	// Remove the obsolete instance and add the new one
	ed.InstanceDelete(instance)
	tile.InstancesAdd(selectedPrefab)
	tile.InstancesRegenerate()

	ed.UpdateCanvasByCoords([]util.Point{coord})
	ed.CommitChanges("Replace Obsolete: " + path)
}

// isObsoletePath checks if a path is one of the obsolete placeholder types
func isObsoletePath(path string) bool {
	return path == "/obj/obselete" || path == "/turf/obselete" || path == "/area/obselete"
}
