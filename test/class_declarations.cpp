#include "class_declarations.hpp"
using namespace std;

NodeProcessingUnit* BlankNPU() {

  NodeStrategos* ns = NodeStrategos::MakeOne(map<string, string>()); 
  NodeProcessingUnit* npu = new NodeProcessingUnit(0, set<int>(), ns);
  return npu;
}
