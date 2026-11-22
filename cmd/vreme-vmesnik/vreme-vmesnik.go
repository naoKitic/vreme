package main

import (
    "fmt"
    "time"
    "strconv"
    "sync"
    "log"
    "os"
    "context"
    "github.com/urfave/cli/v3"
    "github.com/rivo/tview"
    "github.com/gdamore/tcell/v2"
    "github.com/naoKitic/vreme/vreme"
)

var lock sync.Mutex

func main() {
    var minTemp float64 
    var maxTemp float64 
    var stMeritev int 
    
    cmd := &cli.Command{
        
        Flags: []cli.Flag{
            &cli.Float64Flag{
                Name:        "minTemp",
                Value:       -50,
                Usage:       "minimum temperature for a measurement to be valid",
                Destination: &minTemp,
            },
            &cli.Float64Flag{
                Name:        "maxTemp",
                Value:       50,
                Usage:       "maximum temperature for a measurement to be valid",
                Destination: &maxTemp,
            },
            &cli.IntFlag{
                Name:        "stMeritev",
                Value:       3,
                Usage:       "number of measurements needed to calculate statistics",
                Destination: &stMeritev,
            },
        },
        
        Action: func(context.Context, *cli.Command) error {
            vreme.NastaviParametre(minTemp, maxTemp, stMeritev)
            return vremeInterface()
        },
        
        
    }
    
    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
    
}


func vremeInterface() error {
    app := tview.NewApplication()
    pages := tview.NewPages()
    
    //***************************************************************
    //~ Text field for drawing success and failure messages
    
    status := tview.NewTextView().
            SetText("Premikate se lahko z gumbom Tab").
            SetDynamicColors(true)
    status.SetBorder(true).
            SetTitle(" Status ")
    
    //***************************************************************
    //~ Text field for drawing results from vreme package
    
    curr_view := 1
    
    textView := tview.NewTextView().
        SetDynamicColors(true)
        
    go func() {
        for {
            var s string
            var flag bool
            if curr_view == 1 {
                s, flag = vreme.IzpisMeritev()
                
                if !flag {
                    s = "Zaenkrat ni nobene postaje..."
                }
                
            } else {
                s, flag = vreme.IzpisPorocila()
                if !flag {
                    s = "Zaenkrat ni nobenih meritev..."
                }
            }
            
            lock.Lock()
            app.QueueUpdateDraw(func() {
                textView.Clear()
                fmt.Fprintf(textView, "%s", s)
            })
            lock.Unlock()
            time.Sleep(500 * time.Millisecond)
        }
    }()
    
    //****************************************************************
    //~ Button to set if it writes measurements or report
    
    set_view_Btn := tview.NewButton("Spremeni pogled").SetSelectedFunc(func() {
        if curr_view == 1 {
            curr_view = 2 
        } else {
            curr_view = 1
        }
    })
    
    
    //****************************************************************
    //~ Form for adding a new measurement to a station
    new_measurement_form := tview.NewForm()
    new_measurement_form.SetBorder(true).SetTitle(" Dodaj nov meritev ")
    new_measurement_form.AddDropDown("Oznaka", vreme.VrniOznake(), 0, nil)
    new_measurement_form.AddInputField("Meritev", "", 20, nil, nil);
    new_measurement_form.AddButton("Dodaj", func() {
        _, oznaka := new_measurement_form.GetFormItemByLabel("Oznaka").(*tview.DropDown).GetCurrentOption()
        meritevText := new_measurement_form.GetFormItemByLabel("Meritev").(*tview.InputField).GetText()
        meritev, err := strconv.ParseFloat(meritevText, 64)
        
        if err != nil {
            return
        }
        
        lock.Lock()
        err2 := vreme.DodajMeritev(oznaka, meritev)
        
        if err2 != nil {
            error_text := fmt.Sprintf("[red]Prišlo je do napake: %v[-]", err2.Error())
            status.SetText(error_text)
        } else {
            status.SetText("[green]Meritev je bila uspešno dodana![-]")
        }
        
        lock.Unlock()
        
        new_measurement_form.GetFormItemByLabel("Oznaka").(*tview.DropDown).SetCurrentOption(-1)
        new_measurement_form.GetFormItemByLabel("Meritev").(*tview.InputField).SetText("")
        
        pages.SwitchToPage("main")
    })
    new_measurement_form.AddButton("Preklici", func() {
        new_measurement_form.GetFormItemByLabel("Oznaka").(*tview.DropDown).SetCurrentOption(-1)
        new_measurement_form.GetFormItemByLabel("Meritev").(*tview.InputField).SetText("")
        
        pages.SwitchToPage("main")
    })
    new_measurement_Btn := tview.NewButton("Dodaj meritev").SetSelectedFunc(func() {
        pages.ShowPage("new_measurement_form")
        app.SetFocus(new_measurement_form)
    })
    
    //*****************************************************************
    //~ Form for adding a new station
    new_station_form := tview.NewForm()
    new_station_form.SetBorder(true).SetTitle(" Dodaj novo postajo ")
    new_station_form.AddInputField("Oznaka", "", 20, nil, nil)
    new_station_form.AddInputField("Ime", "", 20, nil, nil)
    new_station_form.AddInputField("Kraj", "", 20, nil, nil)
    
    new_station_form.AddButton("Dodaj", func(){
        oznaka := new_station_form.GetFormItemByLabel("Oznaka").(*tview.InputField).GetText()
        ime := new_station_form.GetFormItemByLabel("Ime").(*tview.InputField).GetText()
        kraj := new_station_form.GetFormItemByLabel("Kraj").(*tview.InputField).GetText()
        
        err := vreme.DodajPostajo(oznaka, ime, kraj)
        
        if err != nil {
            error_text := fmt.Sprintf("[red]Prišlo je do napake: %v[-]", err.Error())
            status.SetText(error_text)
        } else {
            status.SetText("[green]Uspešno ste dodali postajo![-]")
        }
        
        new_measurement_form.GetFormItemByLabel("Oznaka").(*tview.DropDown).SetOptions(vreme.VrniOznake(), nil);
        
        new_station_form.GetFormItemByLabel("Oznaka").(*tview.InputField).SetText("")
        new_station_form.GetFormItemByLabel("Ime").(*tview.InputField).SetText("")
        new_station_form.GetFormItemByLabel("Kraj").(*tview.InputField).SetText("")
        
        pages.SwitchToPage("main")
    })
    
    new_station_form.AddButton("Preklici", func(){
        new_station_form.GetFormItemByLabel("Oznaka").(*tview.InputField).SetText("")
        new_station_form.GetFormItemByLabel("Ime").(*tview.InputField).SetText("")
        new_station_form.GetFormItemByLabel("Kraj").(*tview.InputField).SetText("")
        
        pages.SwitchToPage("main")
    })
    
    new_station_Btn := tview.NewButton("Dodaj postajo").SetSelectedFunc(func() {
        pages.ShowPage("new_station_form")
        app.SetFocus(new_station_form)
    })
    
    //******************************************************************
    buttonRow := tview.NewFlex().
            SetDirection(tview.FlexColumn).
            AddItem(set_view_Btn, 0, 1, true).
            AddItem(new_station_Btn, 0, 1, true).
            AddItem(new_measurement_Btn, 0, 1, true)
            
    mainScreen := tview.NewFlex().
            SetDirection(tview.FlexRow).
            AddItem(status, 3, 0, false).
            AddItem(nil, 1, 0, false).
            AddItem(textView, 0, 1, false).
            AddItem(buttonRow, 3, 0, true)

    mainScreen.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyTab {
            switch app.GetFocus() {
                case set_view_Btn:
                    app.SetFocus(new_station_Btn)
                    return nil
                case new_station_Btn:
                    app.SetFocus(new_measurement_Btn)
                    return nil
                case new_measurement_Btn:
                    app.SetFocus(set_view_Btn)
                    return nil
            }
        }
        return event
    })
    
    stationForm := tview.NewFlex().
            SetDirection(tview.FlexRow).
            AddItem(new_station_form, 0, 1, true)
    
    measurementForm := tview.NewFlex().
            SetDirection(tview.FlexRow).
            AddItem(new_measurement_form, 0, 1, true)
            
    pages.AddPage("main", mainScreen, true, true)
    pages.AddPage("new_station_form", stationForm, true, false)
    pages.AddPage("new_measurement_form", measurementForm, true, false)
           
    if err := app.SetRoot(pages, true).Run(); err != nil {
        panic(err)
    }
    return nil
}
