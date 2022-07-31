package start_server

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

//Помощник отправляет пользователю ошибку 500 (внутренняя ошибка)
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//Этот для других 400-х ошибок (кроме 404)
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Этот отправляет 404 - страница не найдена
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
