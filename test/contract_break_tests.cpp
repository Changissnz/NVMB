#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

/*
*/ 
TEST(N1n0_RunContractProcedureThenBreak, CorrectBreakOfCont_Direct) {

    vector<int> x = vector<int>{1,2,3,4,5,6,7,8};
    Network N = Network1WithNodesInSelfMode(x);

    RunNetworkDFSProcessOnActiveNodes(&N); 

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true); 
    MakeContractPlanDefaultForNode(n0, 9); 

    // set mode for node 9 
    N.contents[9]->SetDefaultResponseCF(true,true); 
    N.contents[9]->SetTestModeSpec("make_bond"); 

    // print path info  
    auto pathInfo = N.GetValidPathFromSourceToTarget(0,9);
    cout << "PATH 0-9" << endl; 
    DisplayIterable(pathInfo.first); 

    set<int> act = {0,9}; 
    SetNetworkNodesActiveStatus(&N, act, "active");   

    // run one for contract execution //
    N.RunOneTimestamp(false);

    // run 2 rounds with contract on //
    N.RunOneTimestamp(false); 
    N.RunOneTimestamp(false); 

    // check for contract activity 
    bool b1 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity");  
    bool b2 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 2, "contract activity");  
    cout << "CON-ACT CHECK n0: " << b2 << b1 << endl; 

    b1 = CheckNodeTimestampDataForEvent(N.contents[9], N.contents[9]->GetTimestamp() - 1, "contract activity");  
    b2 = CheckNodeTimestampDataForEvent(N.contents[9], N.contents[9]->GetTimestamp() - 2, "contract activity");  
    cout << "CON-ACT CHECK n9: " << b2 << b1 << endl; 

    // break contract round // 
    BreakContractPlanDefaultForNode(N.contents[0], 9); 
    N.RunOneTimestamp(0); 
    
    //b1 = CheckNodeTimestampDataForEvent(N.contents[0], N.contents[0]->GetTimestamp() - 1, "contract activity");  
    /// TODO: make assertion below. 
    cout << "SUMMMM 0" << endl; 
    N.contents[0]->SummarizeLastNTimestamps(2); 

    cout << "SUMMMM 9" << endl; 
    N.contents[9]->SummarizeLastNTimestamps(2); 

    // post-break contract //
    cout << "RUN ONE NO CONTRACT" << endl; 
    N.RunOneTimestamp(1); // n0 makes invalid move 

    cout << "[2] SUMMMM 0" << endl; 
    N.contents[0]->SummarizeLastNTimestamps(1); 

    cout << "[2] SUMMMM 9" << endl; 
    N.contents[9]->SummarizeLastNTimestamps(1); 

}

TEST(N1n1_RunContractProcedureThenBreak, CorrectBreakOfCont_Direct) {

    // break contract phantom    
    vector<int> x = vector<int>{3,4,5,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);

        /*
    cout << "ACTIVE MODES" << endl;
    for (auto c: N.contents) {
        cout << "PROCESS < " << c.first << " " << (c.second)->GetMode() << endl; 
    }
        */


    RunNetworkDFSProcessOnActiveNodes(&N);
    /*
    cout << "F U" << endl;
    NodeStrategos* ns = (N.contents[0]->npu)->GetStrategy();
    ns->DisplayAllPlans();
    */
    set<int> act = {0,1,2};  
    SetNetworkNodesActiveStatus(&N,act,"active");
    /*
    cout << "ACTIVE MODES" << endl;
    for (auto c: N.contents) {
        cout << "PROCESS < " << c.first << " " << (c.second)->GetMode() << endl; 
    }
    */

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true);
    Plan* p1 = MakeContractPhantomPlanDefaultForNode(n0, 1, 2);
    
    N.RunOneTimestamp(0);
    BankUnit* bu = N.contents[2]->GetBank();

    NVMBNode* n1 = N.contents[1];
    NVMBNode* n2 = N.contents[2];

    n1->SetToManualInput(true);
    n2->SetToManualInput(true);
    N.RunOneTimestamp(0);    
    N.RunOneTimestamp(0);    

    BreakContractPlanDefaultForNode(n1,2);
    N.RunOneTimestamp(0);    

    N.RunOneTimestamp(0);    

    bu->DisplayTimestampUnitHistory();

    bool n1Correct = CheckNodeTimestampDataForEvent(N.contents[1], N.contents[1]->GetTimestamp() - 1, "contract activity");
    bool n2Correct = CheckNodeTimestampDataForEvent(N.contents[2], N.contents[2]->GetTimestamp() - 1, "contract activity");
    
    ASSERT_EQ(n1Correct, false);
    ASSERT_EQ(n2Correct, false);
}

TEST(N1n1_RunContractProcedureThenBreak_Case2, CorrectBreakOfCont_Phantom) {
    

    // break contract phantom    
    vector<int> x = vector<int>{3,4,5,6,7,8,9};
    Network N = Network1WithNodesInSelfMode(x);

    /*
    cout << "ACTIVE MODES" << endl;
    for (auto c: N.contents) {
        cout << "PROCESS < " << c.first << " " << (c.second)->GetMode() << endl; 
    }
    */

    RunNetworkDFSProcessOnActiveNodes(&N);

    set<int> act = {0,1,2};  
    SetNetworkNodesActiveStatus(&N,act,"active");

    NVMBNode* n0 = N.contents[0]; 
    n0->SetToManualInput(true);
    Plan* p1 = MakeContractPhantomPlanDefaultForNode(n0, 1, 2);
    
    N.RunOneTimestamp(0);
    BankUnit* bu = N.contents[0]->GetBank();

    NVMBNode* n1 = N.contents[1];
    NVMBNode* n2 = N.contents[2];

    n1->SetToManualInput(true);
    n2->SetToManualInput(true);
    N.RunOneTimestamp(0);    
    N.RunOneTimestamp(0);    

    // iterate through n1 contracts
    NodeStrategos* ns1 = (n1->npu)->GetStrategy();

    for (auto c: ns1->activeContracts) {
        c->DisplayInfo();
    }

    NodeStrategos* ns2 = (n2->npu)->GetStrategy();

    for (auto c: ns2->activeContracts) {
        c->DisplayInfo();
    }
    
    n0->SetToManualInput(true);
    BreakContractPhantomPlanDefaultForNode(n0,n1,2);
    n1->SetToManualInput(true);
    n2->SetToManualInput(true);

    N.RunOneTimestamp(0);
    N.RunOneTimestamp(0);    

    bool n1Correct = CheckNodeTimestampDataForEvent(N.contents[1], N.contents[1]->GetTimestamp() - 1, "contract activity");
    bool n2Correct = CheckNodeTimestampDataForEvent(N.contents[2], N.contents[2]->GetTimestamp() - 1, "contract activity");

    ASSERT_EQ(n1Correct,0);
    ASSERT_EQ(n2Correct,0);
}