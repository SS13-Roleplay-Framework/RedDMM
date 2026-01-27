package cppalette

import (
	"sdmm/internal/app/ui/component"
	"sdmm/internal/app/window"
	"sdmm/internal/dmapi/dmenv"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"

	"github.com/rs/zerolog/log"
)

// App interface for palette operations
type App interface {
	DoSelectPrefab(*dmmprefab.Prefab)
	DoEditPrefab(prefab *dmmprefab.Prefab)
	LoadedEnvironment() *dmenv.Dme
	HasLoadedEnvironment() bool
	ShowLayout(name string, focus bool)
}

// Category represents a user-created palette category
type Category struct {
	Name   string
	Prefabs []*dmmprefab.Prefab
}

// Palette is the main palette panel component
type Palette struct {
	component.Component

	app App

	categories []*Category
	selectedCategory int
	selectedPrefab *dmmprefab.Prefab
	
	loadedEnvPath string
}

func (p *Palette) Init(app App) {
	p.app = app
	p.categories = []*Category{}
	p.selectedCategory = -1
	log.Print("palette initialized")
}

func (p *Palette) Free() {
	p.categories = nil
	p.selectedPrefab = nil
	p.selectedCategory = -1
}

// AddCategory creates a new palette category
func (p *Palette) AddCategory(name string) {
	p.categories = append(p.categories, &Category{
		Name:    name,
		Prefabs: []*dmmprefab.Prefab{},
	})
	log.Print("palette category added:", name)
	p.Save()
}

// RemoveCategory removes a category by index
func (p *Palette) RemoveCategory(idx int) {
	if idx >= 0 && idx < len(p.categories) {
		name := p.categories[idx].Name
		p.categories = append(p.categories[:idx], p.categories[idx+1:]...)
		log.Print("palette category removed:", name)
		p.Save()
	}
}

// AddPrefabToCategory adds a prefab to a category
func (p *Palette) AddPrefabToCategory(categoryIdx int, prefab *dmmprefab.Prefab) {
	if categoryIdx >= 0 && categoryIdx < len(p.categories) {
		p.categories[categoryIdx].Prefabs = append(p.categories[categoryIdx].Prefabs, prefab)
		log.Print("prefab added to palette:", prefab.Path())
		p.Save()
	}
}

// RemovePrefabFromCategory removes a prefab from a category
func (p *Palette) RemovePrefabFromCategory(categoryIdx, prefabIdx int) {
	if categoryIdx >= 0 && categoryIdx < len(p.categories) {
		cat := p.categories[categoryIdx]
		if prefabIdx >= 0 && prefabIdx < len(cat.Prefabs) {
			cat.Prefabs = append(cat.Prefabs[:prefabIdx], cat.Prefabs[prefabIdx+1:]...)
			p.Save()
		}
	}
}

// Categories returns all palette categories
func (p *Palette) Categories() []*Category {
	return p.categories
}

// SelectedCategory returns the currently selected category index
func (p *Palette) SelectedCategory() int {
	return p.selectedCategory
}

// SetSelectedCategory sets the selected category
func (p *Palette) SetSelectedCategory(idx int) {
	p.selectedCategory = idx
}

func (p *Palette) iconSize() float32 {
	return 32 * window.PointSize()
}
