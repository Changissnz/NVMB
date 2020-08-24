#ifndef CONTRACT_H
#define CONTRACT_H
/// TODO : implement or delete
#define MIN_LENGTH_ACTIVE 5
#define CONTRACT_MULTIPLIER 2.0

#include "extras.hpp"
#include <vector>
#include <iostream>
#include <math.h>

std::pair<float,float> ContractUpdateFunction(float bankValue1, float bankValue2, int distance1, int distance2);

class Contract {

public:
    int idn; // contract id
    int receiver;

    std::vector<int> path;  
    float pathRisk; 
    float averageDelta;
    float prevAverageDelta;

    float averagePathDistance;
    float prevAveragePathDistance;

    int numTransmission;
    int lifespan;
    float transmissionRate;
    float prevTransmissionRate;
    bool alive;

    Contract(int identifier, int rn, std::vector<int> path) {
        idn = identifier;
        receiver = rn;

        averageDelta = 0.0;
        prevAverageDelta = 0.0;

        averagePathDistance = 0.0;
        prevAveragePathDistance = 0.0;

        numTransmission = 0;
        lifespan = 0;

        transmissionRate = 0.0;
        prevTransmissionRate = 0.0;

        alive = true;
    };

    Contract() {
        idn = -1;
    };

    void UpdateContractData(bool alve);

    void UpdateContractData(bool transmissionSuccess, float delta, int pathDistance);

    void DisplayInfo() {
        std::cout << "contract " << idn << " receiver " << receiver << std::endl;
    }
 
    std::string ToString(); 

};

#endif
