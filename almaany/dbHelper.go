package almaany

import (
	"database/sql"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
CREATE TABLE IF NOT EXISTS MAANI (word VARCHAR(40) primary key, wordType VARCHAR(40),explanations TEXT);
CREATE TABLE IF NOT EXISTS USERS(id integer primary key ,firstName VARCHAR(255),lastName VARCHAR(255),username VARCHAR(255))
`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return false
	}
	return true
}

//Adds a user
func AddUser(user *tgbotapi.User) {

	id := user.ID
	firstName := user.FirstName
	if len(firstName) > 200 {
		firstName = firstName[0:200]
	}
	lastName := user.LastName
	if len(lastName) > 200 {
		lastName = lastName[0:200]
	}
	username := user.UserName
	stmt, err := db.Prepare("INSERT INTO USERS(id,firstName,lastName,username) values(?,?,?,?)")
	if err != nil {
		log.Printf("%q: \n", err)
	}

	_, err = stmt.Exec(id, firstName, lastName, username)
	if err != nil {
		log.Printf("%q: \n", err)
	}

}

//Save words to the database
func SaveWords(searched string, words []Manaa) bool {

	stmt, err := db.Prepare("INSERT INTO MAANI(word,wordType,explanations) values(?,?,?)")
	if err != nil {
		log.Printf("%q: \n", err)
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
		}
		searchKeys = append(searchKeys, word)
	}
	stmt, err = db.Prepare("INSERT INTO searchKeys(word,terms) values(?,?)")
	if err != nil {
		log.Printf("%q: \n", err)
	}
	encodedJson, _ := json.Marshal(searchKeys)

	_, err = stmt.Exec(removeTashkeel(searched), string(encodedJson))
	if err != nil {
		log.Printf("%q: \n", err)
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
