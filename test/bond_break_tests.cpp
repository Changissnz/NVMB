#include "network1.h"
#include <gtest/gtest.h>
using namespace std;

TEST(NPU_BestPathsAffectedByNewBondScoresTest, NPU_BestPathsAffectedByNewBondScoresTestCorrect) {

  Network N = Network2();

  // set these variables before calling method
  FUTURE_PLAN_SIZE = 3;
  ACTIVATION_THRESHOLD = 0.5;

  RunNetworkDFSProcessOnNode(&N, 0);

  // get arguments for this method
  NVMBNode* n = N.contents[0];
  NodeStrategos* strategy = n->GetStrategy();

  NodeProcessingUnit* npu = n->npu;

  set<int> affectedNodes = set<int>{};
  set<int> actualAffectedNodes = set<int>{6,7,8};
  vector<pair<int, pair<float,float>>>  affected = npu->BestPathsAffectedByNewBondScores(make_pair(1,6));
  for (auto x: affected) {
    affectedNodes.insert(x.first);
  }

  // checks that affected nodes are correct
  ASSERT_EQ(affectedNodes, actualAffectedNodes);
}

/*
// description
Checks that method BondBreakAdvantageMeasure_PathSkippedNodes runs according to
form.
*/
TEST(BondBreakAdvantageMeasure_SkippedNodesTest, BondBreakAdvantageMeasure_SkippedNodesTestCorrect) {

  Network N = Network2();

    // set these variables before calling method
  FUTURE_PLAN_SIZE = 3;
  ACTIVATION_THRESHOLD = 0.5;

  auto strategy = N.contents[5]->GetStrategy();
  auto npu = N.contents[5]->npu;

  // run DFS on all nodes
  RunNetworkDFSProcessOnAllNodes(N);

  auto brokenPathInfo = npu->BestPathsAffectedByBondBreakSkippedNodes(make_pair(5,6));
  set<int> rc1 = set<int>{0,1,2,8};
  float m1 = BondBreakAdvantageMeasure_PathSkippedNodes(brokenPathInfo, rc1, 1.0, 1.0);
    //
  set<int> rc2 = set<int>{0,1,2};
  float m2 = BondBreakAdvantageMeasure_PathSkippedNodes(brokenPathInfo, rc2, 1.0, 1.0);
    //
  set<int> rc3 = set<int>{6,7,8};
  float m3 = BondBreakAdvantageMeasure_PathSkippedNodes(brokenPathInfo, rc3, 1.0, 1.0);
    //
  set<int> rc4 = set<int>{0,1,2,6};
  float m4 = BondBreakAdvantageMeasure_PathSkippedNodes(brokenPathInfo, rc4, 1.0, 1.0);
    //
  float m5 = BondBreakAdvantageMeasure_PathSkippedNodes(brokenPathInfo, rc3, 1.0, 0.0);
    //
  ASSERT_EQ(m1 > m4, true);
  ASSERT_EQ(m2 == 1.0, true);
  ASSERT_EQ(m3 == 0.5, true);
  ASSERT_EQ(m5 == 0.0, true);
}

/*
// description
tests the method GetCompetitorBreakScore
*/
TEST(GetPairwiseCompetitorIndexScoreTest, GetPairwiseCompetitorIndexScoreTestCorrect) {
  Network N = Network2();

    // set these variables before calling method
  FUTURE_PLAN_SIZE = 3;
  ACTIVATION_THRESHOLD = 0.5;

  auto strategy = N.contents[5]->GetStrategy();
  auto npu = N.contents[5]->npu;

  // run DFS on all nodes
  RunNetworkDFSProcessOnAllNodes(N);

  auto brokenPathInfo = npu->BestPathsAffectedByBondBreakSkippedNodes(make_pair(5,6));

  /// test case 1: no competitors, C = 1.0, G = 0.0
  set<int> rc2 = set<int>();
  // first element: 5 | 6
  pair<float,float> out1 = GetCompetitorBreakScore(brokenPathInfo[0], rc2, 1.0, 0.0);
  // second element: 5 | 6 7
  pair<float,float> out2 = GetCompetitorBreakScore(brokenPathInfo[1], rc2, 1.0, 0.0);
  // third element: 5 | 6 7 8
  pair<float,float> out3 = GetCompetitorBreakScore(brokenPathInfo[2], rc2, 1.0, 0.0);
  float r1 = out1.first / out1.second;
  float r2 = out2.first / out2.second;
  float r3 = out3.first / out3.second;
  ASSERT_EQ(r1 == r2 && r2 == r3, true);

  /// test case 2: no competitors, C = 1.0, G = 1.0
  pair<float,float> out4 = GetCompetitorBreakScore(brokenPathInfo[0], rc2, 1.0, 1.0);
  // second element: 5 | 6 7
  pair<float,float> out5 = GetCompetitorBreakScore(brokenPathInfo[1], rc2, 1.0, 1.0);
  // third element: 5 | 6 7 8
  pair<float,float> out6 = GetCompetitorBreakScore(brokenPathInfo[2], rc2, 1.0, 1.0);
  r1 = out4.first / out4.second;
  r2 = out5.first / out5.second;
  r3 = out6.first / out6.second;
  ASSERT_EQ(r1 == r2 && r2 == r3, true);

  /// test case 2: all are competitors, C = 1.0, G = 0.0
  set<int> rc3 = set<int>{6,7,8,9};
  pair<float,float> out7 = GetCompetitorBreakScore(brokenPathInfo[0], rc3, 1.0, 0.0);
  // second element: 5 | 6 7
  pair<float,float> out8 = GetCompetitorBreakScore(brokenPathInfo[1], rc3, 1.0, 0.0);
  // third element: 5 | 6 7 8
  pair<float,float> out9 = GetCompetitorBreakScore(brokenPathInfo[2], rc3, 1.0, 0.0);
  r1 = out7.first / out7.second;
  r2 = out8.first / out8.second;
  r3 = out9.first / out9.second;
  ASSERT_EQ(r1 == 1.0 && r1 == r2 && r2 == r3, true);

  /// test case 3: all are competitors, C = 1.0, G = 1.0
  pair<float,float> out10 = GetCompetitorBreakScore(brokenPathInfo[0], rc3, 1.0, 1.0);
  // second element: 5 | 6 7
  pair<float,float> out11 = GetCompetitorBreakScore(brokenPathInfo[1], rc3, 1.0, 1.0);
  // third element: 5 | 6 7 8
  pair<float,float> out12 = GetCompetitorBreakScore(brokenPathInfo[2], rc3, 1.0, 1.0);
  r1 = out10.first / out10.second;
  r2 = out11.first / out11.second;
  r3 = out12.first / out12.second;
  ASSERT_EQ(r1 == 0.5 && r1 == r2 && r2 == r3, true);
}
