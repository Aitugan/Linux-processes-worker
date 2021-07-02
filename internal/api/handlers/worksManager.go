package handlers

import (
	"github.com/Aitugan/CodingChallenge/pkg/lib"
	"github.com/google/uuid"
)

type WorkKey struct {
	WorkID, ClientID uuid.UUID
}

type WorksManager struct {
	works map[WorkKey]*lib.Work
}

func NewWorksManager() *WorksManager {
	return &WorksManager{
		works: make(map[WorkKey]*lib.Work),
	}
}
