package vreme

import (
    "fmt"
)

type Postaja struct {
    ime     string
    kraj    string
    meritve []float64
    
}

var MinTemp float64
var MaxTemp float64
var StMeritev int

var postaje = make(map[string]Postaja)

func IzpisPorocila() (string, bool) {
    var s string
    ima_postaje := false
    if len(postaje) > 0 {
        ima_postaje = true
    }
    for oznaka, postaja := range postaje {
        min, max, avg, ok := statistika(oznaka)
        
        if !ok {
            s += fmt.Sprintf("%s: premalo meritev\n", postaja.ime)
            continue
        }
        
        var status string
        
        switch {
            case avg >= 25:
                status = "Vroče"
            case avg >= 15 && avg < 25:
                status = "Toplo"
            case avg >= 5 && avg < 15:
                status = "Hladno"
            case avg < 5:
                status = "Mrzlo"
        }
        
        s += fmt.Sprintf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", postaja.ime, min, max, avg, status)
    }
    return s, ima_postaje
}

func IzpisMeritev() (string, bool)  {
    meritve := "MERITVE:\n"
    ima_meritve := false
    if len(postaje) > 0 {
        ima_meritve = true
    }
    for oznaka, postaja := range postaje {
        
        meritve += fmt.Sprintf("%s - %s (%s): %v\n", oznaka, postaja.ime, postaja.kraj, postaja.meritve)
    }
    return meritve, ima_meritve
}

func statistika(oznaka string) (najnizja, najvisja, povprecna float64, ok bool) {
    postaja, ok := postaje[oznaka]
    
    if !ok || len(postaja.meritve) < StMeritev {
        return 0,0,0,false
    }
    
    var avg float64
    max := MinTemp
    min := MaxTemp
    avg = 0
    
    for _, m := range postaja.meritve {  
        if m > max {
            max = m
        }
        if m < min {
            min = m
        }
        
        avg += m
    }
    
    avg = avg / float64(len(postaja.meritve))
    
    return min, max, avg, true
    
}

func DodajMeritev(oznaka string, temperatura float64) error {
    temp, ok := postaje[oznaka]
    
    if !ok {
        err := fmt.Errorf("Ta postaja ne obstaja!")
        return err
    }
    
    if temperatura < MinTemp || temperatura > MaxTemp {
        err := fmt.Errorf("Podana temperatura ni veljavna!")
        return err
    }
    
    temp.meritve = append(temp.meritve, temperatura)
    postaje[oznaka] = temp

    return nil
}

func DodajPostajo(oznaka string, ime string, kraj string) error {
    if ime == "" {
        err := fmt.Errorf("Ime je prazno")
        return err
    } else if oznaka == "" {
        err := fmt.Errorf("Oznaka je prazna")
        return err
    } else if kraj == "" {
        err := fmt.Errorf("Kraj je prazen")
        return err
    }
    
    var p Postaja
    p.ime = ime
    p.kraj = kraj
    p.meritve = []float64{}
    postaje[oznaka] = p
    
    return nil
}

func VrniOznake() []string {
    oznake := []string{}
    for oznaka, _ := range postaje {
        oznake = append(oznake, oznaka)
    }
    return oznake
}
