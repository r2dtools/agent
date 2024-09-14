package reverter

import (
	"fmt"
	"os"
	"slices"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/unknwon/com"
)

type rollbackError struct {
	err error
}

func (re *rollbackError) Error() string {
	return fmt.Sprintf("rollback failed: %v", re.err)
}

type hostManager interface {
	Enable(configFilePath string) error
	Disable(configFilePath string) error
}

// Reverter reverts change back for configuration files of virtual hosts
type Reverter struct {
	configsToDelete  []string
	configsToRestore map[string]string
	configsToDisable []string
	HostMng          hostManager
	Logger           logger.Logger
}

func (r *Reverter) AddConfigToDeletion(filePath string) {
	r.configsToDelete = append(r.configsToDelete, filePath)
}

func (r *Reverter) BackupConfigs(filePaths []string) error {
	for _, filePath := range filePaths {
		if err := r.BackupConfig(filePath); err != nil {
			return fmt.Errorf("could not make file '%s' backup: %v", filePath, err)
		}
	}

	return nil
}

func (r *Reverter) BackupConfig(filePath string) error {
	bFilePath := r.getBackupConfigPath(filePath)

	if _, ok := r.configsToRestore[filePath]; ok {
		r.Logger.Debug(fmt.Sprintf("file '%s' is already backed up.", filePath))
		return nil
	}

	// Skip file backup if it should be removed
	if slices.Contains(r.configsToDelete, filePath) {
		r.Logger.Debug(fmt.Sprintf("file '%s' will be removed on rollback. Skip its backup.", filePath))
		return nil
	}

	content, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	err = os.WriteFile(bFilePath, content, 0644)

	if err != nil {
		return err
	}

	if r.configsToRestore == nil {
		r.configsToRestore = make(map[string]string)
	}

	r.configsToRestore[filePath] = bFilePath

	return nil
}

func (r *Reverter) AddConfigToDisable(filePath string) {
	r.configsToDisable = append(r.configsToDisable, filePath)
}

func (r *Reverter) Rollback() error {
	// Disable all enabled before sites
	for _, configToDisable := range r.configsToDisable {
		if err := r.HostMng.Disable(configToDisable); err != nil {
			r.Logger.Error(fmt.Sprintf("failed to delete symlink to config file %s", configToDisable))
		}
	}

	// remove created files
	for _, fileToDelete := range r.configsToDelete {
		_, err := os.Stat(fileToDelete)

		if os.IsNotExist(err) {
			r.Logger.Debug(fmt.Sprintf("file '%s' does not exist. Skip its deletion.", fileToDelete))
			continue
		}

		if err != nil {
			return &rollbackError{err}
		}

		err = os.Remove(fileToDelete)

		if err != nil {
			return &rollbackError{err}
		}
	}

	if r.configsToRestore == nil {
		return nil
	}

	// restore the content of backed up files
	for originFilePath, bFilePath := range r.configsToRestore {
		bContent, err := os.ReadFile(bFilePath)

		if err != nil {
			return &rollbackError{err}
		}

		err = os.WriteFile(originFilePath, bContent, 0644)

		if err != nil {
			return &rollbackError{err}
		}

		if err := os.Remove(bFilePath); err != nil {
			r.Logger.Error(fmt.Sprintf("could not remove file '%s' on reverter rollback: %v", bFilePath, err))
		}

		delete(r.configsToRestore, originFilePath)
	}

	return nil
}

func (r *Reverter) Commit() error {
	for filePath, bFilePath := range r.configsToRestore {
		if com.IsFile(bFilePath) {
			if err := os.Remove(bFilePath); err != nil {
				r.Logger.Error(fmt.Sprintf("could not remove file '%s' on reverter commit: %v", bFilePath, err))
			}
		}

		delete(r.configsToRestore, filePath)
	}

	r.configsToDelete = nil

	return nil
}

func (r *Reverter) getBackupConfigPath(filePath string) string {
	return filePath + ".back"
}
