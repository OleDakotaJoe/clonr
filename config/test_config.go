package config

type testConfig struct {
	SourceDir    string
	OutputDir    string
	AnswerKeyDir string
}

func TestConfig() *testConfig {
	this := testConfig{
		SourceDir:    "../.test-project/process_files_test/source_template",
		OutputDir:    "../.test-project/process_files_test/test_output",
		AnswerKeyDir: "../.test-project/process_files_test/answer_key",
	}
	return &this
}
