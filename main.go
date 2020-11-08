package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	htgotts "github.com/hegedustibor/htgo-tts"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	RANDOM_JOKE = "http://official-joke-api.appspot.com/random_joke"
	TEN_RANDOM_JOKE = "http://official-joke-api.appspot.com/random_ten"
	GRAP_BY_TYPE = "http://official-joke-api.appspot.com/jokes/{type}/ten"
)

func getjokes(URL string)  {
	resp,err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}

func getByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	content := vars["type"]
	fmt.Fprint(w, content)
}

func goSpeech(URL string)  {
	resp,err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}else {
		data, _ := ioutil.ReadAll(resp.Body)
		dataStr := string(data)
		fmt.Println(dataStr)

		var content map[string]interface{}
		err := json.Unmarshal(data, &content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		sendContent := fmt.Sprintf("%v", content["setup"])
		sendPunchline := fmt.Sprintf("%v", content["punchline"])
		//direct(sendContent)
		goAloud(sendContent,sendPunchline)
	}
}
func goAloud(text string, punch string)  {
	res,res2 := regexRem(text,punch)
	speak(res,res2)
}

func regexRem(text string, punch string) (a string, b string) {
	reg,err := regexp.Compile("[^a-zA-z0-9]+")

	a = reg.ReplaceAllString(text, "")
	if err != nil {
		fmt.Println(err.Error())
	}

	b = reg.ReplaceAllString(punch, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func speak(text string, punch string) {
	speech := htgotts.Speech{Folder: "contentResult", Language: "en"}
	speech.Speak(text)

	speech2 := htgotts.Speech{Folder: "punchResult", Language: "en"}
	speech2.Speak(punch)
}

func direct(text string) {
	reg, err := regexp.Compile("[^a-zA-z0-9]+")
	if err != nil {
		fmt.Println(err.Error())
	}
	a := reg.ReplaceAllString(text, "")
	speech := htgotts.Speech{Folder: "result", Language: "en"}
	_ = speech.Speak(a)

}

func main() {
	getjokes(RANDOM_JOKE)
	goSpeech(RANDOM_JOKE)

	getjokes(TEN_RANDOM_JOKE)


	router := mux.NewRouter()
	router.HandleFunc(GRAP_BY_TYPE,getByType).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}