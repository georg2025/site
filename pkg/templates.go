package start_server

import (
  "site/pkg/models"
  "html/template"
  "path/filepath"
)

// Тип templateData содержит данные для передачи в HTML шаблон. Т.к. шаблон может
// принимать данные только из 1 источника - он нужен для аккумуляции данных.
type templateData struct {
    Add1 *models.Add1
    Add2 []*models.Add1
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
    // Создаем мапу для хранения кэша
    cache := map[string]*template.Template{}

    // Здесь мы получаем срез всех файлов с расширением '.page.tmpl'. То-есть, по сути, всех страниц приложения.
    pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
    if err != nil {
        return nil, err
    }

    // Перебираем файл шаблона от каждой страницы.
    for _, page := range pages {
        // Извлечение конечное названия файла из полного пути к файлу
        // и присваивание его переменной name.
        name := filepath.Base(page)

        // Обрабатываем итерируемый файл шаблона.
        ts, err := template.ParseFiles(page)
        if err != nil {
            return nil, err
        }

        // Используем метод ParseGlob для добавления всех шаблонов
        ts, err = ts.ParseGlob(filepath.Join(dir, "*.ad.tmpl"))
        if err != nil {
            return nil, err
        }

        // Добавляем полученный набор шаблонов в кэш, используя название страницы
        // в качестве ключа для  мапы.
        cache[name] = ts
    }

    // Возвращаем полученную мапу.
    return cache, nil
}
