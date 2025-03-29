package fileimports

import (
	"fmt"

	"github.com/google/uuid"
)

func getFileImportPath(organizationId uuid.UUID, fileUploadId uuid.UUID, fileName string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", ORGANIZATIONS_SUBPATH, organizationId, FILE_IMPORT_SUBPATH, fileUploadId, fileName)
}
