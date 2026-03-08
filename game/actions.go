package game

import "fmt"

// ExecuteAction processes a player action, logging messages directly to s.Log.
// Returns non-nil when the action paused for a player choice or triggered combat;
// the caller must handle the returned result (store choices or enter combat phase).
// Returns nil when the action completed fully (day already advanced).
func ExecuteAction(s *GameState, a Action) *EventResult {
	switch a {
	case ActionRest:
		msgs := DoRest(s)
		for _, msg := range msgs {
			s.AddLog(msg)
		}
		AdvanceDay(s)

	case ActionSeekNews:
		return doSeekNews(s)

	case ActionSeekFollowers:
		return doSeekFollowers(s)

	case ActionBuyFood:
		msgs := BuyFood(s, 10)
		for _, msg := range msgs {
			s.AddLog(msg)
		}
		AdvanceDay(s)

	case ActionSeekAudience:
		return doSeekAudience(s)

	case ActionSubmitOffering:
		return doSubmitOffering(s)

	case ActionSearchRuins:
		return doSearchRuins(s)

	case ActionSearchCache:
		doSearchCache(s)
		AdvanceDay(s)

	case ActionUseItem:
		return doUseItem(s)

	case ActionHunt:
		_, msgs := HuntForFood(s)
		for _, msg := range msgs {
			s.AddLog(msg)
		}
		AdvanceDay(s)

	case ActionBuyRaft:
		const raftCost = 15
		if s.Gold < raftCost {
			s.AddLog(fmt.Sprintf("You cannot afford a raft (costs %d gold, you have %d).", raftCost, s.Gold))
		} else {
			s.Gold -= raftCost
			s.Prince.AddPossession(PossRaft)
			s.AddLog(fmt.Sprintf("You purchase a raft for %d gold.", raftCost))
		}
		AdvanceDay(s)
	}
	return nil
}

// dispatchActionEvent fires an event for a non-travel action. It logs the
// narrative opener and event messages. If the result needs player input
// (choices or combat) it sets AdvanceDayOnChoice and returns it so the UI can
// pause; otherwise it applies the result and advances the day immediately.
func dispatchActionEvent(s *GameState, opener string, result EventResult) *EventResult {
	s.AddLog(opener)
	for _, msg := range result.Messages {
		s.AddLog(msg)
	}
	if result.CombatTriggered || len(result.Choices) > 0 {
		result.AdvanceDayOnChoice = true
		return &result
	}
	for _, msg := range applyEventResult(s, result) {
		s.AddLog(msg)
	}
	AdvanceDay(s)
	return nil
}

// doSeekNews implements r209: roll 2d6, +1 if W&W >= 5.
func doSeekNews(s *GameState) *EventResult {
	hex := GetHex(s.CurrentHex)
	if hex == nil || !hex.IsSettlement() {
		s.AddLog("No settlement here to seek news.")
		return nil
	}
	s.AddLog("You spend the day seeking news in the taverns and markets...")
	roll := Roll2d6()
	if s.EffectiveWitWiles() >= 5 {
		roll++
	}
	if roll > 12 {
		roll = 12
	}
	var result EventResult
	switch roll {
	case 2:
		result = EventResult{Messages: []string{"You hear nothing of note. The locals have little to say."}}
	case 3:
		result = EventResult{Messages: []string{"You discover a thieves' den — a quick raid nets 50 gold!"}, GoldChange: 50}
	case 4:
		result = EventResult{Messages: []string{"Rumours of secret rites at the nearest temple. Your next offering there is blessed (+1 to roll)."},
			Note: "Temple offering bonus at nearby temple"}
	case 5:
		result = EventResult{Messages: []string{"You feel at home here. The locals warm to you (+1 to future news/hire rolls in this hex)."},
			Note: fmt.Sprintf("Feel at home in %s — +1 to news/hire rolls", hex.Name)}
	case 6:
		result = DispatchEvent(s, "e129", EventContext{}) // large caravan
	case 7:
		result = EventResult{Messages: []string{"You discover an inn with cheaper rates. Lodging and food cost half in this hex today."},
			Note: fmt.Sprintf("Cheap lodgings found in %s", hex.Name)}
	case 8:
		stolen := s.Gold / 2
		result = EventResult{Messages: []string{fmt.Sprintf("A cutpurse relieves you of %d gold!", stolen)}, GoldChange: -stolen}
	case 9:
		result = DispatchEvent(s, "e050", EventContext{}) // constabulary
	case 10:
		result = DispatchEvent(s, "e016", EventContext{}) // local magician
	case 11:
		result = EventResult{Messages: []string{"A thieves' guild contact approaches — for 20 gold you learn of a cache nearby."},
			Note: "Thieves' guild tip: cache rumoured nearby"}
	default: // 12+
		if s.Gold >= 10 {
			result = DispatchEvent(s, "e143", EventContext{}) // secret informant
		} else {
			result = EventResult{Messages: []string{"A secret informant offers news — but wants 10 gold you don't have."}}
		}
	}
	return dispatchActionEvent(s, "", result)
}

// doSeekFollowers implements r210: roll 2d6.
func doSeekFollowers(s *GameState) *EventResult {
	hex := GetHex(s.CurrentHex)
	if hex == nil || !hex.IsSettlement() {
		s.AddLog("No settlement here to hire followers.")
		return nil
	}
	s.AddLog("You spend the day seeking followers in the local market...")
	roll := Roll2d6()
	var result EventResult
	switch roll {
	case 2: // Freeman joins free
		f := Character{Name: "Freeman", Type: TypeGeneric, CombatSkill: 3, MaxEndurance: 4, Morale: 4}
		s.AddFollower(f)
		result = EventResult{Messages: []string{"A freeman offers his sword, asking nothing in return. (CS 3, E 4)"}}
	case 3: // Lancer with horse
		f := Character{Name: "Lancer", Type: TypeMercenary, CombatSkill: 5, MaxEndurance: 5, DailyWage: 3, HasMount: true, MountType: MountHorse, Morale: 5}
		s.AddFollower(f)
		result = EventResult{Messages: []string{"A lancer with a horse joins for 3 gold/day. (CS 5, E 5)"}}
	case 4: // Mercenaries
		count := Roll1d6()%2 + 1
		for i := 0; i < count; i++ {
			name := fmt.Sprintf("Mercenary %d", i+1)
			f := Character{Name: name, Type: TypeMercenary, CombatSkill: 4, MaxEndurance: 4, DailyWage: 2, Morale: 4}
			s.AddFollower(f)
		}
		result = EventResult{Messages: []string{fmt.Sprintf("%d mercenary/mercenaries available at 2 gold/day each. (CS 4, E 4)", count)}}
	case 5: // Horse dealer
		if s.Gold >= 10 {
			s.Gold -= 10
			s.Prince.HasMount = true
			s.Prince.MountType = MountHorse
			result = EventResult{Messages: []string{"A horse dealer sells you a mount for 10 gold."}}
		} else {
			result = EventResult{Messages: []string{"A horse dealer offers mounts at 10 gold each — you can't afford one."}}
		}
	case 6: // Local guide
		f := Character{Name: "Local Guide", Type: TypeGuide, CombatSkill: 2, MaxEndurance: 3, DailyWage: 2, IsGuide: true, Morale: 4}
		s.AddFollower(f)
		result = EventResult{Messages: []string{"A local guide joins for 2 gold/day. (CS 2, E 3, Guide)"}}
	case 7: // Henchmen
		count := Roll1d6()
		for i := 0; i < count; i++ {
			name := fmt.Sprintf("Henchman %d", i+1)
			f := Character{Name: name, Type: TypeGeneric, CombatSkill: 3, MaxEndurance: 3, DailyWage: 1, Morale: 3}
			s.AddFollower(f)
		}
		result = EventResult{Messages: []string{fmt.Sprintf("%d henchman/henchmen available at 1 gold/day each. (CS 3, E 3)", count)}}
	case 8: // Slave market
		result = DispatchEvent(s, "e163", EventContext{})
	case 9: // Nothing — also check news with -1
		result = EventResult{Messages: []string{"No followers available today. You pick up a rumour instead..."}}
		newsRoll := Roll2d6() - 1
		if newsRoll < 2 {
			newsRoll = 2
		}
		s.AddLog(fmt.Sprintf("(News roll: %d)", newsRoll))
	case 10: // Honest horse dealer
		if s.Gold >= 7 {
			s.Gold -= 7
			s.Prince.HasMount = true
			s.Prince.MountType = MountHorse
			result = EventResult{Messages: []string{"An honest horse dealer sells you a fine mount for only 7 gold."}}
		} else {
			result = EventResult{Messages: []string{"An honest horse dealer offers mounts at 7 gold — you can't afford one."}}
		}
	case 11: // Runaway joins free
		f := Character{Name: "Runaway", Type: TypeGeneric, CombatSkill: 1, MaxEndurance: 3, IsEscapee: true, Morale: 3}
		s.AddFollower(f)
		result = EventResult{Messages: []string{"A runaway youth joins your party, asking nothing. (CS 1, E 3) — will flee if you enter a settlement."}}
	default: // 12 — porters + guide
		count := Roll1d6()
		for i := 0; i < count; i++ {
			name := fmt.Sprintf("Porter %d", i+1)
			f := Character{Name: name, Type: TypeGeneric, CombatSkill: 1, MaxEndurance: 2, DailyWage: 1, Morale: 3}
			s.AddFollower(f)
		}
		guide := Character{Name: "Guide", Type: TypeGuide, CombatSkill: 2, MaxEndurance: 3, DailyWage: 2, IsGuide: true, Morale: 4}
		s.AddFollower(guide)
		result = EventResult{Messages: []string{fmt.Sprintf("%d porter(s) at 1 gold/day + a guide at 2 gold/day join you.", count)}}
	}
	return dispatchActionEvent(s, "", result)
}

// doSeekAudience implements r211 per-location 2d6 tables.
func doSeekAudience(s *GameState) *EventResult {
	hex := GetHex(s.CurrentHex)
	if hex == nil || (hex.Structure != StructCastle && hex.Structure != StructKeep) {
		s.AddLog("There is no court here to seek an audience with.")
		return nil
	}
	if s.AudienceBarred[s.CurrentHex] > s.Day {
		s.AddLog(fmt.Sprintf("You are barred from this court until day %d.", s.AudienceBarred[s.CurrentHex]))
		return nil
	}
	s.AddLog(fmt.Sprintf("You request an audience at %s...", hex.Name))

	roll := Roll2d6()
	var eventID EventID

	switch s.CurrentHex {
	case NewHexID(12, 13): // Hulora Castle — Baron of Huldra (r211)
		switch {
		case roll <= 2:
			s.AudienceBarred[s.CurrentHex] = s.Day + 70 // permanently barred
			result := EventResult{Messages: []string{"You have permanently offended the Baron. Audience denied forever."}}
			return dispatchActionEvent(s, "", result)
		case roll == 3:
			eventID = "e154" // Baron's daughter
		case roll == 4:
			eventID = "e149" // learn court manners
		case roll == 5:
			eventID = "e158" // hostile guards
		case roll <= 7:
			s.AudienceBarred[s.CurrentHex] = s.Day + 1
			result := EventResult{Messages: []string{"Audience refused today. You may try again tomorrow."}}
			return dispatchActionEvent(s, "", result)
		case roll == 8:
			eventID = "e153" // Master of Household
		case roll == 9:
			eventID = "e148" // seneschal bribe
		case roll <= 11:
			eventID = "e150" // pay respects to Baron
		default:
			eventID = "e151" // find favour
		}
	case NewHexID(5, 23): // Adrogat Castle — Count Drogat
		switch {
		case roll <= 2:
			eventID = "e061" // next victim
		case roll == 3:
			eventID = "e062" // captain of guard
		case roll == 4:
			eventID = "e154" // count's daughter
		case roll == 5:
			eventID = "e153" // Master of Household
		case roll == 6:
			eventID = "e158" // hostile guards
		case roll <= 8:
			eventID = "e148" // seneschal bribe
		case roll == 9:
			eventID = "e149" // court manners
		case roll == 10:
			eventID = "e151" // find favour
		default:
			eventID = "e161" // audience granted
		}
	case NewHexID(20, 23): // Aeravir Castle — Lady Aeravir
		switch {
		case roll <= 2:
			eventID = "e060" // arrested
		case roll == 3:
			eventID = "e159" // purify yourself
		case roll == 4:
			eventID = "e158" // hostile guards
		case roll == 5:
			eventID = "e149" // court manners
		case roll == 6:
			eventID = "e153" // Master of Household
		case roll <= 8:
			s.AudienceBarred[s.CurrentHex] = s.Day + 1
			result := EventResult{Messages: []string{"Audience refused. You may try again tomorrow."}}
			return dispatchActionEvent(s, "", result)
		case roll == 9:
			eventID = "e148" // seneschal bribe
		case roll == 11:
			eventID = "e154" // Lady's daughter
		default:
			eventID = "e160" // audience granted
		}
	default: // Generic castle / keep
		switch {
		case roll <= 4:
			eventID = "e158" // hostile guards
		case roll == 5:
			eventID = "e153" // Master of Household
		case roll <= 7:
			s.AudienceBarred[s.CurrentHex] = s.Day + 1
			result := EventResult{Messages: []string{"Audience refused. You may try again tomorrow."}}
			return dispatchActionEvent(s, "", result)
		case roll <= 9:
			eventID = "e150" // pay respects
		case roll == 10:
			eventID = "e151" // find favour
		default:
			eventID = "e152" // noble ally
		}
	}

	result := DispatchEvent(s, eventID, EventContext{})
	return dispatchActionEvent(s, "", result)
}

func doSubmitOffering(s *GameState) *EventResult {
	hex := GetHex(s.CurrentHex)
	isTemple := hex != nil && (hex.Structure == StructTemple || s.GetHexFlags(s.CurrentHex).HiddenTemple)
	if !isTemple {
		s.AddLog("There is no temple here.")
		return nil
	}
	const offeringCost = 10
	if s.Gold < offeringCost {
		s.AddLog("You cannot afford to make an offering.")
		return nil
	}
	s.Gold -= offeringCost
	result := DispatchEvent(s, LookupOfferingEvent(Roll1d6()), EventContext{})
	name := "the hidden temple"
	if hex.Structure == StructTemple && hex.Name != "" {
		name = hex.Name
	}
	return dispatchActionEvent(s, fmt.Sprintf("You submit a %d gold offering at %s...", offeringCost, name), result)
}

func doSearchRuins(s *GameState) *EventResult {
	hex := GetHex(s.CurrentHex)
	if hex == nil || hex.Structure != StructRuins {
		s.AddLog("There are no ruins to search here.")
		return nil
	}
	flags := s.GetHexFlags(s.CurrentHex)
	if flags.Searched {
		s.AddLog("You have already thoroughly searched these ruins.")
		return nil
	}
	roll := Roll2d6()
	// Lantern improves the roll (better odds of finding treasure)
	if s.Prince.HasPossession(PossLantern) && Roll1d6() >= 4 {
		if roll < 12 {
			roll++
		}
	}
	result := DispatchEvent(s, LookupRuinsEvent(roll), EventContext{})
	flags.Searched = true
	return dispatchActionEvent(s, fmt.Sprintf("You search the ruins of %s...", hex.Name), result)
}

func doSearchCache(s *GameState) {
	flags := s.GetHexFlags(s.CurrentHex)
	if !flags.CacheHidden || flags.CacheFound {
		s.AddLog("There is no cache to find here.")
		return
	}
	for i, cache := range s.Caches {
		if cache.Location == s.CurrentHex && !cache.Found {
			// r214: 1-4 = found intact, 5 = can't find (try again later), 6 = looted
			roll := Roll1d6()
			switch {
			case roll <= 4:
				s.Caches[i].Found = true
				flags.CacheFound = true
				s.Gold += cache.Gold
				s.AddLog(fmt.Sprintf("You recover your hidden cache: %d gold!", cache.Gold))
			case roll == 5:
				s.AddLog("The landmarks have shifted — you cannot find the cache. You may try again.")
			default: // 6
				s.Caches[i].Found = true
				flags.CacheFound = true
				s.AddLog("You find the hiding spot, but the cache has been looted!")
			}
			return
		}
	}
	s.AddLog("The cache seems to have been found by others.")
}

func doUseItem(s *GameState) *EventResult {
	// Build a list of usable items the prince currently holds
	var usable []PossessionType
	for _, p := range []PossessionType{PossHealingPotion, PossPoisonAntidote} {
		if s.Prince.HasPossession(p) {
			usable = append(usable, p)
		}
	}
	if len(usable) == 0 {
		s.AddLog("You have no usable items.")
		return nil
	}

	choices := make([]string, len(usable))
	for i, p := range usable {
		choices[i] = PossessionName(p)
	}

	handler := func(gs *GameState, choice int) EventResult {
		if choice < 0 || choice >= len(usable) {
			return EventResult{Messages: []string{"Nothing happens."}}
		}
		item := usable[choice]
		gs.Prince.RemovePossession(item)
		switch item {
		case PossHealingPotion:
			healed := Roll1d6()
			if healed > gs.Prince.Wounds {
				healed = gs.Prince.Wounds
			}
			gs.Prince.Wounds -= healed
			return EventResult{Messages: []string{fmt.Sprintf("You drink the Healing Potion and recover %d wound(s).", healed)}}
		case PossPoisonAntidote:
			gs.Prince.PoisonWounds = 0
			return EventResult{Messages: []string{"You drink the Poison Antidote. The poison is purged from your body!"}}
		}
		return EventResult{Messages: []string{"Nothing happens."}}
	}

	result := EventResult{
		Messages:           []string{"You reach into your pack..."},
		Choices:            choices,
		ChoiceHandler:      handler,
		AdvanceDayOnChoice: true,
	}
	s.AddLog("Choose an item to use:")
	return &result
}

// HideCacheHere creates a cache at the current location
func HideCacheHere(s *GameState, gold int) []string {
	flags := s.GetHexFlags(s.CurrentHex)
	if flags.CacheHidden {
		return []string{"There is already a cache hidden here."}
	}
	if s.Gold < gold {
		gold = s.Gold
	}
	if gold <= 0 {
		return []string{"You have no gold to hide."}
	}
	s.Gold -= gold
	s.Caches = append(s.Caches, Cache{Location: s.CurrentHex, Gold: gold})
	flags.CacheHidden = true
	return []string{fmt.Sprintf("You bury %d gold for safekeeping.", gold)}
}
