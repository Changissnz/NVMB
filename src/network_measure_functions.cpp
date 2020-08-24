#include "algorithm"
#include "network.hpp"

float AccumulateContractRisk(float risk, Contract* c) {
  return risk + c->pathRisk;
}

/// TODO: needs error-checking 
/*
// description
iterates through Node's best paths and adds them to pertinent edges.
Two variables are relevant: edgeRisks, nodeRank. 

An edge that is more frequently found in a best path is more risky. 
*/
void Network::CountPertinentEdges_(NVMBNode* n) {

    if (n == nullptr) {return;}

    NodeProcessingUnit* npu = n->npu;
    map<int, vector<pair<vector<int>, float>>> sp = npu->GetBestPaths();
    map<int, vector<pair<vector<int>, float>>>::iterator it;
    string s;
    vector<int> path;

    for (it = sp.begin(); it != sp.end(); it++) {
        // iterate through each path
        auto q = (*it).second;
        for (auto j = q.begin(); j != q.end(); j++) {
            auto z = (*j).first;
            if (z.size() == 0) {
                continue;
            }

            // iterate through path and register edge into edge risks and nodeRank into 
            for (int k = 0; k < z.size() - 1;k++) {
                s = NodePairToString(z[k], z[k+1]);
                edgeRisks[s] = edgeRisks[s] + 1;
                nodeRank[z[k]] += 1;
            }
            nodeRank[z[z.size() -1]] += 1;
        }
    }
}

/*
// description
iterate through shortest paths for each NodeProcessingUnit and count each edge
*/
void Network::CountPertinentEdges() {
  map<int,NVMBNode*>::iterator it;

  for (it = contents.begin(); it != contents.end(); it++) {
    CountPertinentEdges_(it->second);
  }
}

void Network::NormalizeMap(string mapLabel) {
    if (mapLabel == "edge") {
        // get max 
        map<string,float>::iterator best
                = max_element(edgeRisks.begin(), edgeRisks.end(),
                [] (const pair<string,float>& a,
                const pair<string,float>& b)->bool{ return a.second < b.second; } );
        float x = best->second;

        // normalize 
        map<string, float>::iterator it;
        for (it = edgeRisks.begin(); it != edgeRisks.end(); it++) {
            edgeRisks[it->first] = ZeroDiv(edgeRisks[it->first], x, 0.0);
        }
    } else {
        // get max 
        map<int,float>::iterator best
                = max_element(nodeRank.begin(), nodeRank.end(),
                [] (const pair<int,float>& a,
                const pair<int, float>& b)->bool{ return a.second < b.second; } );
        float x = best->second;

        // normalize 
        map<int, float>::iterator it;
        for (it = nodeRank.begin(); it != nodeRank.end(); it++) {
            nodeRank[it->first] = ZeroDiv(nodeRank[it->first], x, 0.0);
        }
    }
}

/*
// description
Counts frequencies of edges and their respective nodes. These frequency values
will be used to derive EdgeRisk and NodeRank scores.

*/
/// CAUTION: no error-handling
void Network::CalculateEdgeRisks() {

  // clear previous data and run count
  edgeRisks = map<string,float>();
  nodeRank = map<int,float>();

  // do frequency count here
  CountPertinentEdges();

  // normalize value
  NormalizeMap("edge");
  NormalizeMap("node");
}

///////////////// MOVE THESE METHODS TO NETWORK-NODE-TRANSFER /////////////////

void Network::TransferEdgeRiskDataToNodes() {
  NodeStrategos::edgeRisks = edgeRisks;
}

void Network::TransferRankToNodes() {
  // transfer node ranks
  for (auto it = nodeRank.begin(); it != nodeRank.end(); it++) {
    /// TODO: added existence check here. 
    if (contents.find(it->first) != contents.end() && contents[it->first] != nullptr) {
        contents[it->first]->UpdateNodeRank(it->second);
    }
  }
}


//// TODO : this needs work!!
/// TODO OPTIMIZATION: instead of re-calculating this data every timestamp, instead
///       keep a running total and keep track of new contracts and recently deleted
///       contracts.
/*
// description

averageContractRisk for each node is determined by the average of the risk
of the best path for each contract.

Network's structure of ContractData is comprised of two parts:
(1) map: nodeId -> pair<number_of_contracts, averageContractRisk>
(2) float - averageContractTransmission (overall)
(3) float - averagePathOfContract (overall)
*/
void Network::CalculateContractData() {
  NVMBNode* n;
  NodeStrategos* ns;
  float risk;
  pair<int, float> p;
  contractData = map<int, pair<int,float>>();

  averageContractTransmission = 0.0;
  averageContractPathLength = 0.0;
  float numContracts = 0;

  for (auto it = contents.begin(); it != contents.end(); it++) {
    // calculate risk

    if (it->second == nullptr) {
        continue; 
    } 

    ns = (it->second)->GetStrategy();
    risk = accumulate(ns->activeContracts.begin(), ns->activeContracts.end(),
          0.0, AccumulateContractRisk);
    risk = ZeroDiv(risk, float(ns->activeContracts.size()), 0.0);

    // make map value
    p = make_pair(ns->activeContracts.size(), risk);
    contractData[it->first] = p;

    // calculate transmission success
    averageContractTransmission = accumulate(ns->activeContracts.begin(), ns->activeContracts.end(),
                  averageContractTransmission,
                  [](float averageContractTransmission, Contract* c) {return averageContractTransmission + c->transmissionRate;}
                );

    // calculate path length
    averageContractPathLength = accumulate(ns->activeContracts.begin(), ns->activeContracts.end(),
                  averageContractPathLength,
                  [](float averageContractPathLength, Contract* c) {return averageContractPathLength + c->path.size();}
                );

    // update number of active contracts
    numContracts += ns->activeContracts.size();
  };

  // update class vars. 
  averageContractTransmission = ZeroDiv(averageContractTransmission, numContracts, 0.0);
  averageContractPathLength = ZeroDiv(averageContractPathLength, numContracts, 0.0);
}

/*
// description
Network's ContractInfo variable provides basic information about each Node's
contract data: # of contracts with other nodes, predicted net gain for each round.
*/
void Network::TransferContractInfoToNodes() {
  NodeStrategos::networkContractInfo = contractData;
  NodeStrategos::averageNodeContractPathDistance = averageContractPathLength;
  NodeStrategos::averageTransmissionMeasure = averageContractTransmission;
}

/*
// description
transfers the following to each node:
- EdgeRisks
- each Node's rank
- # of active contracts of each node
  and their corresponding vector of risks
- net gain of each Node's bank from the previous timestamp
*/
void Network::NetworkDataUpdateToNodes() {
  // EdgeRisk and NodeRank
    CalculateEdgeRisks();
    TransferEdgeRiskDataToNodes();
    TransferRankToNodes();

    // contract data
    CalculateContractData();
    TransferContractInfoToNodes();
}

// TODO: untested
pair<vector<int>, bool> Network::GetValidPathFromSourceToTarget(int sourceNode, int targetNode) {

  if (contents.find(sourceNode) == contents.end()) {
    return make_pair(vector<int>(), false);
  }

  // get the paths from source to target
  auto pathInfo = (contents[sourceNode]->npu)->GetBestPaths()[targetNode];

  for (auto x: pathInfo) {
    if (IsPathValid(x.first)) {
      return make_pair(x.first, true);
    }
  }

  return make_pair(vector<int>(), false);
}

// TODO: untested
bool Network::IsPathValid(vector<int> path) {

  if (path.size() == 0) { return true;};

  if (contents.find(path[0]) == contents.end()) {
    return false;
  }

  NVMBNode* n = contents[path[0]];

  for (int i = 1; i < path.size(); i++) {
    if (!n->NeighborExists(path[i])) {
      return false;
    }
    n = contents[path[i]];
  }
  return true;
}

/*
// description
replies to a Node's request for information
*/
void Network::ReplyToInfoRequest(int NodeId) {
}


set<CommFlare*> Network::RetrieveAllCFOfType(string cfType) {
  set<CommFlare*> heldCF;

  set<CommFlare*> target; 

  for (auto x: contents) {

    if (cfType == "held") {
        target = (x.second->npu)->GetHeldCFlares(); 
    } else if (cfType == "owned") {
        target = (x.second->npu)->GetOwnedCFlares(); 
    } else if (cfType == "processed") {
        target = (x.second->npu)->GetProcessedCFlares(); 
    } else {
      throw invalid_argument( "received invalid cfType" );
    }

    heldCF.insert(target.begin(), target.end()); 
  }

  return heldCF;
}