package start_server

import "net/http"

func (app *application) routes() *http.ServeMux {
    //здесь создается новый сервак и ведется обработка различных страниц. Для каждой страницы - свой HandleFunc
    mux := http.NewServeMux()
    mux.HandleFunc("/", app.Home)
    mux.HandleFunc("/ad", app.Show_ad)
    mux.HandleFunc("/new_ad", app.Ad_create)
    mux.HandleFunc("/latest", app.Latest)

    // далее 2 обработчика статичный файлов из 2 разных папок
    fileServer := http.FileServer(neuteredFileSystem {http.Dir("./ui/static/themes")})
    mux.Handle("/ui/static/themes", http.NotFoundHandler())
    mux.Handle("/themes/", http.StripPrefix("/themes", fileServer))

    fileServer1 := http.FileServer(neuteredFileSystem {http.Dir("./ui/static/bootstrap")})
    mux.Handle("/ui/static/bootstrap", http.NotFoundHandler())
    mux.Handle("/bootstrap/", http.StripPrefix("/bootstrap", fileServer1))

    return mux
}
