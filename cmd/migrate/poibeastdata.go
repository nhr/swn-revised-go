package main

import "database/sql"

// migratePOIData inserts POI type, occupant, and situation data.
func migratePOIData(db *sql.DB) error {
	types := []struct {
		id   int
		name string
		desc string
	}{
		{1, "Deep-space station", "A station situated in the deep void between worlds, often relying on ancient pretech for its continued function."},
		{2, "Asteroid base", "A base carved into or constructed on an asteroid, often used for mining or as a hidden refuge."},
		{3, "Remote moon base", "A base on one of the moons of a gas giant or other large body, isolated and hard to reach."},
		{4, "Ancient orbital ruin", "The remains of an alien or pretech orbital structure, still drifting in the system."},
		{5, "Research base", "A facility dedicated to scientific research, often in a remote or dangerous location."},
		{6, "Asteroid belt", "A field of tumbling rocks and debris, sometimes harboring miners or hidden bases."},
		{7, "Gas giant mine", "A mining operation suspended in the atmosphere of a gas giant, extracting valuable gases."},
		{8, "Refueling station", "A waypoint for ships to take on fuel, often the only sign of civilization in an empty system."},
	}

	occupants := []struct {
		id      int
		typeID  int
		occupant string
	}{
		// Deep-space station
		{101, 1, "Dangerously odd transhumans"},
		{102, 1, "Freeze-dried ancient corpses"},
		{103, 1, "Secretive military observers"},
		{104, 1, "Eccentric oligarch and minions"},
		{105, 1, "Deranged but brilliant scientist"},
		// Asteroid base
		{201, 2, "Zealous religious sectarians"},
		{202, 2, "Failed rebels from another world"},
		{203, 2, "Wage-slave corporate miners"},
		{204, 2, "Independent asteroid prospectors"},
		{205, 2, "Pirates masquerading as otherwise"},
		// Remote moon base
		{301, 3, "Unlucky corporate researchers"},
		{302, 3, "Reclusive hermit genius"},
		{303, 3, "Remnants of a failed colony"},
		{304, 3, "Military listening post"},
		{305, 3, "Lonely overseers and robot miners"},
		// Ancient orbital ruin
		{401, 4, "Robots of dubious sentience"},
		{402, 4, "Trigger-happy scavengers"},
		{403, 4, "Government researchers"},
		{404, 4, "Military quarantine enforcers"},
		{405, 4, "Heirs of the original alien builders"},
		// Research base
		{501, 5, "Experiments that have gotten loose"},
		{502, 5, "Scientists from a major local corp"},
		{503, 5, "Black-ops governmental researchers"},
		{504, 5, "Secret employees of a foreign power"},
		{505, 5, "Aliens studying the human locals"},
		// Asteroid belt
		{601, 6, "Grizzled belter mine laborers"},
		{602, 6, "Ancient automated guardian drones"},
		{603, 6, "Survivors of destroyed asteroid base"},
		{604, 6, "Pirates hiding out among the rocks"},
		{605, 6, "Lonely military patrol base staff"},
		// Gas giant mine
		{701, 7, "Miserable gas-miner slaves or serfs"},
		{702, 7, "Strange robots and their overseers"},
		{703, 7, "Scientists studying the alien life"},
		{704, 7, "Scrappers in the ruined old mine"},
		{705, 7, "Impoverished separatist group"},
		// Refueling station
		{801, 8, "Half-crazed hermit caretaker"},
		{802, 8, "Sordid purveyors of decadent fun"},
		{803, 8, "Extortionate corporate minions"},
		{804, 8, "Religious missionaries to travelers"},
		{805, 8, "Brainless automated vendors"},
	}

	situations := []struct {
		id        int
		typeID    int
		situation string
	}{
		// Deep-space station
		{101, 1, "Systems breaking down"},
		{102, 1, "Foreign sabotage attempt"},
		{103, 1, "Black market for the elite"},
		{104, 1, "Vault for dangerous pretech"},
		{105, 1, "Supply base for pirates"},
		// Asteroid base
		{201, 2, "Life support is threatened"},
		{202, 2, "Base needs a new asteroid"},
		{203, 2, "Dug out something nasty"},
		{204, 2, "Fighting another asteroid"},
		{205, 2, "Hit a priceless vein of ore"},
		// Remote moon base
		{301, 3, "Something dark has awoken"},
		{302, 3, "Criminals trying to take over"},
		{303, 3, "Moon plague breaking out"},
		{304, 3, "Desperate for vital supplies"},
		{305, 3, "Rich but badly-protected"},
		// Ancient orbital ruin
		{401, 4, "Trying to stop it awakening"},
		{402, 4, "Meddling with strange tech"},
		{403, 4, "Impending tech calamity"},
		{404, 4, "A terrible secret is unearthed"},
		{405, 4, "Fighting outside interlopers"},
		// Research base
		{501, 5, "Perilous research underway"},
		{502, 5, "Hideously immoral research"},
		{503, 5, "Held hostage by outsiders"},
		{504, 5, "Science monsters run amok"},
		{505, 5, "Selling black-market tech"},
		// Asteroid belt
		{601, 6, "Ruptured rock released a peril"},
		{602, 6, "Foreign spy ships hide there"},
		{603, 6, "Gold rush for new minerals"},
		{604, 6, "Ancient ruins dot the rocks"},
		{605, 6, "War between rival rocks"},
		// Gas giant mine
		{701, 7, "Things are emerging below"},
		{702, 7, "They need vital supplies"},
		{703, 7, "The workers are in revolt"},
		{704, 7, "Pirates secretly fuel there"},
		{705, 7, "Alien remnants were found"},
		// Refueling station
		{801, 8, "A ship is in severe distress"},
		{802, 8, "Pirates have taken over"},
		{803, 8, "Has corrupt customs agents"},
		{804, 8, "Foreign saboteurs are active"},
		{805, 8, "Deep-space alien signal"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, t := range types {
		if _, err := tx.Exec("INSERT INTO POI_TYPE VALUES (?, ?, ?)", t.id, t.name, t.desc); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, o := range occupants {
		if _, err := tx.Exec("INSERT INTO POI_OCCUPANT VALUES (?, ?, ?)", o.id, o.typeID, o.occupant); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, s := range situations {
		if _, err := tx.Exec("INSERT INTO POI_SITUATION VALUES (?, ?, ?)", s.id, s.typeID, s.situation); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// migrateBeastData inserts xenobestiary table data.
func migrateBeastData(db *sql.DB) error {
	stmts := []string{
		// BEAST_FEATURE (d10: 1-8 unique, 9-10 = roll twice)
		"INSERT INTO BEAST_FEATURE VALUES (1, 'Amphibian, froglike or toadlike')",
		"INSERT INTO BEAST_FEATURE VALUES (2, 'Bird, winged and feathered')",
		"INSERT INTO BEAST_FEATURE VALUES (3, 'Fish, scaled and possibly finned')",
		"INSERT INTO BEAST_FEATURE VALUES (4, 'Insect, beetle-like or fly-winged')",
		"INSERT INTO BEAST_FEATURE VALUES (5, 'Mammal, hairy and fanged')",
		"INSERT INTO BEAST_FEATURE VALUES (6, 'Reptile, scaled and sharp-bodied')",
		"INSERT INTO BEAST_FEATURE VALUES (7, 'Spider, many-legged and fat')",
		"INSERT INTO BEAST_FEATURE VALUES (8, 'Exotic, made of wholly alien elements')",

		// BEAST_BODY_PLAN (d6: 1-5 unique, 6 = roll twice)
		"INSERT INTO BEAST_BODY_PLAN VALUES (1, 'Humanoid')",
		"INSERT INTO BEAST_BODY_PLAN VALUES (2, 'Quadruped')",
		"INSERT INTO BEAST_BODY_PLAN VALUES (3, 'Many-legged')",
		"INSERT INTO BEAST_BODY_PLAN VALUES (4, 'Bulbous')",
		"INSERT INTO BEAST_BODY_PLAN VALUES (5, 'Amorphous')",

		// BEAST_LIMB (d6)
		"INSERT INTO BEAST_LIMB VALUES (1, 'Wings')",
		"INSERT INTO BEAST_LIMB VALUES (2, 'Many joints')",
		"INSERT INTO BEAST_LIMB VALUES (3, 'Tentacles')",
		"INSERT INTO BEAST_LIMB VALUES (4, 'Opposable thumbs')",
		"INSERT INTO BEAST_LIMB VALUES (5, 'Retractable')",
		"INSERT INTO BEAST_LIMB VALUES (6, 'Varying sizes')",

		// BEAST_SKIN (d6)
		"INSERT INTO BEAST_SKIN VALUES (1, 'Hard shell')",
		"INSERT INTO BEAST_SKIN VALUES (2, 'Exoskeleton')",
		"INSERT INTO BEAST_SKIN VALUES (3, 'Odd texture')",
		"INSERT INTO BEAST_SKIN VALUES (4, 'Molts regularly')",
		"INSERT INTO BEAST_SKIN VALUES (5, 'Harmful to touch')",
		"INSERT INTO BEAST_SKIN VALUES (6, 'Wet or slimy')",

		// BEAST_WEAPON (d6)
		"INSERT INTO BEAST_WEAPON VALUES (1, 'Teeth or mandibles')",
		"INSERT INTO BEAST_WEAPON VALUES (2, 'Claws')",
		"INSERT INTO BEAST_WEAPON VALUES (3, 'Poison')",
		"INSERT INTO BEAST_WEAPON VALUES (4, 'Harmful discharge')",
		"INSERT INTO BEAST_WEAPON VALUES (5, 'Pincers')",
		"INSERT INTO BEAST_WEAPON VALUES (6, 'Horns')",

		// BEAST_SIZE (d6)
		"INSERT INTO BEAST_SIZE VALUES (1, 'Cat-sized')",
		"INSERT INTO BEAST_SIZE VALUES (2, 'Wolf-sized')",
		"INSERT INTO BEAST_SIZE VALUES (3, 'Calf-sized')",
		"INSERT INTO BEAST_SIZE VALUES (4, 'Bull-sized')",
		"INSERT INTO BEAST_SIZE VALUES (5, 'Hippo-sized')",
		"INSERT INTO BEAST_SIZE VALUES (6, 'Elephant-sized')",

		// BEAST_BEHAVIOR: composite id = type*100 + roll
		// Predator (type=1)
		"INSERT INTO BEAST_BEHAVIOR VALUES (101, 'Predator', 'Hunts in kin-group packs')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (102, 'Predator', 'Favors ambush attacks')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (103, 'Predator', 'Cripples prey and waits for death')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (104, 'Predator', 'Pack supports alpha-beast attack')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (105, 'Predator', 'Lures other prey with fake distress')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (106, 'Predator', 'Hunts as a lone, powerful hunter')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (107, 'Predator', 'Only hunts during mating frenzy')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (108, 'Predator', 'Mindlessly attacks humans')",
		// Prey (type=2)
		"INSERT INTO BEAST_BEHAVIOR VALUES (201, 'Prey', 'Moves in vigilant herds')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (202, 'Prey', 'Exists in small family groups')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (203, 'Prey', 'They all team up on a single foe')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (204, 'Prey', 'They go berserk when near death')",
		`INSERT INTO BEAST_BEHAVIOR VALUES (205, 'Prey', 'They''re violent in certain seasons')`,
		`INSERT INTO BEAST_BEHAVIOR VALUES (206, 'Prey', 'They''re vicious if threatened')`,
		"INSERT INTO BEAST_BEHAVIOR VALUES (207, 'Prey', 'Symbiotic creature protects them')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (208, 'Prey', 'Breeds at tremendous rates')",
		// Scavenger (type=3)
		"INSERT INTO BEAST_BEHAVIOR VALUES (301, 'Scavenger', 'Never attacks unwounded prey')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (302, 'Scavenger', 'Uses other beasts as fighters')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (303, 'Scavenger', 'Always flees if significantly hurt')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (304, 'Scavenger', 'Poisons prey, waits for it to die')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (305, 'Scavenger', 'Disguises itself as prey')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (306, 'Scavenger', 'Remarkably stealthy')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (307, 'Scavenger', 'Shames predators to abandon prey')",
		"INSERT INTO BEAST_BEHAVIOR VALUES (308, 'Scavenger', 'Steals prey from weaker predator')",

		// BEAST_DISCHARGE (d8)
		"INSERT INTO BEAST_DISCHARGE VALUES (1, 'Acidic secretion, squirts on a hit')",
		"INSERT INTO BEAST_DISCHARGE VALUES (2, 'Toxic spittle or cloud')",
		"INSERT INTO BEAST_DISCHARGE VALUES (3, 'Super-heated or super-chilled spew')",
		"INSERT INTO BEAST_DISCHARGE VALUES (4, 'Sonic shriek or ultrasonic noise')",
		"INSERT INTO BEAST_DISCHARGE VALUES (5, 'Natural laser or plasma discharge')",
		"INSERT INTO BEAST_DISCHARGE VALUES (6, 'Nausea-toxin venting')",
		"INSERT INTO BEAST_DISCHARGE VALUES (7, 'Equipment-melting corrosive')",
		"INSERT INTO BEAST_DISCHARGE VALUES (8, 'Explosive pellets or chemical catalysts')",

		// BEAST_POISON (d6)
		"INSERT INTO BEAST_POISON VALUES (1, 'Death', 'Instant', '1d6 rounds')",
		"INSERT INTO BEAST_POISON VALUES (2, 'Paralysis', '1 round', '1 minute')",
		"INSERT INTO BEAST_POISON VALUES (3, '1d4 damage per onset interval', '1d6 rounds', '10 minutes')",
		"INSERT INTO BEAST_POISON VALUES (4, 'Convulsions', '1 minute', '1 hour')",
		"INSERT INTO BEAST_POISON VALUES (5, 'Blindness', '1d6 minutes', '1d6 hours')",
		"INSERT INTO BEAST_POISON VALUES (6, 'Hallucinations', '1 hour', '1d6 days')",
	}
	return execAll(db, stmts)
}
