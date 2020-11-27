// this file will be responsible for reading and writing a node's data to its
//  file representation.

#include "node_file_reader.hpp"
#include <iostream>
#include <ostream>
#include <assert.h>
#include <string.h>
#include <unistd.h>
#include <fstream>
#include <sstream>
#include <iterator>
using namespace std;

string DATA_DIRECTORY = "/src/data";

vector<string> TIMESTAMP_HEADER = {"time", "event_type", "impact", "role", "nodes", 
    "responses", "currency", "neighbors", "contracts", "transmission", "cfAcceptanceRate",
    "competition", "greed", "negotiation", "growth", "nodeRank", "knownNodes", "competitors",
    "numPathsRatio", "pathsInfo", "competitionMeasure", "contractMeasure", "newContractMeasure",
    "bondAdvantageMeasure", "bondDeletionAdvantageMeasure"};

/*
checks for existence of file at current working directory
*/
int FileExists(const char *path)
{
    struct stat fileStat;
    if ( stat(path, &fileStat) )
    {
        return 0;
    }
    if ( !S_ISREG(fileStat.st_mode) )
    {
        return 0;
    }
    return 1;
}

/*
checks for existence of directory at current working directory
*/
int DirExists(const char *path)
{
    struct stat fileStat;
    if ( stat(path, &fileStat) )
    {
        return 0;
    }
    if ( !S_ISDIR(fileStat.st_mode) )
    {
        return 0;
    }
    return 1;
}

/*
a file name is in correct format if there are no / (directory)
CAUTION: partial check
*/
int IsCorrectFileNameFormat(string fn) {
  // check if directory
  for (int i = 0 ; i < fn.length(); i++)
  {
    if (fn[i] == '/') {
      return 0;
    }
  }
  return 1;
}

/*
splits a path (str) into a vector of tokens
*/
vector<string> SplitStringToVector(string s, const char* delimiter) {
  // convert string to char array
  const char * x2 = s.c_str();
  char x[s.length()];
  strcpy(x, x2);

  vector<string> v;
  char * pch;
  pch = strtok(x, "/");
  while (pch != NULL)
  {
    string str(pch);
    v.push_back(str);
    pch = strtok (NULL, "/");
  }
  return v;
}



/*
makes a string of the form "[../..]..." for moving current working directory back
`numDirBack` directories.
*/
string MakeBackTraversalString(int numDirBack) {

  string x = "../";
  string y = "";
  for (int i = 0; i < numDirBack; i++) {
    y = y + x;
  }
  return y;
}

/*
initializes instance of NodeFileReader
*/
NodeFileReader::NodeFileReader(string fn) {
  // check for correct file format
  assert (IsCorrectFileNameFormat(fn) == true);
  filename = fn;
  InitializeFile();
  datDir = "";
  WriteHeaderRow(TIMESTAMP_HEADER); 
}

NodeFileReader::NodeFileReader(std::string dd, std::string fn) {
  assert (IsCorrectFileNameFormat(dd) == true);
  assert (IsCorrectFileNameFormat(fn) == true);

  filename = fn;
  datDir = dd;
  MakeDataDirectory(); 
  InitializeFile();
  WriteHeaderRow(TIMESTAMP_HEADER);  
}

/*
// description 
this is the first write of file 
*/ 
void NodeFileReader::WriteHeaderRow(vector<string> row) { 

    string s = ""; 

    for (int i = 0; i < row.size() -1; i++) {
        s += row[i] + ","; 
    } 

    s += row[row.size() -1]; 

    fileObj.open(filePath, ofstream::app);
    fileObj << s << endl;  
}

/*
// description 
writes row data to fileObj 
*/ 
void NodeFileReader::WriteRowData(string row) {
    fileObj << row + "\n"; 
}

/*
initializes a file at "./data/fp"
*/
void NodeFileReader::InitializeFile() {

  // travel to p
  TravelToBaseDir();

  // TODO refactor below
  char cwd[256];
  char* buff = getcwd(cwd, sizeof(cwd));
  string b(buff);

  // create file for content storage
  string s = ""; 
  if (datDir != "") {
      s = "/"; 
  }
  filePath = b + DATA_DIRECTORY + "/" + datDir + s + filename;

  ofstream fx(filePath);
  fx.close();
  TravelToPreviousCWD();
}

/*
// description 
*/ 
void NodeFileReader::MakeDataDirectory() {
    TravelToBaseDir();
    mkdir(("src/data/" + datDir).c_str(), 0777); 
}

/*
// description 
*/ 
void NodeFileReader::CloseWrite() {
    fileObj.close(); 
}

/*
iterates through and collects position of delimiters into integer array
*/
vector<int> NodeFileReader::GetDelimiterPositions() {
  // set bytes
  vector<int> x;
  int i = 0;

  std::ifstream input(filePath, std::ios::binary);
  std::string line;
  int pos;
  while (std::getline(input, line))
  {
      std::istringstream iss(line);
      if (line.compare("#--->") == 0) {
        x.push_back(input.tellg());
      }
  }
  return x;
}

int NodeFileReader::TravelToPreviousCWD() {
  int success = chdir(pathToPrevCWD.c_str());
  pathToPrevCWD = "";
  return success;
}


/*
travels to base directory of project,

return:
- 0 if successful, -1 otherwise.
*/
int NodeFileReader::TravelToBaseDir() {
  // get cwd as vector of tokens
  char cwd[256];
  char* buff = getcwd(cwd, sizeof(cwd));
  string b(buff);
  vector<string> x = SplitStringToVector(b, "/");

  // iterate through and find base directory
  int q = x.size() -1;
  int r = -1,y = 0,i = q;
  for (i = q; i > -1; i--, y++) {
    if (x[i] == "NVMB") {
      r = 0;
      break;
    }
  };

  assert (r != -1);

  // traverse back y directories to be at base directory
  string bts = MakeBackTraversalString(y);
  int success = chdir(bts.c_str());

  // set classvar to previous working directory
  if (success != -1) {
    i++;
    pathToPrevCWD = "";
    for (int j = i; j < x.size();j++) {
      pathToPrevCWD += x[j] + "/";
    }
  }
  return success;
}
