package servers

import (
	"github.com/MarkTBSS/071_Admin/modules/middlewares/middlewaresHandlers"
	"github.com/MarkTBSS/071_Admin/modules/middlewares/middlewaresRepositories"
	"github.com/MarkTBSS/071_Admin/modules/middlewares/middlewaresUsecases"
	_pkgModulesMonitorMonitorHandlers "github.com/MarkTBSS/071_Admin/modules/monitor/monitorHandlers"
	"github.com/MarkTBSS/071_Admin/modules/users/usersHandlers"
	"github.com/MarkTBSS/071_Admin/modules/users/usersRepositories"
	"github.com/MarkTBSS/071_Admin/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	handler := middlewaresHandlers.MiddlewaresHandler(usecase, s.cfg)
	return middlewaresHandlers.MiddlewaresHandler(handler, s.cfg)
}

func (m *moduleFactory) MonitorModule() {
	handler := _pkgModulesMonitorMonitorHandlers.MonitorHandler(m.s.cfg)
	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)
	router.Get("/signup-admin", handler.SignUpAdmin)
	router.Get("/admin/secret", handler.GenerateAdminToken)
}
