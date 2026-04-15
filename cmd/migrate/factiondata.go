package main

import "database/sql"

// migrateFactionData inserts faction tag and goal data from the
// Stars Without Number Revised Edition rules.
func migrateFactionData(db *sql.DB) error {
	// Faction goals from Revised Edition p.217
	goals := []struct {
		id   int
		name string
		desc string
	}{
		{1, "Military Conquest", `Destroy a number of Force assets of rival factions equal to your faction's Force rating. Difficulty is 1/2 number of assets destroyed.`},
		{2, "Commercial Expansion", `Destroy a number of Wealth assets of rival factions equal to your faction's Wealth rating. Difficulty is 1/2 number of assets destroyed.`},
		{3, "Intelligence Coup", `Destroy a number of Cunning assets of rival factions equal to your faction's Cunning rating. Difficulty is 1/2 number of assets destroyed.`},
		{4, "Planetary Seizure", `Take control of a planet, becoming the legitimate planetary government. Difficulty equal to half the average of the current ruling faction's Force, Cunning, and Wealth ratings. If the planet somehow lacks any opposing faction to resist the seizure, it counts as Difficulty 1.`},
		{5, "Expand Influence", `Plant a Base of Influence on a new planet. Difficulty 1, +1 if the attempt is contested by a rival faction.`},
		{6, "Blood the Enemy", `Inflict a number of hit points of damage on enemy faction assets or bases equal to your faction's total Force, Cunning, and Wealth ratings. Difficulty 2.`},
		{7, "Peaceable Kingdom", `Don't take an Attack action for four turns. Difficulty 1.`},
		{8, "Destroy the Foe", `Destroy a rival faction. Difficulty equal to 1 plus the average of the faction's Force, Cunning, and Wealth ratings.`},
		{9, "Inside Enemy Territory", `Have a number of stealthed assets on worlds with other planetary governments equal to your Cunning score. Units that are already stealthed on worlds when this goal is adopted don't count. Difficulty 2.`},
		{10, "Invincible Valor", `Destroy a Force asset with a minimum purchase rating higher than your faction's Force rating. Thus, if your Force is 3, you need to destroy a unit that requires Force 4 or higher to purchase. Difficulty 2.`},
		{11, "Wealth of Worlds", `Spend FacCreds equal to four times your faction's Wealth rating on bribes and influence. This money is effectively lost, but the goal is then considered accomplished. The faction's Wealth rating must increase before this goal can be selected again. Difficulty 2.`},
	}

	// Faction tags from Revised Edition
	tags := []struct {
		id     int
		name   string
		desc   string
		effect string
	}{
		{1, "Colonists",
			`This faction is a fresh colony on an otherwise largely untouched planet. It is this brave band of pioneers that will tame the world's wild forces and bring forth a better life for those who come after.`,
			`This faction has all the benefits of the Planetary Government tag for its homeworld, as no other government exists on a fresh colony. The faction's homeworld is also treated as if it had at least tech level 4. Colonies with fewer than 100,000 citizens lack the necessary industrial infrastructure to build Spaceship-type assets.`},
		{2, "Deep Rooted",
			`This faction has been part of a world's life for time out of mind. Most natives can hardly imagine the world without this faction's presence, and the traditional prerogatives and dignities of the group are instinctively respected.`,
			`This faction can roll one additional d10 when defending against attacks on assets on their homeworld. If the faction ever changes homeworlds, this tag is lost.`},
		{3, "Eugenics Cult",
			`The forbidden maltech secrets of advanced human genetic manipulation are known to this faction, and they use them with gusto. Slave-engineered humanoids and "deathless" leadership are just two of the more common alterations these unstable scientists undertake.`,
			`Eugenics Cultists can buy the Gengineered Slaves asset; it's an asset requiring Force 1 with the statistics of 6 HP, 2 FacCred cost, tech level 4 required, with an Attack of Force vs. Force/1d6 damage and a Counterattack of 1d4 damage. Once per turn, the Eugenics Cult can roll an extra d10 on an attack or defense by a Gengineered Slaves asset, regardless of the stat being used. Gengineered Slaves can count as either a Military Unit or Special Forces, determined when the cult first creates a specific asset.`},
		{4, "Exchange Consulate",
			`This faction is either led through an Exchange Consulate or has close ties with that pacifistic society of bankers and diplomats. The sophisticated economic services they provide strengthen the faction.`,
			`When the faction successfully completes a "Peaceable Kingdom" Goal, they may roll 1d6; on a 4+, they gain a bonus experience point. Once per turn, the faction may roll an extra d10 when defending against a Wealth attack.`},
		{5, "Fanatical",
			`The members of this faction just don't know when to quit. No matter how overmatched, the members will keep fighting to the bitter end- and occasionally past it.`,
			`The faction always rerolls any dice that come up as 1. This zealousness leaves them open at times, however; they always lose ties during attacks.`},
		{6, "Imperialists",
			`This faction nurses wild dreams of controlling the sector, whether out of an impulse to bring their local culture and technology to less fortunate worlds or simple lust for dominion. They excel at defeating planetary defenses and standing armies.`,
			`This faction may roll an extra d10 for attacks made as part of a Seize Planet action.`},
		{7, "Machiavellian",
			`This faction's meat and drink is intrigue, its members delighting in every opportunity to scheme. It may be a secret cabal of hidden masters or the decadent court of a fallen stellar empire, but its membership has forgotten more of treachery than most others ever learn.`,
			`Once per turn, this faction can roll an additional d10 when making a Cunning attack.`},
		{8, "Mercenary Group",
			`The faction sells its services to the highest bidder, and is an extremely mobile organization. Vast amounts of men and material can be moved interstellar distances in just a few months.`,
			`All faction assets gain the following special ability: As an action, the asset may move itself to any world within one hex.`},
		{9, "Perimeter Agency",
			`This faction is or is closely tied to an Agency of the enigmatic Perimeter organization. Originally organized by the Terran Mandate to detect and contain maltech outbreaks until Mandate fleet resources could be dispatched, the Perimeter retains numerous ancient override codes for pretech security protocols.`,
			`Once per turn, the faction may roll an additional d10 when making an attack against an asset that requires tech level 5 to purchase. The faction may roll an extra die when making a test to detect Stealthed assets.`},
		{10, "Pirates",
			`This faction is a scourge of the spacelanes, driving up the cost of shipping and terrorizing merchant captains without pity. They steal and refit ships with vicious ingenuity, cobbling together space armadas out of the leavings of their prey.`,
			`Any movement of an asset onto a world that has a Base of Influence for this faction costs one extra FacCred, paid to this faction.`},
		{11, "Planetary Government",
			`This faction is the legitimate government of a planet. Rebel groups and rival factions may have assets on the planet, but control over the instruments of the state is firmly in this faction's hands. The faction may rule openly, or it may simply have an inexorable grasp on the existing authorities.`,
			`The faction's permission is required to buy or import those assets marked as needing government permission. This tag can be acquired multiple times, once for each planet the faction controls.`},
		{12, "Plutocratic",
			`This faction prizes wealth, and its membership strives constantly to expand and maintain personal fortunes. Perhaps it is a ruling council of oligarchs or a star-spanning trade cartel.`,
			`Once per turn, this faction can roll an additional d10 when making a Wealth attack.`},
		{13, "Preceptor Archive",
			`This faction is or has close ties to a Preceptor Archive, a place of learning operated by the learned Preceptors of the Great Archive. These Archives are peaceful institutions dedicated to the spread of practical knowledge and useful engineering to the wider cosmos. Their large numbers of educated personnel make advanced equipment more practical for deployment.`,
			`Purchasing an asset that requires tech level 4 or more costs one fewer FacCred than normal. The Preceptor Archive may also take the special action "Teach Planetary Population", costing 2 FacCreds and allowing them to roll 1d12 for one world. On a 12, the world's tech level permanently becomes 4 for the purposes and purchases of this faction.`},
		{14, "Psychic Academy",
			`Most significant factions are capable of employing psychics, but this faction can actually train their own. They excel at precise and focused application of the psionic disciplines, and can get far more out of their available psychic resources than other factions.`,
			`This faction can provide psionic mentor training to qualified psychics. Once per turn, this faction can also force a rival faction to reroll any one d10, whether or not they're involved in the roll.`},
		{15, "Savage",
			`Whether a proud tribe of neoprimitives struggling against the material limits of their world or a pack of degenerate tomb world cannibals, this faction is accustomed to surviving without the benefits of advanced technology and maximizing local resources.`,
			`Once per turn, this faction can roll an extra die when defending with an asset that requires tech level 0 to purchase.`},
		{16, "Scavengers",
			`This faction might live within the wreckage of a tomb world, salvage the castoffs of some decadent pleasure-world or ply the ruins of an orbital scrapyard. Whatever their territory, this faction knows how to find worth amid seemingly useless trash.`,
			`Whenever the faction destroys an asset or has one of their assets destroyed, they gain 1 FacCred.`},
		{17, "Secretive",
			`This faction is marked by elaborate protocols of secrecy and misdirection. It may be split up into numerous semi-autonomous cells, or consist largely of trained espionage professionals. Finding the assets of such a faction can often be more difficult than destroying them.`,
			`All assets purchased by this faction automatically begin Stealthed. See the list of Cunning assets for details on Stealth.`},
		{18, "Technical Expertise",
			`The faction is staffed by large numbers of expert engineers and scientists. They can turn even the most unpromising labor pool into superb technicians.`,
			`This faction treats all planets on which they have a Base of Influence as if they were at least tech level 4. They can build Starship-type assets on any world with at least ten thousand occupants.`},
		{19, "Theocratic",
			`The faction is fueled by the fierce certainty that God is with them- and with no one else. The tight and occasionally irrational obedience that pervades the organization makes it difficult to infiltrate or subvert effectively.`,
			`Once per turn, this faction can roll an extra d10 when defending against a Cunning attack.`},
		{20, "Warlike",
			`There are factions with a military orientation, and then there are factions that just really love killing things. Whether or not this faction has developed sophisticated military resources and techniques, the membership is zealous in battle and has a positive taste for fighting.`,
			`Once per turn, this faction can roll an additional d10 when making a Force attack.`},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, g := range goals {
		if _, err := tx.Exec("INSERT INTO FACTION_GOAL VALUES (?, ?, ?)", g.id, g.name, g.desc); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, t := range tags {
		if _, err := tx.Exec("INSERT INTO FACTION_TAG VALUES (?, ?, ?, ?)", t.id, t.name, t.desc, t.effect); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
