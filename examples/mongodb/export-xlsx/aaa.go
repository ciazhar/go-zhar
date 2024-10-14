package export_xlsx

import (
	"context"
	"fmt"
	"log"

	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type YourStruct struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	OldName  string `json:"old_name"`
	MapsId   string `json:"maps_id"`
	PostCode string `json:"post_code"`
	Geo      string `json:"geo"`
	Kdc      string `json:"kdc"`
}

func main() {
	// Step 2: Read the XLSX file and extract the data
	file, err := xlsx.OpenFile("area.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	data := []YourStruct{}
	sheet := file.Sheets[0] // Replace 0 with the desired sheet index

	for _, row := range sheet.Rows {
		data = append(data, YourStruct{
			Id:       row.Cells[0].String(),
			ParentId: row.Cells[1].String(),
			Name:     row.Cells[2].String(),
			Level:    row.Cells[3].String(),
			OldName:  row.Cells[4].String(),
			MapsId:   row.Cells[5].String(),
			PostCode: row.Cells[6].String(),
			Geo:      row.Cells[7].String(),
			Kdc:      row.Cells[8].String(),
		})
	}

	// Step 3: Connect to the MongoDB database
	clientOptions := options.Client().ApplyURI("mongodb://182.253.119.10:27017/").SetAuth(options.Credential{
		Username: "admin-tanggap",
		Password: "4dm1n-t4n994p",
	})
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("spgdt_pekalongan") // Replace "your_database" with your database name
	collection := database.Collection("area")       // Replace "your_collection" with your collection name

	// Step 5: Insert the data into the MongoDB collection
	newValue := make([]interface{}, len(data))
	for i := range data {
		newValue[i] = data[i]
	}

	_, err = collection.InsertMany(context.TODO(), newValue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data inserted successfully.")
}
