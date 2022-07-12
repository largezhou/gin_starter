package api

import (
	"github.com/gin-gonic/gin"
	"github.com/largezhou/gin_starter/app/api/userdto"
	"github.com/largezhou/gin_starter/app/middleware"
	"github.com/largezhou/gin_starter/app/model/user"
)

func InitRouter(r *gin.Engine) {
	{
		g := r.Group("/api").Use(
			middleware.Recovery(apiRecover),
		)

		g.GET("/user-list", func(ctx *gin.Context) {
			userList, err := user.GetUserList(ctx)
			if err != nil {
				fail(ctx, err)
				return
			}

			ok(ctx, userList)
		})

		g.POST("/create-user", func(ctx *gin.Context) {
			var dto userdto.CreateDto
			if err := ctx.ShouldBindJSON(&dto); err != nil {
				fail(ctx, err)
				return
			}

			userModel, err := user.Create(ctx, &user.User{
				Username: dto.Username,
				Password: dto.Password,
				Email:    dto.Email,
			})
			if err != nil {
				fail(ctx, err)
				return
			}

			ok(ctx, userModel)
		})
	}
}
