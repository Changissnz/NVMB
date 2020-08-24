#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

/*
// description
tests that dummy plan is executed correctly.
Note that tests requires global variables to remain unchanged.

Test details:
- Network: 2
- Dummy Plan: 0 proposes to break-bond 5-6
*/
/// TODO: bank worth is wrong 
TEST(Node_PostDiscovery_BreakBondDummyPlan_Execution_Test, Node_PostDiscovery_BreakBondDummyPlan_Execution_Test_Correct) {

    Network N = Network2();

    cout << "* TIMEZA 1" << endl; 
    for (auto n: N.contents) {
        cout << "node " << n.first << " bank value " << n.second->GetWorth() << endl; 
    }

    auto s = N.contents[0]->GetStrategy();

    // run DFS on all nodes
    RunNetworkDFSProcessOnAllNodes(N);

    cout << "* TIMEZA 2" << endl; 
    for (auto n: N.contents) {
        cout << "node " << n.first << " bank value " << n.second->GetWorth() << endl; 
    }

    // display timestamp data for one of the nodes 
    auto b = N.contents[3]->GetBank(); 
    cout << "DISPLAY NODE 3" << endl; 
    b->DisplayTimestampUnitHistory(); 

    ///////////////////////////////////////////////////////////////////// 
  // insert dummy plan
  Plan* dummyPlan = new Plan(s->planIndex, "phantom", "break_bond", 5, 6, 0.0, -1);
  (s->planIndex)++;
  s->AddPlanToPriority(dummyPlan);

  // set all nodes except for 0 to inactive and defaultResponseCF to true
  for (auto x: N.contents) {
      if (x.first != 0) {
        (x.second)->SetMode("inactive");
        (x.second)->SetDefaultResponseCF(true, true);
      }
  };

  // collect bank values for each node
  map<int,float> bankValuesBefore;
  for (auto c: N.contents) {
    bankValuesBefore[c.first] = c.second->GetWorth();
  }

  // run the plan and check bank values after
  N.RunOneTimestamp();

  map<int,float> bankValuesAfter;
  for (auto c: N.contents) {
    bankValuesAfter[c.first] = c.second->GetWorth();
  }

  /// TODO: delete this. 
  cout << "POST DISC BBP" << endl; 
  for (auto x: bankValuesBefore) {
    cout << x.first << "\tBEFORE: " << x.second << "\tAFTER: " << bankValuesAfter[x.first] << endl;  
  }

    
  // check bank values before and after 
  ASSERT_EQ(bankValuesBefore[0] > bankValuesAfter[0], true);
  ASSERT_EQ(bankValuesBefore[1] < bankValuesAfter[1], true);
  ASSERT_EQ(bankValuesBefore[2] < bankValuesAfter[2], true);
  ASSERT_EQ(bankValuesBefore[3] < bankValuesAfter[3], true);
  ASSERT_EQ(bankValuesBefore[4] < bankValuesAfter[4], true);

  ASSERT_EQ(bankValuesAfter[2] > bankValuesAfter[1], true);
  ASSERT_EQ(bankValuesAfter[3] > bankValuesAfter[2], true);
  ASSERT_EQ(bankValuesAfter[4] > bankValuesAfter[3], true);

  // sender and receiver nodes
  ASSERT_EQ(bankValuesBefore[5] > bankValuesAfter[5], true);
  //ASSERT_EQ(bankValuesBefore[6] > bankValuesAfter[6], true); 

  /// TODO: delete this. uncomment to display the bank values before and after
  /*
  cout << "BANK DATA" << endl;
  for (auto c: bankValuesBefore) {
    cout << "bank for node " << c.first << "| before " << c.second << " after " << bankValuesAfter[c.first] << endl;
  }
  */

  // check that edge 5-6 no longer exists
  ASSERT_EQ(N.EdgeExists(5,6, false), false);

  // test direct proposition plan
  Plan* dummyPlan2 = new Plan(s->planIndex, "direct", "break_bond", 0, 1, 0.0, -1);
  (s->planIndex)++;
  s->AddPlanToPriority(dummyPlan2);
  N.RunOneTimestamp();
  ASSERT_EQ(N.EdgeExists(0,1, false), false);
}

/*
// description
tests that dummy plan is executed correctly.
Note that tests requires global variables to remain unchanged.

Test details:
- Network: 2
- Dummy Plan: 0 proposes to make-bond 3-7
*/
TEST(Node_PostDiscovery_MakeBondDummyPlan_Execution_Test, Node_PostDiscovery_MakeBondDummyPlan_Execution_Test_Correct) {

  Network N = Network2();

  auto s = N.contents[0]->GetStrategy();

  // run DFS on all nodes
  RunNetworkDFSProcessOnAllNodes(N);

  // insert dummy plan
  Plan* dummyPlan = new Plan(s->planIndex, "phantom", "make_bond", 3, 7, 0.0, -1);
  (s->planIndex)++;
  s->AddPlanToPriority(dummyPlan);

  // set all nodes except for 0 to inactive and defaultResponseCF to true
  for (auto x: N.contents) {
      if (x.first != 0) {
        (x.second)->SetMode("inactive");
        (x.second)->SetDefaultResponseCF(true, true);
      }
  };

  // run the plan and check bank values after
  N.RunOneTimestamp();

  ASSERT_EQ(N.EdgeExists(3,7, false), true);

  Plan* dummyPlan2 = new Plan(s->planIndex, "direct", "make_bond", 0, 5, 0.0, -1);
  (s->planIndex)++;
  s->AddPlanToPriority(dummyPlan2);

  N.RunOneTimestamp();

  ASSERT_EQ(N.EdgeExists(0,5, false), true);
}

/*
after discovery, runs a dummy plan in `node default accept` mode and checks for
completion
*/
TEST(Node_PostDiscovery_MakeContractDummyPlan_Execution_Test, Node_PostDiscovery_MakeContractDummyPlan_Execution_Test_Correct) {

  Network N = Network2();
  auto s = N.contents[0]->GetStrategy();

  // run DFS on all nodes
  RunNetworkDFSProcessOnAllNodes(N);

  // insert dummy plan
  Plan* dummyPlan = new Plan(s->planIndex, "phantom", "make_contract", 3, 7, 0.0, -1);
  (s->planIndex)++;
  s->AddPlanToPriority(dummyPlan);

  // set all nodes except for 0 to inactive and defaultResponseCF to true
  for (auto x: N.contents) {
      if (x.first != 0) {
        (x.second)->SetMode("inactive");
        (x.second)->SetDefaultResponseCF(true,true);
      }
  };

  // run the plan and check bank values after
  N.RunOneTimestamp(true);

  // check active contracts
  ///N.DisplayExistingContracts();
  ASSERT_EQ(N.contents[3]->IsContractedWith(7) && N.contents[7]->IsContractedWith(3), true);
}

/*
// description 
similar to the tests above, except checks for the cases in which
  node 3 rejects
  node 7 rejects
*/
TEST(Node_PostDiscovery_MakeContractDummyPlan_SenderReceiverReject_Execution_Test, Node_PostDiscovery_MakeContractDummyPlan_SenderReceiverReject_Execution_Test_Correct) {

  /// CASE 1: node 3 rejects
  Network N = Network2();
  auto s = N.contents[0]->GetStrategy();

  // run DFS on all nodes
  RunNetworkDFSProcessOnAllNodes(N);
    cout << "[X] TIMESTAMP:\t" << N.contents[3]->GetTimestamp() << endl;  

  // insert dummy plan
  Plan* dummyPlan = new Plan(s->planIndex, "phantom", "make_contract", 3, 7, 0.0, -1);
  (s->planIndex)++;
  s->AddPlanToPriority(dummyPlan);

  // set all nodes except for 0 to inactive and defaultResponseCF to true
  for (auto x: N.contents) {
      if (x.first != 0) {
        (x.second)->SetMode("inactive");
      }
  };

  N.contents[3]->SetDefaultResponseCF(true, false);

  // collect bank values before
  map<int,float> bankValuesBefore;
  for (auto c: N.contents) {
    bankValuesBefore[c.first] = c.second->GetWorth();
  }

    N.RunOneTimestamp(true);
    auto ts = N.contents[3]->GetTimestamp(); 
    cout << "[X] TIMESTAMP:\t" << ts << endl;  

    map<int,float> bankValuesAfter;
    for (auto c: N.contents) {
        bankValuesAfter[c.first] = c.second->GetWorth();
    }

    // check that contract exist 
    auto b1 = N.contents[3]->IsContractedWith(7); 
    auto b2 = N.contents[7]->IsContractedWith(3); 

    ASSERT_EQ(!N.contents[3]->IsContractedWith(7) && !N.contents[7]->IsContractedWith(3), true);

    // check timestamp units for nodes 0 and 3
    BankUnit* bu0 = N.contents[0]->GetBank();
    BankUnit* bu3 = N.contents[3]->GetBank();

    int t0 = bu0->GetNumberOfTimestamps(); 
    int t1 = bu3->GetNumberOfTimestamps(); 

    ASSERT_EQ(bu0->GetNumberOfTimestamps(), 13);
    ASSERT_EQ(bu3->GetNumberOfTimestamps(), 13);
}
