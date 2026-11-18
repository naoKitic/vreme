package vreme

import (
    "fmt"
    "math/rand"
    "math"
)

type Postaja struct {
    ime     string
    kraj    string
    meritve []float64
    
}

func IzpisPorocila(postaje map[string]Postaja) {
    for oznaka, postaja := range postaje {
        min, max, avg, ok := statistika(postaje, oznaka)
        
        if !ok {
            fmt.Printf("%s: premalo meritev\n", postaja.ime)
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
        
        fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", postaja.ime, min, max, avg, status)
    }
}

func IzpisMeritev(postaje map[string]Postaja) {
    fmt.Println("MERITVE:")
    for oznaka, postaja := range postaje {
        fmt.Printf("%s - %s (%s): ", oznaka, postaja.ime, postaja.kraj)
        fmt.Println(postaja.meritve)
    }
}

func statistika(postaje map[string]Postaja, oznaka string) (najnizja, najvisja, povprecna float64, ok bool) {
    postaja, ok := postaje[oznaka]
    
    if !ok || len(postaja.meritve) < 3 {
        return 0,0,0,false
    }
    
    var min float64
    var max float64
    var avg float64
    
    max = -51
    min = 51
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

func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura float64) {
    temp, ok := postaje[oznaka]
    
    if !ok {
        fmt.Println("Ta postaja ne obstaja!")
        return
    }
    
    if temperatura < -50 || temperatura > 50 {
        fmt.Println("Podana temperatura ni veljavna!")
        return
    }
    
    temp.meritve = append(temp.meritve, temperatura)
    postaje[oznaka] = temp

}

func DodajPostajo(oznaka string, ime string, kraj string) {
    var p Postaja
    p.ime = ime
    p.kraj = kraj
    p.meritve = []float64
    meritve[oznaka] = p
}

