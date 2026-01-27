package themes

import (
	"sdmm/internal/util"

	"github.com/SpaiR/imgui-go"
)

var colorMap = map[string]imgui.StyleColorID{
	"Text":                  imgui.StyleColorText,
	"TextDisabled":          imgui.StyleColorTextDisabled,
	"WindowBg":              imgui.StyleColorWindowBg,
	"ChildBg":               imgui.StyleColorChildBg,
	"PopupBg":               imgui.StyleColorPopupBg,
	"Border":                imgui.StyleColorBorder,
	"BorderShadow":          imgui.StyleColorBorderShadow,
	"FrameBg":               imgui.StyleColorFrameBg,
	"FrameBgHovered":        imgui.StyleColorFrameBgHovered,
	"FrameBgActive":         imgui.StyleColorFrameBgActive,
	"TitleBg":               imgui.StyleColorTitleBg,
	"TitleBgActive":         imgui.StyleColorTitleBgActive,
	"TitleBgCollapsed":      imgui.StyleColorTitleBgCollapsed,
	"MenuBarBg":             imgui.StyleColorMenuBarBg,
	"ScrollbarBg":           imgui.StyleColorScrollbarBg,
	"ScrollbarGrab":         imgui.StyleColorScrollbarGrab,
	"ScrollbarGrabHovered":  imgui.StyleColorScrollbarGrabHovered,
	"ScrollbarGrabActive":   imgui.StyleColorScrollbarGrabActive,
	"CheckMark":             imgui.StyleColorCheckMark,
	"SliderGrab":            imgui.StyleColorSliderGrab,
	"SliderGrabActive":      imgui.StyleColorSliderGrabActive,
	"Button":                imgui.StyleColorButton,
	"ButtonHovered":         imgui.StyleColorButtonHovered,
	"ButtonActive":          imgui.StyleColorButtonActive,
	"Header":                imgui.StyleColorHeader,
	"HeaderHovered":         imgui.StyleColorHeaderHovered,
	"HeaderActive":          imgui.StyleColorHeaderActive,
	"Separator":             imgui.StyleColorSeparator,
	"SeparatorHovered":      imgui.StyleColorSeparatorHovered,
	"SeparatorActive":       imgui.StyleColorSeparatorActive,
	"ResizeGrip":            imgui.StyleColorResizeGrip,
	"ResizeGripHovered":     imgui.StyleColorResizeGripHovered,
	"ResizeGripActive":      imgui.StyleColorResizeGripActive,
	"Tab":                   imgui.StyleColorTab,
	"TabHovered":            imgui.StyleColorTabHovered,
	"TabActive":             imgui.StyleColorTabActive,
	"TabUnfocused":          imgui.StyleColorTabUnfocused,
	"TabUnfocusedActive":    imgui.StyleColorTabUnfocusedActive,
	"PlotLines":             imgui.StyleColorPlotLines,
	"PlotLinesHovered":      imgui.StyleColorPlotLinesHovered,
	"PlotHistogram":         imgui.StyleColorPlotHistogram,
	"PlotHistogramHovered":  imgui.StyleColorPlotHistogramHovered,
	"TextSelectedBg":        imgui.StyleColorTextSelectedBg,
	"DragDropTarget":        imgui.StyleColorDragDropTarget,
	"NavHighlight":          imgui.StyleColorNavHighlight,
	"NavWindowingHighlight": imgui.StyleColorNavWindowingHighlight,
	"NavWindowingDimBg":     imgui.StyleColorNavWindowingDimBg,
	"ModalWindowDimBg":      imgui.StyleColorModalWindowDimBg,
}

func Apply(t *Theme) {
	if t.Name == "Light" {
		imgui.StyleColorsLight()
	} else {
		imgui.StyleColorsDark()
	}

	style := imgui.CurrentStyle()
	for _, c := range t.Colors {
		if id, ok := colorMap[c.Name]; ok {
			col := util.ParseColor(c.Value)
			style.SetColor(id, imgui.Vec4{X: col.R(), Y: col.G(), Z: col.B(), W: col.A()})
		}
	}
}
