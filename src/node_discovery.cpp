#include "nvmb_node.hpp"

/*
// description
Outputs a value that determines if Node is required to run a Discovery process.
*/
float NVMBNode::DiscoveryRequirementMeasure() {
   /// TODO: below needs inspection. 
  return npu->GetRoutingScore();
}

/*
// description
Determines if discovery process is running.
*/
bool NVMBNode::IsDiscoveryProcessRunning() {
  return npu->discoveryOn;
}

/*
// description
makes a discovery plan and checks
*/
void NVMBNode::DiscoveryReview(int verbose) {
    // case: discovery already running. done. 
    if (IsDiscoveryProcessRunning()) {
        return; 
    } 

    float score = DiscoveryRequirementMeasure();
    NodeStrategos* ns = GetStrategy();
    Plan* p = ns->MakeDiscoveryPlan(score);
    if (verbose == 2) {cout << "*** Discovery Review for node " << identifier << endl;};
    ns->InsertFuturePlan(p, verbose);
}

void NVMBNode::SetTimestamp(int t) {
    timestamp = t;
    bank->SetTimestamp(t);
}

void NVMBNode::InitializeDiscovery() {
    discoveryCost = -(bank->currency * DISCOVERY_TAX);
    npu->ActivateDF();
}

/*
// description
Adds a neighbor's identifier to instance's neighbors.

// output
True if add False otherwise
*/
bool NVMBNode::AddNeighbor(int neighborId) {
    std::pair<std::set<int>::iterator,bool> ret;   
    ret = neighborIdentifiers.insert(neighborId);
    npu->UpdateNeighborIdentifers(neighborIdentifiers);
    return ret.second;
}

/*
// description
Deletes a neighbor by identifier
*/
void NVMBNode::DeleteNeighbor(int neighborId) {
  neighborIdentifiers.erase(neighborId);
  npu->UpdateNeighborIdentifers(neighborIdentifiers);
}

/*
// description
Determines if neighbor exists
*/
bool NVMBNode::NeighborExists(int neighborId) {
    set<int>::iterator it = neighborIdentifiers.find(neighborId);
    if (it == neighborIdentifiers.end()) {
        return false;
    }
    return true;
}

/*
// description
*/
set<int> NVMBNode::GetNeighbors() {
    return neighborIdentifiers;
}

/*
// description
Finalizes discovery by filtering flares for best paths.
Updates best paths for classvar<strategy>
Logs timestamp
*/
void NVMBNode::FinalizeDiscovery() {
    /// TODO: filter function needs refinement
    npu->FilterFlares(true);

    // log a "discovery finished" timestamp for bank
    TimestampUnit* tsu = new TimestampUnit("discovery", timestamp, discoveryCost);
    bank->LogTimestampUnit(tsu);
    bank->UpdateCurrency(discoveryCost); 
    auto tsu2 = TimestampSummary(); 

    // set discoveryOn to off
    NodeStrategos* strategy = GetStrategy();
    strategy->discoveryOn = false;
    discoveryCost = 0.0;
}

/*
// description 
logs all relevant data for node at timestamp 
*/ 
TimestampUnit* NVMBNode::TimestampSummary() {

    map<string,string> otherDetails; 

    // collect currency,neighbors,contracts 
    BankUnit* bu = GetBank(); 
    otherDetails["currency"] = to_string(bu->currency);

    string neighborsString = NumIterToStringType(GetNeighbors()); 
    otherDetails["neighbors"] = neighborsString; 

    NodeStrategos* ns = GetStrategy(); 

    string contractString = ""; 
    for (auto x: ns->activeContracts) {
        contractString += to_string(x->receiver) + " "; 
    }
    otherDetails["contracts"] = contractString; 

    otherDetails["transmission"] = to_string(ZeroDiv(transmitActual, transmitAttempted, -1.0));
    otherDetails["cfAcceptanceRate"] = to_string(ZeroDiv(cfAccepted, cfEncountered, -1.0)); 

    // collect node nature attributes
    otherDetails["competition"] = to_string(ns->competition); 
    otherDetails["greed"] = to_string(ns->greed); 
    otherDetails["negotiation"] = to_string(ns->negotiation); 
    otherDetails["growth"] = to_string(ns->growth);

    /// TODO: add node strategy functions here! 

    // add nodeRank 
    otherDetails["nodeRank"] = to_string(GetNodeRank());  
    otherDetails["knownNodes"] = NumIterToStringType(knownNodes); 

    // add competitor info 
    otherDetails["competitors"] = NumIterToStringType(ns->competitors); 
    
    // add best paths info
    otherDetails["numPathsRatio"] = to_string(npu->GetBestPathsToKnownRatio()); 
    otherDetails["pathsInfo"] = npu->PathInfoToString(); 

    /// TODO: &^ strategy functions
    otherDetails["competitionMeasure"] = GetStrategy()->competitionMeasure; 
    otherDetails["contractMeasure"] = GetStrategy()->contractMeasure; 
    otherDetails["newContractMeasure"] = GetStrategy()->newContractMeasure; 
    otherDetails["bondAdvantageMeasure"] = GetStrategy()->bondAdvantageMeasure; 
    otherDetails["bondDeletionAdvantageMeasure"] = GetStrategy()->bondDeletionAdvantageMeasure; 

    TimestampUnit* tsu = new TimestampUnit("summary", timestamp, otherDetails); 
    return tsu; 
}

void NVMBNode::LogTimestampAttribute() {
    TimestampUnit* tsu = TimestampSummary();
    bank->LogTimestampUnit(tsu); 
} 