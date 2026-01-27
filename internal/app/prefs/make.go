package prefs

import (
	"math"
	"sort"

	"sdmm/internal/app/themes"
	"sdmm/internal/app/ui/cpwsarea/wsprefs"
	"sdmm/internal/app/window"
)

var themeOptions []string

func init() {
	for k := range themes.Presets {
		themeOptions = append(themeOptions, k)
	}
	sort.Strings(themeOptions)
}

type App interface {
	UpdateScale()
}

func Make(app App, prefs *Prefs) wsprefs.Prefs {
	p := wsprefs.MakePrefs()

	var preferencesPrefabs = map[wsprefs.PrefGroup][]prefPrefab{
		wsprefs.GPEditor: {
			optionPrefPrefab{
				name:    "Save Format",
				desc:    "Controls the format used by the editor to save the map.",
				label:   "##save_format",
				value:   &prefs.Editor.SaveFormat,
				options: SaveFormats,
				help:    SaveFormatHelp,
			},
			optionPrefPrefab{
				name:    "Code Editor",
				desc:    "Controls what code editor is opened when using Go to Definition.",
				label:   "##code_editor",
				value:   &prefs.Editor.CodeEditor,
				options: CodeEditors,
				help:    CodeEditorHelp,
			},
			boolPrefPrefab{
				name:  "Sanitize Variables",
				desc:  "Enables sanitizing for variables which are declared on the map, but has the same value as initial.",
				label: "##sanitize_variables",
				value: &prefs.Editor.SanitizeVariables,
			},
			optionPrefPrefab{
				name:    "Nudge Mode",
				desc:    "Controls which variables will be changed during the nudge.",
				label:   "##nudge_mode",
				value:   &prefs.Editor.NudgeMode,
				options: SaveNudgeModes,
			},
			stringPrefPrefab{
				name:  "Obsolete Object Path",
				desc:  "Type path to use when replacing missing /obj and /mob types. Leave empty to discard.",
				label: "##obsolete_obj_path",
				value: &prefs.Editor.ObsoleteObjectPath,
			},
			stringPrefPrefab{
				name:  "Obsolete Turf Path",
				desc:  "Type path to use when replacing missing /turf types. Leave empty to discard.",
				label: "##obsolete_turf_path",
				value: &prefs.Editor.ObsoleteTurfPath,
			},
			stringPrefPrefab{
				name:  "Obsolete Area Path",
				desc:  "Type path to use when replacing missing /area types. Leave empty to discard.",
				label: "##obsolete_area_path",
				value: &prefs.Editor.ObsoleteAreaPath,
			},
			boolPrefPrefab{
				name:  "Suppress Obsolete Var Transfer Warnings",
				desc:  "If enabled, no warning will be shown when variables cannot be transferred from an obsolete object to its replacement.",
				label: "##suppress_obsolete_var_warnings",
				value: &prefs.Editor.SuppressObsoleteVarWarning,
			},
		},

		wsprefs.GPControls: {
			boolPrefPrefab{
				name:  "Alternative Scroll Behavior",
				desc:  "When enabled, scrolling will do panning. Zoom will be available if a Space key pressed.",
				label: "##alternative_scroll_behavior",
				value: &prefs.Controls.AltScrollBehaviour,
			},
			boolPrefPrefab{
				name:  "Quick Edit: Tile Context Menu",
				desc:  "Controls whether Quick Edit should be shown in the tile context menu.",
				label: "##quick_edit:tile_context_menu",
				value: &prefs.Controls.QuickEditContextMenu,
			},
			boolPrefPrefab{
				name:  "Quick Edit: Map Pane",
				desc:  "Controls whether Quick Edit should be shown on the map pane.",
				label: "##quick_edit:map_pane",
				value: &prefs.Controls.QuickEditMapPane,
			},
		},

		wsprefs.GPInterface: {
			intPrefPrefab{
				name:  "Scale",
				desc:  "Controls the interface scale.",
				label: "%##scale",
				min:   50,
				max:   250,
				value: &prefs.Interface.Scale,
				post: func(int) {
					app.UpdateScale()
				},
			},
			intPrefPrefab{
				name:  "Fps",
				desc:  "Controls the application framerate.",
				label: "##fps",
				min:   30,
				max:   math.MaxInt,
				value: &prefs.Interface.Fps,
				post:  window.SetFps,
			},
			optionPrefPrefab{
				name:    "Theme",
				desc:    "Controls the application theme.",
				label:   "##theme",
				value:   &prefs.Interface.Theme,
				options: themeOptions,
				post: func(val string) {
					if factory, ok := themes.Presets[val]; ok {
						themes.Apply(factory())
					}
				},
			},
		},

		wsprefs.GPApplication: {
			boolPrefPrefab{
				name:  "Check for Updates",
				desc:  "When enabled, the editor will always check for updates on startup.",
				label: "##check_for_updates",
				value: &prefs.Application.CheckForUpdates,
			},
			boolPrefPrefab{
				name:  "Auto Update",
				desc:  "Enables automatic self-update, when a new update is available.",
				label: "##auto_update",
				value: &prefs.Application.AutoUpdate,
			},
		},
	}

	for group, prefabs := range preferencesPrefabs {
		for _, prefab := range prefabs {
			p.Add(group, prefab.make())
		}
	}

	return p
}
