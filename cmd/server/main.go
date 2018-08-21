package main

import (
	"github.com/dmitrymomot/go-test-task-management-tool/app/infrastructure"
	"github.com/dmitrymomot/go-test-task-management-tool/app/interfaces/repositories"
	"github.com/dmitrymomot/go-test-task-management-tool/app/interfaces/webservices"
	"github.com/dmitrymomot/go-test-task-management-tool/app/usecases"
)

func main() {
	// loading of configs from file
	config := loadConfig()

	// init logger
	logger := infrastructure.NewLogger(false)

	// init db connection
	db := infrastructure.NewMySQLHandler(config.DBSource)
	defer db.Close()

	// run migrations
	if err := repositories.Migrate(db); err != nil {
		panic(err)
	}

	// creation of application
	repo := repositories.NewTaskRepository(db)
	interactor := usecases.NewTaskInteractor(repo)
	ws := webservices.NewTasks(interactor, logger)

	// run server
	router := newRouter(ws)
	server := newServer(router.setup(), logger)
	logger.Error(server.run(config.ListenAddress))
}
