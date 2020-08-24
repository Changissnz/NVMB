#ifndef PLAN_HPP
#define PLAN_HPP

#include <string>
#include <vector>
#include <iostream>
#include <assert.h>
#include "proposition.hpp"

class Plan {

public:
  int identifier;
  std::string planClass; // direct || phantom
  std::string planType; // discovery || make_bond || break_bond || make_contract || break_contract
  int sourceNodeIdentifier;
  int targetNodeIdentifier;
  int extraIdentifier;
  float score;

  // look at these two variables
  bool execSuccess;
  bool finished;

  Plan(int idn, std::string pc, std::string pt, int sni, int tni, float planScore, int eid) {
    assert (pt != "discovery" && pt != "invalid");
    identifier = idn;
    planClass = pc;
    planType = pt;
    extraIdentifier = eid; 
    sourceNodeIdentifier = sni;
    targetNodeIdentifier = tni;
    score = planScore;
    execSuccess = true;
  };

  /*
  // description
  declarator used for discovery process
  */
  Plan(int idn, std::string pt, float planScore) {
    assert (pt == "discovery");
    identifier = idn;
    planType = pt;
    score = planScore;
    execSuccess = true;
  };

  /*
  // description 
  invalid plan 
  */ 
  Plan() {planType = "invalid";};

  void DisplayInfo() {
    std::cout << "** plan" << std::endl;
    std::cout << "[x] identifier: " << identifier << std::endl;
    std::cout << "[x] plan class: " << planClass << std::endl;
    std::cout << "[x] plan type: " << planType << std::endl;
    std::cout << "[x] score: " << score << std::endl;
    std::cout << "[x] source: " << sourceNodeIdentifier << " [x] target: "
      << targetNodeIdentifier << std::endl;
    std::cout << "[x] exec success: " << execSuccess << std::endl;
    std::cout << "---------------------------" << std::endl;
  };

  void MarkSuccess(bool success);
};

#endif
