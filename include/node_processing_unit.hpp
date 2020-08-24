/*
The NodeProcessingUnit is responsible for activities such as
- processing an array of Discoveryflares for shortest paths
- receiving and sending propositions
- forming/maintaining contracts
- proposition/contract interactions

The NodeProcessingUnit has a NodeStrategos that aids it in its responsibilities.
*/

#ifndef NODEPROCESSINGUNIT_HPP

#define NODEPROCESSINGUNIT_HPP
#define MAX_GROWTH_RATE 3.0
#define COMMERCE_TAX 0.05
#define NUM_BEST_PATHS 3

#include "discovery_flare.h"
#include "comm_flare.hpp"
#include "node_strategos.hpp" // TODO incorporate node_strategos.hpp
#include "timestamp_unit.hpp"
#include <algorithm>
#include <math.h>
#include <numeric>

extern int MAX_FLARES; // MAX_FLARES per 
extern int MAX_FLARE_MULTIPLIER; // 

class NodeProcessingUnit
{
private:
    int nodeOwnerIdentifier;
    std::set<int> neighborIdentifiers;

    // these flares are used for discovery
    std::set<DiscoveryFlare*> ownedFlares;
    std::set<DiscoveryFlare*> ownedFlaresCache;
    std::set<DiscoveryFlare*> heldFlares;

    // these flares are used for commerce
    std::set<CommFlare*> ownedCFlares;
    std::set<CommFlare*> heldCFlares;
    std::set<CommFlare*> processedCFlares;

    std::set<CommFlare*> finishedCFlares;

    /// TODO: ?!? delete below  
    std::set<int> loggedPIDs; 

    // used to route flares
    std::map<int, std::vector<std::pair<std::vector<int>, float>>> bestPaths;
    std::map<int, std::vector<std::pair<std::vector<int>, float>>> bestPathsCache;

    bool finished;
    NodeStrategos* strategy;

public:
    float worstPathScore;
    int maxPathLength; // this keeps track of the longest node-to-node path
    bool discoveryOn;
    /// TODO: some discoveryType/s have not been coded yet. 
    std::string discoveryType; // flower | elimination | limit 

    NodeProcessingUnit(int noi, std::set<int> ni, NodeStrategos* s) {
        nodeOwnerIdentifier = noi;
        neighborIdentifiers = ni;
        strategy = s;
        finished = false;
        discoveryOn = false;
        discoveryType = "flower";
    }

    void LoadInitialPaths(); 
    int GetNOI();
    std::set<int> GetNeighborIdentifiers();
    std::map<int, std::vector<std::pair<std::vector<int>, float>>> GetBestPaths() {
        return bestPaths;
    }

    int GetNumberOfBestPaths() {
        int x = std::accumulate(bestPaths.begin(), bestPaths.end(),
            0.0,
            [](int t, std::pair<int, std::vector<std::pair<std::vector<int>, float>>> r)
            {
            return t + r.second.size();
            });
        return x;
    }

    void DisplayBestPaths();

    NodeStrategos* GetStrategy() {
        return strategy;
    }

    void UpdateNeighborIdentifers(std::set<int> x);

    // discovery flare methods
    bool ActivateDF();

    bool AddOwnedDF(DiscoveryFlare* df);
    bool AddHeldDF(DiscoveryFlare* df);
    void RemoveHeldDF(DiscoveryFlare* df);
    bool Process1DF(DiscoveryFlare* df);

    void ProcessDF_();
    bool DiscoveryFinished();
    void ProcessDF();

    std::set<DiscoveryFlare*> GetOwnedFlares();
    std::set<DiscoveryFlare*> GetOwnedFlaresCache();
    void ClearDiscoveryData();

    std::set<DiscoveryFlare*> GetHeldFlares();
    std::set<CommFlare*> GetOwnedCFlares() {
        return ownedCFlares;
    }

    std::set<CommFlare*> GetProcessedCFlares() {
        return processedCFlares;
    }

    std::set<CommFlare*> GetHeldCFlares() {
        return heldCFlares;
    }

    void AddOwnedCF(CommFlare* cf) {
        ownedCFlares.insert(cf);
    }

    void AddHeldCF(CommFlare *cf) {
        heldCFlares.insert(cf);
    }

    void RemoveHeldCF(CommFlare* cf) {
        heldCFlares.erase(cf);
    }

    void AddProcessedCF(CommFlare* x) {
        processedCFlares.insert(x);
    }

    void RemoveProcessedCF(CommFlare* x) {
        processedCFlares.erase(x);
    }

    void AddFinishedCF(CommFlare* cf) {
        finishedCFlares.insert(cf);
    }

  /// TODO: error-check this. 
  /*
  cfType - held|owned|processed

  // return 
  [0] CommFlare to remove 
  [1]   [0] erase a complete flare  
        [1] unique flare 
  */ 
    pair<CommFlare*, bool> UpdateOneCompleteFrom(std::string cfType) {

        assert (cfType == "held" || cfType == "owned" || cfType == "processed"); 
        std::set<CommFlare*> scf;
        
        if (cfType == "held") {
            scf = GetHeldCFlares(); 
        } else if (cfType == "owned") {
            scf = GetOwnedCFlares(); 
        } else {
            scf = GetProcessedCFlares(); 
        }
        
        int finishedCFPlanId = -1;
        bool erase = false;
        CommFlare* erasedCF; 

        for (auto it = scf.begin(); it != scf.end(); it++) {
            if ((*it)->finished) {
                erasedCF = *it;  
                erase = true; 
                scf.erase(it); 
                finishedCFPlanId = (*it)->planIdentifier; 

                if (loggedPIDs.find(finishedCFPlanId) == loggedPIDs.end()) {
                    loggedPIDs.insert(finishedCFPlanId);
                }
                break; 
            } 
        }

        if (cfType == "held") {
            heldCFlares = scf; 
        } else if (cfType == "owned") {
            ownedCFlares = scf; 
        } else {
            processedCFlares = scf; 
        }

        if (erase) {
            strategy->FinishExecutedPlan("plan", finishedCFPlanId); 
        }

        return make_pair(erasedCF, erase); 
    }

    void DisplayCF(std::string flareLabel);
    void Update();
    int UpdateHeldFlaresByIteration(string updateType, int index);
    void UpdateHeldFlares(string updateType);
    void UpdateOwnedFlares();
    bool IsFinished();

        // to get best paths based on discovery
    bool PathExists(std::vector<int> path, int key);
    void UpdateMaxPathLength(int pathLength);
    bool AddBetterPath(int key, std::vector<int> path,
        float score);

    bool EliminatePath(int key, vector<int> path, float score);
    bool IsCycle(std::vector<int> path);
    int ConsiderPath(DiscoveryFlare* df, bool excludeCycle);

    int FilterOneFlare(bool excludeCycle, int index);

    void EliminateDiscoveryFlares(); 
    void UpdateBestPathsWithCache();
    void EliminateHeldDF();
    void EliminateFinishedFlares(string dfType);
    void FilterFlares(bool excludeCycle);
    void UpdateWorstPathScore();

    // commerce flare methods
    CommFlare* ProduceCommFlare(Plan* p);
    std::pair<std::vector<int>, bool> GetPathForDestination(int destNode);

    // contract functions
    Contract* GetExistingContract(int nodeId);
    void InitiateContract(CommFlare* cf);
    void BreakContract(CommFlare* cf); 

    // timestamp processing functions 
    std::map<std::string,std::map<std::string,int>> ProcessTimestampSequenceForFrequency(
        std::vector<TimestampUnit*> tsuseq, std::vector<std::string> dataTemplate); 
    std::vector<std::map<std::string,std::string>> GatherTimestampData(std::vector<TimestampUnit*> tsuseq, std::vector<std::string> dataTemplate); 
    
    // file : npu_path_functions
    std::set<int> GetKnownNodes();

    int GetKnownNumberOfNodes() {
        int x = bestPaths.size();
        return x;
    }

    float GetBestPathsToKnownRatio() {
        int numPaths = GetNumberOfBestPaths(); 
        int maxPossible = GetKnownNumberOfNodes() * NUM_BEST_PATHS; 
        return ZeroDiv(float(numPaths), float(maxPossible), 0.0); 
    }

    std::string PathInfoToString() {
        std::vector<float> s = vector<float>();  
        for (auto c: bestPaths) {
            s.push_back(c.first);
            s.push_back(float(c.second.size()) / NUM_BEST_PATHS); 

           float summation = std::accumulate(c.second.begin(), c.second.end(), 0.0,
                            [](float t, std::pair<std::vector<int>, float> r)
                            {
                                return t + r.second; 
                            }); 

            s.push_back(ZeroDiv(summation, c.second.size(), 0.0)); 
        }

        return NumIterToStringType(s);         
    }

    float EdgeScore(float edgeRisk, float competition, float greed);
    float GetPathScore(std::vector<int> path);
    std::pair<float,float> GetPathRiskPotentialScores(std::vector<int> path);

    void UpdatePathScores(int nodeId);
    float GetPathVolatility(std::vector<int> pathInfo);
    float GetTaxationOfPath(int pathLength);
    float CalculateContractWorthMeasure(int otherNode);
    float GetRoutingScore();

    // functions used to gauge worth of creating new bond
    std::vector<int> ShortenPath(std::vector<int> path, std::pair<int,int> newBond);
    std::vector<std::pair<int, std::pair<float,float>>> BestPathsAffectedByNewBondScores(std::pair<int,int> newBond);
    std::vector<int> GetSkippedNodes(std::vector<int> path, std::pair<int,int> newBond);
    std::vector<std::pair<int, std::vector<int>>> BestPathsAffectedByNewBondSkippedNodes(std::pair<int,int> newBond);
    std::pair<int,int> FindBondIndexInPath(std::vector<int> path, std::pair<int,int> bond);
    bool BondExistsInPath(std::vector<int> path, std::pair<int, int> bond);
    std::pair<bool, std::pair<float, float>> GetReplacementPathInfo(int key, std::pair<int,int> oldBond, float worstPathScore);
    std::vector<std::pair<int, std::pair<float, float>>> BestPathsAffectedByBondBreakScores(std::pair<int,int> oldBond, float worstPathScore);
    std::pair<bool, std::pair<std::vector<int>, std::vector<int>>> SplitPathByBond(std::vector<int> path, std::pair<int,int> bond);
    std::vector<std::pair<int, std::pair<std::vector<int>, std::vector<int>>>> BestPathsAffectedByBondBreakSkippedNodes(std::pair<int,int> oldBond);
    std::vector<std::pair<float,float>> BestPathRiskInfoForNewContract(int node);
    std::vector<int> BestPathLengths(int node);

    std::map<int, float> GetBestPathScores();
    int RemoveRecordsOfEdge(std::string edge); 
};

#endif
