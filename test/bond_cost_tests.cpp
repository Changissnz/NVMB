#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

/*
// description
Displays the cost of breaking or adding a bond using a dummy plan

Test details:
- Network: 2
- Dummy Plan: 0 proposes to break-bond 5-6
- Assertions: 
  - relativity tests for bond modification costs for nodes 0,5,6 
*/
TEST(Node_BreakBondCost_DisplayTest, Node_BreakBondCost_DisplayTestCorrect) {

  Network N = Network2();

    // set these variables before calling method
  FUTURE_PLAN_SIZE = 3;
  ACTIVATION_THRESHOLD = 0.5;

  auto s = N.contents[0]->GetStrategy();

  // run DFS on all nodes
  RunNetworkDFSProcessOnAllNodes(N);

  // set connectivity
  N.SetMaxConnectivity();

  // make dummy plan
  Plan* dummyPlan = new Plan(s->planIndex, "phantom", "break_bond", 5, 6, 0.0, -1);

  // convert plan to CF
  NodeProcessingUnit* npu = N.contents[0]->npu;
  CommFlare* cf = npu->ProduceCommFlare(dummyPlan);

  // calculate costs
  auto b = N.SetBondModCosts(cf);

  // relative tests
  ASSERT_EQ(cf->costs[0] < cf->costs[5] && cf->costs[0] < cf->costs[6], true);
  ASSERT_EQ(cf->costs[5] > cf->costs[6], true);
}
