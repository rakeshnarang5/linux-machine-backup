package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"./res"
	"strconv"
	"sync"
	"time"
	"github.com/ctcpip/notifize"
)

type Task struct {
	Name  string
	Count int
}

type TaskMap struct {
	Tasks map[string]Task
	Count int
}

func (t *TaskMap) addTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Name your task")
	name, _ := reader.ReadString('\n')
	t.Count++
	key := strconv.Itoa(t.Count)
	t.Tasks[key] = Task{name[:len(name)-1], 1}
}

func timeconv(pomo int) (int, int) {
	min := pomo * 25
	hh := min / 60
	mm := min % 60
	return hh, mm
}

func (t *TaskMap) status() {
	fmt.Printf("Total tasks: %d\n\n", t.Count)
	var str string
	for key, tsk := range t.Tasks {
		hh, mm := timeconv(tsk.Count)
		//fmt.Printf("Key: %s\nName: %s\nSessions: %d\nTime Spent: %02d:%02d\n\n",key, tsk.Name,tsk.Count, hh, mm)
		str += fmt.Sprintf("Key: %s\nName: %s\nSessions: %d\nTime Spent: %02d:%02d\n\n",key, tsk.Name,tsk.Count, hh, mm)
	}
	fmt.Println(str)
	Notify(str)
}

func (t *TaskMap) update() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press 1 if you work on an existing task\nPress 2 if you worked on a new task\n")
	input, _ := reader.ReadString('\n')
	if input == "1\n" {
		t.status()
		fmt.Println("Which task did you work on?")
		chc,_ := reader.ReadString('\n')
		tsk := t.Tasks[chc[:len(chc)-1]]
		tsk.Count++
		t.Tasks[chc[:len(chc)-1]] = tsk
		fmt.Println("updated!")
		t.status()
	}
	if input == "2\n" {
		t.addTask()
	}
	t.saveState()
	t.timer("break")
}

func Notify(msg string) {
	go 0notifize.Display("", msg, false, "/home/nagarro/workspace/src/timeManager/img/time.jpg")
	res.SendMessage(msg)

}

func (t *TaskMap) timer(tsk string) {
	var d time.Duration
	if tsk == "work" {
		d = time.Second * 20
	} else if tsk == "break" {
		d = time.Second * 10
	}
	timer := time.NewTimer(d)
	<-timer.C 
	if tsk == "work" {
		Notify("Timer finished!\nWhat did you work on?")
		t.update()
	} else if tsk == "break" {
		Notify("Break over!\nGet back to work!")
	}
}

func (t *TaskMap) saveState() {
	m, _ := json.Marshal(t)
	ioutil.WriteFile("./task.json",m,0755)
}

func main() {
	taskmap := initializeTaskMap()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Would you like to continue with saved tasks (Y/n)")
	input, _ := reader.ReadString('\n')
	//fmt.Println([]byte(input))

	if input == "Y\n" {
		data, _ := ioutil.ReadFile("./task.json")
		json.Unmarshal(data, &taskmap)
	}

	fmt.Println(taskmap)

	var wg sync.WaitGroup

	wg.Add(1)

	go menu(taskmap)

	//fmt.Println(taskmap)

	// task1 := Task{"first task", 1}
	// task2 := Task{"second task", 2}
	// task3 := Task{"third task", 3}
	// taskmap.Tasks["1"] = task1
	// taskmap.Tasks["2"] = task2
	// taskmap.Tasks["3"] = task3
	// fmt.Println(taskmap.Tasks)
	// m, _ := json.Marshal(taskmap.Tasks)
	// fmt.Println(string(m))
	// x, _ := json.Marshal(taskmap)
	// fmt.Println(string(x))
	// ioutil.WriteFile("./task.json", x, 0755)
	//unmarshal
	//var y TaskMap
	// data, _ := ioutil.ReadFile("./task.json")
	// json.Unmarshal(data, &y)
	// fmt.Println(y)
	wg.Wait()
}

func menu(t *TaskMap) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Choose an Action:-\nPress 1 to see current tasks\nPress 2 to start working\nPress 3 to end the day")
		input, _ := reader.ReadString('\n')
		if input == "1\n" {
			t.status()
		}
		if input == "2\n" {
			t.timer("work")
		}
		if input == "3\n" {
			t.saveState()
			os.Exit(0)
		}
	}

}

func initializeTaskMap() *TaskMap {
	var taskmap TaskMap
	taskmap.Tasks = make(map[string]Task)
	taskmap.Count = 3
	return &taskmap
}
