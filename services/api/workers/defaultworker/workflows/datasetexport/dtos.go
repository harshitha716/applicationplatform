package datasetexport

import (
	"github.com/google/uuid"
)

type baseParams struct {
	userId uuid.UUID
	orgIds []uuid.UUID
}

func (p baseParams) GetAccessControlParams() (uuid.UUID, []uuid.UUID) {
	return p.userId, p.orgIds
}
