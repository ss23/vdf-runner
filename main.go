package main

import (
	"flag"
	"fmt"
	"os"
	//"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"

	"github.com/galaco/KeyValues"
)

// TODO: Implement this properly...
var variables = []string{
	"%INSTALLDIR%",
	"%CDKEY%",
}

var replacements = []string{
	"S:\\crisp\\steamapps\\common\\Spore",
	"REDACTEDLOL",
}

func main() {
	var vdfPath string
	flag.StringVar(&vdfPath, "vdfPath", "path", "The path to the VDF to parse")

	flag.Parse()

	if vdfPath == "path" {
		println("Please specify a path")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	// Lets try open it
	file, err := os.Open(vdfPath)
	if err != nil {
		panic(err)
	}

	// Lets parse our VDF
	reader := keyvalues.NewReader(file)

	kv, err := reader.Read()
	if err != nil {
		panic(err)
	}

	is, err := kv.Find("installscript")
	if err != nil {
		panic(err)
	}

	// Lets loop over each of the installscript children and perform the relevant tasks
	children, err := is.Children()
	if err != nil {
		panic(err)
	}

	for _, child := range children {
		switch child.Key() {
		case "registry":
			fmt.Println("Running registry install code...")
			runRegistry(child)
		case "run process":
			fmt.Println("Running install processes...")
			runProcess(child)
		default:
			println("Unknown operation specified. Please submit a bug with the .vdf as an example." + child.Key())
		}
	}

	fmt.Println("Completed!")
}

func runRegistry(kv *keyvalues.KeyValue) {
	// Each of the children is the path to the registry "folder" we'll be operating in
	children, err := kv.Children()
	if err != nil {
		panic(err)
	}

	for _, child := range children {
		// We need to "fix" the key. Technically it should be done by the parser, but here we are
		key := strings.ReplaceAll(strings.TrimSpace(child.Key()), "\\\\", "\\")
		root, path := getRegistryPath(key)
		// Attempt to create/open the key
		k, exists, err := registry.CreateKey(root, path, registry.ALL_ACCESS)
		if err != nil {
			panic(err)
		}
		if exists {
			fmt.Println("Opened existing key:", key)
		} else {
			fmt.Println("Created new key:", key)
		}

		// With the key open, we can now loop over the children and perform the required setup
		types, err := child.Children()
		if err != nil {
			panic(err)
		}
		for _, t := range types {
			switch t.Key() {
			case "dword":
				final_keys, err := t.Children()
				if err != nil {
					panic(err)
				}
				for _, final_key := range final_keys {
					// Parse out value
					i, err := final_key.AsInt()
					if err != nil {
						panic(err)
					}
					i2 := uint32(i)
					fmt.Println("Creating registry key:", key+"\\"+final_key.Key(), ":", i2)
					// Create these keys!
					err = k.SetDWordValue(final_key.Key(), i2)
					if err != nil {
						panic(err)
					}
				}
			case "string":
				final_keys, err := t.Children()
				if err != nil {
					panic(err)
				}
				for _, final_key := range final_keys {
					// Parse out value
					s, err := final_key.AsString()
					// Do replacements of variables where required
					s = variableReplacements(s)
					if err != nil {
						if err.Error() == "value is not of type string" {
							// This is a problem with our parsing, but hopefully we can just ignore it
							// We probably should figure out what these kind of keys mean. PRs wanted
							continue
						}
						panic(err)
					}
					fmt.Println("Creating registry key:", key+"\\"+final_key.Key(), ":", s)
					err = k.SetStringValue(final_key.Key(), s)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
	return
}

func getRegistryPath(path string) (registry.Key, string) {
	// Begin by splitting off the first component and match that to an already open registry key
	var k registry.Key
	keys := strings.SplitN(path, "\\", 2)
	switch keys[0] {
	case "hkey_local_machine":
		k = registry.LOCAL_MACHINE
	default:
		panic("Unknown registry key specified")
	}

	return k, keys[1]
}

func variableReplacements(s string) string {
	// Begin by replacing the double slashes
	s = strings.ReplaceAll(s, "\\\\", "\\")

	for i, _ := range variables {
		s = strings.ReplaceAll(s, variables[i], replacements[i])
	}
	if strings.Contains(s, "%") {
		fmt.Println("Found a character that shouldn't be in a string -- bad variable replacement?")
		panic(s)
	}
	return s
}

func runProcess(kv *keyvalues.KeyValue) {
	panic("runProcess not implemented!")
	return
}
