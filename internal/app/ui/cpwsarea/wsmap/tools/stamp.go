package tools

import (
	"sdmm/internal/app/ui/cpwsarea/wsmap/pmap/overlay"
	"sdmm/internal/dmapi/dm"
	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/dmapi/dmmap/dmmdata"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/util"

	"github.com/rs/zerolog/log"
)

type ToolStamp struct {
	tool

	stamp *dmmdata.DmmData

	editedTiles map[util.Point]bool
}

func (ToolStamp) Name() string {
	return TNStamp
}

func newStamp() *ToolStamp {
	return &ToolStamp{
		editedTiles: make(map[util.Point]bool),
	}
}

func (t *ToolStamp) SetStamp(data *dmmdata.DmmData) {
	t.stamp = data
}

func (t *ToolStamp) process() {
	if t.stamp == nil {
		return
	}
	if cs == nil || cs.HoverOutOfBounds() {
		return
	}

	hover := cs.HoveredTile()
	maxX := t.stamp.MaxX
	maxY := t.stamp.MaxY

	endX := hover.X + maxX - 1
	endY := hover.Y + maxY - 1

	bounds := util.Bounds{
		Min: hover,
		Max: util.Point{X: endX, Y: endY, Z: hover.Z},
	}
	ed.OverlayPushArea(bounds, overlay.ColorToolAddTileFill, overlay.ColorToolAddTileBorder)
}

func (t *ToolStamp) onStart(coord util.Point) {
	t.onMove(coord)
}

func (t *ToolStamp) onMove(coord util.Point) {
	if t.stamp == nil || t.editedTiles[coord] {
		return
	}

	t.editedTiles[coord] = true
	t.placeStamp(coord)
}

func (t *ToolStamp) onStop(util.Point) {
	if len(t.editedTiles) != 0 {
		t.editedTiles = make(map[util.Point]bool, len(t.editedTiles))
		go ed.CommitChanges("Place Stamp")
	}
}

func (t *ToolStamp) placeStamp(origin util.Point) {
	d := t.stamp
	var changes []util.Point

	for z := 0; z < d.MaxZ; z++ {
		for y := 0; y < d.MaxY; y++ {
			for x := 0; x < d.MaxX; x++ {
				relPt := util.Point{X: x + 1, Y: y + 1, Z: z + 1}
				targetPt := util.Point{X: origin.X + x, Y: origin.Y + y, Z: origin.Z + z}

				if !ed.Dmm().HasTile(targetPt) {
					continue
				}

				if key, ok := d.Grid[relPt]; ok {
					prefabs, ok := d.Dictionary[key]
					if ok {
						tile := ed.Dmm().GetTile(targetPt)
						for _, p := range prefabs {
							prefabToAdd := p
							// Check if object is obsolete (missing from environment)
							if ed.LoadedEnvironment().Objects[p.Path()] == nil {
								var replacePath string
								if dm.IsPath(p.Path(), "/turf") {
									replacePath = ed.Prefs().Editor.ObsoleteTurfPath
								} else if dm.IsPath(p.Path(), "/area") {
									replacePath = ed.Prefs().Editor.ObsoleteAreaPath
								} else {
									replacePath = ed.Prefs().Editor.ObsoleteObjectPath
								}

								if replacePath != "" {
									log.Printf("Replacing stamp obsolete: %s -> %s", p.Path(), replacePath)
									vars := dmvars.New()
									vars.CopyFrom(p.Vars())
									if !vars.Has("original_type") {
										vars.Set("original_type", p.Path())
									}
									prefabToAdd = dmmprefab.New(dmmprefab.IdNone, replacePath, vars)
								}
							}

							t.basicPrefabAdd(tile, prefabToAdd)
						}
						changes = append(changes, targetPt)
					}
				}
			}
		}
	}
	ed.UpdateCanvasByCoords(changes)
}
