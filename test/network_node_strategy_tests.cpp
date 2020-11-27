#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

/// ABOUT:
/*
file is used to demonstrate different choices for node strategy functions. 
*/

/// NOTE: check when discovery move is activated
TEST(Network__DemonstrateCompetitionMeasures, CompetitionMeasures) {
    vector<int> x = vector<int>{1,2,3,4,0,6,7,8,9}; 
    Network N = Network1WithNodesInSelfMode(x); 
    
    // competition measures
    NodeStrategos* ns = N.contents[5]->GetStrategy();
    ///ns->competitionMeasure = "growth-performance"; 
    ///ns->competitionMeasure = "growth-relative";
    ns->competitionMeasure = "number of contracts";  
    ///ns->competitionMeasure = "number of contracts relative";
    ///ns->competitionMeasure = "path score to max";
    ///ns->competitionMeasure = "path score to mean";
    N.Run(2, 20); 
}

/// TODO: 
TEST(Network__DemonstrateContractMeasures, ContractMeasures) {

}

/// TODO: 
TEST(Network__DemonstrateNewContractMeasures, NewContractMeasures) {

}

/// TODO: 
TEST(Network__DemonstrateBondAdvantageMeasures, BondAdvantageMeasures) {

}

/// TODO:
TEST(Network__DemonstrateBondDeletionMeasures, BondDeletionMeasures) {

}