#include "network1.h"
#include <gtest/gtest.h>

/*
checks bank value before after dfs execution 
*/ 
TEST(NetworkBankDFS, NetworkBankDFSTest) {
  vector<int> x = vector<int>{0,1,2,3,4,6,7,8,9};
  Network N = Network1WithNodesInSelfMode(x);
  NodeProcessingUnit* npu = N.contents[5]->npu;

  BankUnit* bu = N.contents[5]->GetBank();
  ASSERT_EQ(bu->currency, 100.0);

  N.RunOneTimestamp();
  N.contents[5]->SetMode("inactive");

  NodeStrategos* s = (N.contents[5])->GetStrategy();
  Plan* p = s->GetPlansInExecution()[0];

  while (true) {
    N.RunOneTimestamp();
    auto fp = s->GetFinishedPlans();
    auto r1 = find(begin(fp), end(fp), p);
    if (r1 != fp.end()) {
      break;
    }
  }

  ASSERT_EQ(bu->currency, 80.0);
}
