/*
this function constructs all the possible classes of plans a Node is permitted.
*/
#include "nvmb_node.hpp"

/*
// description
this method considers a Node's existing contracts. Contracts that have
not yielded adequate profit have a higher chance of getting broken.
*/
float NVMBNode::ConsiderExistingContract(Contract* C) {
   NodeStrategos* strategy = npu->GetStrategy();
   if (strategy->contractMeasure == "distance") {
     return ContractPerformanceMeasure_Distance(C, npu->maxPathLength);
   } else if (strategy->contractMeasure == "transmission") {
     return ContractPerformanceMeasure_Transmission(C);
   } else if (strategy->contractMeasure == "distance+transmission") {
     return ContractPerformanceMeasure_DistanceAndTransmission(C, npu->maxPathLength);
   } else if (strategy->contractMeasure == "average distance") {
     return ContractPerformanceMeasure_AverageDistance(C, npu->maxPathLength, NodeStrategos::averageNodeContractPathDistance);
   } else if (strategy->contractMeasure == "average transmission") {
     return ContractPerformanceMeasure_AverageTransmission(C, NodeStrategos::averageTransmissionMeasure);
   } else {
     cout << "THROW ERROR HERE" << endl;
   }
   return 0.0;
}

/*
// description
criteria for competitor is:
- growth rate over x past rounds
- Proposition rejected by other node
- number of active contracts other Node has
*/
bool NVMBNode::ConsiderCompetitor(int i) {
  // self cannot be competitor
  if (i == identifier) {
    return false;
  }

  NodeStrategos* strategy = npu->GetStrategy();
  pair<int,float> val = NodeStrategos::networkContractInfo[i];
  float output;

  if (strategy->competitionMeasure == "growth-performance") {
    output = CompetitionMeasureByPredictedGrowthPerformance(val.first, MAX_GROWTH_RATE, val.second);
  } else if (strategy->competitionMeasure == "growth-relative") {
    // calculate predictedGrowthSelf
    float predictedGrowthSelf = PredictedGrowthOfNode(strategy->activeContracts.size(),
          MAX_GROWTH_RATE, NodeStrategos::networkContractInfo[identifier].second);
    output = CompetitionMeasureByPredictedGrowthRelative(predictedGrowthSelf, val.first,
        MAX_GROWTH_RATE, val.second);
  } else if (strategy->competitionMeasure == "number of contracts") {
    output = CompetitionByNumberOfContracts(val.first, npu->GetKnownNumberOfNodes());
  } else if (strategy->competitionMeasure == "number of contracts relative") {
    output = CompetitionByNumberOfContractsRelative(strategy->activeContracts.size(), val.first);
  } else if (strategy->competitionMeasure == "path score to max") {
    auto bps = npu->GetBestPathScores();
    output = CompetitionMeasureByPathRelativeToMax(bps, i);
  } else if (strategy->competitionMeasure == "path score to mean") {
    auto bps = npu->GetBestPathScores();
    output = CompetitionMeasureByPathRelativeToMean(bps, i);
  } else {
    throw invalid_argument("received negative value");
  }
  return InterpretCompetitionMeasure(output, strategy->competition);
}

// TODO : not all relevant methods have been coded. See file `new_contract_functions.hpp`
/*
// description
method
*/
float NVMBNode::ConsiderNewContract(int newNode) {

  /// TODO: 
  //cout << "CONSIDERING NEW CONTRACT" << endl; 
  NodeStrategos* strategy = npu->GetStrategy();
  if (strategy->newContractMeasure == "best paths risk") {
    auto pathsToDestination = npu->BestPathRiskInfoForNewContract(newNode);
    return NewContractBasedOnBestPathsRisk(pathsToDestination);
  } else if (strategy->newContractMeasure == "predicted worth") {
    auto bpl = npu->BestPathLengths(newNode);
    return NewContractBasedOnPredictedWorth(bpl, COMMERCE_TAX);
  } else {
    throw invalid_argument("received invalid measure for new contract");
  }
}


/// TODO: this method needs work 
/*
// description
selects x bonds that are not known to exist,
in which x = 2 * (n -1), for n = number of known nodes at the time.

the q possible new bonds of this Node to another Node will be considered.

and for phantom propositions, x - q arbitrary new bonds

// arguments
bondType -
*/
set<string> NVMBNode::SelectPossibleBondsToAdd(){

  // consider direct bonds first
  set<string> output;
  int maxNumConsideredBonds = (npu->GetKnownNumberOfNodes() - 1) * 2;

  if (npu->GetNumberOfBestPaths() == 0) {
    return set<string>();
  };

  output = SelectPossibleDirectBonds();
  int remainingConsiderations = maxNumConsideredBonds - output.size();
  string edge;

  while (remainingConsiderations > 0) {
    edge = SelectBondFromBestPath("new", 10); // arbitrary recursive limit 
    if (edge != "") {
        auto checkDuplicate = StringToNodePair(edge); 
        if (checkDuplicate.first != checkDuplicate.second) {
            output.insert(edge);
        }  
    }

    remainingConsiderations--;
  }

  return output;
}

/*
// description
Selects possible direct bonds to add. Used for direct bond creation proposition.
*/
set<string> NVMBNode::SelectPossibleDirectBonds() {
    auto kn = GetKnownNodes(); 
    kn = RemoveSet(kn, neighborIdentifiers);
    kn.erase(identifier);
    set<string> output;
    string bond;
    for (auto it = kn.begin(); it != kn.end(); it++) {
        bond = NodePairToString(identifier, *it);
        output.insert(bond);
    }
    return output;
}

/*
// description
selects a random pair of vertices from a random best path.
If mode is set to "old bond", will choose a random edge in a random best path.
If mode is set to "new bond", will choose two random non-adjacent vertices in a
  random best path.
*/
string NVMBNode::SelectBondFromBestPath(string bondType, int recursiveLimit) {

    if (recursiveLimit <= 0) {
        return ""; 
    }

  assert (bondType == "old" || bondType == "new");

  if (npu->GetNumberOfBestPaths() == 0) {
    return "";
  };

  // choose destination vertex for paths
  auto bp = npu->GetBestPaths();
  auto it = bp.begin();
  std::advance(it, rand() % bp.size());

  // has no   
  if (it->second.size() == 0) {
      recursiveLimit--; 
    return SelectBondFromBestPath(bondType, recursiveLimit);
  }

  // choose random bond in best path
  auto path = it->second[0].first;

  // TODO : indexing error here.
  if (bondType == "old") {
    if (path.size() < 2) {
        recursiveLimit--; 
      return SelectBondFromBestPath(bondType, recursiveLimit);
    }
    int index = rand() % (path.size() - 1);
    return NodePairToString(path[index], path[index + 1]);
  } else {
    if (path.size() < 3) {
        recursiveLimit--; 
      return SelectBondFromBestPath(bondType, recursiveLimit);
    }

    int index = rand() % (path.size() - 2);

    vector<int> vec2;
    vec2.assign(path.begin() + index + 2, path.end());
    int index2 = rand() % (vec2.size());

    return NodePairToString(path[index], vec2[index2]);
  }
}

/*
// description
performs measures on bonds.
*/
float NVMBNode::MeasureBondCreation(pair<int,int> newBond, set<int> recognizedCompetitors) {
  NodeStrategos* strategy = npu->GetStrategy();
  float bondMeasure;
  float c = (identifier == newBond.first || identifier == newBond.second) ? strategy->competition : strategy->negotiation;
  if (strategy->bondAdvantageMeasure == "path scores-promise") {
    auto affectedPaths = npu->BestPathsAffectedByNewBondScores(newBond);
    return BondAddAdvantageMeasure_PathScoresPromise(affectedPaths,
        recognizedCompetitors, c, strategy->greed);
  } else if (strategy->bondAdvantageMeasure == "skipped nodes-promise") {
    auto affectedPaths2 = npu->BestPathsAffectedByNewBondSkippedNodes(newBond);
    return BondAddAdvantageMeasure_PathSkippedNodesPromise(affectedPaths2,
          recognizedCompetitors, c, strategy->greed);
  } else {
    throw invalid_argument("invalidavad argumentasaad");
  }
}

/////////////////////////////////////////////// bond deletion methods

/*
// description
Selects instance's known bonds to delete. 
The maximum number of considered bonds is : (KNOWN_NODES| - 1) * 2 
Recursive depth for function<SelectBondFromBestPath> is set to 10.

Procedure is outlined as such: 
  - Considers bonds with neighbors.  
  - Considers bonds between two adjacent nodes in a path. 
*/
set<string> NVMBNode::SelectPossibleBondsToDelete() {

  int maxNumConsideredBonds = (npu->GetKnownNumberOfNodes() - 1) * 2;
  set<string> output;

  // add neighbor bonds first
  for (auto it = neighborIdentifiers.begin(); it != neighborIdentifiers.end(); it++) {
    output.insert(NodePairToString(identifier, *it));
  }

  maxNumConsideredBonds -= output.size();
  string edge;

  while (maxNumConsideredBonds > 0) {
    /// TODO: delete below.
    edge = SelectBondFromBestPath("old", 1000);
    //edge = SelectBondFromBestPath("old", "phantom", 10);
    if (edge != "") {
      output.insert(edge);
    }
    maxNumConsideredBonds--;
  }
  return output;
}

/*
*/
float NVMBNode::MeasureBondDeletion(pair<int,int> oldBond, set<int> recognizedCompetitors) {
  NodeStrategos* strategy = npu->GetStrategy();
  float c = (identifier == oldBond.first || identifier == oldBond.second) ? strategy->competition : strategy->negotiation;

  if (strategy->bondDeletionAdvantageMeasure == "skipped nodes") {
    auto brokenPathInfo = npu->BestPathsAffectedByBondBreakSkippedNodes(oldBond);
    return BondBreakAdvantageMeasure_PathSkippedNodes(brokenPathInfo, recognizedCompetitors,
                      c, strategy->growth);

  } else if (strategy->bondDeletionAdvantageMeasure == "path scores") {
    auto brokenPathInfo = npu->BestPathsAffectedByBondBreakScores(oldBond, npu->worstPathScore);
    return BondBreakAdvantageMeasure_PathScores(brokenPathInfo, npu->worstPathScore, c, strategy->growth);
  } else {
    throw invalid_argument("invalidavad argumentasaad");
  }
}

/*
// description 
Shuts down all active objects and processes related to Node instance. 
Closes NodeFileReader fileObj. 
*/ 
void NVMBNode::ShutDown() {
    // write all remaining timestamps to file 
    WriteOutToFile(false); 
    reader->CloseWrite(); 
}