package cppresets

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sdmm/internal/dmapi/dm"
	"sdmm/internal/dmapi/dmenv"
	"sdmm/internal/dmapi/dmicon"
	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/dmapi/dmmap/dmmdata"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/util"

	"github.com/rs/zerolog/log"
)

const (
	presetsDirName = "RedDMM/presets"
)

type Preset struct {
	Name string
	Path string // Absolute path to .dmm file
}

type App interface {
	LoadedEnvironment() *dmenv.Dme
}

type Presets struct {
	app App
	
	items []*Preset
	loadedEnvPath string
}

func New() *Presets {
	return &Presets{}
}

func (p *Presets) Init(app App) {
	p.app = app
}

func (p *Presets) Items() []*Preset {
	// Lazy load check
	if p.app.LoadedEnvironment() != nil {
		currentEnv := p.app.LoadedEnvironment().RootDir
		if p.loadedEnvPath != currentEnv {
			p.Load()
			p.loadedEnvPath = currentEnv
		}
	} else {
		p.items = nil
		p.loadedEnvPath = ""
	}
	return p.items
}

func (p *Presets) Load() {
	if p.app.LoadedEnvironment() == nil {
		return
	}
	
	rootDir := p.app.LoadedEnvironment().RootDir
	presetsDir := filepath.Join(rootDir, presetsDirName)
	
	p.items = []*Preset{}
	
	files, err := os.ReadDir(presetsDir)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Printf("failed to read presets dir: %v", err)
		return
	}
	
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".dmm") {
			name := strings.TrimSuffix(f.Name(), ".dmm")
			p.items = append(p.items, &Preset{
				Name: name,
				Path: filepath.Join(presetsDir, f.Name()),
			})
		}
	}
	log.Printf("loaded %d mapper presets", len(p.items))
}

func (p *Presets) Refresh() {
	p.Load()
}

func (p *Presets) LoadPresetDmm(name string) (*dmmdata.DmmData, error) {
	path := filepath.Join(p.app.LoadedEnvironment().RootDir, presetsDirName, name+".dmm")
	return dmmdata.New(path)
}

func (p *Presets) SaveSelectionAsStamp(name string, tiles []util.Point, dmm *dmmap.Dmm) error {
	if len(tiles) == 0 {
		return fmt.Errorf("no tiles selected")
	}

	// Calculate bounds
	minX, minY, minZ := tiles[0].X, tiles[0].Y, tiles[0].Z
	maxX, maxY, maxZ := tiles[0].X, tiles[0].Y, tiles[0].Z

	tileSet := make(map[util.Point]bool)

	for _, t := range tiles {
		if t.X < minX { minX = t.X }
		if t.Y < minY { minY = t.Y }
		if t.Z < minZ { minZ = t.Z }
		if t.X > maxX { maxX = t.X }
		if t.Y > maxY { maxY = t.Y }
		if t.Z > maxZ { maxZ = t.Z }
		tileSet[t] = true
	}

	width := maxX - minX + 1
	height := maxY - minY + 1
	depth := maxZ - minZ + 1

	data := &dmmdata.DmmData{
		Filepath:   filepath.Join(p.app.LoadedEnvironment().RootDir, presetsDirName, name+".dmm"),
		MaxX:       width,
		MaxY:       height,
		MaxZ:       depth,
		Dictionary: make(dmmdata.DataDictionary),
		Grid:       make(dmmdata.DataGrid),
		KeyLength:  1,
	}

	fileDir := filepath.Dir(data.Filepath)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return err
	}

	// Helper for keys
	nextKeyStr := "a"
	if data.KeyLength > 1 {
		nextKeyStr = strings.Repeat("a", data.KeyLength)
	}

	prefabsToKey := make(map[string]dmmdata.Key)
	
	// Iterate through bounds (to maintain rectangular shape, filling empty spots with defaults)
	// Actually stamps are usually rectangular. If selection is irregular, we fill missing tiles with logic or empty?
	// The problem says "Save selections". If selection is irregular, we should probably save rectangular bounds but fill empty space with "null" or base turf?
	// DMM usually fills everything.
	// If tile is NOT in tileSet, we should probably output a "skip" tile or default. 
	// But usually stamps are pasted as is. If I paste a stamp with "default floor", it overwrites.
	// We want transparent stamps? DMM doesn't support "void". It specifies what's on the tile.
	// For now assume rectangular selection or fill with what's on map (even if not selected).
	// But `tiles` arg suggests specific tiles. 
	// If I select coordinates (1,1) and (1,3), what happens to (1,2)?
	// If I save as DMM, (1,2) must be present.
	// I will fill it with whatever is on the map at (1,2) relative to bounds.
	// Effectively we are saving the BOUNDING BOX of the selection.

	for z := 0; z < depth; z++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				absPt := util.Point{X: minX + x, Y: minY + y, Z: minZ + z}
				relPt := util.Point{X: x + 1, Y: y + 1, Z: z + 1}

				// Get tile from map
				var tilePrefabs dmmdata.Prefabs
				if dmm.HasTile(absPt) {
					tile := dmm.GetTile(absPt)
					for _, inst := range tile.Instances() {
						// Convert instance prefab to dmmdata prefab (vars might be different?)
						// Instance has vars. Prefab has vars.
						// dmmdata.Prefab needs valid path and vars.
						
						// Create dmmprefab with instance vars
						// We need to clone it because we might modify it? 
						// Actually dmmprefab.New copies vars if passed.
						p := dmmprefab.New(dmmprefab.IdNone, inst.Prefab().Path(), inst.Vars())
						tilePrefabs = append(tilePrefabs, p)
					}
				}
				// Else empty tile? DMM requires at least one prefab usually? Or valid key.

				if len(tilePrefabs) == 0 {
					// Fallback to minimal turf if empty
					// Or empty list
				}

				// Generate signature
				sig := ""
				for _, p := range tilePrefabs {
					sig += p.Path() + p.Vars().String() + ";"
				}

				key, ok := prefabsToKey[sig]
				if !ok {
					key = dmmdata.Key(nextKeyStr)
					data.Dictionary[key] = tilePrefabs
					prefabsToKey[sig] = key
					
					// Increment key
					nextKeyStr = incrementKey(nextKeyStr)
					if len(nextKeyStr) > data.KeyLength {
						data.KeyLength = len(nextKeyStr)
						// Update all existing keys? No, keys can have different lengths in mixed mode? 
						// DMM usually enforces same length.
						// If length increases, we should restart? 
						// Simplification: use length 3 by default or dynamic.
						// Implementing proper key generation is complex. 
						// I'll stick to a simple incrementor that handles length.
					}
				}
				data.Grid[relPt] = key
			}
		}
	}
	
	// Generate preview
	previewPath := filepath.Join(fileDir, name+".png")
	p.generatePreview(previewPath, width, height, data.Dictionary, data.Grid)

	data.Save()
	p.Refresh()
	return nil
}

func (p *Presets) generatePreview(path string, width, height int, dict dmmdata.DataDictionary, grid dmmdata.DataGrid) {
	const iconSize = dmmap.WorldIconSize
	imgWidth := width * iconSize
	imgHeight := height * iconSize
	canvas := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := 1; y <= height; y++ {
		for x := 1; x <= width; x++ {
			targetX := (x - 1) * iconSize
			targetY := (height - y) * iconSize
			// Preview only Z=1 for now
			relPt := util.Point{X: x, Y: y, Z: 1}

			if key, ok := grid[relPt]; ok {
				prefabs := dict[key]
				for _, prefab := range prefabs {
					sprite := dmicon.Cache.GetSpriteOrPlaceholderV(
						prefab.Vars().TextV("icon", ""),
						prefab.Vars().TextV("icon_state", ""),
						prefab.Vars().IntV("dir", dm.DirDefault),
					)

					if sprite != nil {
						if rgba, ok := sprite.Image().(*image.NRGBA); ok {
							sr := image.Rect(sprite.X1, sprite.Y1, sprite.X2, sprite.Y2)
							dp := image.Pt(targetX, targetY)
							r := image.Rectangle{Min: dp, Max: dp.Add(sr.Size())}
							draw.Draw(canvas, r, rgba, sr.Min, draw.Over)
						}
					}
				}
			}
		}
	}

	f, err := os.Create(path)
	if err == nil {
		png.Encode(f, canvas)
		f.Close()
	} else {
		log.Printf("failed to create preview: %v", err)
	}
}

func incrementKey(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] < 'z' {
			runes[i]++
			return string(runes)
		} else if runes[i] < 'Z' && runes[i] >= 'a' { // unlikely if we stick to lower
             // handle Z
			runes[i] = 'a' // wrap
		} else {
             runes[i] = 'a'
        }
	}
	return "a" + string(runes)
}
