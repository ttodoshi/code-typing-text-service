package dto

type GetProgrammingLanguageDto struct {
	UUID string `json:"UUID"`
	Name string `json:"name"`
}

type GetCodeExampleDto struct {
	UUID                    string `json:"UUID"`
	Content                 string `json:"content"`
	ProgrammingLanguageUUID string `json:"programmingLanguageUUID"`
}
