// Package vreme simulates an array of train stations with weather measurements and 
// creates a report. 
//
// The package implements a map of train stations specified by label that contains its name, location and measurements.
// By calling the IzpisPorocila function, we can create a report based on the measurements for each station.
// 
// Example usage:
//  
//  vreme.DodajPostajo("LJ01", "Ljubljana Rudnik", "Ljubljana")
//  vreme.DodajMeritev("LJ01", 11.1)
//  vreme.DodajMeritev("LJ01", 13.2)
//  vreme.DodajMeritev("LJ01", 9.0)
// 
//  vreme.IzpisMeritev()
//  vreme.IzpisPorocila()
//
// The package defines these variables by default:
//   - minimum temperature: -50°C
//   - maximum temperature: 50°C
//   - minimum number of measurements for report: 3
// These values can be edited with the NastaviParametre function.
package vreme

import (
    "fmt"
)
// Postaja represents a single station.
type Postaja struct {
    ime     string
    kraj    string
    meritve []float64
    
}

var minTemp float64 = -50
var maxTemp float64 = 50
var stMeritev int = 3

var postaje = make(map[string]Postaja)

// IzpisPorocila() returns a string report based on measurements for each station.
// If a station has less than the minimum number of measurements, it will say "premalo meritev". 
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

// IzpisMeritev() returns a strign containing information of all the stations.
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
    
    if !ok || len(postaja.meritve) < stMeritev {
        return 0,0,0,false
    }
    
    var avg float64
    max := minTemp
    min := maxTemp
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

// With DodajMeritev(label, temperature) we can add a new temperature to the station
// specified by label.
// The function will return an error in case the temperature is higher/lower than the maximum/minimum temperature.
// It will also return an error if the given label doesn't exist.
func DodajMeritev(oznaka string, temperatura float64) error {
    temp, ok := postaje[oznaka]
    
    if !ok {
        err := fmt.Errorf("Ta postaja ne obstaja!")
        return err
    }
    
    if temperatura < minTemp || temperatura > maxTemp {
        err := fmt.Errorf("Podana temperatura ni veljavna!")
        return err
    }
    
    temp.meritve = append(temp.meritve, temperatura)
    postaje[oznaka] = temp

    return nil
}

// DodajPostajo(label, name, location) adds a new station.
// The function will return an error if any of the parameters are empty or 
// the given label already exists.
func DodajPostajo(oznaka string, ime string, kraj string) error {
    _, ok := postaje[oznaka]
    
    if ime == "" {
        err := fmt.Errorf("Ime je prazno")
        return err
    } else if oznaka == "" {
        err := fmt.Errorf("Oznaka je prazna")
        return err
    } else if kraj == "" {
        err := fmt.Errorf("Kraj je prazen")
        return err
    } else if ok {
        err := fmt.Errorf("Oznaka že obstaja")
        return err
    }
    
    var p Postaja
    p.ime = ime
    p.kraj = kraj
    p.meritve = []float64{}
    postaje[oznaka] = p
    return nil
}

// VrniOznake() returns a slice of all the labels
func VrniOznake() []string {
    oznake := []string{}
    for oznaka, _ := range postaje {
        oznake = append(oznake, oznaka)
    }
    
    return oznake
}

// NastaviParametre(minTemp, maxTemp, avgMeasurements) sets the parameters
// for minimum temperature, maximum temperature and number of measurements needed to
// create a report respectively.
func NastaviParametre(min float64, max float64, avg int) {
    minTemp = min
    maxTemp = max
    stMeritev = avg
}
