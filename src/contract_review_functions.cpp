#include "contract_review_functions.hpp"

/*
// description
compare this value with Node's competition
*/
float ContractPerformanceMeasure_Distance(Contract* C, int networkDiameter) {
  if (networkDiameter <= 0) {
    return 0.0;
  }
  return (networkDiameter - C->averagePathDistance) / float(networkDiameter);
}

/*
// description
compare this value with Node's greed
*/
float ContractPerformanceMeasure_Transmission(Contract* C) {
  return C->transmissionRate;
}

/*
// description
two factors: average path distance and success rate
compare this value with Node's `competition` and `greed`
*/
float ContractPerformanceMeasure_DistanceAndTransmission(Contract* C, int networkDiameter) {

  if (IsMinimallyActive(C) == 1.0) {
    return 1.0;
  }

  float distanceMeasure = ContractPerformanceMeasure_Distance(C, networkDiameter);
  return (distanceMeasure + C->transmissionRate) / 2.0;
}

/*
// description
*/
float ContractPerformanceMeasure_AverageDistance(Contract* C, int networkDiameter, float averageNodeContractPathDistance) {
  float pm = ContractPerformanceMeasure_Distance(C, networkDiameter);
  float pmAverage = (networkDiameter - averageNodeContractPathDistance) / float(networkDiameter);
  return ZeroDiv(pm, (pm + pmAverage), 1.0);
}

/*
// description
*/
float ContractPerformanceMeasure_AverageTransmission(Contract* C, float averageTransmissionMeasure) {
  return ZeroDiv(C->transmissionRate, (C->transmissionRate + averageTransmissionMeasure), 1.0);
}

/*
// description
*/
float IsMinimallyActive(Contract *C) {
  if (C->lifespan < MIN_LENGTH_ACTIVE) {
    return 1.0;
  }
  return 0.0;
}

/// TODO: ?!? check, new addition.
float ContractRouteFailureValue(float nodeWorth) {
    return nodeWorth * CONTRACT_ROUTE_FAILURE_PENALTY; 
} 