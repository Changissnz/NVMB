#ifndef COMMERCEFLARE_H
#define COMMERCEFLARE_H
#define TRAVEL_RATE -1

/// TODO: below not yet implemented
// set below to postive integer n to transport a flare n edges
#define COMM_FLARE_INCREMENT -1
#include <vector>
#include "proposition.hpp"
#include "utility"
#include "assert.h"
#include <iostream>

class CommFlare
{

public:
  std::vector<int> targetPath;
  std::pair<int, int> nullEdge; // default is -1,-1
  Proposition* proposition;
  int currentDestination;
  int vertexLocation;
  int pathIndex;
  bool targetReached;
  bool finished;
  bool valid;
  bool initialized;

  int nodeOwnerIdentifier;
  int planIdentifier;

  std::map<int,float> costs;
  std::map<int,bool> responses; // for non-proposer nodes 


  bool active;

  CommFlare(std::vector<int> tp, Proposition* p, int cd, int pi, int noi) {
    std::cout << "comm flare info" << std::endl; 
    DisplayIterable(tp); 
    assert (tp.size() >= 2);
    targetPath = tp;
    nullEdge = std::make_pair(-1,-1); 
    proposition = p;
    currentDestination = cd;
    vertexLocation = tp[0];
    pathIndex = 0;
    targetReached = false;
    finished = false;
    valid = true;
    initialized = false;

    planIdentifier = pi;
    nodeOwnerIdentifier = noi;
    costs = std::map<int,float>();
    active = true;
  };

  CommFlare(int pi) {
    finished = true;
    valid = false;
    vertexLocation = -1;
    currentDestination = -2;
    planIdentifier = pi;
    initialized = false;
    proposition = nullptr; 
    nullEdge = std::make_pair(-1,-1); 
    costs = std::map<int,float>();
    active = true;

  };

  Proposition* GetProposition() {
    return proposition;
  };

  int GetNextPosition() {
    if (pathIndex + 1 >= targetPath.size()) {
      return -1;
    }
    return targetPath[pathIndex +1];
  };

  int GetPosition() {
    return vertexLocation;
  };

  int GetPreviousPosition() {
    return targetPath[pathIndex -1];
  };

  void UpdatePosition() {
    if (pathIndex + 1 >= targetPath.size()) {
      return;
    }
    vertexLocation = targetPath[pathIndex + 1];
    pathIndex++;
  };

  std::string GetEventType() {
    return (proposition != nullptr) ? proposition->category : "invalid";
  };

  std::string GetProperEventType() {

      if (proposition == nullptr) {
          return "invalid"; 
      }

    if (GetEventType() == "make_contract" || GetEventType() == "break_contract") {
      return GetEventType();
    } 

    if (GetSender() == GetProposer()) {
      return proposition->category + "/direct";
    } else {
      return proposition->category + "/phantom";
    }
  };

  void MarkReached() {
    targetReached = true;
  };

  void MarkInvalid() {
    valid = false; 
  };

  bool IsValid() {
    return valid;
  };

  bool IsReached() {
    return targetReached;
  };

  int GetProposer() {
    return (proposition != nullptr) ? proposition->proposerNode : -1;
  };

  int GetSender() {
    return (proposition != nullptr) ? proposition->referenceNodes.first : -1;
  };

  int GetReceiver() {
    return (proposition != nullptr) ? proposition->referenceNodes.second : -1;
  };

  std::string GetRoleOfNode(int nodeId) {

    if (nodeId != GetReceiver() && nodeId != GetProposer() && nodeId != GetSender()) {
      return "invalid";
    }

    std::string role;
    // direct
    if (GetSender() == GetProposer()) {
      role = (nodeId == GetSender()) ? "sender" : "receiver";
    } else { // phantom
      if (nodeId == GetSender()) {
        role = "sender";
      } else if (nodeId == GetReceiver()) {
        role = "receiver";
      } else {
        role = "proposer";
      }
    }

    return role;
  };

  void DisplayInfo() {
    std::cout << "* comm flare info for proposer " << GetProposer() << std::endl;
    std::cout << "event type: " << GetEventType() << std::endl; 
    std::cout << "sender " << GetSender() << " receiver " << GetReceiver() << std::endl;
    std::cout << "initialized: " << initialized << std::endl;
    std::cout << "target reached: " << targetReached << std::endl;
    std::cout << "finished: " << finished << std::endl;
    std::cout << "valid: " << valid << std::endl;
    std::cout << "location: " << vertexLocation << std::endl; 
    std::cout << "[] COSTS" << std::endl;
    for (auto x: costs) {
      std::cout << "\tnode " << x.first << " | " << x.second << std::endl;
    }
    std::cout << "--------------------------------------" << std::endl;
  };

};


#endif
