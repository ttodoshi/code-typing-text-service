# Speed Typing text service on Golang

speed typing text generation service written on golang with GORM and gin framework

endpoints:

- **/api/v1/texts/** (to get all available regular texts)
- **/api/v1/programming-languages/** (to get available programming languages)
- **/api/v1/code-examples/** (to get all available code examples)
- **/api/v1/code-examples/?programming-language-uuid=?** (to get all available code examples by requested language)
