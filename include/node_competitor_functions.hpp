#ifndef NODECOMPETITORFUNCTIONS_HPP
#define NODECOMPETITORFUNCTIONS_HPP

#include <math.h>
#include <cmath>
#include <climits>
#include <algorithm>
#include <map>
#include <numeric> 
#include "extras.hpp"

float PredictedGrowthOfNode(int numberOfContracts, float growthRate, float averageContractRisk);

float CompetitionMeasureByPredictedGrowthPerformance(int numberOfContracts, float growthRate, float averageContractRisk);

float CompetitionMeasureByPredictedGrowthRelative(float predictedGrowthSelf, int numberOfContracts,
    float growthRate, float averageContractRisk);

float CompetitionByNumberOfContracts(int numberOfContracts, int knownNumberOfNodes);

float CompetitionByNumberOfContractsRelative(int numberOfContractsSelf, int numberOfContractsOther);

float CompetitionMeasureByPathRelativeToMax(std::map<int,float> bestPathScores, int targetNode);

float CompetitionMeasureByPathRelativeToMean(std::map<int,float> bestPathScores, int targetNode);

// TODO : this has not yet been coded.
float CompetitionByRejectedPropositions(int numberOfRelatedPropositions, int numberOfTotalPropositions);
bool InterpretCompetitionMeasure(float measure, float threshold);

#endif
