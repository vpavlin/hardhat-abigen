# Hardhat Abigen

This tool parses content of Hardhat produced `artifacts` directory, extracts all valid ABIs and generates Go bindings using `abigen`.

It assumes `abigen` commandline tool is installed and available on `$PATH`. You can download `abigen` from https://geth.ethereum.org/downloads/ 

## Usage

```
./hardhat-abigen [options] [FILE|DIRECTORY]
  --bindings
	 Use Abigen to generate bindings (default: 'true')
  --exclude
	 List of strings to match aginst paths to exclude them from processing (default: '')
  --outDir
	 Destination directory for extracted ABI json files and generated bindings (default: 'output')
```