#include "extras.hpp"
using namespace std;

string NodePairToString(int v1, int v2) {
  char buffer[100];

  if (v1 < v2) {
    sprintf(buffer, "%d_%d", v1, v2);
  } else {
    sprintf(buffer, "%d_%d", v2, v1);
  }
  string b(buffer);
  return b;
}

/*
// description
outputs -1 if string is invalid integer
*/
int StringToInt(string s) {
  try {
    int q = stoi(s, 0, 10);
    return q;
  } catch (exception& e) {
    return -1;
  }
}

pair<int,int> StringToNodePair(string s) {

  // find the two string node identifiers
  size_t found = s.find("_");
  string s1 = s.substr(0, found);
  string s2 = s.substr(found + 1, s.length());

  // convert them to integers
  int one = StringToInt(s1);
  int two = StringToInt(s2);
  assert (one != -1 && two != -1);
  return make_pair(one,two);
}

float ZeroDiv(float num, float denum, float defaultVal) {
  if (denum == 0.0) {
    return defaultVal;
  }
  return num / denum;
}

float RandomFloat() {
    return ((float) rand()) / (float) RAND_MAX;
}