package almaany

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

//All the harakats
var harakat = [...]string{"ِ", "ُ", "ٓ", "ٰ", "ْ", "ٌ", "ٍ", "ً", "ّ", "َ"}

const dbFile = "./dbFile/maani.db"

//Remove all tashkeels from a word
func removeTashkeel(word string) string {
	for _, harka := range harakat {
		word = strings.Replace(word, harka, "", -1)
	}
	return word

}

var db *sql.DB
var err error

//Open a database connection and create tables.
func InitDatabase() bool {
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Printf("%q: \n", err)
		return false
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS searchKeys ( word VARCHAR(40) PRIMARY KEY, terms TEXT);
CREATE TABLE IF NOT EXISTS MAANI (word VARCHAR(40) primary key, wordType VARCHAR(40),explanations TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return false
	}
	return true
}

//Save words to the database
func SaveWords(searched string, words []Manaa) bool {

	stmt, err := db.Prepare("INSERT INTO MAANI(word,wordType,explanations) values(?,?,?)")
	if err != nil {
		log.Printf("%q: \n", err)
		return false
	}

	var searchKeys []string
	for _, result := range words {
		word := result.Word
		wordType := result.WordType
		encodedJSON, _ := json.Marshal(result.Explanations)
		explanations := string(encodedJSON)
		_, err = stmt.Exec(word, wordType, explanations)
		if err != nil {
			log.Printf("%q: \n", err)
			return false
		}
		searchKeys = append(searchKeys, word)
	}
	stmt, err = db.Prepare("INSERT INTO searchKeys(word,terms) values(?,?)")
	if err != nil {
		log.Printf("%q: \n", err)
		return false
	}
	encodedJson, _ := json.Marshal(searchKeys)

	_, err = stmt.Exec(removeTashkeel(searched), string(encodedJson))
	if err != nil {
		log.Printf("%q: \n", err)
		return false
	}

	return true
}

//Return a list of words from a searched string
func GetSearchedWord(word string) []string {
	var foundWords []string
	stmt, err := db.Prepare("SELECT terms from searchKeys where word=?")
	if err != nil {
		log.Printf("%q: \n", err)
		return foundWords
	}
	rows, err := stmt.Query(removeTashkeel(word))
	if err != nil {
		log.Printf("%q: \n", err)
		return foundWords
	}
	defer rows.Close()
	rows.Next()
	var jsonAnswer string
	err = rows.Scan(&jsonAnswer)
	_ = json.Unmarshal([]byte(jsonAnswer), &foundWords)
	return foundWords

}

// Returns full Manna struct from an exactWord
func GetExplanation(exactWord string) Manaa {
	var explanation Manaa
	stmt, err := db.Prepare("SELECT * from MAANI where word=?")
	if err != nil {
		log.Printf("%q: \n", err)
		return explanation
	}
	rows, err := stmt.Query(exactWord)
	if err != nil {
		log.Printf("%q: \n", err)
		return explanation
	}
	defer rows.Close()
	rows.Next()
	var word string
	var wordType string
	var explanations string

	err = rows.Scan(&word, &wordType, &explanations)
	var explanationsArray []string
	_ = json.Unmarshal([]byte(explanations), &explanationsArray)
	explanation = Manaa{word, wordType, explanationsArray}
	return explanation
}
