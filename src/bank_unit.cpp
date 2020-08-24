#include "bank_unit.hpp"
using namespace std;

std::vector<TimestampUnit*> BankUnit::GetTimestampUnitHistoryAtRange(int min, int max) {

    assert (max > min);
    std::vector<TimestampUnit*> output = std::vector<TimestampUnit*>(); 
    std::vector<TimestampUnit*> tmp; 
    for (int i = min; i < max; i++) { 
        tmp = timestampUnitHistory[i]; 
        output.insert(output.end(), tmp.begin(), tmp.end()); 
    }
    return output; 
} 


/*
// description 
searches through timestamps and checks for delta in range  
*/ 
vector<float> BankUnit::GetDeltaInRange(int min, int max) {
    assert (max > min);
    assert (min >= 0);
    vector<float> output; 

    for (int i = min; i < max; i++) {
        auto x = timestampUnitHistory[i]; 
        for (auto y: x) {
            if (y->eventType == "summary") {
                output.push_back(stof(y->otherDetails2["currency"]));             
            }
        }
    }

    return output; 
} 

/*
// description 
displays all timestamp units 
*/ 
void BankUnit::DisplayTimestampUnitHistory() {

  map<int, vector<TimestampUnit*>>::iterator it2;
  map<int, vector<TimestampUnit*>> x;

  for (auto it = timestampUnitHistory.begin(); it != timestampUnitHistory.end(); it++) {
    vector<TimestampUnit*> x = it->second;
    for (auto y : x) {
      cout << y->ToString() << endl; 
    }
  } 
} 

/*
// description
iterates through timestamps and calculates the effect of node influence on
the current timestamp.
*/
float BankUnit::GetNetChangeAtTimestamp() {
  float sum = accumulate(currentTimestampInfo.begin(), currentTimestampInfo.end(),
                        0, [](const float& a, TimestampUnit* b)
                    {
                      return a + b->impact;
                    });
  return sum;
}

/*
// description
*/ 
float BankUnit::GetNetChangeAtTimestamp(int timestamp) {

    if (timestampUnitHistory.find(timestamp) == timestampUnitHistory.end()) {
        cout << "X timestamp " << timestamp << " does not exist" << endl; 
        return 0; 
    }

    auto vtsu = timestampUnitHistory[timestamp]; 
    return accumulate(vtsu.begin(), vtsu.end(),
                    0, [](const float& a, TimestampUnit* b)
                    {
                      return a + b->impact;
                    });
} 


// description
/*
update's bank's currency with net change from timestamp and logs the vector
of timestamps to classvar<timestampUnitHistory>
*/
void BankUnit::LogTimestampData() {
  timestampUnitHistory[timestamp] = currentTimestampInfo;
  currentTimestampInfo = vector<TimestampUnit*>();
}

int BankUnit::GetNumberOfTimestamps() {

  return accumulate(timestampUnitHistory.begin(), timestampUnitHistory.end(),
            0,
            [](int t, pair<int, vector<TimestampUnit*>> x ) {
              return t + x.second.size();
            });

}