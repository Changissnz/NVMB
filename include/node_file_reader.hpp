#ifndef NODEFILEREADER_HPP
#define NODEFILEREADER_HPP 
#define TIMESTAMP_SIZE_THRESHOLD 1000

#include <sys/stat.h>
#include <string>
#include <fstream>
#include <vector>

int FileExists(const char *path);
int DirExists(char *path);
char* StringToCharArray(std::string s);
std::string CharArrayToString(char* c);
int IsCorrectFileNameFormat(std::string fn);
std::vector<std::string> SplitStringToVector(std::string s, const char* delimiter);

extern std::vector<std::string> TIMESTAMP_HEADER; 

class NodeFileReader
{
private:
    std::string datDir;
    std::string filename;
    std::string filePath;
    std::string tableOfContentsPath;

    std::string pathToPrevCWD; // path to previous working directory
    std::ofstream fileObj; 

public:
    NodeFileReader(std::string fn);

    NodeFileReader(std::string dd, std::string fn);

    void MakeFrame();
    std::string getFilePath() {return filePath;};
    std::vector<int> GetDelimiterPositions();
    std::string getPathToPrevCWD() {return pathToPrevCWD;};

    /// TODO: 
    void WriteHeaderRow(std::vector<std::string> row);
    void WriteRowData(std::string row); 
    void InitializeFile();
    void MakeDataDirectory();
    void CloseWrite();  

    int TravelToBaseDir();
    int TravelToPreviousCWD();

    int UpdateNeighborData(std::string newNeighborData);
    int AppendBankData(std::string newBankData);
    int AppendContractData(std::string newContractData);
    int AppendPropRecv(std::string propRecvData);
    int AppendPropSent(std::string propSentData);
};

#endif
