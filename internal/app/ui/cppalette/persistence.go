package cppalette

import (
	"encoding/json"
	"os"
	"path/filepath"

	"sdmm/internal/dmapi/dmmap"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"github.com/rs/zerolog/log"
)

const (
	paletteDirName  = "RedDMM"
	paletteFileName = "palette.json"
)

type categoryData struct {
	Name    string   `json:"name"`
	Prefabs []string `json:"prefabs"`
}

type paletteData struct {
	Categories []categoryData `json:"categories"`
}

func (p *Palette) Load() {
	if !p.app.HasLoadedEnvironment() {
		return
	}

	rootDir := p.app.LoadedEnvironment().RootDir
	path := filepath.Join(rootDir, paletteDirName, paletteFileName)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return // No palette yet
	}
	if err != nil {
		log.Printf("failed to read palette file: %v", err)
		return
	}

	var pData paletteData
	if err := json.Unmarshal(data, &pData); err != nil {
		log.Printf("failed to unmarshal palette data: %v", err)
		return
	}

	p.categories = make([]*Category, 0, len(pData.Categories))
	for _, catData := range pData.Categories {
		cat := &Category{
			Name:    catData.Name,
			Prefabs: make([]*dmmprefab.Prefab, 0, len(catData.Prefabs)),
		}

		for _, path := range catData.Prefabs {
			// Try to find existing prefab or create one from path
			// For now, we look up in environment objects
			if obj, ok := p.app.LoadedEnvironment().Objects[path]; ok {
				// We need to create a prefab from the object
				// This simplifies things, but we lose specific variable tweaks if we only store paths
				// For a full palette we might need to store full prefab data (vars)
				// But task said "prefabs", usually meaning the base type.
				
				// Using dmmap.PrefabStorage to get/create a clean prefab for the path
				// We might need to construct it if not in storage
				
				// Ideally we should store the serialized prefab if it has custom vars.
				// For this first pass, we assume path-only references.
				
				// Creating a fresh prefab from the environment object path
				prefab := dmmprefab.New(dmmprefab.IdNone, path, obj.Vars)
				cat.Prefabs = append(cat.Prefabs, prefab)
			} else {
				// Unknown path, maybe create a placeholder?
				// Or create a minimal prefab
				prefab := dmmprefab.New(dmmprefab.IdNone, path, nil)
				cat.Prefabs = append(cat.Prefabs, prefab)
			}
		}
		p.categories = append(p.categories, cat)
	}
	log.Print("palette loaded from:", path)
}

func (p *Palette) Save() {
	if !p.app.HasLoadedEnvironment() {
		return
	}

	rootDir := p.app.LoadedEnvironment().RootDir
	dirPath := filepath.Join(rootDir, paletteDirName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("failed to create palette directory: %v", err)
		return
	}

	path := filepath.Join(dirPath, paletteFileName)

	pData := paletteData{
		Categories: make([]categoryData, 0, len(p.categories)),
	}

	for _, cat := range p.categories {
		cData := categoryData{
			Name:    cat.Name,
			Prefabs: make([]string, 0, len(cat.Prefabs)),
		}
		for _, prefab := range cat.Prefabs {
			cData.Prefabs = append(cData.Prefabs, prefab.Path())
		}
		pData.Categories = append(pData.Categories, cData)
	}

	data, err := json.MarshalIndent(pData, "", "\t")
	if err != nil {
		log.Printf("failed to marshal palette data: %v", err)
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Printf("failed to write palette file: %v", err)
		return
	}
	log.Print("palette saved to:", path)
}

func (p *Palette) Import(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var pData paletteData
	if err := json.Unmarshal(data, &pData); err != nil {
		return err
	}

	for _, catData := range pData.Categories {
		cat := &Category{
			Name:    catData.Name,
			Prefabs: make([]*dmmprefab.Prefab, 0, len(catData.Prefabs)),
		}

		for _, path := range catData.Prefabs {
			if obj, ok := p.app.LoadedEnvironment().Objects[path]; ok {
				prefab := dmmprefab.New(dmmprefab.IdNone, path, obj.Vars)
				cat.Prefabs = append(cat.Prefabs, prefab)
			} else {
				prefab := dmmprefab.New(dmmprefab.IdNone, path, nil)
				cat.Prefabs = append(cat.Prefabs, prefab)
			}
		}
		p.categories = append(p.categories, cat)
	}
	
	p.Save()
	return nil
}

func (p *Palette) Export(path string) error {
	pData := paletteData{
		Categories: make([]categoryData, 0, len(p.categories)),
	}

	for _, cat := range p.categories {
		cData := categoryData{
			Name:    cat.Name,
			Prefabs: make([]string, 0, len(cat.Prefabs)),
		}
		for _, prefab := range cat.Prefabs {
			cData.Prefabs = append(cData.Prefabs, prefab.Path())
		}
		pData.Categories = append(pData.Categories, cData)
	}

	data, err := json.MarshalIndent(pData, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
