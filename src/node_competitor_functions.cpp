#include "node_competitor_functions.hpp"
using namespace std;

//// TODO: competitor functions need to be tested!
/*
// description
~
*/
float PredictedGrowthOfNode(int numberOfContracts, float growthRate, float averageContractRisk) {

  float exp = numberOfContracts - (averageContractRisk * numberOfContracts);
  return pow(growthRate, exp);
}

/*
// description
outputs value [0,1]

// arguments
numberOfContracts - regarding other node
growthRate - maximum growth rate set
averageContractRisk - regarding other node
*/
float CompetitionMeasureByPredictedGrowthPerformance(int numberOfContracts, float growthRate, float averageContractRisk) {

  float x = PredictedGrowthOfNode(numberOfContracts, growthRate, averageContractRisk);
  float y = pow(float(growthRate), float(numberOfContracts));
  return x / y;
}

/*
// description
generally, greater the difference between the nodes, greater the competition score.
if scores are on opposing sides of 0, then will always output 1.

// arguments
predictedGrowthSelf := predicted growth rate of the reference node
numberOfContracts := refers to target node
growthRate := refers to target node
averageContractRisk := refers to target node
*/
float CompetitionMeasureByPredictedGrowthRelative(float predictedGrowthSelf, int numberOfContracts,
    float growthRate, float averageContractRisk) {
  float x = PredictedGrowthOfNode(numberOfContracts, growthRate, averageContractRisk);
  float diff = abs(x - predictedGrowthSelf);
  float den = (abs(x) + abs(predictedGrowthSelf));

  return ZeroDiv(diff, den, 0.0);
}

/*
// description
~
*/
float CompetitionByNumberOfContracts(int numberOfContracts, int knownNumberOfNodes) {
  return float(numberOfContracts) / float(knownNumberOfNodes - 1);
}

/*
// description
~
*/
float CompetitionByNumberOfContractsRelative(int numberOfContractsSelf, int numberOfContractsOther) {
  float diff = abs(float(numberOfContractsSelf) - float(numberOfContractsOther));
  float den = ZeroDiv(float(numberOfContractsSelf), float(numberOfContractsOther), 1.0);
  return ZeroDiv(diff, den, 1.0);
}

/*
// description
range of output in [0, float x];
*/
float CompetitionMeasureByPathRelativeToMax(map<int,float> bestPathScores, int targetNode) {
  auto it = max_element(bestPathScores.begin(), bestPathScores.end(),
              [](pair<int,float> x, pair<int,float> x2) {return x.second < x2.second;}
            );

  return ZeroDiv(it->second - bestPathScores[targetNode], it->second, 1.0);
}

float CompetitionMeasureByPathRelativeToMean(map<int,float> bestPathScores, int targetNode) {

  float mean = accumulate(bestPathScores.begin(), bestPathScores.end(), 0.0,
              [](float x, pair<int,float> x2){return x + x2.second;}
            );

  mean = ZeroDiv(mean, float(bestPathScores.size()), 1.0);
  return bestPathScores[targetNode] / mean;
}

/*
// description
outputs true if competitor else false
*/
bool InterpretCompetitionMeasure(float measure, float threshold) {
  if (measure <= threshold) {
    return true;
  }
  return false;
}
