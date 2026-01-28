package obselete

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sdmm/internal/rsc"

	"github.com/rs/zerolog/log"
)

// InsertModule creates the obselete module files in the target directory and updates the DME.
func InsertModule(targetDir, dmePath string) error {
	// Create the obselete directory
	obseletePath := filepath.Join(targetDir, "obselete")
	if err := os.MkdirAll(obseletePath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create obselete directory: %w", err)
	}

	// Write the obselete.dm file
	dmFilePath := filepath.Join(obseletePath, "obselete.dm")
	if err := os.WriteFile(dmFilePath, []byte(rsc.ObsoleteDM), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write obselete.dm: %w", err)
	}
	log.Print("created obselete.dm at:", dmFilePath)

	// Update the DME to include the obselete module
	if err := updateDME(dmePath, targetDir); err != nil {
		return fmt.Errorf("failed to update DME: %w", err)
	}

	return nil
}

// updateDME adds the obselete module include to the DME file.
func updateDME(dmePath, targetDir string) error {
	// Read the DME file
	content, err := os.ReadFile(dmePath)
	if err != nil {
		return fmt.Errorf("failed to read DME: %w", err)
	}

	dmeContent := string(content)

	// Calculate the relative path for the include
	dmeDir := filepath.Dir(dmePath)
	relPath, err := filepath.Rel(dmeDir, filepath.Join(targetDir, "obselete", "obselete.dm"))
	if err != nil {
		return fmt.Errorf("failed to calculate relative path: %w", err)
	}
	// Convert to forward slashes for DM
	relPath = strings.ReplaceAll(relPath, "\\", "/")

	// Check if already included
	includeLine := fmt.Sprintf(`#include "%s"`, relPath)
	if strings.Contains(dmeContent, includeLine) {
		log.Print("obselete module already included in DME")
		return nil
	}

	// Find a good place to insert the include (after other includes or at the end)
	// We'll append to the end of the file
	if !strings.HasSuffix(dmeContent, "\n") {
		dmeContent += "\n"
	}
	dmeContent += includeLine + "\n"

	// Write the updated DME
	if err := os.WriteFile(dmePath, []byte(dmeContent), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write DME: %w", err)
	}
	log.Print("added obselete include to DME:", dmePath)

	return nil
}

// IsModuleInstalled checks if the obselete module is available in the environment.
func IsModuleInstalled(dmeDir string, objPath, turfPath, areaPath string) bool {
	// Simple check: see if any of the configured paths exist
	// This would need to be checked against the environment after parsing
	return objPath != "" || turfPath != "" || areaPath != ""
}
