#include "network1.h"
#include <gtest/gtest.h>
#include<vector>

/*
test does not account for bank
*/
TEST(BaseNetworkDeclarationTest, BaseNetworkDeclared) {
  Network N(1, false);

  N.AddNode(1);
  N.AddNode(2);

  ASSERT_EQ(N.AddNode(1), false);
  ASSERT_EQ(N.Size(), 2);
}

/*
runs two `processDF` on Node 5 and checks for held flares
*/
TEST(NetworkNPUProcess1DF, NPUProcessGood) {
  Network N = Network1();

  // do discovery for node 5
  NVMBNode* n5 = N.contents[5];
  n5->InitializeDiscovery();

  NodeProcessingUnit* npu = n5->npu;
  DiscoveryFlare* df = *(npu->GetOwnedFlares().begin());
  npu->ProcessDF();
  ASSERT_EQ(npu->GetOwnedFlares().size(), 1);
  npu->Update();
  ASSERT_EQ(npu->GetOwnedFlares().size(), 3);

  // assign held flares and check relevant nodes
  N.AssignHeldDFAll();

  ASSERT_EQ(npu->GetHeldFlares().size(), 1); // node 5

  NVMBNode* n0 = N.contents[0];
  npu = n0->npu;
  ASSERT_EQ(npu->GetHeldFlares().size(), 1); // node 0

  NVMBNode* n6 = N.contents[6];
  npu = n6->npu;
  ASSERT_EQ(npu->GetHeldFlares().size(), 1); // node 6

  NVMBNode* n8 = N.contents[8];
  npu = n8->npu;
  ASSERT_EQ(npu->GetHeldFlares().size(), 1); // node 8
}

/*
// description
tests for 1 timestamp on active node 5. Results should include a
Discovery plan
*/
TEST(NetworkRunOneTimestamp, NetworkRunOneTimestampCorrect) {
  vector<int> x = vector<int>{0,1,2,3,4,6,7,8,9};
  Network N = Network1WithNodesInSelfMode(x);
  NodeProcessingUnit* npu = N.contents[5]->npu;

  FUTURE_PLAN_SIZE = 3;
  ACTIVATION_THRESHOLD = 1.01;
  N.RunOneTimestamp();

  /// look at NodeStrategos
  auto strategy = npu->GetStrategy();

  // observe plan variables
  // display the contents of plan modecache
  auto plansCache = strategy->GetPlansInCache();
  auto plansInExecution = strategy->GetPlansInExecution();

  //-----------------------------------------------------------
  // uncomment below to display the contents of plan execution
  /*
  for (auto it = plansInExecution.begin(); it != plansInExecution.end(); it++) {
    cout << "plan exec:\t" << (*it)->planType << "\t" << (*it)->score << endl;
    ///cout << "\t" << "plan info:\t" << (*it)->sourceNodeIdentifier << "\t" << (*it)->targetNodeIdentifier << endl;
  }

  for (auto it2 = plansCache.begin(); it2 != plansCache.end(); it2++) {
    cout << "plan cache:\t" << (*it2)->planType << "\t" << (*it2)->score << endl;
    cout << "\t" << "plan info:\t" << (*it2)->sourceNodeIdentifier << "\t" << (*it2)->targetNodeIdentifier << endl;
  }
  */
  //-----------------------------------------------------------
  ASSERT_EQ(plansInExecution.size(), 1);
  ASSERT_EQ(plansInExecution[0]->planType, "discovery");
  ASSERT_EQ(plansCache.size(), 2);
  ASSERT_EQ(strategy->GetFinishedPlans().size(), 0);
  return;
}


/*
// description
runs a Discovery Process on Node 5 and checks for correct end variables,
which include:
- max path length
- shortest destination distances
*/
TEST(NetworkDiscoveryProcessNode5, NetworkDiscoveryProcessNode5Correct) {
    vector<int> x = vector<int>{0,1,2,3,4,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    NodeProcessingUnit* npu = N.contents[5]->npu;

    cout << "[0] CURRENT TIME " << N.timestamp << "\t" << N.contents[5]->GetTimestamp() << endl;

    N.RunOneTimestamp();
    cout << "[1] CURRENT TIME " << N.timestamp << "\t" << N.contents[5]->GetTimestamp() << endl;

    N.contents[5]->SetMode("inactive");

    /// look at NodeStrategos
    auto strategy = npu->GetStrategy();

    FUTURE_PLAN_SIZE = 3;
    ACTIVATION_THRESHOLD = 0.5;

    // fetch plan for discovery
    auto execPlans = strategy->GetPlansInExecution();
    Plan* p = execPlans[0];

    // run discovery process
    int i = 0;
    while (true) {
        N.RunOneTimestamp();
        auto fp = strategy->GetFinishedPlans();
        auto r1 = find(fp.begin(), fp.end(), p) != fp.end() ;
        i++;
        if (r1) {
        break;
        }
    }

    ASSERT_EQ(strategy->GetFinishedPlans().size(), 1);

    // check for correct max path length
    ASSERT_EQ(npu->maxPathLength, 8);

    // check for correct shortest path distances
    auto spl = npu->GetBestPaths();

    ASSERT_EQ(spl[0][0].first.size() - 1, 1);
    ASSERT_EQ(spl[1][0].first.size() - 1, 2);
    ASSERT_EQ(spl[2][0].first.size() - 1, 3);
    ASSERT_EQ(spl[3][0].first.size() - 1, 2);
    ASSERT_EQ(spl[4][0].first.size() - 1, 2);
    ASSERT_EQ(spl[5][0].first.size() - 1, 0);
    ASSERT_EQ(spl[6][0].first.size() - 1, 1);
    ASSERT_EQ(spl[7][0].first.size() - 1, 3);
    ASSERT_EQ(spl[8][0].first.size() - 1, 1);
    ASSERT_EQ(spl[9][0].first.size() - 1, 2);


    ///// UNCOMMENT BELOW FOR DETAILS 
    /*
    cout << "CURRENT TIME " << N.timestamp << "\t" << N.contents[5]->GetTimestamp() << endl;

    auto tsuh = (N.contents[5]->GetBank())->GetTimestampUnitHistory();

    for (auto s: tsuh) {
        cout << "TIMETH " << s.first << endl; 

        for (auto ss: s.second) {
            ss->DisplayInfo(); 
        }

        cout << "\n\n\n" << endl; 

    }
    N.contents[5]->WriteOutToFile(false); 
    */ 
}

/*
// description
tests for no error when running two activities on active node 5, which involves
- node 5 completes discovery process. 
- node 5 performs another timestamp afterwards. 
*/
TEST(NetworkNode5MoveInitial, NetworkNode5MoveInitialCorrect) {
    vector<int> x = vector<int>{0,1,2,3,4,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    NodeProcessingUnit* npu = N.contents[5]->npu;

    ACTIVATION_THRESHOLD = 0.80;
    N.RunOneTimestamp();

    NodeStrategos* strategy = npu->GetStrategy();

  //-----------------------------------------------------------------------------
  // uncomment to display variable info  
  /*
  cout << "plan cache size:\t" << strategy->GetPlansInCache().size() << endl;
  cout << "plan execution size:\t" << strategy->GetPlansInExecution().size() << endl;
  cout << "finished plan size:\t" << strategy->GetFinishedPlans().size() << endl;
  */ 
  //-----------------------------------------------------------------------------

    ASSERT_EQ(0, strategy->GetPlansInCache().size());
    ASSERT_EQ(1, strategy->GetPlansInExecution().size());
    ASSERT_EQ(0, strategy->GetFinishedPlans().size());

    Plan* p = strategy->GetPlansInExecution()[0];
    int i = 0;
    while (true) {
        N.RunOneTimestamp();
        auto fp = strategy->GetFinishedPlans();
        auto r1 = find(fp.begin(), fp.end(), p) != fp.end() ;
        i++;
        if (r1) {
        break;
        }
    }

    vector<Plan*> rx1 = strategy->GetPlansInExecution(); 

    cout << "[1] DA PLANZ" << endl; 
    for (auto y: rx1) {
        y->DisplayInfo(); 
    }

    N.RunOneTimestamp();

    /*
    cout << "---------------------" << endl;
    cout << "AFTER" << endl; 
    cout << "plan cache size:\t" << strategy->GetPlansInCache().size() << endl;
    cout << "plan execution size:\t" << strategy->GetPlansInExecution().size() << endl;
    cout << "finished plan size:\t" << strategy->GetFinishedPlans().size() << endl;
    */ 
    
    /// TODO: check below assertions; changed from 2,0 to 1,1 

    vector<Plan*> rx = strategy->GetPlansInExecution(); 

    cout << "[2] DA PLANZ" << endl; 
    for (auto y: rx) {
        y->DisplayInfo(); 
    }

    ASSERT_EQ(strategy->GetPlansInExecution().size(), 0);
    ASSERT_EQ(strategy->GetFinishedPlans().size(), 2);
    ASSERT_EQ(strategy->discoveryOn, false);
}

/*
// description
tests DFS process on all nodes

Displays results, no assertions.
*/
TEST(NetworkAllNodesDFS, NetworkAllNodesDFS_Correct) {

  Network N = Network1();

  N.RunOneTimestamp();

  // set all nodes to inactive
  for (auto x: N.contents) {
      (x.second)->SetMode("inactive");
  };

  while (!IsNetworkDiscoveryOff(N)) {
    N.RunOneTimestamp();
  }

  // UNCOMMENT BELOW TO SEE BEST PATHS FOR EACH NODE
  /*
  for (auto x: N.contents) {
    cout << "BEST PATHS FOR NODE: " << x.first << endl;
    DisplayNodeBestPaths(x.second);
    cout << "===================================================" << endl;
  }
  */
}

