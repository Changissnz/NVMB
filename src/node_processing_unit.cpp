#include "node_processing_unit.hpp"
#include "assert.h"
using namespace std;

const int DEFAULT_NUM_THREADS = 15;

/*
loads initial paths for neighbors 
*/ 
void NodeProcessingUnit::LoadInitialPaths() {
    vector<int> path;
    bestPaths = map<int, vector<pair<vector<int>, float>>>(); 
    for (auto n : neighborIdentifiers) {
        path = {GetNOI(), n}; 
        cout << "PATH FOR NODE " << GetNOI() << endl; 
        DisplayIterable(path); 
        bestPaths[n] = vector<pair<vector<int>, float>>(); 
        bestPaths[n].push_back(make_pair(path, 0.0));  
    }

    cout << "AFTER SET INITIAL" << endl; 
}

int NodeProcessingUnit::GetNOI() {
  return nodeOwnerIdentifier;
}

set<int> NodeProcessingUnit::GetNeighborIdentifiers() {
  return neighborIdentifiers;
}

void NodeProcessingUnit::UpdateNeighborIdentifers(set<int> x) {
  neighborIdentifiers = x;
}

/*
// description
displays the best path, if known, for each destination node.
*/
void NodeProcessingUnit::DisplayBestPaths() {

  for (auto x: bestPaths) {
    vector<pair<vector<int>,float>> info = x.second;
    cout << "* best path for: " << x.first << endl;
    if (info.size() > 0) {
      DisplayIterable(info[0].first);
      cout << "score: " << info[0].second << endl;
      cout << "-----------------------------------" << endl;
    }
    cout << "-----" << endl;
  }

}

/*
// description
activates a discovery flare at node's location only if
discovery process is currently not running.
*/
bool NodeProcessingUnit::ActivateDF() {

  discoveryOn = true;
  DiscoveryFlare* df = new DiscoveryFlare(nodeOwnerIdentifier);
  df->UpdateNeighborIdentifers(GetNeighborIdentifiers());

  ownedFlares.insert(df);
  AddHeldDF(df);
  return true;
}

/*
// description
adds an owned DiscoveryFlare to class var
*/
bool NodeProcessingUnit::AddOwnedDF(DiscoveryFlare* df) {
  ownedFlaresCache.insert(df);
  return true;
}

/*
// description
adds a DiscoveryFlare to owned set
*/
bool NodeProcessingUnit::AddHeldDF(DiscoveryFlare* df) {
  heldFlares.insert(df);
  return true;
}

/*
// description
removes a DiscoveryFlare from owned set
*/
void NodeProcessingUnit::RemoveHeldDF(DiscoveryFlare* df) {
  heldFlares.erase(df);
}

/*
// description
moves each owned discovery flare by one untravelled edge, and make flares for all
other untravelled edges.
*/
bool NodeProcessingUnit::Process1DF(DiscoveryFlare* df) {

  // iterates through neighbors for untravelled edges
  int numUntravelled = 0;
  DiscoveryFlare* df2;
  DiscoveryFlare* df3;

    // randomize below ordering 
  vector<int> ni = SetToShuffledVector(df->GetNeighborIdentifiers()); 
  vector<int>::iterator it = ni.begin();  
  int initialLocation = df->GetLocation();

  // get the first untravelled
  df2 = df->MakeCopy();
  while (it != ni.end()) {
    if (!df->HasTravelledEdge(df->GetLocation(), (*it))) {
      numUntravelled++;
      df->UpdateLocation((*it));
      df->MarkEdgeAsTravelled(NodePairToString(initialLocation, (*it)));
      df->AddToPath(*it);
      it++;
      break;
    }
    it++;
  };

  int q = MAX_FLARE_MULTIPLIER; 

  // make flares for all other untravelled edges
  while (it != ni.end()) {
    if (q == 0 && discoveryType == "limit") {
        break; 
    }

    if (ownedFlares.size() + ownedFlaresCache.size() >= MAX_FLARES && discoveryType == "limit") {
        break; 
    }
    
    if (!df->HasTravelledEdge(initialLocation, (*it))) {
        numUntravelled++;
        df3 = df2->MakeCopy();
        df3->UpdateLocation((*it));
        df3->MarkEdgeAsTravelled(NodePairToString(initialLocation, (*it)));

        df3->AddToPath(*it);
        AddOwnedDF(df3); 

        q--; 

    };
    it++;
  };

  if (numUntravelled == 0) {
    return true;
  }
  return false;
}

/*
*/
void NodeProcessingUnit::ProcessDF_() {

  set<DiscoveryFlare*>::iterator it = ownedFlares.begin();
  int sz = ownedFlares.size();

  /// PREVIOUS: before discoveryMode == "limit" 
    for (int i = 0; i < sz; i++) {
        // TODO : caution here
        if ((*it)->finished != true) {
        (*it)->finished = Process1DF((*it));
        }
        it++;
    };
}

/*
// description 
determines if discovery process has finished. 
*/ 
bool NodeProcessingUnit::DiscoveryFinished() {

    if (GetOwnedFlaresCache().size() > 0) {
        return false;
    }

    set<DiscoveryFlare*>::iterator it;
    for (it = ownedFlares.begin(); it != ownedFlares.end(); it++) {
        if ((*it)->finished == false) {
            return false;
        }
    };

    return true;
}

/*
*/
void NodeProcessingUnit::ProcessDF() {
  ProcessDF_();
}

/*
*/
set<DiscoveryFlare*> NodeProcessingUnit::GetOwnedFlares() {
  return ownedFlares;
}

/*
*/
set<DiscoveryFlare*> NodeProcessingUnit::GetOwnedFlaresCache() {
  return ownedFlaresCache;
}

/*
// description
*/
void NodeProcessingUnit::ClearDiscoveryData() {
    ownedFlares = set<DiscoveryFlare*>();
    ownedFlaresCache = set<DiscoveryFlare*>();
    UpdateHeldFlares("finished"); 
}

/*
*/
set<DiscoveryFlare*> NodeProcessingUnit::GetHeldFlares() {
  return heldFlares;
}

void NodeProcessingUnit::DisplayCF(string flareLabel) {

    assert (flareLabel == "owned" || flareLabel == "held" || flareLabel == "processed" || flareLabel == "finished");
    set<CommFlare*> x;
    if (flareLabel == "owned") {
        x = ownedCFlares;
    } else if (flareLabel == "held") {
        x = heldCFlares;
    } else if (flareLabel == "processed") {
        x = processedCFlares; 
    } else {
        x = finishedCFlares;
    }

    cout << "CF " << flareLabel << " for Node " << nodeOwnerIdentifier << endl;
    cout << "===================================" << endl;
    for (auto y: x) {
        y->DisplayInfo();
    }
}

/// TODO: refactor this.
/*
// description
updates variables, including:
- heldFlares: flare neighbor data will be updated
- ownedFlares: ~
*/
void NodeProcessingUnit::Update() {
  // iterate through flares and update their
  UpdateOwnedFlares();
  UpdateHeldFlares("finished");
}


int NodeProcessingUnit::UpdateHeldFlaresByIteration(string updateType, int index) {
  
    assert (updateType == "location" || updateType == "finished"); 
    set<DiscoveryFlare*>::iterator it = heldFlares.begin();
    advance(it, index);
    int index2 = index;
    while (it != heldFlares.end()) {
        if (updateType == "location") {
            if ((*it)->GetLocation() != GetNOI()) {
                RemoveHeldDF(*it);
                return index2;
            } else {
                (*it)->UpdateNeighborIdentifers(GetNeighborIdentifiers());
            }
        } else {
            if ((*it)->finished) {
                RemoveHeldDF(*it);
                return index2;
            }
        }

        index2++;
        it++;
    }
    return -1;
}

/*
// description
updates var<HeldFlare>.
- if the flare belongs to <this>, update neighbor data
- otherwise, remove it.

updateType - location|finished
*/
void NodeProcessingUnit::UpdateHeldFlares(string updateType) {

    int index = 0;
    while (index != -1) {
        index = UpdateHeldFlaresByIteration(updateType, index);
    }
}

/*
// description
updates owned flares with cache
*/
void NodeProcessingUnit::UpdateOwnedFlares() {
    ownedFlares = UnionSet(ownedFlares, ownedFlaresCache);
    ownedFlaresCache = set<DiscoveryFlare*>();
}

bool NodeProcessingUnit::IsFinished() {
  return finished;
}

bool NodeProcessingUnit::PathExists(vector<int> path, int key) {

  if (bestPaths.find(key) == bestPaths.end()) {
    return false;
  }

  auto x = bestPaths[key];
  for (auto y = x.begin(); y != x.end(); y++) {
    if (path == y->first) {
      return true;
    }
  }
  return false;
}

/// TODO: careful, call this method after discovery process is finished. 
void NodeProcessingUnit::UpdateMaxPathLength() {
  /*
  if (pathLength > maxPathLength) {
    maxPathLength = pathLength;
  }
  */ 
  maxPathLength = 0;

  for (auto sp: bestPaths) {
    for (auto sp2: sp.second) {
      if (maxPathLength < sp2.first.size()) {
        maxPathLength = sp2.first.size(); 
      }
    }
  }

  
}

/*
paths are sorted from least to greatest score
*/
bool NodeProcessingUnit::AddBetterPath(int key, vector<int> path, float score) {

    // case: # of best paths for target is 0 
    if (bestPaths[key].size() == 0) {
        bestPaths[key] = vector<pair<vector<int>, float>>();
        bestPaths[key].push_back(make_pair(path, score));
        /*
        cout << "updating max path" << endl; 
        DisplayIterable(path); 
        UpdateMaxPathLength(path.size());
        */ 
        return true;
    }

    vector<pair<vector<int>, float>>::iterator it;
    auto y = bestPaths[key];
    pair<vector<int>,float> p;

    // check duplicate path
    for (it = bestPaths[key].begin(); it != bestPaths[key].end(); it++) {
        if (path == it->first) {
            return false;
        }
    }

    // case: iterate through existing best paths to target and determine if 
    //       `path` is better.
    bool best = false;  
    for (it = bestPaths[key].begin(); it != bestPaths[key].end(); it++) {

        if (score <= it->second) {
            p = make_pair(path, score);
            bestPaths[key].insert(it, p);
            /*
            cout << "updating max path" << endl; 
            DisplayIterable(path); 
            UpdateMaxPathLength(path.size());
            */ 
            best = true; 
            break;
        }
    }

    if (bestPaths[key].size() > NUM_BEST_PATHS) {
        bestPaths[key].pop_back();
    }

    return best;  
}

bool NodeProcessingUnit::EliminatePath(int key, vector<int> path, float score) {

    if (bestPathsCache[key].size() < NUM_BEST_PATHS) {
        bestPathsCache[key].push_back(make_pair(path,score));
        return false; 
    }

    for (auto it = bestPaths[key].begin(); it != bestPaths[key].end(); it++) {

        if (path == it->first) {
            return false; 
        }

        if (score <= it->second) {
            bestPathsCache[key].push_back(make_pair(path,score));
            return false; 
        }
    }

    return true;
}

bool NodeProcessingUnit::IsCycle(vector<int> path) {
  set<int> s(path.begin(), path.end());
  if (s.size() != path.size()) {
    return true;
  }
  return false;
}

/*
return: 
-1 for cycle 
1 for not finished 
0 for finished 
*/ 
int NodeProcessingUnit::ConsiderPath(DiscoveryFlare* df, bool excludeCycle) {

  vector<int> v = df->GetPath();
  float score;
  vector<int> sp;
  bool stat; 

  // iterate through path and check for shortest lengths
  for (int i = 0; i < v.size(); i++) {
    vector<int> sp = vector<int>();
    sp.insert(sp.begin(), v.begin(), v.begin() + i + 1);

    if (excludeCycle) {
      if (IsCycle(sp)) {
        return -1;
      }
    }

    // get the score of the path
    score = GetPathScore(sp);
    AddBetterPath(v[i], sp, score);

  }

  // if the flare is marked finished, then done
  if (df->finished == true) {
    return 0;
  }
  return 1;
}

/*
// description
helper method for classmethod<FilterFlares>
*/
int NodeProcessingUnit::FilterOneFlare(bool excludeCycle, int index) {
    set<DiscoveryFlare*> dfs = GetOwnedFlares();
    int considered;
    int index2 = index;

    set<DiscoveryFlare*>::iterator it = dfs.begin();
    advance(it, index);

    while (it != dfs.end()) {
      considered = ConsiderPath(*it, excludeCycle);

      if (considered == 1) {
        AddOwnedDF(*it);
        ownedFlares.erase(*it);
        return index2;
      }
      // TODO ADDITIONS
      else if (considered == -1) {
        ownedFlares.erase(*it);
        return index2;
      }
      index2++;
      it++;
    }
    return -1;
}
 

/*
// description 
marks all owned flares that possess 
*/
void NodeProcessingUnit::EliminateDiscoveryFlares() {
    
    if (discoveryType != "elimination") {
        return; 
    }

    cout << "ELIMOS" << endl; 
    
    set<DiscoveryFlare*> dfs = GetOwnedFlares();

    int key; 
    vector<int> subpath; 
    float score; 
    bool elim; 

    bestPathsCache.clear(); 

    for (auto d :dfs) {
        vector<int> path = d->GetPath();
        ///cout << "determining elimination" << endl; 
        for (int i = 0; i < path.size(); i++) {
            subpath = vector<int>(); 
            subpath.insert(subpath.begin(), path.begin(), path.begin() + i + 1);
            score = GetPathScore(subpath); 
            elim = EliminatePath(path[i], subpath, score); 
        
            if (elim) {
                d->finished = true;
                break; 
            }
        }
    }

    EliminateFinishedFlares("owned"); 
    EliminateFinishedFlares("cached");

    UpdateBestPathsWithCache();   
}


void NodeProcessingUnit::UpdateBestPathsWithCache() {
 
    for (auto bpi: bestPathsCache) {   
        for (auto pathInfo: bpi.second) {
            AddBetterPath(bpi.first,pathInfo.first, pathInfo.second); 
        }
    }
}

void NodeProcessingUnit::EliminateHeldDF() {

    if (discoveryType != "elimination") {
        return; 
    }

    EliminateFinishedFlares("held"); 


}




/*
eliminates finished flares

dfType - owned | cached |
*/ 
void NodeProcessingUnit::EliminateFinishedFlares(string dfType) {

    assert (!(dfType != "owned" && dfType != "cached" && dfType != "held"));  

    set<DiscoveryFlare*> df;
    if (dfType == "owned") {
        df = GetOwnedFlares(); 
    } else if (dfType == "cached") {
        df = GetOwnedFlaresCache(); 
    } else {
        df = GetHeldFlares(); 
    }
    
    int c = 0; 

    set<DiscoveryFlare*> newDf = set<DiscoveryFlare*>(); 
    for (auto x: df) {
        if (!x->finished) {
            newDf.insert(x); 
        } else {
            c++; 
        }
    }

    if (dfType == "owned") {
        ownedFlares = newDf;  
        cout << "node " << GetNOI() << " eliminates " << c << " owned flares" << endl; 
    } else if (dfType == "cached") {
        ownedFlaresCache = newDf;
        cout << "node " << GetNOI() << " eliminates " << c << " owned flares cache" << endl; 
    } else {
        heldFlares = newDf; 
    }
}


/*
// description
filters flares according to path length and cycle existence,
if path is still valid, exports it to classvar<ownedFlaresCache>,
otherwise, remove it.
*/
void NodeProcessingUnit::FilterFlares(bool excludeCycle) {
  int index = 0;
  maxPathLength = 0;
  while (true) {
    index = FilterOneFlare(excludeCycle, index);
    if (index == -1) {
      break;
    }
  }

  UpdateOwnedFlares();
  UpdateWorstPathScore();
}

/*
// description
iterates through known best paths and outputs the worst (max) score
*/
void NodeProcessingUnit::UpdateWorstPathScore() {

  float score = INFINITY * -1;
  for (auto it = bestPaths.begin(); it != bestPaths.end(); it++) {
    auto paths = it->second;
    for (auto it2 = paths.begin(); it2 != paths.end(); it2++) {
      if (it2->second > score) {
        score = it2->second;
      }
    }
  }
  worstPathScore = score;
}

CommFlare* NodeProcessingUnit::ProduceCommFlare(Plan* p) {

  int destNode = (p->sourceNodeIdentifier == nodeOwnerIdentifier) ? p->targetNodeIdentifier: p->sourceNodeIdentifier;
  Proposition* prop = new Proposition(strategy->propIndex, p->planType, nodeOwnerIdentifier, p->sourceNodeIdentifier, p->targetNodeIdentifier);
  pair<vector<int>, bool> pathInfo = GetPathForDestination(destNode);

  // case: no path exists
  if (!pathInfo.second) {
    return new CommFlare(p->identifier);
  }

  vector<int> path = pathInfo.first;

  cout << nodeOwnerIdentifier << " prop for cf" << endl; 
  prop->DisplayInfo(); 
  CommFlare* cf = new CommFlare(path, prop, destNode, p->identifier, nodeOwnerIdentifier);
  
  ///cf->DisplayInfo(); 
  
  return cf;
};

/*
// description
*/
pair<vector<int>, bool> NodeProcessingUnit::GetPathForDestination(int destNode) {

    auto choices = bestPaths[destNode];

    if (choices.size() == 0) {
        return make_pair(vector<int>(), false);
    };

    vector<int> path = choices[0].first;
    return make_pair(path, true);
}

/*
// description
retrieves contract pertaining to nodeId 
*/
Contract* NodeProcessingUnit::GetExistingContract(int nodeId) {

  vector<Contract*> contracts = strategy->activeContracts;

  for (auto c: contracts) {
    if (c->receiver == nodeId) {
      return c;
    }
  }
  return new Contract();
}

/*
// description
*/
void NodeProcessingUnit::InitiateContract(CommFlare* cf) {

  if (nodeOwnerIdentifier != cf->GetSender() && nodeOwnerIdentifier != cf->GetReceiver()) {
    cout << "cannot initiate contract based on proposition";
    return;
  }

  int targetNode = (nodeOwnerIdentifier == cf->GetSender()) ? cf->GetReceiver(): cf->GetSender();
  // check if contract exists
  Contract* c = GetExistingContract(targetNode);
  if (c->idn == -1) {
    c = new Contract(strategy->contractIndex, targetNode, cf->targetPath);
    (strategy->activeContracts).push_back(c);
    strategy->contractIndex = strategy->contractIndex + 1;
  }
}

// TODO untested
void NodeProcessingUnit::BreakContract(CommFlare* cf) {
  if (nodeOwnerIdentifier != cf->GetSender() && nodeOwnerIdentifier != cf->GetReceiver()) {
    cout << "cannot initiate contract based on proposition";
    return;
  }

  int targetNode = (nodeOwnerIdentifier == cf->GetSender()) ? cf->GetReceiver(): cf->GetSender();
  strategy->DeleteActiveContract(targetNode);
}

/*
uses Timestamp.GetString() 

example queries: 
- eventType 
- eventType, role 
- summary (see timestamp_unit.hpp/otherDetails2) 

*/ 
map<string,map<string,int>> NodeProcessingUnit::ProcessTimestampSequenceForFrequency(vector<TimestampUnit*> tsuseq, vector<string> dataTemplate) {

  vector<map<string,string>> data = GatherTimestampData(tsuseq, dataTemplate);

  // example output 
  map<string,map<string,int>> output = map<string,map<string,int>>(); 
  for (auto x: data) {
      for (auto y: dataTemplate) {   
        if (x[y] != "") {
          output[y][x[y]]++;
        }   
      }
  }

  return output; 
} 

vector<map<string,string>> NodeProcessingUnit::GatherTimestampData(vector<TimestampUnit*> tsuseq, vector<string> dataTemplate) {

  vector<map<string,string>> output = vector<map<string,string>>(); 

  for (auto x: tsuseq) {
    output.push_back(x->ToStringMap(dataTemplate)); 
  }

  return output; 

}
