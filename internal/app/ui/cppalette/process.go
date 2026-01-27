package cppalette

import (
	"fmt"

	"sdmm/internal/app/ui/dialog"
	"github.com/rs/zerolog/log"

	"sdmm/internal/app/render"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/imguiext/icon"
	"sdmm/internal/imguiext/style"
	w "sdmm/internal/imguiext/widget"

	"github.com/SpaiR/imgui-go"
)

func (p *Palette) Process(_ int32) {
	if !p.app.HasLoadedEnvironment() {
		p.loadedEnvPath = "" // Reset
		w.Layout{
			w.TextDisabled("No environment loaded"),
		}.Build()
		return
	}

	// Check if we need to load palette for the current environment
	currentEnv := p.app.LoadedEnvironment().RootDir
	if p.loadedEnvPath != currentEnv {
		p.Load() // This will populate p.categories
		p.loadedEnvPath = currentEnv
	}

	p.showToolbar()
	imgui.Separator()
	p.showCategories()
}

func (p *Palette) showToolbar() {
	w.Layout{
		w.Button(icon.Add, p.doAddCategory).
			Tooltip("Add new category").
			Round(true),
		w.SameLine(),
		w.Disabled(p.selectedCategory < 0, w.Layout{
			w.Button(icon.Delete, p.doRemoveCategory).
				Tooltip("Remove selected category").
				Round(true),
		}),
		w.SameLine(),
		w.Button(icon.Upload, p.doImport).
			Tooltip("Import Palette").
			Round(true),
		w.SameLine(),
		w.Button(icon.Download, p.doExport).
			Tooltip("Export Palette").
			Round(true),
	}.Build()
}

func (p *Palette) doImport() {
	dialog.Open(dialog.NewInput("Import Palette", "File Path:", "", func(path string) {
		if err := p.Import(path); err != nil {
			log.Printf("import failed: %v", err)
			dialog.Open(dialog.TypeInformation{Title: "Error", Information: fmt.Sprint(err)})
		}
	}))
}

func (p *Palette) doExport() {
	dialog.Open(dialog.NewInput("Export Palette", "File Path:", "", func(path string) {
		if err := p.Export(path); err != nil {
			log.Printf("export failed: %v", err)
			dialog.Open(dialog.TypeInformation{Title: "Error", Information: fmt.Sprint(err)})
		}
	}))
}

func (p *Palette) showCategories() {
	if len(p.categories) == 0 {
		w.Layout{
			w.TextDisabled("No categories"),
			w.Text("Click + to add a category"),
		}.Build()
		return
	}

	iconSize := p.iconSize()

	for catIdx, cat := range p.categories {
		// Category header
		open := imgui.TreeNodeV(cat.Name, imgui.TreeNodeFlagsDefaultOpen|imgui.TreeNodeFlagsFramePadding)
		
		if imgui.IsItemClicked() {
			p.selectedCategory = catIdx
		}

		if open {
			// Show prefabs in category
			for prefabIdx, prefab := range cat.Prefabs {
				p.showPrefabItem(catIdx, prefabIdx, prefab, iconSize)
			}

			if len(cat.Prefabs) == 0 {
				imgui.TextDisabled("Empty category")
			}

			imgui.TreePop()
		}
	}
}

func (p *Palette) showPrefabItem(catIdx, prefabIdx int, prefab *dmmprefab.Prefab, iconSize float32) {
	imgui.PushIDInt(catIdx*1000 + prefabIdx)
	defer imgui.PopID()

	// Draw icon
	sprite := render.CreateSpriteOrPlaceholder(prefab.Vars())
	if sprite != nil {
		uv0, uv1 := sprite.UV()
		imgui.ImageV(
			imgui.TextureID(sprite.Texture()),
			imgui.Vec2{X: iconSize, Y: iconSize},
			uv0, uv1,
			style.ColorWhite, imgui.Vec4{},
		)
	} else {
		imgui.Dummy(imgui.Vec2{X: iconSize, Y: iconSize})
	}

	imgui.SameLine()

	// Draw name and type
	imgui.BeginGroup()
	
	// Name at top
	name := prefab.Vars().TextV("name", "")
	if name == "" {
		name = "(unnamed)"
	}
	imgui.Text(name)
	
	// Type at bottom (smaller, dimmed)
	path := prefab.Path()
	imgui.PushStyleColor(imgui.StyleColorText, style.ColorGrey)
	imgui.TextV(path)
	imgui.PopStyleColor()
	
	imgui.EndGroup()

	// Handle selection
	if imgui.IsItemClicked() {
		p.selectedPrefab = prefab
		p.app.DoSelectPrefab(prefab)
	}

	// Context menu for removal
	if imgui.BeginPopupContextItem() {
		if imgui.MenuItem("Remove from Palette") {
			p.RemovePrefabFromCategory(catIdx, prefabIdx)
		}
		imgui.EndPopup()
	}
}

func (p *Palette) doAddCategory() {
	// For now, use a simple incrementing name
	name := fmt.Sprintf("Category %d", len(p.categories)+1)
	p.AddCategory(name)
}

func (p *Palette) doRemoveCategory() {
	if p.selectedCategory >= 0 {
		p.RemoveCategory(p.selectedCategory)
		p.selectedCategory = -1
	}
}
