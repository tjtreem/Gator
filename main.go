package main

import (
	"fmt"
	"os"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/tjtreem/gator/internal/config"
	"github.com/tjtreem/gator/internal/database"

)


func main() {
    
    cfg, err := config.Read()
    if err != nil {
	fmt.Println("Error:", err)
	return
    }

    db, err := sql.Open("postgres", cfg.DBUrl)
    if err != nil{
	log.Fatal(err)
    }
    
    dbQueries := database.New(db)

    state := State{
	Db:  dbQueries,
	Cfg: &cfg,	
}

    cmds := Commands{
	Handlers: make(map[string]func(*State, Command) error),
    }

    cmds.Register("login", HandlerLogin)
    cmds.Register("register", HandlerRegister)
    cmds.Register("reset", HandlerReset)
    cmds.Register("users", HandlerGetUsers)
    cmds.Register("agg", HandlerAgg)


    if len(os.Args) < 2 {
	fmt.Println("Not enough arguments provided")
	os.Exit(1)
    }

    commandName := os.Args[1]
    args := os.Args[2:]

    cmd := Command{
        Name:	commandName,
        Args:	args,
    }
    
    err = cmds.Run(&state, cmd)
    if err != nil {
	fmt.Println("Error:", err)
	os.Exit(1)
    }

}
