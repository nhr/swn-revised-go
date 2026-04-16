package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"

	"swn-revised-go/internal/generator"
	"swn-revised-go/internal/render"
	"swn-revised-go/internal/util"
	"swn-revised-go/internal/wiki"
)

// Handlers holds shared resources for HTTP handlers.
type Handlers struct {
	DB       *sql.DB
	LangData map[string]string
	SwnHTML  string
	BgPNG    []byte
	DotPNG   []byte
	FontData []byte
}

// SeedHandler returns a random seed token.
func (h *Handlers) SeedHandler(w http.ResponseWriter, r *http.Request) {
	token := util.TokenizeSeed(util.RandomSeed())
	resp := map[string]string{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// generateSector seeds the RNG and generates all sector data for the given token.
func (h *Handlers) generateSector(token string) *generator.Sector {
	seed := util.UntokenizeSeed(token)
	rand.Seed(seed)

	gen := generator.New(h.DB, h.LangData)

	sectorName, starMap, worlds := gen.GenSector()
	npcs := gen.GenNPCs(0)
	corps := gen.GenCorps(0)
	rels := gen.GenReligions(0)
	pols := gen.GenPolParties(0)
	aliens := gen.GenAliens(0)
	pois := gen.GenPointsOfInterest(starMap)
	patrons := gen.GenPatrons(0)
	beasts := gen.GenBeasts(0)
	factions := gen.GenFactions(0)
	urbanEnc := gen.GenUrbanEncounters(0)
	wildEnc := gen.GenWildEncounters(0)

	return &generator.Sector{
		Name:     sectorName,
		Token:    token,
		Stars:    starMap,
		Worlds:   worlds,
		NPCs:     npcs,
		Corps:    corps,
		Rels:     rels,
		Pols:     pols,
		Aliens:   aliens,
		POIs:     pois,
		Patrons:  patrons,
		Beasts:   beasts,
		Factions: factions,
		UrbanEnc: urbanEnc,
		WildEnc:  wildEnc,
	}
}

// SectorHandler returns structured sector data as JSON.
func (h *Handlers) SectorHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	sector := h.generateSector(token)

	sort.Slice(sector.NPCs, func(i, j int) bool { return sector.NPCs[i].Sort < sector.NPCs[j].Sort })
	sort.Slice(sector.Corps, func(i, j int) bool { return sector.Corps[i].Name < sector.Corps[j].Name })
	sort.Slice(sector.Rels, func(i, j int) bool { return sector.Rels[i].Name < sector.Rels[j].Name })
	sort.Slice(sector.Pols, func(i, j int) bool { return sector.Pols[i].Name < sector.Pols[j].Name })
	sort.Slice(sector.Worlds, func(i, j int) bool { return sector.Worlds[i].Name < sector.Worlds[j].Name })
	sort.Slice(sector.Aliens, func(i, j int) bool { return sector.Aliens[i].Name < sector.Aliens[j].Name })
	sort.Slice(sector.Patrons, func(i, j int) bool { return sector.Patrons[i].Name < sector.Patrons[j].Name })
	sort.Slice(sector.Beasts, func(i, j int) bool { return sector.Beasts[i].Name < sector.Beasts[j].Name })
	sort.Slice(sector.Factions, func(i, j int) bool { return sector.Factions[i].Name < sector.Factions[j].Name })

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sector)
}

// MapHandler renders a sector map as a PNG image.
func (h *Handlers) MapHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	seed := util.UntokenizeSeed(token)
	rand.Seed(seed)

	gen := generator.New(h.DB, h.LangData)
	sectorName, starMap, _ := gen.GenSector()

	mapData, _ := render.RenderMap(sectorName, starMap, h.BgPNG, h.DotPNG, h.FontData, true)

	w.Header().Set("Content-Type", "image/png")
	w.Write(mapData)
}

// ExportHandler generates a TiddlyWiki ZIP export of the sector.
func (h *Handlers) ExportHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	sector := h.generateSector(token)

	// Apply custom star positions if provided in the request body
	if r.Body != nil && r.ContentLength > 0 {
		var body struct {
			Stars []*generator.Star `json:"stars"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil && len(body.Stars) > 0 {
			sector.Stars = body.Stars

			// Build StarID-to-Cell map from updated positions
			starCells := make(map[int]string)
			for _, s := range sector.Stars {
				starCells[s.ID] = s.Cell
			}

			// Update POI cells to match moved stars
			for _, poi := range sector.POIs {
				if newCell, ok := starCells[poi.StarID]; ok {
					poi.Cell = newCell
				}
			}

			// Update World cells to match moved stars
			for _, w := range sector.Worlds {
				if newCell, ok := starCells[w.StarID]; ok {
					w.Cell = newCell
				}
			}
		}
	}

	zipData := wiki.GenerateWiki(sector, h.SwnHTML, sector.Stars, sector.Worlds,
		h.BgPNG, h.DotPNG, h.FontData, false)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment;filename=SWN_Sector_%s.zip", token))
	w.Write(zipData)
}

