package config

type testConfig struct {
	SourceDir    string
	OutputDir    string
	AnswerKeyDir string
}

func TestConfig() *testConfig {
	this := testConfig{
		SourceDir:    "../testing-resources/process_files_test/source_template",
		OutputDir:    "../testing-resources/process_files_test/test_output",
		AnswerKeyDir: "../testing-resources/process_files_test/answer_key",
	}
	return &this
}
