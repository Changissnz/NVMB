#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

/// TODO: nodes 3 and 9 still need to be complete. 
/* 
// description 
for Node 0 of Network1, tests for making contract b/t nodes 
3, 6, and 9. 

Default response setting has not been set for non-proposer nodes. 

Test details: 
- Network 1, node 0 
- targets are 3,6,9
- Dummy plan: 0 proposes direct contract to each of nodes 3,6,9 
- Assertions: 
        - correct number of finished plans for node 6 
        - 
*/ 
TEST(N1n0_RunDiscovery_ExecuteManualContractInput_Check, CheckManConInp) {

    vector<int> x = vector<int>{1,2,3,4,5,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    RunNetworkDFSProcessOnNode(&N, 0); 

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true); 
    MakeContractPlanDefaultForNode(n0, 6); 
    N.RunOneTimestamp(2);

    NodeStrategos* ns1 = N.contents[0]->GetStrategy();
    ASSERT_EQ(2, ns1->GetFinishedPlans().size()); 
}

/// TODO: nodes 6,9 
/* 
// description 
for Node 0 of Network1, tests for making phantom contract b/t nodes 
3, 6, and 9. 

Test details: 
- Network 1, node 0 
- targets are 3,6,9
- Dummy plan: 0 proposes phantom contract to each of nodes 3,6,9
- default response (non-prop): NONE 
- Assertions: 
        - correct timestamp info. 
        - correct number of finished plans 
*/ 
TEST(N1n0_RunDiscovery_ExecuteManualContractInput2_Check, CheckManConInp2) {

    vector<int> x = vector<int>{1,2,3,4,5,6,7,8,9};

    Network N = Network1WithNodesInSelfMode(x);
    RunNetworkDFSProcessOnActiveNodes(&N); 


    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true);

    Plan* p1 = MakeContractPhantomPlanDefaultForNode(n0, 2, 3); 

    // set active status 
    set<int> inact = {1,2,3,4,5,6,7,8,9};  

    SetNetworkNodesActiveStatus(&N, inact, "inactive");  

    set<int> act = {0};  
    SetNetworkNodesActiveStatus(&N, act, "active"); 

    // execute plan 
    N.RunOneTimestamp(1); 

    // assert correct timestamp 
    bool n0Correct = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "make_contract"); 
    ASSERT_EQ(n0Correct, true);

    bool n2Correct = CheckNodeTimestampDataForEvent(N.contents[2], N.contents[2]->GetTimestamp() - 1, "make_contract"); 
    ASSERT_EQ(n2Correct, true);

    bool n3Correct = CheckNodeTimestampDataForEvent(N.contents[3], N.contents[3]->GetTimestamp() - 1, "make_contract"); 
    ASSERT_EQ(n3Correct, false);

    // assert correct plan counts for each node. 
    auto plan0 = N.contents[0]->GetPlanCounts(); 

    vector<int> ans0 = {0,0,0,2}; 

    
    auto planss = (N.contents[0]->GetStrategy())->GetPlansInExecution(); 

    for (auto p: planss) {
        p->DisplayInfo(); 
    }

    /*
    ASSERT_EQ(plan0, ans0); 

    auto plan2 = N.contents[2]->GetPlanCounts(); 
    vector<int> ans2 = {0,0,0,0}; 
    ASSERT_EQ(plan2, ans2); 

    cout << "PLAN COUNTS FOR NODE 3" << endl; 
    auto plan3 = N.contents[3]->GetPlanCounts(); 
    DisplayIterable(plan3); 
    vector<int> ans3 = {0,0,0,0}; 
    ASSERT_EQ(plan3, ans3); 
    */ 
}

/*
// description 
for Node 0 of Network1, tests for making contract b/t nodes 
3, 6, and 9. 

Test details: 
- Network 1, node 0 
- targets are 3,6,9
- Dummy plan: 0 proposes direct contract to each of nodes 3,6,9
- default response (non-prop): 3 T 6 F 9 T  
- Assertions: 
        - correct timestamp info. 
        - correct number of finished plans 
        - correct response effects (of accepting/rejecting)
            - contract existence 
        - correct cf_tax on path.  
*/ 
TEST(N1n0_RunDiscovery_ExecuteManualContractInput3_Check, CheckManConInp3) {
    
    /// CASE 1: 
    // declare and run network before manual contract 
    vector<int> x = vector<int>{1,2,4,5,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    
    RunNetworkDFSProcessOnActiveNodes(&N); 

    cout << "CURRENT TIME " << N.contents[0]->GetTimestamp() << endl; 

    // set manual contract for node 0 
    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true);
    MakeContractPlanDefaultForNode(n0, 3); 

    // set default response for receiver 3 
    N.contents[3]->SetDefaultResponseCF(true,true); 

    set<int> inact = {1,2,3,4,5,6,7,8,9};  
    SetNetworkNodesActiveStatus(&N, inact, "inactive");  

    set<int> act = {0};  
    SetNetworkNodesActiveStatus(&N, act, "active");   

    // assert correct plan counts pre-timestamp 
    auto x0 = N.contents[0]->GetPlanCounts();
    ASSERT_EQ(x0[3], 1); 

    auto x3 = N.contents[3]->GetPlanCounts(); 
    ASSERT_EQ(x3[3], 1);

    N.RunOneTimestamp(1); 

    bool r0 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "make_contract");
    ASSERT_EQ(r0, true); 
    bool r3 = CheckNodeTimestampDataForEvent(N.contents[3], N.contents[3]->GetTimestamp() - 1, "make_contract");  
    ASSERT_EQ(r3, true); 

    // observe timestamp info for nodes 0 a nd 3 
    /// TODO: delete, for displaying timestamp data 
    cout << "BANK 00" << endl; 
    auto b = N.contents[0]->GetBank(); 
    b->DisplayTimestampUnitHistory();

    cout << "BANK 33" << endl; 
    b = N.contents[3]->GetBank(); 
    b->DisplayTimestampUnitHistory();
    
    /// 
    vector<int> other = {1,2,4,5,6,7,8,9}; 
    for (auto o : other) {
        ASSERT_EQ(false, CheckNodeTimestampDataForEvent(N.contents[o], N.contents[o]->GetTimestamp() - 1, "make_contract"));  
    }

    x0 = N.contents[0]->GetPlanCounts(); 
    ASSERT_EQ(x0[3], 2);  

    x3 = N.contents[3]->GetPlanCounts(); 
    ASSERT_EQ(x3[3], 1); 

} 

/*
// description 
- same as above except phantom contracts instead. 

Test details: 
- Network 1, node 0 
- targets are 3,6,9
- Dummy plan: 0 proposes direct contract to each of nodes 3,6,9
- default response (non-prop): 3 T 6 F 9 T  
- Assertions: 
        - correct timestamp info. 
        - correct number of finished plans 
        - correct response effects (of accepting/rejecting)
            - contract existence 
        - correct cf_tax on path.  
*/ 
TEST(N1n0_RunDiscovery_ExecuteManualContractInput4_Check, CheckManConInp4) {

    //// CASE 1: proposer 0 sender 5 receiver 3, 5 rejects 
    /// 0 5 8 3 
    vector<int> x = vector<int>{1,2,4,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    
    RunNetworkDFSProcessOnActiveNodes(&N); 

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true);
    MakeContractPhantomPlanDefaultForNode(n0, 5, 3); 

    N.contents[3]->SetDefaultResponseCF(true,true); 
    N.contents[5]->SetDefaultResponseCF(true,false); 

    // 
    auto x0 = N.contents[0]->GetPlanCounts();
    auto x3 = N.contents[3]->GetPlanCounts(); 
    auto x5 = N.contents[5]->GetPlanCounts(); 
       
    set<int> act = {0}; 
    SetNetworkNodesActiveStatus(&N, act, "active");   

    N.RunOneTimestamp(1);

    // display timestamp data
    cout << "BANK 00" << endl; 
    auto b = N.contents[0]->GetBank(); 
    b->DisplayTimestampUnitHistory(); 

    cout << "BANK 33" << endl; 
    b = N.contents[3]->GetBank(); 
    b->DisplayTimestampUnitHistory(); 

    cout << "BANK 55" << endl; 
    b = N.contents[5]->GetBank(); 
    b->DisplayTimestampUnitHistory(); 

    // check that nodes 0 and 5 recorded event, and 3 did not 
    auto n0Correct = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1 ,"make_contract"); 
    ASSERT_EQ(n0Correct, true); 

    auto n5Correct = CheckNodeTimestampDataForEvent(N.contents[5], N.contents[5]->GetTimestamp() - 1 ,"make_contract"); 
    ASSERT_EQ(n5Correct, true);

    auto n3Correct = CheckNodeTimestampDataForEvent(N.contents[3], N.contents[3]->GetTimestamp() - 1 ,"make_contract"); 
    ASSERT_EQ(n3Correct, false);

    // comm flare check 
    auto n0c = N.contents[0]->GetCommFlareCounts(); 
    auto n3c = N.contents[3]->GetCommFlareCounts(); 
    auto n5c = N.contents[5]->GetCommFlareCounts(); 

    vector<int> ans = {0,0,0}; 
    ASSERT_EQ(n0c == ans, true); 

    ASSERT_EQ(n3c == ans, true); 

    ASSERT_EQ(n5c == ans, true);
} 

/*
// description 

Network erases an edge used in routing comm flare. 
Tests for comm flare routing failure. 
*/ 
TEST(N1n0_RunDiscovery_ExecuteManualContractInput5_Check, CheckManConInp5) {

    vector<int> x = vector<int>{1,2,4,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    
    RunNetworkDFSProcessOnActiveNodes(&N); 

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true); 
    MakeContractPlanDefaultForNode(n0, 3);  
    
    N.contents[3]->SetDefaultResponseCF(true,true); 

    // print out the best paths of 
    auto pathInfo = N.GetValidPathFromSourceToTarget(0, 3);  

    // know that by default settings, default path is
    //  0 5 8 3 
    // remove edge 5 8 
    N.BreakEdge(5, 8, 0); 

    set<int> act = {0}; 
    SetNetworkNodesActiveStatus(&N, act, "active");   
    N.RunOneTimestamp(1); 

    /// ERROR: bank requires timestamp for contract 
    // check correct timestamp 
    bool n0Correct = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "make_contract"); 
    ASSERT_EQ(n0Correct, true); 

    bool n3Correct = CheckNodeTimestampDataForEvent(N.contents[3], N.contents[3]->GetTimestamp() - 1, "make_contract"); 
    ASSERT_EQ(n3Correct, false); 

    // check best path validity for node 0 
    int remv = (N.contents[0]->npu)->RemoveRecordsOfEdge(NodePairToString(5,8));  

    ASSERT_EQ(remv, 0);
}

/*
// description 
runs network discovery and manual contract execution. 


Then checks for correct 
CF_tax update (bank and timestamp). 

Test details: 
    GO HERE. 
*/ 
TEST(N1n0_ContractExecutionAttributes, CheckContractAttributesCorrect) {
    // set up network and run initial dfs on it 
    vector<int> x = vector<int>{1,2,4,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);
    RunNetworkDFSProcessOnActiveNodes(&N); 

    // set n0 attrib. 
    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true); 
    MakeContractPlanDefaultForNode(n0, 3);  
    
    // set n3 attrib. 
    N.contents[3]->SetTestModeSpec("make_bond"); 
    N.contents[3]->SetDefaultResponseCF(true,true); 

    // print out the best paths of 
    auto pathInfo = N.GetValidPathFromSourceToTarget(0, 3);  

    cout << "PATH INFORENCIO" << endl; 
    DisplayIterable(pathInfo.first); 

    set<int> act = {0,3}; 
    SetNetworkNodesActiveStatus(&N, act, "active");   

    N.RunOneTimestamp(1); 

    // check for cf tax for each non-checkpoint node 
    for (int i = 1; i < pathInfo.first.size() - 1; i++) {
        ASSERT_EQ(true, CheckNodeTimestampDataForEvent(N.contents[pathInfo.first[i]], N.contents[pathInfo.first[i]]->GetTimestamp() - 1, "cf tax")); 
    }

    // check other 
    ASSERT_EQ(true, CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "make_contract")); 
    ASSERT_EQ(true, CheckNodeTimestampDataForEvent(N.contents[3], N.contents[3]->GetTimestamp() - 1, "make_contract")); 

    // run another timestamp and check that contract updates nodes 0 and 3 correctly 
    N.RunOneTimestamp(1); 

    cout << "BP ##" << endl; 
    auto path03 = N.GetValidPathFromSourceToTarget(0, 3); 
    DisplayIterable(path03.first); 

    cout << "BANK 00" << endl; 
    auto b = N.contents[0]->GetBank(); 
    b->DisplayTimestampUnitHistory(); 

    auto c1 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity"); 

    cout << "BANK 33" << endl; 
    b = N.contents[3]->GetBank(); 
    b->DisplayTimestampUnitHistory(); 
    auto c2 = CheckNodeTimestampDataForEvent(N.contents[3], N.contents[3]->GetTimestamp() - 1, "contract activity"); 

    ASSERT_EQ(c1, true); 
    ASSERT_EQ(c2, true); 

    /// TODO: check the other nodes in path for contract tax 

    // get the best paths for each 
    auto pT = N.GetValidPathFromSourceToTarget(0,3);  
    auto pD = N.GetValidPathFromSourceToTarget(3,0);

    // get the counts of each node given the paths 
    map<int,int> counts = map<int,int>(); 
    set<int> non3 = {1,2,4,5,6,7,8,9};  

    //for (auto s: pT.first) {
    // all nodes except first and last 
    for (int i = 1; i < pT.first.size() - 1; i++) {
        counts[pT.first[i]]++; 
    }

    for (int i = 1; i < pD.first.size() - 1; i++) {
        counts[pD.first[i]]++; 
    }

    for (auto x: counts) {
        auto res = CheckNodeTimestampDataForEvent(N.contents[x.first], N.contents[x.first]->GetTimestamp() -1, "contract tax"); 
        ASSERT_EQ(res, true); 
    }
}

/*
// description 
Runs DFS on Network1, create contract between nodes, and check for correct cf tax 
timestamp and impact during contract creation, and during contract activity. 
*/ 
TEST(NetworkNodeContractCFTest, ContractCFTestCorrect) {

    vector<int> x = vector<int>{1,2,3,4,5,7,8};
    Network N = Network1WithNodesInSelfMode(x);
    RunNetworkDFSProcessOnActiveNodes(&N); 

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true);

    Plan* p1 = MakeContractPhantomPlanDefaultForNode(n0, 9, 6);

    ///// display below 
    auto p09 = (N.contents[0]->npu)->GetPathForDestination(9);     
    auto p96 = (N.contents[9]->npu)->GetPathForDestination(6); 

    /// get all nodes 
    map<int,int> counts = map<int,int>(); 
    set<int> non3 = {1,2,4,5,6,7,8,9};  

    // all nodes except first and last 
    for (int i = 1; i < p09.first.size() - 1; i++) {
        counts[p09.first[i]]++; 
    }

    for (int i = 1; i < p96.first.size() - 1; i++) {
        counts[p96.first[i]]++; 
    }

    set<int> inact = {1,2,3,4,5,7,9};  
    SetNetworkNodesActiveStatus(&N, inact, "inactive");  

    set<int> act = {0,8,6};  
    SetNetworkNodesActiveStatus(&N, act, "active");   

    // execute timestamp for contract creation 
    N.RunOneTimestamp(false); 

    // check network and node timestamp
    // check the counts of cf tax for each 
    auto counts2 = GetEventCountPerNodeForTimestamp(&N, N.timestamp - 1, "cf tax"); 
    
    // 5 8 4 
    for (auto c: counts) {
        ASSERT_EQ(counts2[c.first] >= c.second, true); 
    }

    N.RunOneTimestamp(false);     
    // check network and node timestamp 
    p96 = (N.contents[9]->npu)->GetPathForDestination(6); 
    
    // check the counts of cf tax for each 
    auto counts3 = GetEventCountPerNodeForTimestamp(&N, N.timestamp - 1, "contract tax"); 
    
    for (int i = 1; i < p96.first.size() - 1; i++) {
        ASSERT_EQ(counts3[p96.first[i]] >= 1, true); 
    }
}