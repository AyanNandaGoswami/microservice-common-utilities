package main

import (
	"fmt"

	"github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/models"
	"github.com/AyanNandaGoswami/file-sharing-app-common-utilities/v1/outsource"
)

func main() {
	requestBody := models.PermissionValidadtionRequest{
		ValidateBy:      "primitiveUserId",
		PrimitiveUserId: "67b0e9c917ddc6790e248882",
		RequestedUrl:    "/authorization/v1/permission/add/p",
		RequestedMethod: "POST",
	}
	err := outsource.HasPermission(&requestBody)
	if err != nil {
		fmt.Println(err.Error())
	}
}
