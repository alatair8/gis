// Package util содержит вспомогательные функции, не относящиеся к конкретному домену.
package util

import (
	"crypto/rand"
	"encoding/hex"
	"sync/atomic"
)

var counter uint64

// NewID генерирует псевдо-уникальный идентификатор в шестнадцатеричном представлении.
//
// Для демонстрационной системы достаточно комбинации случайных байтов и счётчика.
// В промышленной реализации следует использовать UUID, выдаваемый СУБД или сервисом генерации идентификаторов.
func NewID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		// В крайне маловероятном случае ошибки генерации случайных данных используем только счётчик.
		return hex.EncodeToString([]byte{byte(atomic.AddUint64(&counter, 1))})
	}

	c := atomic.AddUint64(&counter, 1)
	// Добавляем счётчик в конце, чтобы снизить вероятность коллизий при старте.
	idBytes := append(buf, byte(c>>24), byte(c>>16), byte(c>>8), byte(c))
	return hex.EncodeToString(idBytes)
}
