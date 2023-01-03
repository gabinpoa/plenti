package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// EndpointFlag disables the route for a content source by omitting the corresponding svelte template.
var EndpointFlag bool

// SingleTypeFlag create a one time file at the top level of content.
var SingleTypeFlag bool

// typeCmd represents the type command
var typeCmd = &cobra.Command{
	Use:   "type [name]",
	Short: "A content type with structured fields",
	Long: `Types allow you to group content by their data structure.

The following are examples of types you could create that share common fields:
- pages
- blog_posts
- news
- events

You can define any type you'd like, with any field structure you desire.
There are no required fields when creating your new type.

Any individual file within a type can contain variations in its field structure.
Just make sure to account for this in the corresponding '/layouts/content/<your_type>.svelte' file.

Optionally add a _defaults.json file to define the default content for the type when creating a new instance.

Optionally add a _schema.json file to define the input widgets used in the editor.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}
		if len(args) > 1 {
			return errors.New("names cannot have spaces")
		}
		if len(args) == 1 {
			return nil
		}
		return fmt.Errorf("invalid name specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		typeName := args[0]

		// shoud we stop here on error from either?
		if SingleTypeFlag {
			err := createSingleType(typeName)
			if err != nil {
				log.Fatal("Can't create single type %w", err)
			}
		} else {
			err := createMultiType(typeName)
			if err != nil {
				log.Fatal("Can't create new type %w", err)
			}
		}

		if EndpointFlag {

			typeLayoutPath := fmt.Sprintf("layouts/content/%s.svelte", strings.Trim(typeName, " /"))
			if _, err := os.Stat(typeLayoutPath); !os.IsNotExist(err) {
				fmt.Printf("A Type layout with the same name located at \"%s\" already exists\n", typeLayoutPath)
				return
			}

			fmt.Printf("Creating new Type layout: %s\n", typeLayoutPath)
			if _, err := os.OpenFile(typeLayoutPath, os.O_RDONLY|os.O_CREATE, os.ModePerm); err != nil {
				log.Fatal("Can't create layout for type \"%s\": %w", typeName, err)
			}
		}

	},
}

func createMultiType(typeName string) error {
	typeContentPath := fmt.Sprintf("content/%s", strings.Trim(typeName, " /"))

	if _, err := os.Stat(typeContentPath); !os.IsNotExist(err) {
		// The path already exists
		return fmt.Errorf("A Type content source with the same name located at \"%s/\" already exists\n", typeContentPath)
	}

	if _, err := os.Stat(typeContentPath + ".json"); !os.IsNotExist(err) {
		// error or not?
		return fmt.Errorf("A single file Type content source with the same name located at \"%s.json\" already exists\n", typeContentPath)
	}

	fmt.Printf("Creating new Type content source: %s/\n", typeContentPath)
	if err := os.MkdirAll(typeContentPath, os.ModePerm); err != nil {
		return fmt.Errorf("Can't create type named \"%s\": %w\n", typeName, err)
	}
	err := createJSONFile(typeContentPath + "/_defaults.json")
	if err != nil {
		return fmt.Errorf("Can't create _defaults.json for type \"%s\": %w\n", typeName, err)
	}
	err = createJSONFile(typeContentPath + "/_schema.json")
	if err != nil {
		return fmt.Errorf("Can't create _schema.json for type \"%s\": %w\n", typeName, err)
	}
	return nil
}

func createSingleType(typeName string) error {
	singleTypePath := fmt.Sprintf("content/%s.json", typeName)

	if _, err := os.Stat(singleTypePath); !os.IsNotExist(err) {
		return fmt.Errorf("A single type content source with the same name located at \"%s\" already exists\n", singleTypePath)
	}

	fmt.Printf("Creating new single type content source: %s\n", singleTypePath)

	err := createJSONFile(singleTypePath)
	if err != nil {
		return fmt.Errorf("Can't create single type content source for: \"%s\": %w\n", typeName, err)
	}

	return nil
}

func createJSONFile(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Can't create file: \"%s\": %w\n", filePath, err)
	}

	_, err = f.Write([]byte("{}"))
	if err != nil {
		return fmt.Errorf("Can't add empty curly brackets to file: \"%s\": %w\n", filePath, err)
	}
	// can be non-nil error
	return f.Close()
}

func init() {
	newCmd.AddCommand(typeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// typeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// typeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	typeCmd.Flags().BoolVarP(&EndpointFlag, "endpoint", "e", true, "set 'false' to disable route.")
	typeCmd.Flags().BoolVarP(&SingleTypeFlag, "single", "s", false, "set 'true' to generate single content file.")
}
