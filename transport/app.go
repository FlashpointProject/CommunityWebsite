package transport

import (
	"github.com/FlashpointProject/CommunityWebsite/config"
	"github.com/FlashpointProject/CommunityWebsite/service"
	"github.com/FlashpointProject/CommunityWebsite/utils"
)

type App struct {
	Conf    *config.AppConfig
	Service *service.Service
	CC      utils.CookieCutter
	Fpfss   *Fpfss
}
