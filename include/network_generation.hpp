

#ifndef NETWORKGEN_HPP
#define NETWORKGEN_HPP 

#include "network.hpp" 
#include <math.h>

vector<pair<int,int>> AllEdgesInRange(int startNodeId, int endNodeId); 
map<int,vector<int>> GenerateRandomNetworkTemplate(int numberOfNodes, float connectivity); 
Network GenerateRandomNetwork(int nid, bool mortality, string folderPath, int numberOfNodes, float connectivity, string nodePrefix); 
void GenerateNetworkActivity(Network N); 

#endif