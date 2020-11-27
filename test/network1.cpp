#include "network1.h"

//// TODO : make option to add node strategos.
void PrintSpacer() {
  cout << "\n\n\n" << endl; 
  cout << "-----------------------------------------------------" << endl; 
  cout << "-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+" << endl; 
  cout << "-----------------------------------------------------" << endl; 
  cout << "-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+" << endl; 
  cout << "-----------------------------------------------------" << endl; 
  cout << "\n\n\n" << endl; 
}

Network Network1() {

  Network N(1, false, false);

  // declare the nodes
  NVMBNode* n0 = new NVMBNode(0, "node_0");
  n0->AddNeighbor(1);
  n0->AddNeighbor(5);
  N.AddNode(n0);

  NVMBNode* n1 = new NVMBNode(1, "node_1");
  n1->AddNeighbor(0);
  n1->AddNeighbor(2);
  n1->AddNeighbor(7);
  N.AddNode(n1);


  NVMBNode* n2 = new NVMBNode(2, "node_2");
  n2->AddNeighbor(1);
  n2->AddNeighbor(3);
  n2->AddNeighbor(7);
  N.AddNode(n2);

  NVMBNode* n3 = new NVMBNode(3, "node_3");
  n3->AddNeighbor(2);
  n3->AddNeighbor(4);
  n3->AddNeighbor(8);
  N.AddNode(n3);

  NVMBNode* n4 = new NVMBNode(4, "node_4");
  n4->AddNeighbor(3);
  n4->AddNeighbor(6);
  n4->AddNeighbor(9);
  N.AddNode(n4);

  NVMBNode* n5 = new NVMBNode(5, "node_5");
  n5->AddNeighbor(0);
  n5->AddNeighbor(6);
  n5->AddNeighbor(8);
  N.AddNode(n5);

  NVMBNode* n6 = new NVMBNode(6, "node_6");
  n6->AddNeighbor(4);
  n6->AddNeighbor(5);
  N.AddNode(n6);

  NVMBNode* n7 = new NVMBNode(7, "node_7");
  n7->AddNeighbor(1);
  n7->AddNeighbor(2);
  N.AddNode(n7);

  NVMBNode* n8 = new NVMBNode(8, "node_8");
  n8->AddNeighbor(3);
  n8->AddNeighbor(5);
  n8->AddNeighbor(9);
  N.AddNode(n8);

  NVMBNode* n9 = new NVMBNode(9, "node_9");
  n9->AddNeighbor(4);
  n9->AddNeighbor(8);
  N.AddNode(n9);

  return N;
}

Network Network2() {
  Network N(1, false, false);

  // declare the nodes
  NVMBNode* n0 = new NVMBNode(0, "node_0");
  n0->AddNeighbor(1);
  N.AddNode(n0);

  NVMBNode* n1 = new NVMBNode(1, "node_1");
  n1->AddNeighbor(0);
  n1->AddNeighbor(2);
  N.AddNode(n1);

  NVMBNode* n2 = new NVMBNode(2, "node_2");
  n2->AddNeighbor(1);
  n2->AddNeighbor(3);
  N.AddNode(n2);

  NVMBNode* n3 = new NVMBNode(3, "node_3");
  n3->AddNeighbor(2);
  n3->AddNeighbor(4);
  N.AddNode(n3);

  NVMBNode* n4 = new NVMBNode(4, "node_4");
  n4->AddNeighbor(3);
  n4->AddNeighbor(5);
  N.AddNode(n4);

  NVMBNode* n5 = new NVMBNode(5, "node_5");
  n5->AddNeighbor(4);
  n5->AddNeighbor(6);
  N.AddNode(n5);

  NVMBNode* n6 = new NVMBNode(6, "node_6");
  n6->AddNeighbor(5);
  n6->AddNeighbor(7);
  N.AddNode(n6);

  NVMBNode* n7 = new NVMBNode(7, "node_7");
  n7->AddNeighbor(6);
  n7->AddNeighbor(8);
  N.AddNode(n7);

  NVMBNode* n8 = new NVMBNode(8, "node_8");
  n8->AddNeighbor(7);
  N.AddNode(n8);

  return N;
}

Network Network3(string folderName) {

    Network N(1, false, false); 

    NVMBNode* n0 = new NVMBNode(0, folderName, "node_0", MakeRandomNodeNatureVariables("each")); 
    n0->AddNeighbor(1);
    n0->AddNeighbor(5);
    N.AddNode(n0);

    NVMBNode* n1 = new NVMBNode(1, folderName, "node_1", MakeRandomNodeNatureVariables("each")); 
    n1->AddNeighbor(0);
    n1->AddNeighbor(2);
    n1->AddNeighbor(7);
    N.AddNode(n1);

    NVMBNode* n2 = new NVMBNode(2, folderName, "node_2", MakeRandomNodeNatureVariables("each")); 
    n2->AddNeighbor(1);
    n2->AddNeighbor(3);
    n2->AddNeighbor(7);
    N.AddNode(n2);

    NVMBNode* n3 = new NVMBNode(3, folderName, "node_3", MakeRandomNodeNatureVariables("each")); 
    n3->AddNeighbor(2);
    n3->AddNeighbor(4);
    n3->AddNeighbor(8);
    N.AddNode(n3);

    NVMBNode* n4 = new NVMBNode(4, folderName, "node_4", MakeRandomNodeNatureVariables("each")); 
    n4->AddNeighbor(3);
    n4->AddNeighbor(6);
    n4->AddNeighbor(9);
    N.AddNode(n4);

    NVMBNode* n5 = new NVMBNode(5, folderName, "node_5", MakeRandomNodeNatureVariables("each")); 
    n5->AddNeighbor(0);
    n5->AddNeighbor(6);
    n5->AddNeighbor(8);
    N.AddNode(n5);

    NVMBNode* n6 = new NVMBNode(6, folderName, "node_6", MakeRandomNodeNatureVariables("each")); 
    n6->AddNeighbor(4);
    n6->AddNeighbor(5);
    N.AddNode(n6);

    NVMBNode* n7 = new NVMBNode(7, folderName, "node_7", MakeRandomNodeNatureVariables("each")); 
    n7->AddNeighbor(1);
    n7->AddNeighbor(2);
    N.AddNode(n7);

    NVMBNode* n8 = new NVMBNode(8, folderName, "node_8", MakeRandomNodeNatureVariables("each")); 
    n8->AddNeighbor(3);
    n8->AddNeighbor(5);
    n8->AddNeighbor(9);
    N.AddNode(n8);

    NVMBNode* n9 = new NVMBNode(9, folderName, "node_9", MakeRandomNodeNatureVariables("each")); 
    n9->AddNeighbor(4);
    n9->AddNeighbor(8);
    N.AddNode(n9);

    return N; 
}

Network Network4(string folderName, bool mortality) {
    Network N(4, mortality, false); 

    NVMBNode* n0 = new NVMBNode(0, folderName, "node_0", MakeRandomNodeNatureVariables("each")); 
    n0->AddNeighbor(1);
    n0->AddNeighbor(3);
    N.AddNode(n0);

    NVMBNode* n1 = new NVMBNode(1, folderName, "node_1", MakeRandomNodeNatureVariables("each")); 
    n1->AddNeighbor(0);
    n1->AddNeighbor(2);
    N.AddNode(n1);

    NVMBNode* n2 = new NVMBNode(2, folderName, "node_2", MakeRandomNodeNatureVariables("each")); 
    n2->AddNeighbor(1);
    n2->AddNeighbor(3);
    n2->AddNeighbor(5); 
    N.AddNode(n2);

    NVMBNode* n3 = new NVMBNode(3, folderName, "node_3", MakeRandomNodeNatureVariables("each")); 
    n3->AddNeighbor(0);
    n3->AddNeighbor(2);
    n3->AddNeighbor(4); 
    n3->AddNeighbor(5); 
    N.AddNode(n3);

    NVMBNode* n4 = new NVMBNode(4, folderName, "node_4", MakeRandomNodeNatureVariables("each")); 
    n4->AddNeighbor(3);
    n4->AddNeighbor(5); 
    N.AddNode(n4); 

    NVMBNode* n5 = new NVMBNode(5, folderName, "node_5", MakeRandomNodeNatureVariables("each")); 
    n5->AddNeighbor(2);
    n5->AddNeighbor(4); 
    N.AddNode(n5);  

    return N; 
}


Network Network6(string folderName, bool mortality) {
    Network N(6, mortality, false); 

    NVMBNode* n0 = new NVMBNode(0, folderName, "node_0", MakeRandomNodeNatureVariables("each")); 
    n0->AddNeighbor(1);
    N.AddNode(n0);

    NVMBNode* n1 = new NVMBNode(1, folderName, "node_1", MakeRandomNodeNatureVariables("each")); 
    n1->AddNeighbor(0);
    n1->AddNeighbor(2);
    N.AddNode(n1);

    NVMBNode* n2 = new NVMBNode(2, folderName, "node_2", MakeRandomNodeNatureVariables("each")); 
    n2->AddNeighbor(1);
    n2->AddNeighbor(3); 
    N.AddNode(n2);

    NVMBNode* n3 = new NVMBNode(3, folderName, "node_3", MakeRandomNodeNatureVariables("each")); 
    n3->AddNeighbor(2);
    n3->AddNeighbor(4); 
    N.AddNode(n3);

    NVMBNode* n4 = new NVMBNode(4, folderName, "node_4", MakeRandomNodeNatureVariables("each")); 
    n4->AddNeighbor(3);
    n4->AddNeighbor(5); 
    N.AddNode(n4); 

    NVMBNode* n5 = new NVMBNode(5, folderName, "node_5", MakeRandomNodeNatureVariables("each")); 
    n5->AddNeighbor(4); 
    N.AddNode(n5);  

    return N; 
}

/// TODO: delete
Network Network1WithNodesInSelfMode(vector<int> x) {
  Network N = Network1();

  for (int i = 0; i < x.size(); i++) {
      N.contents[x[i]]->SetMode("inactive");
  };

  return N;
}

void SetNetworkNodesActiveStatus(Network* N, set<int> x, string status) {
    assert (status == "active" || status == "inactive");  
    NVMBNode* n;
    for (auto x_: x) {
        (N->contents)[x_]->SetMode(status); 
    } 
}


/*
// description
sets all nodes except for `nodeId` to inactive and runs discovery process.

// arguments
n - a Network at timestamp 0
nodeId - node to set active
*/
void RunNetworkDFSProcessOnNode(Network* N, int nodeId) {

  for (auto it = (N->contents).begin(); it != (N->contents).end(); it++) {
    if (it->first != nodeId) {
      it->second->SetMode("inactive");
    }
  }

  NodeProcessingUnit* npu = (N->contents)[nodeId]->npu;

  N->RunOneTimestamp();

  /// look at NodeStrategos
  auto strategy = npu->GetStrategy();

  // fetch plan for discovery
  auto execPlans = strategy->GetPlansInExecution();
  Plan* p = execPlans[0];

  // run discovery process
  int i = 0;
  while (true) {
    N->RunOneTimestamp();
    auto fp = strategy->GetFinishedPlans();
    auto r1 = find(fp.begin(), fp.end(), p) != fp.end() ;
    i++;
    if (r1) {
      break;
    }
  }
}

/*
// description
runs discovery process on all nodes in Network N

// arguments
N - network is at timestamp 0 with no inactive nodes

*/
void RunNetworkDFSProcessOnAllNodes(Network N) {

  N.RunOneTimestamp();

  // set all nodes to inactive
  for (auto x: N.contents) {
      (x.second)->SetMode("inactive");
  };

    int nr = 0; 
  while (!IsNetworkDiscoveryOff(N)) {
    N.RunOneTimestamp();

    nr++;
  }

  for (auto x: N.contents) {
      (x.second)->SetMode("active");
  };

  cout << "* rounds run:\t" << nr << endl;  
}

void RunNetworkDFSProcessOnActiveNodes(Network* N) { 
    
    N->RunOneTimestamp(); 

    while (!CheckNetworkDiscoveryOff(N)) {
        N->RunOneTimestamp();
    }
}

/*
TODO delete this.
*/
void DisplayNodeBestPaths(NVMBNode* n) {

  NodeProcessingUnit* npu = n->npu;
  auto bestPaths = npu->GetBestPaths();

  /*
  */
  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {

    cout << "PATH INFO: " << (it)->first << endl;
    vector<pair<vector<int>, float>> val = (it)->second;
    for (auto it2: val) {
      cout << "* ";
      DisplayIterable(it2.first);
      cout << "* score: " << it2.second << endl;
      cout << endl;
    }
  }
}

void DisplayNodePlansInfo(Network N, int nodeId) {
  NodeProcessingUnit* npu = N.contents[nodeId]->npu;
  auto strategy = npu->GetStrategy();

  auto plansCache = strategy->GetPlansInCache();
  auto plansInExecution = strategy->GetPlansInExecution();
  auto finishedPlans = strategy->GetFinishedPlans();

  cout << "** Node Plans Data for " << nodeId << endl;

  cout << "PLANS IN CACHE" << endl;
  for (auto it = plansCache.begin(); it != plansCache.end(); it++) {
    (*it)->DisplayInfo();
    cout << endl;
  }

  cout << "PLANS IN EXECUTION" << endl;
  for (auto it2 = plansInExecution.begin(); it2 != plansInExecution.end(); it2++) {
    (*it2)->DisplayInfo();
    cout << endl;
  }

  cout << "FINISHED PLANS" << endl;
  for (auto it3 = finishedPlans.begin(); it3 != finishedPlans.end(); it3++) {
    (*it3)->DisplayInfo();
    cout << endl;
  }
}

/*
*/ 
void DisplayNetworkTimestampUnits(Network N) {
  cout << "TIMESTAMP DATA FOR NETWORK OF TEST" << endl; 

  for (auto x: N.contents) {
    BankUnit* bu = (x.second)->GetBank();
    cout << "\t\t" << "info for node " << x.first << endl;  
    bu->DisplayTimestampUnitHistory(); 
  }
} 


////////////// methods used to make manual plans for testing. 

/* 
*/ 
void LoadPlanForDirectTesting(NVMBNode* n, Plan* p) {

  NodeStrategos* ns = n->GetStrategy(); 
  ns->InsertDefaultPlan(p);  
}

/*
// description 
Makes a direct contract plan for node's default execution. 
*/ 
void MakeContractPlanDefaultForNode(NVMBNode* n, int receiverNode) {
      NodeStrategos* ns = n->GetStrategy(); 
      Plan* p = ns->MakeContractPlan(n->GetId(), n->GetId(), receiverNode, 1.0); 
      ns->InsertDefaultPlan(p); 
}

/*
// description 
Makes a phantom contract plan for node's default execution. 
*/
Plan* MakeContractPhantomPlanDefaultForNode(NVMBNode* n, int senderNode, int receiverNode) {
    NodeStrategos* ns = n->GetStrategy();
    Plan* p = ns->MakeContractPlan(n->GetId(), senderNode, receiverNode, 1.0); 
    (ns->planIndex)++;
    ns->AddPlanToPriority(p);
    return p; 
} 

/// TODO: code below 
void BreakContractPlanDefaultForNode(NVMBNode* n , int contractedNode) {
  // iterate through node contracts and get contract id for target node. 
  NodeStrategos* ns = n->GetStrategy(); 
  int contractId = -1; 
  for (auto x: ns->activeContracts) {
    if (x->receiver == contractedNode) {
      contractId = x->idn; 
      break; 
    }
  }

  // case: contract could not be found
  if (contractId == -1) {
    return; 
  }

  Plan* p = ns->BreakContractPlan(n->GetId(), n->GetId(), contractedNode, 0.0, contractId); 
  ns->AddPlanToPriority(p); 
} 



/// TODO: code for all planTypes 
/*
*/  
void LoadPlanForNodeN(Network N, int n, string planType) {

  NVMBNode* n_ = N.contents[n]; 

  if (planType == "make_contract") {
    // check for available nodes 
    set<int> nidns = n_->GetNodesWithoutContract(); 

    // choose the first if available 
    if (nidns.size() != 0) {
      set<int>::iterator it = nidns.begin(); 
      int target = *it; 

      NodeStrategos* ns = n_->GetStrategy(); 

      Plan* p = ns->MakeContractPlan(n_->GetId(), n_->GetId(), target, 0.0); 
      ns->InsertDefaultPlan(p); 
    }
  } else if (planType == "break_contract") {
    cout << "HAVE NOT CODED" << endl; 
  } else if (planType == "make_bond") {

  } else if (planType == "break_bond") {

  } else if (planType == "discovery") {

  } else {

  }
}

bool IsNetworkDiscoveryOff(Network N) {
  NodeStrategos* ns;
  for (auto x: N.contents) {
    ns = (x.second)->GetStrategy();
    if (ns->discoveryOn) {
      return false;
    }
  }
  return true;
}

bool CheckNetworkDiscoveryOff(Network* N) {
  NodeStrategos* ns;
  for (auto x: N->contents) {
    ns = (x.second)->GetStrategy();
    if (ns->discoveryOn) {
      return false;
    } else {
        (x.second)->SetMode("inactive"); 
    }
  }
  return true;
}

bool CheckNodeTimestampDataForEvent(NVMBNode* n, int timestamp, string eventType) {
    auto bank = n->GetBank(); 
    vector<TimestampUnit*> x =  bank->GetTimestampUnitHistoryAtRange(timestamp, timestamp + 1); 

    for (auto x_: x) {
        if (x_->eventType == eventType) {
            return true; 
        }
    }
    return false; 
}

int GetEventCountForNode(NVMBNode* n, int timestamp, string eventType) {
    /// TRY USING ACCUMULATE
    auto bank = n->GetBank(); 
    vector<TimestampUnit*> tsuv =  bank->GetTimestampUnitHistoryAtRange(timestamp, timestamp + 1); 

    return accumulate(tsuv.begin(), tsuv.end(),
        0.0,
        [eventType](int t, TimestampUnit* r)
        {
            return t + (r->eventType == eventType) ? 1: 0; 
        });
}


map<int,int> GetEventCountPerNodeForTimestamp(Network* N, int timestamp, string eventType) {
    map<int,int> output = map<int,int>(); 

    for (auto x: N->contents) {
        output[x.first] = GetEventCountForNode(x.second, timestamp, eventType); 
    }

    return output; 
}


TimestampUnit* FetchNodeTimestampDataForEvent(NVMBNode* n, int timestamp, string eventType) {
    auto bank = n->GetBank(); 
    vector<TimestampUnit*> x =  bank->GetTimestampUnitHistoryAtRange(timestamp, timestamp + 1); 

    for (auto x_: x) {
        if (x_->eventType == eventType) {
            return x_; 
        }
    }
    return nullptr; 
}


/*
// description 
concatenates best paths from checkpoint to checkpoint 
*/ 
vector<int> GetNodesInPhantomPath(Network* N, vector<int> checkPoints) {
    assert (checkPoints.size() > 1);

    NodeProcessingUnit* npu; 
    vector<int> output = vector<int>();
    
    for (int i = 0; i < checkPoints.size() - 1; i++) {
        npu = N->contents[checkPoints[i]]->npu; 
        auto pathInfo = npu->GetPathForDestination(checkPoints[i + 1]);  
        output.insert(output.end(), pathInfo.first.begin(), pathInfo.first.end()); 
    }
    return output; 
}