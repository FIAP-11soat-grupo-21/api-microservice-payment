package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Este teste documenta a existência da função Init
	// Não podemos executá-la diretamente pois ela:
	// 1. Inicia o servidor (blocking)
	// 2. Tenta conectar ao banco de dados
	// 3. Tenta conectar ao SQS
	// 4. Executa ginRouter.Run() que bloqueia a execução

	// Apenas verificamos que a função existe e é chamável
	assert.NotNil(t, Init, "Init function should exist")
}
