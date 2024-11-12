module taskbot

go 1.19

// это надо для переопределения адреса сервера
// в оригинале урл телеграма в константе, у меня там строка, которую я переопределяю в тесте
replace gopkg.in/telegram-bot-api.v4 => ./local/telegram-bot-api

//replace github.com/go-telegram-bot-api/telegram-bot-api => ./local/telegram-bot-api

require gopkg.in/telegram-bot-api.v4 v4.6.4

require github.com/technoweenie/multipartstreamer v1.0.1 // indirect
