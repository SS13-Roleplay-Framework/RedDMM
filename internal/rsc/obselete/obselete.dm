// RedDMM Obsolete Object Placeholders
// These types store data from objects that were missing from the DME when the map was loaded.
// They preserve the original type path, name, and all variables so the data is not lost.
//
// To restore an obsolete object:
// 1. Add the original type back to your codebase
// 2. Use RedDMM's "Replace Obsolete" tool to convert obsolete objects back
//
// Variables stored on each obsolete type:
// - original_type: The full type path of the missing object (e.g., "/obj/machinery/old_thing")
// - original_name: The name the object had when it was saved
// - original_vars: A list of all variables that were set on the object, as strings

/obj/obselete
	name = "obsolete object"
	icon = 'obselete.dmi'
	icon_state = "obselete_obj"
	var/original_type = null
	var/original_name = null
	var/list/original_vars = list()

/turf/obselete
	name = "obsolete turf"
	icon = 'obselete.dmi'
	icon_state = "obselete_turf"
	var/original_type = null
	var/original_name = null
	var/list/original_vars = list()

/area/obselete
	name = "obsolete area"
	icon = 'obselete.dmi'
	icon_state = "obselete_area"
	var/original_type = null
	var/original_name = null
	var/list/original_vars = list()
