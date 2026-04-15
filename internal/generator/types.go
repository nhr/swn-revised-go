package generator

// Star represents a star system on the sector map.
type Star struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Cell  string `json:"cell"`
	ID    int    `json:"id"`
}

// WorldComponent holds a rolled world attribute (atmosphere, biosphere, etc.)
type WorldComponent struct {
	Roll  int    `json:"roll"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Short string `json:"short"`
}

// WorldTag holds a rolled world tag.
type WorldTag struct {
	Roll          string `json:"roll"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Enemies       string `json:"enemies"`
	Friends       string `json:"friends"`
	Complications string `json:"complications"`
	Things        string `json:"things"`
	Places        string `json:"places"`
	Short         string `json:"short"`
}

// WorldRelationship describes how two worlds in the same system relate.
type WorldRelationship struct {
	Origin       string `json:"origin"`
	Relationship string `json:"relationship"`
	ContactPoint string `json:"contact_point"`
}

// World represents a generated planet.
type World struct {
	Name        string         `json:"name"`
	Atmosphere  WorldComponent `json:"atmosphere"`
	Biosphere   WorldComponent `json:"biosphere"`
	Population  WorldComponent `json:"population"`
	TechLevel   WorldComponent `json:"tech_level"`
	Temperature WorldComponent `json:"temperature"`
	Tags        []WorldTag     `json:"tags"`
	SysCount    int            `json:"sys_count"`
	SysPos      int            `json:"sys_pos"`
	SysName     string         `json:"sys_name"`
	SysNum      string         `json:"sys_num"`
	StarID      int            `json:"star_id"`
	Cell        string         `json:"cell"`
	CellPos     string         `json:"cell_pos"`
	PNum         int                `json:"p_num,omitempty"`
	Relationship *WorldRelationship `json:"relationship,omitempty"`
}

// NPC represents a generated non-player character.
type NPC struct {
	Name       string `json:"name"`
	Sort       string `json:"sort"`
	Age        string `json:"age"`
	Background string `json:"background"`
	Role       string `json:"role"`
	Problem    string `json:"problem"`
	Desire     string `json:"desire"`
	Trait      string `json:"trait"`
}

// Corp represents a generated corporation.
type Corp struct {
	Name       string `json:"name"`
	Business   string `json:"business"`
	Reputation string `json:"reputation"`
}

// Heresy represents a heretical offshoot of a religion.
type Heresy struct {
	Name     string `json:"name"`
	Attitude string `json:"attitude"`
	Founder  string `json:"founder"`
	HerDesc  string `json:"heresy"`
	Quirk    string `json:"quirk"`
}

// Religion represents a generated religion.
type Religion struct {
	Name       string   `json:"name"`
	Origin     string   `json:"origin"`
	Leadership string   `json:"leadership"`
	Evolution  string   `json:"evolution"`
	Offshoots  []Heresy `json:"offshoots"`
}

// PoliticalIssue represents a political issue.
type PoliticalIssue struct {
	Issue string `json:"issue"`
	Tag   string `json:"tag"`
}

// PoliticalParty represents a generated political group.
type PoliticalParty struct {
	Name         string           `json:"name"`
	Issues       []PoliticalIssue `json:"issues"`
	Leadership   string           `json:"leadership"`
	Relationship string           `json:"relationship"`
	Policy       string           `json:"policy"`
}

// Architecture represents generated architectural features.
type Architecture struct {
	Name        string `json:"name"`
	Towers      string `json:"towers"`
	Foundations string `json:"foundations"`
	WallDecor   string `json:"wall_decor"`
	Supports    string `json:"supports"`
	Arches      string `json:"arches"`
	Extras      string `json:"extras"`
}

// Alien represents a generated alien race.
type Alien struct {
	Name   string       `json:"name"`
	Body   [][]string   `json:"body"`
	Lens   [][]string   `json:"lens"`
	Vars   []string     `json:"vars"`
	Social [][]string   `json:"social"`
	Arch   Architecture `json:"arch"`
}

// PointOfInterest represents a generated system feature.
type PointOfInterest struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	OccupiedBy string `json:"occupied_by"`
	Situation  string `json:"situation"`
	StarID     int    `json:"star_id"`
	Cell       string `json:"cell"`
}

// Patron represents a generated quest-giving NPC.
type Patron struct {
	Name            string `json:"name"`
	Trustworthiness string `json:"trustworthiness"`
	Eagerness       string `json:"eagerness"`
	Reward          string `json:"reward"`
	Challenge       string `json:"challenge"`
	Opposition      string `json:"opposition"`
	Complication    string `json:"complication"`
}

// Beast represents a generated creature.
type Beast struct {
	Name      string   `json:"name"`
	Features  []string `json:"features"`
	BodyPlan  string   `json:"body_plan"`
	LimbMod   string   `json:"limb_mod"`
	SkinMod   string   `json:"skin_mod"`
	Weapon    string   `json:"weapon"`
	Size      string   `json:"size"`
	BehavType string   `json:"behavior_type"`
	Behavior  string   `json:"behavior"`
	Discharge string   `json:"discharge,omitempty"`
	Poison    string   `json:"poison,omitempty"`
}

// FactionTag describes a faction characteristic.
type FactionTag struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// Faction represents a generated political faction.
type Faction struct {
	Name      string       `json:"name"`
	Homeworld string       `json:"homeworld"`
	Tags      []FactionTag `json:"tags"`
	Force     int          `json:"force"`
	Cunning   int          `json:"cunning"`
	Wealth    int          `json:"wealth"`
	HP        int          `json:"hp"`
	MaxHP     int          `json:"max_hp"`
	FacCreds  int          `json:"faccreds"`
	Goal      string       `json:"goal"`
}

// UrbanEncounter represents a generated urban encounter.
type UrbanEncounter struct {
	Conflict   string `json:"conflict"`
	Antagonist string `json:"antagonist"`
	Venue      string `json:"venue"`
	PCInvolve  string `json:"pc_involve"`
	Event      string `json:"event"`
	Feature    string `json:"feature"`
}

// WildEncounter represents a generated wilderness encounter.
type WildEncounter struct {
	Range    string `json:"range"`
	Hostile  string `json:"hostile"`
	Weather  string `json:"weather"`
	Nature   string `json:"nature"`
	Friendly string `json:"friendly"`
	Feature  string `json:"feature"`
}

// Sector holds all generated data for a sector.
type Sector struct {
	Name    string            `json:"name"`
	Token   string            `json:"token"`
	Map     interface{}       `json:"map"`
	Stars   []*Star           `json:"stars,omitempty"`
	Worlds  []*World          `json:"worlds,omitempty"`
	NPCs    []*NPC            `json:"npcs,omitempty"`
	Corps   []*Corp           `json:"corps,omitempty"`
	Rels    []*Religion       `json:"rels,omitempty"`
	Pols    []*PoliticalParty `json:"pols,omitempty"`
	Aliens  []*Alien          `json:"aliens,omitempty"`
	POIs    []*PointOfInterest `json:"pois,omitempty"`
	Patrons []*Patron         `json:"patrons,omitempty"`
	Beasts  []*Beast          `json:"beasts,omitempty"`
	Factions []*Faction       `json:"factions,omitempty"`
	UrbanEnc []*UrbanEncounter `json:"urban_encounters,omitempty"`
	WildEnc  []*WildEncounter  `json:"wild_encounters,omitempty"`
}
