package main

import (
	. "basedir/internal/adapter/primary/http"
	. "basedir/internal/adapter/primary/http/order"
	. "basedir/internal/adapter/primary/http/product"
	. "basedir/internal/adapter/secondary/infrastructure/config"
	. "basedir/internal/adapter/secondary/infrastructure/database"
	. "basedir/internal/adapter/secondary/infrastructure/logger"
	. "basedir/internal/adapter/secondary/repository/order"
	. "basedir/internal/adapter/secondary/repository/product"
	. "basedir/internal/core/domain/order"
	. "basedir/internal/core/domain/product"
)

func main() {
	cfg, err := LoadConfig(".env")
	if err != nil {
		panic(err.Error())
	}
	db := NewPostgresDatabase(
		cfg.GetDBHost(),
		cfg.GetDBUser(),
		cfg.GetDBPassword(),
		cfg.GetDBName(),
		cfg.GetDBPort(),
		cfg.GetDBSSLMode(),
		cfg.GetDBTimezone(),
	)
	log := NewLogger(cfg.GetServerEnv())
	baseApiPrefix := cfg.GetServerBaseApiPrefix()

	orderDbRepo := NewOrderDbRepository(db)
	orderSvc := NewOrderService(orderDbRepo)
	orderHttp := NewOrderHttpHandler(orderSvc)

	productDbRepo := NewProductDbRepository(db)
	productSvc := NewProductService(productDbRepo)
	productHttp := NewProductHttpHandler(productSvc)

	routeGroup := NewRouteGroup(orderHttp, productHttp)

	httpServer := NewHttpServer(cfg, log, baseApiPrefix)
	httpServer.SetupRoute(routeGroup)
	httpServer.Start()
}
