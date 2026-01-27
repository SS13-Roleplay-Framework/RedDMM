package cpenvironment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sdmm/internal/app/ui/dialog"
	"sdmm/internal/dmapi/dmenv"
	"sdmm/internal/imguiext/icon"
	"sdmm/internal/imguiext/style"
	w "sdmm/internal/imguiext/widget"
	"sdmm/internal/rsc"

	"github.com/SpaiR/imgui-go"
	"github.com/rs/zerolog/log"
)

func (e *Environment) Process(int32) {
	if e.app.LoadedEnvironment() == nil {
		imgui.TextDisabled("No environment loaded")
	} else {
		e.process()
		e.showControls()
		e.showTree()
		e.postProcess()
	}
}

func (e *Environment) showControls() {
	w.Button(icon.Remove, e.doCollapseAll).
		Tooltip("Collapse All").
		Round(true).
		Build()
	imgui.SameLine()
	e.showTypesFilterButton()
	imgui.SameLine()
	e.showSettingsButton()
	imgui.SameLine()
	w.InputTextWithHint("##filter", "Filter", &e.filter).
		ButtonClear().
		Width(-1).
		OnChange(e.doFilter).
		Build()
	imgui.Separator()
}

func (e *Environment) showTypesFilterButton() {
	var bStyle w.ButtonStyle

	if !e.typesFilterEnabled {
		bStyle = style.ButtonDefault{}
	} else {
		bStyle = style.ButtonGreen{}
	}

	w.Layout{
		w.Button(icon.Eye, e.doToggleTypesFilter).
			Round(true).
			Style(bStyle),
		w.Tooltip(
			w.AlignTextToFramePadding(),
			w.Text("Types Filter"),
			w.SameLine(),
			w.TextFrame("F"),
		),
	}.Build()
}

func (e *Environment) showSettingsButton() {
	w.Layout{
		w.Button(icon.Cog, nil).
			Round(true).
			Tooltip("Settings"),
	}.Build()

	if imgui.BeginPopupContextItemV("environment_settings", imgui.PopupFlagsMouseButtonLeft) {
		imgui.AlignTextToFramePadding()
		imgui.Text("Icon Size")
		imgui.SameLine()
		imgui.SliderInt("##icon_size", &e.config().NodeScale, 50, 300)
		
		imgui.Separator()
		if imgui.MenuItem("Install Obsolete Module...") {
			e.openInstallObsoleteModuleDialog()
		}
		
		imgui.EndPopup()
	}
}

func (e *Environment) openInstallObsoleteModuleDialog() {
	dialog.Open(dialog.NewInstallObsoleteModule(e.installObsoleteModule))
}

func (e *Environment) installObsoleteModule(targetRel string) {
	if e.app.LoadedEnvironment() == nil {
		return
	}

	rootDir := e.app.LoadedEnvironment().RootDir
	targetDir := filepath.Join(rootDir, targetRel)

	log.Printf("installing obsolete module to: %s", targetDir)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Printf("failed to create directory: %v", err)
		dialog.Open(dialog.TypeInformation{
			Title:       "Error",
			Information: fmt.Sprintf("Failed to create directory:\n%v", err),
		})
		return
	}

	dmPath := filepath.Join(targetDir, "obselete.dm")
	if err := os.WriteFile(dmPath, []byte(rsc.ObsoleteDM), 0644); err != nil {
		log.Printf("failed to write obselete.dm: %v", err)
		dialog.Open(dialog.TypeInformation{
			Title:       "Error",
			Information: fmt.Sprintf("Failed to write obselete.dm:\n%v", err),
		})
		return
	}

	dmiPath := filepath.Join(targetDir, "obselete.dmi")
	if err := os.WriteFile(dmiPath, rsc.ObsoleteDMI, 0644); err != nil {
		log.Printf("failed to write obselete.dmi: %v", err)
	}

	// append to DME
	dmePath := e.app.LoadedEnvironment().DmePath
	relPath, err := filepath.Rel(filepath.Dir(dmePath), dmPath)
	if err != nil {
		log.Printf("failed to get relative path: %v", err)
		return
	}
	relPath = strings.ReplaceAll(relPath, "\\", "/")

	f, err := os.OpenFile(dmePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open DME: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("\n#include \"%s\"\n", relPath)); err != nil {
		log.Printf("failed to append to DME: %v", err)
	}
	
	dialog.Open(dialog.TypeInformation{
		Title:       "Success",
		Information: fmt.Sprintf("Successfully installed obsolete module to:\n%s\n\nAdded include to DME.", targetRel),
	})
}

func (e *Environment) showTree() {
	if imgui.BeginChild(fmt.Sprintf("environment_tree_[%d]", e.treeId)) {
		if len(e.filter) == 0 {
			e.showPathBranch("/area")
			e.showPathBranch("/turf")
			e.showPathBranch("/obj")
			e.showPathBranch("/mob")
			e.showMapStampBranch()
		} else {
			e.showFilteredNodes()
		}
	}
	imgui.EndChild()
}

func (e *Environment) showMapStampBranch() {
	presets := e.app.Presets().Items()
	if len(presets) == 0 {
		return
	}

	rootPath := "/map_stamp"

	// Root node
	imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, e.calcTreeNodePadding("map_stamp"))
	
	// Flags
	flags := int(imgui.TreeNodeFlagsSpanAvailWidth | imgui.TreeNodeFlagsOpenOnArrow | imgui.TreeNodeFlagsOpenOnDoubleClick)
	if rootPath == e.selectedPath {
		flags |= int(imgui.TreeNodeFlagsSelected)
	}
	
	opened := imgui.TreeNodeV("map_stamp", imgui.TreeNodeFlags(flags))
	imgui.PopStyleVar()

	// Handle selection for root
	if imgui.IsItemClicked() && e.selectedPath != rootPath {
		// e.app.DoSelectPrefabByPath(rootPath) // Cannot select virtual root yet?
		e.selectedPath = rootPath
	}

	if opened {
		for _, preset := range presets {
			path := rootPath + "/" + preset.Name
			
			imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, e.calcTreeNodePadding(preset.Name))
			
			nodeFlags := int(imgui.TreeNodeFlagsSpanAvailWidth | imgui.TreeNodeFlagsLeaf | imgui.TreeNodeFlagsNoTreePushOnOpen)
			if path == e.selectedPath {
				nodeFlags |= int(imgui.TreeNodeFlagsSelected)
			}
			
			imgui.TreeNodeV(preset.Name, imgui.TreeNodeFlags(nodeFlags))
			imgui.PopStyleVar()
			
			if imgui.IsItemClicked() && e.selectedPath != path {
				e.app.DoSelectPrefabByPath(path)
			}
		}
		imgui.TreePop()
	}
}

func (e *Environment) showFilteredNodes() {
	var clipper imgui.ListClipper
	clipper.Begin(len(e.filteredTreeNodes))
	for clipper.Step() {
		for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
			node := e.filteredTreeNodes[i]
			if e.showAttachment(node) {
				imgui.SameLine()
			}

			imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, e.calcTreeNodePadding(node.name))
			imgui.TreeNodeV(node.orig.Path, e.nodeFlags(node, true))

			imgui.PopStyleVar()

			e.doSelectOnClick(node)
			e.showNodeMenu(node)
		}
	}
}

func (e *Environment) showPathBranch(t string) {
	if atom := e.app.LoadedEnvironment().Objects[t]; atom != nil {
		e.showBranch0(atom)
	}
}

func (e *Environment) showBranch0(object *dmenv.Object) {
	node, ok := e.newTreeNode(object)
	if !ok {
		return
	}

	if e.showAttachment(node) {
		imgui.SameLine()
	}

	imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, e.calcTreeNodePadding(node.name))

	if len(object.DirectChildren) == 0 {
		imgui.AlignTextToFramePadding()
		imgui.TreeNodeV(node.name, e.nodeFlags(node, true))
		imgui.PopStyleVar()
		e.doSelectOnClick(node)
		e.showNodeMenu(node)
		e.scrollToSelectedPath(node)
	} else {
		if e.isPartOfSelectedPath(node) {
			imgui.SetNextItemOpen(true, imgui.ConditionAlways)
		}

		imgui.AlignTextToFramePadding()
		opened := imgui.TreeNodeV(node.name, e.nodeFlags(node, false))
		imgui.PopStyleVar()

		if opened {
			e.doSelectOnClick(node)
			e.showNodeMenu(node)
			e.scrollToSelectedPath(node)

			if e.tmpDoCollapseAll {
				imgui.StateStorage().SetAllInt(0)
			}
			for _, childPath := range object.DirectChildren {
				e.showBranch0(e.app.LoadedEnvironment().Objects[childPath])
			}
			imgui.TreePop()
		} else {
			e.doSelectOnClick(node)
			e.showNodeMenu(node)
		}
	}
}

func (e *Environment) doSelectOnClick(node *treeNode) {
	if imgui.IsItemClicked() && e.selectedPath != node.orig.Path {
		e.app.DoSelectPrefabByPath(node.orig.Path)
		e.app.DoEditPrefabByPath(node.orig.Path)
		e.tmpDoSelectPath = false // we don't need to scroll tree when we select item from the tree itself
	}
}

func (e *Environment) nodeFlags(node *treeNode, leaf bool) imgui.TreeNodeFlags {
	flags := int(imgui.TreeNodeFlagsSpanAvailWidth)
	if node.orig.Path == e.selectedPath {
		flags |= int(imgui.TreeNodeFlagsSelected)
	}
	if leaf {
		flags |= int(imgui.TreeNodeFlagsLeaf | imgui.TreeNodeFlagsNoTreePushOnOpen)
	} else {
		flags |= int(imgui.TreeNodeFlagsOpenOnArrow | imgui.TreeNodeFlagsOpenOnDoubleClick)
	}
	return imgui.TreeNodeFlags(flags)
}

func (e *Environment) isPartOfSelectedPath(node *treeNode) bool {
	return e.tmpDoSelectPath && strings.HasPrefix(e.selectedPath, node.orig.Path)
}

func (e *Environment) scrollToSelectedPath(node *treeNode) {
	if e.tmpDoSelectPath && e.selectedPath == node.orig.Path {
		e.tmpDoSelectPath = false
		imgui.SetScrollHereY(.5)
	}
}

func (e *Environment) showAttachment(node *treeNode) bool {
	if e.typesFilterEnabled {
		e.showVisibilityCheckbox(node)
		return true
	} else {
		e.showIcon(node)
		return false
	}
}

func (e *Environment) showVisibilityCheckbox(node *treeNode) {
	value := e.app.PathsFilter().IsVisiblePath(node.orig.Path)
	vOrig := value

	var hasHiddenChildPath bool
	if value {
		hasHiddenChildPath = e.app.PathsFilter().HasHiddenChildPath(node.orig.Path)
		value = !hasHiddenChildPath
	}

	imgui.PushStyleVarVec2(imgui.StyleVarFramePadding, e.calcTreeNodePadding(node.name))
	if imgui.Checkbox(fmt.Sprint("##node_visibility_", node.orig.Path), &value) {
		e.app.PathsFilter().TogglePath(node.orig.Path)
	}

	imgui.PopStyleVar()

	// Show a dash symbol, if the node has any hidden child.
	if vOrig && hasHiddenChildPath {
		iMin := imgui.ItemRectMin()
		iMax := imgui.ItemRectMax()
		iWidth := iMax.X - iMin.X
		iHeight := iMax.Y - iMin.Y
		mPadding := iWidth * .1  // height
		mHeight := iHeight * .15 // left/right padding
		mMinY := iMin.Y + iHeight/2 - mHeight/2
		mMaxY := iMin.Y + iHeight/2 + mHeight/2
		col := imgui.PackedColorFromVec4(imgui.CurrentStyle().Color(imgui.StyleColorCheckMark))
		imgui.WindowDrawList().AddRectFilled(imgui.Vec2{X: iMin.X + mPadding, Y: mMinY}, imgui.Vec2{X: iMax.X - mPadding, Y: mMaxY}, col)
	}
}

func (e *Environment) showIcon(node *treeNode) {
	s := node.sprite
	w.Image(imgui.TextureID(s.Texture()), e.iconSize(), e.iconSize()).
		TintColor(node.color).
		Uv(imgui.Vec2{X: s.U1, Y: s.V1}, imgui.Vec2{X: s.U2, Y: s.V2}).
		Build()
	imgui.SameLine()
}

func (e *Environment) calcTreeNodePadding(nodeName string) imgui.Vec2 {
	x := imgui.CurrentStyle().FramePadding().X
	textSize := imgui.CalcTextSize(nodeName, false, 0).Y
	y := (e.iconSize() - textSize) / 2
	return imgui.Vec2{
		X: x,
		Y: y,
	}
}
