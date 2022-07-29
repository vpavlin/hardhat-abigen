package extractor

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vpavlin/hardhat-abigen/internal/abigen"
	"github.com/vpavlin/hardhat-abigen/internal/types"
	"github.com/vpavlin/hardhat-abigen/internal/utils"
)

func Extract(in io.Reader) (*types.ABI, error) {
	bytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	abi := new(types.ABI)
	err = json.Unmarshal(bytes, abi)
	if err != nil {
		return nil, err
	}

	if abi.ABI == nil || len(abi.ABI.([]interface{})) == 0 || abi.ByteCode == "" || abi.ByteCode == "0x" {
		return nil, nil
	}

	return abi, nil
}

func Dump(out io.Writer, abi *types.ABI) error {
	marshalled, err := json.Marshal(abi.ABI)
	if err != nil {
		return err
	}

	_, err = out.Write(marshalled)
	if err != nil {
		return err
	}

	return nil
}

func ExtractFromFile(in string, out string) (string, error) {
	path, err := filepath.Abs(in)
	if err != nil {
		return "", nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	abi, err := Extract(fp)
	if err != nil {
		return "", err
	}

	if abi == nil {
		return "", nil
	}

	fpOut, err := os.Create(out)
	if err != nil {
		return "", err
	}
	defer fpOut.Close()

	err = Dump(fpOut, abi)
	if err != nil {
		return "", nil
	}

	return abi.ContractName, nil
}

func ExtractWalk(outDir string, p string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		dir, err := ioutil.ReadDir(p)
		if err != nil {
			return err
		}
		for _, fInfo := range dir {
			newPath := path.Join(p, fInfo.Name())
			filepath.WalkDir(newPath, func(path string, d fs.DirEntry, err error) error {
				return nil //ExtractWalk(path, info, err)
			})
		}
	}

	//outFile := path.Join(outDir, info.Name

	return nil

}

func Excluded(toCheck string, exclude string) bool {
	for _, item := range strings.Split(exclude, ",") {
		if strings.Contains(toCheck, item) {
			return true
		}
	}
	return false
}

func ProcessSingle(inFile string, outDir string, generateBindings bool, exclude string) error {
	if Excluded(inFile, exclude) {
		logrus.Infof("Skipping path %s based on --exclude", inFile)
		return nil
	}

	abiOutDir, err := utils.MkOutDirIfNotExist(outDir, utils.AbiDir)
	if err != nil {
		return err
	}

	if len(inFile) == 0 {
		return fmt.Errorf("You need to provide a patch to the JSON file containing ABI produced by Hardhat")
	}

	abiOutFile := path.Join(abiOutDir, filepath.Base(inFile))

	contractName, err := ExtractFromFile(inFile, abiOutFile)
	if err != nil {
		return fmt.Errorf("Failed to extract abi: %s", err)
	}

	if generateBindings && len(contractName) > 0 {
		bindingsOutDir, err := utils.MkOutDirIfNotExist(outDir, utils.BindingsDir)
		if err != nil {
			return err
		}
		bindingsOutFile := path.Join(bindingsOutDir, fmt.Sprintf("%s.%s", contractName, "go"))
		logrus.Infof("Generating bindings for %s", contractName)
		err = abigen.Run(abiOutFile, contractName, "bindings", bindingsOutFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessAll(inFile string, outDir string, generateBindings bool, exclude string) error {
	return filepath.Walk(inFile, func(p string, info fs.FileInfo, err error) error {
		if filepath.Ext(p) == ".json" {
			return ProcessSingle(p, outDir, generateBindings, exclude)
		}
		return nil
	})
}
