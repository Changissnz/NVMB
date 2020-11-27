

#ifndef NETWORKGEN_HPP
#define NETWORKGEN_HPP 

#include "network.hpp" 
#include <math.h>

vector<pair<int,int>> AllEdgesInRange(int startNodeId, int endNodeId); 
map<int,vector<int>> GenerateRandomNetworkTemplate(int numberOfNodes, float connectivity); 
void GenerateNetworkActivity(Network N); 

#endif