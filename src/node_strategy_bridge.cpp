/*
this file contains methods that allow Node to make and handle the plans proposed by
NodeStrategos
*/
#include <stdlib.h>
#include "nvmb_node.hpp"

/*
// description
executes plan and runs other processes
*/
bool NVMBNode::ProcessOne(int verbose) {

  // execute current plan
  if (mode == "active") {
    ExecutePlan(verbose);
  }

  ContinueDiscovery();
  return true;
}

/*
// description
call this method every round.
*/
void NVMBNode::ContinueDiscovery() {
    if (npu->DiscoveryFinished() && npu->discoveryOn) {
        FinalizeDiscovery(); // TODO requires that Plan keeps record of performance
        MarkDiscoveryFinished();
        npu->ClearDiscoveryData(); 
    } else {
        npu->ProcessDF();
    } 
}

void NVMBNode::UpdateNodeRank(float r) {
  NodeStrategos *s = npu->GetStrategy();
  s->SetRank(r);
}

/// TODO: not used? 
/*
// description
chooses a random neighbor.
*/
int NVMBNode::ChooseRandomNeighbor() {
    int sz = neighborIdentifiers.size();

    if (sz == 0) {
        return -1;
    }

    // choose random index
    set<int>::iterator it = neighborIdentifiers.begin();
    int r = rand() % neighborIdentifiers.size();
    advance(it, r);
    return *it;
}

/*
scrambles node nature variables 
*/ 
void NVMBNode::ScrambleNature() {
    auto ns = GetStrategy(); 

    ns->competition = RandomFloat(); 
    ns->greed = RandomFloat(); 
    ns->negotiation = RandomFloat(); 
    ns->growth = RandomFloat(); 

    ns->CalculateNormedNatureVar();
}

/*
// description
Plan is to execute one of the following:
- CommFlare:
  - send proposition,
      if destination reaches and accepts, then contract sent back to sender
  - send calculator,
      to get "best" path
- DiscoveryFlare:
  -
*/
void NVMBNode::ExecuteGivenPlan(Plan* p, int verbose) {

    if (verbose == 1 || verbose == 2) {
        cout << "* node " << identifier << " executes plan at timestamp " << timestamp << endl;
        p->DisplayInfo();
    }

    NodeStrategos* strategy = GetStrategy(); 
    if (p->planType == "discovery") {
        // TODO : need to erase previous discovery data
        if (!(strategy->discoveryOn)) {
        InitializeDiscovery();
        } else {
        p->execSuccess = false;
        return;
        }
    } else if (p->planType == "invalid") {
        return;
    } else {
        CommFlare* cf = npu->ProduceCommFlare(p);
        npu->AddOwnedCF(cf);
        npu->AddHeldCF(cf);
    }

    // place Plan in Strategos
    strategy->AddExecutedPlan(p);
}

void NVMBNode::MarkDiscoveryFinished() {
  // iterate through executedPlans and mark
  NodeStrategos* strategy = GetStrategy();
  vector<Plan*> plans = strategy->GetPlansInExecution();
  if (plans.size() == 0) {
    return;
  }

  Plan* p;
  int ind = -1;
  for (int i = 0; i < plans.size(); i++) {
    // mark as executed;
    if (plans[i]->planType == "discovery") {
      ind = i;
      break;
    };
  };

  // mark plan as finished
  strategy->FinishExecutedPlan("direct", ind);
  npu->discoveryOn = false;
}

/*
// description
Executes the next plan designated by NodeStrategos.
If there is None, executes the default tactic.
*/
void NVMBNode::ExecutePlan(int verbose) {

    ////// CURRENT
    CalculatePlans(verbose);

    NodeStrategos* ns = GetStrategy();
    ns->AddPlans();

    ////// TODO
    Plan* p = ns->ChoosePlan();
    if (verbose == 1 || verbose == 2) {
        cout << "node " << GetId() << " chooses below plan" << endl; 
    }
    ExecuteGivenPlan(p, verbose);
}

/*
// description 
This method is used for testing; `testModeSpec` specifies the type of plan 
node can make. 
*/ 
void NVMBNode::CalculatePlansTestMode(int verbose) {
  NodeStrategos* ns = npu->GetStrategy();
  if (testModeSpec == "discovery") {
    DiscoveryReview(verbose);
  } else if (testModeSpec == "make_bond") {
    NewBondReview(ns->competitors, verbose);
  } else if (testModeSpec == "break_bond") {
    OldBondReview(ns->competitors, verbose);
  } else if (testModeSpec == "make_contract") {
    NewContractReview(verbose);
  } else if (testModeSpec == "break_contract") {
      ExistingContractReview(verbose);
  } 
}

/*
// description
instance calculates x plans
*/
void NVMBNode::CalculatePlans(int verbose) {
  // preprocessing methods:
    // clear plans cache

  NodeStrategos* ns = npu->GetStrategy();
  ns->ClearFuturePlans();
    // calculate discovery

  if (testModeSpec != "") {
    CalculatePlansTestMode(verbose); 
    return; 
  }

  if (manInput == true) {
    return; 
  }

    if (verbose == 2) { 
        cout << "*********************************************" << endl; 
        cout << "NODE " << GetId() << "\n" << "DISCOVERY REVIEW" << endl;
    };

  if (verbose == 2) {
    cout << "PRESENT COMPETITORS" << endl; 
    DisplayIterable(GetStrategy()->competitors); 
  }
  
  // do review
  DiscoveryReview(verbose);

  if (verbose == 2) { cout << "NEW BOND REVIEW" << endl;};
  NewBondReview(ns->competitors, verbose);

  if (verbose == 2) { cout << "OLD BOND REVIEW" << endl;};
  OldBondReview(ns->competitors, verbose);

  if (verbose == 2) { cout << "EXISTING CONTRACT REVIEW" << endl;};
  ExistingContractReview(verbose);

  if (verbose == 2) { cout << "NEW CONTRACT REVIEW" << endl;};
  NewContractReview(verbose);

  if (verbose == 2) {
    cout << "=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=." << endl;
    cout << "=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=.=." << endl;
  };

}

/*
// description
*/
set<int> NVMBNode::GetKnownNodes() {
    auto bp = npu->GetBestPaths();
    set<int> kn;

    // insert nodes from best paths
    for (auto it = bp.begin(); it != bp.end(); it++) {
        kn.insert(it->first);
    }

    // insert nodes from neighborIdentifiers
    for (auto it2 = neighborIdentifiers.begin(); it2 != neighborIdentifiers.end(); it2++) {
        kn.insert(*it2);
    }

    knownNodes.clear(); 
    knownNodes = kn;
    knownNodes.erase(identifier); 
    return knownNodes;
}

/// TODO: keep a running record of Node's competitors
/*
// description
*/
set<int> NVMBNode::FindCompetitors() {

  GetKnownNodes();
  set<int> competitors = set<int>();
  for (auto it = knownNodes.begin(); it != knownNodes.end(); it++) {
    if (ConsiderCompetitor(*it)) {
      competitors.insert(*it);
    }
  }
  return competitors;
}

///////////////////// methods are used to review ///////////////////////////////

/*
// description
iterates through the set of possible new bonds and outputs the bond with
the lowest scoring bond-add advantage measure.

// output
pair, with [0] a bond
            [1] its score
*/
void NVMBNode::NewBondReview(set<int> recognizedCompetitors, int verbose) {

  // collect bonds to review
  set<string> possibleNewBonds = SelectPossibleBondsToAdd();
  NodeStrategos* strategy = npu->GetStrategy();
  Plan* p;

  if (verbose == 2) { cout << "- New bond review" << endl;};

  for (set<string>::iterator it = possibleNewBonds.begin(); it != possibleNewBonds.end(); it++) {
    auto b = StringToNodePair(*it);
    auto scoreClass = (identifier == b.first || identifier == b.second) ? "make_bond/direct": "make_bond/phantom";
    auto x = MeasureBondCreation(b, recognizedCompetitors);

    /// TODO: error-check below 
    if (scoreClass == "make_bond/direct") {
        if (identifier == b.first) {
            p = strategy->MakeBondPlan("make_bond", identifier, b.first, b.second, x);  
        } else {
            p = strategy->MakeBondPlan("make_bond", identifier, b.second, b.first, x);
        }
    } else {
        p = strategy->MakeBondPlan("make_bond", identifier, b.first, b.second, x);
    }

    if (verbose == 2) {
      p->DisplayInfo();
    }

    strategy->InsertFuturePlan(p);
  }

  if (verbose == 2) {
    cout << "======================================" << endl;
  }
}

/*
// description
calculates the best
*/
void NVMBNode::OldBondReview(set<int> recognizedCompetitors, int verbose) {
  // get options
  set<string> oldBonds = SelectPossibleBondsToDelete();
  NodeStrategos* strategy = npu->GetStrategy();
  Plan* p;

  if (verbose == 2) { cout << "- Old bond review" << endl;};

  for (auto it = oldBonds.begin(); it != oldBonds.end(); it++) {
    auto b = StringToNodePair(*it);
    auto x = MeasureBondDeletion(b, recognizedCompetitors);
    auto scoreClass = (identifier == b.first || identifier == b.second) ? "break_bond/direct": "break_bond/phantom";
    p = strategy->MakeBondPlan("break_bond", identifier, b.first, b.second, x);

    if (verbose == 2) {
      p->DisplayInfo();
    }

    strategy->InsertFuturePlan(p);
  }

  if (verbose == 2) {
    cout << "======================================" << endl;
  }
}

////////////////////////////////////////////////////////////// contract review

void NVMBNode::ExistingContractReview(int verbose) {
  NodeStrategos* ns = npu->GetStrategy();
  Plan* p;

  /*
  cout << "contract info for " << identifier << endl; 
  for (auto c: ns->activeContracts) {
    c->DisplayInfo(); 
  }
  cout << "AFTERR" << endl; 
  */ 

  for (int i = 0; i < ns->activeContracts.size(); i++) {
    auto x = ConsiderExistingContract(ns->activeContracts[i]);
    p = ns->BreakContractPlan(GetId(), GetId(), (ns->activeContracts[i])->receiver, x, (ns->activeContracts[i])->idn);
    /*
    cout << "plan for contract" << endl; 
    (ns->activeContracts[i])->DisplayInfo(); 
    cout << "^^^" << endl; 
    p->DisplayInfo(); 
    */
    ns->InsertFuturePlan(p);
  };
  if (verbose == 2) {cout << "================================" << endl;};
}

/*
// description
*/
set<int> NVMBNode::GetNodesWithoutContract() {
  auto kn = GetKnownNodes();

  // collect all node identifiers with contracts
  NodeStrategos* ns = npu->GetStrategy();

  for (auto it = ns->activeContracts.begin(); it != ns->activeContracts.end(); it++) {
    kn.erase((*it)->receiver);
  };

  kn.erase(GetId()); 

  return kn;
}

/*
// description
*/
void NVMBNode::NewContractReview(int verbose) {
  set<int> nodeIdns = GetNodesWithoutContract();
  Plan* p;

  if (verbose == 2) {
    cout << "** nodes without contract to consider" << endl; 
    DisplayIterable(nodeIdns); 
  }

  NodeStrategos* strategy = npu->GetStrategy();

  for (auto it = nodeIdns.begin(); it != nodeIdns.end(); it++) {
    auto x = ConsiderNewContract(*it);
    p = strategy->MakeContractPlan(identifier, identifier, *it, x); 
    strategy->InsertFuturePlan(p);
  }
  if (verbose == 2) {cout << "================================" << endl;};
}

/*
// description
*/
bool NVMBNode::IsContractedWith(int nodeIdentifier) {
  vector<Contract*> activeContracts = GetStrategy()->activeContracts;
  for (auto x: activeContracts) {
    if (nodeIdentifier == x->receiver) {
      return true;
    }
  }
  return false;
}

/*
// description
updates contract 
*/
void NVMBNode::UpdateContract(int nodeIdentifier, bool transmissionSuccess, float delta, int pathDistance) {

  // iterate through and find contracts
  auto ac = GetStrategy()->activeContracts;

  for (auto x: ac) {
    if (x->receiver == nodeIdentifier) {
      x->UpdateContractData(transmissionSuccess, delta, pathDistance);
      return;
    }
  }
}

void NVMBNode::UpdateContract(int nodeIdentifier, bool alve) {

  // iterate through and find contracts
  auto ac = GetStrategy()->activeContracts;

  for (auto x: ac) {
    if (x->receiver == nodeIdentifier) {
      x->UpdateContractData(alve);
      return;
    }
  }
}

/*
// description 
outputs a vector of size 4: future,cached,execution,finished
*/
vector<int> NVMBNode::GetPlanCounts() {
  auto ns = GetStrategy();
  int f1 = ns->GetFuturePlans().size(); 
  int f2 = ns->GetPlansInCache().size();
  int f3 = ns->GetPlansInExecution().size(); 
  int f4 = ns->GetFinishedPlans().size();   
  
  return {f1,f2,f3,f4}; 
}

/*
// description 
outputs a vector of size 3: owned,cached,held
*/
vector<int> NVMBNode::GetFlareCounts() {
    int of = npu->GetOwnedFlares().size(); 
    int ofc = npu->GetOwnedFlaresCache().size(); 
    int hf = npu->GetHeldFlares().size();
    cout << "GOT"<< endl; 
    return {of,ofc,hf}; 
}


/*
outputs a vector of size 3: held,owned,processed
*/ 
vector<int> NVMBNode::GetCommFlareCounts() {

    int held = (npu->GetHeldCFlares()).size(); 
    int owned = (npu->GetOwnedCFlares()).size(); 
    int processed = (npu->GetProcessedCFlares()).size(); 

    return {held, owned, processed}; 
}

/*
// description
iterates through held flares and either accepts or rejects them.
*/
void NVMBNode::RespondToCommFlares(bool toLog, bool verbose) {

  set<CommFlare*> cfs = npu->GetHeldCFlares();
  bool responseInfo;

  for (auto c: cfs) {
    if (c->valid && c->active) {
        /// cout << "NODE " << GetId() << " LOGOS " << c->planIdentifier << endl;  
      responseInfo = RespondToCommFlare(c, toLog, verbose);
      if (responseInfo) {
        npu->RemoveHeldCF(c);
      }

      npu->AddProcessedCF(c);

      IncrementResponse(c->responses[GetId()]); 
    }
  }
}

/*
// description
responds to contents of CommFlare

Responder node could be one of receiver or sender

// output
[1] finished with CommFlare, remove from held
*/
bool NVMBNode::RespondToCommFlare(CommFlare *cf, bool writeResp, bool verbose) {

  // determine if CommFlare done
  string role = cf->GetRoleOfNode(identifier);
  assert (role != "proposer");
  bool remove = (role == "sender") ? false: true;

  if (cf->finished) {
    cf->active = false;
  }

  // case: return default
  if (defaultResponseCF.first) {

    cf->responses[identifier] = defaultResponseCF.second;
    if (!defaultResponseCF.second) {
        cf->finished = true; 
    }

    if (verbose) {
      cout << "CF_RESPONSE: node " << identifier << " responds to comm flare with " << defaultResponseCF.second << endl;
      cf->DisplayInfo();
    }

    return remove;
  }

  // determine what CommFlare holds
  Proposition* prop = cf->proposition;
  NodeStrategos* strategy = GetStrategy();
  float output;
  float threshold;
  if (prop->category == "make_bond") {
    output = MeasureBondCreation((cf->proposition)->referenceNodes, strategy->competitors);
    if (identifier == cf->GetSender()) { // direct
      threshold = (strategy->competition + strategy->greed) / 2.0;
    } else { // phantom
      threshold = (strategy->negotiation + strategy->greed) / 2.0;
    }
  } else if (prop->category == "break_bond") {
    output = MeasureBondDeletion((cf->proposition)->referenceNodes, strategy->competitors);
    if (identifier == cf->GetSender()) { // direct
      threshold = (strategy->competition + strategy->growth) / 2.0;
    } else { // phantom
      threshold = (strategy->negotiation + strategy->growth) / 2.0;
    }
  } else if (prop->category == "make_contract") {
    assert (identifier == cf->GetSender() || identifier == cf->GetReceiver());
    int targetNode = (identifier == cf->GetSender()) ? cf->GetReceiver(): cf->GetSender();
    output = ConsiderNewContract(targetNode);
    threshold = strategy->growth;
  } else if (prop->category == "break_contract") {
    assert (identifier == cf->GetSender() || identifier == cf->GetReceiver());
    int targetNode = (identifier == cf->GetSender()) ? cf->GetReceiver(): cf->GetSender();
    Contract* c = npu->GetExistingContract(targetNode);
    output = ConsiderExistingContract(c);
    threshold = strategy->greed;
  } else {
    throw invalid_argument("received invalid commflare category");
  }

  // log response to timestamp
    cf->responses[identifier] = (output <= threshold);
    if (!cf->responses[identifier]) {
        cf->finished = true; 
    }

    if (verbose) {
        cout << "CF_RESPONSE: node " << identifier << " responds to comm flare with " << cf->responses[identifier] << endl;
        cf->DisplayInfo();
    };

    return remove;
}

// TODO : untested
/*
// description
calculates response for
*/
void NVMBNode::LogResponse(CommFlare* cf) {
    string eventType = cf->GetProperEventType();
    
    /// TODO: Inspect this. May be burying a bug. 
    if (eventType == "invalid") {
        cout << "invalid event. no response" << endl; 

        return; 
    }

    float impact = cf->costs[identifier];
    string role = cf->GetRoleOfNode(identifier);

    TimestampUnit* tsu = new TimestampUnit(timestamp, role, eventType, cf->responses,
                          cf->costs[identifier], cf->GetProposer(), cf->GetSender(), cf->GetReceiver());

    bank->LogTimestampUnit(tsu);


    // update currency according to response and role. 
    string et = cf->GetEventType(); 
    if (et == "make_contract" || et == "break_contract" || et == "make_bond" || et == "break_bond") {
        // check responses 

            // case: if no response, incur cost 
        if (cf->responses.find(GetId()) == cf->responses.end()) {
            bank->UpdateCurrency(tsu->impact); 
        } else {
            // case: if response, true => incur cost
            if (cf->responses[GetId()] == true) {
                bank->UpdateCurrency(tsu->impact); 
            } 
        }   
    }
}

void NVMBNode::UpdateValue(float c) {
  bank->UpdateCurrency(c);
}

/*
// description 
displays vector of nodes this Node instance is contracted with
*/ 
void NVMBNode::DisplayActiveContracts() {
  vector<int> y;
  vector<Contract*> activeContracts = GetStrategy()->activeContracts;
  for (auto x: activeContracts) {
    y.push_back(x->receiver);
  }
  DisplayIterable(y);
}

void NVMBNode::UpdateFinishedCF() {
        pair<CommFlare*, bool> cfOut; 

        set<CommFlare*> s; 

        while (true) {
            cfOut = npu->UpdateOneCompleteFrom("held"); 
            if (!cfOut.second) {
                break; 
            }
            
            LogResponse(cfOut.first); 
            s.insert(cfOut.first); 
        }

        while (true) {
            cfOut = npu->UpdateOneCompleteFrom("owned"); 
            if (!cfOut.second) {
                break; 
            }

            /// TODO: ?!? 
            if (s.find(cfOut.first) == s.end()) {
                LogResponse(cfOut.first); 
                s.insert(cfOut.first); 
            }
        }
         
        while (true) {
            cfOut = npu->UpdateOneCompleteFrom("processed"); 
            if (!cfOut.second) {
                break; 
            }

            if (s.find(cfOut.first) == s.end()) {
                LogResponse(cfOut.first); 
                s.insert(cfOut.first); 
            }

        }
}

/*
// description 
writes timestamp units out to file if timestamp size threshold reached. 
*/ 
void NVMBNode::WriteOutToFile(bool threshold) {
    // case: number of timestamps has not yet reached write-threshold
    if (threshold) {
        if (bank->GetNumberOfTimestamps() < TIMESTAMP_SIZE_THRESHOLD) {
            return; 
        }
    }

    cout << "SIZE THRESHOLD " << TIMESTAMP_SIZE_THRESHOLD << endl; 
    cout << GetId() << " writing out to file" << endl; 

    // iterate through timestamps and write them out 
    map<int,vector<TimestampUnit*>> tsuh = bank->GetTimestampUnitHistory();
    cout << "size of map: " << tsuh.size() << endl; 
    string s; 
    for (auto c: tsuh) {
        auto vec = c.second; 
        for (auto c_: vec) {
            s = c_->ToDataString();
            reader->WriteRowData(s); 
        }
    }

    // clear bank TimestampUnit history afterwards 
    bank->ClearTimestampUnitHistory(); 
}

/*

*/ 
void NVMBNode::SummarizeLastNTimestamps(int n) {
  vector<TimestampUnit*> tsuh = bank->GetTimestampUnitHistoryAtRange(timestamp - n, timestamp); 
  vector<string> dataTemplate = {"role", "event type", "responses"};//, "summary"}; 
  map<string,map<string,int>> x = npu->ProcessTimestampSequenceForFrequency(tsuh, dataTemplate);

  cout << "TIMESTAMP SUMMARY FOR NODE " << GetId() << endl; 
  for (auto it = x.begin(); it != x.end(); it++) {
    cout << "* " << (*it).first << endl; 
    for (auto z: (*it).second) {
      cout << "\t* " << z.first << " " << z.second << endl;  
    }
    cout << endl; 
  }
}