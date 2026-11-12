package tests

import "os"

var Port = 7540
var DBFile = "./../database/scheduler.db"
var FullNextDate = true
var Search = true

// Token — JWT для тестов с включённой авторизацией.
// Получите его через POST /api/signin и передайте в TODO_TOKEN.
var Token = os.Getenv("TODO_TOKEN")
