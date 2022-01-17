/*
Copyright 2022 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/koderover/zadig/pkg/microservice/aslan/core/label/repository/models"
	"github.com/koderover/zadig/pkg/microservice/aslan/core/label/service"
	internalhandler "github.com/koderover/zadig/pkg/shared/handler"
	e "github.com/koderover/zadig/pkg/tool/errors"
)

func ListLabels(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	listLabelsArgs := make([]*service.ListLabelsArgs, 0)
	if err := c.ShouldBindJSON(&listLabelsArgs); err != nil {
		ctx.Err = err
		return
	}
	ctx.Resp, ctx.Err = service.ListLabels(listLabelsArgs)
}

func createLabelValidate(lbs []*models.Label) error {
	keyValues := sets.NewString()
	for _, v := range lbs {
		if v.Key == "" || v.Value == "" {
			return e.ErrInvalidParam.AddDesc("invalid label args")
		}
		keyValue := fmt.Sprintf("%s-%s", v.Key, v.Value)
		if keyValues.Has(keyValue) {
			return e.ErrInvalidParam.AddDesc(fmt.Sprintf("duplicate key-value:%s", keyValue))
		}
		keyValues.Insert(keyValue)
	}
	return nil
}

func CreateLabels(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	labels := make([]*models.Label, 0)
	if err := c.ShouldBindJSON(&labels); err != nil {
		ctx.Err = err
		return
	}

	if err := createLabelValidate(labels); err != nil {
		ctx.Err = err
		return
	}
	for _, v := range labels {
		v.CreateBy = ctx.UserName
	}
	ctx.Err = service.CreateLabels(labels)
}

//DeleteLabels  can only bulk delete labels which not bind reousrces
func DeleteLabels(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	deleteLabelsArgs := new(service.DeleteLabelsArgs)
	if err := c.ShouldBindJSON(deleteLabelsArgs); err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc("json bind fail")
		return
	}
	force := c.Query("force")
	forceBool, err := strconv.ParseBool(force)
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Err = service.DeleteLabels(deleteLabelsArgs.IDs, forceBool, ctx.Logger)
}

func DeleteLabel(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	id := c.Param("id")
	if id == "" {
		ctx.Err = e.ErrInvalidParam.AddDesc("id must not be empty")
		return
	}
	force := c.Query("force")
	forceBool, err := strconv.ParseBool(force)
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Err = service.DeleteLabel(id, forceBool, ctx.Logger)
}

func ListResourceByLabels(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	listResourceByLabelsReq := new(service.ListResourceByLabelsReq)
	if err := c.ShouldBindJSON(listResourceByLabelsReq); err != nil {
		ctx.Err = err
		return
	}
	ctx.Resp, ctx.Err = service.ListResourcesByLabels(listResourceByLabelsReq.LabelFilters, ctx.Logger)
}

type ListLabelsByResourceReq struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

func ListLabelsByResource(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	listLabelsByResourceReq := new(ListLabelsByResourceReq)
	if err := c.ShouldBindJSON(listLabelsByResourceReq); err != nil {
		ctx.Err = err
		return
	}
	ctx.Resp, ctx.Err = service.ListLabelsByResourceID(listLabelsByResourceReq.ResourceID, listLabelsByResourceReq.ResourceType, ctx.Logger)
}
