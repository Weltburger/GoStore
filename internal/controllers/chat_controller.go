package controllers

import (
	"GoStore/internal/chat"
	"GoStore/pkg/models"
	"GoStore/pkg/statuserror"
	"context"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"sync"
)

type ChatController struct {
	controller *Controller
	hubCash sync.Map
	//channel chan uuid.UUID
	//chatUUID uuid.UUID
}

func (chatController *ChatController) ListenSocket(c *gin.Context) {
	hub := chat.NewHub()
	go hub.Run()
	chatController.hubCash.Store(uuid.NewV1(), hub)
	chat.ServeWs(hub, c.Writer, c.Request, "Guest")
}

func (chatController *ChatController) ShowChats(c *gin.Context) {
	var arr []uuid.UUID
	chatController.hubCash.Range(func(k, v interface{}) bool {
		arr = append(arr, k.(uuid.UUID))
		return true
	})

	c.JSON(http.StatusOK, gin.H{
		"value":   arr,
	})
}

func (chatController *ChatController) SelectChat(c *gin.Context) {
	param := struct {
		UUID string `form:"uuid"`
	}{}

	if err := c.ShouldBind(&param); err != nil {
		panic(err)
	}

	value, ok := c.Get("user")
	if !ok {
		panic(statuserror.NullUser)
	}

	user := value.(*models.User)

	v, ok := chatController.hubCash.Load(uuid.Must(uuid.FromString(param.UUID)))
	if ok {
		hub := v.(*chat.Hub)
		chat.ServeWs(hub, c.Writer, c.Request, user.UUID.String())
	}
}

/*func (chatController *ChatController) ConnectToChat(c *gin.Context) {
	//UUID := <- chatController.channel
	//close(chatController.channel)
	v, ok := chatController.hubCash.Load(chatController.chatUUID)
	if ok {
		hub := v.(*chat.Hub)
		chat.ServeWs(hub, c.Writer, c.Request)
	}
}*/

func (chatController *ChatController) RemoveHub(c *gin.Context) {
	param := struct {
		UUID uuid.UUID `json:"uuid" binding:"required"`
	}{}
	if err := c.Bind(&param); err != nil {
		panic(statuserror.New(404, statuserror.StatusNotFilled, err))
	}

	//////////////////////////////////////////
	v, ok := chatController.hubCash.Load(param.UUID)
	if !ok { return }
	hub := v.(*chat.Hub)
	/*if ok {
		hub := v.(*chat.Hub)
		hub.CloseFile()
	}*/
	ctx := context.Background()
	_, err := chatController.controller.store.ChatRepository().InsertChat(ctx, &models.Chat{
		UUID:     hub.GetUUID(),
		Data:    hub.History(),
	})
	if err != nil {
		panic(statuserror.New(500, statuserror.StatusCodeServerErr, err))
	}
	//////////////////////////////////////////

	chatController.hubCash.Delete(param.UUID)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Hub has been deleted",
		"value": param.UUID,
	})
}

func (chatController *ChatController) RemoveInactiveHubs(c *gin.Context) {
	chatController.hubCash.Range(func(k, v interface{}) bool {
		hub := v.(*chat.Hub)
		//hub.CloseFile()

		ctx := context.Background()
		_, err := chatController.controller.store.ChatRepository().InsertChat(ctx, &models.Chat{
			UUID:     hub.GetUUID(),
			Data:    hub.History(),
		})
		if err != nil {
			panic(statuserror.New(500, statuserror.StatusCodeServerErr, err))
		}

		if len(hub.Clients()) == 0 {
			chatController.hubCash.Delete(k)
		}
		return true
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "Empty hubs have been deleted",
	})
}

