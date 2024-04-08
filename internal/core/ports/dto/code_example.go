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

type GetCustomCodeExampleDto struct {
	UUID                    string `json:"UUID"`
	Content                 string `json:"content"`
	ProgrammingLanguageUUID string `json:"programmingLanguageUUID"`
	UserID                  string `json:"userID"`
}

type CreateCodeExampleDto struct {
	Content                 string `json:"content"`
	ProgrammingLanguageUUID string `json:"programmingLanguageUUID"`
}
