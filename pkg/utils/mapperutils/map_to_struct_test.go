package mapperutils

import (
	"fmt"
	"testing"
	"time"
)

// --- Beispielstrukturen und Demonstrationscode ---

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type User struct {
	ID             int        `json:"id"`
	Date           *time.Time `json:"date"`
	Name           string     `json:"name"`
	PrimaryAddress Address    `json:"primary_address"` // Verschachteltes Struct
	WorkAddress    *Address   `json:"work_address"`    // Zeiger auf Struct
	Tags           []string   `json:"tags"`            // Slice von primitiven Typen
	Friends        []User     `json:"friends"`         // Slice von Structs
	Colleagues     []*User    `json:"colleagues"`      // Slice von Zeigern auf Structs
}

func Test_MapToStruct(t *testing.T) {
	// 1. Quelldaten als map[string]interface{}
	data := map[string]interface{}{
		"id":   123,
		"name": "John Doe",
		"date": "2019-10-10 10:10:10",
		"primary_address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "Anytown",
		},
		"work_address": map[string]interface{}{
			"street": "456 Business Rd",
			"city":   "Workville",
		},
		"tags": []interface{}{"go", "developer", "engineer"},
		"friends": []interface{}{
			map[string]interface{}{
				"id":   201,
				"name": "Jane Smith",
				"primary_address": map[string]interface{}{
					"street": "1 Friend Ave",
					"city":   "Friendship",
				},
			},
		},
		"colleagues": []interface{}{
			map[string]interface{}{
				"id":   301,
				"name": "Peter Jones",
			},
			map[string]interface{}{
				"id":   302,
				"name": "Susan Williams",
			},
		},
	}

	// 2. Ziel-Struct-Instanz erstellen
	var user User

	// 3. Konvertierung durchführen
	err := MapToStructWithTag(data, &user, "json")
	if err != nil {
		fmt.Println("Fehler bei der Konvertierung:", err)
		return
	}

	// 4. Ergebnisse ausgeben und überprüfen
	fmt.Printf("Konvertiertes Struct: %+v\n\n", user)

	// Überprüfung der tiefen Kopie
	fmt.Println("--- Überprüfung der Details ---")
	fmt.Printf("ID: %d, Name: %s\n", user.ID, user.Name)
	fmt.Printf("Primäre Adresse: %s, %s\n", user.PrimaryAddress.Street, user.PrimaryAddress.City)
	if user.WorkAddress != nil {
		fmt.Printf("Arbeitsadresse: %s, %s\n", user.WorkAddress.Street, user.WorkAddress.City)
	}
	fmt.Printf("Tags: %v\n", user.Tags)
	if len(user.Friends) > 0 {
		fmt.Printf("Erster Freund: %+v\n", user.Friends[0])
	}
	if len(user.Colleagues) > 0 {
		fmt.Printf("Erster Kollege (Zeiger): %+v\n", *user.Colleagues[0])
		fmt.Printf("Zweiter Kollege (Zeiger): %+v\n", *user.Colleagues[1])
	}
}
