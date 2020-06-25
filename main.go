package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/kr/pretty"
	"github.com/tealeg/xlsx"
	"googlemaps.github.io/maps"
	"log"
	"os"
	"strconv"
)

func main() {

	file, err := os.Open("api.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 100)
	count, err := file.Read(data)
	apikey := string(data[:count])

	app := app.New()

	w := app.NewWindow("Google Maps Searcher")
	querybox := widget.NewEntry()
	querybox.SetPlaceHolder("Enter Query Here...")
	thresholdbox := widget.NewEntry()
	thresholdbox.SetPlaceHolder("Enter Threshold Here...")
	filenamebox := widget.NewEntry()
	filenamebox.SetPlaceHolder("Enter Filename Here...")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("your api is: "+apikey),
		querybox,
		thresholdbox,
		filenamebox,
		widget.NewButton("Run", func() {
			query := querybox.Text
			threshold, _ := strconv.ParseFloat(thresholdbox.Text, 64)
			filename := filenamebox.Text
			println(query)
			println(threshold)
			println(filename)
			filteredResults := search(apikey, query, float32(threshold))
			save(filename,filteredResults)
			println("Finished with this query!")
			println("Results can be found in "+filename+".xlsx")
		}),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()

}

func search(apikey string, query string, threshold float32) [][]string {
	c, err := maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.TextSearchRequest{
		Query:     query,
		Location:  nil,
		Radius:    0,
		Language:  "",
		MinPrice:  "",
		MaxPrice:  "",
		OpenNow:   false,
		Type:      "",
		PageToken: "",
		Region:    "",
	}
	answer, _ := c.TextSearch(context.Background(), r)
	var filteredResults [][]string
	//pretty.Println(answer.Results[0].Rating)
	for _, e := range answer.Results {
		if e.Rating < threshold {
			//println(e.Name)
			details := &maps.PlaceDetailsRequest{
				PlaceID:      e.PlaceID,
				Language:     "",
				Fields:       nil,
				SessionToken: maps.PlaceAutocompleteSessionToken{},
				Region:       "",
			}
			f, _ := c.PlaceDetails(context.Background(), details)
			//pretty.Println(f)
			pretty.Println(f.Name, f.FormattedPhoneNumber, f.FormattedAddress, f.Rating, f.UserRatingsTotal)
			filteredResults = append(filteredResults, []string{f.Name, f.FormattedPhoneNumber, f.FormattedAddress, strconv.FormatFloat(float64(f.Rating),'f',-1,32), strconv.Itoa(f.UserRatingsTotal)})
		}
	}
	return filteredResults
}

func save(filename string, entries [][]string) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	if err != nil {
		fmt.Printf(err.Error())
	}
	for _,e := range entries{
		row = sheet.AddRow()
		for _,ee :=range e{
			cell = row.AddCell()
			cell.Value = ee
		}
	}
	err = file.Save(filename+".xlsx")
}
