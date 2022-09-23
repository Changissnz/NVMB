#ifndef NETWORK1_H
#define NETWORK1_H

#include "network_generation.hpp"
using namespace std;

void PrintSpacer(); 
Network Network1();
Network Network2();

Network Network3(string folderName); 
Network Network4(string folderName, bool mortality); 
Network Network5(string folderName, bool mortality); 
Network Network6(string folderName, bool mortality); 

void RunNetworkDFSProcessOnNode(Network* N, int nodeId);
void RunNetworkDFSProcessOnAllNodes(Network N); 
void RunNetworkDFSProcessOnActiveNodes(Network* N); 

/// TODO: below method to delete, use `SetNetworkNodesActiveStatus`
Network Network1WithNodesInSelfMode(vector<int> x);
void SetNetworkNodesActiveStatus(Network* N, set<int> x, string status);
void DisplayNodeBestPaths(NVMBNode* n);

void DisplayNodePlansInfo(Network N, int nodeId); 
void DisplayNetworkTimestampUnits(Network N);  
void LoadPlanForDirectTesting(NVMBNode* n, Plan* p); 

// makes a contract plan  
void MakeContractPlanDefaultForNode(NVMBNode* n, int receiverNode); 
Plan* MakeContractPhantomPlanDefaultForNode(NVMBNode* n, int senderNode, int receiverNode); 

// break contract plan 
void BreakContractPlanDefaultForNode(NVMBNode* n , int contractedNode); 
void BreakContractPhantomPlanDefaultForNode(NVMBNode* n, NVMBNode* n2, int receiverNode);

void LoadPlanForNodeN(Network N, int n);  

bool IsNetworkDiscoveryOff(Network N);
/// TODO: below will replace above. 
bool CheckNetworkDiscoveryOff(Network* N); 

        // performs basic checking of timestamp log for node, 
        // which includes summary, discovery, and 
bool CheckNodeTimestampDataForEvent(NVMBNode* n, int timestamp, string eventType);
int GetEventCountForNode(NVMBNode* n, int timestamp, string eventType); 
map<int,int> GetEventCountPerNodeForTimestamp(Network* N, int timestamp, string eventType); 

TimestampUnit* FetchNodeTimestampDataForEvent(NVMBNode* n, int timestamp, string eventType); 
vector<int> GetNodesInPhantomPath(Network* N, vector<int> checkPoints);  

#endif