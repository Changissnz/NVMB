#ifndef CONTRACTREVIEWFUNCTIONS_HPP
#define CONTRACTREVIEWFUNCTIONS_HPP
#define CONTRACT_ROUTE_FAILURE_PENALTY 0.2
#include "contract.hpp"

/*
compare this value with Node's competition
*/
float ContractPerformanceMeasure_Distance(Contract* C, int networkDiameter);

/*
compare this value with Node's greed
*/
float ContractPerformanceMeasure_Transmission(Contract* C);
float ContractPerformanceMeasure_DistanceAndTransmission(Contract* C, int networkDiameter);
float ContractPerformanceMeasure_AverageDistance(Contract* C, int networkDiameter, float averageNodeContractPathDistance);
float ContractPerformanceMeasure_AverageTransmission(Contract* C, float averageTransmissionMeasure);
float IsMinimallyActive(Contract *C);
float ContractRouteFailureValue(float nodeWorth); 
#endif
