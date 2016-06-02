package dice

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/itsabot/abot/shared/datatypes"
	"github.com/itsabot/abot/shared/nlp"
	"github.com/itsabot/abot/shared/plugin"
)

var p *dt.Plugin

func init() {

	rand.Seed(time.Now().UTC().UnixNano())

	// Abot should route messages to this plugin that contain any combination
	// of the below words. The stems of the words below are used, so you don't
	// need to include duplicates (e.g. there's no need to include both "stock"
	// and "stocks"). Everything will be lowercased as well, so there's no
	// difference between "ETF" and "etf".
	trigger := &nlp.StructuredInput{
		Commands: []string{"roll"},
		Objects:  []string{"die", "dice"},
	}

	// Tell Abot how this plugin will respond to new conversations and follow-up
	// requests.
	fns := &dt.PluginFns{Run: Run, FollowUp: FollowUp}

	// Create the plugin.
	var err error
	pluginPath := "github.com/njaremko/abot_roll_die"
	p, err = plugin.New(pluginPath, trigger, fns)
	if err != nil {
		log.Fatalln("building", err)
	}

	// Add vocab handlers to the plugin
	p.Vocab = dt.NewVocab(
		dt.VocabHandler{
			Fn:      findNumDieSides,
			Trigger: trigger,
		},
	)
}

// Abot calls Run the first time a user interacts with a plugin
func Run(in *dt.Msg) (string, error) {
	return FollowUp(in)
}

// Later add the ability for people to say "reroll" and such
func FollowUp(in *dt.Msg) (string, error) {
	return p.Vocab.HandleKeywords(in), nil
}

func findNumDieSides(in *dt.Msg) (resp string) {
	for _, obj := range in.Tokens {
		if sides, err := strconv.Atoi(obj); err == nil {
			return rollDie(sides)
		}
	}
	return rollDie(6)
}

func rollDie(sides int) string {
	var result int
	result = rand.Intn(sides + 1)
	if result == 0 {
		return rollDie(sides)
	} else {
		return "I rolled a " + strconv.Itoa(result) + "."
	}
}
