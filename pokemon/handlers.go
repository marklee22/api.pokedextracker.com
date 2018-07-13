package pokemon

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pokedextracker/api.pokedextracker.com/application"
	"github.com/pokedextracker/api.pokedextracker.com/dexnumbers"
)

func listHandler(c *gin.Context) {
	app := c.MustGet("app").(*application.App)

	var pokemon []*Pokemon
	var dexNumbers []*dexnumbers.GameFamilyDexNumber
	var evolutions []*Evolution

	app.DB.Order("national_order ASC").Find(&pokemon)
	app.DB.Find(&dexNumbers)
	app.DB.
		Joins("INNER JOIN pokemon AS evolved ON evolutions.evolved_pokemon_id = evolved.id").
		Joins("INNER JOIN pokemon AS evolving ON evolutions.evolving_pokemon_id = evolving.id").
		Order("CASE WHEN trigger = 'breed' THEN evolving.national_id ELSE evolved.national_id END, trigger DESC, evolved.national_order ASC").
		Find(&evolutions)

	for _, p := range pokemon {
		p.LoadDexNumbers(dexNumbers)
		p.LoadEvolutions(evolutions)
	}

	c.JSON(http.StatusOK, pokemon)
}

func retrieveHandler(c *gin.Context) {
	app := c.MustGet("app").(*application.App)

	id := c.Param("id")

	var p Pokemon
	var dexNumbers []*dexnumbers.GameFamilyDexNumber
	var evolutions []*Evolution

	app.DB.Where("id = ?", id).First(&p)
	app.DB.Where("pokemon_id = ?", id).Find(&dexNumbers)
	app.DB.
		Joins("INNER JOIN pokemon AS evolved ON evolutions.evolved_pokemon_id = evolved.id").
		Joins("INNER JOIN pokemon AS evolving ON evolutions.evolving_pokemon_id = evolving.id").
		Where("evolutions.evolution_family_id = ?", p.EvolutionFamilyID).
		Order("CASE WHEN trigger = 'breed' THEN evolving.national_id ELSE evolved.national_id END, trigger DESC, evolved.national_order ASC").
		Find(&evolutions)

	p.LoadDexNumbers(dexNumbers)
	p.LoadEvolutions(evolutions)

	c.JSON(http.StatusOK, p)
}