package tools

import (
	"math/rand"
	"strconv"

	"math/rand"
	"strconv"

	"sdmm/internal/dmapi/dm"
	"sdmm/internal/dmapi/dmenv"
	"sdmm/internal/dmapi/dmicon"
	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/util"
)

// Tool is a basic interface for tools in the panel.
type Tool interface {
	Name() string

	IgnoreBounds() bool
	Stale() bool
	AltBehaviour() bool
	setAltBehaviour(bool)

	// OnDeselect gees when the current tool is deselected.
	OnDeselect()

	// Goes every app cycle to handle stuff like pushing overlays etc.
	process()
	// Goes when user clicks on the map.
	onStart(coord util.Point)
	// Goes when user clicked and, while holding the mouse button, move the mouse.
	onMove(coord util.Point)
	// Goes when user releases the mouse button.
	onStop(coord util.Point)
}

// Tool is a basic interface for tools in the panel.
type tool struct {
	altBehaviour bool
}

func (tool) IgnoreBounds() bool {
	return false
}

func (tool) Stale() bool {
	return true
}

func (t *tool) AltBehaviour() bool {
	return t.altBehaviour
}

func (t *tool) setAltBehaviour(altBehaviour bool) {
	t.altBehaviour = altBehaviour
}

func (tool) process() {
}

//nolint:unused
func (tool) onStart(util.Point) {
}

func (tool) onMove(util.Point) {
}

func (tool) onStop(util.Point) {
}

func (tool) OnDeselect() {
}

// Direction constants for randomization
var _allDirs = []int{dm.DirNorth, dm.DirSouth, dm.DirEast, dm.DirWest, dm.DirNortheast, dm.DirSouthwest, dm.DirNorthwest, dm.DirSoutheast}

// A basic behaviour add.
// Adds object above and tile with a replacement.
// Mirrors that behaviour in the alt mode.
func (t *tool) basicPrefabAdd(tile *dmmap.Tile, prefab *dmmprefab.Prefab) {
	if !t.altBehaviour {
		if dm.IsPath(prefab.Path(), "/area") {
			tile.InstancesRemoveByPath("/area")
		} else if dm.IsPath(prefab.Path(), "/turf") {
			tile.InstancesRemoveByPath("/turf")
		}
	} else if dm.IsPath(prefab.Path(), "/obj") {
		tile.InstancesRemoveByPath("/obj")
	}

	// Apply random direction if enabled in preferences
	prefabToAdd := prefab
	if ed != nil && ed.Prefs().Editor.RandomizeDirection {
		prefabToAdd = applyRandomDirection(prefab)
	}

	tile.InstancesAdd(prefabToAdd)
	tile.InstancesRegenerate()
}

// applyRandomDirection creates a new prefab with a randomized direction
// based on the available directions in the icon
func applyRandomDirection(prefab *dmmprefab.Prefab) *dmmprefab.Prefab {
	vars := prefab.Vars()
	icon := vars.TextV("icon", "")
	iconState := vars.TextV("icon_state", "")

	// Get max directions from the icon
	state, err := dmicon.Cache.GetState(icon, iconState)
	if err != nil || state.Dirs <= 1 {
		return prefab // No directions to randomize
	}

	// Get available directions based on icon dirs
	availableDirs := _allDirs[:state.Dirs]

	// Pick a random direction
	randomDir := availableDirs[rand.Intn(len(availableDirs))]

	// Create new prefab with the random direction
	newVars := dmvars.Set(vars, "dir", strconv.Itoa(randomDir))
	return dmmprefab.New(dmmprefab.IdNone, prefab.Path(), newVars)
}

