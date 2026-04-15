// Command migrate updates the embedded SQLite database for the
// Stars Without Number Revised Edition.
//
// Usage:
//
//	go run ./cmd/migrate [path/to/swn.sqlite]
//
// If no path is given it defaults to static/swn.sqlite.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "static/swn.sqlite"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Enable WAL mode for better write performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		log.Fatalf("set journal mode: %v", err)
	}

	steps := []struct {
		name string
		fn   func(*sql.DB) error
	}{
		{"World Population", migrateWorldPopulation},
		{"World Temperature", migrateWorldTemperature},
		{"World Tech Level", migrateWorldTechLevel},
		{"NPC tables", migrateNPCTables},
		{"Alien Body Types", migrateAlienBodyTypes},
		{"Alien Lenses", migrateAlienLenses},
		{"Alien Structure", migrateAlienStructure},
		{"World Tags", migrateWorldTags},
		{"Patron tables", migratePatronTables},
		{"Urban Encounter tables", migrateUrbanEncounterTables},
		{"Wilderness Encounter tables", migrateWildEncounterTables},
		{"POI tables", migratePOITables},
		{"Beast tables", migrateBeastTables},
		{"Faction tables", migrateFactionTables},
		{"World Relationship table", migrateWorldRelationship},
	}

	for _, s := range steps {
		fmt.Printf("  migrating %s...\n", s.name)
		if err := s.fn(db); err != nil {
			log.Fatalf("migrate %s: %v", s.name, err)
		}
	}

	// Compact the database
	if _, err := db.Exec("VACUUM"); err != nil {
		log.Fatalf("vacuum: %v", err)
	}

	fmt.Println("migration complete")
}

// execAll runs every statement in the slice inside a transaction.
func execAll(db *sql.DB, stmts []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, s := range stmts {
		if _, err := tx.Exec(s); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec %q: %w", s[:min(len(s), 80)], err)
		}
	}
	return tx.Commit()
}

// insertRows drops/creates a table and bulk-inserts rows.
func insertRows(db *sql.DB, drop, create string, insert string, rows [][]interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(drop); err != nil {
		tx.Rollback()
		return fmt.Errorf("drop: %w", err)
	}
	if _, err := tx.Exec(create); err != nil {
		tx.Rollback()
		return fmt.Errorf("create: %w", err)
	}
	stmt, err := tx.Prepare(insert)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()
	for _, r := range rows {
		if _, err := stmt.Exec(r...); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert row %v: %w", r[0], err)
		}
	}
	return tx.Commit()
}

// ---------------------------------------------------------------------------
// World Population (IDs 4-12 changed)
// ---------------------------------------------------------------------------

func migrateWorldPopulation(db *sql.DB) error {
	updates := []struct {
		id    int
		name  string
		desc  string
		short string
	}{
		{4, "Fewer than a million inhabitants",
			"Fewer than a million inhabitants are often found on newly-colonized worlds, or those with limited arable land. This class has the widest range of results, from societies barely large enough to avoid inbreeding to minor worlds of several hundred thousand locals struggling against hostile local conditions.\nFor worlds with a very small population, settlements tend to be clustered close to one another for mutual support and efficient exploitation of land. The atmosphere might be something like that of a village-dotted countryside, with too few people to support large, specialized cities. The major settlement of the planet is usually built around the spaceport.\nPolitics on planets with so few people tend to be very personal in nature. Individual leaders can address the entire population of a world at a single meeting, and networks of family, friends, and employees can dominate the local social life. This small population combined with limited outside contact can result in somewhat... unusual societies more strongly influenced by a leader's personal quirks than a more prosperous planet might be. Some such worlds can grow decidedly unhealthy in their customs and traditions.",
			"<1M"},
		{5, "Fewer than a million inhabitants",
			"Fewer than a million inhabitants are often found on newly-colonized worlds, or those with limited arable land. This class has the widest range of results, from societies barely large enough to avoid inbreeding to minor worlds of several hundred thousand locals struggling against hostile local conditions.\nFor worlds with a very small population, settlements tend to be clustered close to one another for mutual support and efficient exploitation of land. The atmosphere might be something like that of a village-dotted countryside, with too few people to support large, specialized cities. The major settlement of the planet is usually built around the spaceport.\nPolitics on planets with so few people tend to be very personal in nature. Individual leaders can address the entire population of a world at a single meeting, and networks of family, friends, and employees can dominate the local social life. This small population combined with limited outside contact can result in somewhat... unusual societies more strongly influenced by a leader's personal quirks than a more prosperous planet might be. Some such worlds can grow decidedly unhealthy in their customs and traditions.",
			"<1M"},
		{6, "Several million inhabitants",
			"Several million inhabitants make up the population of most frontier worlds, usually ones with adequate supplies of arable land and surface conditions amenable to a human without a vacc suit or filter mask. This level of population is also usually the most that a world at TL4 can support using only artificial hydroponic complexes and sealed agricultural systems, so less hospitable planets rarely have more than this many inhabitants.\nAt this level of population, city-states and hub settlements start to form, and some political divisions might exist beneath the overall planetary government. Regional variations start to become present as towns and smaller settlements start to assert their own cultural traits. Young colonies rarely have much of a framework for resolving these new tensions, and it's not uncommon for outsiders to end up as catspaws for varying factions.\nThis level of population is normally the smallest that can create spike drive-equipped starships. Less populous worlds rarely have the broad range of specialists and workers necessary to build such ships, even if they possess the necessary technical information.",
			"Millions"},
		{7, "Several million inhabitants",
			"Several million inhabitants make up the population of most frontier worlds, usually ones with adequate supplies of arable land and surface conditions amenable to a human without a vacc suit or filter mask. This level of population is also usually the most that a world at TL4 can support using only artificial hydroponic complexes and sealed agricultural systems, so less hospitable planets rarely have more than this many inhabitants.\nAt this level of population, city-states and hub settlements start to form, and some political divisions might exist beneath the overall planetary government. Regional variations start to become present as towns and smaller settlements start to assert their own cultural traits. Young colonies rarely have much of a framework for resolving these new tensions, and it's not uncommon for outsiders to end up as catspaws for varying factions.\nThis level of population is normally the smallest that can create spike drive-equipped starships. Less populous worlds rarely have the broad range of specialists and workers necessary to build such ships, even if they possess the necessary technical information.",
			"Millions"},
		{8, "Several million inhabitants",
			"Several million inhabitants make up the population of most frontier worlds, usually ones with adequate supplies of arable land and surface conditions amenable to a human without a vacc suit or filter mask. This level of population is also usually the most that a world at TL4 can support using only artificial hydroponic complexes and sealed agricultural systems, so less hospitable planets rarely have more than this many inhabitants.\nAt this level of population, city-states and hub settlements start to form, and some political divisions might exist beneath the overall planetary government. Regional variations start to become present as towns and smaller settlements start to assert their own cultural traits. Young colonies rarely have much of a framework for resolving these new tensions, and it's not uncommon for outsiders to end up as catspaws for varying factions.\nThis level of population is normally the smallest that can create spike drive-equipped starships. Less populous worlds rarely have the broad range of specialists and workers necessary to build such ships, even if they possess the necessary technical information.",
			"Millions"},
		{9, "Hundreds of millions of inhabitants",
			"Hundreds of millions of inhabitants are found on some of the more populous frontier worlds, those planets with an excellent climate and soil for agriculture. These worlds teem with human life compared to most of their neighbors, and if they have sufficient technological expertise to build spike drives, they probably are one of the more powerful planets in the sector. The sheer amount of labor and expertise they can put into their projects dwarfs the capabilities of most other worlds.\nSome worlds of this type fail to coalesce under a unified planetary government, however, and turn most of their attention to intra-planetary quarrels between different nations or city-states. Full-scale planetary wars have been known to break out on such worlds, with some being reduced to burnt-out cinders by some maltech weapon or by nuclear attacks against regions stripped of their nuke snuffers.",
			"100M+"},
		{10, "Billions of inhabitants",
			"Billions of inhabitants on a world is a freak occurrence in most sectors, a circumstance dependent on almost perfect planetary conditions, a First Wave date of colonization, and a very large initial colonial expedition. These worlds have been colonized for almost a thousand years, and have likely gone through substantial cultural development and change since the original founding.\nPlanets with this many inhabitants tend to come in two main varieties. The more common is that of regional hegemon, leveraging its enormous reserves of labor and scientific personnel to build hundreds or thousands of starships. Some hegemons might even have the transport weight and personnel available to colonize or conquer their less powerful neighbors outright, though most hegemons prefer the safer route of docile client states rather than outright conquests.\nNot all such worlds have the necessary resources to build starships, even when the labor is cheap and abundant. These worlds tend to take a different route, becoming balkanized, conflict-torn planets not unlike a mid-20th-century Earth.",
			"1B+"},
		{11, "Billions of inhabitants",
			"Billions of inhabitants on a world is a freak occurrence in most sectors, a circumstance dependent on almost perfect planetary conditions, a First Wave date of colonization, and a very large initial colonial expedition. These worlds have been colonized for almost a thousand years, and have likely gone through substantial cultural development and change since the original founding.\nPlanets with this many inhabitants tend to come in two main varieties. The more common is that of regional hegemon, leveraging its enormous reserves of labor and scientific personnel to build hundreds or thousands of starships. Some hegemons might even have the transport weight and personnel available to colonize or conquer their less powerful neighbors outright, though most hegemons prefer the safer route of docile client states rather than outright conquests.\nNot all such worlds have the necessary resources to build starships, even when the labor is cheap and abundant. These worlds tend to take a different route, becoming balkanized, conflict-torn planets not unlike a mid-20th-century Earth.",
			"1B+"},
		{12, "Alien inhabitants",
			"Alien civilizations are even rarer than massively human-populated planets, but they can be found in some sectors. On these worlds, any human presence is either nonexistent or simply a token contingent of traders and diplomats. The actual number of aliens may range from a few tens of thousands to teeming worlds of billions.\nAlien civilizations appear to be no more immune to conflict, disaster, and decadence than that of humans, and most alien worlds are isolated planets. Some sectors were known to be dominated by alien empires before the Scream, but most inhabited worlds found by colonial surveyors consisted of indigenous, low-tech sapients who evolved on the world, or else were the decaying remnants of a former empire that had fallen apart into isolated colonial worlds.",
			"Alien"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, u := range updates {
		if _, err := tx.Exec(
			"UPDATE WORLD_POPULATION SET name = ?, desc = ?, short = ? WHERE id = ?",
			u.name, u.desc, u.short, u.id,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ---------------------------------------------------------------------------
// World Temperature (IDs 3-5, 9-11 changed)
// ---------------------------------------------------------------------------

func migrateWorldTemperature(db *sql.DB) error {
	updates := []struct {
		id    int
		name  string
		desc  string
		short string
	}{
		{3, "Cold",
			"Cold worlds are uncomfortable, but a human can survive on them in nothing more than heavy clothing. The worst of the cold worlds are similar in condition to Earth's Antarctic regions, barren wastes of ice and wind. The more clement ones have brief warmer seasons or equatorial bands that get enough solar radiation to support substantial agriculture.\nOn cold worlds, agricultural land is at a premium. The population is unlikely to be willing to limit its numbers to what hydroponic farms can provide, and they may lack the technology for wide-scale artificial production methods. Savage wars may have broken out over control of fertile growing regions, and populations have perhaps been driven away from the arable land into the cold zones and a lingering death.",
			"Cold"},
		{4, "Variable cold with temperate places",
			"Variable temperature worlds tend to show a greater distribution of climates than other worlds. This may be a climate that changes world-wide when a long, slow orbit brings the planet into proximity with a stellar primary, or it might be a world that has substantially different climatic zones spread across its surface. The north pole of one world might be a sun-blasted desert that moderates to a cool, wet equator, or an icy world might be warmed here and there by complex channels of geothermally-heated subsurface rivers.\nVariable temperature worlds tend to have savage weather. The mixing of hot and cold air can send ferocious cyclones and raging hurricanes across the surface of the world, some large enough to consume most of a hemisphere.",
			"Var. Cold"},
		{5, "Variable cold with temperate places",
			"Variable temperature worlds tend to show a greater distribution of climates than other worlds. This may be a climate that changes world-wide when a long, slow orbit brings the planet into proximity with a stellar primary, or it might be a world that has substantially different climatic zones spread across its surface. The north pole of one world might be a sun-blasted desert that moderates to a cool, wet equator, or an icy world might be warmed here and there by complex channels of geothermally-heated subsurface rivers.\nVariable temperature worlds tend to have savage weather. The mixing of hot and cold air can send ferocious cyclones and raging hurricanes across the surface of the world, some large enough to consume most of a hemisphere.",
			"Var. Cold"},
		{9, "Variable warm with temperate places",
			"Variable temperature worlds tend to show a greater distribution of climates than other worlds. This may be a climate that changes world-wide when a long, slow orbit brings the planet into proximity with a stellar primary, or it might be a world that has substantially different climatic zones spread across its surface. The north pole of one world might be a sun-blasted desert that moderates to a cool, wet equator, or an icy world might be warmed here and there by complex channels of geothermally-heated subsurface rivers.\nVariable temperature worlds tend to have savage weather. The mixing of hot and cold air can send ferocious cyclones and raging hurricanes across the surface of the world, some large enough to consume most of a hemisphere.",
			"Var. Warm"},
		{10, "Warm",
			"Warm worlds come in two main flavors, depending on the prevalence of water. Desert worlds are hot enough or arid enough to be deprived of most surface water. Any life on such a planet has to be capable of extracting water from living prey, air currents or moisture supplies deep underground. Humans can survive on warm worlds without more than appropriate clothing, but desert-world humans must learn sophisticated techniques of water prospecting and well-drilling in order to maintain their agricultural systems. Warfare often revolves around these water systems.\nThe other common variety of warm world is one blessed with an abundance of surface water. These worlds tend to be covered with plant and animal life, teeming with thick jungles and vast mats of sargasso on seas that boil with piscids and other alien life analogs. These worlds can be quite prosperous for humans if the native life is edible, but inhabitants are often forced to deal with large predators and a vividly ingenious native ecology of disease.",
			"Warm"},
		{11, "Warm",
			"Warm worlds come in two main flavors, depending on the prevalence of water. Desert worlds are hot enough or arid enough to be deprived of most surface water. Any life on such a planet has to be capable of extracting water from living prey, air currents or moisture supplies deep underground. Humans can survive on warm worlds without more than appropriate clothing, but desert-world humans must learn sophisticated techniques of water prospecting and well-drilling in order to maintain their agricultural systems. Warfare often revolves around these water systems.\nThe other common variety of warm world is one blessed with an abundance of surface water. These worlds tend to be covered with plant and animal life, teeming with thick jungles and vast mats of sargasso on seas that boil with piscids and other alien life analogs. These worlds can be quite prosperous for humans if the native life is edible, but inhabitants are often forced to deal with large predators and a vividly ingenious native ecology of disease.",
			"Warm"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, u := range updates {
		if _, err := tx.Exec(
			"UPDATE WORLD_TEMPERATURE SET name = ?, desc = ?, short = ? WHERE id = ?",
			u.name, u.desc, u.short, u.id,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ---------------------------------------------------------------------------
// World Tech Level (IDs 2, 4-11 changed)
// ---------------------------------------------------------------------------

func migrateWorldTechLevel(db *sql.DB) error {
	updates := []struct {
		id    int
		name  string
		desc  string
		short string
	}{
		{2, "TL0, neolithic-level technology",
			"Tech level 0 represents a world with technologies similar to neolithic-era humanity. Very few worlds collapse this far without dying out completely, and so a world this primitive is very rare in human space. The world is so profoundly devoid of useful resources that the natives may simply not have anything better to work with than rocks and native vegetation.\nNatives of a world at this tech level might retain elaborate cultural artifacts and a very sophisticated society. Existing resources tend to be exploited to the hilt, with large populations producing massive stone cities, human-powered engines for irrigation, and vast displays of foodstuffs and craft work for the ruling class.",
			"0"},
		{4, "TL2, early Industrial Age tech",
			"Tech level 2 worlds are blessed with the presence of fossil fuels or a similarly energy-dense substance can be found. Internal combustion engines are crude things compared to the power of a fusion plant, but they can be built with primitive tools and a basic degree of education. The mechanization that results from this resource vastly increases the efficiency of local farming, the speed of travel, and the productivity of factory workers.\nMost tech level 2 worlds are knockdowns, planets that once had a much higher degree of sophistication but were recently crippled by disease, disaster, or invasion. The natives tend to be keenly aware of the glories of former generations, and work fiercely to redevelop lost techniques and reclaim old knowledge.",
			"2"},
		{5, "TL2, early Industrial Age tech",
			"Tech level 2 worlds are blessed with the presence of fossil fuels or a similarly energy-dense substance can be found. Internal combustion engines are crude things compared to the power of a fusion plant, but they can be built with primitive tools and a basic degree of education. The mechanization that results from this resource vastly increases the efficiency of local farming, the speed of travel, and the productivity of factory workers.\nMost tech level 2 worlds are knockdowns, planets that once had a much higher degree of sophistication but were recently crippled by disease, disaster, or invasion. The natives tend to be keenly aware of the glories of former generations, and work fiercely to redevelop lost techniques and reclaim old knowledge.",
			"2"},
		{6, "TL4, modern postech",
			"Tech level 4 worlds are the most common in human space, and their technical expertise is the baseline for modern post-Silence postech. These worlds can create spike drives rated up to drive-3, fusion power plants, grav vehicles, simple energy weapons, and medicines that extend human life to a hundred years of vigorous good health. They can manage sophisticated gengineering on simple life forms, and some tech level 4 worlds have even attempted to improve the genetic structures of human life itself.\nMost worlds with regular interstellar contact and the necessary raw resources eventually gravitate towards this level of technological expertise. It may take decades, or even centuries in the case of more retrograde worlds, but this tech level is the highest that can readily be achieved by most planets.",
			"4"},
		{7, "TL4, modern postech",
			"Tech level 4 worlds are the most common in human space, and their technical expertise is the baseline for modern post-Silence postech. These worlds can create spike drives rated up to drive-3, fusion power plants, grav vehicles, simple energy weapons, and medicines that extend human life to a hundred years of vigorous good health. They can manage sophisticated gengineering on simple life forms, and some tech level 4 worlds have even attempted to improve the genetic structures of human life itself.\nMost worlds with regular interstellar contact and the necessary raw resources eventually gravitate towards this level of technological expertise. It may take decades, or even centuries in the case of more retrograde worlds, but this tech level is the highest that can readily be achieved by most planets.",
			"4"},
		{8, "TL4, modern postech",
			"Tech level 4 worlds are the most common in human space, and their technical expertise is the baseline for modern post-Silence postech. These worlds can create spike drives rated up to drive-3, fusion power plants, grav vehicles, simple energy weapons, and medicines that extend human life to a hundred years of vigorous good health. They can manage sophisticated gengineering on simple life forms, and some tech level 4 worlds have even attempted to improve the genetic structures of human life itself.\nMost worlds with regular interstellar contact and the necessary raw resources eventually gravitate towards this level of technological expertise. It may take decades, or even centuries in the case of more retrograde worlds, but this tech level is the highest that can readily be achieved by most planets.",
			"4"},
		{9, "TL3, tech like that of present-day Earth",
			"Tech level 3 worlds are those that have developed to a level much like that of the more advanced nations of late twentieth-century Earth. The main difference between tech level 3 worlds and level 4 worlds is the lack of fusion power or spike drive manufacture.\nThese worlds come in two major varieties. The first is a knockdown world that formerly had tech level 4 technologies before some catastrophe or pirate attack destroyed their knowledge and industrial base. The second major variety is a world that suffered a severe disaster shortly after colonization and has been forced to painstakingly build up their technological base ever since.",
			"3"},
		{10, "TL4+, postech with specialties",
			"Tech level 4 with specialties is an unusual case of a normal tech level 4 world that has retained some pretech industries or has developed their own local technical expertise beyond baseline postech in certain specific fields.\nMost such specialties are relatively narrow in scope: grav tech, medicine, hydroponics, force field generation, or some field of roughly similar breadth. These specialties are usually either the product of a few irreplaceable pretech manufactories or the result of some unique local resource that serves amazingly well for the purpose at hand. The world will jealously guard the tech, and much local conflict may relate to control over these resources.",
			"4+"},
		{11, "TL4+, postech with specialties",
			"Tech level 4 with specialties is an unusual case of a normal tech level 4 world that has retained some pretech industries or has developed their own local technical expertise beyond baseline postech in certain specific fields.\nMost such specialties are relatively narrow in scope: grav tech, medicine, hydroponics, force field generation, or some field of roughly similar breadth. These specialties are usually either the product of a few irreplaceable pretech manufactories or the result of some unique local resource that serves amazingly well for the purpose at hand. The world will jealously guard the tech, and much local conflict may relate to control over these resources.",
			"4+"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, u := range updates {
		if _, err := tx.Exec(
			"UPDATE WORLD_TECH_LEVEL SET name = ?, desc = ?, short = ? WHERE id = ?",
			u.name, u.desc, u.short, u.id,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ---------------------------------------------------------------------------
// NPC tables — complete replacement
// ---------------------------------------------------------------------------

func migrateNPCTables(db *sql.DB) error {
	stmts := []string{
		// Drop old tables
		"DROP TABLE IF EXISTS NPC_GENDER",
		"DROP TABLE IF EXISTS NPC_HEIGHT",
		"DROP TABLE IF EXISTS NPC_MOTIVATION",
		"DROP TABLE IF EXISTS NPC_QUIRKS",
		"DROP TABLE IF EXISTS NPC_AGE",
		"DROP TABLE IF EXISTS NPC_PROBLEMS",

		// Create new tables
		"CREATE TABLE NPC_AGE (id INTEGER PRIMARY KEY, age VARCHAR NOT NULL)",
		"CREATE TABLE NPC_BACKGROUND (id INTEGER PRIMARY KEY, background VARCHAR NOT NULL)",
		"CREATE TABLE NPC_ROLE (id INTEGER PRIMARY KEY, role VARCHAR NOT NULL)",
		"CREATE TABLE NPC_PROBLEM (id INTEGER PRIMARY KEY, problem VARCHAR NOT NULL)",
		"CREATE TABLE NPC_DESIRE (id INTEGER PRIMARY KEY, desire VARCHAR NOT NULL)",
		"CREATE TABLE NPC_TRAIT (id INTEGER PRIMARY KEY, trait VARCHAR NOT NULL)",

		// NPC_AGE (d4)
		"INSERT INTO NPC_AGE VALUES (1, 'Unusually young or old for their role')",
		"INSERT INTO NPC_AGE VALUES (2, 'Young adult')",
		"INSERT INTO NPC_AGE VALUES (3, 'Mature prime')",
		"INSERT INTO NPC_AGE VALUES (4, 'Middle-aged or elderly')",

		// NPC_BACKGROUND (d6)
		"INSERT INTO NPC_BACKGROUND VALUES (1, 'The local underclass or poorest natives')",
		"INSERT INTO NPC_BACKGROUND VALUES (2, 'Common laborers or cube workers')",
		"INSERT INTO NPC_BACKGROUND VALUES (3, 'Aspiring bourgeoisie or upper class')",
		"INSERT INTO NPC_BACKGROUND VALUES (4, 'The elite of this society')",
		"INSERT INTO NPC_BACKGROUND VALUES (5, 'Minority or foreigners; reroll on 1d4')",
		"INSERT INTO NPC_BACKGROUND VALUES (6, 'Offworlders or exotics; reroll on 1d4')",

		// NPC_ROLE (d8)
		"INSERT INTO NPC_ROLE VALUES (1, 'Criminal, thug, thief, swindler')",
		"INSERT INTO NPC_ROLE VALUES (2, 'Menial, cleaner, retail worker, servant')",
		"INSERT INTO NPC_ROLE VALUES (3, 'Unskilled heavy labor, porter, construction')",
		"INSERT INTO NPC_ROLE VALUES (4, 'Skilled trade, electrician, mechanic, pilot')",
		"INSERT INTO NPC_ROLE VALUES (5, 'Idea worker, programmer, writer')",
		"INSERT INTO NPC_ROLE VALUES (6, 'Merchant, business owner, trader, banker')",
		"INSERT INTO NPC_ROLE VALUES (7, 'Official, bureaucrat, courtier, clerk')",
		"INSERT INTO NPC_ROLE VALUES (8, 'Military, soldier, enforcer, law officer')",

		// NPC_PROBLEM (d10)
		"INSERT INTO NPC_PROBLEM VALUES (1, 'They have significant debt or money woes')",
		"INSERT INTO NPC_PROBLEM VALUES (2, 'A loved one is in trouble; reroll for it')",
		"INSERT INTO NPC_PROBLEM VALUES (3, 'Romantic failure with a desired person')",
		"INSERT INTO NPC_PROBLEM VALUES (4, 'Drug or behavioral addiction')",
		"INSERT INTO NPC_PROBLEM VALUES (5, 'Their superior dislikes or resents them')",
		"INSERT INTO NPC_PROBLEM VALUES (6, 'They have a persistent sickness')",
		"INSERT INTO NPC_PROBLEM VALUES (7, 'They hate their job or life situation')",
		"INSERT INTO NPC_PROBLEM VALUES (8, 'Someone dangerous is targeting them')",
		`INSERT INTO NPC_PROBLEM VALUES (9, 'They''re pursuing a disastrous purpose')`,
		"INSERT INTO NPC_PROBLEM VALUES (10, 'They have no problems worth mentioning')",

		// NPC_DESIRE (d12)
		"INSERT INTO NPC_DESIRE VALUES (1, 'They want a particular romantic partner')",
		"INSERT INTO NPC_DESIRE VALUES (2, 'They want money for them or a loved one')",
		"INSERT INTO NPC_DESIRE VALUES (3, 'They want a promotion in their job')",
		"INSERT INTO NPC_DESIRE VALUES (4, 'They want answers about a past trauma')",
		"INSERT INTO NPC_DESIRE VALUES (5, 'They want revenge on an enemy')",
		"INSERT INTO NPC_DESIRE VALUES (6, 'They want to help a beleaguered friend')",
		"INSERT INTO NPC_DESIRE VALUES (7, 'They want an entirely different job')",
		"INSERT INTO NPC_DESIRE VALUES (8, 'They want protection from an enemy')",
		"INSERT INTO NPC_DESIRE VALUES (9, 'They want to leave their current life')",
		"INSERT INTO NPC_DESIRE VALUES (10, 'They want fame and glory')",
		"INSERT INTO NPC_DESIRE VALUES (11, 'They want power over those around them')",
		"INSERT INTO NPC_DESIRE VALUES (12, 'They have everything they want from life')",

		// NPC_TRAIT (d20)
		"INSERT INTO NPC_TRAIT VALUES (1, 'Ambition')",
		"INSERT INTO NPC_TRAIT VALUES (2, 'Avarice')",
		"INSERT INTO NPC_TRAIT VALUES (3, 'Bitterness')",
		"INSERT INTO NPC_TRAIT VALUES (4, 'Courage')",
		"INSERT INTO NPC_TRAIT VALUES (5, 'Cowardice')",
		"INSERT INTO NPC_TRAIT VALUES (6, 'Curiosity')",
		"INSERT INTO NPC_TRAIT VALUES (7, 'Deceitfulness')",
		"INSERT INTO NPC_TRAIT VALUES (8, 'Determination')",
		"INSERT INTO NPC_TRAIT VALUES (9, 'Devotion to a cause')",
		"INSERT INTO NPC_TRAIT VALUES (10, 'Filiality')",
		"INSERT INTO NPC_TRAIT VALUES (11, 'Hatred')",
		"INSERT INTO NPC_TRAIT VALUES (12, 'Honesty')",
		"INSERT INTO NPC_TRAIT VALUES (13, 'Hopefulness')",
		"INSERT INTO NPC_TRAIT VALUES (14, 'Love of a person')",
		"INSERT INTO NPC_TRAIT VALUES (15, 'Nihilism')",
		"INSERT INTO NPC_TRAIT VALUES (16, 'Paternalism')",
		"INSERT INTO NPC_TRAIT VALUES (17, 'Pessimism')",
		"INSERT INTO NPC_TRAIT VALUES (18, 'Protectiveness')",
		"INSERT INTO NPC_TRAIT VALUES (19, 'Resentment')",
		"INSERT INTO NPC_TRAIT VALUES (20, 'Shame')",
	}
	return execAll(db, stmts)
}

// ---------------------------------------------------------------------------
// Alien Body Types — reorder and rename Humanlike → Mammalian
// ---------------------------------------------------------------------------

func migrateAlienBodyTypes(db *sql.DB) error {
	updates := []struct {
		id   int
		name string
		desc string
	}{
		{1, "Avian",
			"Avian aliens favor birdlike features. Feathered tails, wings, beaks, drastic coloring differences between males and females, skin either partially or fully covered in feathers, hollow bones, and clawed feet and hands are all traits that might be present in an avian-flavored alien. Such avians might have remarkably good vision or show unusual talent as pilots thanks to an evolutionary predilection for thinking in three dimensions."},
		{2, "Reptilian",
			"Reptilian aliens can range from the classic lizard man look to a more subtle combination of pebble-scaled skin, membrane frills, broad bodies, cold-bloodedness, nictitating membranes, an exceptional sense of smell, thick claws, powerful builds, sharp teeth, and fine scales. Bodies may be broad and thick like that of a lizard, or perhaps sinuous like that of a snake. Reptilian aliens might be found especially on damp, watery worlds, and might be outright amphibious over all or part of their life spans."},
		{3, "Insectile",
			`Insectile aliens are usually possessed of chitinous skin and one or more pairs of compound eyes. Membranous wings aren't unknown, but are probably vestigial and useless on all but the smallest or lightest aliens. Insectile aliens might have sucking mouthparts if their primary food source is plant sap or blood, or might use chewing mandibles instead for a more solid diet.`},
		{4, "Mammalian",
			"Mammalian aliens might have very human-like appearances, or they might have fur, muzzles, claws, or quadrupedal modes of locomotion. Mammalian aliens are easy to conceptualize and describe, but if you want to play up the alienness of them, you'll likely need to make sure they have some dramatic social or cultural difference, or some unique element of physiology."},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, u := range updates {
		if _, err := tx.Exec(
			"UPDATE ALIEN_BODY_TYPE SET name = ?, desc = ? WHERE id = ?",
			u.name, u.desc, u.id,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ---------------------------------------------------------------------------
// Alien Lenses — reorder (Fear moves from 6 to 20)
// ---------------------------------------------------------------------------

func migrateAlienLenses(db *sql.DB) error {
	// Read all existing lens descriptions keyed by name
	rows, err := db.Query("SELECT name, desc FROM ALIEN_LENSES ORDER BY id")
	if err != nil {
		return err
	}
	lensDescs := make(map[string]string)
	for rows.Next() {
		var name, desc string
		if err := rows.Scan(&name, &desc); err != nil {
			rows.Close()
			return err
		}
		lensDescs[name] = desc
	}
	rows.Close()

	// Revised ordering
	revisedOrder := []string{
		"Collectivity", "Curiosity", "Despair", "Domination", "Faith",
		"Gluttony", "Greed", "Hate", "Honor", "Journeying",
		"Joy", "Pacifism", "Pride", "Sagacity", "Subtlety",
		"Tradition", "Treachery", "Tribalism", "Wrath", "Fear",
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for i, name := range revisedOrder {
		id := i + 1
		desc := lensDescs[name]
		if _, err := tx.Exec(
			"UPDATE ALIEN_LENSES SET name = ?, desc = ? WHERE id = ?",
			name, desc, id,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// ---------------------------------------------------------------------------
// Alien Structure — rename Monarchic, add Competitive (rows 7-8)
// ---------------------------------------------------------------------------

func migrateAlienStructure(db *sql.DB) error {
	stmts := []string{
		"UPDATE ALIEN_STRUCTURE SET name = 'Monarchy' WHERE id = 2",
		`INSERT OR REPLACE INTO ALIEN_STRUCTURE (id, name, desc) VALUES (7, 'Competitive', 'Individual members or small groups vie for dominance through direct competition, whether through combat, economic contests, intellectual challenges, or other metrics of personal excellence. Status and leadership positions are earned through demonstrated superiority rather than inheritance, election, or tradition. These societies tend to be dynamic and innovative but also prone to internal conflict and instability as ambitious individuals constantly challenge the established order.')`,
		`INSERT OR REPLACE INTO ALIEN_STRUCTURE (id, name, desc) VALUES (8, 'Competitive', 'Individual members or small groups vie for dominance through direct competition, whether through combat, economic contests, intellectual challenges, or other metrics of personal excellence. Status and leadership positions are earned through demonstrated superiority rather than inheritance, election, or tradition. These societies tend to be dynamic and innovative but also prone to internal conflict and instability as ambitious individuals constantly challenge the established order.')`,
	}
	return execAll(db, stmts)
}

// ---------------------------------------------------------------------------
// World Tags — replace all with revised 99-tag set
// ---------------------------------------------------------------------------

func migrateWorldTags(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM WORLD_TAG"); err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO WORLD_TAG (id, name, desc, enemies, friends, complications, things, places, short) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, tag := range worldTags {
		short := tag.Name
		if len(short) > 32 {
			short = short[:32]
		}
		if _, err := stmt.Exec(tag.ID, tag.Name, tag.Desc, tag.Enemies, tag.Friends, tag.Complications, tag.Things, tag.Places, short); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert tag %d (%s): %w", tag.ID, tag.Name, err)
		}
	}

	return tx.Commit()
}

// ---------------------------------------------------------------------------
// Patron tables
// ---------------------------------------------------------------------------

func migratePatronTables(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS PATRON_TRUST (id INTEGER PRIMARY KEY, trust VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS PATRON_EAGERNESS (id INTEGER PRIMARY KEY, eagerness VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS PATRON_REWARD (id INTEGER PRIMARY KEY, reward VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS PATRON_CHALLENGE (id INTEGER PRIMARY KEY, challenge VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS PATRON_OPPOSITION (id INTEGER PRIMARY KEY, opposition VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS PATRON_COMPLICATION (id INTEGER PRIMARY KEY, complication VARCHAR NOT NULL)",

		// PATRON_TRUST (d6)
		"INSERT INTO PATRON_TRUST VALUES (1, 'They intend to totally screw the PCs')",
		`INSERT INTO PATRON_TRUST VALUES (2, 'They won''t pay unless forced to do so')`,
		`INSERT INTO PATRON_TRUST VALUES (3, 'They''ll pay slowly or reluctantly')`,
		`INSERT INTO PATRON_TRUST VALUES (4, 'They''ll pay, but discount for mistakes')`,
		`INSERT INTO PATRON_TRUST VALUES (5, 'They''ll pay without quibbling')`,
		`INSERT INTO PATRON_TRUST VALUES (6, 'They''ll pay more than they promised')`,

		// PATRON_EAGERNESS (d4)
		`INSERT INTO PATRON_EAGERNESS VALUES (1, 'Cautious, but can be convinced to hire')`,
		"INSERT INTO PATRON_EAGERNESS VALUES (2, 'Willing to promise standard rates')",
		"INSERT INTO PATRON_EAGERNESS VALUES (3, 'Eager, willing to offer a bonus')",
		`INSERT INTO PATRON_EAGERNESS VALUES (4, 'Desperate, might offer what they can''t pay')`,

		// PATRON_REWARD (d12)
		"INSERT INTO PATRON_REWARD VALUES (1, 'Government official favors owed')",
		"INSERT INTO PATRON_REWARD VALUES (2, 'Property in the area')",
		"INSERT INTO PATRON_REWARD VALUES (3, 'An item very valuable on another world')",
		"INSERT INTO PATRON_REWARD VALUES (4, 'Pretech mod components')",
		"INSERT INTO PATRON_REWARD VALUES (5, 'Useful pretech artifact')",
		"INSERT INTO PATRON_REWARD VALUES (6, 'Information the PCs need')",
		"INSERT INTO PATRON_REWARD VALUES (7, 'Membership in a powerful group')",
		"INSERT INTO PATRON_REWARD VALUES (8, 'Black market access')",
		"INSERT INTO PATRON_REWARD VALUES (9, 'Use of restricted facilities or shipyards')",
		"INSERT INTO PATRON_REWARD VALUES (10, 'Shares in a profitable business')",
		"INSERT INTO PATRON_REWARD VALUES (11, 'Maps to a hidden or guarded treasure')",
		"INSERT INTO PATRON_REWARD VALUES (12, 'Illegal but valuable weapons or gear')",

		// PATRON_CHALLENGE (d8)
		"INSERT INTO PATRON_CHALLENGE VALUES (1, 'Kill somebody who might deserve it')",
		"INSERT INTO PATRON_CHALLENGE VALUES (2, 'Kidnap someone dangerous')",
		"INSERT INTO PATRON_CHALLENGE VALUES (3, 'Steal a well-guarded object')",
		"INSERT INTO PATRON_CHALLENGE VALUES (4, 'Arson or sabotage on a place')",
		"INSERT INTO PATRON_CHALLENGE VALUES (5, 'Get proof of some misdeed')",
		"INSERT INTO PATRON_CHALLENGE VALUES (6, 'Protect someone from an immediate threat')",
		"INSERT INTO PATRON_CHALLENGE VALUES (7, 'Transport someone through danger')",
		"INSERT INTO PATRON_CHALLENGE VALUES (8, 'Guard an object being transported')",

		// PATRON_OPPOSITION (d10)
		"INSERT INTO PATRON_OPPOSITION VALUES (1, 'A treacherous employer or subordinate')",
		"INSERT INTO PATRON_OPPOSITION VALUES (2, 'An open and known enemy of the patron')",
		"INSERT INTO PATRON_OPPOSITION VALUES (3, 'Official governmental meddling')",
		"INSERT INTO PATRON_OPPOSITION VALUES (4, 'An unknown rival of the patron')",
		"INSERT INTO PATRON_OPPOSITION VALUES (5, 'The macguffin itself opposes them')",
		"INSERT INTO PATRON_OPPOSITION VALUES (6, 'Very short time frame allowed')",
		"INSERT INTO PATRON_OPPOSITION VALUES (7, 'The job is spectacularly illegal')",
		"INSERT INTO PATRON_OPPOSITION VALUES (8, 'A participant would profit by their failure')",
		"INSERT INTO PATRON_OPPOSITION VALUES (9, 'The patron is badly wrong about a thing')",
		"INSERT INTO PATRON_OPPOSITION VALUES (10, 'The locals are against the patron')",

		// PATRON_COMPLICATION (d20)
		"INSERT INTO PATRON_COMPLICATION VALUES (1, 'An ambush is laid somewhere')",
		"INSERT INTO PATRON_COMPLICATION VALUES (2, 'PC involvement is leaked to the enemy')",
		"INSERT INTO PATRON_COMPLICATION VALUES (3, 'The patron gives faulty aid somehow')",
		"INSERT INTO PATRON_COMPLICATION VALUES (4, 'Failing would be extremely unhealthy')",
		"INSERT INTO PATRON_COMPLICATION VALUES (5, 'The job IDs them as allies of a local faction')",
		"INSERT INTO PATRON_COMPLICATION VALUES (6, 'The macguffin is physically dangerous')",
		"INSERT INTO PATRON_COMPLICATION VALUES (7, 'An important location is hard to get into')",
		"INSERT INTO PATRON_COMPLICATION VALUES (8, 'Succeeding would be morally distasteful')",
		"INSERT INTO PATRON_COMPLICATION VALUES (9, 'A supposed ally is very unhelpful or stupid')",
		"INSERT INTO PATRON_COMPLICATION VALUES (10, 'The patron badly misunderstood the PCs')",
		"INSERT INTO PATRON_COMPLICATION VALUES (11, 'The job changes suddenly partway through')",
		"INSERT INTO PATRON_COMPLICATION VALUES (12, 'An unexpected troublemaker is involved')",
		"INSERT INTO PATRON_COMPLICATION VALUES (13, 'Critical gear will fail partway through')",
		"INSERT INTO PATRON_COMPLICATION VALUES (14, 'An unrelated accident complicates things')",
		"INSERT INTO PATRON_COMPLICATION VALUES (15, 'Payment comes in a hard-to-handle form')",
		"INSERT INTO PATRON_COMPLICATION VALUES (16, 'Someone is turning traitor on the patron')",
		"INSERT INTO PATRON_COMPLICATION VALUES (17, 'A critical element has suddenly moved')",
		"INSERT INTO PATRON_COMPLICATION VALUES (18, 'Payment is in avidly-pursued hot goods')",
		"INSERT INTO PATRON_COMPLICATION VALUES (19, 'The true goal is a subsidiary part of the job')",
		`INSERT INTO PATRON_COMPLICATION VALUES (20, 'No complications; it''s just as it seems to be')`,
	}
	return execAll(db, stmts)
}

// ---------------------------------------------------------------------------
// Urban Encounter tables
// ---------------------------------------------------------------------------

func migrateUrbanEncounterTables(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS URBAN_CONFLICT (id INTEGER PRIMARY KEY, conflict VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS URBAN_ANTAGONIST (id INTEGER PRIMARY KEY, antagonist VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS URBAN_VENUE (id INTEGER PRIMARY KEY, venue VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS URBAN_PC_INVOLVE (id INTEGER PRIMARY KEY, involve VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS URBAN_EVENT (id INTEGER PRIMARY KEY, event VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS URBAN_FEATURE (id INTEGER PRIMARY KEY, feature VARCHAR NOT NULL)",

		// URBAN_CONFLICT (d4)
		"INSERT INTO URBAN_CONFLICT VALUES (1, 'Money, extortion, payment due, debts')",
		"INSERT INTO URBAN_CONFLICT VALUES (2, 'Respect, submission to social authority')",
		"INSERT INTO URBAN_CONFLICT VALUES (3, 'Grudges, ethnic resentment, gang payback')",
		"INSERT INTO URBAN_CONFLICT VALUES (4, 'Politics, religion, or other ideology')",

		// URBAN_ANTAGONIST (d12)
		"INSERT INTO URBAN_ANTAGONIST VALUES (1, 'A local bully and their thugs')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (2, 'A ruthless political boss and their zealots')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (3, 'Violent criminals')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (4, 'Religious fanatics')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (5, 'A blisteringly obnoxious offworlder')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (6, 'Corrupt or over-strict government official')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (7, 'A mob of intoxicated locals')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (8, 'A ranting demagogue and their followers')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (9, 'A stupidly bull-headed local grandee')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (10, 'A very capable assassin or strong-arm')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (11, 'A self-centered local scion of power')",
		"INSERT INTO URBAN_ANTAGONIST VALUES (12, 'A confused foreigner or backwoodsman')",

		// URBAN_VENUE (d6)
		"INSERT INTO URBAN_VENUE VALUES (1, 'In the middle of the street')",
		"INSERT INTO URBAN_VENUE VALUES (2, 'In a public plaza')",
		"INSERT INTO URBAN_VENUE VALUES (3, 'Down a side alley')",
		"INSERT INTO URBAN_VENUE VALUES (4, 'Inside a local business')",
		"INSERT INTO URBAN_VENUE VALUES (5, 'Next to or in a public park')",
		"INSERT INTO URBAN_VENUE VALUES (6, 'At a mass-transit station')",

		// URBAN_PC_INVOLVE (d8)
		"INSERT INTO URBAN_PC_INVOLVE VALUES (1, 'A sympathetic participant appeals to them')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (2, 'Ways around it are all dangerous/blocked')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (3, 'It happens immediately around them')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (4, 'A valuable thing looks snatchable amid it')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (5, 'A participant offers a reward for help')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (6, 'Someone mistakenly involves the PCs in it')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (7, 'The seeming way out just leads deeper in')",
		"INSERT INTO URBAN_PC_INVOLVE VALUES (8, 'Responsibility is somehow pinned on them')",

		// URBAN_EVENT (d10)
		"INSERT INTO URBAN_EVENT VALUES (1, 'A parade or festival is being disrupted')",
		"INSERT INTO URBAN_EVENT VALUES (2, 'Innocents are being assaulted')",
		"INSERT INTO URBAN_EVENT VALUES (3, 'An establishment is being robbed')",
		"INSERT INTO URBAN_EVENT VALUES (4, 'A disturbance over local politics happens')",
		"INSERT INTO URBAN_EVENT VALUES (5, 'Someone is being blamed for something')",
		"INSERT INTO URBAN_EVENT VALUES (6, 'Fires or building collapses are happening')",
		"INSERT INTO URBAN_EVENT VALUES (7, 'A medical emergency is happening')",
		"INSERT INTO URBAN_EVENT VALUES (8, 'Someone is trying to cheat the PCs')",
		"INSERT INTO URBAN_EVENT VALUES (9, 'A vehicle accident is happening')",
		"INSERT INTO URBAN_EVENT VALUES (10, 'A religious ceremony is being disrupted')",

		// URBAN_FEATURE (d20)
		"INSERT INTO URBAN_FEATURE VALUES (1, 'Heavy traffic running through the place')",
		"INSERT INTO URBAN_FEATURE VALUES (2, 'Music blaring at deafening volumes')",
		"INSERT INTO URBAN_FEATURE VALUES (3, 'Two groups present that detest each other')",
		"INSERT INTO URBAN_FEATURE VALUES (4, 'Large delivery taking place right there')",
		"INSERT INTO URBAN_FEATURE VALUES (5, 'Swarm of schoolkids or feral youth')",
		"INSERT INTO URBAN_FEATURE VALUES (6, 'Insistent soapbox preacher here')",
		"INSERT INTO URBAN_FEATURE VALUES (7, 'Several pickpockets working the crowd')",
		"INSERT INTO URBAN_FEATURE VALUES (8, 'A kiosk is tipping over and spilling things')",
		"INSERT INTO URBAN_FEATURE VALUES (9, 'Streetlights are out or visibility is low')",
		"INSERT INTO URBAN_FEATURE VALUES (10, 'A cop patrol is here and reluctant to act')",
		"INSERT INTO URBAN_FEATURE VALUES (11, 'PC-hostile reporters are recording here')",
		`INSERT INTO URBAN_FEATURE VALUES (12, 'Someone''s trying to sell something to PCs')`,
		"INSERT INTO URBAN_FEATURE VALUES (13, 'Feral dogs or other animals crowd here')",
		"INSERT INTO URBAN_FEATURE VALUES (14, 'Unrelated activists are protesting here')",
		"INSERT INTO URBAN_FEATURE VALUES (15, 'Street kids are trying to steal from the PCs')",
		"INSERT INTO URBAN_FEATURE VALUES (16, 'GPS maps are dangerously wrong here')",
		"INSERT INTO URBAN_FEATURE VALUES (17, 'Downed power lines are a danger here')",
		"INSERT INTO URBAN_FEATURE VALUES (18, 'Numerous open manholes and utility holes')",
		`INSERT INTO URBAN_FEATURE VALUES (19, 'The street''s blockaded by something')`,
		"INSERT INTO URBAN_FEATURE VALUES (20, 'Crowds so thick one can barely move')",
	}
	return execAll(db, stmts)
}

// ---------------------------------------------------------------------------
// Wilderness Encounter tables
// ---------------------------------------------------------------------------

func migrateWildEncounterTables(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS WILD_RANGE (id INTEGER PRIMARY KEY, range_desc VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS WILD_HOSTILE (id INTEGER PRIMARY KEY, hostile VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS WILD_WEATHER (id INTEGER PRIMARY KEY, weather VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS WILD_NATURE (id INTEGER PRIMARY KEY, nature VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS WILD_FRIENDLY (id INTEGER PRIMARY KEY, friendly VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS WILD_FEATURE (id INTEGER PRIMARY KEY, feature VARCHAR NOT NULL)",

		// WILD_RANGE (d4)
		"INSERT INTO WILD_RANGE VALUES (1, 'Visible from a long distance away')",
		"INSERT INTO WILD_RANGE VALUES (2, 'Noticed 1d4 hundred meters away')",
		"INSERT INTO WILD_RANGE VALUES (3, 'Noticed only within 1d6 x 10 meters')",
		"INSERT INTO WILD_RANGE VALUES (4, 'Noticed only when adjacent to the event')",

		// WILD_HOSTILE (d12)
		"INSERT INTO WILD_HOSTILE VALUES (1, 'Bandits in their wilderness hideout')",
		"INSERT INTO WILD_HOSTILE VALUES (2, 'Dangerous locals looking for easy marks')",
		"INSERT INTO WILD_HOSTILE VALUES (3, 'Rabid or diseased large predator')",
		"INSERT INTO WILD_HOSTILE VALUES (4, 'Pack of hungry hunting beasts')",
		"INSERT INTO WILD_HOSTILE VALUES (5, 'Herd of potentially dangerous prey animals')",
		"INSERT INTO WILD_HOSTILE VALUES (6, 'Swarm of dangerous vermin')",
		"INSERT INTO WILD_HOSTILE VALUES (7, 'Criminal seeking to evade the law')",
		"INSERT INTO WILD_HOSTILE VALUES (8, 'Brutal local landowner and their men')",
		"INSERT INTO WILD_HOSTILE VALUES (9, 'Crazed hermit seeking enforced solitude')",
		"INSERT INTO WILD_HOSTILE VALUES (10, 'Friendly-seeming guide into lethal danger')",
		"INSERT INTO WILD_HOSTILE VALUES (11, 'Harmless-looking but dangerous beast')",
		"INSERT INTO WILD_HOSTILE VALUES (12, 'Confidence man seeking to gull the PCs')",

		// WILD_WEATHER (d6)
		"INSERT INTO WILD_WEATHER VALUES (1, 'Takes place in daylight and clear weather')",
		"INSERT INTO WILD_WEATHER VALUES (2, 'Daylight, but fog, mist, rain or the like')",
		"INSERT INTO WILD_WEATHER VALUES (3, 'Daylight, but harsh seasonal weather')",
		"INSERT INTO WILD_WEATHER VALUES (4, 'Night encounter, but clear weather')",
		"INSERT INTO WILD_WEATHER VALUES (5, 'Night, with rain or other obscuring effects')",
		"INSERT INTO WILD_WEATHER VALUES (6, 'Night, with terrible weather and wind')",

		// WILD_NATURE (d8)
		"INSERT INTO WILD_NATURE VALUES (1, 'Attack by pack of hostiles')",
		"INSERT INTO WILD_NATURE VALUES (2, 'Ambush by single lone hostile')",
		`INSERT INTO WILD_NATURE VALUES (3, 'Meet people who don''t want to be met')`,
		"INSERT INTO WILD_NATURE VALUES (4, 'Encounter people in need of aid')",
		"INSERT INTO WILD_NATURE VALUES (5, 'Encounter hostile creatures')",
		"INSERT INTO WILD_NATURE VALUES (6, 'Nearby feature is somehow dangerous')",
		"INSERT INTO WILD_NATURE VALUES (7, 'Nearby feature promises useful loot')",
		`INSERT INTO WILD_NATURE VALUES (8, 'Meet hostiles that aren''t immediately so')`,

		// WILD_FRIENDLY (d10)
		"INSERT INTO WILD_FRIENDLY VALUES (1, 'Affable but reclusive hermit')",
		"INSERT INTO WILD_FRIENDLY VALUES (2, 'Local herd animal let loose to graze')",
		"INSERT INTO WILD_FRIENDLY VALUES (3, 'Government ranger or circuit judge')",
		"INSERT INTO WILD_FRIENDLY VALUES (4, 'Curious local animal')",
		"INSERT INTO WILD_FRIENDLY VALUES (5, 'Remote homesteader and family')",
		"INSERT INTO WILD_FRIENDLY VALUES (6, 'Working trapper or hunter')",
		"INSERT INTO WILD_FRIENDLY VALUES (7, 'Back-country villager or native')",
		"INSERT INTO WILD_FRIENDLY VALUES (8, 'Hiker or wilderness tourist')",
		"INSERT INTO WILD_FRIENDLY VALUES (9, 'Religious recluse or holy person')",
		"INSERT INTO WILD_FRIENDLY VALUES (10, 'Impoverished social exile')",

		// WILD_FEATURE (d20)
		"INSERT INTO WILD_FEATURE VALUES (1, 'Overgrown homestead')",
		"INSERT INTO WILD_FEATURE VALUES (2, 'Stream prone to flash-flooding')",
		"INSERT INTO WILD_FEATURE VALUES (3, 'Narrow bridge or beam over deep cleft')",
		"INSERT INTO WILD_FEATURE VALUES (4, 'Box canyon with steep sides')",
		"INSERT INTO WILD_FEATURE VALUES (5, 'Unstable hillside that slides if disturbed')",
		"INSERT INTO WILD_FEATURE VALUES (6, 'Long-lost crash site of a gravflyer')",
		"INSERT INTO WILD_FEATURE VALUES (7, 'Once-inhabited cave or tunnel')",
		"INSERT INTO WILD_FEATURE VALUES (8, 'Steep and dangerous cliff')",
		"INSERT INTO WILD_FEATURE VALUES (9, 'Quicksand-laden swamp or dust pit')",
		"INSERT INTO WILD_FEATURE VALUES (10, 'Ruins of a ghost town or lost hamlet')",
		"INSERT INTO WILD_FEATURE VALUES (11, 'Hunting cabin with necessities')",
		"INSERT INTO WILD_FEATURE VALUES (12, 'Ill-tended graveyard of a lost family stead')",
		`INSERT INTO WILD_FEATURE VALUES (13, 'Narrow pass that''s easily blocked')`,
		"INSERT INTO WILD_FEATURE VALUES (14, 'Dilapidated resort building')",
		"INSERT INTO WILD_FEATURE VALUES (15, 'Remote government monitoring outpost')",
		"INSERT INTO WILD_FEATURE VALUES (16, 'Illicit substance farm or processing center')",
		"INSERT INTO WILD_FEATURE VALUES (17, 'Old and forgotten battleground')",
		"INSERT INTO WILD_FEATURE VALUES (18, 'Zone overrun by dangerous plants')",
		"INSERT INTO WILD_FEATURE VALUES (19, 'Thick growth that lights up at a spark')",
		"INSERT INTO WILD_FEATURE VALUES (20, 'Abandoned vehicle')",
	}
	return execAll(db, stmts)
}

// ---------------------------------------------------------------------------
// POI tables — placeholder, will be populated from PDF data
// ---------------------------------------------------------------------------

func migratePOITables(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS POI_TYPE (id INTEGER PRIMARY KEY, name VARCHAR(64) NOT NULL, desc CLOB NOT NULL)",
		"CREATE TABLE IF NOT EXISTS POI_OCCUPANT (id INTEGER PRIMARY KEY, type_id INTEGER NOT NULL, occupant VARCHAR(128) NOT NULL)",
		"CREATE TABLE IF NOT EXISTS POI_SITUATION (id INTEGER PRIMARY KEY, type_id INTEGER NOT NULL, situation VARCHAR(128) NOT NULL)",
	}
	if err := execAll(db, stmts); err != nil {
		return err
	}
	return migratePOIData(db)
}

// ---------------------------------------------------------------------------
// Beast tables — placeholder, will be populated from PDF data
// ---------------------------------------------------------------------------

func migrateBeastTables(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS BEAST_FEATURE (id INTEGER PRIMARY KEY, feature VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_BODY_PLAN (id INTEGER PRIMARY KEY, plan VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_LIMB (id INTEGER PRIMARY KEY, limb VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_SKIN (id INTEGER PRIMARY KEY, skin VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_WEAPON (id INTEGER PRIMARY KEY, weapon VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_SIZE (id INTEGER PRIMARY KEY, size VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_BEHAVIOR (id INTEGER PRIMARY KEY, type VARCHAR NOT NULL, behavior VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_DISCHARGE (id INTEGER PRIMARY KEY, discharge VARCHAR NOT NULL)",
		"CREATE TABLE IF NOT EXISTS BEAST_POISON (id INTEGER PRIMARY KEY, poison VARCHAR NOT NULL, onset VARCHAR NOT NULL, duration VARCHAR NOT NULL)",
	}
	if err := execAll(db, stmts); err != nil {
		return err
	}
	return migrateBeastData(db)
}

// ---------------------------------------------------------------------------
// Faction tables — placeholder, will be populated from PDF data
// ---------------------------------------------------------------------------

func migrateFactionTables(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS FACTION_TAG (id INTEGER PRIMARY KEY, name VARCHAR(64) NOT NULL, desc CLOB NOT NULL, effect CLOB NOT NULL)",
		"CREATE TABLE IF NOT EXISTS FACTION_GOAL (id INTEGER PRIMARY KEY, name VARCHAR(64) NOT NULL, desc CLOB NOT NULL)",
	}
	if err := execAll(db, stmts); err != nil {
		return err
	}
	return migrateFactionData(db)
}

// ---------------------------------------------------------------------------
// World Relationship table
// ---------------------------------------------------------------------------

func migrateWorldRelationship(db *sql.DB) error {
	stmts := []string{
		"CREATE TABLE IF NOT EXISTS WORLD_RELATIONSHIP (id INTEGER PRIMARY KEY, origin CLOB NOT NULL, relationship CLOB NOT NULL, contact_point CLOB NOT NULL)",

		`INSERT INTO WORLD_RELATIONSHIP VALUES (1, 'Recent colony from the primary world', 'Confirmed hatred of each other', 'Trade in vital goods')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (2, 'Refuge for exiles from primary', 'Active cold war between them', 'Shared religion')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (3, 'Founded ages ago by a different group', 'Old grudges or resentments', 'Mutual language')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (4, 'Founded long before the primary world', 'Cultural disgust and avoidance', 'Entertainment content')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (5, 'Lost ancient colony of the primary', 'Polite interchange and trade', 'Shared research')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (6, 'Colony recently torn free of the primary', 'Cultural admiration for primary', 'Threat to both of them')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (7, 'Long-standing cooperative colony world', 'Long-standing friendship', 'Shared elite families')`,
		`INSERT INTO WORLD_RELATIONSHIP VALUES (8, 'Recent interstellar colony from elsewhere', 'Unflinching mutual loyalty', 'Exploiting shared resource')`,
	}
	return execAll(db, stmts)
}
