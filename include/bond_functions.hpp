#ifndef BONDFUNCTIONS_HPP
#define BONDFUNCTIONS_HPP

#define BOND_CREATION_COST_RATIO 2

#include <math.h>
#include "node_strategos.hpp"
#include <algorithm>
using namespace std;

float BondModificationCost1(float nodeRank, int connectivity, int maxConnectivity, float currency);
float BondModificationCost2(float nodeRank, float connectivity, float totalNumNodes, float currency);

float BondAddAdvantageMeasure_PathScoresPromise(vector<pair<int, pair<float,float>>> affectedPaths,
    set<int> competitors, float coeff1, float coeff2);

float BondAddAdvantageMeasure_PathSkippedNodesPromise(vector<pair<int, vector<int>>> pathInfo,
      set<int> competitors, float coeff1, float coeff2);

float BondDestructionCost1();
map<int, bool> FindCompetitorIndices(vector<int> path, set<int> competitors);

float GetPairwiseCompetitorIndexScore(int index1, int index2, int size, bool isComp1, bool isComp2, float coeff1, float coeff2);
pair<float,float> GetCompetitorBreakScore(pair<int, pair<vector<int>, vector<int>>> brokenPathUnit, set<int> competitors, float coeff1, float coeff2);

float BondBreakAdvantageMeasure_PathSkippedNodes(vector<pair<int, pair<vector<int>, vector<int>>>> brokenPathInfo, set<int> competitors,
                  float coeff1, float coeff2);
float BondBreakAdvantageMeasure_PathScores(vector<pair<int, pair<float, float>>> brokenPathInfo, float worstPathScore, float coeff1, float coeff2);

#endif
