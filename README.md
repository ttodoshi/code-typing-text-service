# Speed Typing text service on Golang

speed typing text generation service written on golang with GORM and gin framework

endpoints:

- **GET /api/v1/texts/** (to get all available regular texts)
- **GET /api/v1/texts/programming-languages/** (to get available programming languages)
- **GET /api/v1/texts/code-examples/** (to get all available code examples)
- **GET /api/v1/texts/code-examples/?programming-language-name=?** (to get all available code examples by requested
  language)
