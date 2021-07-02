package handlers

import (
	"errors"
	"fmt"

	"github.com/Aitugan/CodingChallenge/pkg/lib"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (wm *WorksManager) getClientID(ctx *gin.Context) (uuid.UUID, error) {
	clientIDStr := ctx.Request.TLS.VerifiedChains[0][0].Subject.CommonName
	if clientIDStr == "" {
		return uuid.Nil, errors.New("The client certificate's common name should be a uuid")
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return uuid.Nil, errors.New("Common name of the client certificate must be a valid uuid")
	}
	return clientID, nil
}

type StartCommandRequest struct {
	Command []string `json:"command"`
}

func (wm *WorksManager) Start(ctx *gin.Context) {
	var startCommandReq StartCommandRequest
	if err := ctx.ShouldBindJSON(&startCommandReq); err != nil {
		ctx.JSON(401, "Incorrect body format")
		return
	}
	clientID, err := wm.getClientID(ctx)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}
	work := lib.NewWork(clientID, startCommandReq.Command)
	wm.works[WorkKey{work.ID, clientID}] = work
	var processErr error
	go func() {
		work.Process.Stdout = &work.Output
		work.Process.Stderr = &work.Output
		err = work.Process.Start()
		if err != nil {
			ctx.JSON(400, "Some error")
			return
		}
		work.Status = lib.ActiveStatus
		processErr = work.Process.Wait()
		if processErr == nil {
			work.Status = lib.KilledStatus
			return
		}
		work.Status = lib.ExitedStatus
	}()
	ctx.JSON(200, work.ID)
}

func (wm *WorksManager) Log(ctx *gin.Context) {
	cmdID, err := uuid.Parse(ctx.Param("id"))

	clientID, err := wm.getClientID(ctx)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}

	if _, ok := wm.works[WorkKey{cmdID, clientID}]; !ok {
		ctx.JSON(400, "Such work does not exist in your history")
		return
	}

	ctx.JSON(200, wm.works[WorkKey{cmdID, clientID}].Output.String())
}

func (wm *WorksManager) Stop(ctx *gin.Context) {
	cmdID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, "Work id must be a valid uuid")
		return
	}

	clientID, err := wm.getClientID(ctx)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}

	if _, ok := wm.works[WorkKey{cmdID, clientID}]; !ok {
		ctx.JSON(400, "Such work does not exist in your history")
		return
	}
	if wm.works[WorkKey{cmdID, clientID}].Status != lib.ActiveStatus {
		ctx.JSON(400, "This work is not running already")
		return
	}
	if err := wm.works[WorkKey{cmdID, clientID}].Process.Process.Kill(); err != nil {
		ctx.JSON(400, "Failed to kill the process")
		return
	}
	wm.works[WorkKey{cmdID, clientID}].Status = lib.StoppedStatus
	fmt.Println(wm.works[WorkKey{cmdID, clientID}])
	ctx.JSON(200, "The process is killed")
}

func (wm *WorksManager) QueryStatus(ctx *gin.Context) {
	cmdID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, "Work id must be a valid uuid")
		return
	}

	clientID, err := wm.getClientID(ctx)
	if err != nil {
		ctx.JSON(401, err.Error())
		return
	}

	if _, ok := wm.works[WorkKey{cmdID, clientID}]; !ok {
		ctx.JSON(400, "Such work does not exist in your history")
		return
	}
	fmt.Println(wm.works[WorkKey{cmdID, clientID}].Status)
	ctx.JSON(200, "The process status is "+wm.works[WorkKey{cmdID, clientID}].Status)
}
