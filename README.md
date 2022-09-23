# NVMB
- WORK IN PROGRESS
- Stands for Network of Virtual Machine Blocks.
- A simulation of a decentralized undirected smart network system.
- Decentralization is implemented by partitioning of global(network) and local (node) data into discrete parts.

# Dependencies
- Recommended libraries to have are:
  * OpenMP
- Required libraries are:
  * GoogleTest
  * Armadillo 
    * `install these libraries for optimization`: 
      - LAPACK 
      - BLAS 
- Data files used for testing: 
  - SOURCE: https://www.kaggle.com/mlg-ulb/creditcardfraud
  - NAME IT: `creditcard.csv`
  - PUT IT IN: `golem/golem_base/datos/other_data/`

  - SOURCE: https://www.kaggle.com/uciml/glass 
  - NAME IT: `glass.csv` 
  - PUT IT IN: `golem/golem_base/datos/other_data/`

  - SOURCE: https://www.kaggle.com/prathamtripathi/drug-classification
  - NAME IT: `drug200.csv`
  - PUT IT IN: `golem/golem_base/datos/other_data/`

  - SOURCE: https://www.kaggle.com/arshid/iris-flower-dataset
  - NAME IT: `IRIS.csv`
  - PUT IT IN: `golem/golem_base/datos/other_data/`

# Build instructions
- Open up command line.
- `cd` into `NVMB`. 
- Enable automatic `go get` commands from make file:
`go env -w GO111MODULE=auto`
- `cd` into the directory `build`
- type in command `cmake ..`. 
- type in command `make`.
- binaries have been built.
- run binaries with `./BINARY_NAME`

# Test instructions
After successfully performing build instructions,
- test the C++ portion:
    - `cd` into `build/test`.
    - run `./testy`.
- test the Golang portion by
    - downloading and placing the data files into the appropriate folders as mentioned in the section `dependencies`.
    - `cd` into `golem`.
    - uncomment the code that generates the csv data.
    - run `go run main.go`.
    - test by `go test ./...`.

## UPDATE: 9/19/22
- will continue working on automating functionalities within these next few months.

## Update: 5/8/2021 
- Planned restructuring and re-write of some components of network. 

## Update: 5/7/2021 
- Required libraries include Armadillo.


## Update: Machine-learning 11/26/2020
- `golem.DForest` class has been coded but not yet perfected.
- `golem.Sampler` class needs to be optimized.

# [NOTE] 
- Machine learning algorithms and other functionalities will be pushed
  whenever contributor/s feel necessary and/or like it.
- Also, some refactoring is necessary. 
- Plenty of suboptimal parts to this software. Feel free to improve it. 