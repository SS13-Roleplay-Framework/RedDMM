package dialog

import (
	w "sdmm/internal/imguiext/widget"

	"github.com/SpaiR/imgui-go"
)

type TypeInput struct {
	simple
	Label     string
	Value     string
	OnConfirm func(string)
}

func NewInput(title, label, defaultValue string, onConfirm func(string)) *TypeInput {
	return &TypeInput{
		simple: simple{
			name: title,
		},
		Label:     label,
		Value:     defaultValue,
		OnConfirm: onConfirm,
	}
}

func (d *TypeInput) Process() {
	if imgui.Button("Cancel") {
		imgui.CloseCurrentPopup()
	}
	imgui.SameLine()
	if imgui.Button("OK") {
		if d.OnConfirm != nil {
			d.OnConfirm(d.Value)
		}
		imgui.CloseCurrentPopup()
	}

	imgui.Separator()
	w.Layout{
		w.Text(d.Label),
		w.InputText("##input_value", &d.Value).Width(300),
	}.Build()
}
