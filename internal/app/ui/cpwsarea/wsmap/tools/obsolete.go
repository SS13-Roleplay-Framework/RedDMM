package tools

import (
	"sdmm/internal/app/ui/cpwsarea/wsmap/pmap/overlay"
	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/util"

	"github.com/rs/zerolog/log"
)

const (
	TNViewObsolete    = "View Obsolete"
	TNReplaceObsolete = "Replace Obsolete"
)

// ObsoleteInfo holds information about the currently viewed obsolete object.
type ObsoleteInfo struct {
	OriginalPath string
	OriginalVars string
	ShowPopup    bool
}

var currentObsoleteInfo ObsoleteInfo

// GetObsoleteInfo returns the current obsolete info and clears the popup flag.
func GetObsoleteInfo() ObsoleteInfo {
	info := currentObsoleteInfo
	currentObsoleteInfo.ShowPopup = false
	return info
}

// isObsoletePath checks if a path matches the configured obsolete paths.
func isObsoletePath(path string) bool {
	if ed == nil {
		return false
	}
	prefs := ed.Prefs()
	return path == prefs.Editor.ObsoleteObjectPath ||
		path == prefs.Editor.ObsoleteTurfPath ||
		path == prefs.Editor.ObsoleteAreaPath
}

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

func (t *ToolViewObsolete) process() {
	if ed == nil || cs == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	if isObsoletePath(instance.Prefab().Path()) {
		coord := instance.Coord()
		ed.OverlayPushTile(coord, overlay.ColorToolPickTileFill, overlay.ColorToolPickTileBorder)
	}
}

func (t *ToolViewObsolete) onStart(coord util.Point) {
	if ed == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	path := instance.Prefab().Path()
	if !isObsoletePath(path) {
		return
	}

	// Get original path and vars from the obsolete object
	vars := instance.Prefab().Vars()
	originalPath, _ := vars.Value("original_path")
	originalVars, _ := vars.Value("original_vars")

	ed.InstanceSelect(instance)

	currentObsoleteInfo = ObsoleteInfo{
		OriginalPath: originalPath,
		OriginalVars: originalVars,
		ShowPopup:    true,
	}
	log.Print("viewing obsolete object:", originalPath)
}

// ToolReplaceObsolete allows replacing obsolete objects with a selected type.
type ToolReplaceObsolete struct {
	tool
	replacedCount int
}

func newReplaceObsolete() *ToolReplaceObsolete {
	return &ToolReplaceObsolete{}
}

func (*ToolReplaceObsolete) Name() string {
	return TNReplaceObsolete
}

func (t *ToolReplaceObsolete) process() {
	if ed == nil || cs == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	if isObsoletePath(instance.Prefab().Path()) {
		coord := instance.Coord()
		ed.OverlayPushTile(coord, overlay.ColorToolReplaceAltTileFill, overlay.ColorToolReplaceAltTileBorder)
	}
}

func (t *ToolReplaceObsolete) onStart(coord util.Point) {
	t.onMove(coord)
}

func (t *ToolReplaceObsolete) onMove(coord util.Point) {
	if ed == nil || cs.HoverOutOfBounds() {
		return
	}

	instance := ed.HoveredInstance()
	if instance == nil {
		return
	}

	path := instance.Prefab().Path()
	if !isObsoletePath(path) {
		return
	}

	selectedPrefab, ok := ed.SelectedPrefab()
	if !ok {
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
	if tile == nil {
		return
	}

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

func (t *ToolReplaceObsolete) createReplacementPrefab(basePrefab *dmmprefab.Prefab, originalVars map[string]string) *dmmprefab.Prefab {
	newVars := basePrefab.Vars()

	for key, value := range originalVars {
		if key == "name" {
			newVars = dmvars.Set(newVars, "name", value)
			continue
		}
		if key != "original_path" && key != "original_vars" {
			newVars = dmvars.Set(newVars, key, value)
		}
	}

	return dmmprefab.New(dmmprefab.IdNone, basePrefab.Path(), newVars)
}

func (t *ToolReplaceObsolete) OnDeselect() {
	t.replacedCount = 0
}

// Direction mapping for random selection
var _relativeIndexToDir = map[int32]int{
	1: 1,  // NORTH
	2: 2,  // SOUTH
	3: 4,  // EAST
	4: 8,  // WEST
	5: 5,  // NORTHEAST
	6: 6,  // SOUTHWEST
	7: 9,  // NORTHWEST
	8: 10, // SOUTHEAST
}

