package engine

import (
	"database/sql"
	"log"

	db "github.com/bogunenko/template/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	store  *db.Store
	router *gin.Engine
}

func NewEngine(dbSource string) (*Engine, error) {

	conn, err := sql.Open("mysql", dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	engine := &Engine{
		store: db.NewStore(conn),
	}

	engine.setupRouter()
	return engine, nil
}

func (engine *Engine) setupRouter() {
	router := gin.Default()

	router.POST("/account", engine.createAccount)
	router.GET("/account/:id", engine.getAccount)
	router.POST("/deposit", engine.deposit)
	router.POST("/withdraw", engine.withdraw)
	router.POST("/transfer", engine.transfer)

	engine.router = router
}

func (engine *Engine) Start() error {
	return engine.router.Run(":8080")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
