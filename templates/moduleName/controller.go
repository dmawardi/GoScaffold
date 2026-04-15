package moduleName

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/goTemplate/internal/common"
	"github.com/gin-gonic/gin"
)

type ModuleNameController interface {
	FindAll(c *gin.Context)
	Find(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type moduleNameController struct {
	service ModuleNameService
}

func NewModuleNameController(service ModuleNameService) ModuleNameController {
	return &moduleNameController{service}
}

// Used to init the query params for easy extraction in controller
// Returns: map[string]string{"name": "string"}
func ModuleNameConditionQueryParams() map[string]string {
	return map[string]string{
		"name": "string",
	}
}

// API/MODULENAMES
// @Summary      Find a list of moduleNames
// @Description  Accepts limit, offset, order, search (added as non-case sensitive LIKE) and field names (eg. name=) query parameters to find a list of moduleNames. Search is applied to all string fields.
// @Tags         ModuleName
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        offset   query      int  false  "offset"
// @Param        order   query      int  false  "order by eg. (asc) \"id\" (desc) \"id_desc\" )"
// @Param        search   query      string  false  "search (added to all string conditions as LIKE SQL search)"
// @Param        name     query      string  false  "name"
// @Success      200 {object} common.PaginatedResult[ModuleNameResponse]
// @Failure      400 {string} string "Can't find moduleNames"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Failure      400 {string} string "Error extracting query params"
// @Router       /moduleNames [get]
func (c moduleNameController) FindAll(ctx *gin.Context) {
	// Grab basic query params set defaults as needed
	baseQueryParams, err := common.ExtractBasicFindAllQueryParams(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error extracting query params"})
		return
	}

	// Generate query params to extract
	queryParamsToExtract := ModuleNameConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := common.ExtractSearchAndConditionParams(ctx.Request, queryParamsToExtract)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error extracting query params"})
		return
	}

	// Query database for all moduleNames using query params
	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Can't find moduleNames"})
		return
	}

	ctx.JSON(http.StatusOK, found)
}

// @Summary      Find ModuleName
// @Description  Find a moduleName by ID
// @Tags         ModuleName
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ModuleName ID"
// @Success      200 {object} ModuleNameResponse
// @Failure      400 {string} string "Can't find moduleName with ID: {id}"
// @Router       /moduleNames/{id} [get]
func (c moduleNameController) Find(ctx *gin.Context) {
	// Grab URL parameter
	stringParameter := ctx.Param("id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	found, err := c.service.FindById(idParameter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Can't find moduleName with ID: %d", idParameter)})
		return
	}

	ctx.JSON(http.StatusOK, found)
}

// @Summary      Create ModuleName
// @Description  Create a moduleName
// @Tags         ModuleName
// @Accept       json
// @Produce      json
// @Param        ModuleName  body      CreateModuleName  true  "ModuleName entity that needs to be added"
// @Success      201 {object} map[string]interface{} "message: ModuleName creation successful!"
// @Failure      400 {string} string "Failed to create moduleName"
// @Failure      400 {string} string "Invalid JSON"
// @Router       /moduleNames [post]
func (c moduleNameController) Create(ctx *gin.Context) {
	var createData CreateModuleName
	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&createData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Create moduleName
	created, err := c.service.Create(&createData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create moduleName"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ModuleName creation successful!",
		"data":    created,
	})
}

// @Summary      Update ModuleName
// @Description  Update a moduleName
// @Tags         ModuleName
// @Accept       json
// @Produce      json
// @Param        id          path      int                 true  "ModuleName ID"
// @Param        ModuleName  body      UpdateModuleName    true  "ModuleName entity that needs to be updated"
// @Success      200 {object} map[string]interface{} "message: ModuleName update successful!"
// @Failure      400 {string} string "Can't find moduleName with ID: {id}"
// @Failure      400 {string} string "Invalid ID"
// @Failure      400 {string} string "Invalid JSON"
// @Router       /moduleNames/{id} [put]
func (c moduleNameController) Update(ctx *gin.Context) {
	// Grab URL parameter
	stringParameter := ctx.Param("id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updateData UpdateModuleName
	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Update moduleName
	updated, err := c.service.Update(idParameter, &updateData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Can't find moduleName with ID: %d", idParameter)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ModuleName update successful!",
		"data":    updated,
	})
}

// @Summary      Delete ModuleName
// @Description  Delete a moduleName
// @Tags         ModuleName
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ModuleName ID"
// @Success      200 {object} map[string]interface{} "message: ModuleName deletion successful!"
// @Failure      400 {string} string "Can't find moduleName with ID: {id}"
// @Failure      400 {string} string "Invalid ID"
// @Router       /moduleNames/{id} [delete]
func (c moduleNameController) Delete(ctx *gin.Context) {
	// Grab URL parameter
	stringParameter := ctx.Param("id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Delete moduleName
	err = c.service.Delete(idParameter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Can't find moduleName with ID: %d", idParameter)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ModuleName deletion successful!"})
}