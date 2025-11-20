package main

import (
	"log"
	"reviewer-api/internal/app/config"
	http "reviewer-api/internal/app/http-server/handlers"
	"reviewer-api/internal/app/repository/postgres"
	pkg "reviewer-api/internal/pkg/app"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pg, err := postgres.NewPostgers(cfg.GetDSN())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	// baseSlice := []string{"1", "2", "3"}
	// randomSlice := utils.GetRandomSlice(baseSlice)
	// log.Println(randomSlice)

	handl, _ := http.NewHandler(pg)
	router := gin.Default()
	app := pkg.NewApplication(cfg, router, handl)
	app.RunApplication()

}

// // TEST
// t, _ := pg.AddTeam(dto.TeamDTO{
// 	Name: "Team",
// 	Members: []dto.UserDTO{
// 		{ID: "user-68a77abf-ddc1-4398-a766-d7408ed2add9", Name: "name-1", IsActive: true},
// 		{ID: "user-13362a66-f1cb-4668-87f8-03a89f7bebb5", Name: "name-145", IsActive: false},
// 		{ID: "user-68a77abf", Name: "name-21", IsActive: true},
// 	},
// })

// js, _ := json.Marshal(t)
// log.Println(string(js))

// t, _ = pg.GetTeam("Team")
// js, _ = json.Marshal(t)
// log.Println(string(js))
// // TEST
