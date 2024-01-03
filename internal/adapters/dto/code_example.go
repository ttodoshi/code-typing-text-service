package dto

type GetProgrammingLanguageDto struct {
	UUID string `json:"UUID"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type GetCodeExampleDto struct {
	UUID                    string `json:"UUID"`
	Content                 string `json:"content"`
	ProgrammingLanguageUUID string `json:"programmingLanguageUUID"`
}
