package password

import (
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err1 := GenerateSalt(16)
	salt2, err2 := GenerateSalt(16)

	if err1 != nil || err2 != nil {
		t.Fatalf("Ошибка генерации соли: %v, %v", err1, err2)
	}

	if salt1 == salt2 {
		t.Errorf("Сгенерированная соль не должна быть одинаковой: %s == %s", salt1, salt2)
	}

	if len(salt1) != 32 || len(salt2) != 32 { // 16 байт = 32 символа в hex
		t.Errorf("Длина соли должна быть 32 символа, но получено: %d и %d", len(salt1), len(salt2))
	}
}

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "securepassword"

	// Генерация хэша и соли
	hash, salt, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Ошибка хеширования пароля: %v", err)
	}

	// Проверка, что пароль проходит валидацию
	if !CheckPassword(password, salt, hash) {
		t.Errorf("CheckPassword должен возвращать true для корректного пароля")
	}

	// Проверка с неправильным паролем
	if CheckPassword("wrongpassword", salt, hash) {
		t.Errorf("CheckPassword должен возвращать false для неверного пароля")
	}

	// Проверка с неправильной солью
	if CheckPassword(password, "wrongSalt", hash) {
		t.Errorf("CheckPassword должен возвращать false для неверной соли")
	}

	// Проверка с неправильным хэшем
	if CheckPassword(password, salt, "wrongHash") {
		t.Errorf("CheckPassword должен возвращать false для неверного хэша")
	}
}
