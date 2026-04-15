# SWN Revised Edition Sector Generator — Design Specification

This document specifies all modifications needed to port the Stars Without Number
sector generator (`swn-go`) to support the **Stars Without Number: Revised
Edition** ruleset. The original codebase lives in `../swn-go`; this revised
version targets `swn-revised-go`.

The spec is organized into two parts:

- **Part 1 — Required Changes**: modifications to existing generators where the
  Revised Edition changed the underlying tables or structure.
- **Part 2 — New Generators**: optional additions based on new content in the
  Revised Edition that did not exist in the original.

Source references use the format `file.go:NN` relative to `internal/`.

---

## Table of Contents

1. [Part 1: Required Changes](#part-1-required-changes)
   1. [1.1 World Tags — Expansion from 60 to 100](#11-world-tags)
   2. [1.2 World Population Table](#12-world-population-table)
   3. [1.3 World Temperature Table](#13-world-temperature-table)
   4. [1.4 World Tech Level Table](#14-world-tech-level-table)
   5. [1.5 NPC Generation — Complete Restructure](#15-npc-generation)
   6. [1.6 Alien Body Types — Reorder and Rename](#16-alien-body-types)
   7. [1.7 Alien Lenses — Reorder](#17-alien-lenses)
   8. [1.8 Alien Social Structure — Expand from d6 to d8](#18-alien-social-structure)
   9. [1.9 Alien Variations — Remove Per-Body-Type Sub-Tables](#19-alien-variations)
   10. [1.10 Seed Compatibility Note](#110-seed-compatibility)
2. [Part 2: New Generators](#part-2-new-generators)
   1. [2.1 System Points of Interest](#21-system-points-of-interest)
   2. [2.2 Other Worlds (Inter-World Relationships)](#22-other-worlds)
   3. [2.3 One-Roll Patrons](#23-one-roll-patrons)
   4. [2.4 Xenobestiary (Beast Generation)](#24-xenobestiary)
   5. [2.5 Faction Generation](#25-faction-generation)
   6. [2.6 Expanded Name Generators](#26-expanded-name-generators)
   7. [2.7 One-Roll Urban Encounters](#27-one-roll-urban-encounters)
   8. [2.8 One-Roll Wilderness Encounters](#28-one-roll-wilderness-encounters)
3. [Part 3: Cross-Cutting Concerns](#part-3-cross-cutting-concerns)
   1. [3.1 Database Migration Strategy](#31-database-migration-strategy)
   2. [3.2 API Changes](#32-api-changes)
   3. [3.3 TiddlyWiki Template Updates](#33-tiddlywiki-template-updates)
   4. [3.4 Frontend Updates](#34-frontend-updates)
   5. [3.5 Testing Strategy](#35-testing-strategy)

---

# Part 1: Required Changes

These changes reflect tables and mechanics that **differ** between the original
and Revised Edition. All existing generators that use these tables must be
updated.

---

## 1.1 World Tags

### What Changed

The original uses a 1d6+1d10 system producing 60 unique tags (IDs `101`–`610`).
The Revised Edition uses a d100 system producing 100 entries (1–99 are unique
tags; 100 means "roll twice and combine").

Of the original 60 tags, all but one are retained (some renamed). Approximately
40 new tags were added. "Ancient Ruins" was renamed to "Alien Ruins."

### New Tags (Not in Original)

```
Anarchists, Anthropomorphs, Battleground, Beastmasters, Cheap Life,
Cultural Power, Cybercommunists, Cyborgs, Cyclical Doom, Doomed World,
Dying Race, Fallen Hegemon, Former Warriors, Great Work, Hivemind,
Holy War, Immortals, Mandate Base, Maneaters, Megacorps, Mercenaries,
Night World, Nomads, Pleasure World, Post-Scarcity, Prison Planet,
Refugees, Revanchists, Revolutionaries, Rising Hegemon, Ritual Combat,
Robots, Shackled World, Societal Despair, Sole Supplier, Taboo Treasure,
Terraform Failure, Urbanized Surface, Utopia
```

### Database Changes

- **Drop and recreate `WORLD_TAG` table** with new ID scheme (integer 1–99).
- Each row retains the same schema: `id, name, desc, enemies, friends,
  complications, things, places, short`.
- All 99 tag entries must be populated from the Revised Edition PDF (p.132–159).
  Each tag entry in the PDF includes: a description paragraph, plus lists of
  Enemies, Friends, Complications, Things, and Places.
- The `short` column should contain an abbreviated name for table display, same
  as the original convention.

### Code Changes

**`generator/generator.go`**:

1. **`initTagIDs()`** (line 75–86): Replace the d6×d10 composite ID generation
   with a simple sequence of `"1"` through `"99"`.

   ```go
   // BEFORE: 6×10 grid producing IDs like "101", "205", "610"
   for d6 := 1; d6 < 7; d6++ {
       for d10 := 1; d10 < 11; d10++ {
           id := fmt.Sprintf("%d%02d", d6, d10)
           ...

   // AFTER: straight 1–99
   for i := 1; i <= 99; i++ {
       id := fmt.Sprintf("%d", i)
       ...
   ```

2. **`genWorldTag()`** (line 286–341): Update `totTags` from 60 to 99. The
   existing selection logic (random index, skip duplicates, wrap around) works
   unchanged — only the pool size changes.

3. **Handle "roll twice" (tag 100)**: This is a tabletop-only convenience (the
   generator already assigns two tags per world). No code change needed — simply
   do not include ID 100 in the database.

---

## 1.2 World Population Table

### What Changed

The population scale shifted upward and descriptions were reworded.

| Roll | Original | Revised |
|------|----------|---------|
| 2 | Failed colony | Failed colony *(unchanged)* |
| 3 | Outpost | Outpost *(unchanged)* |
| 4–5 | Tens of thousands of inhabitants | Fewer than a million inhabitants |
| 6–8 | Hundreds of thousands of inhabitants | Several million inhabitants |
| 9 | Millions of inhabitants | Hundreds of millions of inhabitants |
| 10 | Millions of inhabitants | Billions of inhabitants |
| 11 | Billions of inhabitants | Billions of inhabitants |
| 12 | Alien civilization | Alien inhabitants |

### Database Changes

Update `WORLD_POPULATION` rows for IDs 4–12: new `name`, `desc`, and `short`
values from Revised Edition (p.166–167).

### Code Changes

None — `genWorldComponent()` already looks up by roll ID.

---

## 1.3 World Temperature Table

### What Changed

The positions of "Variable" and "Cold"/"Warm" entries swapped.

| Roll | Original | Revised |
|------|----------|---------|
| 2 | Frozen | Frozen *(unchanged)* |
| 3 | Variable; Cold-to-Temperate | Cold |
| 4–5 | Cold | Variable cold with temperate places |
| 6–8 | Temperate | Temperate *(unchanged)* |
| 9 | Warm | Variable warm with temperate places |
| 10–11 | Variable; Temperate-to-Warm | Warm |
| 12 | Burning | Burning *(unchanged)* |

### Database Changes

Update `WORLD_TEMPERATURE` rows for IDs 3, 4, 5, 9, 10, 11: new `name`, `desc`,
and `short` values from Revised Edition (p.162–163).

### Code Changes

None.

---

## 1.4 World Tech Level Table

### What Changed

This is the most significant world-table change. The Revised Edition reorders
tech levels so that TL4 (postech) is the baseline at rolls 6–8, and TL3
(present-day Earth tech) appears at roll 9 as a rarer result. This reflects the
setting logic that most inhabited worlds recovered to postech, while
present-day-equivalent tech is actually uncommon.

| Roll | Original | Revised |
|------|----------|---------|
| 2 | TL0. Stone-age technology. | TL0, neolithic-level technology |
| 3 | TL1. Medieval technology. | TL1, medieval technology *(minor rewording)* |
| 4 | TL2. 19th-century technology. | TL2, early Industrial Age tech |
| 5 | TL3. 20th-century technology. | TL2, early Industrial Age tech |
| 6 | TL3. 20th-century technology. | TL4, modern postech |
| 7 | TL4. Baseline postech. | TL4, modern postech |
| 8 | TL4. Baseline postech. | TL4, modern postech |
| 9 | TL4. Baseline postech. | TL3, tech like that of present-day Earth |
| 10 | TL4. Baseline postech. | TL4+, postech with specialties |
| 11 | TL4 w/ specialities / pretech. | TL4+, postech with specialties |
| 12 | TL5. Pretech, pre-Silence tech. | TL5, pretech with surviving infrastructure |

### Database Changes

Update `WORLD_TECH_LEVEL` rows for IDs 2, 4–11: new `name`, `desc`, and `short`
values from Revised Edition (p.168–169).

### Code Changes

None.

---

## 1.5 NPC Generation

### What Changed

The NPC generation system was **completely restructured**. The original used six
tables with specific physical/personality attributes. The Revised Edition
replaces these with a "One-Roll NPC" system using the same die progression
(d4/d6/d8/d10/d12/d20) but entirely different categories.

| Die | Original Category | Revised Category |
|-----|-------------------|------------------|
| d4 | Gender (Male/Female) | Age (Young/Old for role, Young adult, Mature, Middle-aged/elderly) |
| d6 | Age (Young, Middle-aged, Old) | Background (Underclass, Laborers, Bourgeoisie, Elite, Minority, Offworlders) |
| d8 | Height (Very Short → Very Tall) | Role (Criminal, Menial, Unskilled labor, Skilled trade, Idea worker, Merchant, Official, Military) |
| d10 | Problem (10 specific problems) | Problem (10 new problems, different content) |
| d12 | Motivation (12 entries) | Desire (12 entries, different content) |
| d20 | Quirk (20 entries, some gender-specific) | Trait (20 single-word traits: Ambition, Courage, etc.) |

### Database Changes

Replace all six NPC tables with revised content:

| Original Table | New Table | Die | Entries |
|----------------|-----------|-----|---------|
| `NPC_GENDER` | `NPC_AGE` (repurpose) | d4 | 4 |
| `NPC_AGE` | `NPC_BACKGROUND` (new) | d6 | 6 |
| `NPC_HEIGHT` | `NPC_ROLE` (new) | d8 | 8 |
| `NPC_PROBLEMS` | `NPC_PROBLEM` (rename) | d10 | 10 |
| `NPC_MOTIVATION` | `NPC_DESIRE` (new) | d12 | 12 |
| `NPC_QUIRKS` | `NPC_TRAIT` (new) | d20 | 20 |

Drop: `NPC_GENDER`, `NPC_HEIGHT`, `NPC_MOTIVATION`, `NPC_QUIRKS`.
Rename/recreate: `NPC_AGE`, `NPC_PROBLEMS`.
Create: `NPC_BACKGROUND`, `NPC_ROLE`, `NPC_DESIRE`, `NPC_TRAIT`.

Schema for each new table:

```sql
CREATE TABLE NPC_AGE (id INTEGER PRIMARY KEY, age VARCHAR NOT NULL);
-- IDs 1–4, values from Revised p.244

CREATE TABLE NPC_BACKGROUND (id INTEGER PRIMARY KEY, background VARCHAR NOT NULL);
-- IDs 1–6, values from Revised p.244

CREATE TABLE NPC_ROLE (id INTEGER PRIMARY KEY, role VARCHAR NOT NULL);
-- IDs 1–8, values from Revised p.244

CREATE TABLE NPC_PROBLEM (id INTEGER PRIMARY KEY, problem VARCHAR NOT NULL);
-- IDs 1–10, values from Revised p.244

CREATE TABLE NPC_DESIRE (id INTEGER PRIMARY KEY, desire VARCHAR NOT NULL);
-- IDs 1–12, values from Revised p.244

CREATE TABLE NPC_TRAIT (id INTEGER PRIMARY KEY, trait VARCHAR NOT NULL);
-- IDs 1–20, values from Revised p.244
```

### Type Changes (`generator/types.go`)

```go
// BEFORE
type NPC struct {
    Name    string `json:"name"`
    Sort    string `json:"sort"`
    Gender  string `json:"gender"`
    Age     string `json:"age"`
    Height  string `json:"height"`
    Problem string `json:"problem"`
    Motive  string `json:"motive"`
    Quirk   string `json:"quirk"`
}

// AFTER
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
```

### Generator Changes (`generator/generator.go`)

Rewrite `GenNPCs()` (lines 344–396):

```go
func (g *Generator) GenNPCs(max int) []*NPC {
    if max == 0 {
        max = 100
    }
    npcs := make([]*NPC, 0, max)
    for i := 0; i < max; i++ {
        ageR  := dice.DInt(4)
        bgR   := dice.DInt(6)
        roleR := dice.DInt(8)
        probR := dice.DInt(10)
        desR  := dice.DInt(12)
        trR   := dice.DInt(20)

        age   := g.queryString("SELECT age FROM NPC_AGE WHERE id = ?", ageR)
        bg    := g.queryString("SELECT background FROM NPC_BACKGROUND WHERE id = ?", bgR)
        role  := g.queryString("SELECT role FROM NPC_ROLE WHERE id = ?", roleR)
        prob  := g.queryString("SELECT problem FROM NPC_PROBLEM WHERE id = ?", probR)
        des   := g.queryString("SELECT desire FROM NPC_DESIRE WHERE id = ?", desR)
        tr    := g.queryString("SELECT trait FROM NPC_TRAIT WHERE id = ?", trR)

        // Gender is no longer a rolled attribute. For name generation,
        // randomly pick male or female (50/50) to select from name lists.
        gender := "male"
        if dice.DInt(2) == 2 {
            gender = "female"
        }
        name := g.pickName(gender)

        nameParts := strings.SplitN(name, " ", 2)
        sortName := name
        if len(nameParts) == 2 {
            sortName = nameParts[1] + " " + nameParts[0]
        }

        npcs = append(npcs, &NPC{
            Name:       name,
            Sort:       sortName,
            Age:        age,
            Background: bg,
            Role:       role,
            Problem:    prob,
            Desire:     des,
            Trait:      tr,
        })
    }
    return npcs
}
```

Key decisions:
- Gender is no longer a visible attribute but is still used internally to select
  from gendered name lists.
- The gender-specific quirk handling (lines 367–375 in original) is removed
  entirely since Trait is a single-word value with no gendered variants.

### Wiki Changes (`wiki/wiki.go`)

**`writeNPCArticles()`** (lines 658–700): Update the NPC listing table and
article markup.

Listing table header changes from:
```
|!Name|!Gender|!Age|!Height|h
```
to:
```
|!Name|!Age|!Role|!Trait|h
```

Individual NPC article fields change from
`Gender/Age/Height/Quirk/Motive/Problem` to
`Age/Background/Role/Trait/Desire/Problem`.

Tag categories change from `Gender:X Age:X Height:X` to
`Age:X Background:X Role:X`.

---

## 1.6 Alien Body Types

### What Changed

The body types were reordered and "Humanlike" was renamed to "Mammalian."

| Roll | Original | Revised |
|------|----------|---------|
| 1 | Humanlike | Avian |
| 2 | Avian | Reptilian |
| 3 | Reptilian | Insectile |
| 4 | Insectile | Mammalian |
| 5 | Exotic | Exotic |
| 6 | Hybrid | Hybrid |

### Database Changes

Update `ALIEN_BODY_TYPE` rows 1–4 with new names and descriptions from Revised
Edition (p.203).

### Code Changes

**`generator/generator.go`** — Update the `alienVarTables` map (line 706–712):

```go
// BEFORE
var alienVarTables = map[string]string{
    "Humanlike":  "ALIEN_VAR_MAMMAL",
    "Avian":      "ALIEN_VAR_AVIAN",
    "Reptilian":  "ALIEN_VAR_REPTILE",
    "Insectile":  "ALIEN_VAR_INSECT",
    "Exotic":     "ALIEN_VAR_EXOTIC",
}

// AFTER
var alienVarTables = map[string]string{
    "Mammalian":  "ALIEN_VAR_MAMMAL",
    "Avian":      "ALIEN_VAR_AVIAN",
    "Reptilian":  "ALIEN_VAR_REPTILE",
    "Insectile":  "ALIEN_VAR_INSECT",
    "Exotic":     "ALIEN_VAR_EXOTIC",
}
```

Also in `genAlien()` (line 717–720): The hybrid roll logic uses `dice.DInt(5)`
for each component body type. The Revised Edition keeps 6 = Hybrid and 1–5 as
the non-hybrid types, so this logic is unchanged. However, because the order
changed, the existing body type rolls (1–5 mapping to specific types) will
automatically produce the correct results after the database is updated, since
the code looks up by ID.

---

## 1.7 Alien Lenses

### What Changed

Same 20 lenses, but Fear moved from position 6 to position 20, shifting entries
7–20 down by one position.

| ID | Original | Revised |
|----|----------|---------|
| 1 | Collectivity | Collectivity |
| 2 | Curiosity | Curiosity |
| 3 | Despair | Despair |
| 4 | Domination | Domination |
| 5 | Faith | Faith |
| 6 | **Fear** | **Gluttony** |
| 7 | Gluttony | Greed |
| 8 | Greed | Hate |
| 9 | Hate | Honor |
| 10 | Honor | Journeying |
| 11 | Journeying | Joy |
| 12 | Joy | Pacifism |
| 13 | Pacifism | Pride |
| 14 | Pride | Sagacity |
| 15 | Sagacity | Subtlety |
| 16 | Subtlety | Tradition |
| 17 | Tradition | Treachery |
| 18 | Treachery | Tribalism |
| 19 | Tribalism | Wrath |
| 20 | Wrath | **Fear** |

### Database Changes

Update `ALIEN_LENSES` rows 6–20 with new ID-to-name mappings. The descriptions
for each lens are unchanged — only the ordering shifted. The simplest approach:
delete all rows and reinsert with the revised ordering.

### Code Changes

None — lenses are looked up by roll ID from the database.

---

## 1.8 Alien Social Structure

### What Changed

The die changed from d6 to d8, and a new entry "Competitive" was added.

| Roll | Original | Revised |
|------|----------|---------|
| 1 | Democratic | Democratic |
| 2 | Monarchic | Monarchy |
| 3 | Tribal | Tribal |
| 4 | Oligarchic | Oligarchic |
| 5–6 | Multipolar | Multipolar |
| 7–8 | *(n/a)* | **Competitive** |

### Database Changes

- Update row 2: "Monarchic" → "Monarchy" (minor rename).
- Add rows 7 and 8 to `ALIEN_STRUCTURE` with name "Competitive" and an
  appropriate description. The Revised Edition leaves this intentionally
  open-ended for GM interpretation; use a description such as: "Individual
  members or small groups vie for dominance through direct competition, whether
  through combat, economic contests, intellectual challenges, or other metrics of
  personal excellence."

### Code Changes

**`generator/generator.go` — `genAlien()`** (line 757):

```go
// BEFORE
sroll = dice.DInt(6)

// AFTER
sroll = dice.DInt(8)
```

Also update the multipolar threshold logic (lines 757–759). In the original,
`sroll > 4` triggers "multipolar" (multiple structures). In the revised edition,
only `sroll == 5 || sroll == 6` should trigger multipolar, while `sroll >= 7`
is "Competitive" (single structure):

```go
// BEFORE
totSoc := 1
if sroll > 4 {
    totSoc = 1 + dice.DInt(3)
}
for i := 0; i < totSoc; i++ {
    roll := sroll
    if sroll > 4 {
        roll = dice.DInt(4)
    }

// AFTER
totSoc := 1
if sroll == 5 || sroll == 6 {
    totSoc = 1 + dice.DInt(3)
}
for i := 0; i < totSoc; i++ {
    roll := sroll
    if sroll == 5 || sroll == 6 {
        roll = dice.DInt(4)
    }
```

The token logic for deduplication (line 746–748) should also be updated. In the
original, `sroll == 5` and `sroll == 6` both map to the same "Multipolar"
concept, so the token normalizes them. In the revised edition, we should
normalize only 5→5 and 6→5 (still Multipolar), while 7 and 8 remain distinct
as Competitive:

```go
// BEFORE
stok := sroll
if sroll == 6 {
    stok = 5
}

// AFTER
stok := sroll
if sroll == 6 {
    stok = 5
} else if sroll == 8 {
    stok = 7
}
```

---

## 1.9 Alien Variations

### What Changed

The original edition had d20 variation sub-tables for each body type (stored in
`ALIEN_VAR_MAMMAL`, `ALIEN_VAR_AVIAN`, `ALIEN_VAR_REPTILE`, `ALIEN_VAR_INSECT`,
`ALIEN_VAR_EXOTIC`). The Revised Edition **removes these sub-tables entirely**,
replacing them with prose guidance for GMs.

### Decision Required

Two options:

**Option A — Keep the original variation tables.** The existing variations add
flavor that the prose guidance doesn't replicate procedurally. Since this is a
generator (not a tabletop reference), keeping the variation tables provides
richer automated output. The data is not contradicted by the revised rules — it
was simply moved to prose.

**Option B — Remove variation generation.** Strip out the `Vars` field from the
`Alien` struct and the associated variation-rolling code. This matches the
Revised Edition exactly but produces less detailed alien descriptions.

**Recommendation: Option A.** Keep the variation tables and the existing
generation logic. The data is compatible and enhances the generator's output.

---

## 1.10 Seed Compatibility Note

Changing die rolls, table sizes, and generation order means that **the same seed
will produce different sectors** in the revised generator compared to the
original. This is expected and unavoidable. The revised generator is a distinct
application, not a backward-compatible update.

If desired, a `?edition=original` query parameter could be added to preserve
access to the original generation logic, but this is not recommended due to the
complexity of maintaining two codepaths.

---

# Part 2: New Generators

These are entirely new generation systems from the Revised Edition that could
enrich the sector generator's output. They are listed in priority order —
highest-impact additions first.

---

## 2.1 System Points of Interest

### Description

The Revised Edition adds a system for generating non-planetary points of
interest within each star system: deep-space stations, asteroid bases, moon
bases, ancient orbital ruins, research bases, asteroid belts, gas giant mines,
and refueling stations.

This is the highest-value addition because it directly enriches the existing
sector map — each star system can now contain more than just planets.

### Data Model

```go
// PointOfInterest represents a generated system feature.
type PointOfInterest struct {
    Name       string `json:"name"`
    Type       string `json:"type"`       // "Deep-space station", etc.
    OccupiedBy string `json:"occupied_by"`
    Situation  string `json:"situation"`
    StarID     int    `json:"star_id"`
    Cell       string `json:"cell"`
}
```

### Database Tables

```sql
CREATE TABLE POI_TYPE (
    id INTEGER PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    desc CLOB NOT NULL
);
-- 8 rows (d8): Deep-space station, Asteroid base, Remote moon base,
-- Ancient orbital ruin, Research base, Asteroid belt, Gas giant mine,
-- Refueling station

CREATE TABLE POI_OCCUPANT (
    id INTEGER PRIMARY KEY,       -- composite: type_id * 100 + roll
    type_id INTEGER NOT NULL,     -- FK to POI_TYPE.id
    occupant VARCHAR(128) NOT NULL
);
-- 5 entries per type × 8 types = 40 rows
-- Each entry covers 2 d10 values (1-2, 3-4, 5-6, 7-8, 9-0)

CREATE TABLE POI_SITUATION (
    id INTEGER PRIMARY KEY,       -- composite: type_id * 100 + roll
    type_id INTEGER NOT NULL,
    situation VARCHAR(128) NOT NULL
);
-- 5 entries per type × 8 types = 40 rows
```

### Generation Logic

Add `GenPointsOfInterest()` to the generator. For each star system, roll 1d6:
on a 5–6, the system has a point of interest (tunable probability). If present:

1. Roll 1d8 for point type.
2. Roll 1d10 for occupant (map to 1-2/3-4/5-6/7-8/9-0 ranges).
3. Roll 1d10 for situation (same mapping).
4. Assign a name using existing cosmic name generator.

Default: generate 0–1 POIs per system. The GM version of the wiki should include
POI articles; the PC version should show only the type and name (not occupant or
situation).

### Sector Integration

Add `POIs []*PointOfInterest` to the `Sector` struct. System articles in the
wiki should list associated POIs alongside planets.

---

## 2.2 Other Worlds (Inter-World Relationships)

### Description

For star systems with 2+ worlds, the Revised Edition provides a d8 table to
generate the relationship between worlds: their origin, current relationship,
and primary contact point.

### Data Model

```go
// WorldRelationship describes how two worlds in the same system relate.
type WorldRelationship struct {
    Origin       string `json:"origin"`
    Relationship string `json:"relationship"`
    ContactPoint string `json:"contact_point"`
}
```

This should be attached to multi-world star systems. Add a `Relationship` field
to the `World` struct (populated only for the 2nd and 3rd worlds in a system,
describing their relationship to the primary world).

### Database Table

```sql
CREATE TABLE WORLD_RELATIONSHIP (
    id INTEGER PRIMARY KEY,   -- 1–8 (d8)
    origin CLOB NOT NULL,
    relationship CLOB NOT NULL,
    contact_point CLOB NOT NULL
);
```

### Generation Logic

In `GenSector()`, after generating worlds for a multi-world system, roll 1d8 for
each secondary world to determine its relationship to the primary world. Store
on the secondary world's struct.

### Wiki Integration

In system articles and planet articles for multi-world systems, add a
"Relationship to [Primary World]" section showing origin, current relationship,
and contact point. This is GM-only information (red text).

---

## 2.3 One-Roll Patrons

### Description

A new generator for quest-giving NPCs. Uses six dice
(d6/d4/d12/d8/d10/d20) to generate trustworthiness, eagerness, non-cash
rewards, job challenge, countervailing force, and complications.

### Data Model

```go
type Patron struct {
    Name            string `json:"name"`
    Trustworthiness string `json:"trustworthiness"`
    Eagerness       string `json:"eagerness"`
    Reward          string `json:"reward"`
    Challenge       string `json:"challenge"`
    Opposition      string `json:"opposition"`
    Complication    string `json:"complication"`
}
```

### Database Tables

```sql
CREATE TABLE PATRON_TRUST (id INTEGER PRIMARY KEY, trust VARCHAR NOT NULL);
-- 6 rows (d6)

CREATE TABLE PATRON_EAGERNESS (id INTEGER PRIMARY KEY, eagerness VARCHAR NOT NULL);
-- 4 rows (d4)

CREATE TABLE PATRON_REWARD (id INTEGER PRIMARY KEY, reward VARCHAR NOT NULL);
-- 12 rows (d12)

CREATE TABLE PATRON_CHALLENGE (id INTEGER PRIMARY KEY, challenge VARCHAR NOT NULL);
-- 8 rows (d8)

CREATE TABLE PATRON_OPPOSITION (id INTEGER PRIMARY KEY, opposition VARCHAR NOT NULL);
-- 10 rows (d10)

CREATE TABLE PATRON_COMPLICATION (id INTEGER PRIMARY KEY, complication VARCHAR NOT NULL);
-- 20 rows (d20)
```

### Generation Logic

Add `GenPatrons(max int) []*Patron` to the generator. Default count: 24.
Generate a name using existing name lists. Roll each die and look up the result.

### Sector Integration

Add `Patrons []*Patron` to `Sector`. Patrons are **GM-only** content — they
should not appear in the PC wiki at all. The GM wiki gets a "Patron Listing"
article with a summary table and individual patron articles.

---

## 2.4 Xenobestiary (Beast Generation)

### Description

A creature generation system for alien wildlife. Uses multiple dice to build
composite creatures with distinct features, weapons, and behaviors.

### Data Model

```go
type Beast struct {
    Name      string   `json:"name"`
    Features  []string `json:"features"`   // d10: Amphibian, Bird, etc.
    BodyPlan  string   `json:"body_plan"`  // d6
    LimbMod   string   `json:"limb_mod"`   // d6
    SkinMod   string   `json:"skin_mod"`   // d6
    Weapon    string   `json:"weapon"`     // d6
    Size      string   `json:"size"`       // d6
    BehavType string   `json:"behavior_type"` // Predator/Prey/Scavenger
    Behavior  string   `json:"behavior"`   // d8 within type
    Discharge string   `json:"discharge"`  // d8 (if weapon = "Harmful discharge")
    Poison    string   `json:"poison"`     // d6 (if weapon = "Poison")
}
```

### Database Tables

```sql
CREATE TABLE BEAST_FEATURE (id INTEGER PRIMARY KEY, feature VARCHAR NOT NULL);
-- 8 rows (d10: 1-8 unique, 9-10 = "roll twice")

CREATE TABLE BEAST_BODY_PLAN (id INTEGER PRIMARY KEY, plan VARCHAR NOT NULL);
-- 6 rows (d6: 1-5 unique, 6 = "roll twice")

CREATE TABLE BEAST_LIMB (id INTEGER PRIMARY KEY, limb VARCHAR NOT NULL);
-- 6 rows

CREATE TABLE BEAST_SKIN (id INTEGER PRIMARY KEY, skin VARCHAR NOT NULL);
-- 6 rows

CREATE TABLE BEAST_WEAPON (id INTEGER PRIMARY KEY, weapon VARCHAR NOT NULL);
-- 6 rows

CREATE TABLE BEAST_SIZE (id INTEGER PRIMARY KEY, size VARCHAR NOT NULL);
-- 6 rows

CREATE TABLE BEAST_BEHAVIOR (
    id INTEGER PRIMARY KEY,       -- composite: type * 100 + roll
    type VARCHAR NOT NULL,        -- "Predator", "Prey", "Scavenger"
    behavior VARCHAR NOT NULL
);
-- 8 entries × 3 types = 24 rows

CREATE TABLE BEAST_DISCHARGE (id INTEGER PRIMARY KEY, discharge VARCHAR NOT NULL);
-- 8 rows

CREATE TABLE BEAST_POISON (
    id INTEGER PRIMARY KEY,
    poison VARCHAR NOT NULL,
    onset VARCHAR NOT NULL,
    duration VARCHAR NOT NULL
);
-- 6 rows
```

### Generation Logic

Add `GenBeasts(max int) []*Beast` to the generator. Default count: 12.

1. Roll d10 for basic feature. On 9–10, roll twice and combine (hybrid).
2. Roll d6 for body plan. On 6, roll twice and combine.
3. Roll d6 each for limb mod, skin mod, weapon, size.
4. Roll d3 to pick behavior type (1=Predator, 2=Prey, 3=Scavenger), then
   d8 within that type.
5. If weapon = "Harmful discharge" (roll 4), roll d8 on discharge table.
6. If weapon = "Poison" (roll 3), roll d6 on poison table.
7. Generate name using conflux alien name generator.

### Sector Integration

Add `Beasts []*Beast` to `Sector`. Beasts are GM-only. The GM wiki gets a
"Xenobestiary" article listing all generated creatures with their attributes.

---

## 2.5 Faction Generation

### Description

The Revised Edition's most complex new system. Factions are organizations that
compete for control across the sector. Each has three ratings (Force, Cunning,
Wealth), hit points, assets, goals, and tags.

### Data Model

```go
type Faction struct {
    Name       string       `json:"name"`
    Homeworld  string       `json:"homeworld"`  // planet name
    Tags       []FactionTag `json:"tags"`
    Force      int          `json:"force"`
    Cunning    int          `json:"cunning"`
    Wealth     int          `json:"wealth"`
    HP         int          `json:"hp"`
    MaxHP      int          `json:"max_hp"`
    FacCreds   int          `json:"faccreds"`
    Goal       string       `json:"goal"`
    Assets     []string     `json:"assets"`
}

type FactionTag struct {
    Name string `json:"name"`
    Desc string `json:"desc"`
}
```

### Database Tables

```sql
CREATE TABLE FACTION_TAG (
    id INTEGER PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    desc CLOB NOT NULL,
    effect CLOB NOT NULL
);
-- 20 rows: Colonists, Deep Rooted, Eugenics Cult, Exchange Consulate,
-- Fanatical, Imperialists, Machiavellian, Mercenary Group,
-- Perimeter Agency, Pirates, Planetary Government, Plutocratic,
-- Preceptor Archive, Psychic Academy, Savage, Scavengers, Secretive,
-- Technical Expertise, Theocratic, Warlike

CREATE TABLE FACTION_GOAL (
    id INTEGER PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    desc CLOB NOT NULL
);
-- 11 rows: Military Conquest, Commercial Expansion, Intelligence Coup,
-- Planetary Seizure, Expand Influence, Blood the Enemy,
-- Peaceable Kingdom, Destroy the Foe, Inside Enemy Territory,
-- Invincible Valor, Wealth of Worlds
```

### Generation Logic

Add `GenFactions(max int) []*Faction` to the generator. Default count: 6–8
(approximately 1 faction per 3–4 star systems).

For each faction:
1. Assign a homeworld from the generated worlds (pick from populated worlds
   with population roll >= 6).
2. Roll Force, Cunning, Wealth ratings. Use a simple scheme:
   - One primary attribute at 3–5 (1d3+2)
   - Two secondary attributes at 1–3 (1d3)
3. Calculate HP = 4 + sum of rating XP values.
4. Select 1 faction tag from the d20 table.
5. Select 1 faction goal from the d11 table.
6. Select 1–3 starting assets appropriate to the primary attribute and rating.
7. Generate a name. For factions with the "Planetary Government" tag, use the
   homeworld name + "Government". Otherwise, use the political party name
   generator or a custom faction name scheme.

### Sector Integration

Add `Factions []*Faction` to `Sector`. Factions are **GM-only**. The GM wiki
gets a "Sector Factions" listing article and individual faction articles showing
all attributes, assets, goals, and tags. Faction homeworlds should be
cross-linked to their planet articles.

### Complexity Note

The full faction system includes detailed asset stat blocks, attack/defense
mechanics, and turn-by-turn resolution rules. For the generator, we only need to
**create** factions with initial state — the mechanical simulation is a tabletop
activity. The generator should produce a starting set of factions that the GM can
then run using the rulebook.

---

## 2.6 Expanded Name Generators

### Description

The Revised Edition provides d100 name tables for 10 cultures: Arabic, Chinese,
English, Greek, Indian, Japanese, Latin, Nigerian, Russian, and Spanish. Each
culture has 25 Male names, 25 Female names, 25 Surnames, and 50 Place names.

### Integration Approach

The current app uses pre-compiled name lists in `names/data.go` (1.1 MB of
`[][]string` data). The revised name data can either:

**Option A — Replace `names/data.go`** with culture-specific name sets. This
changes the character of generated names but aligns with the Revised Edition.

**Option B — Add culture-specific names as an additional source** alongside the
existing lists, used for specific purposes (e.g., place names for system names,
culture-tagged NPCs).

**Recommendation: Option B.** Keep the existing name lists for general NPC/star
names (they work well and are extensive), but add the Revised Edition's
culture-specific lists as a supplementary data source. This enables:
- Culture-tagged NPCs where a world's cultural background influences naming.
- More thematic place names for systems and sectors.

### Data Model

```go
// CultureNames holds name lists for a specific culture.
type CultureNames struct {
    Culture  string   `json:"culture"`
    Male     []string `json:"male"`
    Female   []string `json:"female"`
    Surnames []string `json:"surnames"`
    Places   []string `json:"places"`
}
```

### Storage

Add a new file `names/cultures.go` containing the 10 culture name sets, sourced
from Revised Edition p.239–243. Total data: 10 cultures × (25+25+25+50) =
1,250 name entries. This is small enough to embed directly in Go source.

### Generation Integration

No immediate integration required for the core generator. Expose via:
- A new API endpoint: `GET /api/names/{culture}` returning random names from
  that culture.
- Optional: when generating NPCs, allow a `culture` parameter that biases name
  selection toward a specific culture list.

---

## 2.7 One-Roll Urban Encounters

### Description

A generator for random urban encounter scenarios, using six dice
(d4/d12/d6/d8/d10/d20).

### Data Model

```go
type UrbanEncounter struct {
    Conflict   string `json:"conflict"`    // d4
    Antagonist string `json:"antagonist"`  // d12
    Venue      string `json:"venue"`       // d6
    PCInvolve  string `json:"pc_involve"`  // d8
    Event      string `json:"event"`       // d10
    Feature    string `json:"feature"`     // d20
}
```

### Database Tables

```sql
CREATE TABLE URBAN_CONFLICT (id INTEGER PRIMARY KEY, conflict VARCHAR NOT NULL);
CREATE TABLE URBAN_ANTAGONIST (id INTEGER PRIMARY KEY, antagonist VARCHAR NOT NULL);
CREATE TABLE URBAN_VENUE (id INTEGER PRIMARY KEY, venue VARCHAR NOT NULL);
CREATE TABLE URBAN_PC_INVOLVE (id INTEGER PRIMARY KEY, involve VARCHAR NOT NULL);
CREATE TABLE URBAN_EVENT (id INTEGER PRIMARY KEY, event VARCHAR NOT NULL);
CREATE TABLE URBAN_FEATURE (id INTEGER PRIMARY KEY, feature VARCHAR NOT NULL);
```

### Generation Logic

Add `GenUrbanEncounters(max int) []*UrbanEncounter`. Default count: 12.
Roll each die and look up the result.

### Sector Integration

Add `UrbanEncounters []*UrbanEncounter` to `Sector`. GM-only content.
The GM wiki gets an "Urban Encounters" reference article.

---

## 2.8 One-Roll Wilderness Encounters

### Description

A companion to urban encounters for wilderness scenarios, using six dice
(d4/d12/d6/d8/d10/d20).

### Data Model

```go
type WildEncounter struct {
    Range     string `json:"range"`       // d4
    Hostile   string `json:"hostile"`     // d12
    Weather   string `json:"weather"`     // d6
    Nature    string `json:"nature"`      // d8
    Friendly  string `json:"friendly"`    // d10
    Feature   string `json:"feature"`     // d20
}
```

### Database Tables

```sql
CREATE TABLE WILD_RANGE (id INTEGER PRIMARY KEY, range_desc VARCHAR NOT NULL);
CREATE TABLE WILD_HOSTILE (id INTEGER PRIMARY KEY, hostile VARCHAR NOT NULL);
CREATE TABLE WILD_WEATHER (id INTEGER PRIMARY KEY, weather VARCHAR NOT NULL);
CREATE TABLE WILD_NATURE (id INTEGER PRIMARY KEY, nature VARCHAR NOT NULL);
CREATE TABLE WILD_FRIENDLY (id INTEGER PRIMARY KEY, friendly VARCHAR NOT NULL);
CREATE TABLE WILD_FEATURE (id INTEGER PRIMARY KEY, feature VARCHAR NOT NULL);
```

### Generation Logic

Add `GenWildEncounters(max int) []*WildEncounter`. Default count: 12.

### Sector Integration

Add `WildEncounters []*WildEncounter` to `Sector`. GM-only content.
The GM wiki gets a "Wilderness Encounters" reference article.

---

# Part 3: Cross-Cutting Concerns

---

## 3.1 Database Migration Strategy

The simplest approach: **create a new SQLite database from scratch** rather than
migrating the existing one. The revised database has enough structural changes
(new tables, renamed tables, changed schemas, reordered data) that a clean build
is cleaner than incremental migration.

### Steps

1. Write a SQL schema file (`schema.sql`) defining all tables.
2. Write a data population script or SQL file (`data.sql`) with all INSERT
   statements, sourced from the Revised Edition PDF.
3. Build the new `swn.sqlite` using these files.
4. Embed the new database in the binary via `//go:embed`.

### Tables in the New Database

**Retained from original (with data updates):**
- `WORLD_ATMOSPHERE` — unchanged
- `WORLD_BIOSPHERE` — unchanged
- `WORLD_POPULATION` — updated data
- `WORLD_TECH_LEVEL` — updated data
- `WORLD_TEMPERATURE` — updated data
- `WORLD_TAG` — expanded from 60 to 99 rows, new ID scheme
- `ALIEN_BODY_TYPE` — reordered/renamed
- `ALIEN_LENSES` — reordered
- `ALIEN_STRUCTURE` — expanded (2 new rows)
- `ALIEN_VAR_MAMMAL`, `ALIEN_VAR_AVIAN`, `ALIEN_VAR_REPTILE`,
  `ALIEN_VAR_INSECT`, `ALIEN_VAR_EXOTIC` — retained as-is (see 1.9)
- `QUICK_ARCHITECTURE` — retained as-is (not in free Revised PDF; keep original)
- `QUICK_CORP_*` — retained as-is
- `QUICK_RELIGION_*` — retained as-is
- `QUICK_HERESY_*` — retained as-is
- `QUICK_POLITICAL_*` — retained as-is
- `QUICK_ROOM` — retained as-is

**Replaced:**
- `NPC_GENDER` → dropped
- `NPC_AGE` → recreated with new schema
- `NPC_HEIGHT` → dropped
- `NPC_PROBLEMS` → renamed to `NPC_PROBLEM`, new data
- `NPC_MOTIVATION` → dropped, replaced by `NPC_DESIRE`
- `NPC_QUIRKS` → dropped, replaced by `NPC_TRAIT`

**New tables:**
- `NPC_BACKGROUND`, `NPC_ROLE`, `NPC_DESIRE`, `NPC_TRAIT`
- `POI_TYPE`, `POI_OCCUPANT`, `POI_SITUATION`
- `WORLD_RELATIONSHIP`
- `PATRON_TRUST`, `PATRON_EAGERNESS`, `PATRON_REWARD`, `PATRON_CHALLENGE`,
  `PATRON_OPPOSITION`, `PATRON_COMPLICATION`
- `BEAST_*` (8 tables)
- `FACTION_TAG`, `FACTION_GOAL`
- `URBAN_*` (6 tables)
- `WILD_*` (6 tables)

---

## 3.2 API Changes

The existing API surface remains, with additions:

| Method | Path | Change |
|--------|------|--------|
| `GET` | `/api/seed` | Unchanged |
| `GET` | `/api/sector/{token}` | Response JSON expanded with new fields |
| `GET` | `/api/sector/{token}/map` | Unchanged |
| `POST` | `/api/sector/{token}/export` | Wiki output expanded |
| `GET` | `/api/names/{culture}` | **New** — returns random names from culture |

The `Sector` JSON response gains these fields:

```go
type Sector struct {
    Name      string              `json:"name"`
    Token     string              `json:"token"`
    Map       interface{}         `json:"map"`
    Stars     []*Star             `json:"stars,omitempty"`
    Worlds    []*World            `json:"worlds,omitempty"`
    NPCs      []*NPC              `json:"npcs,omitempty"`        // changed struct
    Corps     []*Corp             `json:"corps,omitempty"`
    Rels      []*Religion         `json:"rels,omitempty"`
    Pols      []*PoliticalParty   `json:"pols,omitempty"`
    Aliens    []*Alien            `json:"aliens,omitempty"`
    // New in Revised Edition:
    POIs      []*PointOfInterest  `json:"pois,omitempty"`
    Patrons   []*Patron           `json:"patrons,omitempty"`
    Beasts    []*Beast            `json:"beasts,omitempty"`
    Factions  []*Faction          `json:"factions,omitempty"`
    UrbanEnc  []*UrbanEncounter   `json:"urban_encounters,omitempty"`
    WildEnc   []*WildEncounter    `json:"wild_encounters,omitempty"`
}
```

The `handlers.go` `generateSector()` function must call the new generator
methods and populate the expanded Sector struct.

---

## 3.3 TiddlyWiki Template Updates

The TiddlyWiki template (`static/Includes/swn.html`) needs new placeholder
markers for the additional content sections. Add these placeholders:

```
<!--$$POI_INFO-->
<!--$$PATRON_INFO-->
<!--$$BEAST_INFO-->
<!--$$FACTION_INFO-->
<!--$$ENCOUNTER_INFO-->
```

The `wiki.go` file needs new `write*Articles()` functions:

- `writePOIArticles()` — listing and individual POI articles
- `writePatronArticles()` — GM-only patron listing and articles
- `writeBeastArticles()` — GM-only xenobestiary listing
- `writeFactionArticles()` — GM-only faction listing and articles
- `writeEncounterArticles()` — GM-only encounter reference tables

The bookmark section in the wiki should be updated to include links to
these new sections:

```go
pcBkmk := strings.Join([]string{
    "----",
    "[[Alien Races]]",
    "[[Political Groups]]",
    "[[Corporations|Corporation Listing]]",
    "[[Religions|Sector Religions]]",
}, "\n")

gmBkmk := pcBkmk + "\n" + strings.Join([]string{
    "[[NPCs|NPC Listing]]",
    "[[Patrons|Patron Listing]]",
    "[[Factions|Sector Factions]]",
    "[[Xenobestiary]]",
    "[[Points of Interest]]",
    "[[Urban Encounters]]",
    "[[Wilderness Encounters]]",
}, "\n")
```

The GM Info article (`writeGMInfo()`) should be updated to describe the new
content sections.

---

## 3.4 Frontend Updates

The `static/index.html` frontend is a jQuery-based single page app. Updates
needed:

1. **Sector display**: Add sections in the sector data view for new content
   types (POIs, Patrons, Beasts, Factions, Encounters). These can follow the
   existing pattern of collapsible sections.

2. **NPC display**: Update the NPC table columns from
   Gender/Age/Height/Problem/Motive/Quirk to Age/Background/Role/Problem/
   Desire/Trait.

3. **New content toggles**: Add checkboxes or controls to include/exclude
   optional content types from the export (e.g., "Include Factions",
   "Include Xenobestiary").

4. **Title/branding**: Update any references from "Stars Without Number" to
   "Stars Without Number: Revised Edition" to distinguish from the original app.

---

## 3.5 Testing Strategy

### Unit Tests

Each new or modified generator function should have unit tests verifying:

1. **Output count**: `GenPatrons(24)` returns exactly 24 patrons.
2. **No panics on all roll values**: Seed-sweep test that generates 100 sectors
   with different seeds, ensuring no nil dereferences or database lookup
   failures.
3. **Field population**: Every field in every generated struct is non-empty.
4. **Uniqueness**: Where the generator prevents duplicates (tags, names),
   verify no duplicates appear.
5. **Die range coverage**: For modified tables, verify that every valid roll ID
   (e.g., 2–12 for 2d6 tables) returns a valid result from the database.

### Integration Tests

1. **Full sector generation**: Generate a sector with all content types and
   verify the JSON is valid and all expected fields are present.
2. **Wiki export**: Generate a wiki export and verify the ZIP contains valid
   HTML files with all expected sections.
3. **Backward compatibility**: The original `swn-go` should continue to
   function independently — no shared state or dependencies between the two
   apps.

### Manual Verification

1. Start the dev server and generate several sectors via the web UI.
2. Download and open both GM and PC wiki exports in a browser.
3. Verify all new sections render correctly in TiddlyWiki.
4. Verify PC wiki does not contain any GM-only content.
5. Verify cross-links between factions and their homeworlds work.
6. Verify POIs appear in the correct system articles.

---

## Appendix A: File Change Summary

| File | Change Type | Description |
|------|-------------|-------------|
| `static/swn.sqlite` | Rebuild | New database with all schema/data changes |
| `internal/generator/types.go` | Modify | Update NPC struct; add POI, Patron, Beast, Faction, Encounter types |
| `internal/generator/generator.go` | Modify | Update tag init, NPC gen, alien gen; add new generator methods |
| `internal/handlers/handlers.go` | Modify | Call new generators; update Sector struct population |
| `internal/wiki/wiki.go` | Modify | Update NPC articles; add POI/Patron/Beast/Faction/Encounter articles |
| `internal/names/cultures.go` | New | Culture-specific name data (10 cultures) |
| `static/Includes/swn.html` | Modify | Add template placeholders for new content |
| `static/index.html` | Modify | Update NPC display; add new content sections |
| `main.go` | Minor | Update title/branding; no structural changes |

## Appendix B: Corporation/Religion/Heresy/Political/Architecture Tables

These tables were moved to the **Deluxe Edition** of the Revised rules and are
not present in the free Revised Edition PDF. The original tables remain valid
and are not contradicted by the Revised Edition.

**Decision: Retain all original tables unchanged.** If the Deluxe Edition PDF
becomes available, these tables should be compared and updated as a follow-up
task. The tables in question:

- `QUICK_CORP_NAME`, `QUICK_CORP_BUSINESS`, `QUICK_CORP_REPUTATION`
- `QUICK_RELIGION_ORIGIN`, `QUICK_RELIGION_LEADERSHIP`, `QUICK_RELIGION_EVOLUTION`
- `QUICK_HERESY_ATTITUDE`, `QUICK_HERESY_FOUNDER`, `QUICK_HERESY_HERESY`,
  `QUICK_HERESY_QUIRK`
- `QUICK_POLITICAL_NAME`, `QUICK_POLITICAL_ISSUES`, `QUICK_POLITICAL_LEADERSHIP`,
  `QUICK_POLITICAL_OUTSIDERS`, `QUICK_POLITICAL_POLICY`
- `QUICK_ARCHITECTURE`

## Appendix C: Suggested Implementation Order

1. **Copy `swn-go` codebase** into `swn-revised-go` as the starting point.
2. **Database rebuild** — create new `swn.sqlite` with all schema changes. This
   is the foundation everything else depends on.
3. **World table updates** (Population, Temperature, Tech Level) — smallest
   changes, good for verifying the database integration works.
4. **World tag expansion** — moderate code change, large data entry effort.
5. **NPC restructure** — touches types, generator, and wiki; do this as a
   cohesive unit.
6. **Alien updates** (body types, lenses, social structure) — moderate changes
   across generator and database.
7. **System Points of Interest** — highest-value new feature.
8. **Other Worlds relationships** — small addition that enriches existing
   multi-world systems.
9. **Patrons** — self-contained new generator.
10. **Xenobestiary** — self-contained new generator.
11. **Factions** — most complex new feature; do last.
12. **Encounter generators** — simple additions.
13. **Culture name lists** — low priority, additive.
14. **Frontend/wiki updates** — integrate incrementally as each generator is
    completed.
