package dialog

import (
	"path/filepath"

	"sdmm/internal/app/window"
	"sdmm/internal/imguiext/style"
	w "sdmm/internal/imguiext/widget"

	"github.com/SpaiR/imgui-go"
)

type InstallObsoleteModule struct {
	simple

	TargetDir string
	OnInstall func(path string)

	errMessage string
}

func NewInstallObsoleteModule(onInstall func(string)) *InstallObsoleteModule {
	return &InstallObsoleteModule{
		simple: simple{
			name: "Install Obsolete Module",
		},
		TargetDir: "code/modules/obsolete", // Default
		OnInstall: onInstall,
	}
}

func (d *InstallObsoleteModule) Process() {
	if imgui.Button("Cancel") {
		imgui.CloseCurrentPopup()
	}
	
	imgui.SameLine()
	
	if imgui.Button("Install") {
		if d.TargetDir == "" {
			d.errMessage = "Target directory cannot be empty"
		} else {
			d.OnInstall(d.TargetDir)
			imgui.CloseCurrentPopup()
		}
	}

	imgui.Separator()

	w.Layout{
		w.Text("Target Directory (relative to project root):"),
		w.InputText("##target_dir", &d.TargetDir),
		w.TextDisabled("Files internal/rsc/obselete/obselete.dm and .dmi will be copied here."),
	}.Build()

	if d.errMessage != "" {
		w.TextColored(style.ColorError, d.errMessage)
	}
}
