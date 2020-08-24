#ifndef NEWCONTRACTFUNCTIONS_HPP
#define NEWCONTRACTFUNCTIONS_HPP

#include <vector>
#include <utility>
#include <math.h>
#include <numeric>
#include "extras.hpp"

/*
Calculates a score in [0,1] for a new contract. Low scores means
contract is more attractive.
*/
float NewContractBasedOnBestPathsRisk(std::vector<std::pair<float, float>> pathsToDestination);
float NewContractBasedOnPredictedWorth(std::vector<int> bestPathLengths, float commerceTax);

/*
// description
Score relies on currency growth record of Node.
If predicted worth of new contract greater than the average contract worth of
Node, contract has greater likelihood of getting chosen.
*/
float NewContractBasedOnCurrentGrowth();

#endif
