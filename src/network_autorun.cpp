#include "network.hpp"
using namespace std;

/// TODO: 
/*
below needs to update node strategy functions as well.
*/ 

/// TODO: refactor below 
void Network::RunRandomNodeNatureVars(int verbose, int numRounds, float switchFrequency) {
    
    PrerunProcess(); 
    if (numRounds > 0) {
        while (IsAlive() && numRounds) {
            if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            UpdateNodeNaturesAtRandom(switchFrequency, "each", "reset");
            numRounds--; 
        }
    } else {
        while (IsAlive()) {
             if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            UpdateNodeNaturesAtRandom(switchFrequency, "each", "reset"); 
        }        
    }
    
    WriteOutToFile(false); 
    ShutDown();  
};


void Network::RunRandomPatternStyleNodeNatureVars(int verbose, int numRounds, float switchFrequency, string randomNatureType, string updateType) {

    PrerunProcess(); 
    if (numRounds > 0) {
        while (IsAlive() && numRounds) {
            if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            UpdateNodeNaturesAtRandom(switchFrequency, randomNatureType, "delta");
            numRounds--; 
        }
    } else {
        while (IsAlive()) {
             if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            UpdateNodeNaturesAtRandom(switchFrequency, randomNatureType, "delta"); 
        }        
    }

    WriteOutToFile(false); 
    ShutDown();  
};