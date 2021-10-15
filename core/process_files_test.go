package core

import (
	"github.com/oledakotajoe/clonr/config"
	"github.com/oledakotajoe/clonr/types"
	log "github.com/sirupsen/logrus"
	"testing"
)

func Test_setup(t *testing.T) {
	// This is not really a test, it is just a setup function.
	config.ConfigureLogger()
}

func Test_givenTemplateFile_getFileMapFromTemplate(t *testing.T) {
	sourceDir := config.TestConfig().SourceDir

	testSettings := &types.FileProcessorSettings{
		ConfigFilePath: sourceDir,
		StringInputReader: func(input string) string {
			return input
		},
	}

	processClonrConfig(testSettings)
	exampleFileMap := types.FileMap{
		"../.testing-.resources/process_files_test/source_template/sub-dir/another-test.txt": types.ClonrVarMap{"file_sub_dir_multi_diff_1": "file_sub_dir_multi_diff_1", "file_sub_dir_multi_diff_2": "file_sub_dir_multi_diff_2"},
		"../.testing-.resources/process_files_test/source_template/test.txt":                 types.ClonrVarMap{"file_in_root_multi": "file_in_root_multi"},
	}

	for key, value := range testSettings.MainTemplateMap {
		log.Infof("key: %s, value: %s", key, value)
		if exampleFileMap[key] == nil {
			t.Fatalf("Maps were not equivalent")
		}
		for k, v := range value {
			log.Infof("key: %s, value: %s", k, v)
			if exampleFileMap[key][k] != v {
				t.Fatalf("Maps were not equivalent")
			}
		}
	}
}
