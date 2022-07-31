package start_server

import (
  "errors"
  "strconv"
  "net/http"
  "fmt"
  "html/template"
  "site/pkg/models"

)

//Домашняя страница

func (app *application) Home (w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    app.notFound(w)
    return
  }
  files:=[]string {
    "./ui/html/main_page.tmpl",
    "./ui/html/bricks/carousel.tmpl",
    "./ui/html/bricks/sidebar.tmpl",
    "./ui/html/bricks/header.tmpl",
    "./ui/html/bricks/footer.tmpl",
    "./ui/html/bricks/latest_ad.tmpl",
    "./ui/html/bricks/featured.tmpl",
  }
  ts, err := template.ParseFiles(files...)
	if err != nil {
    app.serverError(w, err)
  }
	err = ts.Execute(w, nil)
	if err != nil {
    app.serverError(w, err)
	}
}

//страница для создания объявления. Метод - только post

func (app *application) Ad_create(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }
    // Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
  	title := "История про улитку"
  	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
  	expired := "7"
    adress := "Вологда"
    var region uint8
    region = 22
    price := 410000

  	// Передаем данные в метод Add1Model.Insert(), получая обратно
  	// ID только что созданной записи в базу данных.
  	id, err := app.add1.Insert(title, content, expired, adress, region, price)
  	if err != nil {
  		app.serverError(w, err)
  		return
  	}

  	// Перенаправляем пользователя на соответствующую страницу заметки.
  	http.Redirect(w, r, fmt.Sprintf("/ad?id=%d", id), http.StatusSeeOther)

}

//страница для показа объявлений. Принимает html запрос с id объявления и выдает инфу конкретного объявления

func (app *application) Show_ad(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    s, err := app.add1.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }
    data := &templateData{Add1: s}
    // Инициализируем срез, содержащий путь к файлу show_ad.tmpl
    files := []string{
        "./ui/html/show_ad.tmpl",
        "./ui/html/bricks/sidebar.tmpl",
        "./ui/html/bricks/header.tmpl",
        "./ui/html/bricks/footer.tmpl",
    }

    // Парсинг файлов шаблонов...
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // А затем выполняем их.
    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }

}

func (app *application) Latest(w http.ResponseWriter, r *http.Request) {
    s, err := app.add1.Latest()
    if err != nil {
      app.serverError(w, err)
      return
    }

    for _, add := range s {
      fmt.Fprintf(w, "%v\n", add)
    }
}
