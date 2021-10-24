package config

type testConfig struct {
	ProcessFilesTestSource string
	ProcessFilesTestOutput string
	ProcessFilesTestAnswer string
	ValidationTestSource   string
	ValidationTestOutput   string
	ValidationTestAnswer   string
}

func TestConfig() *testConfig {
	this := testConfig{
		ProcessFilesTestSource: "../.test-project/process_files_test/source_template",
		ProcessFilesTestOutput: "../.test-project/process_files_test/test_output",
		ProcessFilesTestAnswer: "../.test-project/process_files_test/answer_key",
		ValidationTestSource:   "../.test-project/validation_test/source_template",
		ValidationTestOutput:   "../.test-project/validation_test/test_output",
		ValidationTestAnswer:   "../.test-project/validation_test/answer_key",
	}
	return &this
}
