#ifndef NETWORK_HPP
#define NETWORK_HPP

// set below to float in [0, 1] to set max connectivity of network
#define MAX_CONNECTIVITY -1.0
#define TRANSPORT_TAX 0.02

#include "nvmb_node.hpp"

float AccumulateContractRisk(float risk, Contract* c);
map<string,string> MakeRandomNodeNatureVariables();


class Network
{

private:
    int identifier;

public:

    map<int, NVMBNode*> contents;
    vector<int> deceased; 
    bool mortality; 

    // these values are to be transmitted to Nodes for their own strategic analysis.
    map<string, float> edgeRisks; // key is edge, value is float
    map<int, float> nodeRank; // traffic measure
    map<int, pair<int, float>> contractData;
    float averageContractTransmission;
    float averageContractPathLength;

    /// TODO: not yet coded.
    map<int,int> nodeTimestampDeaths; 
    int timestamp;
    int maximallyConnected;

    Network(int i, bool m) {
        contents = map<int, NVMBNode*>(); 
        identifier = i;
        mortality = m; 
        timestamp = 0; 
        maximallyConnected = 0; 
    }

    //////////// file network.cpp
    //// update and run functions
    bool NodeExists(int nodeIdentifier);
    bool AddNode(NVMBNode* n);
    bool AddNode(int i);
    void BreakEdge(int v1, int v2, int directedBreak);
    void AddEdge(int v1, int v2, int directedAdd);
    bool EdgeExists(int v1, int v2, int directed);
    int Size();
    void UpdateNodeData();
    void UpdateTime(); 
    void UpdateNodes(int verbose);
    void ClearNullNodes(); 
    void RunOneTimestamp(int verbose = 0);
    void Run(int verbose = 0, int numRounds = -1);
    void RunRandomNodeNatureVars(int verbose, int numRounds, float switchFrequency);
    void UpdateNodeNaturesAtRandom(float freq); 
    void UpdateNodeNatureAtRandom(int nodeId, float freq);
    void PrerunProcess(); 
    void ShutDown(); 
    void Reopen(); 
    void ActivitySummaryUpdate();
    void WriteOutToFile(bool threshold); 
    bool IsAlive();
    vector<int> ActiveNodes(); 

    //////////// file network_flare_functions.cpp
    //// discovery flare
    void AssignHeldDF(NVMBNode* n);
    void AssignHeldDFAll();
    void UpdateNodeDFOwned();
    void UpdateNodeDFNeighbors();

    //// commerce flare
    bool IncrementCFlare(CommFlare* cf);
    void MoveCFlare(CommFlare* cf);
    void UpdateCFlare(CommFlare* cf, int nodeSource, int verbose);
    void TransferCFlare(CommFlare* cf, int nodeSource);
    void ProcessNodeCF(NVMBNode *n, int verbose);
    void SetHeldCFToActive();
    float IncrementNodeBankByCF(CommFlare *cf);
    void UpdateCompleteCFlares(); 
    void MisrouteUpdate(CommFlare* cf, int nodeSource); 
    void ProcessNodeCFAll_(bool logResp, int verbose);
    void ProcessNodeCFAll(int verbose);
    void MarkExecutedPlansAsFinished();
    void ProcessNodeCFResponses(NodeProcessingUnit* npu);

    void ModifyNetworkByCFBond(CommFlare* cf);
    void ModifyNetworkByCFContract(CommFlare* cf);

    void SetMaxConnectivity();
    bool SetBondModCosts(CommFlare *cf);

    void DisplayExistingContracts();
    void ActiveContractExecution(); 
    bool ConductContract(Contract* C, int nodeSource); 
    pair<pair<int,int>,bool> RouteContract(Contract* C, vector<int> path, int nodeSource); 


    //////////// file : network_measure_functions.cpp
    // calculate the edge risks
    void CountPertinentEdges_(NVMBNode* n);
    void CountPertinentEdges();
    void NormalizeMap(std::string mapLabel);
    void CalculateEdgeRisks();
    void TransferEdgeRiskDataToNodes();
    void TransferRankToNodes();
    /* calculates average trnsmission data, contract path distance,
      and average contract risk
    */
    void CalculateContractData();
    void TransferContractInfoToNodes();
    void NetworkDataUpdateToNodes();
    void ReplyToInfoRequest(int NodeId);

    std::set<CommFlare*> RetrieveAllCFOfType(std::string cfType); 
    std::pair<std::vector<int>, bool> GetValidPathFromSourceToTarget(int sourceNode, int targetNode);
    bool IsPathValid(std::vector<int> path);
};

#endif
