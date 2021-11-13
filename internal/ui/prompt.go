package ui

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/kyokomi/emoji"
	"github.com/manifoldco/promptui"
)

const (
	otherLabel = "...Other..."
)

func runSelectPrompt(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  30,
	}
	_, retval, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return retval, nil
}

func promptUserSelect() (string, error) {
	return runSelectPrompt(emoji.Sprintf(":man_technologist: Select remote user"), []string{
		"ec2-user",
		"centos",
		"ubuntu",
		otherLabel,
	})
}

func promptUserInput() string {
	prompt := promptui.Prompt{
		Label: emoji.Sprintf(":man_technologist: Insert remote user"),
	}
	user, err := prompt.Run()
	if err != nil {
		return ""
	}
	return user
}

func promptKeysSelect(awsKeyname string) (string, error) {
	keys := []string{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/.ssh/", homedir))
	if err != nil {
		return "", err
	}
	for _, f := range files {
		keys = append(keys, f.Name())
	}

	label := emoji.Sprintf("Select Private Key | :computer: %s", awsKeyname)
	return runSelectPrompt(label, keys)
}

func promptPortInput() string {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid port number")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    emoji.Sprintf("Insert SSH Port"),
		Validate: validate,
		Default:  "22",
	}
	port, err := prompt.Run()
	if err != nil {
		return ""
	}
	return port
}

func startPrompt(awsKeyname string) (string, string, string, error) {
	port := promptPortInput()
	user, err := promptUserSelect()
	if err != nil {
		return "", "", "", err
	}
	if user == otherLabel {
		user = promptUserInput()
	}
	pkey, err := promptKeysSelect(awsKeyname)
	if err != nil {
		return "", "", "", err
	}
	return port, user, pkey, nil
}
