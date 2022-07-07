package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	//"ioutil"
)

var flagId string
var flagOperation string
var flagItem string
var flagFilename string

func init() {
	flag.StringVar(&flagId, "id", "", "help message for flagname")
	flag.StringVar(&flagOperation, "operation", "", "help message for flagname")
	flag.StringVar(&flagItem, "item", "", "help message for flagname")
	flag.StringVar(&flagFilename, "fileName", "", "help message for flagname")
}

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Items []Item

func parseArgs() Arguments {
	arg := make(map[string]string)
	flag.Parse()
	arg["id"] = flagId
	arg["operation"] = flagOperation
	arg["item"] = flagItem
	arg["fileName"] = flagFilename
	return arg
}

func Add(flagFilename string, flagItem string, writer io.Writer) error {
	var it Items
	var j Item
	file, err := os.OpenFile(flagFilename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) != 0 {
		err = json.Unmarshal(b, &it)
		if err != nil {
			return err
		}
	}
	err = json.Unmarshal([]byte(flagItem), &j)
	if err != nil {
		return err
	}
	for i := 0; i < len(it); i++ {
		if it[i].Id == j.Id {
			s := fmt.Sprintf("Item with id %s already exists", j.Id)
			writer.Write([]byte(s))
		}
	}
	it = append(it, j)
	f, err := json.Marshal(&it)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(flagFilename, f, 0777)
	if err != nil {
		return err
	}
	return nil
}

func Remove(flagFilename string, flagId string, writer io.Writer) error {
	var it Items
	t := false
	file, err := os.OpenFile(flagFilename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) != 0 {
		err = json.Unmarshal(b, &it)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(it); i++ {
		if it[i].Id == flagId {
			it[i] = it[len(it)-1]
			it = it[:len(it)-1]
			t = true
			f, err := json.Marshal(&it)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(flagFilename, f, 0777)
			if err != nil {
				return err
			}
		}
	}
	if t == false {
		s := fmt.Sprintf("Item with id %s not found", flagId)
		writer.Write([]byte(s))
	}
	return nil
}

func FindById(flagFilename string, flagId string, writer io.Writer) error {
	var it Items
	t := false
	file, err := os.OpenFile(flagFilename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) != 0 {
		err = json.Unmarshal(b, &it)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(it); i++ {
		if it[i].Id == flagId {
			t = true
			f, err := json.Marshal(&it[i])
			writer.Write(f)
			if err != nil {
				return err
			}
		}
	}
	if t == false {
		writer.Write([]byte(""))
	}
	return nil
}

func List(flagFilename string, writer io.Writer) error {
	var it Items
	t := false
	file, err := os.OpenFile(flagFilename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) != 0 {
		t = true
		err = json.Unmarshal(b, &it)
		if err != nil {
			return err
		}
	}
	if t == true {
		f, err := json.Marshal(&it)
		writer.Write(f)
		if err != nil {
			return err
		}
		//for i := 0; i < len(it); i++ {
		//	f, err := json.Marshal(&it[i])
		//	writer.Write(f)
		//	if err != nil {
		//		return err
		//	}
		//}
	}
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	var err error
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	switch args["operation"] {
	case "":
		return fmt.Errorf("-operation flag has to be specified")
	case "add":
		{
			if args["item"] == "" {
				//writer.Write([]byte("-item flag has to be specified"))
				return fmt.Errorf("-item flag has to be specified")
				//return err
			} else {
				err = Add(args["fileName"], args["item"], writer)
			}
		}
	case "remove":
		{
			if args["id"] == "" {
				//writer.Write([]byte("-id flag has to be specified"))
				return fmt.Errorf("-id flag has to be specified")
				//return err
			} else {
				err = Remove(args["fileName"], args["id"], writer)
			}
		}
	case "findById":
		{
			if args["id"] == "" {
				//writer.Write([]byte("-id flag has to be specified"))
				return fmt.Errorf("-id flag has to be specified")
				//return err
			} else {
				err = FindById(args["fileName"], args["id"], writer)
			}
		}
	case "list":
		err = List(args["fileName"], writer)
	default:
		{
			s := fmt.Sprintf("Operation %s not allowed!", args["operation"])
			return fmt.Errorf(s)
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
