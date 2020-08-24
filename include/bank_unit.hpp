#ifndef BANKUNIT_HPP
#define BANKUNIT_HPP

#include <map>
#include <vector>
#include <utility>
#include <iterator>
#include <iostream>
#include "assert.h"
#include "timestamp_unit.hpp"
#define DEFAULT_WORTH 100.0

// node file reader will be in charge of exporting bank unit data to file
class BankUnit {

private:
    std::map<int, std::vector<TimestampUnit*>> timestampUnitHistory;
    std::map<int,float> netChangePerTimestamp;
    std::vector<TimestampUnit*> currentTimestampInfo;
    int timestamp;

public:

    float currency;
    float currencyPrior; 
    float delta;

    BankUnit() {
        currency = DEFAULT_WORTH;
        timestamp = 0;
        delta = 0.0;
    };

    BankUnit(float c) {
        currency = c;
        timestamp = 0;
        delta = 0.0;
    };

    BankUnit(int ts) {
        currency = DEFAULT_WORTH;
        timestamp = ts;
        delta = 0.0;
    };

    std::map<int, std::vector<TimestampUnit*>> GetTimestampUnitHistory() {
        return timestampUnitHistory;
    };

    void ClearTimestampUnitHistory() {
        timestampUnitHistory.clear(); 
    }

    std::vector<TimestampUnit*> GetCurrentTimestampInfo() {
        return currentTimestampInfo; 
    }
  
    std::vector<TimestampUnit*> GetTimestampUnitHistoryAtRange(int min, int max); 

    vector<float> GetDeltaInRange(int min, int max); 

    void DisplayTimestampUnitHistory();

    void SetTimestamp(int t) {
        timestamp = t;
    }; 

    void UpdateCurrency(float x) {
        currencyPrior = currency; 
        currency = currency + x;
    };

    float GetNetChangeAtTimestamp();
    float GetNetChangeAtTimestamp(int timestamp);

    /// TODO: warning, no checks coded 
    /*
    // description 
    adds timestamp to cache 
    */
    void LogTimestampUnit(TimestampUnit* tsu) {
        currentTimestampInfo.push_back(tsu);
    };

    void LogTimestampData();

    int GetNumberOfTimestamps();
};

#endif
