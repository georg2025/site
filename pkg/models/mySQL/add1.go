package mysql

import (
	"errors"
	"database/sql"
	"site/pkg/models"
)

// Делаем тип, который обертывает пул подключений к БД
type Add1Model struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *Add1Model) Insert(title, content, expired, adress string, region uint8, price int) (int, error) {
	// Это сам SQL запрос, который нужно выполнить
	// Все данные передаются через знак ?, чтобы исключить возможность sql-инъекций
	stmt := `INSERT INTO add1 (title, content, created, expired, adress, region, price)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?, ?, ?)`

// Дальше с помощью Exec() передается сам запрос (stmt) и информация о нем. Метод возвращает result с информацией о запросе
	result, err := m.DB.Exec(stmt, title, content, expired, adress, region, price)
	if err != nil {
		return 0, err
	}

	// С помощью метода LastInsertId() узнаем айдишник послденей записи (которую только что добавили)
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Возвращаем айдишник и конвертируем его в int.
	return int(id), nil
	}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *Add1Model) Get(id int) (*models.Add1, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT id, title, content, created, expired, adress, region, price FROM add1
    WHERE expired > UTC_TIMESTAMP() AND id = ?`

	// Метод QueryRow() выполняет SQL запрос, id передается не напрямую, а через знак ?,
	// чтобы исключить возможность внедрения ненадежных данных

	row := m.DB.QueryRow(stmt, id)

	// Инициализируем указатель на новую структуру.
	s := &models.Add1{}

	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в
	// соответствующее поле в структуре Add1.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired, &s.Adress, &s.Region, &s.Price)
	if err != nil {
		// Проверка на наличие ошибок
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект.
	return s, nil
}

// Latest - Метод возвращает последние.
func (m *Add1Model) Latest() ([]*models.Add1, error) {
	// SQL запрос, который хотим выполнить.
	stmt := `SELECT id, title, content, created, expired, region, price, adress FROM add1
	WHERE expired > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 9`

	// Передаем запрос через метод Query()
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Создаем дефер rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest().
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Add1
	var adds []*models.Add1

	// Используем rows.Next() для перебора результатов.
	for rows.Next() {
		// Создаем указатель на новую структуру Add1
		s := &models.Add1{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Опять же, аргументы предоставленные в row.Scan()
		// должны быть указателями на место, куда требуется скопировать данные и
		// количество аргументов должно быть точно таким же, как количество
		// столбцов из таблицы базы данных, возвращаемых вашим SQL запросом.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired, &s.Region, &s.Price, &s.Adress)
		if err != nil {
			return nil, err
		}
	// Добавляем структуру в срез.
		adds = append(adds, s)
	}

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return adds, nil

}
