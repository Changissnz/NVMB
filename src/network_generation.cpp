/*
this file will be used for generating networks. 
*/
#include "network_generation.hpp" 
using namespace std; 

vector<pair<int,int>> AllEdgesInRange(int startNodeId, int endNodeId) {
    auto output = vector<pair<int,int>>(); 
    
    for (int i = startNodeId; i < endNodeId - 1; i++) {
        for (int j = i + 1; j < endNodeId; j++) {
            auto p = make_pair(i,j); 
            output.push_back(p); 
        }
    }

    return output; 
};

/// CAUTION: to use this method for large graphs `AllEdgesInRange` needs to be modified.  
map<int,vector<int>> GenerateRandomNetworkTemplate(int numberOfNodes, float connectivity) {
    // initialize map of nodes 
    auto x = map<int,vector<int>>(); 
    for (int i = 0; i < numberOfNodes; i++) {
        x[i] = vector<int>();  
    }

    // get all possible edges 
    auto allEdges = AllEdgesInRange(0, numberOfNodes); 

    // choose y possible edges
    float y = float(ceil(double(float(allEdges.size()) * connectivity)));
    
    for (int y_ = 0; y_ < y; y_++) {
        auto index = rand() % allEdges.size(); /// RandomIndex(0, allEdges.size());
        auto edge = allEdges[index];
        allEdges.erase(allEdges.begin()+index); 
        
        x[edge.first].push_back(edge.second);
        x[edge.second].push_back(edge.first); 
    }

    return x; 
}; 

Network GenerateRandomNetwork(int nid, bool mortality, string folderPath, int numberOfNodes, float connectivity, string nodePrefix) {
    Network N(nid, mortality);
    auto networkTemplate = GenerateRandomNetworkTemplate(numberOfNodes, connectivity);

    for (auto x: networkTemplate) {
        NVMBNode* n0 = new NVMBNode(x.first, folderPath, nodePrefix + "_" +  to_string(x.first), MakeRandomNodeNatureVariables()); 
        for (int y: x.second) {
            n0->AddNeighbor(y); 
        };

        N.AddNode(n0); 
    }

    return N; 
} 