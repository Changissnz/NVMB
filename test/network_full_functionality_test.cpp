/*
file will contain code to test accuracy of executing network functionality. 
*/ 

#include "network1.h"
#include "network_generation.hpp"
#include <gtest/gtest.h>
#include <time.h> 
using namespace std;


/*
all nodes except 0 are in self-mode 
*/ 
TEST(NetWorkRunNTimestampsFullDisplay, NetWorkRunNTimestampsFullDisplay) {
  vector<int> x = vector<int>{1,2,3,4,5,6,7,8,9};
  Network N = Network1WithNodesInSelfMode(x);

  // run n rounds 
    // SEE HERE: set value below for number of rounds 
  int i = 3; 
  for (int j = 0; j < i; j++) {
    N.RunOneTimestamp(1); 
    PrintSpacer();
  }
}

TEST(NetworkRunTimestampWithNodeDeletion, RunWithNodeDeletionCorrect) {

    cout << "FULL RUN TEST" << endl; 
    Network N = Network1(); 
    N.mortality = true; 
    int i = 100; 

    N.RunOneTimestamp(1);

    // iterate through and show timestamp 
    for (auto c: N.contents) {
        auto bu = (c.second)->GetBank(); 
        cout << "TIMESTAMP FOR NODE " << c.first << endl; 
        bu->DisplayTimestampUnitHistory();
        cout << "SIZE:\t" << bu->GetCurrentTimestampInfo().size() << endl; 
        cout << "\n\n" << endl;
    }

    /// TODO: delete below. 
        /*
    for (auto x: N.contents) {
        cout << "SUMMARY FOR NODE " << x.first << endl; 
        (x.second)->SummarizeLastNTimestamps(1000); 
    } 
        */ 

    // mark 0, 3 as bankrupt, run and check. 
    auto b0 = N.contents[0]->GetBank(); 
    b0->currency = float(-106725); 

    auto b3 = N.contents[3]->GetBank(); 
    b3->currency = float(-106725); 
    
    cout << "RUNNING ANOTHER" << endl; 
    N.RunOneTimestamp(1);  
}

/*
// description 
Initializes a contract b/t nodes 0 and 9. Bankrupts node 9 and checks if contract 
is still effective. 
*/ 
TEST(NetworkRunTimestampNodeDeletionAndContract, NodeDeletionContractCorrect) {

    vector<int> x = vector<int>{1,2,3,4,5,6,7,8}; 
    Network N = Network1WithNodesInSelfMode(x);
    RunNetworkDFSProcessOnActiveNodes(&N); 
    N.mortality = true; 

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true); 
    MakeContractPlanDefaultForNode(n0, 9); 

    // set mode for node 9 
    N.contents[9]->SetDefaultResponseCF(true,true); 
    N.contents[9]->SetTestModeSpec("make_bond"); 

    set<int> act = {0,9}; 
    SetNetworkNodesActiveStatus(&N, act, "active"); 

    // run timestamp for making contract 
    N.RunOneTimestamp(1);
    n0->SetToManualInput(false); 

    // run timestamp for contract activity 
    auto bankBefore = N.contents[0]->GetWorth();  
    N.RunOneTimestamp(1); 
    auto bankAfter = N.contents[0]->GetWorth();     
    (N.contents[0]->GetBank())->DisplayTimestampUnitHistory(); 
    ASSERT_EQ(bankAfter > bankBefore, true); 

    bool e1 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity"); 
    ASSERT_EQ(e1,true);  

    // bankrupt node 9 and check 
    auto b9 = N.contents[9]->GetBank(); 
    b9->currency = float(-106725);
    n0->SetTestModeSpec("discovery"); 
    N.RunOneTimestamp(); 

        // below contracts will not activate, but node 0 does not know node 9 is dead. 
    e1 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity"); 
    
    auto tsuc = FetchNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity"); 
    ASSERT_EQ(tsuc != nullptr,true);
    ASSERT_EQ(tsuc->impact, 0);  

    N.RunOneTimestamp(); 

    e1 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity"); 
    tsuc = FetchNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity"); 
    ASSERT_EQ(tsuc != nullptr,true);
    ASSERT_EQ(tsuc->impact, 0);   
    return; 
}



/*
tests for correct recording of timestamp data to file using NodeFileReader 
*/ 
TEST(NetworkRunRecord, CorrectNetworkRunRecord) {
    //+!+ 
    ////return; 

    /// TODO: declare a new network with a special folder names
    Network N = Network1(); 
    N.mortality = true; 

    // check file write 

    /// TODO: for method `CheckNullNodes`, have to code writeToFile for dead node. 
    /// run n rounds 
    int numRounds = 100; 
    
    while (N.IsAlive() && numRounds > 0) {
        cout << "** round " << 100 - numRounds + 1 << endl; 
        N.RunOneTimestamp(1);
        numRounds -= 1; 
    }

    cout << "rounds run:\t" << 100 - numRounds << endl; 
    DisplayIterable(N.ActiveNodes()); 
    for (auto c: N.contents) {
        cout << "node " << c.first << "\tbank " << c.second->GetWorth() << endl; 
    }   

    /// TODO: clean up and write-to-file. 
    cout << "** deceased" << endl; 
    DisplayIterable(N.deceased); 
}

/*
*/ 
TEST(NetworkRunRecord_RandomNodeNatures, RunRecordRandomNatureCorrect) {
    srand (100);

    Network N1 = Network3("folder1");
    N1.mortality = true; 
    N1.Run(0, 1000);
    cout << "[0] rounds run: " << N1.timestamp << endl; 
    cout << "** deceased:\t"; 
    DisplayIterable(N1.deceased);

    Network N2 = Network3("folder2"); 
    N2.mortality = true; 
    N2.Run(0, 1000);
    cout << "[1] rounds run: " << N2.timestamp << endl; 
    cout << "** deceased:\t"; 
    DisplayIterable(N2.deceased); 

    Network N3 = Network3("folder3"); 
    N3.mortality = true; 
    N3.Run(0, 1000); 
    cout << "[2] rounds run: " << N3.timestamp << endl; 
    cout << "** deceased:\t"; 
    DisplayIterable(N3.deceased);  
}

TEST(Network_DFS_FlareCount, FlareCountCorrect) {

    // case 1: Network4 
    Network N4 = Network6("folder4", true);

        // set all nodes to inactive except for node 0
    set<int> inactive = {1,2,3,4,5};   
    
    SetNetworkNodesActiveStatus(&N4, inactive, "inactive"); 

        // set node 0 to discovery mode 

        // # run 3 timestamps 
    for (int i = 0; i < 3; i++) {
        N4.RunOneTimestamp(0);
        cout << "[" << i << "] flare count for 0" << endl; 
        DisplayIterable(N4.contents[0]->GetFlareCounts());
    }

    /*
    // case 2: Network5 
    Network N5 = Network5("folderDFS_0", true);  

        // set all nodes to inactive except for node 0
    set<int> inactive = {1,2,3,4,5,6,7,8};   
    SetNetworkNodesActiveStatus(&N5, inactive, "inactive"); 

        // set node 0 to discovery mode 
    N5.contents[0]->SetTestModeSpec("discovery"); 

    N5.RunOneTimestamp(0); 
    cout << "[0] flare count for 0" << endl; 
    DisplayIterable(N5.contents[0]->GetFlareCounts());
    N5.RunOneTimestamp(0);
    cout << "[1] flare count for 0" << endl; 
    DisplayIterable(N5.contents[0]->GetFlareCounts());

    N5.RunOneTimestamp(0); 
    cout << "[2] flare count for 0" << endl; 
    DisplayIterable(N5.contents[0]->GetFlareCounts());
    */ 
}

/*
this method is used to collect data for randomly-generated networks. 
*/ 
TEST(Network__CollectData__, CollectingData) {

    srand (873);


     // generate the random network.  
    Network N = GenerateRandomNetwork(0, true, "folder14", 20, 0.5, "fode"); 

    // set all nodes to discoveryType limit mode 
    for (auto c: N.contents) {
        c.second->SetDiscoveryType("limit");        
    }

    N.PrerunProcess(); 

    time_t t1; 
    time(&t1);  /* get current time; same as: timer = time(NULL)  */
    N.Run(2,-1);   
    
    time_t t2; 
    time(&t2);

    cout << "NUM ROUNDS " << N.timestamp << endl;  
    cout << "ELAPSED " << difftime(t2, t1) << endl; 
    printf ("%.f seconds", difftime(t2,t1));

    cout << "node timestamp deaths" << endl;
    DisplayIterable(N.deceased);
}

TEST(Network__CollectData_RandomNetwork_RunRandomNodeNatureVars, CollectingData2) {

     // generate the random network.  
    Network N = GenerateRandomNetwork(0, true, "rnnv_0", 20, 0.5, "node"); 

    // set all nodes to discoveryType limit mode 
    for (auto c: N.contents) {
        c.second->SetDiscoveryType("limit");        
    }

    time_t t1; 
    time(&t1);  /* get current time; same as: timer = time(NULL)  */
    time_t t2; 
    time(&t2);

    cout << "NUM ROUNDS " << N.timestamp << endl;  
    cout << "ELAPSED " << difftime(t2, t1) << endl; 
    printf ("%.f seconds", difftime(t2,t1));
    cout << "node timestamp deaths" << endl;
    DisplayIterable(N.deceased);
}

TEST(Network__RunTestGraph, RunSuccess) { 
  Network N(1, false);

  NVMBNode* n0 = new NVMBNode(0, "node_0");
  NVMBNode* n1 = new NVMBNode(1, "node_1");
  N.AddNode(n0); 
  N.AddNode(n1); 

  N.AddEdge(0,1,false); 

  N.Run(1, 100); 
}