package ctrl

import (
	"fmt"
)

const (

	actionsPath       = "actions"
	gitPath           = "git"
	blueprintPath     = "blueprint"
)

type ServiceLoc struct {
	rootDir string
}

func NewServeLoc(dir string) *ServiceLoc {
	if dir == "" {
		dir = "/usr/local/var/ctrl"
	}

	return &ServiceLoc{
		rootDir: dir,
	}
}

func (sl *ServiceLoc) ActionFolderPath() string {
	return fmt.Sprintf("%s/%s", sl.rootDir, actionsPath)
}

func (sl *ServiceLoc) ActionPath(name string) string {
	return fmt.Sprintf("%s/%s/%s", sl.rootDir, actionsPath, name)
}

func (sl *ServiceLoc) GitPath() string {
	return fmt.Sprintf("%s/%s", sl.rootDir, gitPath)
}

func (sl *ServiceLoc) BlueprintPath() string {
	return fmt.Sprintf("%s/%s", sl.rootDir, blueprintPath)
}

