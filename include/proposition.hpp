#ifndef PROPOSITION_HPP
#define PROPOSITION_HPP

#include <string>
#include <utility>
#include <functional>
#include "contract.hpp"
#include <iostream> 

class Proposition {

public:
    int propositionId;
    std::string category; // make_bond, break_bond, make_contract, break_contract

    int proposerNode;
    std::pair<int,int> referenceNodes; // first is sender, second is receiver

    Proposition(int pid, std::string cat, int pn, int senderNode, int recvNode) {
        propositionId = pid;
        assert (cat == "make_bond" || cat == "break_bond" || cat == "make_contract" || cat == "break_contract");
        category = cat;
        proposerNode = pn;
        referenceNodes = std::make_pair(senderNode, recvNode);
    };

    void DisplayInfo(); 

  Contract* ConvertToContract();

};

#endif
