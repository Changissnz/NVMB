/*
this file contains code to run flares and calculate best best paths
*/
#include "node_processing_unit.hpp"
using namespace std;

set<int> NodeProcessingUnit::GetKnownNodes() {

  set<int> output;

  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {
    output.insert(it->first);
  }

  output.erase(nodeOwnerIdentifier);
  return output;
}

/*
competition positively correlates to shorter paths
greed positively correlates to less risk.

// arguments
competition && greed: normalized to sum to 1
*/
float NodeProcessingUnit::EdgeScore(float edgeRisk, float competition, float greed) {
    return competition + (edgeRisk * greed);
}

/*
// description
calculates the score of path used for bestPaths variable.
Considers the two attributes `normedC` (competition) and `normedG` (greed)
*/
float NodeProcessingUnit::GetPathScore(vector<int> path) {

  int sz = path.size() - 1;
  string edge;
  float risk;
  float score = 0;
  for (int i = 0; i < sz; i++) {
    edge = NodePairToString(path[i], path[i + 1]);
    risk = NodeStrategos::edgeRisks[edge];
    score += EdgeScore(risk, strategy->normedC, strategy->normedG);
  }
  return score;
}

/*
// description
used to calculate path's actual risk score (with regards to greed) and
potential risk (maximum risk score possible given path)

// output
pair<float,float> : [0] is actual risk of choosing path, [1] is potential risk
*/
pair<float,float> NodeProcessingUnit::GetPathRiskPotentialScores(vector<int> path) {

  int sz = path.size() - 1;
  float actualRisk = 0.0, potRisk = 0.0;
  string edge;
  float risk;

  for (int i = 0; i < sz; i++) {
    edge = NodePairToString(path[i], path[i + 1]);
    risk = NodeStrategos::edgeRisks[edge];
    actualRisk += (risk * strategy->normedG);
    potRisk += strategy->normedG;
  };

  return make_pair(actualRisk, potRisk);
}

/// TODO
// use this with caution. may result in sync errors
void NodeProcessingUnit::UpdatePathScores(int nodeId) {
  auto bp = bestPaths[nodeId];
  float score;
  for (int i = 0; i != bp.size(); i++) {
    score = GetPathScore(bp[i].first);
    bp[i].second = score;
  }
  // TODO: sort bp from least to greatest
  sort(bp.begin(), bp.end(),
    [](pair<vector<int>, float> p1, pair<vector<int>, float> p2) {
      return p1.second < p2.second;
    });
  bestPaths[nodeId] = bp;
}

/*
// description
calculates the cumulative risk of taking the path.
score is not influenced by index of edge within path (edge risks are weighted
equally).
*/
float NodeProcessingUnit::GetPathVolatility(vector<int> path) {
  float pathRisk;
  int sz = path.size() -1;

  // zero case
  if (sz <= 0) {
    return 0.0;
  }

  string edge;
  for (int i = 0; i < sz; i++) {
    edge = NodePairToString(path[i], path[i+1]);
    pathRisk += strategy->edgeRisks[edge];
  }
  return pathRisk / sz;
}

float NodeProcessingUnit::GetTaxationOfPath(int pathLength) {
  float tax = 1.0;
  for (int i = 0; i < pathLength; i++) {
    tax = tax * COMMERCE_TAX;
  }
  return tax;
}

/*
// description
estimates the worth of a contract based on current edge risks and path length
*/
/// TODO : estimate using log-based function instead, in which the best path
///        has a higher weight to it!
/// TODO : this is a temporary function
float NodeProcessingUnit::CalculateContractWorthMeasure(int otherNode) {
  UpdatePathScores(otherNode);

  // get the average path length and average volatility
  float volatility;
  float untaxed;

  auto bp = bestPaths[otherNode];

  if (bp.size() == 0) {
    return 0.0;
  }

  for (auto c = bp.begin(); c != bp.end(); c++) {
    volatility += GetPathVolatility(c->first);
    untaxed += (1 - GetTaxationOfPath(c->first.size() -1));
  }
  volatility = volatility / bp.size();
  untaxed = untaxed / bp.size();
  return ((1 - volatility) + untaxed) / 2;
}

/*
// description
calculates the ratio of number of valid best paths to maximum possible number of
best paths.
*/
float NodeProcessingUnit::GetRoutingScore() {

  float numerator = accumulate(bestPaths.begin(), bestPaths.end(),
                0.0,
                [](float t, pair<int, vector<pair<vector<int>, float>>> p) {return t + p.second.size();}
              );

  float denum = float(bestPaths.size() * NUM_BEST_PATHS);
  return ZeroDiv(numerator, denum, 0.0);
}

vector<int> NodeProcessingUnit::ShortenPath(vector<int> path, pair<int,int> newBond) {

  // error-check: no self-edges
  assert(newBond.first != newBond.second);

  // find the indices of the elements
  int index1 = -1, index2 = -1;
  for (int i = 0; i < path.size(); i++) {
    if (path[i] == newBond.first && index1 == -1) {
      index1 = i;
    } else if (path[i] == newBond.second && index2 == -1) {
      index2 = i;
    }
  }

  if (index1 == -1 || index2 == -1) {
    return path;
  }

  int indexMin = (index1 < index2) ? index1 : index2;
  int indexMax = (index1 > index2) ? index1 : index2;

  vector<int> output;

  // insert the first half
  copy(path.begin(), path.begin() + indexMin, back_inserter(output));
  output.push_back(path[indexMin]);
  copy(path.begin() + indexMax, path.end(), back_inserter(output));

  return output;
}

/*
// description
Calculates a vector of data relating to best paths resulting from a new bond,
from this Node instance to another Node. Only the top path to each Node is considered.

// return
a vector of elements, each element is a pair,
  [0] is the destination node,
  [1] is a pair, first is path score before, second is path score after bond.
*/
vector<pair<int, pair<float,float>>> NodeProcessingUnit::BestPathsAffectedByNewBondScores(pair<int,int> newBond) {

  vector<int> p;
  vector<int> p2;

  vector<pair<int, pair<float,float>>> output;
  pair<int, pair<float,float>> nodeAffectedInfo;
  pair<float,float> scorePair;

  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {
    if ((it->second).size() == 0) {
      continue;
    }

    p = it->second[0].first;
    p2 = ShortenPath(p, newBond);

    if (p.size() == p2.size()) {
      continue;
    }

    auto x = GetPathScore(p);
    auto newScore = GetPathScore(p2);
    scorePair = make_pair(x, newScore);
    nodeAffectedInfo = make_pair(it->first, scorePair);
    output.push_back(nodeAffectedInfo);
  }
  return output;
}

/*
// description
helper method for classmethod<BestPathsAffectedByNewBondSkippedNodes>
*/
vector<int> NodeProcessingUnit::GetSkippedNodes(vector<int> path, pair<int,int> newBond) {
  vector<int> output;
  int ind1 = -1, ind2 = -1;
  for (int i = 0; i < path.size(); i++) {
    if (path[i] == newBond.first && ind1 == -1) {
      ind1 = i;
    } else if (path[i] == newBond.second && ind2 == -1) {
      ind2 = i;
    }

    if (ind1 != -1 && ind2 != -1) {
      break;
    }
  }

  if (ind1 == -1 || ind2 == -1) {
    return output;
  }

  int indexMin = (ind1 < ind2) ? ind1 : ind2;
  int indexMax = (ind1 > ind2) ? ind1 : ind2;
  copy(path.begin() + indexMin + 1, path.begin() + indexMax, back_inserter(output));
  return output;
}

/*
// return
vector of elements, each element is a pair, with
[0] the destination node
[1] the vector of vertices skipped because of the new connection.
*/
vector<pair<int, vector<int>>> NodeProcessingUnit::BestPathsAffectedByNewBondSkippedNodes(pair<int,int> newBond) {

  vector<pair<int, vector<int>>> output;
  vector<int> p;
  pair<int, vector<int>> element;

  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {
    if ((it->second).size() == 0) {
      continue;
    }

    p = it->second[0].first;
    auto out = GetSkippedNodes(p, newBond);

    if (out.size() == 0) {
      continue;
    }

    element = make_pair(it->first, out);
    output.push_back(element);
  };

  return output;
}

/*
// return
the index of the first occurrence of bond in path 
*/ 
 pair<int,int> NodeProcessingUnit::FindBondIndexInPath(vector<int> path, pair<int,int> bond) {
     int s = path.size() - 1; 
  
     if (s <= 0) {
        return make_pair(-1,-1); 
     }

     for (int i = 0; i < s; i ++) {

        if (path[i] == bond.first && path[i + 1] == bond.second) {
            return make_pair(i, i + 1); 
        } 
        else if (path[i] == bond.second && path[i + 1] == bond.first) {
            return make_pair(i, i + 1); 
        }
     } 

     return make_pair(-1,-1); 
} 


/*
// description
helper method for classmethod<BestPathsAffectedByBondBreakScores>
*/
bool NodeProcessingUnit::BondExistsInPath(vector<int> path, pair<int, int> bond) {
  pair<int,int> index = FindBondIndexInPath(path, bond);

  if (index.first == -1 || index.second == -1) {
    return false;
  };

  if (abs(index.first - index.second) == 1) {
    return true;
  };

  return false;
}

/*
// description
outputs path score information for a destination node in the event that a bond is broken.

// output
[0]: bond exists in path
[1]:
  [0]: path score before bond broke
  [1]: path score after bond broke
*/
pair<bool, pair<float, float>>  NodeProcessingUnit::GetReplacementPathInfo(int key, pair<int,int> oldBond, float worstPathScore) {
  auto paths = bestPaths[key];

  // does not exist
  if (paths.size() == 0) {
    return make_pair(false, pair<float,float>());
  }

  if (BondExistsInPath(paths[0].first, oldBond)) {
    int j = -1;
    for (int i = 1; i < paths.size(); i++) {
      if (!BondExistsInPath(paths[i].first, oldBond)) {
        j = i;
        break;
      }
    };

    if (j != -1) {
      return make_pair(true, make_pair(paths[0].second, paths[j].second));
    } else {
      return make_pair(true, make_pair(paths[0].second, worstPathScore));
    }
  } else {
    return make_pair(false, pair<float,float>());
  }


}

/*
// return
a vector of elements, each element is a pair, with
[0] the destination node
[1] [0] best path score
    [1] replacement path score (methodvar<worstPathScore> if no replacement)
*/
vector<pair<int, pair<float, float>>> NodeProcessingUnit::BestPathsAffectedByBondBreakScores(pair<int,int> oldBond, float worstPathScore) {

  vector<pair<int, pair<float, float>>> output;

  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {
    auto x = GetReplacementPathInfo(it->first, oldBond, worstPathScore);
    if (x.first == true) {
      output.push_back(x);
    }
  }

  return output;
}

/*
// description
helper method for classmethod<BestPathsAffectedByBondBreakSkippedNodes>
*/
pair<bool, pair<vector<int>, vector<int>>> NodeProcessingUnit::SplitPathByBond(vector<int> path, pair<int,int> bond) {

  pair<int,int> index = FindBondIndexInPath(path, bond);
  pair<vector<int>, vector<int>> x;

  if (index.first == -1 || index.second == -1) {
    return make_pair(false, x);
  }

  int minIndex = (index.first < index.second) ? index.first: index.second;
  int maxIndex = (index.first > index.second) ? index.first: index.second;

  vector<int> p1, p2;

  copy(path.begin(), path.begin() + minIndex + 1, back_inserter(p1));
  copy(path.begin() + maxIndex, path.end(), back_inserter(p2));
  x = make_pair(p1,p2);
  return make_pair(true, x);
}

/*
// description
Outputs a vector of elements, each element is a pair, with
[0] the destination node
[1] a pair of vectors regarding the best path,
  the first vector is the subpath before the bond break,
  the second vector is the subpath after the bond break.
*/
vector<pair<int, pair<vector<int>, vector<int>>>> NodeProcessingUnit::BestPathsAffectedByBondBreakSkippedNodes(pair<int,int> oldBond) {

  vector<pair<int, pair<vector<int>, vector<int>>>> output;
  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {

    if (it->second.size() == 0) {
      continue;
    }

    auto x = SplitPathByBond(it->second[0].first, oldBond);
    if (x.first == true) {
      auto q = make_pair(it->first, x.second);
      output.push_back(make_pair(it->first, x.second));
    }
  }

  return output;
}

vector<pair<float,float>> NodeProcessingUnit::BestPathRiskInfoForNewContract(int node) {

  auto bp = bestPaths[node];
  vector<pair<float,float>> output;
  for (auto it = bp.begin(); it != bp.end(); it++) {
    output.push_back(GetPathRiskPotentialScores(it->first));
  }

  return output;
}

vector<int> NodeProcessingUnit::BestPathLengths(int node) {
  vector<int> output;
  auto bp = bestPaths[node];
  for (auto it = bp.begin(); it != bp.end(); it++) {
    output.push_back(it->first.size());
  };

  return output;
}

map<int, float> NodeProcessingUnit::GetBestPathScores() {
  map<int, float> bp;
  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {
    bp[it->first] = it->second.size() > 0 ? (it->second)[0].second : 0;
  }
  return bp;
}

/*
// description 
remove all info pertaining to missing edge 
*/ 
int NodeProcessingUnit::RemoveRecordsOfEdge(string edge) {
    // remove all paths that have edge 
    pair<int,int> p = StringToNodePair(edge); 
    int recRemoved = 0; 
    for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {

        auto paths = (*it).second;
        bool cnt = false; 
        while (cnt) {
            cnt = false; 
            for (auto it2 = paths.begin(); it2 != paths.end(); it2++) {
                if (BondExistsInPath((*it2).first, p)) {
                    cnt = true;
                    paths.erase(it2);
                    recRemoved++;   
                    break; 
                }
            }
        }
        bestPaths[(*it).first] = paths; 
    }
    return recRemoved; 
}
