#include "node_file_reader.hpp"
#include <gtest/gtest.h>
using namespace std;


/// TODO: below methods need to have assertions
/*
TMP test to be deleted.
*/
TEST(TmpGetCwd, DisplayCwd) {
  char cwd[256];
  char* buff = getcwd(cwd, sizeof(cwd));
  std::cout << "[DisplayCwd] Current path is " << buff << '\n';
  string b(buff);

  ASSERT_EQ(1,1);
}


TEST(TmpChangeCWD, DisplayNewCwd) {
    string newDir = "../";
    int rc = chdir(newDir.c_str());

    // display new directory
    char cwd[256];
    char* buff = getcwd(cwd, sizeof(cwd));
    std::cout << "[DisplayNewCwd] Current path is " << buff << '\n';
    //---
    ASSERT_EQ(1,1);
}

TEST(SplitStringToVectorTest, ProperTokens) {
  string y = "dir1/dir2/dir234/text.txt";
  cout << "[ProperTokens] " << "FIRST" << endl;
  vector<string> y2 = SplitStringToVector(y, "/");
  for (int i = 0; i < y2.size(); i++) {
    std::cout << "[ProperTokens] " << y2[i] << endl;
  }
  ASSERT_EQ(4, y2.size());
}

TEST(TravelToBaseDirTest, AtBaseDir) {
  char cwd1[256];
  char* buff1 = getcwd(cwd1, sizeof(cwd1));
  std::cout << "[AtBaseDir] Initial path:\t" << buff1 << '\n';

  NodeFileReader nfr("sample");
  nfr.TravelToBaseDir();
  char cwd[256];
  char* buff = getcwd(cwd, sizeof(cwd));
  std::cout << "[AtBaseDir] Current path:\t" << buff << '\n';
  std::cout << "[AtBaseDir] Previous working directory:\t" << nfr.getPathToPrevCWD() << '\n';

  nfr.TravelToPreviousCWD();
  char cwd2[256];
  char* buff2 = getcwd(cwd, sizeof(cwd));
  std::cout << "[AtBaseDir] go back to initial path:\t" << buff2 << std::endl;
}

TEST(InitializeFileTest, CreateNodeFileSuccess) {
  NodeFileReader nfr("sample");// = new NodeFileReader("sample");

  /*
  nfr.InitializeFile();
  vector<int> x = nfr.GetDelimiterPositions();
  std::cout << "[CreateNodeFileSuccess] " << x.size() << endl;
  for (int i =0; i < 5;i++) {
    std::cout << "[CreateNodeFileSuccess] " << x[i] << endl;
  }
  */
}

// TODO GetDelimiterPositions test here

TEST(NodeFileReaderDeclareTest, FileExistence)
    {
    NodeFileReader nfr("sample");// = new NodeFileReader("sample");

    string str = "../../src/bank_history.cpp";
    //string str = "network1.h";
    //str = "gavn";
    //str = "node_test.cpp";
    //const char * c = str.c_str();

    cout << "THERE " << FileExists(str.c_str()) << endl;

    string fp = "../src/node_file_reader_test.cpp"; 
    /*
    char * fpc = fp.c_str(); 
    ASSERT_EQ(1, FileExists(fpc));// FileExists(StringToCharArray(nfr.getFilePath())));
    */ 

    char *cstr = new char[fp.length() + 1];
    strcpy(cstr, fp.c_str()); 

    /*
    cout << "HERE "
    ASSERT_EQ(1, FileExists(fp.c_str()));// FileExists(StringToCharArray(nfr.getFilePath())));
    */ 

}


/// TODO: merge this with above and make assertion 
TEST(NodeFileReaderSpecTest, SpecCorrect) {
    auto nfr = NodeFileReader("MARKOS", "nolo_1"); 

    // data dir creation works 
    // 
    nfr.InitializeFile(); 
    // ASSERT THE FILE EXISTS HERE 
}