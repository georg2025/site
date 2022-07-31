package start_server

import (
  "database/sql"
  "flag"
  "log"
  "net/http"
  "path/filepath"
  "os"
  "html/template"
  _ "github.com/go-sql-driver/mysql"
  "site/pkg/models/mySQL"
)
//Структура, позволяющая сделать разные объекты доступным в других функциях
type application struct {
  errorLog *log.Logger
  infoLog  *log.Logger
  add1 *mysql.Add1Model
  templateCache map[string]*template.Template
}

func Start() {
    //сетевой адрес, на котором запускается сервак. По умолчанию 8080
    addr := flag.String("addr", ":8080", "Сетевой адрес HTTP")
    //старт соединения с mysql сервером
    dsn := flag.String("dsn", "web:localhost@/site?parseTime=true", "Название MySQL источника данных")

    flag.Parse()
    //Это логеры для выведения ошибок и информации
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
 	  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    //dsn передается в функцию openDB
    db, err := openDB(*dsn)
      if err != nil {
      errorLog.Fatal(err)
    }
    //на всякий пожарный - дефер
    defer db.Close()
    // Кэш для хранения шаблонов
    templateCache, err := newTemplateCache("./ui/html/")
    if err != nil {
        errorLog.Fatal(err)
    }


    //Структура для логгеров
    app := &application{
      errorLog: errorLog,
      infoLog:  infoLog,
      add1: &mysql.Add1Model{DB: db},
      templateCache: templateCache,
    }


    //Структура с параметрами конфигурации для сервака. Через нее потом запускается listenAndServe
    srv := &http.Server{
      Addr:     *addr,
      ErrorLog: errorLog,
      Handler:  app.routes(),
    }

    //принт этот не нужен. Можно его убрать при запуске. Данная функция запускает сервак
    infoLog.Printf("Запуск веб-сервера на %s", *addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}

// Ниже тип и функция нужны для того, чтобы нельзя было извне получить доступ к статическимм файлам. Она не дает загружать html код, если в папке нет файла index.html, либо, если на файл явно не ссылаются
type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }

            return nil, err
        }
    }

    return f, nil
}

//функция для открытия соединений с sql
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
