// this file will contain bond/contract functions
#include "bond_functions.hpp"
#include <iostream>

/*
// description
only proposer will incur this cost. Based on node connectivity, maximally
connect node, and bank. The more connected a Node is, the more expensive the
next bond will cost.

Arguments are with respect to the sender node.

// arguments
connectivity - number of edges belonging to the Node of interest.
maxConnectivity - maximum possible edges a Node can be connected to.
currency - non-negative
*/
float BondModificationCost1(float nodeRank, int connectivity, int maxConnectivity, float currency) {
  return -(float(connectivity) /maxConnectivity * currency) * nodeRank;
}

/*
// description
the sender and the receiver will each incur this cost, calculated by their
set of arguments.

Lower-ranking nodes will incur a higher proportional cost, and the more total
number of nodes there are, the cheaper the bond will be.

// arguments
nodeRank - traffic rank of node
connectivity - number of neighbor nodes
totalNumNodes - total number of nodes
currency - value of node
//
*/
float BondModificationCost2(float nodeRank, float connectivity, float totalNumNodes, float currency) {
  float cost = currency * (1 - nodeRank);
  return -cost * connectivity / totalNumNodes;
}

/*
// description
Does not take into consideration the cost of making a bond.
Calculates the advantageousness of adding a bond by their path scores.
*lower sums are better*

// arguments
- affectedPaths := best paths affected by bond addition
- competitors := set of competitors
- coeff1 := variable coefficient value 1
- coeff2 := variable coefficient value 2
*/
float BondAddAdvantageMeasure_PathScoresPromise(vector<pair<int, pair<float,float>>> affectedPaths,
    set<int> competitors, float coeff1, float coeff2) {

  float comp = 0.0, nc = 0.0;
  for (auto it = affectedPaths.begin(); it != affectedPaths.end(); it++) {
    if (competitors.find(it->first) != competitors.end()) {
      comp += ((it->second.first - it->second.second) * coeff1);
    } else {
      nc += ((it->second.first - it->second.second) * coeff2);
    }
  }

  return ZeroDiv(comp, nc + comp, 1.0);
}

/*
// description
Calculates the advantageousness of adding a bond by considering their skipped
scores.

// arguments
- pathInfo := each element in vector is a pair with
              [0] destination node
              [1] nodes in previous paths which will no longer be taken in new paths
                  because of new edge.
- competitors := set of competitors
- coeff1 := competition if direct bond, o.w. negotiation
- coeff2 := greed
*/
float BondAddAdvantageMeasure_PathSkippedNodesPromise(vector<pair<int, vector<int>>> pathInfo,
      set<int> competitors, float coeff1, float coeff2) {

  float comp = 0.0, nc = 0.0;
  vector<int> path;
  for (auto it = pathInfo.begin(); it != pathInfo.end(); it++) {
    path = it->second;
    for (auto it2 = path.begin(); it2 != path.end(); it2++) {
      if (competitors.find(*it2) != competitors.end()) {
        comp += coeff1;
      } else {
        nc += coeff2;
      }
    }
  }
  return ZeroDiv(comp, nc + comp, 1.0);
}

/*
// description
outputs a map corresponding to path elements, in which key is path node, and value
is competitor status
*/
map<int, bool> FindCompetitorIndices(vector<int> path, set<int> competitors) {

  map<int, bool> indices;
  for (int i = 0; i < path.size(); i++) {
    if (competitors.find(path[i]) != competitors.end()) {
      indices[i] = true;
    } else {
      indices[i] = false;
    }
  }

  return indices;
}

/*
// description
score in range [0,1]

description of scoring:
1 => (index2 - index1) && (both are competitors)
(1 - negotiationMeasure) => (index2 - index1) && (one is competitor, other is not)
  *higher the negotiation, lower the score*

// arguments
- index1 := of the first element with respect to path
- index2 := of the second element with respect to path
- size := number of edges of path
- negotiationMeasure := float
*/
float GetPairwiseCompetitorIndexScore(int index1, int index2, int size, bool isComp1, bool isComp2, float coeff1, float coeff2) {

  if (!isComp1 && !isComp2) {
    return 0.0;
  }

  int distance = abs(index1 - index2); // 1 is min, size is max
  float coeff = (isComp1 && isComp2) ? (coeff1 + coeff2) : coeff1;
  return coeff * (size - distance + 1) / float(size);
}

/*
// description
Helper method for filemethod<BondBreakAdvantageMeasure_PathSkippedNodes> below..
Calculates the CompetitorBreakScore for each split path.
*/
pair<float,float> GetCompetitorBreakScore(pair<int, pair<vector<int>, vector<int>>> brokenPathUnit, set<int> competitors, float coeff1, float coeff2) {

  auto p1 = brokenPathUnit.second.first;
  auto p2 = brokenPathUnit.second.second;

  map<int,bool> ind1 = FindCompetitorIndices(p1, competitors);
  map<int,bool> ind2 = FindCompetitorIndices(p2, competitors);
  bool c1, c2;
  int totalSz = p1.size() + p2.size() - 1; // number of edges of path
  float output = 0.0, pOutput = 0.0, x = 0.0, y = 0.0;

  for (int i = 0; i < p1.size(); i++) {
    c1 = ind1[i];
    for (int j = 0; j < p2.size(); j++) {
      c2 = ind2[j];
      x = GetPairwiseCompetitorIndexScore(i, p1.size() + j, totalSz, c1, c2, coeff1, coeff2);
      y = GetPairwiseCompetitorIndexScore(i, p1.size() + j, totalSz, true, true, coeff1, coeff2);
      output += x;
      pOutput += y;
    }
  }

  return make_pair(output, pOutput);
}

/*
// description
this measures the negative impact of breaking a bond by estimating the disruption
in competitor's paths to other nodes. The greater the degree in disruption in
competitor paths to other nodes, the lower the score.
*/
float BondBreakAdvantageMeasure_PathSkippedNodes(vector<pair<int, pair<vector<int>, vector<int>>>> brokenPathInfo, set<int> competitors,
                  float coeff1, float coeff2) {

    float numerator = 0.0, denumerator = 0.0;
    for (auto it = brokenPathInfo.begin(); it != brokenPathInfo.end(); it++) {
      auto breakScore = GetCompetitorBreakScore(*it, competitors, coeff1, coeff2);
      numerator += breakScore.first;
      denumerator += breakScore.second;
    }
    return 1.0 - ZeroDiv(numerator, denumerator, 0.0);
}

/*
// description
This measures the negative impact of breaking a bond by considering score differences
between affected best paths and their alternative paths.
*/
/// TODO: this measure does not take into account the entire network (all nodes),
///       only the affected paths. This could pose a problem.
float BondBreakAdvantageMeasure_PathScores(vector<pair<int, pair<float, float>>> brokenPathInfo, float worstPathScore, float coeff1, float coeff2) {
  float totalWorstPathScore = worstPathScore * brokenPathInfo.size();
  float cumulativeChange = 0.0;
  for (auto it = brokenPathInfo.begin(); it != brokenPathInfo.end(); it++) {
    cumulativeChange += (it->second.second - it->second.first);
  }
  return ZeroDiv(cumulativeChange, totalWorstPathScore, 1.0);
}
