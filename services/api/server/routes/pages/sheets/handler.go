package sheets

import (
	"net/http"

	"github.com/Zampfi/application-platform/services/api/core/sheets"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/server/routes/pages/sheets/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleGetSheetsAll(c *gin.Context, sheetService sheets.SheetsService) {

	pageIdStr := c.Param("pageId")
	pageId, err := uuid.Parse(pageIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid page id"})
		return
	}

	sheets, err := sheetService.GetSheetsByPageId(c, pageId)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, sheets)
}

func handleGetSheetByID(c *gin.Context, sheetService sheets.SheetsService) {
	sheetIDStr := c.Param("sheetId")
	sheetID, err := uuid.Parse(sheetIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid sheet id"})
		return
	}

	sheet, err := sheetService.GetSheetById(c, sheetID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	sheetResponse := dtos.SheetResponse{}
	err = sheetResponse.NewSheetResponse(sheet)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, sheetResponse)
}

func handleGetSheetFilters(c *gin.Context, sheetService sheets.SheetsService) {
	_, _, orgIds := apicontext.GetAuthFromContext(c)
	if len(orgIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no organization ids found"})
		return
	}

	orgId := orgIds[0]

	sheetIdStr := c.Param("sheetId")
	_, err := uuid.Parse(sheetIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid sheet id"})
		return
	}

	sheetId, err := uuid.Parse(sheetIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid sheet id"})
		return
	}

	filterConfig, err := sheetService.GetSheetFilterConfigById(c, orgId, sheetId)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	sheetFilterConfigDto := dtos.SheetFilterConfig{}
	sheetFilterConfigDto.NewSheetFilterConfig(filterConfig)

	c.JSON(200, sheetFilterConfigDto)
}
