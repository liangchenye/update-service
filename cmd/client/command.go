package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/liangchenye/update-service/utils"
)

var addCommand = cli.Command{
	Name:  "add",
	Usage: "add a repository url",

	Action: func(context *cli.Context) error {
		proto := context.Args().Get(0)
		url := context.Args().Get(1)
		ucc, _ := DefaultUpdateClientConfig()
		if err := ucc.Add(proto, url); err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Success in adding %s %s.\n", proto, url)
		return nil
	},
}

var removeCommand = cli.Command{
	Name:  "remove",
	Usage: "remove a repository url",

	Action: func(context *cli.Context) error {
		proto := context.Args().Get(0)
		url := context.Args().Get(1)
		ucc, _ := DefaultUpdateClientConfig()
		if err := ucc.Remove(proto, url); err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Success in removing %s.\n", proto, url)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list the saved repositories or appliances of a certain repository",

	Action: func(context *cli.Context) error {
		if len(context.Args()) == 0 {
			ucc, _ := DefaultUpdateClientConfig()
			for _, repo := range ucc.Repos {
				fmt.Println(repo)
			}
		} else if len(context.Args()) == 2 {
			proto := context.Args().Get(0)
			url := context.Args().Get(1)
			repo, _ := NewUpdateClientRepo(proto, url)
			apps, err := repo.List()
			if err != nil {
				fmt.Println(err)
				return err
			}

			for _, app := range apps {
				fmt.Println(app)
			}

			ucc, _ := DefaultUpdateClientConfig()
			ucc.Add(proto, url)
		}
		return nil
	},
}

var pushCommand = cli.Command{
	Name:  "push",
	Usage: "push a file to a repository",

	Action: func(context *cli.Context) error {
		//TODO: we can have a default repo
		if len(context.Args()) != 3 {
			err := errors.New("wrong syntax: push 'proto' 'repo url' 'local filepath'")
			fmt.Println(err)
			return err
		}

		proto := context.Args().Get(0)
		url := context.Args().Get(1)
		file := context.Args().Get(2)
		repo, _ := NewUpdateClientRepo(proto, url)

		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = repo.Put(filepath.Base(file), content)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	},
}

var pullCommand = cli.Command{
	Name:  "pull",
	Usage: "pull a file from a repository",

	Action: func(context *cli.Context) error {
		//TODO: we can have a default repo
		if len(context.Args()) != 3 {
			err := errors.New("wrong syntax: pull  'proto' 'repo url' 'name' ")
			fmt.Println(err)
			return err
		}

		proto := context.Args().Get(0)
		url := context.Args().Get(1)
		name := context.Args().Get(2)
		repo, _ := NewUpdateClientRepo(proto, url)
		ucc, _ := DefaultUpdateClientConfig()
		repo.SetCacheDir(ucc.GetCacheDir())

		fmt.Println("start to download and verify meta data")
		err := repo.Sync()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("success in downloading and verifying meta data")

		fmt.Println("start to download file")
		savedURL, err := repo.Get(name)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("file downloaded to: ", savedURL)

		fmt.Println("start to compare the hash value")
		sha, err := repo.GetSHAS(name)
		if err != nil {
			fmt.Println(err)
			return err
		}

		data, _ := ioutil.ReadFile(savedURL)
		shaCal, _ := utils.SHA512(data)

		if sha != shaCal {
			message := fmt.Sprintf("The downloaded file is invalid, expected sha: <%s>, but get: <%s>", sha, shaCal)
			return errors.New(message)
		}
		fmt.Println("the downloaded file is valid.")
		return nil
	},
}
