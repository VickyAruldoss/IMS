package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vickyaruldoss/ims/controller"
	"github.com/vickyaruldoss/ims/repository"
	"github.com/vickyaruldoss/ims/service"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	memberRepo := repository.NewPostgresRepository(db)
	memberSvc := service.NewMemberService(memberRepo)
	memberCtrl := controller.NewMemberController(memberSvc)

	v1 := r.Group("/api/v1")
	{
		members := v1.Group("/members")
		{
			members.POST("", memberCtrl.CreateMember)
			members.GET("", memberCtrl.GetAllMembers)
			members.GET("/:id", memberCtrl.GetMember)
			members.PUT("/:id", memberCtrl.UpdateMember)
			members.DELETE("/:id", memberCtrl.DeleteMember)
		}
	}

	return r
}
