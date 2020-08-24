#include "new_contract_functions.hpp"
using namespace std;

/*
// description
calculates new-contract score by best path's risk criteria.

// arguments
pathsToDestination - vector containing pairs, each pair is data that corresponds
                     to a path of the form
                     [0] actual risk score
                     [1] potential risk score
*/
float NewContractBasedOnBestPathsRisk(vector<pair<float, float>> pathsToDestination) {

  //
  float numerator = accumulate(pathsToDestination.begin(), pathsToDestination.end(),
                0.0,
                [](float t, pair<float,float> p) {return t + p.first;}
              );
  float denumerator = accumulate(pathsToDestination.begin(), pathsToDestination.end(),
                0.0,
                [](float t, pair<float,float> p) {return t + p.second;}
              );

  return ZeroDiv(numerator, denumerator, 1.0);
}

/*
// description
calculates the worthiness of a contract by the lengths of its n best paths.
Each paths are scaled by a coefficient:
  for an index i in [0, n -1], its coefficient is (n - i).
*/
float NewContractBasedOnPredictedWorth(vector<int> bestPathLengths, float commerceTax) {
  int sz = bestPathLengths.size();
  float bestPossibleWorth = 1.0 - commerceTax;

  float numerator = 0.0, denumerator = 0.0;
  for (int i = 0; i < bestPathLengths.size(); i++) {
    auto coeff1 = pow(bestPossibleWorth, bestPathLengths[i]);
    numerator += ((sz - i) * coeff1);
    denumerator += ((sz - i) * bestPossibleWorth);
  };

  return ZeroDiv(numerator, denumerator, 1.0);
}

/*
// description
*/
float NewContractBasedOnGrowthRecord() {
  return 1.0;
}
