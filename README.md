# vreme
Package vreme simulates an array of train stations with weather measurements and 
creates a report. 

The package implements a map of train stations specified by label that contains its name, location and measurements.
By calling the IzpisPorocila function, we can create a report based on the measurements for each station.
 
Example usage:
```go
 vreme.DodajPostajo("LJ01", "Ljubljana Rudnik", "Ljubljana")
 vreme.DodajMeritev("LJ01", 11.1)
 vreme.DodajMeritev("LJ01", 13.2)
 vreme.DodajMeritev("LJ01", 9.0)
 
 vreme.IzpisMeritev()
 vreme.IzpisPorocila()
```

The package defines these variables by default:
  - minimum temperature: -50°C
  - maximum temperature: 50°C
  - minimum number of measurements for report: 3
  
These values can be edited with the NastaviParametre function.
