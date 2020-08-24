/*
A timestamp unit is used to record transaction costs between some reference
node and a map of target nodes.

info:
key : node id
value : <event_type, source, worth>

event_type : one of discovery info (source), proposition
*/

#ifndef TIMESTAMPUNIT_HPP
#define TIMESTAMPUNIT_HPP

#include <string>
#include <numeric>
#include <utility>
#include <assert.h>
#include <iostream>
#include "contract.hpp"
using namespace std;

/*
A TimestampUnit is used to record one action pertaining to a Node at a timestamp.
*/
class TimestampUnit {

public:

  int time;

  // NOTE: role `proposer` used only if proposer and sender are different
  string role; // proposer | sender | receiver

  ///types of events
  /*
  - cf tax
  - discovery
  - make_bond/phantom
  - make_bond/direct
  - break_contract
  - make_contract
  - contract activity
  - contract activity failure 
  - contract tax
  */
  string eventType;
  map<int,bool> response;
  float impact; // TODO : this value will have to be calculated by Node at a later timestamp
  vector<int> nodeIdentifiers; 
  map<string, int> otherDetails; // use -1 and 1 for some keys
  map<string,string> otherDetails2;

  /*
  standard constructor.
  CAUTION full-check assertions not enabled.
  */
  TimestampUnit(int t, string r, string et, map<int,bool> resp,
      float imp, int n1, int n2, int n3) {
    assert (t >= 0);
    time = t;

    assert (r == "proposer" || r == "sender" || r == "receiver" || r == "invalid");
    role = r;

    assert (et == "make_bond/phantom" ||
      et == "make_bond/direct"|| et == "break_bond/phantom" ||
      et == "break_bond/direct" || et == "break_contract" ||
      et == "make_contract" || et == "invalid");

    eventType = et;
    response = resp;
    impact = imp;
    nodeIdentifiers = {n1,n2,n3};

  };

  /*
  constructor specialized for discovery and cf tax processes
  */
  TimestampUnit(string et, int t, float imp) {
    assert (et == "discovery");
    eventType = et;
    time = t;
    impact = imp;
  }

  /*
  for `contract activity`
  ----
  float `delta change` = 1 | -1
  float `average path distance change` = 1 | -1
  float `transmission rate change` = 1 | -1
  ____________
  for `contract dead`
  ----
  empty
  */
  TimestampUnit(string et, int t, float imp, map<string, string> od) {
    assert (et == "cf tax" || et == "contract activity" 
        || et == "contract activity failure" || et == "contract tax"); 
    eventType = et;
    time = t;
    impact = imp;

    if (et == "cf tax") {
        nodeIdentifiers = {stoi(od["nodes"])};  
    } else {
        otherDetails2 = od;
    }
  }

  TimestampUnit(string et, int t, map<string,string> details) {
    assert (et == "summary"); 
    eventType = et;
    time = t; 
    impact = 0; 
    otherDetails2 = details; 
  }

  void DisplayInfo() {
    cout << "TIMESTAMP INFO" << endl;
    cout << "timestamp: " << time << endl;
    cout << "role: " << role << endl;
    cout << "event type: " << eventType << endl;
    if (eventType == "summary") {
        cout << "currency: " << otherDetails2["currency"] << endl; 
    }

    cout << "responses: " << endl;
    for (auto x: response) {
      cout << "\tnode " << x.first << " | " << x.second << endl;
    }
    cout << "impact: " << impact << endl;
    cout << "node identifiers:\t"; 
    DisplayIterable(nodeIdentifiers); 
    cout << "------------------------------" << endl;
  }

  string ToString() {
    string output = ""; 
    output += "timestamp: " + to_string(time) + "\n"; 
    output += "role: " + role + "\n"; 
    output += "event type: " + eventType + "\n"; 
    output += "impact: " + to_string(impact) + "\n"; 

    output += "responses: \n"; 
    for (auto x: response) {
      output += "- node " + to_string(x.first) + " | " + to_string(x.second) + "\n"; 
    }

    output += "details\n";
    output += "---------";  
    if (eventType == "summary") {
      for (auto x: otherDetails2) {
        output += "* " + x.first + "\t* " + x.second + "\n"; 
      }

    } else {
      for (auto x: otherDetails) {
        output += "* " + x.first + "\t* " + to_string(x.second) + "\n"; 
      }
    }
    output += "============================="; 
    return output; 
  }; 

  void DetailsString(int detailsNo) {
    assert (detailsNo == 0 || detailsNo == 1); 
    
    if (detailsNo == 0) {
      DisplayMap(otherDetails); 
    } else {
      DisplayMap(otherDetails2); 
    }
  };

    /*
   // description 
    converts timestamp to csv string-data form. See node_file_reader::TIMESTAMP_HEADER 
    for column details. 
   */ 
  string ToDataString() {
        string s = ""; 

        s += to_string(time) + ","; 
        s += eventType + ",";
        s += to_string(impact) + ","; 
        s += role + ","; 
        s += NumIterToStringType(nodeIdentifiers) + ","; 
        for (auto x: response) {
            s += to_string(x.first) + " " + to_string(x.second); 
        }
        s += ",";

        s += otherDetails2["currency"] + ","; 
        s += otherDetails2["neighbors"] + ","; 
        s += otherDetails2["contracts"] + ","; 
        s += otherDetails2["transmission"] + ",";
        s += otherDetails2["cfAcceptanceRate"] + ","; 
        s += otherDetails2["competition"] + ","; 
        s += otherDetails2["greed"] + ","; 
        s += otherDetails2["negotiation"] + ","; 
        s += otherDetails2["growth"] + ",";
        s += otherDetails2["nodeRank"] + ",";  
        s += otherDetails2["knownNodes"] + ",";

        s += otherDetails2["competitors"] + ","; 
        s += otherDetails2["numPathsRatio"] + ","; 
        s += otherDetails2["pathsInfo"];

        return s;  
  }; 
  
  /*
retrieves a variable 
  */ 
  string GetVar(string key) {
    if (key == "time") {
      return to_string(time); 
    } 

    if (key == "role") {
      return role; 
    }

    if (key == "event type") {
      return eventType; 
    }

    if (key == "event type") {
      return eventType; 
    }

    if (key == "responses") {
      string output = ""; 
      
      for (auto x: response) {
        output += to_string(x.first) + " " + to_string(x.second) + " | "; 
      }; 
      return output;
    } 

    if (key == "summary") {
      string output = "";
      for (auto x: otherDetails2) {
        output += x.first + " " + x.second + " | "; 
      }
      return output; 
    }

    if (key == "summary extra") {
      string output = ""; 
      for (auto x: otherDetails) {
          output += x.first + " " + to_string(x.second) + " | "; 
      } 
      return output;
    }

    throw invalid_argument("invalid key input for timestamp value"); 
  }; 

  /*
  will collect values for unique keys of map 
  */ 
  map<string,string> ToStringMap(vector<string> keys) {

    map<string,string> output = map<string,string>(); 

    set<string> got = set<string>(); 
    set<string>::iterator itg; 
    
    for (auto x: keys) {
      itg = got.find(x);

      // duplicate key 
      if (itg != got.end()) {
        continue; 
      }

      got.insert(x); 
      output[x] = GetVar(x); 
    } 
    
    return output; 
  } 

};

#endif