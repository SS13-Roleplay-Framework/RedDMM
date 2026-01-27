package dialog

import (
	"fmt"
	"sort"

	"sdmm/internal/app/window"
	"sdmm/internal/dmapi/dmvars"
	w "sdmm/internal/imguiext/widget"

	"github.com/SpaiR/imgui-go"
)

type ObsoleteInfo struct {
	simple

	OriginalType string
	OriginalName string
	OriginalVars *dmvars.Variables
}

func NewObsoleteInfo(vars *dmvars.Variables) *ObsoleteInfo {
	d := &ObsoleteInfo{
		simple: simple{
			name: "Obsolete Object Info",
		},
		OriginalType: vars.TextV("original_type", "unknown"),
		OriginalName: vars.TextV("original_name", "unknown"),
		OriginalVars: vars,
	}
	return d
}

func (d *ObsoleteInfo) Process() {
	if imgui.Button("Close") {
		imgui.CloseCurrentPopup()
	}

	imgui.Separator()

	w.Layout{
		w.TextFrame("Original Type:"),
		w.SameLine(),
		w.Text(d.OriginalType),
		w.TextFrame("Original Name:"),
		w.SameLine(),
		w.Text(d.OriginalName),
	}.Build()

	imgui.Separator()
	w.Text("Original Variables:")

	imgui.BeginChildV("vars_scroll", imgui.Vec2{X: 400 * window.PointSize(), Y: 300 * window.PointSize()}, true, 0)
	
	// Access the original vars list if possible. 
	// Since original_vars is stored as a list in DM, we might need to parse it or iterate differently.
	// For now, let's just show key-value pairs if present.
	
	// Assuming OriginalVars might contain other info, let's just list what we have
	// In a real implementation we'd need to properly decode the 'original_vars' list from the DM var
	
	keys := d.OriginalVars.Keys()
	sort.Strings(keys)
	
	for _, k := range keys {
		// Skip the tracking vars themselves
		if k == "original_type" || k == "original_name" || k == "icon" || k == "icon_state" {
			continue
		}
		
		val := d.OriginalVars.ValueV(k, "")
		w.Layout{
			w.TextFrame(k + " ="),
			w.SameLine(),
			w.Text(val),
		}.Build()
	}
	
	imgui.EndChild()
}
