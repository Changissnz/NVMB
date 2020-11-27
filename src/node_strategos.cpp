/*
node strategos will be responsible for node's strategy
*/
#include "node_strategos.hpp"
using namespace std;

int FUTURE_PLAN_SIZE = 1;
float ACTIVATION_THRESHOLD = 0.5;

/// class static variables
float NodeStrategos::averageNodeContractPathDistance = 0;
float NodeStrategos::averageTransmissionMeasure = 0;

map<string, float> NodeStrategos::edgeRisks;
map<int, pair<int,float>> NodeStrategos::networkContractInfo;

/*
// description
Makes a NodeStrategos instance using both user-defined variables and
default values in the case user does not define.
*/
NodeStrategos* NodeStrategos::MakeOne(map<string, string> userDefinedVariables) {
  // set default values here
  float comp = 0.5;
  float greed = 0.5;
  float negotiation = 1.0;
  float growth = 1.0;

  // default values a1re arbitrary
  string compMeasure = "path score to max"; //"growth-performance"; //
  string contMeasure = "distance";
  string ncontMeasure = "best paths risk";
  string bam = "path scores-promise";
  string bdam = "skipped nodes";

  if (userDefinedVariables["competition"].length() != 0) {
    comp = stof (userDefinedVariables["competition"]);
  }

  if (userDefinedVariables["greed"].length() != 0) {
    greed = stof (userDefinedVariables["greed"]);
  }

  if (userDefinedVariables["negotiation"].length() != 0) {
    negotiation = stof (userDefinedVariables["negotiation"]);
  }

  if (userDefinedVariables["growth"].length() != 0) {
    growth = stof (userDefinedVariables["growth"]);
  }

  if (userDefinedVariables["competition measure"].length() != 0) {
    compMeasure = userDefinedVariables["competition measure"];
  }

  if (userDefinedVariables["contract measure"].length() != 0) {
    contMeasure = userDefinedVariables["contract measure"];
  }

  if (userDefinedVariables["new contract measure"].length() != 0) {
    ncontMeasure = userDefinedVariables["new contract measure"];
  }

  if (userDefinedVariables["bond advantage measure"].length() != 0) {
    bam = userDefinedVariables["bond advantage measure"];
  }

  if (userDefinedVariables["bond deletion advantage measure"].length() != 0) {
    bdam = userDefinedVariables["bond deletion advantage measure"];
  }

  NodeStrategos* s = new NodeStrategos(comp, greed, negotiation, growth, compMeasure,
    contMeasure, ncontMeasure, bam, bdam);
  return s;
};

void NodeStrategos::ChangeNodeNS(map<string, string> delta, string changeType) {
  //// TODO: chunk assumes not all attributes have keys
  if (changeType != "reset" && changeType != "delta") {
    throw invalid_argument("invalid change type"); 
  }

  float x; 
  for (auto it = delta.begin(); it != delta.end(); it++) {
    
    if (it->first == "competition") {
      x = stof(it->second); 
      if (x < 0 || x > 1) {
        throw invalid_argument("invalid value for competition");
      }

      if (changeType == "reset") {
        competition = x;
        continue; 
      }

      competition += x;

      if (competition > 1) {
        competition = competition - 1; 
      }
    } else if (it->first == "greed") {

      x = stof(it->second); 
      if (x < 0 || x > 1) {
        throw invalid_argument("invalid value for greed");
      }

      if (changeType == "reset") {
        greed = x;
        continue; 
      }

      greed += x;
      if (greed > 1) {
        greed = greed - 1; 
      }

    } else if (it->first == "negotiation") {

      x = stof(it->second); 
      if (x < 0 || x > 1) {
        throw invalid_argument("invalid value for negotiation");
      }

      if (changeType == "reset") {
        negotiation = x;
        continue; 
      }

      negotiation += x;

      if (negotiation > 1) {
        negotiation = negotiation - 1; 
      }
    } else if (it->first == "growth") {

      x = stof(it->second); 
      if (x < 0 || x > 1) {
        throw invalid_argument("invalid value for growth");
      }

      if (changeType == "reset") {
        growth = x;
        continue; 
      }

      growth += x;
      if (growth > 1) {
        growth = growth - 1; 
      } //////// strategy functions here, caution (no arg-check)
    } else if (it->first == "competition-measure") {
      competitionMeasure = it->second; 
    } else if (it->first == "contract-measure") {
      contractMeasure = it->second; 
    } else if (it->first == "newContract-measure") {
      newContractMeasure = it->second; 
    } else if (it->first == "bondAdvantage-measure") {
      bondAdvantageMeasure = it->second; 
    } else if (it->first == "bondDeletionAdvantage-measure") {
      bondDeletionAdvantageMeasure = it->second; 
    } else {
      throw invalid_argument("invalid key!");
    }
  }

  CalculateNormedNatureVar(); 
}; 

/*
// description
adds a plan to the front of `nodePlansCache`
*/
void NodeStrategos::AddPlanToPriority(Plan* p) {
  nodePlansCache.insert(nodePlansCache.begin(), p);
}


/*
// description
Adds a new plan to cache. Plan has not been executed.
*/
void NodeStrategos::AddPlans() {
  copy(futureNodePlans.begin(),futureNodePlans.end(),back_inserter(nodePlansCache));
}

/*
// description
Adds a plan to execution cache. Plan undergoing execution.
*/
void NodeStrategos::AddExecutedPlan(Plan* p) {
  plansInExecution.push_back(p);

  if (p->planType == "discovery") {
    SwitchDiscoveryMode();
  }
};

/*
// description
*/
int NodeStrategos::FindPlanByIdentifier(string planStorage, int planId) {
  assert (planStorage == "cache" || planStorage == "execution" || planStorage == "finished");
  vector<Plan*> targetStorage;

  if (planStorage == "cache") {
    targetStorage = nodePlansCache;
  } else if (planStorage == "execution") {
    targetStorage = plansInExecution;
  } else {
    targetStorage = finishedPlans;
  }

  for (int i = 0; i < targetStorage.size(); i++) {
    if (targetStorage[i]->identifier == planId) {
      return i;
    }
  }
  return -1;
}

/*
// description
deletes a plan of a specific type by its identifier
*/
void NodeStrategos::DeletePlan(string planType, int identifier) {
  int j = FindPlanByIdentifier(planType, identifier);
  if (planType == "cache") {
    nodePlansCache.erase(nodePlansCache.begin() + j);
  } else if (planType == "execution") {
    plansInExecution.erase(plansInExecution.begin() + j);
  } else {
    finishedPlans.erase(finishedPlans.begin() + j);
  }
}


/*
// description
marks a plan in execution to finished storage.

// arguments
indexType - direct (direct indexing) | plan (identifier)
i - index
*/
void NodeStrategos::FinishExecutedPlan(string indexType, int i) {
  assert (indexType == "direct" || indexType == "plan");
  Plan* p;
  if (indexType != "direct") {
    int j = FindPlanByIdentifier("execution", i);
    i = j;
  }

  if (i == -1) {
    return;
  }

  p = plansInExecution[i];
  plansInExecution.erase(plansInExecution.begin()+i);
  finishedPlans.push_back(p);
}

Plan* NodeStrategos::ChoosePlan() {
  if (nodePlansCache.size() == 0) {
    return new Plan();
  }
  Plan* p = nodePlansCache[0];
  nodePlansCache.erase(nodePlansCache.begin());
  return p;
};

/*

*/ 
void NodeStrategos::InsertDefaultPlan(Plan* p) {
  vector<Plan*>::iterator it = nodePlansCache.begin(); 
  nodePlansCache.insert(it, p); 
}


void NodeStrategos::SwitchDiscoveryMode() {
  if (discoveryOn == true) {
    discoveryOn = false;
  } else {
    discoveryOn = true;
  }
};

/*
// description
*/
Plan* NodeStrategos::MakeBondPlan(string bondType, int referenceNode, int sourceNode, int targetNode, float planScore) {
  // determine if direct or phantom
  assert (bondType == "make_bond" || bondType == "break_bond");
  string planClass;

  // plan type
  if (referenceNode == sourceNode || referenceNode == targetNode) {
    planClass = "direct";
  } else {
    planClass = "phantom";
  }

  Plan* p = new Plan(planIndex, planClass, bondType, sourceNode, targetNode, planScore, -1);
  planIndex++;
  return p;
}

/*
// description
constructs a make-contract plan 
*/
Plan* NodeStrategos::MakeContractPlan(int nodeId, int referenceNode, int otherNode, float planScore) {
  //string planClass = "direct";
  string planClass = (referenceNode == nodeId) ? "direct":"phantom"; 
  Plan* p = new Plan(planIndex, planClass, "make_contract", referenceNode, otherNode, planScore, -1);
  planIndex++;
  return p;
}

/*
// description 
constructs a break-contract plan 
*/ 
Plan* NodeStrategos::BreakContractPlan(int nodeId, int referenceNode, int otherNode, float planScore, int breakPlanIdentifier) {
  string planClass = (referenceNode == nodeId) ? "direct":"phantom"; 
  Plan* p = new Plan(planIndex, planClass, "break_contract", referenceNode, otherNode, planScore, breakPlanIdentifier);
  planIndex++; 
  return p; 
}

/*
// description
*/
Plan* NodeStrategos::MakeDiscoveryPlan(float planScore) {
  Plan* p = new Plan(planIndex, "discovery", planScore);
  planIndex++;
  return p;
}

/*
// description
inserts a future plan based on its score performance
*/
void NodeStrategos::InsertFuturePlan(Plan* p, int verbose) {

  // plan does not register
  if (p->score >= ACTIVATION_THRESHOLD) {
    if (verbose == 2) {cout << "discovery not permitted" <<  endl;};
    return ;
  }

  // find position to insert
  int j = futureNodePlans.size();
  for (int i = 0; i < futureNodePlans.size(); i++) {
    if (futureNodePlans[i]->score > p->score) {
      j = i;
      break;
    }
  }

  if (verbose == 2) {
    p->DisplayInfo();
  }

  futureNodePlans.insert(futureNodePlans.begin() + j, p);

  // reduce variable size to global limit
  while (futureNodePlans.size() > FUTURE_PLAN_SIZE) {
    futureNodePlans.pop_back();
  }
}

void NodeStrategos::DeleteActiveContract(int nodeId) {
  int index = -1;

  for (int i = 0; i < activeContracts.size(); i++) {
    if (activeContracts[i]->receiver == nodeId) {
      index = i;
      break;
    }
  }

  if (index != -1) {
    activeContracts.erase(activeContracts.begin() + index);
  } else {
    cout << "active contract for node " << nodeId << " not found" << endl;
  }
}

void NodeStrategos::DisplayVariables() {

  cout << "** NODE STRATEGY VARIABLES" << endl;
  cout << "competition:\t" << competition << endl;
  cout << "greed:\t" << greed << endl;
  cout << "normed competition:\t" << normedC << endl;
  cout << "normed greed:\t" << normedG << endl;
  cout << "negotiation:\t" << negotiation << endl;
  cout << "growth:\t" << growth << endl;
  cout << "---------------------------------------" << endl;
}

void NodeStrategos::DisplayAllPlans() {

  cout << "^ PLANS IN CACHE" << endl;
  for (auto x: nodePlansCache) {
    x->DisplayInfo();
  }

  cout << "^ PLANS IN EXECUTION" << endl;
  for (auto x: plansInExecution) {
    x->DisplayInfo();
  }

  cout << "^ PLANS FINISHED" << endl;
  for (auto x: finishedPlans) {
    x->DisplayInfo();
  }
}
