package ui

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/manifoldco/promptui"
)

const (
	otherLabel = "...Other..."
)

func runSelectPrompt(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, retval, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return retval, nil
}

func promptUserSelect() (string, error) {
	return runSelectPrompt("Select user", []string{
		"ec2-user",
		"centos",
		"ubuntu",
		otherLabel,
	})
}

func promptUserInput() string {
	prompt := promptui.Prompt{
		Label: "Insert remote user",
	}
	user, err := prompt.Run()
	if err != nil {
		return ""
	}
	return user
}

func promptKeysSelect() (string, error) {
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
	return runSelectPrompt("Select private key", keys)
}

func startPrompt() (string, string, error) {
	user, err := promptUserSelect()
	if err != nil {
		return "", "", err
	}
	if user == otherLabel {
		user = promptUserInput()
	}
	pkey, err := promptKeysSelect()
	if err != nil {
		return "", "", err
	}
	return user, pkey, nil
}
