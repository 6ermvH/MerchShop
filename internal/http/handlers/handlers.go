//go:generate openapi-generator generate -i ../../../schema.yaml -g go-gin-server -o ../../../gen -p apiPath=openapi,interfaceOnly=true,packageName=openapi,hideGenerationTimestamp=true
//go:generate mockgen -source=../../../gen/openapi/api_default.go -destination=../../../gen/mock/handlers/mock_handlers.go -package=mock_handlers

package handlers

import (
	"github.com/6ermvH/MerchShop/internal/http/middleware"
	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
)

type API struct {
	repos repo.MerchRepo
	hs    jwtutil.JWT
}

func NewAPI(repo repo.MerchRepo, hs jwtutil.JWT) *API {
	return &API{
		repos: repo,
		hs:    hs,
	}
}

func (api *API) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/auth", api.ApiAuthPost)

	apiG := r.Group("/api", middleware.Auth(api.hs, api.repos))
	{
		apiG.GET("/buy/:item", api.ApiBuyItemGet)
		apiG.GET("/info", api.ApiInfoGet)
		apiG.POST("/sendCoin", api.ApiSendCoinPost)
	}
}
