#ifndef NVMB_NODE_HPP
#define NVMB_NODE_HPP
#define DISCOVERY_TAX 0.2
#define NUM_THREADS 

#include <assert.h>
#include <stdexcept>
#include "node_file_reader.hpp"
#include "node_processing_unit.hpp"
#include "bank_unit.hpp"
#include "node_competitor_functions.hpp"
#include "contract_review_functions.hpp"
#include "bond_functions.hpp"
#include "new_contract_functions.hpp"
using namespace std;

class NVMBNode
{
private:
    int identifier;
    NodeFileReader* reader;
    BankUnit* bank;

    set<int> neighborIdentifiers;
    // specifies if node is active 
    string mode;
    
    set<int> knownNodes; 

    /*
    options for testModeSpec: 
    make_bond|break_bond|make_contract|break_contract|discovery 
    */ 
    string testModeSpec; 
    bool manInput; 

    bool initialized;
    int timestamp;

    // variables for Node actions
    /// TODO : refactor and delete this variable
    float discoveryCost;

    // if [0] is true, then default response is [1] 
    pair<bool, bool> defaultResponseCF;

    int transmitActual; 
    int transmitAttempted;

    int cfAccepted; 
    int cfEncountered; 

public:
    NodeProcessingUnit* npu;

    NVMBNode(int id, string fn); 

    NVMBNode(int id, string fn, string m); 

    NVMBNode(int id, string fn, map<string,string> m); 

    NVMBNode(int id, string dd, string fn, map<string,string> m); 

    ///////// getter and setter header functions
    int GetId() {return identifier;};

    string GetMode() {return mode;};

    BankUnit* GetBank() {
      return bank;
    }

    float GetWorth() {
      return bank->currency;
    }

    NodeStrategos* GetStrategy() {
      return npu->GetStrategy();
    }

    int GetTimestamp() {
      return timestamp;
    }

    void SetMode(string m) {
      assert (m == "active" || m == "inactive");
      mode = m;
    }

    void SetTestModeSpec(string m) { 
      assert (m == "" || m == "discovery" || m == "make_bond" || m == "break_bond" || m == "make_contract" || m == "break_contract"); 
      testModeSpec = m; 
    }

    void SetDiscoveryType(string m) {
        assert (m == "flower" || m == "elimination" || m == "limit"); 
        npu->discoveryType = m; 
    }

    void SetToManualInput(bool manInp) {
      manInput = manInp; 
    }

    void SetDefaultResponseCF(bool defaultOn, bool def) {
      defaultResponseCF.first = defaultOn;
      defaultResponseCF.second = def;
    }


    void IncrementTransmission(bool success) {
        if (success) {
            transmitActual++;
        } 
        transmitAttempted++; 
    }

    void IncrementResponse(bool accept) {
        if (accept) {
            cfAccepted++; 
        } 
        cfEncountered++;         
    }

    void ClearTransmissionAndResponse() {
        transmitActual = 0; 
        transmitAttempted = 0;

        cfAccepted = 0; 
        cfEncountered = 0; 
    }

    pair<int,int> GetTransmission() {
        return make_pair(transmitActual, transmitAttempted); 
    }

    ///////// file node_discovery
    float DiscoveryRequirementMeasure();
    bool IsDiscoveryProcessRunning();
    void DiscoveryReview(int verbose);
    void SetTimestamp(int t);
    void InitializeDiscovery();

    // neighbor methods
    bool AddNeighbor(int neighborId);
    void DeleteNeighbor(int neighborId);
    bool NeighborExists(int neighborId);
    int GetNumberOfNeighbors() {return neighborIdentifiers.size();};
    set<int> GetNeighbors();
    void FinalizeDiscovery();

    TimestampUnit* TimestampSummary(); 
    void LogTimestampAttribute(); 

    ///////////////// methods for `node_planner`
    ////////////////    file will contain functions that handle a Node's decision (set of functions for plans)
    float ConsiderExistingContract(Contract* C);
    bool ConsiderCompetitor(int i);
    float ConsiderNewContract(int newNode);
    set<string> SelectPossibleBondsToAdd();
    set<string> SelectPossibleDirectBonds();
    string SelectBondFromBestPath(string bondType, int recursiveLimit);

    float MeasureBondCreation(pair<int,int> newBond, set<int> recognizedCompetitors);
    set<string> SelectPossibleBondsToDelete();
    float MeasureBondDeletion(pair<int,int> oldBond, set<int> recognizedCompetitors);

    void UpdateConductedContract(int contractReceiver) {
      NodeStrategos* ns = GetStrategy();
      (ns->contractsConducted).insert(contractReceiver);
    }

    void ClearConductedContracts() {
      NodeStrategos* ns = GetStrategy();
      ns->contractsConducted = std::set<int>();
    }

    bool HasConductedContractWith(int contractReceiver) {
      NodeStrategos* ns = GetStrategy();
      return (ns->contractsConducted).find(contractReceiver) == (ns->contractsConducted).end();
    }

    vector<TimestampUnit*> FetchBankHistoryInRange(int minumum, int maximum) {
      return bank->GetTimestampUnitHistoryAtRange(minumum, maximum); 
    } 

    void ShutDown(); 

    //////////////// methods for `node_strategy_bridge`

    bool ProcessOne(int verbose);
    void ContinueDiscovery();
    float GetNodeRank() {
      NodeStrategos* ns = npu->GetStrategy();
      return ns->GetRank();
    }
    void UpdateNodeRank(float r);
    int ChooseRandomNeighbor(); 
    void ScrambleNature(); 
    void ExecuteGivenPlan(Plan* p, int verbose);
    void MarkDiscoveryFinished();
    void ExecutePlan(int verbose);
    void CalculatePlansTestMode(int verbose); 
    void CalculatePlans(int verbose);
    set<int> GetKnownNodes();
    set<int> FindCompetitors();
    void NewBondReview(set<int> recognizedCompetitors, int verbose);
    void OldBondReview(set<int> recognizedCompetitors, int verbose);
    void ExistingContractReview(int verbose);
    set<int> GetNodesWithoutContract();
    void NewContractReview(int verbose);
    bool IsContractedWith(int nodeIdentifier);
    void UpdateContract(int nodeIdentifier, bool transmissionSuccess, float delta, int pathDistance); 
    void UpdateContract(int nodeIdentifier, bool alve);

    vector<int> GetPlanCounts();

    vector<int> GetFlareCounts();  
    vector<int> GetCommFlareCounts(); 
    void RespondToCommFlares(bool toLog, bool verbose);
    bool RespondToCommFlare(CommFlare *cf, bool writeResp, bool verbose); 
    void LogResponse(CommFlare* cf);

    void UpdateValue(float c);
    void DisplayActiveContracts();

    void UpdateFinishedCF(); 
 
    void WriteOutToFile(bool threshold = true); 
    void SummarizeLastNTimestamps(int n); 
};

#endif
