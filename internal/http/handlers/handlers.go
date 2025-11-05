//go:generate openapi-generator generate -i ../../../schema.yaml -g go-gin-server -o ../../../gen -p apiPath=openapi,interfaceOnly=true,packageName=openapi,hideGenerationTimestamp=true

package httptransport

type API struct {
	auth *service.AuthService
	shop *service.ShopService
	jwt  *jwtutil.HS
}
