package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Project struct {
	Path  string `json:"project_path"`
	Title string `json:"project_name"`
}

type Config struct {
	MaxLen      int     `json:"max_len"`
	SelectColor string  `json:"select_color"`
	H           float64 `json:"h"`
	S           float64 `json:"s"`
	L           float64 `json:"l"`
}

type jsonLoad struct {
	Settings Config    `json:"settings"`
	Projects []Project `json:"projects"`
}

func getPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path: %v", err)
		return "", err
	}

	exeDir := filepath.Dir(exePath)
	jsonPath := filepath.Join(exeDir, "carpoon.json")

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		log.Printf("json file not found, creating default 'carpoon.json'\n")

		defaultData := jsonLoad{
			Settings: Config{
				MaxLen:      10,
				SelectColor: "#6C3BAA",
				H:           266,
				S:           48,
				L:           45,
			},
			Projects: []Project{},
		}

		dObj, err := json.MarshalIndent(defaultData, "", "  ")
		if err != nil {
			log.Fatal("json object could not be created")
			return "", err
		}

		if err := os.WriteFile(jsonPath, dObj, 0644); err != nil {
			log.Fatal("cant create json, try manually creating the json file")
			return "", err
		}
	}
	return jsonPath, nil
}

func (data *jsonLoad) Init() error {
	jsonPath, err := getPath()
	if err != nil {
		return err
	}

	byteVal, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("cant open file %s: %v", jsonPath, err)
		return err
	}

	if err := json.Unmarshal(byteVal, data); err != nil {
		log.Fatalf("cant unmarshal json data: %v", err)
		return err
	}

	return nil
}

func (data *jsonLoad) changeColor(hex string, H float64, S float64, L float64) error {
	data.Settings.SelectColor = hex
	data.Settings.H = H
	data.Settings.S = S
	data.Settings.L = L
	return data.setData()
}

func (data *jsonLoad) setData() error {
	jsonPath, err := getPath()
	if err != nil {
		return err
	}

	updatedJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("cant marshal json")
		return err
	}

	err = os.WriteFile(jsonPath, updatedJson, 0644)
	if err != nil {
		log.Fatal("cant write json")
		return err
	}
	return nil
}

func (data *jsonLoad) addProject(path string, title string) error {
	if len(data.Projects) >= data.Settings.MaxLen {
		return fmt.Errorf("config max length reached (%d)", data.Settings.MaxLen)
	}

	newProject := Project{
		Path:  path,
		Title: title,
	}

	data.Projects = append(data.Projects, newProject)
	return data.setData()
}

func (data *jsonLoad) removeProject(indexToRemove int) error {
	if len(data.Projects) == 0 {
		return fmt.Errorf("no projects to remove")
	}

	if indexToRemove < 0 || indexToRemove >= len(data.Projects) {
		return fmt.Errorf("invalid index: %d", indexToRemove)
	}

	data.Projects = append(data.Projects[:indexToRemove], data.Projects[indexToRemove+1:]...)

	return data.setData()
}
