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

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    // Здесь мы извлекаем из кэша шаблон. Если в кэше нет такого шаблона - вызывается метод ошибки serverError().
	  ts, ok := app.templateCache[name]
    if !ok {
        app.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
        return
    }

    // Рендерим файлы шаблона, передавая динамические данные из переменной `td`.
    err := ts.Execute(w, td)
    if err != nil {
        app.serverError(w, err)
    }
}
