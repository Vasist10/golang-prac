package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Task struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`

}
const taskfile = "tasks.json"

func LoadTasks() ([]Task,error){
	var tasks []Task
	
	if _, err := os.Stat(taskfile); os.IsNotExist(err) {
		return tasks, nil
	}
	data,err :=ioutil.ReadFile(taskfile)
	if err != nil {	
		return nil, err
	}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(taskfile, data, 0644)
}
