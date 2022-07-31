package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

//структура данных объявлений из таблицы add1 mysql
type Add1 struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expired time.Time
  Region uint8
  Price int
  Adress string
}
