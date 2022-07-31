package start_server

import "site/pkg/models"

// Тип templateData содержит данные для передачи в HTML шаблон. Т.к. шаблон может
// принимать данные только из 1 источника - он нужен для аккумуляции данных.
type templateData struct {
    Add1 *models.Add1
}
