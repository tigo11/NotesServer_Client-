package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"client_hw2/models/dto"
	"net/http"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Println("Пожалуйста, выберите действие:\n1. Добавить новую заметку\n2. Получить заметку по ID\n3. Обновить заметку по ID\n4. Удалить заметку по ID\n5. Получить все заметки\n6. Выйти")
		var action string
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		action, _ = reader.ReadString('\n')
		action = strings.TrimSpace(action)
		switch action {
		case "1":
			noteAdd()
		case "2":
			NoteGet()
		case "3":
			noteUpdate()
		case "4":
			noteDelete()
		case "5":
			notesGetAll()
		case "6":
			return
		default:
			fmt.Println("Wrong action")
		}
	}
}

func noteAdd() {
	note := dto.NewNote()
	fmt.Println("Заполните данные:")

	for note.Name == "" {
		fmt.Println("Имя:")
		reader := bufio.NewReader(os.Stdin)
		note.Name, _ = reader.ReadString('\n')
		note.Name = strings.TrimSpace(note.Name)
	}

	for note.LastName == "" {
		fmt.Println("Фамилия:")
		reader := bufio.NewReader(os.Stdin)
		note.LastName, _ = reader.ReadString('\n')
		note.LastName = strings.TrimSpace(note.LastName)
	}

	for note.Note == "" {
		fmt.Println("Заметка:")
		reader := bufio.NewReader(os.Stdin)
		note.Note, _ = reader.ReadString('\n')
		note.Note = strings.TrimSpace(note.Note)
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error json.Marshal():", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)
}

func NoteGet() {
	note := dto.NewNote()

	fmt.Println("Укажите ID заметки, которую Вы хотите увидеть:")
	fmt.Scanln(&note.ID)

	if note.ID < 1 {
		fmt.Println("Error: ID Must be valid")
		return
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)
}

func noteUpdate() {
	note := dto.NewNote()

	for note.ID < 1 {
		fmt.Println("Укажите ID заметки, которую Вы хотите изменить:")
		fmt.Scanln(&note.ID)
	}

	fmt.Println("Заполните данные:")

	fmt.Println("Имя:")
	reader := bufio.NewReader(os.Stdin)
	note.Name, _ = reader.ReadString('\n')
	note.Name = strings.TrimSpace(note.Name)

	fmt.Println("Фамилия:")
	reader = bufio.NewReader(os.Stdin)
	note.LastName, _ = reader.ReadString('\n')
	note.LastName = strings.TrimSpace(note.LastName)

	fmt.Println("Заметка:")
	reader = bufio.NewReader(os.Stdin)
	note.Note, _ = reader.ReadString('\n')
	note.Note = strings.TrimSpace(note.Note)

	if note.Name == "" && note.LastName == "" && note.Note == "" {
		fmt.Println("Error: all fields are empty")
		return
	}

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/update", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Response Body:", err)
		return
	}

	ResponseHandler(body)
}

func noteDelete() {
	note := dto.NewNote()
	fmt.Print("Введите ID заметки, которую Вы хотите удалить: ")
	fmt.Scanln(&note.ID)

	jsonData, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/delete", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)
}

func notesGetAll() {
	resp, err := http.Post("http://localhost:8080/get-all", "application/json", bytes.NewBuffer(nil))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	_, err = io.ReadFull(resp.Body, body)
	if err != nil {
		fmt.Println("Error Reading Body:", err)
		return
	}

	ResponseHandler(body)

}

func ResponseHandler(body []byte) {
	resp := dto.Response{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Error in response:", err)
		return
	}

	if resp.Error != "" {
		fmt.Println(resp.Result)
		return
	}

	fmt.Println("")

	if resp.Data != nil {
		data := []dto.Note{}
		err = json.Unmarshal(resp.Data, &data)
		if err != nil {
			data := dto.Note{}
			err = nil
			err = json.Unmarshal(resp.Data, &data)
			if err != nil {
				fmt.Println("Error in response:", err)
				return
			}
			PrintNote(data)
			fmt.Println()
			fmt.Println()
			return
		}

		for _, note := range data {
			PrintNote(note)
		}
	}
	fmt.Println()
	fmt.Println()
}

func PrintNote(note dto.Note) {
	fmt.Printf("ID заметки:\nID: %d\nИмя: %s\nФамилия: %s\nЗаметка: %s\n\n", note.ID, note.Name, note.LastName, note.Note)
}