package themes

var Default = RedDMM()

var Presets = map[string]func() *Theme{
	"RedDMM": RedDMM,
	"Dark":   Dark,
	"Light":  Light,
}

func Dark() *Theme {
	return &Theme{
		Name: "Dark",
		Colors: []Color{
			// Standard ImGui Dark colors (approximate or just Name it so manager knows to apply StyleColorsDark base)
			// Actually manager logic: if Name == "Dark", apply StyleColorsDark first?
			// Or we just define overrides.
		},
	}
}

func Light() *Theme {
	return &Theme{
		Name: "Light",
		Colors: []Color{
			
		},
	}
}

func RedDMM() *Theme {
	return &Theme{
		Name: "RedDMM",
		Colors: []Color{
			{Name: "Text", Value: "#e6e6e6ff"},
			{Name: "TextDisabled", Value: "#808080ff"},
			{Name: "WindowBg", Value: "#1a1a2eff"},
			{Name: "ChildBg", Value: "#161621ff"},
			{Name: "PopupBg", Value: "#1a1a2ef0"},
			{Name: "Border", Value: "#e9456080"},
			{Name: "BorderShadow", Value: "#00000000"},
			{Name: "FrameBg", Value: "#16213eff"},
			{Name: "FrameBgHovered", Value: "#1a1a2eff"},
			{Name: "FrameBgActive", Value: "#e9456060"},
			{Name: "TitleBg", Value: "#16213eff"},
			{Name: "TitleBgActive", Value: "#1a1a2eff"},
			{Name: "TitleBgCollapsed", Value: "#16213e80"},
			{Name: "MenuBarBg", Value: "#16213eff"},
			{Name: "ScrollbarBg", Value: "#00000000"},
			{Name: "ScrollbarGrab", Value: "#e9456060"},
			{Name: "ScrollbarGrabHovered", Value: "#e94560a0"},
			{Name: "ScrollbarGrabActive", Value: "#e94560e0"},
			{Name: "CheckMark", Value: "#e94560ff"},
			{Name: "SliderGrab", Value: "#e94560ff"},
			{Name: "SliderGrabActive", Value: "#e94560ff"},
			{Name: "Button", Value: "#e9456060"},
			{Name: "ButtonHovered", Value: "#e94560a0"},
			{Name: "ButtonActive", Value: "#e94560e0"},
			{Name: "Header", Value: "#e9456060"},
			{Name: "HeaderHovered", Value: "#e94560a0"},
			{Name: "HeaderActive", Value: "#e94560e0"},
			{Name: "Separator", Value: "#e9456080"},
			{Name: "SeparatorHovered", Value: "#e94560a0"},
			{Name: "SeparatorActive", Value: "#e94560e0"},
			{Name: "ResizeGrip", Value: "#e9456040"},
			{Name: "ResizeGripHovered", Value: "#e94560a0"},
			{Name: "ResizeGripActive", Value: "#e94560e0"},
			{Name: "Tab", Value: "#16213eff"},
			{Name: "TabHovered", Value: "#e94560a0"},
			{Name: "TabActive", Value: "#e94560e0"},
			{Name: "TabUnfocused", Value: "#16213eff"},
			{Name: "TabUnfocusedActive", Value: "#1a1a2eff"},
			{Name: "PlotLines", Value: "#e94560ff"},
			{Name: "PlotLinesHovered", Value: "#ff6e85ff"},
			{Name: "PlotHistogram", Value: "#e94560ff"},
			{Name: "PlotHistogramHovered", Value: "#ff6e85ff"},
			{Name: "TextSelectedBg", Value: "#e9456060"},
			{Name: "DragDropTarget", Value: "#ffff00e6"},
			{Name: "NavHighlight", Value: "#e94560ff"},
			{Name: "NavWindowingHighlight", Value: "#ffffffb3"},
			{Name: "NavWindowingDimBg", Value: "#33333333"},
			{Name: "ModalWindowDimBg", Value: "#3333335a"},
		},
	}
}
