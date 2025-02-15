package configload

import (
	"os"
	"testing"
)

const testConfig = `
server:
  port: 8080

database:
  host: "localhost"
  port: "5432"
  user: "testuser"
  password: "testpassword"
  dbname: "testdb"
  sslmode: "disable"

jwt:
  secret: "testsecret"
`

func TestLoadConfig(t *testing.T) {
	// Создаем временный файл для тестовой конфигурации
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удаляем после теста

	// Записываем тестовую конфигурацию
	if _, err := tmpFile.Write([]byte(testConfig)); err != nil {
		t.Fatalf("Ошибка записи в файл: %v", err)
	}
	tmpFile.Close() // Закрываем файл перед использованием

	// Загружаем конфигурацию
	config := LoadConfig(tmpFile.Name())

	// Проверяем значения
	if config.Server.Port != 8080 {
		t.Errorf("Ожидался порт 8080, но получен %d", config.Server.Port)
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Ожидался host 'localhost', но получен %s", config.Database.Host)
	}
	if config.Database.Port != "5432" {
		t.Errorf("Ожидался port '5432', но получен %s", config.Database.Port)
	}
	if config.Database.User != "testuser" {
		t.Errorf("Ожидался user 'testuser', но получен %s", config.Database.User)
	}
	if config.Database.Password != "testpassword" {
		t.Errorf("Ожидался пароль 'testpassword', но получен %s", config.Database.Password)
	}
	if config.Database.DBName != "testdb" {
		t.Errorf("Ожидался dbname 'testdb', но получен %s", config.Database.DBName)
	}
	if config.Database.SSLMode != "disable" {
		t.Errorf("Ожидался sslmode 'disable', но получен %s", config.Database.SSLMode)
	}
	if config.JWT.Secret != "testsecret" {
		t.Errorf("Ожидался jwt secret 'testsecret', но получен %s", config.JWT.Secret)
	}
}
