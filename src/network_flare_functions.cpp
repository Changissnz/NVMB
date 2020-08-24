#include "network.hpp"
using namespace std;

///////////////// discovery flare functions
/*
for all neighbor identifiers of owned flares, update the nodes that hold them.
*/
///
void Network::AssignHeldDF(NVMBNode* n) {
    NodeProcessingUnit* npu = n->npu;
    set<DiscoveryFlare*> of = npu->GetOwnedFlares();
    set<DiscoveryFlare*>::iterator it = of.begin();
    NVMBNode* n2;
    NodeProcessingUnit* npu2;

    while (it != of.end()) {
        auto it2 = contents.find((*it)->GetLocation()); 
        if (it2 == contents.end()) {
            it++; 
            continue; 
        } 

        n2 = (it2)->second; 
        npu2 = n2->npu;
        npu2->AddHeldDF((*it));
        it++;
  }
}

void Network::AssignHeldDFAll() {
  map<int, NVMBNode*>::iterator it;
  NodeProcessingUnit* npu;

  for ( it = contents.begin(); it != contents.end(); it++ ) {
    AssignHeldDF(it->second);
  };
}

void Network::UpdateNodeDFOwned() {
  // run update for each npu
  map<int, NVMBNode*>::iterator it;
  NVMBNode* n;
  NodeProcessingUnit* npu;

  //#pragma omp parallel num_threads(DEFAULT_NUM_THREADS)
  //{
    //#pragma omp for
    for (it = contents.begin(); it != contents.end(); it++) {
        n = it->second;
        npu = n->npu;
        npu->UpdateOwnedFlares();
        ///cout << "UPDATE NODE DF " << it->first << endl; 
    };
  //} 
}

void Network::UpdateNodeDFNeighbors() {
    NVMBNode* n;
    NodeProcessingUnit* npu;

    //#pragma omp parallel num_threads(DEFAULT_NUM_THREADS)
    //{
    //    #pragma omp for
        for (map<int, NVMBNode*>::iterator it = contents.begin(); it != contents.end(); it++) {
            ///cout << "Update node df " << it->first << endl; 
            n = it->second;
            npu = n->npu;
            npu->UpdateHeldFlares("location");
        };
    //}
}


////////////////// commerce flare functions

/*
// description
increments a CommFlare by one edge
*/
bool Network::IncrementCFlare(CommFlare* cf) {
  if (!cf->IsValid() || cf->finished) {
    return false;
  }

  // try updating the next position only if neighbor
  int currentPos = cf->GetPosition();
  int next = cf->GetNextPosition();

  // target reached
  if (next == -1) {
    cf->UpdatePosition();
    cf->MarkReached();
    return false;
  }

  // check that it is neighbor
  NVMBNode* refNode = contents[currentPos];
  
  // case: reference node no longer exists 
  if (refNode == nullptr) {
      cf->MarkInvalid(); 
      cf->nullEdge = make_pair(cf->GetPreviousPosition(), currentPos); 
      return false; 
  }
    
  if (refNode->NeighborExists(next) && contents[next] != nullptr) {
    cf->UpdatePosition();
  } else {
    cf->MarkInvalid();
    cf->nullEdge = make_pair(currentPos, next); 
    return false;
  }

  IncrementNodeBankByCF(cf);
  return true;
}

// TODO : should return message for Node to update its best paths
/*
// description
moves a commFlare until destination reached or route failure
*/
void Network::MoveCFlare(CommFlare* cf) {
  cf->initialized = true;
  bool moved = true;
  int source = cf->currentDestination;
    ///cout << "source here:\t" << source << endl; 
  if (COMM_FLARE_INCREMENT == -1) {
    while (moved) {
      moved = IncrementCFlare(cf);
    }
  } else {
      /// TODO: below is a bug. DO NOT USE
    for (int x = 0; x < COMM_FLARE_INCREMENT; x++) {
      moved = IncrementCFlare(cf);
      if (!moved) {
        break;
      }
    }
  }

  cf->active = false; 
}

/*
// description
updates CFlare attribute.
If destination has been reached, and
  receiver is destination, then done,
  receiver is not destination, then update,
Else
  done
*/
void Network::UpdateCFlare(CommFlare* cf, int nodeSource, int verbose) {

    if (verbose == 1 || verbose == 2) {
        cout << "DISPLAYING COMM FLARE INFO" << endl; 
        cf->DisplayInfo(); 
    }

    if (cf->finished) {
        cout << "FINISHED" << endl; 
        return; 
    }

    // check success of routing
    bool valid = (cf->vertexLocation == cf->currentDestination) ? true : false;
    cf->valid = valid;

    // case : invalid route
    if (!cf->valid) {
        if (verbose == 1 || verbose == 2) {
            cout << "CF_UPDATE 2: invalid flare at location " << nodeSource << " w/ dest. " << cf->currentDestination << endl;
            cout << "null edge is " << cf->nullEdge.first << " " << cf->nullEdge.second << endl;        
        };

        // log response
        cf->finished = true;

        /// CHECK: transmission failure 
        contents[nodeSource]->IncrementTransmission(false); 
        return;
    }

    /// CHECK: transmission success 
    contents[nodeSource]->IncrementTransmission(true); 

    // check for direct or indirect comm
    auto prop = cf->GetProposition();

    cf->targetReached = (cf->vertexLocation == (prop->referenceNodes).second) ? true : false;
    if (cf->targetReached && cf->valid) {
        cf->finished = true; 
    }

    // update current destination : sender (indirect), receiver (direct)
    TransferCFlare(cf, nodeSource);

    int prevDest = cf->currentDestination;
    if (!cf->finished) {

        cf->currentDestination = (prop->referenceNodes).second;
    // update path to current destination
        /// TODO: re-factor this chunk with below.
        if (contents[prevDest] != nullptr) {
            NodeProcessingUnit* npu = contents[prevDest]->npu;
            auto pathInfo = npu->GetPathForDestination(cf->currentDestination);

            if (pathInfo.second) {     // path does not exist for next destination
                cf->targetPath = pathInfo.first;
                cf->pathIndex = 0;
            } 
        } else { 
            cf->finished = true; 
            cf->valid = false;
        }
    }

    if (verbose == 1 || verbose == 2) {
        cout << "CF_UPDATE" << endl;
        cout << "previous dest: " << prevDest << " current dest: " << cf->currentDestination << endl;
        cout << "finished: " << cf->finished << " valid: " << cf->valid << endl;
        cout << "path to current dest:\t";
        DisplayIterable(cf->targetPath);
        cout << "----------------------" << endl;
    }
}

// TODO : check the CommFlare variables `finished` and `valid`
/*
// description
transfers possessive ownership of CommFlare from source to current destination.
CAUTION : no error-checking enabled.
*/
void Network::TransferCFlare(CommFlare* cf, int nodeSource) {


  // remove the comm flare from source held
    auto n = contents[nodeSource];
    if (n != nullptr) {
        auto npu = n->npu;
        npu->RemoveHeldCF(cf);
    }


    // add comm flare to vertex location held
    auto n2 = contents[cf->currentDestination];
    if (n2 != nullptr) {
        auto npu2 = n2->npu;
        cf->active = false;
        npu2->AddHeldCF(cf);
    }
}

/*
// description
Processes all held commflares belonging to node n. 
Requires two calls per timestamp, for commflares. 
*/
void Network::ProcessNodeCF(NVMBNode* n, int verbose) {

    if (n == nullptr) {
        return; 
    }

    NodeProcessingUnit* npu = n->npu;
    set<CommFlare*> heldCFlares = npu->GetHeldCFlares();
  
    int c = 0; 
    for (set<CommFlare*>::iterator it = heldCFlares.begin(); it != heldCFlares.end(); it++) {
    // if comm flare has not yet been initialized, set its costs 
    auto xx = (*it)->proposition; 
    if ((*it)->proposition == nullptr) {
        (*it)->MarkInvalid();
        continue; 
    }

    /// TODO:!!
    if (!(*it)->valid) {
        cout << "NOT VALID" << endl; 
        continue; 
    }

    if ((*it)->initialized == false) {
      SetBondModCosts(*it);
    }

    /// TODO: 
    /*
    deleted pointer check here. 

    test that held, owned, and processed cflares get cleared after 
    every round. 
    */

    // flares are in cached status
    if (!(*it)->active) {
      continue;
    }

    ///+
    if (((*it)->finished) || !(*it)->valid) {
      continue; 
    }

    c++; 

    // move CommFlare along path and make necessary changes to flare variable
    MoveCFlare(*it);
    UpdateCFlare(*it, n->GetId(), verbose);
  }
}

/*
// description
sets the held CF of all nodes to active status
*/
void Network::SetHeldCFToActive() {

  set<CommFlare*> heldCF; //? 
  NVMBNode* n;
  for (auto x: contents) {
    n = x.second;

    if (n == nullptr) { continue;}

    NodeProcessingUnit* npu = n->npu;
    set<CommFlare*> heldCFlares = npu->GetHeldCFlares();
    for (set<CommFlare*>::iterator it = heldCFlares.begin(); it != heldCFlares.end(); it++) {
      (*it)->active = true;
    }
  }
}

/*
// description
increments the value of a node's bank given the cf cost of the proposer
*/
float Network::IncrementNodeBankByCF(CommFlare *cf) {
    int pos = cf->GetPosition();

    // node no longer exists
    if (contents.find(pos) == contents.end() || contents[pos] == nullptr) {
        return 0.0;
    }

    // destination nodes do not tax
    if (cf->GetPosition() == cf->currentDestination) {
        return 0.0;
    }

    // calculate the tax of transportation and apply it to relevant nodes
    float taxCost = cf->costs[cf->GetProposer()] * TRANSPORT_TAX;
    contents[pos]->UpdateValue(-1 * taxCost); // node gets gain
    cf->costs[cf->GetProposer()] += taxCost; // proposer incurs loss

    // record transaction into bank data
    NVMBNode* n = contents[cf->GetPosition()];
    map<string,string> od = map<string,string>(); 
    od["nodes"] = to_string(cf->GetPreviousPosition()); 
    TimestampUnit* tsu = new TimestampUnit("cf tax", n->GetTimestamp(), -1 * taxCost, od);

    BankUnit* bu = n->GetBank();
    bu->LogTimestampUnit(tsu); 
    bu->UpdateCurrency(tsu->impact); 

    return taxCost;
}


void Network::UpdateCompleteCFlares() {
    NodeProcessingUnit* npu; 
    for (auto x: contents) {        
        if (x.second != nullptr) { (x.second)->UpdateFinishedCF();}         
    }
}

/*
performs necessary updates for node in the event of a misroute. 
*/ 
void Network::MisrouteUpdate(CommFlare* cf, int nodeSource) {

    // check that comm flare's null edge exists 
    if (cf->nullEdge == make_pair(-1,-1)) {
        cout << "X misroute update cannot be performed, no edge specified" << endl; 
        return ; 
    }

    // check that null edge is valid 
    if (EdgeExists(cf->nullEdge.first, cf->nullEdge.second, false)) {
        cout << "X misroute update: edge " << NodePairToString(cf->nullEdge.first, cf->nullEdge.second) << " still exists" << endl; 
        return; 
    }

    // update node info
    auto npu = contents[nodeSource]->npu; 
    npu->RemoveRecordsOfEdge(NodePairToString(cf->nullEdge.first, cf->nullEdge.second)); 
    return ;
}

/*
// description

*/
/// TODO: ! logResp needs to be set to True
void Network::ProcessNodeCFAll_(bool logResp, int verbose) {
  // get the max number of neighbors
  /// TODO: check for actual contents size (after node deletion)
    SetMaxConnectivity();
    NVMBNode* n;

    // process held cf
    //#pragma omp parallel num_threads(DEFAULT_NUM_THREADS)
    //{
    //    #pragma omp for 
        for (map<int,NVMBNode*>::iterator it = contents.begin(); it != contents.end(); it++) {
            n = it->second;
            ProcessNodeCF(n, verbose);
        };
    //}


    SetHeldCFToActive();

    // process each node's response
    //#pragma omp parallel num_threads(DEFAULT_NUM_THREADS)
    //{
    //    #pragma omp for 

        for (map<int,NVMBNode*>::iterator it = contents.begin(); it != contents.end(); it++) {
            n = it->second;
            if (n != nullptr) {
                n->RespondToCommFlares(logResp, verbose);

            }
        };
    //}

    // execute responses
    //#pragma omp parallel num_threads(DEFAULT_NUM_THREADS)
    //{
    //    #pragma omp for 

        for (map<int,NVMBNode*>::iterator it = contents.begin(); it != contents.end(); it++) {
            n = it->second;

            if (n != nullptr) {
                ProcessNodeCFResponses(n->npu);
            }
        };
    //}   

}

/*
// description
processes all CommFlares, direct and phantom
*/
void Network::ProcessNodeCFAll(int verbose) {

    ProcessNodeCFAll_(true, verbose);

    UpdateCompleteCFlares();
    ProcessNodeCFAll_(true, verbose);
    UpdateCompleteCFlares(); 

    // iterate through plans of interest 
    ProcessNodeCFAll_(true, verbose);
    UpdateCompleteCFlares(); 
}

/*
// description
marks all executed plans that are finished as finished.

iterates through held comm flares of each
*/
void Network::MarkExecutedPlansAsFinished() {

  // iterate through processed CF and move their plans from executed to finished status
    set<CommFlare*> scf;

    for (map<int, NVMBNode*>::iterator it = contents.begin(); it != contents.end(); it++) {

        auto npu = (it->second)->npu;
        scf = npu->GetProcessedCFlares();
        for (auto cf: scf) {
        // flare is finished
        if (cf->finished) {
            // check if node owner exists
            if (!NodeExists(cf->nodeOwnerIdentifier)) {
            continue;
            }

            // find executed plan (node owner) corresponding to cflare and mark it as finished
            auto strat = (contents[cf->nodeOwnerIdentifier])->GetStrategy();
            strat->FinishExecutedPlan("plan", cf->planIdentifier);
        }
        }
    }
}

/*
// description
processes all responses for each node's received CommFlares, including 
modifying by bonds/contracts and marking finished. 
*/
void Network::ProcessNodeCFResponses(NodeProcessingUnit* npu) {
  set<CommFlare*> pcf = npu->GetProcessedCFlares();

  for (auto x: pcf) {
    // response is true and flare is finished
    if (x->responses[npu->GetNOI()] && x->finished && x->valid) {
        ModifyNetworkByCFBond(x);
        ModifyNetworkByCFContract(x); 
    }

    // response is reject, finished comm flare
    if (!x->responses[npu->GetNOI()]) {
      x->finished = true;
    }
  }
}

/*
// description
given a CommFlare that is accepted by both sender and receiver,
adds an undirected edge between them.
*/
void Network::ModifyNetworkByCFBond(CommFlare* cf) {
  if (cf->GetEventType() == "make_bond") {
    AddEdge(cf->GetSender(), cf->GetReceiver(), false);
  } else if (cf->GetEventType() == "break_bond") {
    BreakEdge(cf->GetSender(), cf->GetReceiver(), false);
  }
}

/// TODO : untested
/*
// description
*/
void Network::ModifyNetworkByCFContract(CommFlare* cf) {
  NodeProcessingUnit* npu;
  if (cf->GetEventType() == "make_contract") {
    auto npu1 = contents[cf->GetSender()]->npu;
    npu1->InitiateContract(cf);
    auto npu2 = contents[cf->GetReceiver()]->npu;
    npu2->InitiateContract(cf);
  } else if (cf->GetEventType() == "break_contract") {
    auto npu1 = contents[cf->GetSender()]->npu;
    npu1->BreakContract(cf);
    auto npu2 = contents[cf->GetReceiver()]->npu;
    npu2->BreakContract(cf);
  }
}

/*
// description
calculates the max connection a Node can have
*/
void Network::SetMaxConnectivity() {
   maximallyConnected = 0; 
   for (auto c: contents) {
       if (c.second == nullptr) {continue;} 
       int z = c.second->GetNumberOfNeighbors(); 
       if (z > maximallyConnected) {maximallyConnected = z;} 
   }

}

/*
// description
calculates the bond modification costs for a CommFlare that has just been initialized.

// output
true if sets bond modification costs successfully, else false
*/
bool Network::SetBondModCosts(CommFlare *cf) {
    if (!cf->valid) {
        return false; 
    }

  // set up costs map
  Proposition* prop = cf->proposition;
    if (prop == nullptr) { 
        cout << "warning: could not activate null CF" << endl; 
        return false; 
    }

  cf->costs[prop->proposerNode] = 0;
  cf->costs[prop->referenceNodes.first] = 0;
  cf->costs[prop->referenceNodes.second] = 0;

    bool stat = true; 
    if (contents[prop->proposerNode] != nullptr && contents.find(prop->proposerNode) != contents.end()) {
        int conn1 = contents[prop->proposerNode]->GetNumberOfNeighbors();
        float curr1 = contents[prop->proposerNode]->GetWorth();
        float bondModCost1 = BondModificationCost1(contents[prop->proposerNode]->GetNodeRank(), conn1, maximallyConnected, curr1);
        float bondModCost11 = BondModificationCost2(contents[prop->proposerNode]->GetNodeRank(), contents[prop->proposerNode]->GetNumberOfNeighbors(), contents.size(), curr1);
        cf->costs[prop->proposerNode] += bondModCost1 + bondModCost11;
    } else {
        stat = false; 
    }
    
  // calculate cost 2:
  if (contents[prop->referenceNodes.first] != nullptr && contents.find(prop->referenceNodes.first) != contents.end()) {
    float nr2 = contents[prop->referenceNodes.first]->GetNodeRank();
    float curr2 = (contents[prop->referenceNodes.first])->GetWorth();
    float bondModCost2 = BondModificationCost2(nr2, contents[prop->referenceNodes.first]->GetNumberOfNeighbors(), contents.size(), curr2);
    cf->costs[prop->referenceNodes.first] += bondModCost2;
  } else {
        stat = false; 
    }

  // calculate cost 3
  if (contents[prop->referenceNodes.second] != nullptr && contents.find(prop->referenceNodes.second) != contents.end()) {
    float nr3 = contents[prop->referenceNodes.second]->GetNodeRank();
    float curr3 = (contents[prop->referenceNodes.second])->GetWorth();
    float bondModCost3 = BondModificationCost2(nr3, contents[prop->referenceNodes.second]->GetNumberOfNeighbors(), contents.size(), curr3);
    cf->costs[prop->referenceNodes.second] += bondModCost3;
  } else {
        stat = false; 
    }


  return stat;
}

void Network::DisplayExistingContracts() {
  for (auto x: contents) {
    cout << "** contracts for node " << x.first << endl;
    (x.second)->DisplayActiveContracts();
    cout << "-------------------------" << endl;
  }
}

/*
// description 
Performs activities for all active contracts. 
*/ 
void Network::ActiveContractExecution() {
    //
    for (auto x: contents) {
        auto strategy = x.second->GetStrategy(); 
        auto contracts = strategy->activeContracts; 
        for (auto y: strategy->activeContracts) {
            ConductContract(y, x.first); 
        }
    } 

}

/*
// description 
Network will update bank of nodeSource according to successful transmission of 
contract. 

Any missing edges experienced during routing will be updated to `nodeSource`. 

*/ 
bool Network::ConductContract(Contract* C, int nodeSource) { 

    if (contents.find(nodeSource) == contents.end()) {
        cout << "invalid node source " << nodeSource << " . Contract could not be updated" << endl; 
        return false; 
    }

    // find the best paths 
    auto npuX = contents[nodeSource]->npu;
    auto pathInfo = npuX->GetBestPaths()[C->receiver];

    vector<pair<int,int>> nullEdges = vector<pair<int,int>>(); 

    bool success = false; 
    for (auto x: pathInfo) {
        auto nullEdge = RouteContract(C, x.first, nodeSource); 

        // success 
        if (nullEdge.first.first == -1 && nullEdge.second) {
            success = true; 
            break; 
        } else if (nullEdge.first.first == -1 && !nullEdge.second) {
            continue; 
        } else {/// TODO: fail, have to update Node's missing edge  
            nullEdges.push_back(nullEdge.first); 
            npuX->RemoveRecordsOfEdge(NodePairToString(nullEdge.first.first, nullEdge.first.second));  
        } 
    }

    // update source node's edge records 
    for (auto x: nullEdges) {
        npuX->RemoveRecordsOfEdge(NodePairToString(x.first, x.second));  
    }

    if (!success) {
        map<string,string> m = {{"success", to_string(0)}}; 
        auto tsu = new TimestampUnit("contract activity", contents[nodeSource]->GetTimestamp(), 0, m); 
        auto bu = contents[nodeSource]->GetBank(); 
        bu->LogTimestampUnit(tsu); 
        bu->UpdateCurrency(tsu->impact); 
    }


    /// CHECK: transmission
    contents[nodeSource]->IncrementTransmission(success); 
    return success; 
}

/*
// description 
Routes a contract along its `path` from nodeSource. 
Logs TimestampUnits according to `contract tax` or `contract activity` 

// return
[0] null edge 
[1] success of routing contract 
*/ 
pair<pair<int,int>, bool> Network::RouteContract(Contract* C, vector<int> path, int nodeSource) {

    int location = nodeSource;

    // case: destination does not exist 
    if (contents.find(C->receiver) == contents.end()) {
        return make_pair(make_pair(-1,-1), false); 
    }

    /// TODO: below is design-choice. Worth dependent on nodeSource and not receiver 
    float worth = contents[nodeSource]->GetWorth(); 
    float gains = (worth * MAX_GROWTH_RATE) - worth;

    // route along path 
    pair<int,int> nullEdge = make_pair(-1,-1);  
    float contractTax = 0.0;
    int i = 0;

    for (i = 0; i < path.size() - 2; i++) {
        // check for edge 
        if (!EdgeExists(path[i], path[i + 1], false)) {
            nullEdge = make_pair(path[i], path[i + 1]); 
            break; 
        }

        contractTax = gains * TRANSPORT_TAX; 
        gains -= contractTax;    
        map<string,string> m = {{"success", to_string(1)}}; 
        auto tsu = new TimestampUnit("contract tax", contents[path[i + 1]]->GetTimestamp(), -contractTax, m); 
        BankUnit* bu = contents[path[i + 1]]->GetBank(); 
        bu->LogTimestampUnit(tsu);
    }

    // path route was invalid 
    if (nullEdge.first != -1) {
        return make_pair(nullEdge, false);  
    } 

    // route along the last edge
    bool success = true; /// TODO: check this 

    TimestampUnit* tsu; 

    if (!EdgeExists(path[i], path[i + 1], false)) {
        nullEdge = make_pair(path[i], path[i + 1]); 
        success = false;

        auto loss = ContractRouteFailureValue(contents[nodeSource]->GetWorth()); 
        tsu = new TimestampUnit("contract activity failure", contents[nodeSource]->GetTimestamp(), loss, map<string,string>());  

    } else { // reached destination 
        map<string,string> m = {{"success", to_string(1)}}; 
        tsu = new TimestampUnit("contract activity", contents[nodeSource]->GetTimestamp(), gains, m); 
    }
    
    // log timestamp and its impact 
    BankUnit* bu = contents[nodeSource]->GetBank();
    bu->LogTimestampUnit(tsu);
    bu->UpdateCurrency(tsu->impact); 
    return make_pair(nullEdge, success);   
}