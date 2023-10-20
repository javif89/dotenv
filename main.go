package dotenv

import (
	"bufio"
	"os"
	"strings"
)

type EnvValue struct {
	Key string
	Value string
}

type EnvFile struct {
	Values []EnvValue
	Path string
}

func (e *EnvFile) Add(key string, value string) {
	key = standardizeKey(key)
	value = cleanString(value)
	if e.Has(key) {
		e.Set(key, value)
		return
	}
	e.Values = append(e.Values, EnvValue{Key: key, Value: value})
}

func (e *EnvFile) Remove(key string) {
	for i, v := range e.Values {
		if v.Key == key {
			e.Values = append(e.Values[:i], e.Values[i+1:]...)
		}
	}
}

func (e *EnvFile) Has(key string) bool {
	for _, v := range e.Values {
		if v.Key == key {
			return true
		}
	}
	return false
}

func (e *EnvFile) Get(key string) string {
	key = standardizeKey(key)
	for _, v := range e.Values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func (e *EnvFile) Set(key string, value string) {
	for i, v := range e.Values {
		if v.Key == key {
			e.Values[i].Value = cleanString(value)
		}
	}
}

func (e *EnvFile) Save() {
	// Check if file exists if not create it
	if _, err := os.Stat(e.Path); os.IsNotExist(err) {
		file, err := os.Create(e.Path)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	file, err := os.OpenFile(e.Path, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, v := range e.Values {
		file.WriteString(v.Key + "=" + formatValueForPrint(v.Value) + "\n")
	}
}

func New(path string) *EnvFile {
	return &EnvFile{Path: path}
}

func LoadOrCreate(path string) *EnvFile {
	if !fileExists(path) {
		return New(path)
	}

	return Load(path)
}

func Load(path string) *EnvFile {
	file, e := os.Open(path)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	env := EnvFile{Path: path}

	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}

		kv := strings.Split(scanner.Text(), "=")

		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		value := kv[1]

		env.Add(key, value)
	}

	return &env
}