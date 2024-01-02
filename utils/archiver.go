package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/StephanUllmann/hypermedia-templ/db"
	models "github.com/StephanUllmann/hypermedia-templ/model"
)

type Archiver struct {
	Status          string
	Progress        float64
	ArchiveFilePath string
}

var A Archiver

func (a *Archiver) createFile() {
	var contacts []models.Contact

	rows, err := db.DB.Query("SELECT * FROM contacts;")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		contact := models.Contact{}
		err = rows.Scan(&contact.ID, &contact.First, &contact.Last, &contact.Phone, &contact.Email)
		if err != nil {
			fmt.Println(err)
		}
		contacts = append(contacts, contact)
	}

	j, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("public/contacts.json", j, 0644)
	if err != nil {
		fmt.Println(err)
	}
	a.ArchiveFilePath = "/public/contacts.json"
}

func (a *Archiver) Init() {
	a.Status = "Waiting"
	a.Progress = 0.0
	a.ArchiveFilePath = "/public/contacts.json"
}

func (a *Archiver) Run() {
	a.Status = "Running"
	a.Progress = 0.0
	a.ArchiveFilePath = "/public/contacts.json"

	go a.createFile()

}

func (a *Archiver) GetStatus() string {
	if a.Progress == 0.0 && a.Status == "Waiting" {
		return "Waiting"
	}
	if a.Progress < 1.0 {
		return "Running"
	} else {
		a.Status = "Complete"
		return "Complete"
	}
}

func (a *Archiver) GetProgress() float64 {
	if a.Progress == 0.0 && a.Status == "Waiting" {
		return 0.0
	}
	if a.Progress < 1.0 {
		a.Progress += 0.175
		return a.Progress
	} else {
		a.Progress = 1.0
		return a.Progress
	}
}

func (a *Archiver) GetArchiveFilePath() string {
	return a.ArchiveFilePath
}

// Approach with goroutines

// import (
// 	// "models
// 	"time"
// )

// type Archiver struct {
// 	Status          string
// 	Progress        float64
// 	ArchiveFilePath string
// }

// func (a *Archiver) Run() {
// 	// Create a channel to communicate progress updates
// 	progressChan := make(chan float64)

// 	go func() {
// 		// Simulate a long-running job
// 		for i := 0; i <= 100; i++ {
// 			time.Sleep(100 * time.Millisecond) // simulate work

// 			// Send progress update to the channel
// 			progressChan <- float64(i) / 100.0
// 		}

// 		close(progressChan) // close the channel when done
// 	}()

// 	// Update the status and progress in another goroutine
// 	go func() {
// 		for progress := range progressChan {
// 			a.Status = "Running"
// 			a.Progress = progress

// 			if progress == 1 {
// 				a.Status = "Complete"
// 				a.ArchiveFilePath = "/public/test.json"
// 			}
// 			// fmt.Printf("From Archiver: %v, %v\n", a.GetStatus(), a.GetProgress())
// 		}
// 	}()
// }

// func (a *Archiver) GetStatus() string {
// 	return a.Status
// }

// func (a *Archiver) GetProgress() float64 {
// 	return a.Progress
// }

// func (a *Archiver) GetArchiveFilePath() string {
// 	return a.ArchiveFilePath
// }