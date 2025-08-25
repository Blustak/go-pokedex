package pokemon

import (
	"math/rand"
    "fmt"
)

// < 40, then max chance to catch
// > 300, then min chance to catch
const minChanceToCatch = float64(0.15)
const maxChanceToCatch = float64(0.85)
const lowExpThreshhold = int(50)
const hiExpThreshhold = int(250)

type Pokemon struct {
    Name string `json:"name"`
    BaseExperience int `json:"base_experience"`
    Height int `json:"height"`
    Weight int `json:"weight"`
    Stats []struct{
        BaseStat int `json:"base_stat"`
        Stat struct{
            Name string `json:"name"`
        } `json:"stat"`
    } `json:"stats"`
    Types []struct{
        Type struct{
            Name string `json:"name"`
        } `json:"type"`
    } `json:"types"`
}

func (p Pokemon) TryCatch() bool {
    r := rand.Float64()
    var chance float64
    if p.BaseExperience <= lowExpThreshhold {
        chance = maxChanceToCatch
    }
    if p.BaseExperience >= hiExpThreshhold {
        chance = minChanceToCatch
    }
    if p.BaseExperience > lowExpThreshhold && p.BaseExperience < hiExpThreshhold {
        chance = maxChanceToCatch - (maxChanceToCatch - minChanceToCatch)*(float64(p.BaseExperience - lowExpThreshhold)/float64(hiExpThreshhold - lowExpThreshhold))
    }
    return r < chance
}

func (p Pokemon) Print() {
    fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", p.Name,p.Height,p.Weight)
    for _, s := range p.Stats {
        fmt.Printf("  - %s: %d\n", s.Stat.Name, s.BaseStat)
    }
    fmt.Println("Types:")
    for _, t := range p.Types {
        fmt.Printf("  - %s\n", t.Type.Name)
    }
}
