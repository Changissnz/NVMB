#include "contract.hpp"
using namespace std;
/*

// description
calculates the updated value of two nodes' bank values if they are contracted to
each other by the following:
----------------------------------------------------------------------------
d = (distance1 + distance2) / 2
e = (1 + transportTax) ** d

updatedBankValue1 = (CONTRACT_MULTIPLIER * bankValue1 + bankValue2) / e
updatedBankValue2 = (bankValue1 + CONTRACT_MULTIPLIER * bankValue2) / e
----------------------------------------------------------------------------
*/
pair<float,float> ContractUpdateFunction(float bankValue1, float bankValue2, int distance1, int distance2, float transportTax) {

  float averageDistance = float(distance1 + distance2) / 2.0;

  float denum = pow(1.0 + transportTax, averageDistance);

  float update1 = (CONTRACT_MULTIPLIER * bankValue1 + bankValue2) / denum;
  float update2 = (bankValue1 + CONTRACT_MULTIPLIER * bankValue2) / denum;
  return make_pair(update1, update2);
}

void Contract::UpdateContractData(bool alve) {
  alive = alve;
}

/// TODO: check this
void Contract::UpdateContractData(bool transmissionSuccess, float delta, int pathDistance) {
  lifespan++;
  float t = averagePathDistance * numTransmission;

  prevAverageDelta = averageDelta;
  prevTransmissionRate = transmissionRate;
  prevAveragePathDistance = averagePathDistance;

  if (transmissionSuccess) {
    numTransmission++;
    averageDelta = (prevAverageDelta * (numTransmission - 1) + delta) / float(numTransmission);
    averagePathDistance = (t + pathDistance) / float(lifespan);
  } else {
    averageDelta = (prevAverageDelta * (numTransmission - 1)) / float(numTransmission);
  }

  transmissionRate = float(numTransmission) / float(lifespan);
}

/*
// description 
output string will be the following: 
  idn averageDelta averagePathDistance numTransmission lifespan transmission_rate
*/
string Contract::ToString() {
  string output = ""; 

  output += to_string(idn); 
  output += " "; 
  output +=  to_string(averageDelta); 
  output += " "; 
  output +=  to_string(averagePathDistance);  
  output += " "; 
  output +=  to_string(numTransmission); 
  output += " ";
  output +=  to_string(lifespan); 
  output += " "; 
  output +=  to_string(transmissionRate); 
  
  return output; 
} 