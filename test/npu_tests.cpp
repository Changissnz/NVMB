#include "class_declarations.hpp"
#include <gtest/gtest.h>
using namespace std;

TEST(NPUShortenPath, NPUShortenPathCorrect) {

  NodeProcessingUnit* npu = BlankNPU();
  vector<int> path = vector<int>{1,4,5,6,7,8,10,12,13,14};

  // bond.first at [0]
  pair<int,int> p1 = make_pair(1, 10);
  auto np1 = npu->ShortenPath(path, p1);
  auto np1a = vector<int>{1,10,12,13,14};
  ASSERT_EQ(np1, np1a);

  pair<int,int> p2 = make_pair(13, 1);
  auto np2 = npu->ShortenPath(path, p2);
  auto np2a = vector<int>{1,13,14};
  ASSERT_EQ(np2, np2a);

  pair<int,int> p3 = make_pair(6, 12);
  auto np3 = npu->ShortenPath(path, p3);
  auto np3a = vector<int>{1,4,5,6,12,13,14};
  ASSERT_EQ(np3, np3a);
}