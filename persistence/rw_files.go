package persistence

import (
	"os"
)

// Esse módulo é responssável pelo controle de arquivos
func ReadFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	return file, err
}

// sobreescreve o arquivo indicado com os dados enviados
func OverwriteFile(filename string, data []byte) (int, error) {
	file, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close() // só posso fechar depois de confirmar que não teve erros
	b, err := file.Write(data)
	return b, err
}
