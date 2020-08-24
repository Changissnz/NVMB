#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

/*
*/ 
TEST(N1n0_RunContractProcedureThenBreak, CorrectBreakOfCont) {

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

/// TODO: make better display of each node's move per timestamp 
/// TODO: run 100 rounds and check data 
//----------------------------------------------------------------
